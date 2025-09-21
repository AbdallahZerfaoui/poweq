package solver

import (
	"math"
)

func f(x, n, m, K float64) float64 {
	return n*math.Log(x) - math.Log(K) - x*math.Log(m)
}

func fPrime(x, n, m float64) float64 {
	return n/x - math.Log(m)
}