package solver

import (
	"math"
)

func handleEdgeCases(job Job) (bool, float64) {
	n, m, K := job.N, job.M, job.K
	a, b := job.A, job.B

	// Case m = 1
	if m == 1 {
		solution := math.Pow(K, 1/n)
		if solution < a || solution > b {
			solution = -1.0
		}
		return true, solution
	}

	// Case n = 0
	if n == 0 {
		solution := -1.0 * math.Log(K) / math.Log(m)
		if solution < a || solution > b {
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
