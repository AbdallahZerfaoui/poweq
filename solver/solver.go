package solver

import (
	"errors"
	"fmt"
	"math"
)

// Equation to solve: x^n = K * m^x
// equivalently: f(x) = n * ln(x) - ln(K) - x * ln(m) = 0
// f'(x) = n/x - ln(m)

func Solve(job Job, method string) []Result {
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
				fmt.Println("Error:", result.Err)
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
				fmt.Println("Error:", result.Err)
				solutions = append(solutions, Result{Id: job.Id, X: -1.0, Steps: 0, Err: result.Err})
			} else {
				// fmt.Printf("Found solution x = %.6f in %d steps\n", result.X, result.Steps)
				solutions = append(solutions, result)
			}
		}
	default:
		fmt.Println("Unknown method:", method)
	}

	return solutions
}

func f(x, n, m, K float64) float64 {
	return n*math.Log(x) - math.Log(K) - x*math.Log(m)
}

func fPrime(x, n, m float64) float64 {
	return n/x - math.Log(m)
}

// Newton-Raphson method
func NewtonSolve(job Job, x0 float64) Result {
	a, b := job.A, job.B
	n, m := job.N, job.M
	K, tol, maxIter := job.K, job.Tol, job.MaxIter

	for i := range maxIter {
		fx := f(x0, n, m, K)
		fpx := fPrime(x0, n, m)

		if fpx == 0 {
			return Result{job.Id, 0, i, errors.New("derivative is zero")}
		}

		x1 := x0 - fx/fpx // Newton-Raphson update

		if math.Abs(x1-x0) < tol {
			// We check only the last value to see if it's within bounds
			// because if the initial guess is within bounds and the method converges,
			// it should remain within bounds.
			// However, if an intermediate value is out of bounds, we might still converge to a valid solution.
			// Thus, we only check the final result.
			if x1 < a || x1 > b {
				return Result{job.Id, 0, i + 1, errors.New("solution out of bounds")}
			}
			return Result{job.Id, x1, i + 1, nil}
		}

		x0 = x1
	}

	return Result{job.Id, 0, maxIter, errors.New("maximum iterations reached without convergence")}
}

func BisectionSolve(job Job, lower float64, upper float64) Result {
	a, b := job.A, job.B
	n, m := job.N, job.M
	K, tol, maxIter := job.K, job.Tol, job.MaxIter

	fa := f(lower, n, m, K)
	fb := f(upper, n, m, K)

	if fa*fb > 0 {
		return Result{job.Id, 0, 0, errors.New("f(a) and f(b) must have opposite signs")}
	}

	for i := range maxIter {
		c := (lower + upper) / 2
		fc := f(c, n, m, K)

		if math.Abs(fc) < tol || (b-a)/2 < tol {
			return Result{job.Id, c, i + 1, nil}
		}

		if fa*fc < 0 {
			upper = c
		} else {
			lower = c
		}
	}

	return Result{job.Id, 0, 0, errors.New("maximum iterations reached without convergence")}
}

func handleEdgeCases(job Job) (bool, float64) {
	n, m, K := job.N, job.M, job.K
	a, b := job.A, job.B

	// Case m = 1
	if m == 1 {
		solution := math.Pow(K, 1/n)
		if solution < a && solution > b {
			solution = -1.0
		}
		return true, solution
	}

	// Case n = 0
	if n == 0 {
		solution := -1.0 * math.Log(K) / math.Log(m)
		if solution < a && solution > b {
			solution = -1.0
		}
		return true, solution
	}

	return false, -1.0
}

func GetInitValues(job Job) []float64 {
	a, b := job.A, job.B
	n, m := job.N, job.M

	// x_limit is the point where f'(x) = 0
	// f'(x) = n/x - ln(m) = 0  =>  x = n / ln(m)
	x_limit := n / math.Log(m)
	if x_limit >= b {
		return []float64{(a + b) / 2}
	}
	// i divide by 10 to avoid starting too close to the null point of the derivative
	// which can cause very large steps and divergence
	// similarly, i take the midpoint between x_limit and b to avoid being too close to the null point
	// this is a heuristic choice to improve convergence chances
	return []float64{a + x_limit/10, (x_limit + b) / 2}
}

func getIntervals(job Job) [][2]float64 {
	a, b := job.A, job.B
	n, m := job.N, job.M

	// x_limit is the point where f'(x) = 0
	// f'(x) = n/x - ln(m) = 0  =>  x = n / ln(m)
	x_limit := n / math.Log(m)
	if x_limit >= b {
		return [][2]float64{{a, b}}
	}
	return [][2]float64{{a, x_limit}, {x_limit, b}}
}
