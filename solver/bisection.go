package solver

import (
	"errors"
	"math"
)

// Bisection method

func BisectionSolve(job Job, lower float64, upper float64) Result {
	// a, b := job.A, job.B
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

		if math.Abs(fc) < tol || (upper-lower)/2 < tol {
			return Result{job.Id, c, i + 1, nil}
		}

		if fa*fc < 0 {
			upper = c
		} else {
			lower = c
		}
		fa = f(lower, n, m, K) // Update fa for the new interval
	}

	return Result{job.Id, 0, 0, errors.New("maximum iterations reached without convergence")}
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
