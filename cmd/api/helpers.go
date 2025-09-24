package main

import (
	"errors"

	"github.com/AbdallahZerfaoui/poweq/solver"

)

func Solve4API(req SolveRequest) (SolveResponse, error) {
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
	solutions := solver.Solve(job, req.Algorithm, logger)
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
