package main

import (
	"flag"
	"fmt"
	"github.com/AbdallahZerfaoui/poweq/solver"
	"os"
	"time"
	"runtime"
	"errors"
	"encoding/csv"
	"strings"
)

func main() {
	var mStart, mEnd runtime.MemStats
	runtime.ReadMemStats(&mStart)

	start := time.Now()

	// scannerFlagSet := flag.NewFlagSet("scan", flag.ExitOnError)


	// Scanner Flags
	// inFile := scannerFlagSet.String("in", "jobs.csv", "Input file containing jobs to solve")
	// outFile := scannerFlagSet.String("out", "solutions.csv", "Output file to write solutions")
	
	// Before this step, n, m and K are default values
	switch os.Args[1] {
	case "scan":
		// err := scannerFlagSet.Parse(os.Args[2:])
		// if err != nil {
		// 	fmt.Println("Error parsing flags:", err)
		// 	return
		// }
		_, err := scanCommand(os.Args[2:])
		if err != nil {
			fmt.Println("Error scanning:", err)
			return
		}

	case "solve":

		solutions, err := solveCommand(os.Args[2:])
		if err != nil {
			fmt.Println("Error solving:", err)
			return
		}
		displaySolutions(solutions)

	default:
		fmt.Println("Unknown command:", os.Args[1])
		fmt.Println("Available commands: solve, scan")
		return
	}

	elapsed := time.Since(start)
	runtime.ReadMemStats(&mEnd)
	usedMemory := (mEnd.Alloc - mStart.Alloc) / 1024
	fmt.Printf("Used memory: %d KB\n", usedMemory)
	fmt.Printf("Execution time: %s\n", elapsed)
}

func solveCommand(args []string) ([]solver.Result, error) {
	solverFlagSet := flag.NewFlagSet("poweq", flag.ExitOnError)

	// Create flag set for the "solve" command
	n := solverFlagSet.Float64("n", 1.0, "The exponent n in the equation x^n = K m^x")
	m := solverFlagSet.Float64("m", 1.0, "The base m in the equation x^n = K m^x")
	K := solverFlagSet.Float64("K", 1.0, "The coefficient K in the equation x^n = K m^x")
	a := solverFlagSet.Float64("a", 1e-6, "Lowwer bound of the interval to search for a solution")
	b := solverFlagSet.Float64("b", 1e6, "Upper bound of the interval to search for a solution")
	tolence := solverFlagSet.Float64("tol", 1e-6, "Tolerance for the solution")
	maxIter := solverFlagSet.Int("maxIter", 100, "Maximum number of iterations")
	algorithm := solverFlagSet.String("alg", "newton", "Algorithm to use: 'newton' or 'bisection'")

	// Parse flags and execute solving logic
	err := solverFlagSet.Parse(args)
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return nil, err
	}

	newJob := solver.Job{Id: 0, N: *n, M: *m, K: *K,
		A: *a, B: *b,
		Tol: *tolence, MaxIter: *maxIter}

	if err := solver.ValidateJob(newJob); err != nil {
		fmt.Println("Invalid job parameters:", err)
		return nil, err
	}
	// Use the right solver functions from the solver package

	fmt.Printf("Solving the equation x^%.2f = %.2f * %.2f^x\n", *n, *K, *m)
	fmt.Printf("Searching for a solution in the interval [%.2f, %.2f] with tolerance %.2e and max iterations %d\n", *a, *b, *tolence, *maxIter)

	if !solver.SolutionsExist(newJob) {
		fmt.Println("No solutions exist for the given parameters.")
		return nil, errors.New("no solutions exist for the given parameters")
	}

	solutions := solver.Solve(newJob, *algorithm)

	return solutions, nil
}

func displaySolutions(solutions []solver.Result) {
	for _, result := range solutions {
		if result.Err != nil {
			fmt.Println("Error:", result.Err)
		} else {
			fmt.Printf("Found solution x = %.6f in %d steps\n", result.X, result.Steps)
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
		fmt.Println("Error parsing flags:", err)
		return solver.Batch{}, err
	}

	// Create a Batch instance
	batch := solver.Batch{InFile: *in, OutFile: *out,
		Jobs: []solver.Job{}, Results: []solver.Result{}}

	// Open files
	inFile, err := os.Open(*in)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return solver.Batch{}, err
	}
	defer inFile.Close()

	outFile, err := os.Create(*out)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return solver.Batch{}, err
	}
	defer outFile.Close()

	// Read jobs from input file
	jobs, err := readJobsFromCSV(inFile)
	if err != nil {
		fmt.Println("Error reading jobs from input file:", err)
		return solver.Batch{}, err
	}
	batch.Jobs = jobs
	// Build a map of jobs by their IDs for easy lookup
	jobsMap := buildJobsMap(jobs)

	fmt.Println("[debug] Jobs loaded:", len(batch.Jobs))
	// Solve each job and collect results
	for _, job := range batch.Jobs {
		if err := solver.ValidateJob(job); err != nil {
			fmt.Println("Invalid job parameters:", err)
			batch.Results = append(batch.Results, solver.Result{Id: job.Id, X: -1.0, Steps: 0, Err: err})
			continue
		}
		if !solver.SolutionsExist(job) {
			fmt.Printf("[debug] Solving job %d: N=%.2f, M=%.2f, K=%.2f, A=%.2f, B=%.2f, Tol=%.2e, MaxIter=%d\n", job.Id, job.N, job.M, job.K, job.A, job.B, job.Tol, job.MaxIter)
			fmt.Println("[scan] No solutions exist for the given parameters.")
			batch.Results = append(batch.Results, solver.Result{Id: job.Id, X: -1.0, Steps: 0, Err: errors.New("no solutions exist for the given parameters")})
			continue
		}
		solutions := solver.Solve(job, "newton") // or "bisection"
		if len(solutions) > 0 {
			batch.Results = append(batch.Results, solutions...)
		} else {
			batch.Results = append(batch.Results, solver.Result{Id: job.Id, X: -1.0, Steps: 0, Err: errors.New("no solutions found")})
		}
	}

	// Write results to output file
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{"Id", "N", "M", "K", "A", "B", "Tol", "MaxIter", "X", "Steps", "Error"})
	if err != nil {
		fmt.Println("Error writing header:", err)
		return batch, err
	}

	// Write records
	for _, result := range batch.Results {
		job := jobsMap[result.Id]
		err = writer.Write([]string{
			fmt.Sprintf("%d", job.Id),
			fmt.Sprintf("%.2f", job.N),
			fmt.Sprintf("%.2f", job.M),
			fmt.Sprintf("%.2f", job.K),
			fmt.Sprintf("%.2f", job.A),
			fmt.Sprintf("%.2f", job.B),
			fmt.Sprintf("%.2e", job.Tol),
			fmt.Sprintf("%d", job.MaxIter),
			fmt.Sprintf("%.6f", result.X),
			fmt.Sprintf("%d", result.Steps),
			fmt.Sprintf("%v", result.Err),
		})
		if err != nil {
			fmt.Println("Error writing record:", err)
			return batch, err
		}
	}

	return batch, nil
}

func readJobsFromCSV(file *os.File) ([]solver.Job, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return nil, err
	}

	// Build jobs from records
	var jobs []solver.Job
	for _, record := range records[1:] { // Skip header
		// fmt.Println("[debug] record:", record)
		if len(record) != 8 {
			fmt.Println("Invalid record length:", record)
			continue
		}
		var job solver.Job

		_, err := fmt.Sscanf(strings.Join(record, ","), "%d,%f,%f,%f,%f,%f,%f,%d",
			&job.Id, &job.N, &job.M, &job.K, &job.A, &job.B, &job.Tol, &job.MaxIter)
		if err != nil {
			fmt.Println("Error parsing record:", record, err)
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func buildJobsMap(jobs []solver.Job) map[int]solver.Job {
	jobsMap := make(map[int]solver.Job)
	for _, job := range jobs {
		jobsMap[job.Id] = job
	}
	return jobsMap
}