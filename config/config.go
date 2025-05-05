package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type CollectConfig struct {
	GithubRepo string
	MaxPages   int
	MaxWorkers int
	OutputPath string
}

type AnalysisConfig struct {
	InputPath         string
	OutputPath        string
	AnalysisStartDate *time.Time
	AnalysisEndDate   *time.Time
}

func LoadCollectConfig() CollectConfig {
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
	outputPath := getEnv("COLLECT_OUTPUT_PATH", "runs.csv")
	githubRepo := getEnv("GITHUB_REPO", "")
	if githubRepo == "" {
		log.Fatal("No Github repo provided")
	}
	return CollectConfig{
		GithubRepo: githubRepo,
		MaxPages:   maxPages,
		MaxWorkers: maxWorkers,
		OutputPath: outputPath,
	}
}

func LoadAnalysisConfig() AnalysisConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	inputPath := getEnv("ANALYSIS_INPUT_PATH", "runs.csv")
	outputPath := getEnv("ANALYSIS_OUTPUT_PATH", "analysis.csv")
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
	return AnalysisConfig{
		InputPath:         inputPath,
		OutputPath:        outputPath,
		AnalysisStartDate: analysisStartDateParsed,
		AnalysisEndDate:   analysisEndDateParsed,
	}
}

func getEnv(key, fallback string) string {
	if _, exists := os.LookupEnv(key); !exists {
		panic("Missing environment variable: " + key)
	}
	return fallback
}
