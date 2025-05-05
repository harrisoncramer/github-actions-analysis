package analysis

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type JobStats struct {
	Durations []int
}

type NewAnalyzerParams struct {
	InputPath  string
	OutputPath string
	StartDate  *time.Time
	EndDate    *time.Time
}

type Analyzer struct {
	inputPath  string
	outputPath string
	startDate  *time.Time
	endDate    *time.Time
}

type JobDataHeaderIdxs struct {
	RunIdIdx        int
	WorkflowNameIdx int
	JobNameIdx      int
	StatusIdx       int
	ConclusionIdx   int
	StartedAtIdx    int
	CompletedAtIdx  int
	DurationIdx     int
}

func NewAnalyzer(params NewAnalyzerParams) *Analyzer {
	return &Analyzer{
		startDate:  params.StartDate,
		endDate:    params.EndDate,
		inputPath:  params.InputPath,
		outputPath: params.OutputPath,
	}
}

func (a *Analyzer) Analyze() error {

	file, err := os.Open(fmt.Sprintf("data/%s", a.inputPath))
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headerIndexes, err := a.readHeaders(reader)

	if err != nil {
		return fmt.Errorf("failed to read run data: %w", err)
	}

	durations := a.collectDurations(collectDurationParams{
		r:            reader,
		startedAtIdx: headerIndexes.StartedAtIdx,
		durationIdx:  headerIndexes.DurationIdx,
		jobNameIdx:   headerIndexes.JobNameIdx,
	})

	analysis := a.performAnalysis(durations)

	return a.writeAnalysisToFile(analysis)
}

func (a *Analyzer) writeAnalysisToFile(records [][]string) error {
	fmt.Println("Writing analysis...")

	outFile, err := os.Create(filepath.Join("data", a.outputPath))
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

func (a *Analyzer) readHeaders(r *csv.Reader) (*JobDataHeaderIdxs, error) {
	fmt.Println("Reading headers from file:", a.inputPath)

	headers, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read headers: %w", err)
	}

	runIdIdx, workflowNameIdx, jobNameIdx, statusIdx, conclusionIdx, startedAtIdx, completedAtIdx, durationIdx := -1, -1, -1, -1, -1, -1, -1, -1
	for i, h := range headers {
		switch h {
		case "run_id":
			runIdIdx = i
		case "workflow_name":
			workflowNameIdx = i
		case "job_name":
			jobNameIdx = i
		case "status":
			statusIdx = i
		case "conclusion":
			conclusionIdx = i
		case "started_at":
			startedAtIdx = i
		case "completed_at":
			completedAtIdx = i
		case "duration_seconds":
			durationIdx = i
		}
	}
	if runIdIdx == -1 || workflowNameIdx == -1 || jobNameIdx == -1 || statusIdx == -1 || conclusionIdx == -1 || startedAtIdx == -1 || completedAtIdx == -1 || durationIdx == -1 {
		return nil, fmt.Errorf("expected all specified columns")
	}

	return &JobDataHeaderIdxs{
		RunIdIdx:        runIdIdx,
		WorkflowNameIdx: workflowNameIdx,
		JobNameIdx:      jobNameIdx,
		StatusIdx:       statusIdx,
		ConclusionIdx:   conclusionIdx,
		StartedAtIdx:    startedAtIdx,
		CompletedAtIdx:  completedAtIdx,
		DurationIdx:     durationIdx,
	}, nil
}
