package analysis

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

type JobStats struct {
	Durations []int
}

func percentile(sorted []int, p float64) int {
	if len(sorted) == 0 {
		return 0
	}
	k := int(float64(len(sorted)-1) * p)
	return sorted[k]
}

func writeAnalysisToFile(jobDurations map[string]*JobStats, outputPath string) error {
	outFile, err := os.Create(outputPath)
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

	for job, stats := range jobDurations {
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

type AnalyzeParams struct {
	Outfile string
}

func Analyze(params AnalyzeParams) {
	file, err := os.Open(params.Outfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, _ := reader.Read() // skip header
	jobNameIdx, durationIdx := -1, -1
	for i, h := range headers {
		switch h {
		case "job_name":
			jobNameIdx = i
		case "duration_seconds":
			durationIdx = i
		}
	}
	if jobNameIdx == -1 || durationIdx == -1 {
		panic("Expected 'job_name' and 'duration_seconds' columns")
	}

	jobDurations := map[string]*JobStats{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		job := record[jobNameIdx]
		dur, err := strconv.Atoi(record[durationIdx])
		if err != nil {
			continue
		}

		if _, ok := jobDurations[job]; !ok {
			jobDurations[job] = &JobStats{}
		}
		jobDurations[job].Durations = append(jobDurations[job].Durations, dur)
	}

	if err := writeAnalysisToFile(jobDurations, "data/analysis.csv"); err != nil {
		panic(err)
	}
}
