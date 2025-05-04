package collect

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type CollectParams struct {
	GithubRepo string
	MaxWorkers int
	MaxPages   int
	Outfile    string
}

// Collect fetches workflow runs and jobs from GitHub and writes them to a CSV file.
func Collect(params CollectParams) {
	client := NewGitHubClient(GithubClientParams{
		repo: params.GithubRepo,
	})

	outFile, err := os.Create(params.Outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	writer.Write([]string{
		"run_id", "workflow_name", "job_name",
		"status", "conclusion", "started_at", "completed_at", "duration_seconds",
	})

	jobChan := make(chan JobRecord, 1000)
	var wg sync.WaitGroup
	sem := make(chan struct{}, params.MaxWorkers)

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

	for record := range jobChan {
		duration := int(record.Job.CompletedAt.Sub(record.Job.StartedAt).Seconds())
		writer.Write([]string{
			strconv.FormatInt(record.RunID, 10),
			record.WorkflowName,
			record.Job.Name,
			record.Job.Status,
			record.Job.Conclusion,
			record.Job.StartedAt.Format(time.RFC3339),
			record.Job.CompletedAt.Format(time.RFC3339),
			strconv.Itoa(duration),
		})
	}
}
