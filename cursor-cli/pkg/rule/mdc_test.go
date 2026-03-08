package rule

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMDC(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectError   bool
		expectedDesc  string
		expectedGlobs []string
		expectedApply bool
	}{
		{
			name: "basic valid mdc",
			input: `---
description: "Test rule"
globs:
  - "**/*.ts"
alwaysApply: false
---

# Test Content`,
			expectError:   false,
			expectedDesc:  "Test rule",
			expectedGlobs: []string{"**/*.ts"},
			expectedApply: false,
		},
		{
			name: "always apply rule",
			input: `---
description: "Global rule"
globs: []
alwaysApply: true
---

# Global Content`,
			expectError:   false,
			expectedDesc:  "Global rule",
			expectedGlobs: []string{},
			expectedApply: true,
		},
		{
			name: "with tags and version",
			input: `---
description: "Advanced rule"
globs:
  - "**/*.js"
  - "**/*.jsx"
alwaysApply: false
tags:
  - javascript
  - react
version: "1.2.3"
---

# Advanced Content`,
			expectError:   false,
			expectedDesc:  "Advanced rule",
			expectedGlobs: []string{"**/*.js", "**/*.jsx"},
			expectedApply: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, err := ParseMDC([]byte(tt.input), "test.mdc")

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedDesc, rule.Frontmatter.Description)
			assert.Equal(t, tt.expectedGlobs, rule.Frontmatter.Globs)
			assert.Equal(t, tt.expectedApply, rule.Frontmatter.AlwaysApply)
			assert.NotEmpty(t, rule.Content)
		})
	}
}

func TestParseMDCFile(t *testing.T) {
	testdataDir := "../../testdata"

	t.Run("parse valid basic file", func(t *testing.T) {
		filePath := filepath.Join(testdataDir, "valid", "basic.mdc")
		rule, err := ParseMDCFile(filePath)

		require.NoError(t, err)
		assert.Equal(t, "Basic TypeScript coding standards", rule.Frontmatter.Description)
		assert.Contains(t, rule.Frontmatter.Globs, "**/*.ts")
		assert.Contains(t, rule.Frontmatter.Globs, "**/*.tsx")
		assert.False(t, rule.Frontmatter.AlwaysApply)
		assert.Contains(t, rule.Frontmatter.Tags, "typescript")
		assert.Equal(t, "1.0.0", rule.Frontmatter.Version)
		assert.Contains(t, rule.Content, "TypeScript Coding Standards")
	})

	t.Run("parse always-apply file", func(t *testing.T) {
		filePath := filepath.Join(testdataDir, "valid", "always-apply.mdc")
		rule, err := ParseMDCFile(filePath)

		require.NoError(t, err)
		assert.Equal(t, "Always-on global coding conventions", rule.Frontmatter.Description)
		assert.Empty(t, rule.Frontmatter.Globs)
		assert.True(t, rule.Frontmatter.AlwaysApply)
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := ParseMDCFile("non-existent.mdc")
		assert.Error(t, err)
	})
}

func TestMDCMarshal(t *testing.T) {
	rule := &MDCRule{
		Frontmatter: MDCFrontmatter{
			Description: "Test marshal",
			Globs:       []string{"**/*.go"},
			AlwaysApply: false,
			Tags:        []string{"golang"},
			Version:     "1.0.0",
		},
		Content:  "# Test Content\n\nThis is a test.",
		FilePath: "test.mdc",
	}

	data, err := rule.Marshal()
	require.NoError(t, err)

	// Parse it back
	parsedRule, err := ParseMDC(data, "test.mdc")
	require.NoError(t, err)

	assert.Equal(t, rule.Frontmatter.Description, parsedRule.Frontmatter.Description)
	assert.Equal(t, rule.Frontmatter.Globs, parsedRule.Frontmatter.Globs)
	assert.Equal(t, rule.Frontmatter.AlwaysApply, parsedRule.Frontmatter.AlwaysApply)
	assert.Equal(t, rule.Frontmatter.Tags, parsedRule.Frontmatter.Tags)
	assert.Equal(t, rule.Frontmatter.Version, parsedRule.Frontmatter.Version)
}

func TestMDCWriteTo(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test-rule.mdc")

	rule := &MDCRule{
		Frontmatter: MDCFrontmatter{
			Description: "Test write",
			Globs:       []string{"**/*.test.ts"},
			AlwaysApply: false,
		},
		Content: "# Test\n\nContent here.",
	}

	err := rule.WriteTo(filePath)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(filePath)
	require.NoError(t, err)

	// Parse it back
	parsedRule, err := ParseMDCFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, rule.Frontmatter.Description, parsedRule.Frontmatter.Description)
	assert.Equal(t, rule.Frontmatter.Globs, parsedRule.Frontmatter.Globs)
}

func TestMDCName(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{"simple name", "/path/to/typescript.mdc", "typescript"},
		{"with dashes", "/path/to/react-hooks.mdc", "react-hooks"},
		{"nested path", "/a/b/c/my-rule.mdc", "my-rule"},
		{"empty path", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &MDCRule{FilePath: tt.filePath}
			assert.Equal(t, tt.expected, rule.Name())
		})
	}
}

func TestMDCHasGlob(t *testing.T) {
	rule := &MDCRule{
		Frontmatter: MDCFrontmatter{
			Globs: []string{"**/*.ts", "**/*.tsx", "src/**/*.js"},
		},
	}

	assert.True(t, rule.HasGlob("**/*.ts"))
	assert.True(t, rule.HasGlob("**/*.tsx"))
	assert.True(t, rule.HasGlob("src/**/*.js"))
	assert.False(t, rule.HasGlob("**/*.go"))
	assert.False(t, rule.HasGlob("*.ts"))
}
