package main

import (
	"fmt"
	"errors"
	"flag"
	"github.com/AbdallahZerfaoui/poweq/solver"
	"os"
)

const (
	DEFAULT_SOLUTIONS_ALGO = "newton"
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
	algorithm := solverFlagSet.String("alg", "newton", "Algorithm to use: 'newton' or 'bisection'")

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
