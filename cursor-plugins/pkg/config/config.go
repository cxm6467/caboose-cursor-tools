package config

// ProjectConfig represents the complete configuration for a project
type ProjectConfig struct {
	Version      int                    `json:"version"`
	ProjectDir   string                 `json:"project_dir"`
	Languages    map[string]LangConfig  `json:"languages"`
	ToolLinters  map[string]ToolLinter  `json:"tool_linters,omitempty"`
	HookSettings HookSettings           `json:"hook_settings"`
}

// LangConfig represents configuration for a specific language
type LangConfig struct {
	Enabled bool                   `json:"enabled"`
	Linters map[string]LinterConfig `json:"linters"`
}

// LinterConfig represents configuration for a single linter
type LinterConfig struct {
	Enabled          bool     `json:"enabled"`
	Command          string   `json:"command"`
	AutofixCommand   string   `json:"autofix_command,omitempty"`
	FilePatterns     []string `json:"file_patterns"`
	ExcludePatterns  []string `json:"exclude_patterns,omitempty"`
	Timeout          int      `json:"timeout"`                      // seconds
	WholeProjectOnly bool     `json:"whole_project_only,omitempty"` // e.g., for brakeman, clippy
}

// ToolLinter represents language-agnostic linters like semgrep
type ToolLinter struct {
	Enabled         bool     `json:"enabled"`
	Command         string   `json:"command"`
	AutofixCommand  string   `json:"autofix_command,omitempty"`
	Rulesets        []string `json:"rulesets,omitempty"`
	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
	Timeout         int      `json:"timeout"`
	HookMode        string   `json:"hook_mode,omitempty"` // "per-file" or "project-wide"
}

// HookSettings represents global hook behavior settings
type HookSettings struct {
	StopOnFirstFailure bool `json:"stop_on_first_failure"`
	AutofixBeforeLint  bool `json:"autofix_before_lint"`
	SkipGeneratedFiles bool `json:"skip_generated_files"`
}
