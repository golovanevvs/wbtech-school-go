package service

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.2/mygrep/internal/model"
)

// GrepService provides grep functionality
type GrepService struct {
	config *model.Config
}

// NewGrepService creates a new grep service
func NewGrepService(config *model.Config) *GrepService {
	return &GrepService{
		config: config,
	}
}

// ExecuteGrep performs grep operation based on configuration
func (s *GrepService) ExecuteGrep() error {
	if len(s.config.Files) > 0 {
		for _, filename := range s.config.Files {
			if err := s.processFile(filename); err != nil {
				return fmt.Errorf("error processing file %s: %v", filename, err)
			}
		}
	} else {
		if err := s.processStream(os.Stdin, "stdin"); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteDistributedGrep performs distributed grep operation
func (s *GrepService) ExecuteDistributedGrep() error {
	fmt.Println("Starting distributed grep...")

	if len(s.config.Files) == 0 {
		return fmt.Errorf("distributed mode requires file arguments")
	}

	// Process each file in distributed mode
	for _, filename := range s.config.Files {
		if err := s.processFileDistributed(filename); err != nil {
			return fmt.Errorf("distributed processing error for file %s: %v", filename, err)
		}
	}

	return nil
}

// processFile processes a single file
func (s *GrepService) processFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open file %s: %v", filename, err)
	}
	defer file.Close()

	return s.processStream(file, filename)
}

// processFileDistributed processes a file in distributed mode
func (s *GrepService) processFileDistributed(filename string) error {
	fmt.Printf("Processing file %s in distributed mode\n", filename)

	// Count total lines in file
	totalLines, err := s.countLinesInFile(filename)
	if err != nil {
		return err
	}

	fmt.Printf("File contains %d lines\n", totalLines)

	// Determine number of servers (current + peers)
	numServers := 1 + len(s.config.Peers)
	if numServers == 1 {
		// Only one server, process locally
		return s.processFile(filename)
	}

	// Split file and distribute jobs
	return s.distributeAndProcess(filename, totalLines, numServers)
}

// countLinesInFile counts lines in a file
func (s *GrepService) countLinesInFile(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
	}

	return lines, scanner.Err()
}

// distributeAndProcess distributes jobs between servers
func (s *GrepService) distributeAndProcess(filename string, totalLines, numServers int) error {
	linesPerServer := totalLines / numServers

	// Create jobs for each server
	jobs := make([]map[string]interface{}, 0, numServers)

	for i := 0; i < numServers; i++ {
		startLine := i*linesPerServer + 1
		endLine := (i + 1) * linesPerServer

		// Last server gets remaining lines
		if i == numServers-1 {
			endLine = totalLines
		}

		job := map[string]interface{}{
			"job_id":     fmt.Sprintf("%s-%d", s.config.ServerID, i),
			"server_id":  s.config.ServerID,
			"pattern":    s.config.Pattern,
			"start_line": startLine,
			"end_line":   endLine,
			"filename":   filename,
			"flags":      s.config.Flags,
			"created_at": time.Now(),
		}
		jobs = append(jobs, job)
	}

	// Execute jobs locally (integration with transport layer will be added)
	return s.executeJobsLocally(jobs)
}

// processStream processes a data stream
func (s *GrepService) processStream(reader io.Reader, sourceName string) error {
	scanner := bufio.NewScanner(reader)
	lineNumber := 0
	matches := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		result, found := s.searchInLine(line, lineNumber)
		if found {
			matches++
			if !s.config.Flags.Count {
				if len(s.config.Files) > 1 {
					fmt.Fprintf(s.config.Output, "%s:", sourceName)
				}
				s.printResult(result)
			}
		}
	}

	if s.config.Flags.Count {
		if len(s.config.Files) > 1 {
			fmt.Fprintf(s.config.Output, "%s:", sourceName)
		}
		fmt.Fprintf(s.config.Output, "%d\n", matches)
	}

	return nil
}

// SearchInLine searches for pattern in a line
func (s *GrepService) SearchInLine(line string, lineNumber int) (*model.GrepResult, bool) {
	return s.searchInLine(line, lineNumber)
}

// searchInLine performs pattern search in a line
func (s *GrepService) searchInLine(line string, lineNumber int) (*model.GrepResult, bool) {
	pattern := regexp.QuoteMeta(s.config.Pattern)
	if s.config.Flags.WholeLine {
		pattern = "^" + pattern + "$"
	}

	flags := ""
	if s.config.Flags.IgnoreCase {
		flags = "(?i)"
	}

	fullPattern := flags + pattern
	re, err := regexp.Compile(fullPattern)
	if err != nil {
		return nil, false
	}

	matches := re.FindAllString(line, -1)

	if s.config.Flags.InvertMatch {
		if len(matches) == 0 {
			return &model.GrepResult{
				LineNumber: lineNumber,
				Line:       line,
				Match:      line,
			}, true
		}
		return nil, false
	}

	if len(matches) > 0 {
		if s.config.Flags.OnlyMatching {
			for _, match := range matches {
				return &model.GrepResult{
					LineNumber: lineNumber,
					Line:       match,
					Match:      match,
				}, true
			}
		}

		return &model.GrepResult{
			LineNumber: lineNumber,
			Line:       line,
			Match:      matches[0],
		}, true
	}

	return nil, false
}

// ExecuteJobLocally executes a job locally and returns result
func (s *GrepService) ExecuteJobLocally(job map[string]interface{}) (*model.JobResult, error) {
	filename := job["filename"].(string)
	startLine := int(job["start_line"].(float64))
	endLine := int(job["end_line"].(float64))
	pattern := job["pattern"].(string)
	jobID := job["job_id"].(string)

	fmt.Printf("Executing local job %s: file=%s, lines=%d-%d, pattern=%s\n",
		jobID, filename, startLine, endLine, pattern)

	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return &model.JobResult{
			JobID:       jobID,
			ServerID:    s.config.ServerID,
			Processed:   0,
			Success:     false,
			Error:       err.Error(),
			CompletedAt: time.Now(),
		}, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	matches := make([]model.GrepResult, 0)

	for scanner.Scan() {
		currentLine++

		// Skip lines before startLine
		if currentLine < startLine {
			continue
		}

		// Stop after endLine
		if currentLine > endLine {
			break
		}

		line := scanner.Text()

		// Search in line
		result, found := s.searchInLineForJob(line, currentLine, job)
		if found {
			matches = append(matches, *result)
		}
	}

	if err := scanner.Err(); err != nil {
		return &model.JobResult{
			JobID:       jobID,
			ServerID:    s.config.ServerID,
			Processed:   currentLine - startLine + 1,
			Success:     false,
			Error:       err.Error(),
			CompletedAt: time.Now(),
		}, nil
	}

	return &model.JobResult{
		JobID:       jobID,
		ServerID:    s.config.ServerID,
		Matches:     matches,
		Processed:   currentLine - startLine + 1,
		Success:     true,
		CompletedAt: time.Now(),
	}, nil
}

// searchInLineForJob searches in line for specific job
func (s *GrepService) searchInLineForJob(line string, lineNumber int, job map[string]interface{}) (*model.GrepResult, bool) {
	pattern, _ := job["pattern"].(string)
	flags, _ := job["flags"].(model.GrepFlags)

	// Create temporary config for search
	tempConfig := &model.Config{
		Pattern: pattern,
		Flags:   flags,
		Output:  s.config.Output,
	}

	// Use existing search logic
	result, found := s.searchInLineWithConfig(line, lineNumber, tempConfig)
	return result, found
}

// searchInLineWithConfig searches with given config
func (s *GrepService) searchInLineWithConfig(line string, lineNumber int, config *model.Config) (*model.GrepResult, bool) {
	pattern := regexp.QuoteMeta(config.Pattern)
	if config.Flags.WholeLine {
		pattern = "^" + pattern + "$"
	}

	flags := ""
	if config.Flags.IgnoreCase {
		flags = "(?i)"
	}

	fullPattern := flags + pattern
	re, err := regexp.Compile(fullPattern)
	if err != nil {
		return nil, false
	}

	matches := re.FindAllString(line, -1)

	if config.Flags.InvertMatch {
		if len(matches) == 0 {
			return &model.GrepResult{
				LineNumber: lineNumber,
				Line:       line,
				Match:      line,
			}, true
		}
		return nil, false
	}

	if len(matches) > 0 {
		if config.Flags.OnlyMatching {
			for _, match := range matches {
				return &model.GrepResult{
					LineNumber: lineNumber,
					Line:       match,
					Match:      match,
				}, true
			}
		}

		return &model.GrepResult{
			LineNumber: lineNumber,
			Line:       line,
			Match:      matches[0],
		}, true
	}

	return nil, false
}

// executeJobsLocally executes jobs locally
func (s *GrepService) executeJobsLocally(jobs []map[string]interface{}) error {
	fmt.Printf("Executing %d jobs locally\n", len(jobs))

	// Execute each job locally
	for i, job := range jobs {
		filename := job["filename"].(string)
		startLine := int(job["start_line"].(float64))
		endLine := int(job["end_line"].(float64))

		fmt.Printf("Job %d: file=%s, lines=%d-%d\n",
			i+1, filename, startLine, endLine)

		// Execute grep on specified line range
		if err := s.executeJobLocally(job); err != nil {
			fmt.Printf("Error executing job %d: %v\n", i+1, err)
			return err
		}
	}

	return nil
}

// executeJobLocally executes a single job locally
func (s *GrepService) executeJobLocally(job map[string]interface{}) error {
	filename := job["filename"].(string)
	startLine := int(job["start_line"].(float64))
	endLine := int(job["end_line"].(float64))

	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	matches := 0

	for scanner.Scan() {
		currentLine++

		// Skip lines before startLine
		if currentLine < startLine {
			continue
		}

		// Stop after endLine
		if currentLine > endLine {
			break
		}

		line := scanner.Text()

		// Search in line
		result, found := s.searchInLineForJob(line, currentLine, job)
		if found {
			matches++

			// Print result
			if len(s.config.Files) > 1 {
				fmt.Fprintf(s.config.Output, "%s:", filename)
			}
			s.printResultForJob(result)
		}
	}

	if s.config.Flags.Count {
		if len(s.config.Files) > 1 {
			fmt.Fprintf(s.config.Output, "%s:", filename)
		}
		fmt.Fprintf(s.config.Output, "%d\n", matches)
	}

	return scanner.Err()
}

// printResult prints a grep result
func (s *GrepService) printResult(result *model.GrepResult) {
	if s.config.Flags.LineNumber {
		fmt.Fprintf(s.config.Output, "%d:", result.LineNumber)
	}

	if s.config.Flags.OnlyMatching {
		fmt.Fprintf(s.config.Output, "%s\n", result.Match)
	} else {
		fmt.Fprintf(s.config.Output, "%s\n", result.Line)
	}
}

// printResultForJob prints result for job
func (s *GrepService) printResultForJob(result *model.GrepResult) {
	if s.config.Flags.LineNumber {
		fmt.Fprintf(s.config.Output, "%d:", result.LineNumber)
	}

	if s.config.Flags.OnlyMatching {
		fmt.Fprintf(s.config.Output, "%s\n", result.Match)
	} else {
		fmt.Fprintf(s.config.Output, "%s\n", result.Line)
	}
}
