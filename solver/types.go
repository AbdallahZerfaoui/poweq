package solver

type Job struct {
	N, M, K float64
	A, B    float64
	Tol     float64
	MaxIter int
}

type Result struct {
	X     float64
	Steps int
	Err   error
}