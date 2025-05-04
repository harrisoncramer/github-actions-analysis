package collect

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type CollectParams struct {
	GithubRepo string
	MaxWorkers int
	MaxPages   int
	Outfile    string
	PerPage    int
}

// Collect fetches workflow runs and jobs from GitHub and writes them to a CSV file.
func Collect(params CollectParams) {
	client := NewGitHubClient(GithubClientParams{
		repo:       params.GithubRepo,
		maxWorkers: params.MaxWorkers,
		perPage:    params.PerPage,
	})

	outFile, err := os.Create(filepath.Join("data", params.Outfile))
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	// Do NOT defer writer.Flush() here anymore

	// Write headers before goroutine starts
	writer.Write([]string{
		"run_id", "workflow_name", "job_name",
		"status", "conclusion", "started_at", "completed_at", "duration_seconds",
	})

	jobChan := make(chan JobRecord, 1000)
	var wg sync.WaitGroup
	sem := make(chan struct{}, params.MaxWorkers)

	// Writer waitgroup
	var writerWg sync.WaitGroup
	writerWg.Add(1)
	go func() {
		defer writerWg.Done()
		for record := range jobChan {
			duration := int(record.job.CompletedAt.Sub(record.job.StartedAt).Seconds())
			writer.Write([]string{
				strconv.FormatInt(record.runID, 10),
				record.workflowName,
				record.job.Name,
				record.job.Status,
				record.job.Conclusion,
				record.job.StartedAt.Format(time.RFC3339),
				record.job.CompletedAt.Format(time.RFC3339),
				strconv.Itoa(duration),
			})
		}
	}()

	for page := 1; page <= params.MaxPages; page++ {
		runs, err := client.FetchWorkflowRuns(page)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching workflow runs: %v\n", err)
			continue
		}
		if len(runs) == 0 {
			break
		}

		for _, run := range runs {
			wg.Add(1)
			sem <- struct{}{}

			go func(run WorkflowRun) {
				defer wg.Done()
				defer func() {
					<-sem
				}()
				records, err := client.FetchJobsForRun(run)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error fetching jobs for run %d: %v\n", run.ID, err)
					return
				}
				for _, record := range records {
					jobChan <- record
				}
			}(run)
		}
	}

	go func() {
		wg.Wait()
		close(jobChan)
	}()

	writerWg.Wait()
	writer.Flush()
}
