package main

import (
	"encoding/csv"
	"fmt"
	"github.com/AbdallahZerfaoui/poweq/solver"
	"os"
	"strings"
)

func readJobsFromCSV(file *os.File) ([]solver.Job, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return nil, err
	}

	// Build jobs from records
	var jobs []solver.Job
	for _, record := range records[1:] { // Skip header
		// fmt.Println("[debug] record:", record)
		if len(record) != 8 {
			fmt.Println("Invalid record length:", record)
			continue
		}
		var job solver.Job

		_, err := fmt.Sscanf(strings.Join(record, ","), "%d,%f,%f,%f,%f,%f,%f,%d",
			&job.Id, &job.N, &job.M, &job.K, &job.A, &job.B, &job.Tol, &job.MaxIter)
		if err != nil {
			fmt.Println("Error parsing record:", record, err)
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func buildJobsMap(jobs []solver.Job) map[int]solver.Job {
	jobsMap := make(map[int]solver.Job)
	for _, job := range jobs {
		jobsMap[job.Id] = job
	}
	return jobsMap
}

func writeResultsToCSV(outFile *os.File, batch solver.Batch, jobsMap map[int]solver.Job) (solver.Batch, error) {
	// Write results to output file
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Write header
	err := writer.Write([]string{"Id", "N", "M", "K", "A", "B", "Tol", "MaxIter", "X", "Steps", "Error"})
	if err != nil {
		fmt.Println("Error writing header:", err)
		return batch, err
	}

	// Write records
	for _, result := range batch.Results {
		if result.Id == 0 {
			fmt.Println("Skipping result with no associated job ID")
			continue // Skip results with no associated job ID
		}
		job := jobsMap[result.Id]
		err = writer.Write([]string{
			fmt.Sprintf("%d", job.Id),
			fmt.Sprintf("%.2f", job.N),
			fmt.Sprintf("%.2f", job.M),
			fmt.Sprintf("%.2f", job.K),
			fmt.Sprintf("%.2f", job.A),
			fmt.Sprintf("%.2f", job.B),
			fmt.Sprintf("%.2e", job.Tol),
			fmt.Sprintf("%d", job.MaxIter),
			fmt.Sprintf("%.6f", result.X),
			fmt.Sprintf("%d", result.Steps),
			fmt.Sprintf("%v", result.Err),
		})
		if err != nil {
			fmt.Println("Error writing record:", err)
			return batch, err
		}
	}
	return batch, nil
}