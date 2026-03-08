package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caboose/cursor-cli/pkg/installer"
	"github.com/caboose/cursor-cli/pkg/rule"
	"github.com/spf13/cobra"
)

// Init implements: crules init
func Init(cmd *cobra.Command, args []string) error {
	targetDir, err := getTargetDir(cmd)
	if err != nil {
		return err
	}

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	verbose, _ := cmd.Flags().GetBool("verbose")

	inst := installer.NewInstaller(targetDir, dryRun, verbose)

	if err := inst.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	if !dryRun {
		fmt.Printf("Initialized rules directory: %s\n", targetDir)
	}

	return nil
}

// New implements: crules new <rule-name>
func New(cmd *cobra.Command, args []string) error {
	ruleName := args[0]

	targetDir, err := getTargetDir(cmd)
	if err != nil {
		return err
	}

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Get template type from flag
	templateType, _ := cmd.Flags().GetString("template")

	// Create target file path
	fileName := ruleName
	if !strings.HasSuffix(fileName, ".mdc") {
		fileName = fileName + ".mdc"
	}
	targetPath := filepath.Join(targetDir, fileName)

	// Check if file already exists
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("rule file already exists: %s", targetPath)
	}

	// Generate from template using the generator
	// Try to find template directory
	templateDir := ""
	possiblePaths := []string{
		"templates",
		"../../templates",
		"/home/caboose/dev/caboose-cursor-rules/cursor-cli/templates",
	}
	for _, p := range possiblePaths {
		if _, err := os.Stat(p); err == nil {
			templateDir = p
			break
		}
	}

	generator := rule.NewGenerator(templateDir)

	// Prepare template data
	data := rule.TemplateData{
		Name:        rule.FormatName(ruleName),
		Description: fmt.Sprintf("Rules for %s", rule.FormatName(ruleName)),
		Globs:       []string{}, // Will be set by template type
		Tags:        []string{strings.ToLower(ruleName)},
		Author:      rule.GetCurrentUser(),
	}

	// Set globs based on template type
	switch templateType {
	case "glob-based":
		data.Globs = []string{"**/*"}
	case "basic":
		// Basic template uses alwaysApply, but we'll use the always-apply template instead
		templateType = "always-apply"
	}

	if dryRun {
		fmt.Printf("Would create: %s\n", targetPath)
		if verbose {
			fmt.Printf("  Template: %s\n", templateType)
			fmt.Printf("  Name: %s\n", data.Name)
			fmt.Printf("  Description: %s\n", data.Description)
		}
		return nil
	}

	// Ensure directory exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate from template
	if err := generator.Generate(rule.TemplateType(templateType), data, targetPath); err != nil {
		return fmt.Errorf("failed to generate rule: %w", err)
	}

	fmt.Printf("Created rule: %s\n", targetPath)
	if verbose {
		fmt.Printf("  Template: %s\n", templateType)
		fmt.Printf("  Description: %s\n", data.Description)
	}

	return nil
}

// Validate implements: crules validate [file]
func Validate(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	verbose, _ := cmd.Flags().GetBool("verbose")

	// Check if path is a file or directory
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path not found: %w", err)
	}

	validator := rule.NewValidator(verbose)

	var resultsMap map[string]*rule.ValidationResult

	if info.IsDir() {
		// Validate all .mdc files in directory
		var err error
		resultsMap, err = validator.ValidateDirectory(path)
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	} else {
		// Validate single file
		if filepath.Ext(path) != ".mdc" {
			return fmt.Errorf("not a .mdc file: %s", path)
		}

		result, err := validator.ValidateFile(path)
		if err != nil {
			return fmt.Errorf("failed to validate file: %w", err)
		}

		resultsMap = map[string]*rule.ValidationResult{path: result}
	}

	// Print results
	errorCount := 0
	warningCount := 0

	for filePath, result := range resultsMap {
		if len(result.Errors) == 0 {
			if verbose {
				fmt.Printf("✓ %s - Valid\n", filePath)
			}
			continue
		}

		fmt.Printf("\n%s:\n", filePath)
		for _, valErr := range result.Errors {
			if valErr.Severity == "error" {
				errorCount++
				fmt.Printf("  ✗ ERROR: %s\n", valErr.Message)
			} else {
				warningCount++
				fmt.Printf("  ⚠ WARNING: %s\n", valErr.Message)
			}
			if valErr.Field != "" {
				fmt.Printf("    Field: %s\n", valErr.Field)
			}
		}
	}

	// Summary
	fmt.Printf("\nValidated %d file(s)\n", len(resultsMap))
	if errorCount > 0 {
		fmt.Printf("  Errors: %d\n", errorCount)
	}
	if warningCount > 0 {
		fmt.Printf("  Warnings: %d\n", warningCount)
	}

	if errorCount > 0 {
		return fmt.Errorf("validation failed with %d error(s)", errorCount)
	}

	if errorCount == 0 && warningCount == 0 {
		fmt.Println("  All valid ✓")
	}

	return nil
}

// List implements: crules list
func List(cmd *cobra.Command, args []string) error {
	targetDir, err := getTargetDir(cmd)
	if err != nil {
		return err
	}

	verbose, _ := cmd.Flags().GetBool("verbose")

	inst := installer.NewInstaller(targetDir, false, verbose)

	rules, err := inst.List()
	if err != nil {
		return fmt.Errorf("failed to list rules: %w", err)
	}

	if len(rules) == 0 {
		fmt.Println("No rules installed")
		fmt.Printf("\nInstall rules with: crules install <path-or-name>\n")
		return nil
	}

	fmt.Printf("Installed rules (%d):\n\n", len(rules))

	// Load lock file to get metadata
	lockFilePath, err := getLockFilePath(cmd)
	if err != nil {
		return err
	}

	lockFile, err := installer.LoadLockFile(lockFilePath)
	if err != nil {
		// Non-fatal, just list files
		for _, ruleName := range rules {
			fmt.Printf("  • %s\n", ruleName)
		}
		return nil
	}

	// Show with metadata from lock file
	for _, ruleName := range rules {
		fmt.Printf("  • %s", ruleName)

		// Strip .mdc extension for lock file lookup
		lookupName := strings.TrimSuffix(ruleName, ".mdc")
		if lock, ok := lockFile.GetRule(lookupName); ok {
			if lock.Version != "" {
				fmt.Printf(" (v%s)", lock.Version)
			}
			if lock.Source != "" {
				fmt.Printf(" [%s]", lock.Source)
			}
		}
		fmt.Println()

		// Show details in verbose mode
		if verbose {
			rulePath := filepath.Join(targetDir, ruleName)
			if mdcRule, err := rule.ParseMDCFile(rulePath); err == nil {
				fmt.Printf("      %s\n", mdcRule.Frontmatter.Description)
				if len(mdcRule.Frontmatter.Tags) > 0 {
					fmt.Printf("      Tags: %s\n", strings.Join(mdcRule.Frontmatter.Tags, ", "))
				}
			}
		}
	}

	return nil
}

// Uninstall implements: crules uninstall <rule-name>
func Uninstall(cmd *cobra.Command, args []string) error {
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

	inst := installer.NewInstaller(targetDir, dryRun, verbose)

	// Uninstall the file
	if err := inst.Uninstall(ruleName); err != nil {
		return fmt.Errorf("failed to uninstall: %w", err)
	}

	// Update lock file
	if !dryRun {
		lockFile, err := installer.LoadLockFile(lockFilePath)
		if err != nil {
			// Non-fatal
			fmt.Printf("Warning: failed to load lock file: %v\n", err)
		} else {
			// Strip .mdc extension for lock file
			lookupName := strings.TrimSuffix(ruleName, ".mdc")
			lockFile.RemoveRule(lookupName)
			if err := lockFile.Save(); err != nil {
				fmt.Printf("Warning: failed to update lock file: %v\n", err)
			}
		}
	}

	if !dryRun {
		fmt.Printf("Uninstalled: %s\n", ruleName)
	}

	return nil
}
