package solver

import (
	"errors"
	"math"
)

func ValidateJob(job Job) error {
	if job.N < 0 {
		return errors.New("n must be positive or zero")
	}
	if job.M < 1 {
		return errors.New("m must be at least 1")
	}
	if job.K <= 0 {
		return errors.New("value K must be positive")
	}
	if job.A < 0 || job.B <= 0 {
		return errors.New("values A and B must be positive")
	}
	if job.A >= job.B {
		return errors.New("value A must be less than B")
	}
	if job.Tol <= 0 {
		return errors.New("tolerance must be positive")
	}
	if job.MaxIter <= 0 {
		return errors.New("maxIter must be positive")
	}
	return nil
}

func SolutionsExist(job Job) bool {
	n, m, K := job.N, job.M, job.K

	// If the highest point of f(x) is below 0, there is no solution
	// f'(x) = n/x - ln(m) = 0  =>  x = n / ln(m)
	// this is valid only if m > 1 because ln(m) must be positive
	x_limit := n / math.Log(m)
	if y0 := f(x_limit, n, m, K); y0 < 0 {
		return false
	}
	return true
}