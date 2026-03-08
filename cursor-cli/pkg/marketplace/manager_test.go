package marketplace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	registries := []Registry{
		{Name: "test1", URL: "https://example.com/registry1"},
		{Name: "test2", URL: "https://example.com/registry2"},
	}

	manager := NewManager(registries)

	assert.NotNil(t, manager)
	assert.Equal(t, 2, len(manager.registries))
	assert.NotNil(t, manager.client)
	assert.NotNil(t, manager.indices)
}

func TestGetRegistry(t *testing.T) {
	registries := []Registry{
		{Name: "test1", URL: "https://example.com/registry1"},
		{Name: "test2", URL: "https://example.com/registry2"},
	}

	manager := NewManager(registries)

	t.Run("Find existing registry", func(t *testing.T) {
		reg, err := manager.GetRegistry("test1")
		assert.NoError(t, err)
		assert.NotNil(t, reg)
		assert.Equal(t, "test1", reg.Name)
	})

	t.Run("Registry not found", func(t *testing.T) {
		reg, err := manager.GetRegistry("nonexistent")
		assert.Error(t, err)
		assert.Nil(t, reg)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestListRegistries(t *testing.T) {
	registries := []Registry{
		{Name: "test1", URL: "https://example.com/registry1"},
		{Name: "test2", URL: "https://example.com/registry2"},
	}

	manager := NewManager(registries)
	listed := manager.ListRegistries()

	assert.Equal(t, 2, len(listed))
	assert.Equal(t, "test1", listed[0].Name)
	assert.Equal(t, "test2", listed[1].Name)
}
