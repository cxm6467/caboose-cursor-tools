# Cursor Plugins - Ported from claude-plugins

> 🙏 **Attribution:** These plugins are **direct ports** of the excellent [claude-plugins](https://github.com/ericboehs/claude-plugins) repository by **[Eric Boehs](https://github.com/ericboehs)**. All original ideas, workflows, and designs belong to Eric. This port simply adapts the implementation from Ruby/Bash to Go and from Claude Code CLI to Cursor IDE.
>
> **⭐ Please star the original repository:** https://github.com/ericboehs/claude-plugins

## About This Port

This repository contains **Cursor IDE adaptations** of Eric Boehs' claude-plugins. The original repository is designed for the Claude Code CLI, while this port targets the Cursor IDE environment.

### What Changed
- Ruby/Bash → Go (cross-platform binaries)
- `.claude-plugin` → `.cursor-plugin` format
- Claude Code CLI paths → Cursor IDE paths
- Claude session JSONL → Cursor SQLite session storage

### What Stayed the Same
- All workflow logic and behavior
- Skill definitions and instructions
- Command-line interfaces
- User experience and documentation

See [CREDITS.md](./CREDITS.md) for detailed attribution.

---

## Available Plugins

**Status: 7/7 Fully Implemented ✅**

### 🔍 session-improver
**Port of:** [claude-plugins/session-improver](https://github.com/ericboehs/claude-plugins/tree/main/plugins/session-improver)

Optimize your Cursor AI workflow by analyzing session history to identify inefficiencies. Detects token waste from repeated file reads, linter loops from unfixed errors, and missing automations.

```bash
parse-session --current        # Analyze current session
parse-session <session-id>     # Analyze specific session
```

**Skill:** `/improve-session` - Full analysis with recommendations

**Type:** Go Binary (7.9 MB) + Skill
**Port:** Ruby script parsing JSONL → Go binary parsing SQLite

---

### 🔧 code-lint
**Port of:** [claude-plugins/code-lint](https://github.com/ericboehs/claude-plugins/tree/main/plugins/code-lint)

Automatic code quality enforcement with multi-language linting. Runs linters automatically on file edits via Cursor hooks, providing instant feedback. Supports ESLint, Pylint, RuboCop, golint, and more with per-project `.linter.json` configuration.

```bash
# Automatically runs on file edits via .cursor/hooks.json
# No manual invocation needed
```

**Type:** Go Binary (3.5 MB) - Hook Daemon
**Port:** Bash hook script → Go daemon with config detection

---

### ⏱️ watch-ci
**Port of:** [claude-plugins/watch-ci](https://github.com/ericboehs/claude-plugins/tree/main/plugins/watch-ci)

Real-time GitHub Actions CI monitoring with auto-exit on completion. Watch your CI pipeline's progress with live status updates, job tracking, and color-coded results. Stay in flow while CI runs.

```bash
watch-ci          # Watch current branch (10s intervals)
watch-ci 30       # Custom polling interval (30s)
```

**Type:** Go Binary (3.5 MB)
**Port:** Bash script with jq → Go binary with native JSON parsing

---

### 🤖 gh-copilot-review
**Port of:** [claude-plugins/gh-copilot-review](https://github.com/ericboehs/claude-plugins/tree/main/plugins/gh-copilot-review)

Streamline GitHub Copilot code review workflow with automated feedback monitoring. Watches for Copilot PR reviews, displays inline comments, and enables AI-assisted fixes with automatic thread resolution.

```bash
watch-copilot-reviews           # Watch current branch PR
watch-copilot-reviews 42        # Watch specific PR #42
```

**Skill:** `/gh-copilot-review` - Full automated workflow (wait → display → fix → resolve → push)

**Type:** Go Binary (3.6 MB) + Skill
**Port:** Bash script → Go binary with GitHub API integration

---

### 🔀 git-utils
**Port of:** [claude-plugins/git-utils](https://github.com/ericboehs/claude-plugins/tree/main/plugins/git-utils)

Streamlined git workflow automation with intelligent skills. Automates common git patterns: atomic commits with semantic messages, PR merging with branch cleanup, and worktree management.

**Skills:**
- `/commit-and-push` - Stage, commit with semantic message, and push
- `/merge-and-cleanup` - Merge PR, delete branches, switch to main, clean worktrees

**Type:** Skills-Only (wraps git/gh commands)
**Port:** 1:1 skill definitions with enhanced descriptions

---

### 💬 slack
**Port of:** [claude-plugins/slack](https://github.com/ericboehs/claude-plugins/tree/main/plugins/slack)

Complete Slack integration for Cursor workflow. Check unread messages, read channels, search conversations, set status, and send messages without leaving your editor.

**Skills:**
- `/slack` - Check unread, read channels, search messages
- `/setup-slack` - Guided installation of slk CLI

**Type:** Skills-Only (wraps slk Ruby gem)
**Port:** 1:1 skill definitions with comprehensive command reference

---

### 📝 gist
**Port of:** [claude-plugins/gist](https://github.com/ericboehs/claude-plugins/tree/main/plugins/gist)

Effortlessly create and update GitHub Gists with AI-generated documentation. Automatically generates descriptive README comments for your gists using Cursor AI's code understanding.

**Skills:**
- `/gist-create <file>` - Create gist with AI-generated README
- `/gist-update <file>` - Update gist with refreshed docs

**Type:** Skills-Only (wraps gh gist commands)
**Port:** 1:1 skill definitions with AI documentation workflow

---

## Installation

### Quick Install

```bash
cd cursor-plugins
make build        # Build all binaries
make install      # Install to /usr/local/bin
```

### Verify Installation

```bash
which parse-session lint-daemon watch-ci watch-copilot-reviews
```

### Using in Cursor

See [USAGE.md](./USAGE.md) for detailed instructions on:
- Command-line usage
- Setting up hooks for auto-linting
- Adding skills to `.cursorrules`
- Integrating with Cursor AI

## Publishing to Marketplace

See [MARKETPLACE.md](./MARKETPLACE.md) for instructions on publishing these plugins to the official Cursor Marketplace.

## Architecture

```
cursor-plugins/
├── .cursor-plugin/
│   └── marketplace.json       # Marketplace manifest
├── plugins/                   # Individual plugin directories
│   ├── session-improver/
│   ├── code-lint/
│   ├── watch-ci/
│   ├── git-utils/
│   ├── slack/
│   ├── gh-copilot-review/
│   └── gist/
├── cmd/                       # Go command binaries
│   ├── parse-session/
│   ├── lint-daemon/
│   ├── watch-ci/
│   └── watch-copilot-reviews/
├── pkg/                       # Shared Go packages
│   ├── session/
│   └── linter/
└── bin/                       # Built binaries
```

## Requirements

### For Binary Plugins
- Go 1.21+ (for building from source)
- `gh` CLI (GitHub integration - used by watch-ci, gh-copilot-review)
  - Install: `brew install gh` or see [GitHub CLI docs](https://cli.github.com/)
  - Authenticate: `gh auth login`

### For Skill-Only Plugins
- **git-utils**: git and gh CLI
- **slack**: Ruby + slk gem (`gem install slk`) - guided setup via `/setup-slack`
- **gist**: gh CLI (built-in gist support)

### Runtime Dependencies
None of the binaries require external dependencies at runtime. They're statically compiled Go binaries.

## Building

```bash
make build       # Build all binaries
make test        # Run tests
make install     # Install to /usr/local/bin
make clean       # Remove binaries
```

## Contributing

### To the Original Repository
If you want to improve the core functionality, **please contribute to the original repository:**
https://github.com/ericboehs/claude-plugins

### To This Port
For Cursor-specific improvements:
1. Fork this repository
2. Create a feature branch
3. Submit a PR

We'll gladly port improvements from the original back to this Cursor version.

## Credits

**Original Author:** [Eric Boehs](https://github.com/ericboehs)
**Original Repository:** https://github.com/ericboehs/claude-plugins
**Port Maintainer:** [Caboose AI](https://github.com/cxm6467/caboose-ai)

See [CREDITS.md](./CREDITS.md) for detailed attribution.

## License

This port maintains compatibility with the original repository's licensing. Please refer to the original repository for license details.

## Support

- **Original claude-plugins issues:** https://github.com/ericboehs/claude-plugins/issues
- **Cursor-specific issues:** https://github.com/cxm6467/caboose-ai/issues

## Related Projects

- **claude-plugins** (original): https://github.com/ericboehs/claude-plugins
- **Cursor IDE**: https://cursor.com
- **Claude Code CLI**: https://docs.anthropic.com/claude/docs/claude-code

---

**🙏 Thank you Eric Boehs for the original claude-plugins!** If you find these plugins useful, please star the original repository.
