package marketplace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeRegistryURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "GitHub repo URL",
			input:    "https://github.com/user/repo",
			expected: "https://raw.githubusercontent.com/user/repo/main/index.json",
		},
		{
			name:     "GitHub repo with trailing slash",
			input:    "https://github.com/user/repo/",
			expected: "https://raw.githubusercontent.com/user/repo/main/index.json",
		},
		{
			name:     "GitHub repo with branch",
			input:    "https://github.com/user/repo/tree/develop",
			expected: "https://raw.githubusercontent.com/user/repo/develop/index.json",
		},
		{
			name:     "Direct index.json URL",
			input:    "https://example.com/index.json",
			expected: "https://example.com/index.json",
		},
		{
			name:     "Non-GitHub URL",
			input:    "https://example.com/marketplace",
			expected: "https://example.com/marketplace/index.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeRegistryURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchRules(t *testing.T) {
	client := NewClient()

	index := &RegistryIndex{
		Name: "test-registry",
		Rules: []RegistryRule{
			{
				Name:        "typescript-standards",
				Description: "TypeScript coding standards",
				Tags:        []string{"typescript", "standards"},
			},
			{
				Name:        "react-best-practices",
				Description: "React best practices and patterns",
				Tags:        []string{"react", "javascript"},
			},
			{
				Name:        "python-style",
				Description: "Python style guide",
				Tags:        []string{"python", "style"},
			},
		},
	}

	tests := []struct {
		name          string
		query         string
		expectedCount int
		expectedNames []string
	}{
		{
			name:          "Search by name",
			query:         "typescript",
			expectedCount: 1,
			expectedNames: []string{"typescript-standards"},
		},
		{
			name:          "Search by description",
			query:         "best practices",
			expectedCount: 1,
			expectedNames: []string{"react-best-practices"},
		},
		{
			name:          "Search by tag",
			query:         "python",
			expectedCount: 1,
			expectedNames: []string{"python-style"},
		},
		{
			name:          "Search with no matches",
			query:         "golang",
			expectedCount: 0,
			expectedNames: []string{},
		},
		{
			name:          "Case insensitive search",
			query:         "TYPESCRIPT",
			expectedCount: 1,
			expectedNames: []string{"typescript-standards"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := client.SearchRules(index, tt.query)
			assert.Equal(t, tt.expectedCount, len(results))

			if tt.expectedCount > 0 {
				for i, name := range tt.expectedNames {
					assert.Equal(t, name, results[i].Name)
				}
			}
		})
	}
}

func TestRegistryRule(t *testing.T) {
	rule := RegistryRule{
		Name:        "test-rule",
		Description: "Test rule description",
		Version:     "1.0.0",
		URL:         "https://example.com/test-rule.mdc",
		Tags:        []string{"test", "example"},
		Author:      "Test Author",
	}

	assert.Equal(t, "test-rule", rule.Name)
	assert.Equal(t, "1.0.0", rule.Version)
	assert.Contains(t, rule.Tags, "test")
}
