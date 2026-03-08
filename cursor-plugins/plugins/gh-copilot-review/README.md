# gh-copilot-review

> **Port of:** [ericboehs/claude-plugins](https://github.com/ericboehs/claude-plugins) - See [CREDITS.md](../../CREDITS.md) for attribution details.


Wait for GitHub Copilot PR reviews, address feedback, resolve threads, and push fixes — all in one command.

## Usage

```
/gh-copilot-review
```

## What it does

1. **Watches for Copilot review** — polls until Copilot submits its review (auto-requests @copilot as reviewer if needed)
2. **Fetches inline comments** — extracts all Copilot review comments with file paths, line numbers, and suggestion blocks
3. **Addresses each comment** — applies suggested changes or makes appropriate fixes
4. **Replies to comments** — posts a reply on each comment explaining what was done
5. **Resolves review threads** — marks all addressed threads as resolved via GraphQL
6. **Commits and pushes** — stages changes, commits, and pushes to origin

## Requirements

- `gh` (GitHub CLI) — authenticated
- `jq` — for JSON parsing
- A PR must exist for the current branch

## Binaries

- `watch-copilot-reviews` — standalone binary that monitors Copilot review status

### Building

```bash
go build -o bin/watch-copilot-reviews ./cmd/watch-copilot-reviews
```

### Installation

```bash
sudo cp bin/watch-copilot-reviews /usr/local/bin/
```

### Usage

```bash
# Watch current branch's PR
watch-copilot-reviews

# Watch specific PR number
watch-copilot-reviews 42

# Custom polling interval
watch-copilot-reviews 42 30
```

## Integration with Cursor

This skill can be invoked via Cursor's command palette or by mentioning it in chat:
- "copilot review"
- "wait for copilot"
- "/gh-copilot-review"

The AI will wait for Copilot's review, apply all suggested changes, and push fixes automatically.

## References

- [GitHub CLI](https://cli.github.com/)
- [GitHub Copilot](https://github.com/features/copilot)
