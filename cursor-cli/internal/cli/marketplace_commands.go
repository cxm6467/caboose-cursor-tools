package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caboose/cursor-cli/pkg/config"
	"github.com/caboose/cursor-cli/pkg/installer"
	"github.com/caboose/cursor-cli/pkg/marketplace"
	"github.com/spf13/cobra"
)

// getTargetDir determines the target directory based on flags
func getTargetDir(cmd *cobra.Command) (string, error) {
	globalFlag, _ := cmd.Flags().GetBool("global")
	projectFlag, _ := cmd.Flags().GetBool("project")

	if globalFlag && projectFlag {
		return "", fmt.Errorf("cannot specify both --global and --project")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	if globalFlag {
		return filepath.Join(home, ".cursor", "rules"), nil
	}

	// Default to project
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	return filepath.Join(cwd, ".cursor", "rules"), nil
}

// getLockFilePath determines the lock file path based on flags
func getLockFilePath(cmd *cobra.Command) (string, error) {
	globalFlag, _ := cmd.Flags().GetBool("global")

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	if globalFlag {
		return filepath.Join(home, ".cursor", "crules-lock.json"), nil
	}

	// Default to project
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	return filepath.Join(cwd, ".cursor", "crules-lock.json"), nil
}

// MarketplaceAdd implements: crules marketplace add <url>
func MarketplaceAdd(cmd *cobra.Command, args []string) error {
	url := args[0]

	// Extract name from URL (use last part of path)
	parts := strings.Split(strings.TrimSuffix(url, "/"), "/")
	name := parts[len(parts)-1]

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Verify the registry is accessible
	client := marketplace.NewClient()
	index, err := client.FetchIndex(url)
	if err != nil {
		return fmt.Errorf("failed to fetch marketplace index: %w", err)
	}

	description := index.Description
	if description == "" {
		description = index.Name
	}

	// Add to config
	if err := cfg.AddMarketplace(name, url, description); err != nil {
		return err
	}

	// Save config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Added marketplace: %s (%s)\n", name, url)
	fmt.Printf("  %s\n", description)
	fmt.Printf("  %d rules available\n", len(index.Rules))

	return nil
}

// MarketplaceList implements: crules marketplace list
func MarketplaceList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Marketplaces) == 0 {
		fmt.Println("No marketplaces configured")
		fmt.Println("\nAdd a marketplace with: crules marketplace add <url>")
		return nil
	}

	fmt.Printf("Configured marketplaces (%d):\n\n", len(cfg.Marketplaces))
	for _, m := range cfg.Marketplaces {
		fmt.Printf("  %s\n", m.Name)
		fmt.Printf("    URL: %s\n", m.URL)
		if m.Description != "" {
			fmt.Printf("    %s\n", m.Description)
		}
		fmt.Println()
	}

	return nil
}

// MarketplaceRemove implements: crules marketplace remove <name>
func MarketplaceRemove(cmd *cobra.Command, args []string) error {
	name := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.RemoveMarketplace(name); err != nil {
		return err
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Removed marketplace: %s\n", name)
	return nil
}

// MarketplaceUpdate implements: crules marketplace update
func MarketplaceUpdate(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Marketplaces) == 0 {
		fmt.Println("No marketplaces configured")
		return nil
	}

	manager := marketplace.NewManager(cfg.Marketplaces)

	fmt.Println("Updating marketplace indices...")
	if err := manager.UpdateIndices(); err != nil {
		return fmt.Errorf("failed to update indices: %w", err)
	}

	fmt.Printf("Updated %d marketplace(s)\n", len(cfg.Marketplaces))
	return nil
}

// Search implements: crules search <query>
func Search(cmd *cobra.Command, args []string) error {
	query := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Marketplaces) == 0 {
		fmt.Println("No marketplaces configured")
		fmt.Println("\nAdd a marketplace with: crules marketplace add <url>")
		return nil
	}

	manager := marketplace.NewManager(cfg.Marketplaces)

	// Update indices
	if err := manager.UpdateIndices(); err != nil {
		return fmt.Errorf("failed to update indices: %w", err)
	}

	// Search
	results, err := manager.Search(query)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Printf("No rules found matching '%s'\n", query)
		return nil
	}

	fmt.Printf("Found %d rule(s):\n\n", len(results))
	for _, result := range results {
		fmt.Printf("  %s (v%s)\n", result.Rule.Name, result.Rule.Version)
		fmt.Printf("    %s\n", result.Rule.Description)
		fmt.Printf("    Marketplace: %s\n", result.RegistryName)
		if len(result.Rule.Tags) > 0 {
			fmt.Printf("    Tags: %s\n", strings.Join(result.Rule.Tags, ", "))
		}
		fmt.Println()
	}

	return nil
}

// InstallFromLocal implements: crules install <path> (local mode)
func InstallFromLocal(cmd *cobra.Command, args []string) error {
	sourcePath := args[0]

	targetDir, err := getTargetDir(cmd)
	if err != nil {
		return err
	}

	lockFilePath, err := getLockFilePath(cmd)
	if err != nil {
		return err
	}

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Create installer
	inst := installer.NewInstaller(targetDir, dryRun, verbose)

	// Install
	result, err := inst.Install(sourcePath)
	if err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors occurred:")
		for _, e := range result.Errors {
			fmt.Printf("  - %v\n", e)
		}
		return fmt.Errorf("installation failed with errors")
	}

	// Update lock file for each installed file
	if !dryRun && len(result.InstalledFiles) > 0 {
		lockFile, err := installer.LoadLockFile(lockFilePath)
		if err != nil {
			return fmt.Errorf("failed to load lock file: %w", err)
		}

		for _, file := range result.InstalledFiles {
			// Extract rule name from file path
			ruleName := filepath.Base(file)
			ruleName = strings.TrimSuffix(ruleName, ".mdc")

			lockFile.AddRule(installer.RuleLock{
				Name:    ruleName,
				Version: "",
				Source:  "local",
			})
		}

		if err := lockFile.Save(); err != nil {
			return fmt.Errorf("failed to save lock file: %w", err)
		}
	}

	fmt.Printf("\nInstalled %d rule(s):\n", len(result.InstalledFiles))
	for _, file := range result.InstalledFiles {
		fmt.Printf("  %s\n", file)
	}

	return nil
}

// InstallFromMarketplace implements: crules install <rule-name> (marketplace mode)
func InstallFromMarketplace(cmd *cobra.Command, args []string) error {
	ruleName := args[0]

	targetDir, err := getTargetDir(cmd)
	if err != nil {
		return err
	}

	lockFilePath, err := getLockFilePath(cmd)
	if err != nil {
		return err
	}

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Marketplaces) == 0 {
		return fmt.Errorf("no marketplaces configured. Add one with: crules marketplace add <url>")
	}

	// Create manager
	manager := marketplace.NewManager(cfg.Marketplaces)

	// Update indices
	fmt.Println("Updating marketplace indices...")
	if err := manager.UpdateIndices(); err != nil {
		return fmt.Errorf("failed to update indices: %w", err)
	}

	// Download rule
	fmt.Printf("Installing %s...\n", ruleName)
	content, ruleInfo, err := manager.InstallRule(ruleName)
	if err != nil {
		return err
	}

	// Create installer
	inst := installer.NewInstaller(targetDir, dryRun, verbose)

	// Install from content
	result, err := inst.InstallFromContent(ruleInfo.Name, content)
	if err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors occurred:")
		for _, e := range result.Errors {
			fmt.Printf("  - %v\n", e)
		}
		return fmt.Errorf("installation failed with errors")
	}

	// Update lock file
	if !dryRun {
		lockFile, err := installer.LoadLockFile(lockFilePath)
		if err != nil {
			return fmt.Errorf("failed to load lock file: %w", err)
		}

		lockFile.AddRule(installer.RuleLock{
			Name:         ruleInfo.Name,
			Version:      ruleInfo.Version,
			Source:       "marketplace",
			RegistryName: "", // Could be extracted from search result
		})

		if err := lockFile.Save(); err != nil {
			return fmt.Errorf("failed to save lock file: %w", err)
		}
	}

	fmt.Printf("\nInstalled %s (v%s)\n", ruleInfo.Name, ruleInfo.Version)
	for _, file := range result.InstalledFiles {
		fmt.Printf("  %s\n", file)
	}

	return nil
}

// Update implements: crules update [rule-name]
func Update(cmd *cobra.Command, args []string) error {
	lockFilePath, err := getLockFilePath(cmd)
	if err != nil {
		return err
	}

	targetDir, err := getTargetDir(cmd)
	if err != nil {
		return err
	}

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Load config and lock file
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	lockFile, err := installer.LoadLockFile(lockFilePath)
	if err != nil {
		return fmt.Errorf("failed to load lock file: %w", err)
	}

	manager := marketplace.NewManager(cfg.Marketplaces)

	// Check for updates
	var updates []marketplace.UpdateCheck
	if len(args) == 0 {
		// Update all rules
		fmt.Println("Checking for updates...")
		updates, err = manager.CheckUpdates(lockFile)
		if err != nil {
			return err
		}
	} else {
		// Update specific rule
		ruleName := args[0]
		update, err := manager.CheckUpdate(ruleName, lockFile)
		if err != nil {
			return err
		}
		if update.Available {
			updates = []marketplace.UpdateCheck{*update}
		}
	}

	if len(updates) == 0 {
		fmt.Println("All rules are up to date")
		return nil
	}

	// Show available updates
	fmt.Printf("Updates available for %d rule(s):\n\n", len(updates))
	for _, update := range updates {
		if update.Available {
			fmt.Printf("  %s: %s -> %s\n", update.RuleName, update.CurrentVersion, update.LatestVersion)
		}
	}

	if dryRun {
		return nil
	}

	// Perform updates
	fmt.Println("\nUpdating...")
	inst := installer.NewInstaller(targetDir, dryRun, verbose)

	for _, update := range updates {
		if !update.Available {
			continue
		}

		content, _, err := manager.InstallRule(update.RuleName)
		if err != nil {
			fmt.Printf("  Failed to update %s: %v\n", update.RuleName, err)
			continue
		}

		result, err := inst.InstallFromContent(update.RuleName, content)
		if err != nil || len(result.Errors) > 0 {
			fmt.Printf("  Failed to update %s\n", update.RuleName)
			continue
		}

		// Update lock file
		if lock, ok := lockFile.GetRule(update.RuleName); ok {
			lock.Version = update.LatestVersion
			lockFile.AddRule(*lock)
		}

		fmt.Printf("  Updated %s to v%s\n", update.RuleName, update.LatestVersion)
	}

	// Save lock file
	if err := lockFile.Save(); err != nil {
		return fmt.Errorf("failed to save lock file: %w", err)
	}

	return nil
}
