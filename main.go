package main

import (
	"flag"
	"fmt"
	"github.com/AbdallahZerfaoui/poweq/solver"
	"os"
	"time"
	"runtime"
	"errors"
)

func main() {
	var mStart, mEnd runtime.MemStats
	runtime.ReadMemStats(&mStart)

	start := time.Now()
	// solverFlagSet := flag.NewFlagSet("poweq", flag.ExitOnError)
	scannerFlagSet := flag.NewFlagSet("scan", flag.ExitOnError)

	// Solver Flags: it returns a pointer to the variable
	// n := solverFlagSet.Float64("n", 1.0, "The exponent n in the equation x^n = K m^x")
	// m := solverFlagSet.Float64("m", 1.0, "The base m in the equation x^n = K m^x")
	// K := solverFlagSet.Float64("K", 1.0, "The coefficient K in the equation x^n = K m^x")
	// a := solverFlagSet.Float64("a", 1e-6, "Lowwer bound of the interval to search for a solution")
	// b := solverFlagSet.Float64("b", 1e6, "Upper bound of the interval to search for a solution")
	// tolence := solverFlagSet.Float64("tol", 1e-6, "Tolerance for the solution")
	// maxIter := solverFlagSet.Int("maxIter", 100, "Maximum number of iterations")
	// algorithm := solverFlagSet.String("alg", "newton", "Algorithm to use: 'newton' or 'bisection'")

	// Scanner Flags
	// inFile := scannerFlagSet.String("in", "jobs.csv", "Input file containing jobs to solve")
	// outFile := scannerFlagSet.String("out", "solutions.csv", "Output file to write solutions")
	
	// Before this step, n, m and K are default values
	if os.Args[1] == "scan" {
		err := scannerFlagSet.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing flags:", err)
			return
		}
		// newBatch := solver.Batch{inFile: *inFile, outFile: *outFile, Jobs: nil, Results: nil}
	} else if os.Args[1] == "solve" {
		// err := solverFlagSet.Parse(os.Args[2:])
		// if err != nil {
		// 	fmt.Println("Error parsing flags:", err)
		// 	return
		// }
		solutions, err := solveCommand(os.Args[2:])
		if err != nil {
			fmt.Println("Error solving:", err)
			return
		}
		displaySolutions(solutions)
	}

	// newJob := solver.Job{N: *n, M: *m, K: *K,
	// 	A: *a, B: *b,
	// 	Tol: *tolence, MaxIter: *maxIter}

	// if err := solver.ValidateJob(newJob); err != nil {
	// 	fmt.Println("Invalid job parameters:", err)
	// 	return
	// }

	// fmt.Printf("Solving the equation x^%.2f = %.2f * %.2f^x\n", *n, *K, *m)
	// fmt.Printf("Searching for a solution in the interval [%.2f, %.2f] with tolerance %.2e and max iterations %d\n", *a, *b, *tolence, *maxIter)

	// if !solver.SolutionsExist(newJob) {
	// 	fmt.Println("No solutions exist for the given parameters.")
	// 	return
	// }

	// solutions := solver.Solve(newJob, *algorithm)

	// for _, result := range solutions {
	// 	if result.Err != nil {
	// 		fmt.Println("Error:", result.Err)
	// 	} else {
	// 		fmt.Printf("Found solution x = %.6f in %d steps\n", result.X, result.Steps)
	// 	}
	// }

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

	newJob := solver.Job{N: *n, M: *m, K: *K,
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