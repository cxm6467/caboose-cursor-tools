# Code Lint Plugin

> **Port of:** [ericboehs/claude-plugins](https://github.com/ericboehs/claude-plugins) - See [CREDITS.md](../../CREDITS.md) for attribution details.


Multi-language linting via Cursor hooks. Automatically runs configured linters when files are edited.

## Supported Languages & Linters

| Language   | Linters |
|------------|---------|
| Ruby       | rubocop, reek, brakeman, standardrb |
| JavaScript/TS | eslint, biome, prettier |
| Python     | ruff, mypy, flake8, black |
| Go         | golangci-lint, gofmt |
| Rust       | clippy, rustfmt |
| Markdown   | markdownlint, prettier |
| HTML/CSS   | htmlhint, prettier |
| Shell      | shellcheck |
| Any        | semgrep |

## Setup

1. Build the lint daemon:
   ```bash
   go build -o bin/lint-daemon ./cmd/lint-daemon
   ```

2. Install to PATH:
   ```bash
   sudo cp bin/lint-daemon /usr/local/bin/
   ```

3. Configure Cursor hooks in `.cursor/hooks.json`:
   ```json
   {
     "afterFileEdit": {
       "command": "lint-daemon"
     }
   }
   ```

4. Create project config at `~/.cursor/code-lint/<project-hash>/config.json`:
   ```json
   {
     "languages": {
       "javascript": {
         "enabled": true,
         "linters": [
           {
             "name": "eslint",
             "command": "eslint {file}",
             "enabled": true,
             "timeout": "5s"
           }
         ]
       },
       "python": {
         "enabled": true,
         "linters": [
           {
             "name": "ruff",
             "command": "ruff check {file}",
             "enabled": true
           }
         ]
       }
     }
   }
   ```

## How It Works

1. When you edit/write a file in Cursor, the `afterFileEdit` hook triggers
2. `lint-daemon` reads the file path from hook input (stdin)
3. Detects the file's language by extension or shebang
4. Loads config from `~/.cursor/code-lint/<project-hash>/config.json`
5. Runs all enabled linters for that language
6. Outputs results to Cursor

## Configuration

Config location: `~/.cursor/code-lint/<project-hash>/config.json`

Project hash is the project path with slashes replaced by dashes:
- `/Users/foo/Code/bar` → `Users-foo-Code-bar`

### Config Schema

```json
{
  "languages": {
    "<language>": {
      "enabled": true|false,
      "linters": [
        {
          "name": "linter-name",
          "command": "linter-cmd {file}",
          "enabled": true|false,
          "timeout": "5s"
        }
      ]
    }
  },
  "tool_linters": [
    {
      "name": "semgrep",
      "command": "semgrep --config auto .",
      "enabled": true
    }
  ]
}
```

## Examples

### TypeScript + ESLint

```json
{
  "languages": {
    "javascript": {
      "enabled": true,
      "linters": [
        {
          "name": "eslint",
          "command": "eslint --format compact {file}",
          "enabled": true
        },
        {
          "name": "prettier",
          "command": "prettier --check {file}",
          "enabled": false
        }
      ]
    }
  }
}
```

### Python + Ruff

```json
{
  "languages": {
    "python": {
      "enabled": true,
      "linters": [
        {
          "name": "ruff",
          "command": "ruff check {file}",
          "enabled": true
        },
        {
          "name": "mypy",
          "command": "mypy {file}",
          "enabled": false
        }
      ]
    }
  }
}
```

## Troubleshooting

- **No linting happening**: Check that `lint-daemon` is in PATH and `.cursor/hooks.json` is configured
- **Config not found**: Run setup or manually create config at the correct path
- **Linter not running**: Ensure linter is installed and command is correct in config
- **Permission denied**: Make sure `lint-daemon` has execute permissions

## References

- [Cursor Hooks Documentation](https://cursor.com/docs/hooks)
