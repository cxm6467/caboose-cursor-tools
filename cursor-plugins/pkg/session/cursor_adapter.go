package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// CursorMessage represents a message in Cursor's JSONL format
type CursorMessage struct {
	Role    string `json:"role"`
	Message struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"message"`
}

// CursorAdapter adapts Cursor's JSONL transcript format to our parser
type CursorAdapter struct {
	projectsPath string
}

// NewCursorAdapter creates a new Cursor transcript adapter
func NewCursorAdapter() (*CursorAdapter, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	projectsPath := filepath.Join(home, ".cursor", "projects")
	return &CursorAdapter{projectsPath: projectsPath}, nil
}

// FindRecentTranscript finds the most recent Cursor transcript
func (a *CursorAdapter) FindRecentTranscript() (string, error) {
	type transcriptInfo struct {
		path    string
		modTime time.Time
	}

	var transcripts []transcriptInfo

	// Walk through all projects
	projects, err := os.ReadDir(a.projectsPath)
	if err != nil {
		return "", fmt.Errorf("failed to read projects directory: %w", err)
	}

	for _, project := range projects {
		if !project.IsDir() {
			continue
		}

		transcriptsDir := filepath.Join(a.projectsPath, project.Name(), "agent-transcripts")
		if _, err := os.Stat(transcriptsDir); os.IsNotExist(err) {
			continue
		}

		// List all transcript UUIDs
		uuids, err := os.ReadDir(transcriptsDir)
		if err != nil {
			continue
		}

		for _, uuid := range uuids {
			if !uuid.IsDir() {
				continue
			}

			jsonlPath := filepath.Join(transcriptsDir, uuid.Name(), uuid.Name()+".jsonl")
			info, err := os.Stat(jsonlPath)
			if err != nil {
				continue
			}

			transcripts = append(transcripts, transcriptInfo{
				path:    jsonlPath,
				modTime: info.ModTime(),
			})
		}
	}

	if len(transcripts) == 0 {
		return "", fmt.Errorf("no transcripts found")
	}

	// Sort by modification time (most recent first)
	sort.Slice(transcripts, func(i, j int) bool {
		return transcripts[i].modTime.After(transcripts[j].modTime)
	})

	return transcripts[0].path, nil
}

// ParseTranscript parses a Cursor JSONL transcript file
func (a *CursorAdapter) ParseTranscript(path string) (*SessionSummary, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open transcript: %w", err)
	}
	defer file.Close()

	// Extract session ID from filename
	sessionID := filepath.Base(path)
	sessionID = strings.TrimSuffix(sessionID, ".jsonl")

	summary := &SessionSummary{
		SessionID:       sessionID,
		FileReadTracker: make(map[string]int),
		TokenUsage:      TokenUsage{},
	}

	scanner := bufio.NewScanner(file)
	// Increase buffer size for large transcript lines (default is 64KB, use 10MB)
	const maxCapacity = 10 * 1024 * 1024 // 10MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	lineNum := 0

	for scanner.Scan() {
		lineNum++
		var msg CursorMessage
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			// Skip malformed lines
			continue
		}

		if msg.Role == "user" {
			summary.TotalTurns++
		} else if msg.Role == "assistant" {
			summary.TotalAssistantTurns++
		}

		// Extract text content
		for _, content := range msg.Message.Content {
			if content.Type == "text" {
				a.analyzeContent(summary, content.Text, msg.Role)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading transcript: %w", err)
	}

	return summary, nil
}

// analyzeContent analyzes message content for patterns
func (a *CursorAdapter) analyzeContent(summary *SessionSummary, text, role string) {
	// Detect linter patterns
	for linter, pattern := range LinterPatterns {
		if pattern.MatchString(text) {
			found := false
			for i := range summary.LinterLoops {
				if summary.LinterLoops[i].Linter == linter {
					summary.LinterLoops[i].Occurrences++
					found = true
					break
				}
			}
			if !found {
				summary.LinterLoops = append(summary.LinterLoops, LinterLoop{
					Linter:      linter,
					Occurrences: 1,
					Smells:      []string{},
				})
			}
		}
	}

	// Detect file reads (look for file path patterns)
	if strings.Contains(text, "Reading") || strings.Contains(text, "file:") {
		// Simple heuristic: extract paths
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			if strings.Contains(line, "/") && (strings.Contains(line, ".ts") ||
				strings.Contains(line, ".js") || strings.Contains(line, ".go") ||
				strings.Contains(line, ".py") || strings.Contains(line, ".json")) {

				// Extract potential file path
				parts := strings.Fields(line)
				for _, part := range parts {
					if strings.Contains(part, "/") && strings.Contains(part, ".") {
						summary.FileReadTracker[part]++

						if summary.FileReadTracker[part] >= 3 {
							found := false
							for i := range summary.LargeReads {
								if summary.LargeReads[i].File == part {
									summary.LargeReads[i].ReadCount = summary.FileReadTracker[part]
									found = true
									break
								}
							}
							if !found {
								summary.LargeReads = append(summary.LargeReads, LargeRead{
									File:      part,
									ReadCount: summary.FileReadTracker[part],
								})
							}
						}
					}
				}
			}
		}
	}

	// Detect tool failures (error patterns)
	if strings.Contains(text, "error") || strings.Contains(text, "failed") ||
		strings.Contains(text, "Error:") {
		// Extract error context
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			lowerLine := strings.ToLower(line)
			if strings.Contains(lowerLine, "error") || strings.Contains(lowerLine, "failed") {
				// Simple categorization
				var tool string
				if strings.Contains(lowerLine, "bash") || strings.Contains(lowerLine, "command") {
					tool = "Bash"
				} else if strings.Contains(lowerLine, "read") {
					tool = "Read"
				} else if strings.Contains(lowerLine, "write") || strings.Contains(lowerLine, "edit") {
					tool = "Write/Edit"
				} else {
					tool = "Unknown"
				}

				found := false
				for i := range summary.ToolFailures {
					if summary.ToolFailures[i].Tool == tool {
						summary.ToolFailures[i].Occurrences++
						summary.ToolFailures[i].LastError = line
						found = true
						break
					}
				}
				if !found && tool != "Unknown" {
					summary.ToolFailures = append(summary.ToolFailures, ToolFailure{
						Tool:        tool,
						Occurrences: 1,
						LastError:   line,
					})
				}
				break // Only count once per message
			}
		}
	}
}
