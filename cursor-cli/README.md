# Cursor Rules CLI (`crules`)

A Go-based CLI tool for managing Cursor IDE's `.cursorrules` and `.mdc` rule files with marketplace discovery, installation, team sharing, and version control.

Ported from [claude-plugins](https://github.com/ericboehs/claude-plugins) marketplace system and adapted for Cursor IDE's rule ecosystem.

## Features

### ✅ Phase 1: Core MVP (Complete)

- **CLI Framework**: Full Cobra-based CLI with all planned commands
- **MDC Parser**: YAML frontmatter + markdown parser for `.mdc` files
- **Validation**: Comprehensive rule validation with error/warning reporting
- **Local Installation**: Install rules from local files or directories
- **Templates**: Generate new rules from templates (basic, glob-based, always-apply)
- **Init & New**: Initialize projects and create rules from templates
- **List & Uninstall**: Manage installed rules
- **Unit Tests**: >80% coverage for parser and validator
- **All CLI Commands**: Fully implemented and functional

### ✅ Phase 2: Marketplace (Complete)

- **Registry Support**: GitHub-based marketplace registries with `index.json`
- **Remote Installation**: Install rules directly from marketplace URLs
- **Search & Discovery**: Search for rules across configured marketplaces
- **Lock Files**: `crules-lock.json` for reproducible installations
- **Update Mechanism**: Check and update installed rules to latest versions
- **Version Management**: Semantic versioning support with update detection

### 📋 Phase 3: Advanced Features (Planned)

- Import from awesome-cursorrules
- Import from cursor.directory
- Team synchronization
- Publishing workflow

## Installation

```bash
# Clone the repository
git clone https://github.com/caboose/cursor-cli.git
cd cursor-cli

# Build the binary
go build -o bin/crules ./cmd/crules

# Optional: Install to PATH
sudo cp bin/crules /usr/local/bin/
```

## Quick Start

### Local Rules

```bash
# Initialize a new project
crules init

# Create a new rule from template
crules new typescript-standards

# Validate a rule file
crules validate my-rule.mdc

# Install a local rule
crules install ./my-rules/typescript.mdc

# List installed rules
crules list

# Uninstall a rule
crules uninstall typescript.mdc
```

### Marketplace Usage

```bash
# Add a marketplace
crules marketplace add https://github.com/username/cursor-rules-marketplace

# Search for rules
crules search typescript

# Install a rule from marketplace
crules install typescript-standards

# Check for updates
crules update

# Update a specific rule
crules update typescript-standards

# List configured marketplaces
crules marketplace list
```

## Commands

### Core Commands

```bash
crules init                    # Initialize .cursor/rules/ directory
crules new <name>             # Create new rule from template
crules validate [file]        # Validate .mdc syntax
crules install <path>         # Install rule from local path
crules uninstall <name>       # Remove installed rule
crules list                   # List installed rules
```

### Marketplace Commands

```bash
crules marketplace add <url>   # Add marketplace registry
crules marketplace list        # List configured marketplaces
crules marketplace remove <name> # Remove a marketplace
crules marketplace update      # Update marketplace indices
crules search <query>          # Search for rules
crules install <rule-name>     # Install from marketplace
crules update [name]           # Update rule(s) to latest versions
```

### Team Commands (Coming Soon)

```bash
crules team export            # Export rules as bundle
crules team import <bundle>   # Import rule bundle
crules team sync              # Sync with remote
```

### Import Commands (Coming Soon)

```bash
crules import awesome         # Import from awesome-cursorrules
crules import directory       # Import from cursor.directory
crules migrate <file>         # Convert legacy .cursorrules
```

## Global Flags

```bash
-g, --global      # Apply to global rules (~/.cursor/rules/)
-p, --project     # Apply to project rules (.cursor/rules/)
-c, --config      # Custom config path
-v, --verbose     # Verbose output
    --dry-run     # Show what would happen without making changes
```

## MDC File Format

`.mdc` files are Cursor's modern rule format with YAML frontmatter and markdown content:

```markdown
---
description: "TypeScript coding standards"
globs:
  - "**/*.ts"
  - "**/*.tsx"
alwaysApply: false
tags:
  - typescript
  - coding-standards
version: "1.0.0"
---

# TypeScript Coding Standards

## Guidelines

- Use strict mode
- Prefer explicit types
- Follow naming conventions

## Examples

\```typescript
// Good
const getUserById = (id: number): User => { ... }
\```
```

### Frontmatter Fields

| Field | Required | Description |
|-------|----------|-------------|
| `description` | Yes | Brief description of the rule |
| `globs` | Yes* | File patterns to match |
| `alwaysApply` | No | Apply to all files (default: false) |
| `tags` | No | Categorization tags |
| `version` | No | Semantic version |
| `author` | No | Rule author |

\* `globs` required unless `alwaysApply: true`

## Directory Structure

```
~/.cursor/
├── config.json                # CLI configuration
├── marketplaces/              # Marketplace caches (future)
├── rules/                     # Global rules (--global)
│   └── *.mdc
└── installed/                 # Package metadata (future)

<project>/.cursor/
├── rules/                     # Project rules
│   └── *.mdc
├── crules.json                # Project config (future)
└── crules-lock.json           # Version lock (future)
```

## Development

### Project Structure

```
cursor-cli/
├── cmd/crules/                # CLI entry point
├── pkg/
│   ├── config/                # Configuration management
│   ├── marketplace/           # Marketplace client
│   ├── rule/                  # MDC parser, validator, templates
│   ├── installer/             # Installation logic
│   ├── importer/              # Import from external sources
│   ├── team/                  # Team sync
│   └── util/                  # Utilities
├── internal/cli/              # CLI commands
├── templates/                 # Rule templates
├── testdata/                  # Test fixtures
└── go.mod
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test -v ./pkg/rule/...

# Run with race detection
go test -race ./...
```

### Building

```bash
# Development build
go build -o bin/crules ./cmd/crules

# Production build with version info
go build -ldflags="-X main.version=1.0.0 -X main.commit=$(git rev-parse HEAD) -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/crules ./cmd/crules

# Cross-compile for multiple platforms
GOOS=linux GOARCH=amd64 go build -o bin/crules-linux-amd64 ./cmd/crules
GOOS=darwin GOARCH=amd64 go build -o bin/crules-darwin-amd64 ./cmd/crules
GOOS=windows GOARCH=amd64 go build -o bin/crules-windows-amd64.exe ./cmd/crules
```

## Architecture

### MDC Parser

The MDC parser (`pkg/rule/mdc.go`) handles:
- YAML frontmatter parsing using `github.com/adrg/frontmatter`
- Markdown content extraction
- Marshaling back to `.mdc` format
- File I/O operations

### Validator

The validator (`pkg/rule/validator.go`) provides:
- Required field validation
- Glob pattern validation using `doublestar`
- Content quality checks
- Semantic versioning validation
- Directory-level batch validation

### Installer

The installer (`pkg/installer/installer.go`) supports:
- Local file installation
- Directory installation (batch)
- Dry-run mode
- Validation before installation
- List/uninstall operations

### Templates

Three template types are provided:
- **basic**: General-purpose rule template
- **glob-based**: File-pattern specific rules
- **always-apply**: Global rules that apply to all files

### Marketplace

The marketplace system (`pkg/marketplace/`) provides:
- **Registry Client**: HTTP client for fetching `index.json` from registries
- **Manager**: Orchestrates multi-registry search and installation
- **URL Normalization**: Converts GitHub URLs to raw content URLs
- **Search**: Full-text search across rule names, descriptions, and tags
- **Updater**: Semantic version comparison and update detection

### Lock File

The lock file (`pkg/installer/lockfile.go`) tracks:
- Installed rules with versions
- Installation source (local vs marketplace)
- Registry metadata
- Installation and update timestamps

## Creating Your Own Marketplace

A marketplace is simply a GitHub repository (or any web-accessible location) with an `index.json` file. Here's how to create one:

### 1. Create an `index.json` file

```json
{
  "name": "My Cursor Rules Marketplace",
  "description": "Custom rules for my team",
  "homepage": "https://github.com/username/cursor-rules",
  "version": "1.0",
  "rules": [
    {
      "name": "my-rule",
      "description": "Description of the rule",
      "author": "Your Name",
      "version": "1.0.0",
      "url": "https://raw.githubusercontent.com/username/repo/main/rules/my-rule.mdc",
      "tags": ["tag1", "tag2"]
    }
  ]
}
```

### 2. Host your `.mdc` files

Place your rule files in a publicly accessible location (e.g., GitHub raw URLs).

### 3. Share your marketplace

Users can add your marketplace with:

```bash
crules marketplace add https://github.com/username/cursor-rules
```

The CLI will automatically convert GitHub URLs to raw content URLs.

## Roadmap

- [x] Phase 1: Core MVP
  - [x] CLI framework
  - [x] MDC parser
  - [x] Validator
  - [x] Local installer
  - [x] Templates
  - [x] Unit tests

- [x] Phase 2: Marketplace
  - [x] Marketplace registry client
  - [x] Remote installation
  - [x] GitHub-based registries
  - [x] Lock file support
  - [x] Update mechanism
  - [x] Search functionality
  - [x] Unit tests

- [ ] Phase 3: Advanced Features
  - [ ] awesome-cursorrules importer
  - [ ] cursor.directory importer
  - [ ] Team sync (git/bundle)
  - [ ] Publishing workflow
  - [ ] TUI improvements

## Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

[To be determined]

## Acknowledgments

- Forked from [claude-plugins](https://github.com/ericboehs/claude-plugins) by Eric Boehs
- Adapted for Cursor IDE's `.mdc` rule system
- Built with Go and the Cobra CLI framework

## Support

For issues, questions, or feature requests, please open an issue on GitHub.
