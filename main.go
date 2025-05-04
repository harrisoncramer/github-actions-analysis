package main

import (
	"analysis/analysis"
	"analysis/collect"
	"analysis/config"
	"log"
)

func main() {
	c := config.LoadConfig()

	collect.Collect(collect.CollectParams{
		GithubRepo: c.GithubRepo,
		MaxWorkers: c.MaxWorkers,
		MaxPages:   c.MaxPages,
		Outfile:    c.CollectOutfile,
		PerPage:    100,
	})

	analyzer := analysis.NewAnalyzer(analysis.AnalyzeParams{
		InputPath:  c.CollectOutfile,
		OutputPath: c.AnalysisOutfile,
	})
	err := analyzer.Analyze()
	if err != nil {
		log.Fatalf("Failed to perform analysis: %v", err)
	}

}
