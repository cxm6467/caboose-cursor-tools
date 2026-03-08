# Caboose Cursor Tools

A collection of tools and plugins for Cursor IDE, ported from [claude-plugins](https://github.com/ericboehs/claude-plugins) by Eric Boehs.

## Projects

### 🔧 cursor-cli
Go-based CLI tool for managing Cursor IDE's `.cursorrules` and `.mdc` rule files with marketplace discovery, installation, and version control.

[**Read cursor-cli documentation →**](./cursor-cli/README.md)

**Features:**
- MDC parser and validator
- Local and remote rule installation
- Marketplace support with search and updates
- Template generation
- Lock file management

**Installation:**
```bash
cd cursor-cli
go build -o bin/crules ./cmd/crules
sudo cp bin/crules /usr/local/bin/
```

### 🔌 cursor-plugins
Complete port of claude-plugins to Cursor IDE with Go implementations and skill definitions.

[**Read cursor-plugins documentation →**](./cursor-plugins/README.md)

**7 Plugins Available:**

#### Binary Plugins (Go):
1. **session-improver** (7.6 MB) - Analyze sessions for efficiency improvements
2. **code-lint** (3.4 MB) - Multi-language automatic linting
3. **watch-ci** (3.4 MB) - Real-time GitHub Actions monitoring
4. **gh-copilot-review** (3.6 MB) - Automated Copilot review workflow

#### Skill-Only Plugins:
5. **git-utils** - Semantic commits and PR automation
6. **slack** - Complete Slack integration
7. **gist** - AI-generated gist documentation

**Installation:**
```bash
cd cursor-plugins
make build
make install
```

## Quick Start

### Install cursor-cli
```bash
cd cursor-cli
go build -o bin/crules ./cmd/crules
./bin/crules init
./bin/crules new my-rule
```

### Install cursor-plugins
```bash
cd cursor-plugins
make build        # Build all binaries
make install      # Install to /usr/local/bin

# Test the tools
watch-ci          # Monitor CI status
parse-session --current  # Analyze current session
```

## Requirements

- **Go 1.21+** (for building)
- **gh CLI** (for GitHub integration)
- **Ruby + slk gem** (optional, for Slack plugin)

## Attribution

- **cursor-plugins** are direct ports of [ericboehs/claude-plugins](https://github.com/ericboehs/claude-plugins)
- Original author: [Eric Boehs](https://github.com/ericboehs)
- ⭐ **Please star the original repository:** https://github.com/ericboehs/claude-plugins

## Documentation

- [cursor-cli README](./cursor-cli/README.md)
- [cursor-cli Implementation Status](./cursor-cli/IMPLEMENTATION_STATUS.md)
- [cursor-plugins README](./cursor-plugins/README.md)
- [cursor-plugins Status](./cursor-plugins/STATUS.md)
- [cursor-plugins Usage Guide](./cursor-plugins/USAGE.md)
- [cursor-plugins Marketplace Publishing](./cursor-plugins/MARKETPLACE.md)

## License

See individual project directories for license information.

## Contributing

Contributions welcome! For improvements to the original plugin concepts, please contribute to [ericboehs/claude-plugins](https://github.com/ericboehs/claude-plugins). For Cursor-specific adaptations, open an issue or PR here.

## Support

- **cursor-cli issues:** Open an issue in this repository
- **cursor-plugins issues:** Open an issue in this repository
- **Original claude-plugins:** https://github.com/ericboehs/claude-plugins/issues
