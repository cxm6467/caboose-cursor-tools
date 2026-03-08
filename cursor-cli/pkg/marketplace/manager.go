package marketplace

import (
	"fmt"
	"sync"
)

// Manager handles operations across multiple marketplace registries
type Manager struct {
	client     *Client
	registries []Registry
	indices    map[string]*RegistryIndex // Cache of fetched indices
	mu         sync.RWMutex
}

// NewManager creates a new marketplace manager
func NewManager(registries []Registry) *Manager {
	return &Manager{
		client:     NewClient(),
		registries: registries,
		indices:    make(map[string]*RegistryIndex),
	}
}

// UpdateIndices fetches and caches indices from all configured registries
func (m *Manager) UpdateIndices() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []error
	for _, registry := range m.registries {
		index, err := m.client.FetchIndex(registry.URL)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to fetch %s: %w", registry.Name, err))
			continue
		}
		m.indices[registry.Name] = index
	}

	if len(errors) > 0 && len(errors) == len(m.registries) {
		return fmt.Errorf("failed to fetch all registries: %v", errors)
	}

	return nil
}

// Search searches for rules across all registries
func (m *Manager) Search(query string) ([]SearchResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []SearchResult
	for registryName, index := range m.indices {
		matches := m.client.SearchRules(index, query)
		for _, match := range matches {
			results = append(results, SearchResult{
				RegistryName: registryName,
				Rule:         match,
			})
		}
	}

	return results, nil
}

// FindRule finds a specific rule by name across all registries
func (m *Manager) FindRule(ruleName string) (*SearchResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for registryName, index := range m.indices {
		for _, rule := range index.Rules {
			if rule.Name == ruleName {
				return &SearchResult{
					RegistryName: registryName,
					Rule:         rule,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("rule '%s' not found in any marketplace", ruleName)
}

// InstallRule downloads and returns the content of a rule
func (m *Manager) InstallRule(ruleName string) ([]byte, *RegistryRule, error) {
	result, err := m.FindRule(ruleName)
	if err != nil {
		return nil, nil, err
	}

	content, err := m.client.FetchRule(result.Rule.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download rule: %w", err)
	}

	return content, &result.Rule, nil
}

// GetRegistry returns a specific registry by name
func (m *Manager) GetRegistry(name string) (*Registry, error) {
	for _, reg := range m.registries {
		if reg.Name == name {
			return &reg, nil
		}
	}
	return nil, fmt.Errorf("registry '%s' not found", name)
}

// ListRegistries returns all configured registries
func (m *Manager) ListRegistries() []Registry {
	return m.registries
}

// SearchResult represents a search result with registry context
type SearchResult struct {
	RegistryName string
	Rule         RegistryRule
}
