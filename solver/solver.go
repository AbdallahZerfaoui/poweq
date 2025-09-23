package solver

import (
	"log"
)

// Equation to solve: x^n = K * m^x
// equivalently: f(x) = n * ln(x) - ln(K) - x * ln(m) = 0
// f'(x) = n/x - ln(m)

func Solve(job Job, method string, logger *log.Logger) []Result {
	var solutions []Result

	// Handle edge cases first
	if done, solution := handleEdgeCases(job); done {
		if solution != -1.0 {
			return append(solutions, Result{Id: job.Id, X: solution, Steps: 0, Err: nil})
		}
	}

	switch method {
	case "newton":
		for _, x0 := range GetInitValues(job) {
			result := NewtonSolve(job, x0)
			if result.Err != nil {
				logger.Println("Error:", result.Err)
				solutions = append(solutions, Result{Id: job.Id, X: -1.0, Steps: 0, Err: result.Err})
			} else {
				// fmt.Printf("Found solution x = %.6f in %d steps\n", result.X, result.Steps)
				solutions = append(solutions, result)
			}
		}
	case "bisection":
		// For bisection, we need an interval [lower, upper]
		// Here, we use x0 as the midpoint and create a small interval around it
		for _, interval := range getIntervals(job) {
			result := BisectionSolve(job, interval[0], interval[1])
			if result.Err != nil {
				logger.Println("Error:", result.Err)
				solutions = append(solutions, Result{Id: job.Id, X: -1.0, Steps: 0, Err: result.Err})
			} else {
				// fmt.Printf("Found solution x = %.6f in %d steps\n", result.X, result.Steps)
				solutions = append(solutions, result)
			}
		}
	case "auto":
		// First try Newton-Raphson with multiple initial guesses
		for _, x0 := range GetInitValues(job) {
			result := NewtonSolve(job, x0)
			if result.Err == nil {
				solutions = append(solutions, result)
			}
		}
		// If no solutions found, fall back to Bisection method
		if len(solutions) == 0 {
			for _, interval := range getIntervals(job) {
				result := BisectionSolve(job, interval[0], interval[1])
				if result.Err == nil {
					solutions = append(solutions, result)
				}
			}
		}
	default:
		logger.Println("Unknown method:", method)
	}

	return solutions
}



