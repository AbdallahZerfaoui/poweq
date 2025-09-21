package solver

type Job struct {
	Id      int
	N, M, K float64
	A, B    float64
	Tol     float64
	MaxIter int
}

type Result struct {
	Id    int
	X     float64
	Steps int
	Err   error
}

type Batch struct {
	InFile  string
	OutFile string
	Jobs    []Job
	Results []Result
}
