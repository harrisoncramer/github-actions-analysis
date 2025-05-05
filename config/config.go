package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Repo              string
	MaxPages          int
	MaxWorkers        int
	CollectOutfile    string
	GithubRepo        string
	AnalysisOutfile   string
	AnalysisStartDate *time.Time
	AnalysisEndDate   *time.Time
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

	analysisStartDate := getEnv("ANALYSIS_START_DATE", "")
	var analysisStartDateParsed *time.Time
	if analysisStartDate != "" {
		parsedTime, err := time.Parse(time.RFC3339, analysisStartDate)
		if err == nil {
			analysisStartDateParsed = &parsedTime
		}
	}

	analysisEndDate := getEnv("ANALYSIS_END_DATE", "")
	var analysisEndDateParsed *time.Time
	if analysisEndDate != "" {
		parsedTime, err := time.Parse(time.RFC3339, analysisEndDate)
		if err == nil {
			analysisEndDateParsed = &parsedTime
		}
	}

	return Config{
		MaxPages:          maxPages,
		MaxWorkers:        maxWorkers,
		CollectOutfile:    collectOutfile,
		GithubRepo:        githubRepo,
		AnalysisStartDate: analysisStartDateParsed,
		AnalysisEndDate:   analysisEndDateParsed,
		AnalysisOutfile:   analysisOutfile,
	}

}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
