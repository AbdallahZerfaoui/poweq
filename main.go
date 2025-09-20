package main

import (
	"flag"
	"fmt"
	"github.com/AbdallahZerfaoui/poweq/solver"
	"os"
)

func main() {
	optsFlagSet := flag.NewFlagSet("poweq", flag.ExitOnError)

	// Flags: it returns a pointer to the variable
	n := optsFlagSet.Float64("n", 1.0, "The exponent n in the equation x^n = K m^x")
	m := optsFlagSet.Float64("m", 1.0, "The base m in the equation x^n = K m^x")
	K := optsFlagSet.Float64("K", 1.0, "The coefficient K in the equation x^n = K m^x")
	a := optsFlagSet.Float64("a", 1e-6, "Lowwer bound of the interval to search for a solution")
	b := optsFlagSet.Float64("b", 1e6, "Upper bound of the interval to search for a solution")
	tolence := optsFlagSet.Float64("tol", 1e-6, "Tolerance for the solution")
	maxIter := optsFlagSet.Int("maxIter", 100, "Maximum number of iterations")

	// Before this step, n, m and K are default values
	err := optsFlagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return
	}

	newJob := solver.Job{N: *n, M: *m, K: *K,
		A: *a, B: *b,
		Tol: *tolence, MaxIter: *maxIter}

	if err := solver.ValidateJob(newJob); err != nil {
		fmt.Println("Invalid job parameters:", err)
		return
	}

	fmt.Printf("Solving the equation x^%.2f = %.2f * %.2f^x\n", *n, *K, *m)
	fmt.Printf("Searching for a solution in the interval [%.2f, %.2f] with tolerance %.2e and max iterations %d\n", *a, *b, *tolence, *maxIter)

	for _, x0 := range solver.GetInitValues(newJob) {
		result := solver.Solve(newJob, x0)
		if result.Err != nil {
			fmt.Println("Error:", result.Err)
		} else {
			fmt.Printf("Found solution x = %.6f in %d steps\n", result.X, result.Steps)
		}
	}
}
