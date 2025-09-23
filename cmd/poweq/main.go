package main

import (
	"os"
	"runtime"
	"time"
	"log"
)

const(
	KILO = 1024
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	var mStart, mEnd runtime.MemStats
	// logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	runtime.ReadMemStats(&mStart)

	start := time.Now()
	// CRASH if no arguments!
	if len(os.Args) < 2 {
		logger.Println("no command provided", "usage", "Available commands: solve, scan")
		return
	}
	// Before this step, n, m and K are default values
	switch os.Args[1] {
	case "scan":
		_, err := scanCommand(os.Args[2:])
		if err != nil {
			logger.Println("scan failed", "error", err)
			return
		}

	case "solve":
		solutions, err := solveCommand(os.Args[2:])
		if err != nil {
			logger.Println("solve failed", "error", err)
			return
		}
		displaySolutions(solutions)

	case "generate":
		err := generateCommand(os.Args[2:])
		if err != nil {
			logger.Println("generate failed", "error", err)
			return
		}

	default:
		logger.Println("unknown command", "command", os.Args[1])
		logger.Println("Available commands: solve, scan")
		return
	}

	elapsed := time.Since(start)
	runtime.ReadMemStats(&mEnd)
	usedMemory := (mEnd.Alloc - mStart.Alloc) / KILO
	logger.Println("memory usage", "used", usedMemory)
	logger.Println("execution time", "duration", elapsed)
}


