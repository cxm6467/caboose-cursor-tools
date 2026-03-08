package rule

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidator_ValidateFile(t *testing.T) {
	testdataDir := "../../testdata"
	validator := NewValidator(false)

	t.Run("valid basic file", func(t *testing.T) {
		filePath := filepath.Join(testdataDir, "valid", "basic.mdc")
		result, err := validator.ValidateFile(filePath)

		require.NoError(t, err)
		assert.True(t, result.Valid)
		assert.False(t, result.HasErrors())
	})

	t.Run("valid always-apply file", func(t *testing.T) {
		filePath := filepath.Join(testdataDir, "valid", "always-apply.mdc")
		result, err := validator.ValidateFile(filePath)

		require.NoError(t, err)
		assert.True(t, result.Valid)
		assert.False(t, result.HasErrors())
	})

	t.Run("missing description", func(t *testing.T) {
		filePath := filepath.Join(testdataDir, "invalid", "missing-description.mdc")
		result, err := validator.ValidateFile(filePath)

		require.NoError(t, err)
		assert.False(t, result.Valid)
		assert.True(t, result.HasErrors())

		// Should have error about missing description
		foundError := false
		for _, err := range result.Errors {
			if err.Field == "description" && err.Severity == "error" {
				foundError = true
				break
			}
		}
		assert.True(t, foundError, "Expected error for missing description")
	})

	t.Run("empty content", func(t *testing.T) {
		filePath := filepath.Join(testdataDir, "invalid", "empty-content.mdc")
		result, err := validator.ValidateFile(filePath)

		require.NoError(t, err)
		assert.False(t, result.Valid)
		assert.True(t, result.HasErrors())

		// Should have error about empty content
		foundError := false
		for _, err := range result.Errors {
			if err.Field == "content" && err.Severity == "error" {
				foundError = true
				break
			}
		}
		assert.True(t, foundError, "Expected error for empty content")
	})
}

func TestValidator_Validate(t *testing.T) {
	validator := NewValidator(false)

	t.Run("valid rule", func(t *testing.T) {
		rule := &MDCRule{
			Frontmatter: MDCFrontmatter{
				Description: "Valid rule",
				Globs:       []string{"**/*.ts"},
				AlwaysApply: false,
			},
			Content:  "# Valid Content\n\nThis is a valid rule with enough content.",
			FilePath: "test.mdc",
		}

		result := validator.Validate(rule)
		assert.True(t, result.Valid)
		assert.False(t, result.HasErrors())
	})

	t.Run("missing description", func(t *testing.T) {
		rule := &MDCRule{
			Frontmatter: MDCFrontmatter{
				Description: "",
				Globs:       []string{"**/*.ts"},
				AlwaysApply: false,
			},
			Content:  "# Content",
			FilePath: "test.mdc",
		}

		result := validator.Validate(rule)
		assert.False(t, result.Valid)
		assert.True(t, result.HasErrors())
	})

	t.Run("missing globs when not always apply", func(t *testing.T) {
		rule := &MDCRule{
			Frontmatter: MDCFrontmatter{
				Description: "Test",
				Globs:       []string{},
				AlwaysApply: false,
			},
			Content:  "# Content",
			FilePath: "test.mdc",
		}

		result := validator.Validate(rule)
		assert.False(t, result.Valid)
		assert.True(t, result.HasErrors())
	})

	t.Run("always apply without globs is valid", func(t *testing.T) {
		rule := &MDCRule{
			Frontmatter: MDCFrontmatter{
				Description: "Global rule",
				Globs:       []string{},
				AlwaysApply: true,
			},
			Content:  "# Global content here",
			FilePath: "test.mdc",
		}

		result := validator.Validate(rule)
		assert.True(t, result.Valid)
		assert.False(t, result.HasErrors())
	})

	t.Run("empty content", func(t *testing.T) {
		rule := &MDCRule{
			Frontmatter: MDCFrontmatter{
				Description: "Test",
				Globs:       []string{"**/*.ts"},
				AlwaysApply: false,
			},
			Content:  "",
			FilePath: "test.mdc",
		}

		result := validator.Validate(rule)
		assert.False(t, result.Valid)
		assert.True(t, result.HasErrors())
	})
}

func TestValidator_ValidateGlobs(t *testing.T) {
	validator := NewValidator(false)

	tests := []struct {
		name         string
		globs        []string
		expectErrors bool
		expectWarns  bool
	}{
		{
			name:         "valid globs",
			globs:        []string{"**/*.ts", "src/**/*.js"},
			expectErrors: false,
			expectWarns:  false,
		},
		{
			name:         "empty glob",
			globs:        []string{""},
			expectErrors: true,
			expectWarns:  false,
		},
		{
			name:         "overly broad pattern",
			globs:        []string{"**/*"},
			expectErrors: false,
			expectWarns:  true,
		},
		{
			name:         "absolute path warning",
			globs:        []string{"/absolute/path/*.ts"},
			expectErrors: false,
			expectWarns:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &MDCRule{
				Frontmatter: MDCFrontmatter{
					Description: "Test",
					Globs:       tt.globs,
					AlwaysApply: false,
				},
				Content:  "# Content",
				FilePath: "test.mdc",
			}

			result := validator.Validate(rule)

			if tt.expectErrors {
				assert.True(t, result.HasErrors(), "Expected errors for globs: %v", tt.globs)
			} else {
				assert.False(t, result.HasErrors(), "Did not expect errors for globs: %v", tt.globs)
			}

			if tt.expectWarns {
				assert.True(t, result.HasWarnings(), "Expected warnings for globs: %v", tt.globs)
			}
		})
	}
}

func TestValidator_ValidateDirectory(t *testing.T) {
	testdataDir := "../../testdata"
	validator := NewValidator(false)

	t.Run("validate valid directory", func(t *testing.T) {
		validDir := filepath.Join(testdataDir, "valid")
		results, err := validator.ValidateDirectory(validDir)

		require.NoError(t, err)
		assert.NotEmpty(t, results)

		// All valid files should pass
		for path, result := range results {
			assert.True(t, result.Valid, "File %s should be valid", path)
		}
	})

	t.Run("validate invalid directory", func(t *testing.T) {
		invalidDir := filepath.Join(testdataDir, "invalid")
		results, err := validator.ValidateDirectory(invalidDir)

		require.NoError(t, err)
		assert.NotEmpty(t, results)

		// At least some invalid files should have errors
		hasErrors := false
		for _, result := range results {
			if result.HasErrors() {
				hasErrors = true
				break
			}
		}
		assert.True(t, hasErrors, "Expected at least one file with errors in invalid directory")
	})
}

func TestIsValidSemver(t *testing.T) {
	tests := []struct {
		version string
		valid   bool
	}{
		{"1.0.0", true},
		{"0.1.0", true},
		{"10.20.30", true},
		{"1.0", false},
		{"1", false},
		{"1.0.0.0", false},
		{"v1.0.0", false},
		{"1.0.a", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			result := isValidSemver(tt.version)
			assert.Equal(t, tt.valid, result, "Version: %s", tt.version)
		})
	}
}

func TestValidationResult(t *testing.T) {
	t.Run("add error", func(t *testing.T) {
		vr := &ValidationResult{Valid: true}
		vr.AddError("field", "message", "path")

		assert.False(t, vr.Valid)
		assert.True(t, vr.HasErrors())
		assert.Len(t, vr.Errors, 1)
		assert.Equal(t, "error", vr.Errors[0].Severity)
	})

	t.Run("add warning", func(t *testing.T) {
		vr := &ValidationResult{Valid: true}
		vr.AddWarning("field", "message", "path")

		assert.True(t, vr.Valid) // Warnings don't invalidate
		assert.False(t, vr.HasErrors())
		assert.True(t, vr.HasWarnings())
		assert.Len(t, vr.Errors, 1) // Stored in same array
		assert.Equal(t, "warning", vr.Errors[0].Severity)
	})
}
