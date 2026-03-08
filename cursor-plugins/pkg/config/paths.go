package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectHash generates a hash from a project directory path
// Maintains compatibility with original claude-plugins format
// Example: /foo/bar → -foo-bar
func ProjectHash(projectDir string) string {
	return strings.ReplaceAll(projectDir, "/", "-")
}

// GetConfigDir returns the base configuration directory
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".cursor"), nil
}

// GetPluginConfigDir returns the plugin configuration directory
func GetPluginConfigDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "plugins"), nil
}

// GetProjectConfigPath returns the full path to a project's config file
func GetProjectConfigPath(projectDir string) (string, error) {
	pluginDir, err := GetPluginConfigDir()
	if err != nil {
		return "", err
	}

	hash := ProjectHash(projectDir)
	return filepath.Join(pluginDir, hash, "config.json"), nil
}

// LoadProjectConfig loads configuration for a specific project
func LoadProjectConfig(projectDir string) (*ProjectConfig, error) {
	configPath, err := GetProjectConfigPath(projectDir)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config not found for project %s (expected: %s)", projectDir, configPath)
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config ProjectConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// SaveProjectConfig saves configuration for a specific project
func SaveProjectConfig(projectDir string, config *ProjectConfig) error {
	configPath, err := GetProjectConfigPath(projectDir)
	if err != nil {
		return err
	}

	// Ensure directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
