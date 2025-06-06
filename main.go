package main

import (
	"github-actions-analysis/analysis"
	"github-actions-analysis/collect"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "github-actions-analysis",
		Short: "CLI for data collection and analysis",
	}

	var collectCmd = &cobra.Command{
		Use:   "collect",
		Short: "Collect data from the specified GitHub repository",
		Run: func(cmd *cobra.Command, args []string) {
			c := collect.LoadCollectConfig()
			err := collect.Collect(collect.CollectParams{
				GithubRepo: c.GithubRepo,
				MaxWorkers: c.MaxWorkers,
				MaxPages:   c.MaxPages,
				Outfile:    c.OutputPath,
				PerPage:    100,
			})
			if err != nil {
				log.Fatalf("Failed to collect data: %v", err)
			}

		},
	}

	var analyzeCmd = &cobra.Command{
		Use:   "analyze",
		Short: "Analyze the collected data",
		Run: func(cmd *cobra.Command, args []string) {
			c := analysis.LoadAnalysisConfig()
			analyzer := analysis.NewAnalyzer(analysis.NewAnalyzerParams{
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
