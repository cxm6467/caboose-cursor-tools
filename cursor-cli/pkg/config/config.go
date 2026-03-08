package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/caboose/cursor-cli/pkg/marketplace"
)

// Config represents the crules CLI configuration
type Config struct {
	Marketplaces []marketplace.Registry `json:"marketplaces"`
	DefaultScope string                 `json:"defaultScope"` // "global" or "project"
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".cursor", "config.json"), nil
}

// Load loads the configuration from disk
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Marketplaces: []marketplace.Registry{},
			DefaultScope: "project",
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to disk
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// AddMarketplace adds a new marketplace to the configuration
func (c *Config) AddMarketplace(name, url, description string) error {
	// Check if marketplace already exists
	for _, m := range c.Marketplaces {
		if m.Name == name {
			return fmt.Errorf("marketplace '%s' already exists", name)
		}
		if m.URL == url {
			return fmt.Errorf("marketplace with URL '%s' already exists", url)
		}
	}

	c.Marketplaces = append(c.Marketplaces, marketplace.Registry{
		Name:        name,
		URL:         url,
		Description: description,
	})

	return nil
}

// RemoveMarketplace removes a marketplace from the configuration
func (c *Config) RemoveMarketplace(name string) error {
	for i, m := range c.Marketplaces {
		if m.Name == name {
			c.Marketplaces = append(c.Marketplaces[:i], c.Marketplaces[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("marketplace '%s' not found", name)
}

// GetMarketplace retrieves a marketplace by name
func (c *Config) GetMarketplace(name string) (*marketplace.Registry, error) {
	for _, m := range c.Marketplaces {
		if m.Name == name {
			return &m, nil
		}
	}
	return nil, fmt.Errorf("marketplace '%s' not found", name)
}
