package main

import (
	"analysis/analysis"
	"analysis/collect"
	"analysis/config"
)

func main() {
	c := config.LoadConfig()

	collect.Collect(collect.CollectParams{
		GithubRepo: c.GithubRepo,
		MaxPages:   c.MaxPages,
		MaxWorkers: c.MaxWorkers,
		Outfile:    c.CollectOutfile,
	})
	analysis.Analyze(analysis.AnalyzeParams{
		Outfile: c.AnalysisOutfile,
	})
}
