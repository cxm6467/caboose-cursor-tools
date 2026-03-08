package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/caboose/cursor-cli/pkg/rule"
)

// Installer handles rule installation and management
type Installer struct {
	TargetDir  string // .cursor/rules/ directory
	DryRun     bool
	Verbose    bool
	Validator  *rule.Validator
}

// NewInstaller creates a new installer
func NewInstaller(targetDir string, dryRun, verbose bool) *Installer {
	return &Installer{
		TargetDir: targetDir,
		DryRun:    dryRun,
		Verbose:   verbose,
		Validator: rule.NewValidator(false),
	}
}

// InstallResult represents the result of an installation
type InstallResult struct {
	InstalledFiles []string
	SkippedFiles   []string
	Errors         []error
}

// Install installs a rule from a local path (file or directory)
func (i *Installer) Install(sourcePath string) (*InstallResult, error) {
	// Check if source exists
	info, err := os.Stat(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("source path not found: %w", err)
	}

	// Handle file or directory
	if info.IsDir() {
		return i.installFromDirectory(sourcePath)
	}

	return i.installFromFile(sourcePath)
}

// installFromFile installs a single .mdc file
func (i *Installer) installFromFile(filePath string) (*InstallResult, error) {
	result := &InstallResult{}

	// Validate it's a .mdc file
	if filepath.Ext(filePath) != ".mdc" {
		return nil, fmt.Errorf("not a .mdc file: %s", filePath)
	}

	// Parse and validate the rule
	mdcRule, err := rule.ParseMDCFile(filePath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("failed to parse %s: %w", filePath, err))
		return result, nil
	}

	// Validate
	validationResult := i.Validator.Validate(mdcRule)
	if validationResult.HasErrors() {
		for _, valErr := range validationResult.Errors {
			if valErr.Severity == "error" {
				result.Errors = append(result.Errors, fmt.Errorf("validation error in %s: %s", filePath, valErr.Message))
			}
		}
		if len(result.Errors) > 0 {
			return result, nil
		}
	}

	// Determine target path
	targetPath := filepath.Join(i.TargetDir, filepath.Base(filePath))

	// Check if already exists
	if _, err := os.Stat(targetPath); err == nil {
		if i.Verbose {
			fmt.Printf("File already exists, overwriting: %s\n", targetPath)
		}
	}

	// Install (copy) the file
	if err := i.copyFile(filePath, targetPath); err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("failed to install %s: %w", filePath, err))
		return result, nil
	}

	result.InstalledFiles = append(result.InstalledFiles, targetPath)

	if i.Verbose {
		fmt.Printf("Installed: %s -> %s\n", filePath, targetPath)
	}

	return result, nil
}

// installFromDirectory installs all .mdc files from a directory
func (i *Installer) installFromDirectory(dirPath string) (*InstallResult, error) {
	result := &InstallResult{}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .mdc files
		if filepath.Ext(path) != ".mdc" {
			return nil
		}

		// Install individual file
		fileResult, err := i.installFromFile(path)
		if err != nil {
			result.Errors = append(result.Errors, err)
			return nil // Continue with other files
		}

		// Merge results
		result.InstalledFiles = append(result.InstalledFiles, fileResult.InstalledFiles...)
		result.SkippedFiles = append(result.SkippedFiles, fileResult.SkippedFiles...)
		result.Errors = append(result.Errors, fileResult.Errors...)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return result, nil
}

// Uninstall removes an installed rule
func (i *Installer) Uninstall(ruleName string) error {
	// Add .mdc extension if not present
	if filepath.Ext(ruleName) != ".mdc" {
		ruleName = ruleName + ".mdc"
	}

	targetPath := filepath.Join(i.TargetDir, ruleName)

	// Check if exists
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return fmt.Errorf("rule not found: %s", ruleName)
	}

	// Dry run
	if i.DryRun {
		fmt.Printf("Would remove: %s\n", targetPath)
		return nil
	}

	// Remove the file
	if err := os.Remove(targetPath); err != nil {
		return fmt.Errorf("failed to remove %s: %w", targetPath, err)
	}

	if i.Verbose {
		fmt.Printf("Removed: %s\n", targetPath)
	}

	return nil
}

// List returns all installed rules
func (i *Installer) List() ([]string, error) {
	var rules []string

	// Ensure target directory exists
	if _, err := os.Stat(i.TargetDir); os.IsNotExist(err) {
		return rules, nil // Empty list if directory doesn't exist
	}

	// Read directory
	entries, err := os.ReadDir(i.TargetDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// Filter .mdc files
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) == ".mdc" {
			rules = append(rules, entry.Name())
		}
	}

	return rules, nil
}

// copyFile copies a file from src to dst
func (i *Installer) copyFile(src, dst string) error {
	if i.DryRun {
		fmt.Printf("Would copy: %s -> %s\n", src, dst)
		return nil
	}

	// Ensure target directory exists
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Read source file
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source: %w", err)
	}

	// Write to destination
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return fmt.Errorf("failed to write destination: %w", err)
	}

	return nil
}

// Initialize creates the target directory structure
func (i *Installer) Initialize() error {
	if i.DryRun {
		fmt.Printf("Would create directory: %s\n", i.TargetDir)
		return nil
	}

	if err := os.MkdirAll(i.TargetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if i.Verbose {
		fmt.Printf("Created directory: %s\n", i.TargetDir)
	}

	return nil
}

// IsInstalled checks if a rule is already installed
func (i *Installer) IsInstalled(ruleName string) bool {
	if filepath.Ext(ruleName) != ".mdc" {
		ruleName = ruleName + ".mdc"
	}

	targetPath := filepath.Join(i.TargetDir, ruleName)
	_, err := os.Stat(targetPath)
	return err == nil
}

// InstallFromContent installs a rule from raw content (used for marketplace installs)
func (i *Installer) InstallFromContent(ruleName string, content []byte) (*InstallResult, error) {
	result := &InstallResult{}

	// Ensure .mdc extension
	if filepath.Ext(ruleName) != ".mdc" {
		ruleName = ruleName + ".mdc"
	}

	// Parse and validate the rule content
	mdcRule, err := rule.ParseMDC(content, ruleName)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("failed to parse rule: %w", err))
		return result, nil
	}

	// Validate
	validationResult := i.Validator.Validate(mdcRule)
	if validationResult.HasErrors() {
		for _, valErr := range validationResult.Errors {
			if valErr.Severity == "error" {
				result.Errors = append(result.Errors, fmt.Errorf("validation error: %s", valErr.Message))
			}
		}
		if len(result.Errors) > 0 {
			return result, nil
		}
	}

	// Determine target path
	targetPath := filepath.Join(i.TargetDir, ruleName)

	// Check if already exists
	if _, err := os.Stat(targetPath); err == nil {
		if i.Verbose {
			fmt.Printf("File already exists, overwriting: %s\n", targetPath)
		}
	}

	// Dry run
	if i.DryRun {
		fmt.Printf("Would install: %s\n", targetPath)
		result.InstalledFiles = append(result.InstalledFiles, targetPath)
		return result, nil
	}

	// Ensure target directory exists
	if err := os.MkdirAll(i.TargetDir, 0755); err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("failed to create directory: %w", err))
		return result, nil
	}

	// Write the content
	if err := os.WriteFile(targetPath, content, 0644); err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("failed to write file: %w", err))
		return result, nil
	}

	result.InstalledFiles = append(result.InstalledFiles, targetPath)

	if i.Verbose {
		fmt.Printf("Installed: %s\n", targetPath)
	}

	return result, nil
}
