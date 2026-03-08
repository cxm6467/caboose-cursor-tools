package installer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LockFile represents the crules-lock.json file
type LockFile struct {
	Version  string               `json:"version"` // Lock file schema version
	Rules    map[string]RuleLock  `json:"rules"`   // Key is rule name
	FilePath string               `json:"-"`       // Path to lock file (not serialized)
}

// RuleLock represents metadata about an installed rule
type RuleLock struct {
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	Source       string    `json:"source"`       // "local", "marketplace", URL
	RegistryName string    `json:"registryName,omitempty"` // Marketplace name if applicable
	InstalledAt  time.Time `json:"installedAt"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
	Checksum     string    `json:"checksum,omitempty"` // SHA-256 of content
}

// NewLockFile creates a new lock file instance
func NewLockFile(lockFilePath string) *LockFile {
	return &LockFile{
		Version:  "1.0",
		Rules:    make(map[string]RuleLock),
		FilePath: lockFilePath,
	}
}

// LoadLockFile loads a lock file from disk
func LoadLockFile(lockFilePath string) (*LockFile, error) {
	// If lock file doesn't exist, return new empty lock file
	if _, err := os.Stat(lockFilePath); os.IsNotExist(err) {
		return NewLockFile(lockFilePath), nil
	}

	data, err := os.ReadFile(lockFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	var lf LockFile
	if err := json.Unmarshal(data, &lf); err != nil {
		return nil, fmt.Errorf("failed to parse lock file: %w", err)
	}

	lf.FilePath = lockFilePath

	// Initialize Rules map if nil
	if lf.Rules == nil {
		lf.Rules = make(map[string]RuleLock)
	}

	return &lf, nil
}

// Save writes the lock file to disk
func (lf *LockFile) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(lf.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal lock file: %w", err)
	}

	if err := os.WriteFile(lf.FilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write lock file: %w", err)
	}

	return nil
}

// AddRule adds or updates a rule in the lock file
func (lf *LockFile) AddRule(lock RuleLock) {
	// Check if rule already exists
	if existing, ok := lf.Rules[lock.Name]; ok {
		lock.InstalledAt = existing.InstalledAt
		lock.UpdatedAt = time.Now()
	} else {
		lock.InstalledAt = time.Now()
	}

	lf.Rules[lock.Name] = lock
}

// RemoveRule removes a rule from the lock file
func (lf *LockFile) RemoveRule(ruleName string) {
	delete(lf.Rules, ruleName)
}

// GetRule retrieves a rule lock entry
func (lf *LockFile) GetRule(ruleName string) (*RuleLock, bool) {
	lock, ok := lf.Rules[ruleName]
	return &lock, ok
}

// HasRule checks if a rule is in the lock file
func (lf *LockFile) HasRule(ruleName string) bool {
	_, ok := lf.Rules[ruleName]
	return ok
}

// ListRules returns all rule names in the lock file
func (lf *LockFile) ListRules() []string {
	var names []string
	for name := range lf.Rules {
		names = append(names, name)
	}
	return names
}
