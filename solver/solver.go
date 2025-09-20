package solver

import (
	"errors"
	"math"
	// "fmt"
)

// Equation to solve: x^n = K * m^x
// equivalently: f(x) = n * ln(x) - ln(K) - x * ln(m) = 0
// f'(x) = n/x - ln(m)

func f(x, n, m, K float64) float64 {
	return n * math.Log(x) - math.Log(K) - x * math.Log(m)
}

func fPrime(x, n, m float64) float64 {
	return n/x - math.Log(m)
}

// Newton-Raphson method
func Solve(job Job, x0 float64) Result {
	a, b := job.A, job.B
	n, m := job.N, job.M
	K, tol, maxIter := job.K, job.Tol, job.MaxIter

	for i := range maxIter {
		fx := f(x0, n, m, K)
		fpx := fPrime(x0, n, m)
		
		if fpx == 0 {
			return Result{0, i, errors.New("derivative is zero")}
		}

		x1 := x0 - fx/fpx // Newton-Raphson update

		if math.Abs(x1 - x0) < tol {
			return Result{x1, i + 1, nil}
		}
		// fmt.Println("f(x):", fx, "f'(x):", fpx)
		// fmt.Println("Current approximation:", x1)
		if x1 < a || x1 > b {
			return Result{0, i + 1, errors.New("solution out of bounds")}
		}

		x0 = x1
	}

	return Result{0, maxIter, errors.New("maximum iterations reached without convergence")}
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
	return []float64{a + x_limit / 10, (x_limit + b) / 2}
}

