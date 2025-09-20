package main

import (
	"fmt"
	"os"
	"flag"
)

func main() {
	optsFlagSet := flag.NewFlagSet("poweq", flag.ExitOnError)
	
	// Flags
	n := optsFlagSet.Float64("n", 1.0, "The exponent n in the equation x^n = A m^x")
	m := optsFlagSet.Float64("m", 1.0, "The base m in the equation x^n = A m^x")
	A := optsFlagSet.Float64("A", 1.0, "The coefficient A in the equation x^n = A m^x")

	// Before this step, n, m and A are default values
	err := optsFlagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return
	}

	fmt.Printf("Solving the equation x^%.2f = %.2f * %.2f^x\n", *n, *A, *m)
}