package solver

import (
	"errors"
	"math"
)

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
