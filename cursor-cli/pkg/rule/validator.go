package rule

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// ValidationError represents a validation error with severity and location
type ValidationError struct {
	Severity string // "error" or "warning"
	Field    string
	Message  string
	FilePath string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("[%s] %s: %s (file: %s)", e.Severity, e.Field, e.Message, e.FilePath)
}

// ValidationResult contains the result of validating an MDC rule
type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

// AddError adds an error-level validation issue
func (vr *ValidationResult) AddError(field, message, filePath string) {
	vr.Valid = false
	vr.Errors = append(vr.Errors, ValidationError{
		Severity: "error",
		Field:    field,
		Message:  message,
		FilePath: filePath,
	})
}

// AddWarning adds a warning-level validation issue
func (vr *ValidationResult) AddWarning(field, message, filePath string) {
	vr.Errors = append(vr.Errors, ValidationError{
		Severity: "warning",
		Field:    field,
		Message:  message,
		FilePath: filePath,
	})
}

// HasErrors returns true if there are any error-level issues
func (vr *ValidationResult) HasErrors() bool {
	for _, err := range vr.Errors {
		if err.Severity == "error" {
			return true
		}
	}
	return false
}

// HasWarnings returns true if there are any warning-level issues
func (vr *ValidationResult) HasWarnings() bool {
	for _, err := range vr.Errors {
		if err.Severity == "warning" {
			return true
		}
	}
	return false
}

// Validator validates MDC rules
type Validator struct {
	StrictMode bool
}

// NewValidator creates a new validator
func NewValidator(strictMode bool) *Validator {
	return &Validator{
		StrictMode: strictMode,
	}
}

// ValidateFile validates a .mdc file by path
func (v *Validator) ValidateFile(filePath string) (*ValidationResult, error) {
	rule, err := ParseMDCFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return v.Validate(rule), nil
}

// Validate validates an MDC rule
func (v *Validator) Validate(rule *MDCRule) *ValidationResult {
	result := &ValidationResult{Valid: true}

	v.validateFrontmatter(rule, result)
	v.validateContent(rule, result)
	v.validateGlobs(rule, result)

	if result.HasErrors() {
		result.Valid = false
	}

	return result
}

// validateFrontmatter validates the YAML frontmatter
func (v *Validator) validateFrontmatter(rule *MDCRule, result *ValidationResult) {
	fm := rule.Frontmatter

	// Required fields
	if fm.Description == "" {
		result.AddError("description", "description is required", rule.FilePath)
	}

	if len(fm.Globs) == 0 && !fm.AlwaysApply {
		result.AddError("globs", "globs is required when alwaysApply is false", rule.FilePath)
	}

	// Validate description length
	if len(fm.Description) > 200 {
		result.AddWarning("description", "description is very long (>200 chars), consider shortening", rule.FilePath)
	}

	// Validate description quality
	if len(fm.Description) > 0 && len(fm.Description) < 10 {
		result.AddWarning("description", "description is very short (<10 chars), consider adding more detail", rule.FilePath)
	}

	// Validate version format (if provided)
	if fm.Version != "" && !isValidSemver(fm.Version) {
		result.AddWarning("version", fmt.Sprintf("version '%s' does not follow semver format", fm.Version), rule.FilePath)
	}
}

// validateContent validates the markdown content
func (v *Validator) validateContent(rule *MDCRule, result *ValidationResult) {
	content := strings.TrimSpace(rule.Content)

	if content == "" {
		result.AddError("content", "markdown content is empty", rule.FilePath)
		return
	}

	// Check content length
	if len(content) < 50 && v.StrictMode {
		result.AddWarning("content", "content is very short (<50 chars)", rule.FilePath)
	}

	// Check for proper markdown structure
	if !strings.Contains(content, "#") && v.StrictMode {
		result.AddWarning("content", "content has no markdown headers", rule.FilePath)
	}
}

// validateGlobs validates glob patterns
func (v *Validator) validateGlobs(rule *MDCRule, result *ValidationResult) {
	for i, glob := range rule.Frontmatter.Globs {
		// Check if glob pattern is valid
		if glob == "" {
			result.AddError("globs", fmt.Sprintf("glob pattern at index %d is empty", i), rule.FilePath)
			continue
		}

		// Test if glob pattern is valid by trying to match
		_, err := doublestar.Match(glob, "test.txt")
		if err != nil {
			result.AddError("globs", fmt.Sprintf("invalid glob pattern '%s': %v", glob, err), rule.FilePath)
		}

		// Warning for overly broad patterns
		if glob == "**/*" || glob == "*" {
			result.AddWarning("globs", fmt.Sprintf("glob pattern '%s' is very broad and may match too many files", glob), rule.FilePath)
		}

		// Warning for absolute paths (should be relative)
		if filepath.IsAbs(glob) {
			result.AddWarning("globs", fmt.Sprintf("glob pattern '%s' uses absolute path, use relative paths instead", glob), rule.FilePath)
		}
	}
}

// ValidateDirectory validates all .mdc files in a directory
func (v *Validator) ValidateDirectory(dirPath string) (map[string]*ValidationResult, error) {
	results := make(map[string]*ValidationResult)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .mdc files
		if filepath.Ext(path) != ".mdc" {
			return nil
		}

		result, err := v.ValidateFile(path)
		if err != nil {
			return fmt.Errorf("failed to validate %s: %w", path, err)
		}

		results[path] = result
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

// isValidSemver checks if a version string follows semantic versioning
// Simple check for X.Y.Z format
func isValidSemver(version string) bool {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return false
	}

	for _, part := range parts {
		if part == "" {
			return false
		}

		// Check if all characters are digits
		for _, c := range part {
			if c < '0' || c > '9' {
				return false
			}
		}
	}

	return true
}
