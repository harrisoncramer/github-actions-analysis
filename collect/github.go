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
	perPage    = 100
	apiBaseURL = "https://api.github.com"
	maxWorkers = 10
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
	RunID        int64
	WorkflowName string
	Job          Job
}

type GitHubClient struct {
	Repo       string
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type GithubClientParams struct {
	repo string
}

func NewGitHubClient(params GithubClientParams) *GitHubClient {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable not set")
	}
	return &GitHubClient{
		Repo:       params.repo,
		BaseURL:    apiBaseURL,
		Token:      token,
		HTTPClient: &http.Client{},
	}
}

func (c *GitHubClient) get(endpoint string) ([]byte, error) {
	url := c.BaseURL + endpoint
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/vnd.github+json")

	fmt.Printf("Fetching %s\n", url)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *GitHubClient) FetchWorkflowRuns(page int) ([]WorkflowRun, error) {
	endpoint := fmt.Sprintf("/repos/%s/actions/runs?per_page=%d&page=%d", c.Repo, perPage, page)
	data, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}
	var result WorkflowRunsResponse
	err = json.Unmarshal(data, &result)
	return result.WorkflowRuns, err
}

func (c *GitHubClient) FetchJobsForRun(run WorkflowRun) ([]JobRecord, error) {
	endpoint := fmt.Sprintf("/repos/%s/actions/runs/%d/jobs", c.Repo, run.ID)
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
			RunID:        run.ID,
			WorkflowName: run.Name,
			Job:          job,
		})
	}
	return records, nil
}
