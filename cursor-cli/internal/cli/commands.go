package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCommand creates the root command for crules CLI
func NewRootCommand(version, commit, date string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "crules",
		Short: "Cursor Rules CLI - Manage Cursor IDE rules from the command line",
		Long: `crules is a CLI tool for managing Cursor IDE's .cursorrules and .mdc files.
It provides marketplace discovery, installation, team sharing, and version control for Cursor rules.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	// Add global flags
	rootCmd.PersistentFlags().BoolP("global", "g", false, "Apply to global rules (~/.cursor/rules/)")
	rootCmd.PersistentFlags().BoolP("project", "p", false, "Apply to project rules (.cursor/rules/)")
	rootCmd.PersistentFlags().StringP("config", "c", "", "Custom config path")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Show what would happen without making changes")

	// Add subcommands
	rootCmd.AddCommand(newInitCommand())
	rootCmd.AddCommand(newNewCommand())
	rootCmd.AddCommand(newValidateCommand())
	rootCmd.AddCommand(newInstallCommand())
	rootCmd.AddCommand(newUninstallCommand())
	rootCmd.AddCommand(newListCommand())
	rootCmd.AddCommand(newUpdateCommand())
	rootCmd.AddCommand(newSearchCommand())
	rootCmd.AddCommand(newMarketplaceCommand())
	rootCmd.AddCommand(newTeamCommand())
	rootCmd.AddCommand(newImportCommand())
	rootCmd.AddCommand(newMigrateCommand())
	rootCmd.AddCommand(newPublishCommand())

	return rootCmd
}

func newInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize .cursor/rules/ directory in current project",
		Long:  `Creates a .cursor/rules/ directory and initializes crules.json configuration.`,
		RunE:  Init,
	}
}

func newNewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new <rule-name>",
		Short: "Create a new rule from template",
		Long:  `Generates a new .mdc rule file from a template with proper frontmatter.`,
		Args:  cobra.ExactArgs(1),
		RunE:  New,
	}
	cmd.Flags().StringP("template", "t", "basic", "Template type (basic, glob-based, always-apply)")
	return cmd
}

func newValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate [file-or-directory]",
		Short: "Validate .mdc rule syntax",
		Long:  `Validates YAML frontmatter and markdown format of .mdc files.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  Validate,
	}
}

func newInstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "install <rule-name-or-path>",
		Short: "Install a rule from marketplace or local path",
		Long:  `Downloads and installs a rule from a marketplace or installs from a local path.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ruleNameOrPath := args[0]

			// Check if it's a local file or directory
			if _, err := os.Stat(ruleNameOrPath); err == nil {
				// Local installation
				return InstallFromLocal(cmd, args)
			}

			// Marketplace installation
			return InstallFromMarketplace(cmd, args)
		},
	}
}

func newUninstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "uninstall <rule-name>",
		Aliases: []string{"remove", "rm"},
		Short:   "Uninstall a rule",
		Long:    `Removes an installed rule from the project or global rules directory.`,
		Args:    cobra.ExactArgs(1),
		RunE:    Uninstall,
	}
}

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List installed rules",
		Long:    `Shows all installed rules in the current project or globally.`,
		RunE:    List,
	}
}

func newUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update [rule-name]",
		Short: "Update installed rule(s)",
		Long:  `Updates one or all installed rules to their latest compatible versions.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  Update,
	}
}

func newSearchCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search for rules across marketplaces",
		Long:  `Searches all configured marketplaces for rules matching the query.`,
		Args:  cobra.ExactArgs(1),
		RunE:  Search,
	}
}

func newMarketplaceCommand() *cobra.Command {
	marketplaceCmd := &cobra.Command{
		Use:   "marketplace",
		Short: "Manage marketplace registries",
		Long:  `Add, list, remove, and update marketplace registries.`,
	}

	marketplaceCmd.AddCommand(&cobra.Command{
		Use:   "add <url>",
		Short: "Add a marketplace registry",
		Args:  cobra.ExactArgs(1),
		RunE:  MarketplaceAdd,
	})

	marketplaceCmd.AddCommand(&cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List configured marketplaces",
		RunE:    MarketplaceList,
	})

	marketplaceCmd.AddCommand(&cobra.Command{
		Use:     "remove <name>",
		Aliases: []string{"rm"},
		Short:   "Remove a marketplace",
		Args:    cobra.ExactArgs(1),
		RunE:    MarketplaceRemove,
	})

	marketplaceCmd.AddCommand(&cobra.Command{
		Use:   "update",
		Short: "Update marketplace indices",
		RunE:  MarketplaceUpdate,
	})

	return marketplaceCmd
}

func newTeamCommand() *cobra.Command {
	teamCmd := &cobra.Command{
		Use:   "team",
		Short: "Team collaboration commands",
		Long:  `Export, import, and sync team rules.`,
	}

	teamCmd.AddCommand(&cobra.Command{
		Use:   "export",
		Short: "Export project rules as a bundle",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("team export (not yet implemented)")
			return nil
		},
	})

	teamCmd.AddCommand(&cobra.Command{
		Use:   "import <bundle>",
		Short: "Import a rule bundle",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bundle := args[0]
			fmt.Printf("team import - importing: %s (not yet implemented)\n", bundle)
			return nil
		},
	})

	teamCmd.AddCommand(&cobra.Command{
		Use:   "sync",
		Short: "Sync team rules from remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("team sync (not yet implemented)")
			return nil
		},
	})

	return teamCmd
}

func newImportCommand() *cobra.Command {
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import rules from existing marketplaces",
		Long:  `Import rules from awesome-cursorrules, cursor.directory, or files.`,
	}

	importCmd.AddCommand(&cobra.Command{
		Use:   "awesome",
		Short: "Import from awesome-cursorrules repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("import awesome (not yet implemented)")
			return nil
		},
	})

	importCmd.AddCommand(&cobra.Command{
		Use:   "directory",
		Short: "Import from cursor.directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("import directory (not yet implemented)")
			return nil
		},
	})

	importCmd.AddCommand(&cobra.Command{
		Use:   "file <path>",
		Short: "Import from a .cursorrules file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			fmt.Printf("import file - importing: %s (not yet implemented)\n", path)
			return nil
		},
	})

	return importCmd
}

func newMigrateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate <file>",
		Short: "Migrate legacy .cursorrules to .mdc format",
		Long:  `Converts a legacy .cursorrules file to modern .cursor/rules/*.mdc format.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			file := args[0]
			fmt.Printf("migrate command - migrating: %s (not yet implemented)\n", file)
			return nil
		},
	}
}

func newPublishCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "publish",
		Short: "Publish a rule package to marketplace",
		Long:  `Validates and publishes your rule package to a marketplace registry.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("publish command (not yet implemented)")
			return nil
		},
	}
}
