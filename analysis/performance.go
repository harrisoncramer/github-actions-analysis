package analysis

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"
)

type DurationLookup map[string]*JobStats

type collectDurationParams struct {
	r            *csv.Reader
	startedAtIdx int
	durationIdx  int
	jobNameIdx   int
}

func (a *Analyzer) collectDurations(params collectDurationParams) DurationLookup {
	fmt.Println("Collecting durations...")

	jobDurations := make(map[string]*JobStats)

	for {
		record, err := params.r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Skipping record due to error: %v\n", err)
			continue
		}

		if a.startDate != nil {
			startedAt, err := time.Parse(time.RFC3339, record[params.startedAtIdx])
			if err != nil || startedAt.Before(*a.startDate) {
				continue
			}
		}

		if a.endDate != nil {
			startedAt, err := time.Parse(time.RFC3339, record[params.startedAtIdx])
			if err != nil || startedAt.After(*a.endDate) {
				continue
			}
		}

		dur, err := strconv.Atoi(record[params.durationIdx])
		if err != nil {
			fmt.Printf("Unable to read record duration: %v\n", err)
			continue
		}

		job := record[params.jobNameIdx]
		if _, ok := jobDurations[job]; !ok {
			jobDurations[job] = &JobStats{}
		}
		jobDurations[job].Durations = append(jobDurations[job].Durations, dur)
	}

	return jobDurations
}

func (a *Analyzer) performAnalysis(durations DurationLookup) [][]string {
	fmt.Printf("Performing analysis of %d records...\n", len(durations))

	records := [][]string{}

	// Add headers to analysis CSV
	records = append(records, []string{
		"Job Name",
		"Count",
		"Avg (seconds)",
		"Min (seconds)",
		"Max (seconds)",
		"P90 (seconds)",
		"P99 (seconds)",
	})

	for job, stats := range durations {
		durs := stats.Durations
		sort.Ints(durs)

		count := len(durs)
		sum := 0
		for _, d := range durs {
			sum += d
		}
		avg := float64(sum) / float64(count)
		min := durs[0]
		max := durs[len(durs)-1]
		p90 := percentile(durs, 0.90)
		p99 := percentile(durs, 0.99)

		record := []string{
			job,
			strconv.Itoa(count),
			fmt.Sprintf("%.2f", avg),
			strconv.Itoa(min),
			strconv.Itoa(max),
			strconv.Itoa(p90),
			strconv.Itoa(p99),
		}
		records = append(records, record)
	}

	return records
}
