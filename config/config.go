package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	maxPages, err := strconv.Atoi(getEnv("MAX_PAGES", "1"))
	if err != nil {
		maxPages = 10
	}
	maxWorkers, err := strconv.Atoi(getEnv("MAX_WORKERS", "10"))
	if err != nil {
		maxWorkers = 10
	}
	collectOutfile := getEnv("COLLECT_OUTFILE", "runs.csv")
	analysisOutfile := getEnv("ANALYSIS_OUTFILE", "analysis.csv")

	githubRepo := getEnv("GITHUB_REPO", "")
	if githubRepo == "" {
		log.Fatal("No Github repo provided")
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
