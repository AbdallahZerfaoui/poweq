package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

const(
	KILO = 1024
)

func main() {
	var mStart, mEnd runtime.MemStats
	runtime.ReadMemStats(&mStart)

	start := time.Now()
	// Before this step, n, m and K are default values
	switch os.Args[1] {
	case "scan":
		_, err := scanCommand(os.Args[2:])
		if err != nil {
			fmt.Println("Error scanning:", err)
			return
		}

	case "solve":
		solutions, err := solveCommand(os.Args[2:])
		if err != nil {
			fmt.Println("Error solving:", err)
			return
		}
		displaySolutions(solutions)

	default:
		fmt.Println("Unknown command:", os.Args[1])
		fmt.Println("Available commands: solve, scan")
		return
	}

	elapsed := time.Since(start)
	runtime.ReadMemStats(&mEnd)
	usedMemory := (mEnd.Alloc - mStart.Alloc) / KILO
	fmt.Printf("Used memory: %d KB\n", usedMemory)
	fmt.Printf("Execution time: %s\n", elapsed)
}


