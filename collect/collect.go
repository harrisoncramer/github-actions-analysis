package collect

import (
	"encoding/csv"
	"fmt"
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
func Collect(params CollectParams) error {
	client := NewGitHubClient(GithubClientParams{
		repo:       params.GithubRepo,
		maxWorkers: params.MaxWorkers,
		perPage:    params.PerPage,
	})

	outFile, err := os.Create(filepath.Join("data", params.Outfile))
	if err != nil {
		return fmt.Errorf("failed to open outfile: %w", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	err = writer.Write([]string{
		"run_id", "workflow_name", "job_name",
		"status", "conclusion", "started_at", "completed_at", "duration_seconds",
	})
	if err != nil {
		return fmt.Errorf("failed to write headers to outfile: %w", err)
	}

	jobChan := make(chan JobRecord, 1000)
	var wg sync.WaitGroup
	sem := make(chan struct{}, params.MaxWorkers)
	errChan := make(chan error, params.MaxWorkers*params.MaxPages)

	var writerWg sync.WaitGroup
	writerWg.Add(1)
	go func() {
		defer writerWg.Done()
		for record := range jobChan {
			duration := int(record.job.CompletedAt.Sub(record.job.StartedAt).Seconds())
			err := writer.Write([]string{
				strconv.FormatInt(record.runID, 10),
				record.workflowName,
				record.job.Name,
				record.job.Status,
				record.job.Conclusion,
				record.job.StartedAt.Format(time.RFC3339),
				record.job.CompletedAt.Format(time.RFC3339),
				strconv.Itoa(duration),
			})
			if err != nil {
				errChan <- fmt.Errorf("failed writing record to output file: %w", err)
			}
		}
	}()

	for page := 1; page <= params.MaxPages; page++ {
		runs, err := client.FetchWorkflowRuns(page)
		if err != nil {
			return fmt.Errorf("error fetching workflow runs: %v", err)
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
					errChan <- fmt.Errorf("error fetching jobs for run %d: %w", run.ID, err)
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
		close(errChan)
	}()

	writerWg.Wait()
	writer.Flush()

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("job fetch errors: %v", errs)
	}

	return nil
}
