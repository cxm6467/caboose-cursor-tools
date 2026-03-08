package marketplace

import (
	"fmt"

	"github.com/caboose/cursor-cli/pkg/installer"
	"github.com/hashicorp/go-version"
)

// UpdateCheck represents an available update for a rule
type UpdateCheck struct {
	RuleName       string
	CurrentVersion string
	LatestVersion  string
	Available      bool
	RegistryName   string
	DownloadURL    string
}

// CheckUpdates checks for updates for all installed rules
func (m *Manager) CheckUpdates(lockFile *installer.LockFile) ([]UpdateCheck, error) {
	var updates []UpdateCheck

	// Refresh indices
	if err := m.UpdateIndices(); err != nil {
		return nil, fmt.Errorf("failed to update indices: %w", err)
	}

	// Check each installed rule
	for ruleName, ruleLock := range lockFile.Rules {
		// Only check marketplace-installed rules
		if ruleLock.Source != "marketplace" {
			continue
		}

		// Find the rule in the marketplace
		result, err := m.FindRule(ruleName)
		if err != nil {
			// Rule no longer in marketplace
			continue
		}

		// Compare versions
		updateAvailable, err := isUpdateAvailable(ruleLock.Version, result.Rule.Version)
		if err != nil {
			// Skip if version comparison fails
			continue
		}

		updates = append(updates, UpdateCheck{
			RuleName:       ruleName,
			CurrentVersion: ruleLock.Version,
			LatestVersion:  result.Rule.Version,
			Available:      updateAvailable,
			RegistryName:   result.RegistryName,
			DownloadURL:    result.Rule.URL,
		})
	}

	return updates, nil
}

// CheckUpdate checks if an update is available for a specific rule
func (m *Manager) CheckUpdate(ruleName string, lockFile *installer.LockFile) (*UpdateCheck, error) {
	ruleLock, ok := lockFile.GetRule(ruleName)
	if !ok {
		return nil, fmt.Errorf("rule '%s' not installed", ruleName)
	}

	if ruleLock.Source != "marketplace" {
		return nil, fmt.Errorf("rule '%s' was not installed from marketplace", ruleName)
	}

	// Refresh indices
	if err := m.UpdateIndices(); err != nil {
		return nil, fmt.Errorf("failed to update indices: %w", err)
	}

	// Find the rule in the marketplace
	result, err := m.FindRule(ruleName)
	if err != nil {
		return nil, fmt.Errorf("rule not found in marketplace: %w", err)
	}

	// Compare versions
	updateAvailable, err := isUpdateAvailable(ruleLock.Version, result.Rule.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to compare versions: %w", err)
	}

	return &UpdateCheck{
		RuleName:       ruleName,
		CurrentVersion: ruleLock.Version,
		LatestVersion:  result.Rule.Version,
		Available:      updateAvailable,
		RegistryName:   result.RegistryName,
		DownloadURL:    result.Rule.URL,
	}, nil
}

// isUpdateAvailable compares two semantic versions
func isUpdateAvailable(current, latest string) (bool, error) {
	// Handle empty versions
	if current == "" || latest == "" {
		return false, nil
	}

	currVer, err := version.NewVersion(current)
	if err != nil {
		return false, fmt.Errorf("invalid current version: %w", err)
	}

	latestVer, err := version.NewVersion(latest)
	if err != nil {
		return false, fmt.Errorf("invalid latest version: %w", err)
	}

	return latestVer.GreaterThan(currVer), nil
}
