package analysis

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type AnalysisConfig struct {
	InputPath         string
	OutputPath        string
	AnalysisStartDate *time.Time
	AnalysisEndDate   *time.Time
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
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	fmt.Printf("%s not found, returning default '%s'\n", key, fallback)
	return fallback
}
