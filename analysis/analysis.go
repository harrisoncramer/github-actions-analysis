package analysis

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type JobStats struct {
	Durations []int
}

type AnalyzeParams struct {
	InputPath  string // formerly Outfile, renamed for clarity
	OutputPath string
	StartDate  *time.Time
	EndDate    *time.Time
}

type Analyzer struct {
	inputPath    string
	outputPath   string
	startDate    *time.Time
	endDate      *time.Time
	jobDurations map[string]*JobStats
}

func NewAnalyzer(params AnalyzeParams) *Analyzer {
	return &Analyzer{
		startDate:    params.StartDate,
		endDate:      params.EndDate,
		inputPath:    params.InputPath,
		outputPath:   params.OutputPath,
		jobDurations: make(map[string]*JobStats),
	}
}

func (a *Analyzer) Analyze() error {
	file, err := os.Open(fmt.Sprintf("data/%s", a.inputPath))
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	jobNameIdx, durationIdx, startedAtIdx := -1, -1, -1
	for i, h := range headers {
		switch h {
		case "job_name":
			jobNameIdx = i
		case "duration_seconds":
			durationIdx = i
		case "started_at":
			startedAtIdx = i
		}
	}
	if jobNameIdx == -1 || durationIdx == -1 || startedAtIdx == -1 {
		return fmt.Errorf("expected 'job_name', 'duration_seconds', and 'started_at' columns")
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) <= startedAtIdx {
			continue
		}

		if a.startDate != nil && a.endDate != nil {
			startedAt, err := time.Parse(time.RFC3339, record[startedAtIdx])
			if err != nil || startedAt.Before(*a.startDate) || startedAt.After(*a.endDate) {
				continue
			}
		}

		job := record[jobNameIdx]
		dur, err := strconv.Atoi(record[durationIdx])
		if err != nil {
			continue
		}

		if _, ok := a.jobDurations[job]; !ok {
			a.jobDurations[job] = &JobStats{}
		}
		a.jobDurations[job].Durations = append(a.jobDurations[job].Durations, dur)
	}

	return a.writeAnalysisToFile()
}

func (a *Analyzer) writeAnalysisToFile() error {
	outFile, err := os.Create(filepath.Join("data", a.outputPath))
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	headers := []string{"Job Name", "Count", "Avg (s)", "Min (s)", "Max (s)", "P90 (s)", "P99 (s)"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for job, stats := range a.jobDurations {
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
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

func percentile(sorted []int, p float64) int {
	if len(sorted) == 0 {
		return 0
	}
	k := int(float64(len(sorted)-1) * p)
	return sorted[k]
}
