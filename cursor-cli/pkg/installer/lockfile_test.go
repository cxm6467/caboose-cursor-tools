package installer

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLockFile(t *testing.T) {
	lf := NewLockFile("/tmp/test-lock.json")

	assert.NotNil(t, lf)
	assert.Equal(t, "1.0", lf.Version)
	assert.NotNil(t, lf.Rules)
	assert.Equal(t, "/tmp/test-lock.json", lf.FilePath)
}

func TestAddRule(t *testing.T) {
	lf := NewLockFile("/tmp/test-lock.json")

	rule := RuleLock{
		Name:    "test-rule",
		Version: "1.0.0",
		Source:  "marketplace",
	}

	lf.AddRule(rule)

	assert.Equal(t, 1, len(lf.Rules))
	assert.True(t, lf.HasRule("test-rule"))

	retrieved, ok := lf.GetRule("test-rule")
	assert.True(t, ok)
	assert.Equal(t, "test-rule", retrieved.Name)
	assert.Equal(t, "1.0.0", retrieved.Version)
}

func TestUpdateRule(t *testing.T) {
	lf := NewLockFile("/tmp/test-lock.json")

	// Add initial rule
	rule := RuleLock{
		Name:    "test-rule",
		Version: "1.0.0",
		Source:  "marketplace",
	}
	lf.AddRule(rule)

	// Sleep briefly to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	// Update rule
	updatedRule := RuleLock{
		Name:    "test-rule",
		Version: "2.0.0",
		Source:  "marketplace",
	}
	lf.AddRule(updatedRule)

	// Should still only have one rule
	assert.Equal(t, 1, len(lf.Rules))

	// But version should be updated
	retrieved, ok := lf.GetRule("test-rule")
	assert.True(t, ok)
	assert.Equal(t, "2.0.0", retrieved.Version)

	// UpdatedAt should be set
	assert.False(t, retrieved.UpdatedAt.IsZero())
}

func TestRemoveRule(t *testing.T) {
	lf := NewLockFile("/tmp/test-lock.json")

	rule := RuleLock{
		Name:    "test-rule",
		Version: "1.0.0",
		Source:  "marketplace",
	}
	lf.AddRule(rule)

	assert.True(t, lf.HasRule("test-rule"))

	lf.RemoveRule("test-rule")

	assert.False(t, lf.HasRule("test-rule"))
	assert.Equal(t, 0, len(lf.Rules))
}

func TestListRules(t *testing.T) {
	lf := NewLockFile("/tmp/test-lock.json")

	lf.AddRule(RuleLock{Name: "rule1", Version: "1.0.0", Source: "marketplace"})
	lf.AddRule(RuleLock{Name: "rule2", Version: "1.0.0", Source: "local"})
	lf.AddRule(RuleLock{Name: "rule3", Version: "2.0.0", Source: "marketplace"})

	rules := lf.ListRules()
	assert.Equal(t, 3, len(rules))
	assert.Contains(t, rules, "rule1")
	assert.Contains(t, rules, "rule2")
	assert.Contains(t, rules, "rule3")
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	lockFilePath := filepath.Join(tmpDir, "crules-lock.json")

	// Create and populate lock file
	lf := NewLockFile(lockFilePath)
	lf.AddRule(RuleLock{
		Name:         "test-rule",
		Version:      "1.0.0",
		Source:       "marketplace",
		RegistryName: "default",
	})

	// Save
	err := lf.Save()
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(lockFilePath)
	require.NoError(t, err)

	// Load
	loaded, err := LoadLockFile(lockFilePath)
	require.NoError(t, err)

	assert.Equal(t, lf.Version, loaded.Version)
	assert.Equal(t, 1, len(loaded.Rules))
	assert.True(t, loaded.HasRule("test-rule"))

	retrieved, ok := loaded.GetRule("test-rule")
	assert.True(t, ok)
	assert.Equal(t, "test-rule", retrieved.Name)
	assert.Equal(t, "1.0.0", retrieved.Version)
	assert.Equal(t, "marketplace", retrieved.Source)
	assert.Equal(t, "default", retrieved.RegistryName)
}

func TestLoadNonExistentLockFile(t *testing.T) {
	tmpDir := t.TempDir()
	lockFilePath := filepath.Join(tmpDir, "nonexistent.json")

	// Loading non-existent file should return empty lock file
	lf, err := LoadLockFile(lockFilePath)
	require.NoError(t, err)
	assert.NotNil(t, lf)
	assert.Equal(t, 0, len(lf.Rules))
	assert.Equal(t, lockFilePath, lf.FilePath)
}
