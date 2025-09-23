package main

import (
	"fmt"
	"errors"
	"flag"
	"github.com/AbdallahZerfaoui/poweq/solver"
	"os"
	"encoding/csv"
	"math/rand"
)

const (
	DEFAULT_SOLUTIONS_ALGO = "auto"
	DEFAULT_ERROR_SOLUTION = -1.0
)

// logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

func solveCommand(args []string) ([]solver.Result, error) {
	solverFlagSet := flag.NewFlagSet("poweq", flag.ExitOnError)

	// Create flag set for the "solve" command
	n := solverFlagSet.Float64("n", 1.0, "The exponent n in the equation x^n = K m^x")
	m := solverFlagSet.Float64("m", 2.718281828, "The base m in the equation x^n = K m^x")
	K := solverFlagSet.Float64("K", 1.0, "The coefficient K in the equation x^n = K m^x")
	a := solverFlagSet.Float64("a", 1e-6, "Lowwer bound of the interval to search for a solution")
	b := solverFlagSet.Float64("b", 1e6, "Upper bound of the interval to search for a solution")
	tolence := solverFlagSet.Float64("tol", 1e-6, "Tolerance for the solution")
	maxIter := solverFlagSet.Int("maxIter", 100, "Maximum number of iterations")
	algorithm := solverFlagSet.String("alg", "auto", "Algorithm to use: 'newton' or 'bisection'")

	// Parse flags and execute solving logic
	err := solverFlagSet.Parse(args)
	if err != nil {
		logger.Println("Error parsing flags", "error", err)
		return nil, err
	}

	newJob := solver.Job{Id: 0, N: *n, M: *m, K: *K,
		A: *a, B: *b,
		Tol: *tolence, MaxIter: *maxIter}

	if err := solver.ValidateJob(newJob); err != nil {
		logger.Println("Invalid job parameters", "error", err)
		return nil, err
	}
	// Use the right solver functions from the solver package

	logger.Println("Solving the equation", "equation", fmt.Sprintf("x^%.2f = %.2f * %.2f^x", *n, *K, *m))
	logger.Println("Searching for a solution", "interval", fmt.Sprintf("[%.2f, %.2f]", *a, *b), "tolerance", *tolence, "max iterations", *maxIter)

	if !solver.SolutionsExist(newJob) {
		logger.Println("No solutions exist for the given parameters")
		return nil, errors.New("no solutions exist for the given parameters")
	}

	solutions := solver.Solve(newJob, *algorithm, logger)

	return solutions, nil
}

func displaySolutions(solutions []solver.Result) {
	for _, result := range solutions {
		if result.Err != nil {
			logger.Println("Error", "error", result.Err)
		} else {
			logger.Println("Found solution", "x", result.X, "steps", result.Steps)
		}
	}
}

func scanCommand(args []string) (solver.Batch, error) {
	scannerFlagSet := flag.NewFlagSet("scan", flag.ExitOnError)

	// Create flag set for the "scan" command
	in := scannerFlagSet.String("in", "jobs.csv", "Input file containing jobs to solve")
	out := scannerFlagSet.String("out", "solutions.csv", "Output file to write solutions")

	// Parse flags
	err := scannerFlagSet.Parse(args)
	if err != nil {
		logger.Println("Error parsing flags", "error", err)
		return solver.Batch{}, err
	}

	// Create a Batch instance
	batch := solver.Batch{InFile: *in, OutFile: *out,
		Jobs: []solver.Job{}, Results: []solver.Result{}}

	// Open files
	inFile, err := os.Open(*in)
	if err != nil {
		logger.Println("Error opening input file", "error", err)
		return solver.Batch{}, err
	}
	defer inFile.Close()

	outFile, err := os.Create(*out)
	if err != nil {
		logger.Println("Error creating output file", "error", err)
		return solver.Batch{}, err
	}
	defer outFile.Close()

	// Read jobs from input file
	jobs, err := readJobsFromCSV(inFile)
	if err != nil {
		logger.Println("Error reading jobs from input file", "error", err)
		return solver.Batch{}, err
	}
	batch.Jobs = jobs
	// Build a map of jobs by their IDs for easy lookup
	jobsMap := buildJobsMap(jobs)

	// fmt.Println("[debug] Jobs loaded:", len(batch.Jobs))
	// Solve each job and collect results
	for _, job := range batch.Jobs {
		if err := solver.ValidateJob(job); err != nil {
			logger.Println("Invalid job parameters", "error", err)
			batch.Results = append(batch.Results, solver.Result{Id: job.Id, X: DEFAULT_ERROR_SOLUTION, Steps: 0, Err: err})
			continue
		}
		if !solver.SolutionsExist(job) {
			batch.Results = append(batch.Results, solver.Result{Id: job.Id, X: DEFAULT_ERROR_SOLUTION, Steps: 0, Err: errors.New("no solutions exist for the given parameters")})
			continue
		}
		solutions := solver.Solve(job, DEFAULT_SOLUTIONS_ALGO, logger) // or "bisection"
		if len(solutions) > 0 {
			batch.Results = append(batch.Results, solutions...)
		} else {
			batch.Results = append(batch.Results, solver.Result{Id: job.Id, X: DEFAULT_ERROR_SOLUTION, Steps: 0, Err: errors.New("no solutions found")})
		}
	}

	// Write results to output file using the helper function
	batch, err = writeResultsToCSV(outFile, batch, jobsMap)
	if err != nil {
		logger.Println("Error writing results to output file", "error", err)
		return batch, err
	}
	return batch, nil
}

func generateCommand(args []string) error {
    generateFlagSet := flag.NewFlagSet("generate", flag.ExitOnError)

    N := generateFlagSet.Int("N", 50, "Number of jobs to generate")
    out := generateFlagSet.String("out", "jobs.csv", "Output file to write jobs")

    err := generateFlagSet.Parse(args)
    if err != nil {
        return err
    }

    outFile, err := os.Create(*out)
    if err != nil {
        return err
    }
    defer outFile.Close()

    writer := csv.NewWriter(outFile)
    defer writer.Flush()

    // Write header
    writer.Write([]string{"Id", "N", "M", "K", "A", "B", "Tol", "MaxIter"})

    for i := 1; i <= *N; i++ {
        job := solver.Job{
            Id:      i,
            N:       float64(rand.Intn(1000) + 10) / 100.0,       // n between 0.1 and 10.0
            M:       float64(rand.Intn(500) + 110) / 100.0,        // m between 1.1 and 6.0
            K:       float64(rand.Intn(1000000) + 1) / 100.0,      // K between 0.01 and 10000
            A:       1e-6,
            B:       float64(rand.Intn(1000000) + 10), // b between 0.1 and 1000010
            Tol:     1e-6,
            MaxIter: rand.Intn(91) + 10,               // MaxIter between 10 and 100
        }
		if solver.SolutionsExist(job) { // TODO: Should i keep this check?
			record := []string{
				fmt.Sprintf("%d", job.Id),
				fmt.Sprintf("%.2f", job.N),
				fmt.Sprintf("%.2f", job.M),
				fmt.Sprintf("%.2f", job.K),
				fmt.Sprintf("%.2f", job.A),
				fmt.Sprintf("%.2f", job.B),
				fmt.Sprintf("%.2e", job.Tol),
				fmt.Sprintf("%d", job.MaxIter),
			}
			writer.Write(record)
		} else {
			i-- // Retry this iteration
		}
	}

    fmt.Printf("Generated %d jobs into %s\n", *N, *out)
    return nil
}
