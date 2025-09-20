package solver

import (
	"errors"
)

func ValidateJob(job Job) error {
	if job.N <= 0 {
		return errors.New("n must be positive")
	}
	if job.M <= 0 {
		return errors.New("m must be positive")
	}
	if job.K <= 0 {
		return errors.New("K must be positive")
	}
	if job.A <= 0 || job.B <= 0 {
		return errors.New("A and B must be positive")
	}
	if job.A >= job.B {
		return errors.New("A must be less than B")
	}
	if job.Tol <= 0 {
		return errors.New("tolerance must be positive")
	}
	if job.MaxIter <= 0 {
		return errors.New("maxIter must be positive")
	}
	return nil
}