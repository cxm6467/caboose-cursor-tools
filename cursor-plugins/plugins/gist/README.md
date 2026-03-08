# gist

> **Port of:** [ericboehs/claude-plugins](https://github.com/ericboehs/claude-plugins) - See [CREDITS.md](../../CREDITS.md) for attribution details.


Create and update GitHub Gists with auto-generated README comments.

## Usage

- `/gist-create <filepath>` — Create a public gist with a README comment
- `/gist-update <filepath>` — Update an existing gist and its README comment

## Features

- Creates public gists by default (pass "private" for secret gists)
- Auto-generates a comprehensive README as the first gist comment
- Updates existing gists by matching filename
- README includes install instructions, usage examples, and dependencies

## Requirements

- `gh` (GitHub CLI), authenticated

## Installation

Install the GitHub CLI:
```bash
brew install gh
gh auth login
```

## Integration with Cursor

These skills can be invoked via Cursor's command palette or by mentioning them in chat:
- "create a gist"
- "update the gist"
- "/gist-create"
- "/gist-update"

The AI will create or update gists with automatically generated README comments.

## References

- [GitHub CLI](https://cli.github.com/)
- [GitHub Gists](https://docs.github.com/en/github/writing-on-github/editing-and-sharing-content-with-gists)
