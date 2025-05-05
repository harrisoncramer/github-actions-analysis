package collect

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type CollectConfig struct {
	GithubRepo string
	MaxPages   int
	MaxWorkers int
	OutputPath string
}

func LoadCollectConfig() CollectConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	maxPages, err := strconv.Atoi(getEnv("COLLECT_MAX_PAGES", "1"))
	if err != nil {
		maxPages = 10
	}
	maxWorkers, err := strconv.Atoi(getEnv("COLLECT_MAX_WORKERS", "10"))
	if err != nil {
		maxWorkers = 10
	}
	outputPath := getEnv("COLLECT_OUTPUT_PATH", "runs.csv")
	githubRepo := os.Getenv("COLLECT_GITHUB_REPO")
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

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	fmt.Printf("%s not found, returning default '%s'\n", key, fallback)
	return fallback
}
