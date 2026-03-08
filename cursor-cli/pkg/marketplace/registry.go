package marketplace

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Registry represents a marketplace registry
type Registry struct {
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Description string    `json:"description,omitempty"`
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
}

// RegistryIndex represents the index.json file in a registry
type RegistryIndex struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Homepage    string         `json:"homepage,omitempty"`
	Rules       []RegistryRule `json:"rules"`
	Version     string         `json:"version"` // Registry schema version
}

// RegistryRule represents a rule entry in the registry index
type RegistryRule struct {
	Name        string   `json:"name"`        // e.g., "typescript-standards"
	Description string   `json:"description"` // Brief description
	Author      string   `json:"author,omitempty"`
	Version     string   `json:"version"`     // Latest version (semver)
	URL         string   `json:"url"`         // URL to the .mdc file
	Tags        []string `json:"tags,omitempty"`
	Homepage    string   `json:"homepage,omitempty"`
	Repository  string   `json:"repository,omitempty"`
}

// Client handles marketplace operations
type Client struct {
	httpClient *http.Client
	userAgent  string
}

// NewClient creates a new marketplace client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: "crules-cli/1.0",
	}
}

// FetchIndex fetches the index.json from a registry URL
func (c *Client) FetchIndex(registryURL string) (*RegistryIndex, error) {
	// Normalize URL: if it's a GitHub repo, convert to raw content URL
	indexURL := normalizeRegistryURL(registryURL)

	req, err := http.NewRequest("GET", indexURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch index: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch index: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var index RegistryIndex
	if err := json.Unmarshal(body, &index); err != nil {
		return nil, fmt.Errorf("failed to parse index: %w", err)
	}

	return &index, nil
}

// FetchRule downloads a rule file from a URL
func (c *Client) FetchRule(ruleURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", ruleURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rule: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch rule: status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// SearchRules searches for rules matching a query across all rules in an index
func (c *Client) SearchRules(index *RegistryIndex, query string) []RegistryRule {
	query = strings.ToLower(query)
	var results []RegistryRule

	for _, rule := range index.Rules {
		// Search in name, description, and tags
		if strings.Contains(strings.ToLower(rule.Name), query) ||
			strings.Contains(strings.ToLower(rule.Description), query) {
			results = append(results, rule)
			continue
		}

		// Search in tags
		for _, tag := range rule.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				results = append(results, rule)
				break
			}
		}
	}

	return results
}

// normalizeRegistryURL converts GitHub repo URLs to raw content URLs for index.json
// Examples:
// - https://github.com/user/repo -> https://raw.githubusercontent.com/user/repo/main/index.json
// - https://github.com/user/repo/tree/branch -> https://raw.githubusercontent.com/user/repo/branch/index.json
func normalizeRegistryURL(url string) string {
	// If it's already a direct URL to index.json, return as-is
	if strings.HasSuffix(url, "index.json") {
		return url
	}

	// Handle GitHub URLs
	if strings.Contains(url, "github.com") {
		// Remove trailing slash
		url = strings.TrimSuffix(url, "/")

		// Extract owner/repo and optional branch
		parts := strings.Split(url, "github.com/")
		if len(parts) != 2 {
			return url + "/index.json"
		}

		path := parts[1]
		branch := "main"

		// Check if URL contains /tree/branch
		if strings.Contains(path, "/tree/") {
			treeParts := strings.Split(path, "/tree/")
			path = treeParts[0]
			if len(treeParts) > 1 {
				branch = treeParts[1]
			}
		}

		// Convert to raw.githubusercontent.com URL
		return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/index.json", path, branch)
	}

	// For non-GitHub URLs, assume index.json is at the root
	return strings.TrimSuffix(url, "/") + "/index.json"
}
