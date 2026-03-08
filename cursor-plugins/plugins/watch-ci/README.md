# watch-ci Plugin

> **Port of:** [ericboehs/claude-plugins](https://github.com/ericboehs/claude-plugins) - See [CREDITS.md](../../CREDITS.md) for attribution details.


Monitor GitHub Actions CI status for the current branch. Automatically watches until all checks pass or fail, then exits.

## Installation

1. Build the binary:
   ```bash
   go build -o bin/watch-ci ./cmd/watch-ci
   ```

2. Install to PATH:
   ```bash
   sudo cp bin/watch-ci /usr/local/bin/
   ```

## Requirements

- **gh** (GitHub CLI) - `brew install gh`
- **jq** - `brew install jq`
- Git repository with a remote on GitHub
- GitHub Actions enabled

## Usage

```bash
# Watch CI for current branch (10s interval)
watch-ci

# Custom polling interval (30 seconds)
watch-ci 30

# Help
watch-ci --help
```

## How It Works

1. Detects the current git branch
2. Checks for an open PR first (using `gh pr checks`)
3. Falls back to workflow runs if no PR exists
4. Polls every N seconds (default: 10s, minimum: 5s)
5. Displays live status with color coding:
   - 🟢 **PASS** - Check succeeded
   - 🔴 **FAIL** - Check failed
   - 🟡 **WAIT** - Check still running
   - 🟡 **SKIP** - Check skipped/cancelled
6. Exits when all checks complete:
   - **Exit code 0** - All checks passed
   - **Exit code 1** - Some checks failed
7. Shows failure logs automatically when checks fail

## Output Example

```
CI Status for branch: feature/new-thing
Last updated: 2026-03-07 20:15:30
Press Ctrl+C to stop watching

CI Check Results:
====================
PASS Build
PASS Test - Unit
WAIT Test - E2E - running...
PASS Lint

Checks still running...
```

## Integration with Cursor

Use in a Cursor skill/rule:

```markdown
---
name: watch-ci
description: Monitor CI status for current branch
tools: Bash
---

Run `watch-ci` and wait for it to exit. Report results to user.

If checks fail, offer to help fix the issues.
```

## Tips

- Run before pushing to production
- Use longer intervals (30-60s) for slow CI pipelines
- Press Ctrl+C to stop watching early
- Works with both PR checks and direct workflow runs

## Troubleshooting

- **"gh not found"**: Install with `brew install gh`
- **"jq not found"**: Install with `brew install jq`
- **"Not in a git repository"**: Run from within a git repo
- **"No CI runs found"**: Push your branch to GitHub first
- **"Not on a branch"**: Checkout a branch (not detached HEAD)

## References

- [GitHub CLI](https://cli.github.com/)
- [GitHub Actions](https://docs.github.com/en/actions)
