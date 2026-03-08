package manifest

// Author represents plugin author information
type Author struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Plugin represents a single plugin in the marketplace
type Plugin struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      Author `json:"author"`
	Source      string `json:"source,omitempty"`   // For marketplace.json
	Category    string `json:"category,omitempty"` // For marketplace.json
}

// Marketplace represents the root marketplace configuration
type Marketplace struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Plugins     []Plugin `json:"plugins"`
}

// PluginManifest represents an individual plugin's manifest (plugin.json)
type PluginManifest struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      Author `json:"author"`
}
