# Caboose Cursor Tools

Two powerful toolsets for Cursor IDE development workflows.

## What's Inside

### 🔧 cursor-cli
Go-based CLI tool for managing Cursor IDE's `.cursorrules` and `.mdc` rule files with marketplace discovery, installation, and version control.

[**Read cursor-cli documentation →**](./cursor-cli/README.md)

---

### 🔌 cursor-plugins
**Cursor IDE port of [claude-plugins](https://github.com/ericboehs/claude-plugins) by Eric Boehs.**

The originals are Ruby/Bash scripts designed for Claude Code CLI. This port reimplements them as **Go binaries** using Cursor's `.cursor-plugin` format, so they run seamlessly inside Cursor IDE.

[**Read cursor-plugins documentation →**](./cursor-plugins/README.md)

#### What It Is

A **bundle of 7 sub-plugins** that enhance your Cursor workflow:

| Plugin | What It Does | Type |
|--------|-------------|------|
| **session-improver** | Analyzes Cursor session history (SQLite) for inefficiencies, token waste, linter loops | Go Binary + Skill |
| **code-lint** | Runs linters automatically when you edit files via Cursor hooks | Go Binary (Hook) |
| **watch-ci** | Watches GitHub Actions CI for current branch with live status updates | Go Binary |
| **gh-copilot-review** | Watches for Copilot PR reviews, displays comments, helps apply/resolve feedback | Go Binary + Skill |
| **git-utils** | Git workflows: semantic commit + push, merge PR + cleanup branches | Skills |
| **slack** | Slack integration via slk CLI: unread, channels, search, status, messages | Skills |
| **gist** | Create/update GitHub Gists with AI-generated README documentation | Skills |

#### How You Use It

**1. Commands (Slash or Natural Language)**
```bash
/commit-and-push          # Invoke skill by name
/merge-and-cleanup        # Merge PR and cleanup
/slack                    # Check Slack messages
/gist-create myfile.js    # Create gist with AI docs
/improve-session          # Analyze session efficiency
/gh-copilot-review        # Process Copilot feedback
```

Or just say things naturally:
- "commit and push"
- "watch ci"
- "check slack"
- "create a gist from this file"

**2. Skills**
Each sub-plugin defines skills (instructions in `skills/*/SKILL.md`). When you invoke a skill, Cursor AI reads the SKILL instructions and executes the workflow — often by running the plugin's binaries or CLI commands.

**3. Hooks (code-lint)**
`code-lint` uses Cursor's hook system. The plugin ships a `hooks.json` that runs `lint-daemon` on `afterFileEdit`. When you (or the AI) edit a file, Cursor automatically runs the hook, lints the file, and feeds results back into the chat.

**4. Binaries**
Some features rely on Go binaries: `parse-session`, `lint-daemon`, `watch-ci`, `watch-copilot-reviews`. After `make install`, these are available in your PATH for both manual use and skill automation.

#### Installation
```bash
cd cursor-plugins
make build        # Build all Go binaries
make install      # Install to /usr/local/bin

# Binaries are now available:
watch-ci          # Monitor CI status
parse-session --current  # Analyze current session
```

## Quick Start

### Option 1: cursor-cli (Rules Management)
```bash
cd cursor-cli
go build -o bin/crules ./cmd/crules
sudo cp bin/crules /usr/local/bin/

# Use it
crules init                    # Initialize .cursor/rules/
crules new typescript-rules    # Create new rule
crules marketplace add <url>   # Add marketplace
crules search react            # Search for rules
crules install react-patterns  # Install from marketplace
```

### Option 2: cursor-plugins (Workflow Automation)
```bash
cd cursor-plugins
make build        # Build all Go binaries
make install      # Install to /usr/local/bin

# Use directly
watch-ci                      # Watch CI for current branch
parse-session --current       # Analyze session

# Or use via Cursor AI
# Just say: "commit and push"
# Or type: /slack
```

### Using Skills in Cursor

Once installed, you can:
1. **Type slash commands:** `/commit-and-push`, `/slack`, `/gist-create myfile.js`
2. **Speak naturally:** "commit and push these changes", "check my slack messages"
3. **Let hooks work:** Edit a file → `code-lint` auto-runs (if hooks configured)

See [cursor-plugins/USAGE.md](./cursor-plugins/USAGE.md) for detailed setup.

## Requirements

### For cursor-cli
- **Go 1.21+** (for building)

### For cursor-plugins
- **Go 1.21+** (for building binaries)
- **gh CLI** (for GitHub integration: `watch-ci`, `gh-copilot-review`)
- **Ruby + slk gem** (optional, for Slack plugin)
- **Git** (for git-utils skills)

## Attribution & Port Details

### cursor-plugins
**Direct port of [ericboehs/claude-plugins](https://github.com/ericboehs/claude-plugins)**

- **Original author:** [Eric Boehs](https://github.com/ericboehs)
- ⭐ **Please star the original:** https://github.com/ericboehs/claude-plugins

#### What Changed in the Port

| Aspect | Original (claude-plugins) | This Port (cursor-plugins) |
|--------|---------------------------|----------------------------|
| **Target IDE** | Claude Code CLI | Cursor IDE |
| **Language** | Ruby/Bash scripts | Go binaries |
| **Plugin Format** | `.claude-plugin` | `.cursor-plugin` |
| **Session Storage** | JSONL files | SQLite database |
| **File Paths** | `~/.claude/` | `~/.cursor/` |
| **Rules Files** | `CLAUDE.md` | `.cursorrules` |

#### What Stayed the Same

- ✅ All workflow logic and behavior
- ✅ Skill definitions and instructions
- ✅ Command-line interfaces
- ✅ User experience

The port maintains **100% functional compatibility** while adapting to Cursor's architecture.

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
