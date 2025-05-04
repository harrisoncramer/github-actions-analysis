package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Repo            string
	MaxPages        int
	MaxWorkers      int
	CollectOutfile  string
	AnalysisOutfile string
	GithubRepo      string
}

func LoadConfig() Config {
	maxPages, err := strconv.Atoi(getEnv("MAX_PAGES", "1"))
	if err != nil {
		maxPages = 1
	}
	maxWorkers, err := strconv.Atoi(getEnv("MAX_WORKERS", "10"))
	if err != nil {
		maxWorkers = 10
	}
	collectOutfile := getEnv("COLLECT_OUTFILE", "data/runs.csv")
	analysisOutfile := getEnv("ANALYSIS_OUTFILE", "data/analysis.csv")

	githubRepo := getEnv("GITHUB_REPO", "")
	if githubRepo == "" {
		log.Fatal("No GIthub repo provided")
	}

	return Config{
		MaxPages:        maxPages,
		MaxWorkers:      maxWorkers,
		CollectOutfile:  collectOutfile,
		AnalysisOutfile: analysisOutfile,
		GithubRepo:      githubRepo,
	}

}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
