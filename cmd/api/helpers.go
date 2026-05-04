package main

import (
	"errors"

	"github.com/AbdallahZerfaoui/poweq/solver"

)

func (req SolveRequest) Solve4API() (SolveResponse, error) {
	var resp SolveResponse

	// Call the solver function
	job := solver.Job{
		Id:      0, // Or random ID for the job
		N:       req.N,
		M:       req.M,
		K:       req.K,
		A:       req.A,
		B:       req.B,
		Tol:     req.Tolerance,
		MaxIter: req.MaxIter,
	}
	solutions := job.Solve(req.Algorithm, logger)
	if len(solutions) == 0 {
		return resp, errors.New("no solutions found")
	}

	resp.Solutions = make([]APISolution, len(solutions))
	for i, sol := range solutions {
		resp.Solutions[i] = APISolution{
			X:     sol.X,
			Steps: sol.Steps,
			Error: sol.Err,
		}
	}

	return resp, nil
}
