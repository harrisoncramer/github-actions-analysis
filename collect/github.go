package collect

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	apiBaseURL = "https://api.github.com"
)

type WorkflowRun struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type WorkflowRunsResponse struct {
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

type Job struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
	RunAttempt  int       `json:"run_attempt"`
}

type JobsResponse struct {
	Jobs []Job `json:"jobs"`
}

type JobRecord struct {
	runID        int64
	workflowName string
	job          Job
}

type GitHubClient struct {
	repo    string
	baseURL string
	token   string
	client  *http.Client
	perPage int
}

type GithubClientParams struct {
	repo       string
	maxWorkers int
	perPage    int
}

func NewGitHubClient(params GithubClientParams) *GitHubClient {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable not set")
	}
	return &GitHubClient{
		repo:    params.repo,
		baseURL: apiBaseURL,
		token:   token,
		perPage: params.perPage,
		client:  &http.Client{},
	}
}

func (c *GitHubClient) get(endpoint string) ([]byte, error) {
	url := c.baseURL + endpoint
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")

	fmt.Printf("Fetching %s\n", url)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *GitHubClient) FetchWorkflowRuns(page int) ([]WorkflowRun, error) {
	endpoint := fmt.Sprintf("/repos/%s/actions/runs?per_page=%d&page=%d", c.repo, c.perPage, page)
	data, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}
	var result WorkflowRunsResponse
	err = json.Unmarshal(data, &result)
	return result.WorkflowRuns, err
}

func (c *GitHubClient) FetchJobsForRun(run WorkflowRun) ([]JobRecord, error) {
	endpoint := fmt.Sprintf("/repos/%s/actions/runs/%d/jobs", c.repo, run.ID)
	data, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}
	var result JobsResponse
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	var records []JobRecord
	for _, job := range result.Jobs {
		records = append(records, JobRecord{
			runID:        run.ID,
			workflowName: run.Name,
			job:          job,
		})
	}
	return records, nil
}
