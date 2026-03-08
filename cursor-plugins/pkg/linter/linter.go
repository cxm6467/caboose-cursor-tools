package linter

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Config represents the linter configuration for a project
type Config struct {
	Languages   map[string]LanguageConfig `json:"languages"`
	ToolLinters []ToolLinter              `json:"tool_linters,omitempty"`
}

// LanguageConfig holds linter configuration for a specific language
type LanguageConfig struct {
	Enabled bool             `json:"enabled"`
	Linters []LinterSettings `json:"linters"`
}

// LinterSettings defines a single linter's configuration
type LinterSettings struct {
	Name    string `json:"name"`
	Command string `json:"command"`
	Enabled bool   `json:"enabled"`
	Timeout string `json:"timeout,omitempty"` // e.g., "5s"
}

// ToolLinter represents project-wide linters (not per-file)
type ToolLinter struct {
	Name    string `json:"name"`
	Command string `json:"command"`
	Enabled bool   `json:"enabled"`
}

// LanguageDetector maps file extensions to languages
type LanguageDetector struct{}

// Detect determines the language of a file based on extension or shebang
func (ld *LanguageDetector) Detect(filePath string) string {
	ext := filepath.Ext(filePath)
	base := filepath.Base(filePath)

	// Extension-based detection
	switch strings.TrimPrefix(ext, ".") {
	case "rb", "rake", "gemspec":
		return "ruby"
	case "js", "jsx", "mjs", "cjs":
		return "javascript"
	case "ts", "tsx", "mts", "cts":
		return "javascript"
	case "py", "pyi":
		return "python"
	case "go":
		return "go"
	case "rs":
		return "rust"
	case "md", "mdx":
		return "markdown"
	case "html", "htm", "css":
		return "html"
	case "sh", "bash", "zsh":
		return "shell"
	}

	// Check shebang for extensionless files
	if ext == "" || ext == base {
		shebang, err := getShebang(filePath)
		if err == nil {
			if strings.Contains(shebang, "bash") || strings.Contains(shebang, "/sh") {
				return "shell"
			}
			if strings.Contains(shebang, "zsh") {
				return "shell"
			}
			if strings.Contains(shebang, "ruby") {
				return "ruby"
			}
			if strings.Contains(shebang, "python") {
				return "python"
			}
			if strings.Contains(shebang, "node") {
				return "javascript"
			}
		}
	}

	return ""
}

// getShebang reads the first line of a file
func getShebang(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 {
		return lines[0], nil
	}

	return "", fmt.Errorf("empty file")
}

// Runner runs linters on files
type Runner struct {
	config      *Config
	projectRoot string
}

// NewRunner creates a new linter runner
func NewRunner(projectRoot string) (*Runner, error) {
	config, err := loadConfig(projectRoot)
	if err != nil {
		return nil, err
	}

	return &Runner{
		config:      config,
		projectRoot: projectRoot,
	}, nil
}

// Run executes linters for a given file
func (r *Runner) Run(filePath string) ([]LintResult, error) {
	// Detect language
	detector := &LanguageDetector{}
	lang := detector.Detect(filePath)
	if lang == "" {
		return nil, nil // No linters for unknown file types
	}

	// Get language config
	langConfig, ok := r.config.Languages[lang]
	if !ok || !langConfig.Enabled {
		return nil, nil // Language not configured or disabled
	}

	var results []LintResult

	// Run each enabled linter
	for _, linter := range langConfig.Linters {
		if !linter.Enabled {
			continue
		}

		result, err := r.runLinter(linter, filePath)
		if err != nil {
			results = append(results, LintResult{
				Linter:  linter.Name,
				Success: false,
				Output:  fmt.Sprintf("Failed to run: %v", err),
			})
			continue
		}

		results = append(results, *result)
	}

	return results, nil
}

// runLinter executes a single linter
func (r *Runner) runLinter(linter LinterSettings, filePath string) (*LintResult, error) {
	// Replace $1 or {file} with actual file path
	cmd := strings.ReplaceAll(linter.Command, "$1", filePath)
	cmd = strings.ReplaceAll(cmd, "{file}", filePath)

	// Execute command
	execCmd := exec.Command("bash", "-c", cmd)
	execCmd.Dir = r.projectRoot

	output, err := execCmd.CombinedOutput()

	result := &LintResult{
		Linter:  linter.Name,
		Success: err == nil,
		Output:  string(output),
	}

	return result, nil
}

// LintResult represents the result of running a linter
type LintResult struct {
	Linter  string `json:"linter"`
	Success bool   `json:"success"`
	Output  string `json:"output"`
}

// loadConfig loads the linter configuration for a project
func loadConfig(projectRoot string) (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Compute project hash (same logic as claude-plugins)
	projectHash := strings.ReplaceAll(projectRoot, "/", "-")
	projectHash = strings.TrimPrefix(projectHash, "-")

	configPath := filepath.Join(home, ".cursor", "code-lint", projectHash, "config.json")

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("no config found at %s, run setup-lint first", configPath)
	}

	// Read and parse config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the linter configuration
func SaveConfig(projectRoot string, config *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	projectHash := strings.ReplaceAll(projectRoot, "/", "-")
	projectHash = strings.TrimPrefix(projectHash, "-")

	configDir := filepath.Join(home, ".cursor", "code-lint", projectHash)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
