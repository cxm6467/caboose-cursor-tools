package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caboose/cursor-cli/pkg/rule"
)

// Importer handles importing rules from various sources
type Importer struct {
	Verbose bool
}

// NewImporter creates a new importer
func NewImporter(verbose bool) *Importer {
	return &Importer{
		Verbose: verbose,
	}
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	Source         string
	ImportedRules  []string
	SkippedRules   []string
	Errors         []error
	TotalProcessed int
}

// MigrateLegacyCursorrules converts a legacy .cursorrules file to .mdc format
func (imp *Importer) MigrateLegacyCursorrules(inputPath, outputPath string) (*rule.MDCRule, error) {
	// Read the legacy file
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse the content - legacy files are just markdown
	contentStr := string(content)

	// Try to infer metadata from the content
	description := inferDescription(contentStr)
	tags := inferTags(contentStr)

	// Create MDC rule
	mdcRule := &rule.MDCRule{
		Frontmatter: rule.MDCFrontmatter{
			Description: description,
			AlwaysApply: true, // Legacy files typically apply to all files
			Tags:        tags,
			Version:     "1.0.0",
			Author:      rule.GetCurrentUser(),
		},
		Content:  contentStr,
		FilePath: outputPath,
	}

	return mdcRule, nil
}

// MigrateDirectory migrates all .cursorrules files in a directory
func (imp *Importer) MigrateDirectory(inputDir, outputDir string) (*ImportResult, error) {
	result := &ImportResult{
		Source: inputDir,
	}

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Only process .cursorrules files
		if filepath.Ext(path) != ".cursorrules" && filepath.Base(path) != ".cursorrules" {
			return nil
		}

		result.TotalProcessed++

		// Determine output filename
		baseName := strings.TrimSuffix(filepath.Base(path), ".cursorrules")
		if baseName == "" {
			baseName = "migrated-rule"
		}
		outputPath := filepath.Join(outputDir, baseName+".mdc")

		// Migrate the file
		mdcRule, err := imp.MigrateLegacyCursorrules(path, outputPath)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to migrate %s: %w", path, err))
			result.SkippedRules = append(result.SkippedRules, path)
			return nil
		}

		// Write the migrated rule
		if err := mdcRule.WriteTo(outputPath); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to write %s: %w", outputPath, err))
			result.SkippedRules = append(result.SkippedRules, path)
			return nil
		}

		result.ImportedRules = append(result.ImportedRules, outputPath)

		if imp.Verbose {
			fmt.Printf("Migrated: %s -> %s\n", path, outputPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return result, nil
}

// inferDescription tries to extract a description from the content
func inferDescription(content string) string {
	lines := strings.Split(content, "\n")

	// Look for the first heading or meaningful line
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// If it's a heading, use it
		if strings.HasPrefix(line, "#") {
			desc := strings.TrimSpace(strings.TrimLeft(line, "#"))
			if len(desc) > 0 {
				return desc
			}
		}

		// Otherwise, use the first non-empty line (truncated)
		if len(line) > 0 {
			if len(line) > 100 {
				return line[:100] + "..."
			}
			return line
		}
	}

	return "Migrated cursor rules"
}

// inferTags tries to extract tags from the content
func inferTags(content string) []string {
	tags := make(map[string]bool)

	// Common programming language keywords
	keywords := []string{
		"typescript", "javascript", "python", "go", "rust", "java", "c++", "cpp",
		"react", "vue", "angular", "node", "deno", "bun",
		"test", "testing", "security", "performance", "style", "lint",
	}

	contentLower := strings.ToLower(content)
	for _, keyword := range keywords {
		if strings.Contains(contentLower, keyword) {
			tags[keyword] = true
		}
	}

	// Convert to slice
	result := make([]string, 0, len(tags))
	for tag := range tags {
		result = append(result, tag)
	}

	// Add default if no tags found
	if len(result) == 0 {
		result = []string{"general"}
	}

	return result
}
