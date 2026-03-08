package rule

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateData holds the data for rendering a rule template
type TemplateData struct {
	Name        string
	Description string
	Globs       []string
	Tags        []string
	Author      string
}

// TemplateType represents the type of template
type TemplateType string

const (
	TemplateBasic        TemplateType = "basic"
	TemplateGlobBased    TemplateType = "glob-based"
	TemplateAlwaysApply  TemplateType = "always-apply"
)

// Generator generates rules from templates
type Generator struct {
	TemplateDir string
}

// NewGenerator creates a new template generator
func NewGenerator(templateDir string) *Generator {
	return &Generator{
		TemplateDir: templateDir,
	}
}

// Generate creates a new rule from a template
func (g *Generator) Generate(templateType TemplateType, data TemplateData, outputPath string) error {
	// Get template content
	templateContent, err := g.getTemplate(templateType)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	// Parse template
	tmpl, err := template.New(string(templateType)).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Ensure output directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// getTemplate retrieves template content
func (g *Generator) getTemplate(templateType TemplateType) (string, error) {
	templateName := string(templateType) + ".mdc.tmpl"

	// Try custom template directory first
	if g.TemplateDir != "" {
		externalPath := filepath.Join(g.TemplateDir, templateName)
		data, err := os.ReadFile(externalPath)
		if err == nil {
			return string(data), nil
		}
	}

	// Fall back to default template directory (relative to project root)
	// Assumes templates/ is in the project root
	defaultPaths := []string{
		filepath.Join("templates", templateName),
		filepath.Join("..", "..", "templates", templateName), // From pkg/rule
		filepath.Join("../../templates", templateName),       // Alternative
	}

	for _, path := range defaultPaths {
		data, err := os.ReadFile(path)
		if err == nil {
			return string(data), nil
		}
	}

	return "", fmt.Errorf("template %s not found", templateName)
}

// GenerateDefault creates a rule with sensible defaults
func GenerateDefault(name, description string, globs []string, outputPath string) error {
	generator := NewGenerator("")

	data := TemplateData{
		Name:        FormatName(name),
		Description: description,
		Globs:       globs,
		Tags:        inferTags(name, globs),
		Author:      GetCurrentUser(),
	}

	// Choose template type based on input
	templateType := TemplateBasic
	if len(globs) == 0 {
		templateType = TemplateAlwaysApply
	} else if len(globs) > 0 {
		templateType = TemplateGlobBased
	}

	return generator.Generate(templateType, data, outputPath)
}

// FormatName converts a rule name to proper title case
func FormatName(name string) string {
	// Replace dashes and underscores with spaces
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")

	// Title case
	words := strings.Fields(name)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

// inferTags attempts to infer tags from the rule name and globs
func inferTags(name string, globs []string) []string {
	tags := make(map[string]bool)

	// Extract from name
	nameParts := strings.FieldsFunc(strings.ToLower(name), func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})
	for _, part := range nameParts {
		tags[part] = true
	}

	// Extract from globs
	for _, glob := range globs {
		// Extract file extensions
		if strings.Contains(glob, ".") {
			parts := strings.Split(glob, ".")
			if len(parts) > 0 {
				ext := parts[len(parts)-1]
				// Remove trailing patterns like * or }
				ext = strings.TrimRight(ext, "*}")
				if ext != "" && len(ext) < 10 {
					tags[ext] = true
				}
			}
		}
	}

	// Convert map to slice
	result := make([]string, 0, len(tags))
	for tag := range tags {
		if len(tag) > 1 { // Skip single-char tags
			result = append(result, tag)
		}
	}

	return result
}

// GetCurrentUser returns the current user's name or a default
func GetCurrentUser() string {
	if user := os.Getenv("USER"); user != "" {
		return user
	}
	if user := os.Getenv("USERNAME"); user != "" {
		return user
	}
	return "Unknown"
}
