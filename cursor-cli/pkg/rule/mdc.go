package rule

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

// MDCFrontmatter represents the YAML frontmatter in a .mdc file
type MDCFrontmatter struct {
	Description string   `yaml:"description"`
	Globs       []string `yaml:"globs"`
	AlwaysApply bool     `yaml:"alwaysApply"`
	Tags        []string `yaml:"tags,omitempty"`
	Version     string   `yaml:"version,omitempty"`
	Author      string   `yaml:"author,omitempty"`
}

// MDCRule represents a complete .mdc rule file
type MDCRule struct {
	Frontmatter MDCFrontmatter
	Content     string
	FilePath    string
}

// ParseMDCFile parses a .mdc file and returns an MDCRule
func ParseMDCFile(filePath string) (*MDCRule, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return ParseMDC(data, filePath)
}

// ParseMDC parses MDC content (YAML frontmatter + markdown) from bytes
func ParseMDC(data []byte, filePath string) (*MDCRule, error) {
	var fm MDCFrontmatter
	var content string

	// Use frontmatter library to parse YAML frontmatter
	rest, err := frontmatter.Parse(bytes.NewReader(data), &fm)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	content = string(rest)

	rule := &MDCRule{
		Frontmatter: fm,
		Content:     content,
		FilePath:    filePath,
	}

	return rule, nil
}

// Marshal converts an MDCRule back to .mdc format (YAML frontmatter + markdown)
func (r *MDCRule) Marshal() ([]byte, error) {
	var buf bytes.Buffer

	// Write YAML frontmatter
	buf.WriteString("---\n")

	fmBytes, err := yaml.Marshal(&r.Frontmatter)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	buf.Write(fmBytes)
	buf.WriteString("---\n\n")

	// Write markdown content
	buf.WriteString(r.Content)

	return buf.Bytes(), nil
}

// WriteTo writes the MDC rule to a file
func (r *MDCRule) WriteTo(filePath string) error {
	data, err := r.Marshal()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	r.FilePath = filePath
	return nil
}

// Name returns the rule name derived from the file path
func (r *MDCRule) Name() string {
	if r.FilePath == "" {
		return ""
	}

	base := filepath.Base(r.FilePath)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// HasGlob checks if the rule matches a given file pattern
func (r *MDCRule) HasGlob(pattern string) bool {
	for _, glob := range r.Frontmatter.Globs {
		if glob == pattern {
			return true
		}
	}
	return false
}

// String returns a human-readable representation of the rule
func (r *MDCRule) String() string {
	return fmt.Sprintf("MDCRule{Name: %s, Description: %s, Globs: %v, AlwaysApply: %v}",
		r.Name(),
		r.Frontmatter.Description,
		r.Frontmatter.Globs,
		r.Frontmatter.AlwaysApply,
	)
}
