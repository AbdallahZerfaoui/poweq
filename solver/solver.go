package solver

import (
	"errors"
	"math"
)

// Equation to solve: x^n = A * m^x
// equivalently: f(x) = n * ln(x) - ln(A) - x * ln(m) = 0
// f'(x) = n/x - ln(m)

func f(x, n, m, A float64) float64 {
	return n * math.Log(x) - math.Log(A) - x * math.Log(m)
}

func fPrime(x, n, m float64) float64 {
	return n/x - math.Log(m)
}

// Newton-Raphson method
func Solve(job Job) Result {
	a, b := job.A, job.B
	n, m := job.N, job.M
	tol, maxIter := job.Tol, job.MaxIter

	x0 := 1.0
	for i := 0; i < maxIter; i++ {
		fx := f(x0, n, m, a)
		fpx := fPrime(x0, n, m)
		
		if fpx == 0 {
			return Result{0, i, errors.New("derivative is zero")}
		}

		x1 := x0 - fx/fpx // Newton-Raphson update

		if math.Abs(x1 - x0) < tol {
			return Result{x1, i + 1, nil}
		}

		if x1 < a || x1 > b {
			return Result{0, i + 1, errors.New("solution out of bounds")}
		}

		x0 = x1
	}

	return Result{0, maxIter, errors.New("maximum iterations reached without convergence")}
}

