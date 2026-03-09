package session

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// LinterPatterns maps linter names to their output patterns
var LinterPatterns = map[string]*regexp.Regexp{
	"reek":       regexp.MustCompile(`(?i)warning|smell|--\s+\d+\s+warning`),
	"rubocop":    regexp.MustCompile(`(?i)offenses? detected|C:|W:|E:|no offenses`),
	"eslint":     regexp.MustCompile(`(?i)\d+ problems?|\d+ errors?.*\d+ warnings?`),
	"prettier":   regexp.MustCompile(`(?i)Code style issues found|Forgot to run Prettier`),
	"ruff":       regexp.MustCompile(`(?i)Found \d+ errors?|ruff check`),
	"standardrb": regexp.MustCompile(`(?i)standard.*offenses?`),
}

// ReekSmellPattern matches reek output lines
var ReekSmellPattern = regexp.MustCompile(`\[([^\]]*?):(\d+)\]:\s+(\w+):\s+(.+)`)

// RubocopOffensePattern matches rubocop offense lines
var RubocopOffensePattern = regexp.MustCompile(`([A-Z]\w+\/\w+):\s+(.+)`)

// SessionSummary contains analyzed session data
type SessionSummary struct {
	SessionID          string              `json:"session_id"`
	Project            string              `json:"project"`
	DurationMinutes    float64             `json:"duration_minutes"`
	TotalTurns         int                 `json:"total_turns"`
	TotalAssistantTurns int                `json:"total_assistant_turns"`
	TokenUsage         TokenUsage          `json:"token_usage"`
	LinterLoops        []LinterLoop        `json:"linter_loops"`
	ToolFailures       []ToolFailure       `json:"tool_failures"`
	RepeatedSequences  []RepeatedSequence  `json:"repeated_sequences"`
	LargeReads         []LargeRead         `json:"large_reads"`
	PermissionEvents   []PermissionEvent   `json:"permission_events"`
	HookFailures       []HookFailure       `json:"hook_failures"`
	EditCount          int                 `json:"edit_count"`
	ToolCallCount      int                 `json:"tool_call_count"`
	AgentSpawnCount    int                 `json:"agent_spawn_count"`
	FileReadTracker    map[string]int      `json:"-"`
}

// TokenUsage tracks token consumption
type TokenUsage struct {
	Input         int `json:"input"`
	Output        int `json:"output"`
	CacheRead     int `json:"cache_read"`
	CacheCreation int `json:"cache_creation"`
}

// LinterLoop represents a detected linter loop
type LinterLoop struct {
	Linter      string   `json:"linter"`
	File        string   `json:"file"`
	Occurrences int      `json:"occurrences"`
	Smells      []string `json:"smells"`
}

// ToolFailure represents a failed tool execution
type ToolFailure struct {
	Tool        string `json:"tool"`
	Occurrences int    `json:"occurrences"`
	LastError   string `json:"last_error"`
}

// RepeatedSequence represents a repeated workflow pattern
type RepeatedSequence struct {
	Pattern     string `json:"pattern"`
	Occurrences int    `json:"occurrences"`
}

// LargeRead represents frequently read files
type LargeRead struct {
	File      string `json:"file"`
	ReadCount int    `json:"read_count"`
}

// PermissionEvent represents a permission request
type PermissionEvent struct {
	Tool    string `json:"tool"`
	Command string `json:"command"`
	Count   int    `json:"count"`
}

// HookFailure represents a hook failure
type HookFailure struct {
	Hook  string `json:"hook"`
	Error string `json:"error"`
	Count int    `json:"count"`
}

// Parser parses Cursor session data from SQLite
type Parser struct {
	dbPath string
}

// NewParser creates a new session parser
// NOTE: This is for Claude Code JSONL format. For Cursor, use NewCursorAdapter() instead.
func NewParser() (*Parser, error) {
	// Determine Cursor DB path based on OS
	dbPath, err := getCursorDBPath()
	if err != nil {
		return nil, err
	}

	return &Parser{dbPath: dbPath}, nil
}

// ParseCursorSession parses a Cursor JSONL transcript
func ParseCursorSession(sessionPath string) (*SessionSummary, error) {
	adapter, err := NewCursorAdapter()
	if err != nil {
		return nil, err
	}

	return adapter.ParseTranscript(sessionPath)
}

// FindRecentCursorSession finds the most recent Cursor transcript
func FindRecentCursorSession() (string, error) {
	adapter, err := NewCursorAdapter()
	if err != nil {
		return "", err
	}

	return adapter.FindRecentTranscript()
}

// ParseSession parses a Cursor session and returns a summary
func (p *Parser) ParseSession(sessionID string) (*SessionSummary, error) {
	db, err := sql.Open("sqlite3", p.dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Query session data from cursorDiskKV table
	// This is a placeholder - actual schema needs to be discovered
	rows, err := db.Query(`
		SELECT key, value
		FROM cursorDiskKV
		WHERE key LIKE ?
	`, fmt.Sprintf("%%session:%s%%", sessionID))
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()

	summary := &SessionSummary{
		SessionID:       sessionID,
		FileReadTracker: make(map[string]int),
		TokenUsage:      TokenUsage{},
	}

	// Process rows and build summary
	// This is a placeholder implementation
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			continue
		}

		// Parse the session data
		// Actual parsing logic depends on Cursor's SQLite schema
		p.processEntry(summary, value)
	}

	return summary, nil
}

// processEntry processes a single session entry
func (p *Parser) processEntry(summary *SessionSummary, data string) {
	// Parse JSON data
	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		return
	}

	// Extract metadata, token usage, detect patterns, etc.
	// This is placeholder logic - needs actual Cursor session format
	p.detectLinterLoops(summary, entry)
	p.detectToolFailures(summary, entry)
	p.detectLargeReads(summary, entry)
}

// detectLinterLoops detects linter loop patterns
func (p *Parser) detectLinterLoops(summary *SessionSummary, entry map[string]interface{}) {
	// Check for linter output patterns
	content, ok := entry["content"].(string)
	if !ok {
		return
	}

	for linter, pattern := range LinterPatterns {
		if pattern.MatchString(content) {
			// Track linter occurrence
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
}

// detectToolFailures detects repeated tool failures
func (p *Parser) detectToolFailures(summary *SessionSummary, entry map[string]interface{}) {
	// Placeholder - detect tool execution failures
	toolName, _ := entry["tool"].(string)
	error, hasError := entry["error"].(string)

	if hasError && toolName != "" {
		found := false
		for i := range summary.ToolFailures {
			if summary.ToolFailures[i].Tool == toolName {
				summary.ToolFailures[i].Occurrences++
				summary.ToolFailures[i].LastError = error
				found = true
				break
			}
		}
		if !found {
			summary.ToolFailures = append(summary.ToolFailures, ToolFailure{
				Tool:        toolName,
				Occurrences: 1,
				LastError:   error,
			})
		}
	}
}

// detectLargeReads tracks frequently read files
func (p *Parser) detectLargeReads(summary *SessionSummary, entry map[string]interface{}) {
	filePath, ok := entry["file_path"].(string)
	if !ok {
		return
	}

	summary.FileReadTracker[filePath]++

	// If read 3+ times, add to large reads
	if summary.FileReadTracker[filePath] >= 3 {
		found := false
		for i := range summary.LargeReads {
			if summary.LargeReads[i].File == filePath {
				summary.LargeReads[i].ReadCount = summary.FileReadTracker[filePath]
				found = true
				break
			}
		}
		if !found {
			summary.LargeReads = append(summary.LargeReads, LargeRead{
				File:      filePath,
				ReadCount: summary.FileReadTracker[filePath],
			})
		}
	}
}

// getCursorDBPath returns the path to Cursor's SQLite database
func getCursorDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Try different OS paths
	paths := []string{
		filepath.Join(home, "Library", "Application Support", "Cursor", "User", "globalStorage", "state.vscdb"),
		filepath.Join(home, ".config", "Cursor", "User", "globalStorage", "state.vscdb"),
		filepath.Join(os.Getenv("APPDATA"), "Cursor", "User", "globalStorage", "state.vscdb"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("could not find Cursor database")
}

// FindRecentSession finds the most recent session ID
func (p *Parser) FindRecentSession() (string, error) {
	db, err := sql.Open("sqlite3", p.dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Query for recent sessions
	// Placeholder - actual query depends on Cursor's schema
	var sessionID string
	err = db.QueryRow(`
		SELECT key
		FROM cursorDiskKV
		WHERE key LIKE 'session:%'
		ORDER BY rowid DESC
		LIMIT 1
	`).Scan(&sessionID)

	if err != nil {
		return "", fmt.Errorf("no recent session found: %w", err)
	}

	// Extract session ID from key
	parts := strings.Split(sessionID, ":")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid session key format")
	}

	return parts[1], nil
}
