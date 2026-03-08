package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/cxm6467/caboose-ai/cursor-plugins/pkg/linter"
)

// HookInput represents the input from Cursor hooks
type HookInput struct {
	ToolName  string                 `json:"tool_name"`
	ToolInput map[string]interface{} `json:"tool_input"`
}

func main() {
	// Read hook input from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read input: %v\n", err)
		os.Exit(0) // Silent exit on error (hook convention)
	}

	var hookInput HookInput
	if err := json.Unmarshal(input, &hookInput); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse input: %v\n", err)
		os.Exit(0)
	}

	// Extract file path
	filePath, ok := hookInput.ToolInput["file_path"].(string)
	if !ok || filePath == "" {
		os.Exit(0) // No file path, silent exit
	}

	// Make absolute path
	if filePath[0] != '/' {
		cwd, _ := os.Getwd()
		filePath = fmt.Sprintf("%s/%s", cwd, filePath)
	}

	// Check file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Exit(0)
	}

	// Determine project root
	projectRoot := os.Getenv("CURSOR_PROJECT_DIR")
	if projectRoot == "" {
		// Fall back to git root or file directory
		projectRoot = findProjectRoot(filePath)
	}

	// Create linter runner
	runner, err := linter.NewRunner(projectRoot)
	if err != nil {
		// No config = silent exit (no-op)
		os.Exit(0)
	}

	// Run linters
	results, err := runner.Run(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Linter error: %v\n", err)
		os.Exit(0)
	}

	// Output results
	if len(results) > 0 {
		for _, result := range results {
			if !result.Success && result.Output != "" {
				fmt.Printf("=== %s ===\n%s\n", result.Linter, result.Output)
			}
		}
	}
}

// findProjectRoot attempts to find the project root directory
func findProjectRoot(filePath string) string {
	// Try git root
	// This is a simplified version - could use exec.Command("git", "rev-parse", "--show-toplevel")
	dir := filePath
	for {
		parent := dir + "/.."
		if _, err := os.Stat(parent + "/.git"); err == nil {
			return parent
		}
		if parent == "/" {
			break
		}
		dir = parent
	}

	// Fall back to file directory
	return filePath
}
