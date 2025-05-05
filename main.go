package main

import (
	"gh-analysis/analysis"
	"gh-analysis/collect"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gh-analysis",
		Short: "CLI for data collection and analysis",
	}

	var collectCmd = &cobra.Command{
		Use:   "collect",
		Short: "Collect data from the specified GitHub repository",
		Run: func(cmd *cobra.Command, args []string) {
			c := collect.LoadCollectConfig()
			collect.Collect(collect.CollectParams{
				GithubRepo: c.GithubRepo,
				MaxWorkers: c.MaxWorkers,
				MaxPages:   c.MaxPages,
				Outfile:    c.OutputPath,
				PerPage:    100,
			})
		},
	}

	var analyzeCmd = &cobra.Command{
		Use:   "analyze",
		Short: "Analyze the collected data",
		Run: func(cmd *cobra.Command, args []string) {
			c := analysis.LoadAnalysisConfig()
			analyzer := analysis.AnalyzePerformance(analysis.AnalyzePerformanceParams{
				InputPath:  c.InputPath,
				OutputPath: c.OutputPath,
				StartDate:  c.AnalysisStartDate,
				EndDate:    c.AnalysisEndDate,
			})
			err := analyzer.Analyze()
			if err != nil {
				log.Fatalf("Failed to perform analysis: %v", err)
			}
		},
	}

	rootCmd.AddCommand(collectCmd)
	rootCmd.AddCommand(analyzeCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
