# Using cursor-plugins in Cursor IDE

The binaries are now installed! Here's how to use these plugins in Cursor.

## Quick Start

### 1. Command-Line Tools (Available Now)

You can run these directly from your terminal:

```bash
# Watch CI status for current branch
watch-ci

# Watch for Copilot review on current PR
watch-copilot-reviews

# Parse a Cursor session
parse-session --current

# Lint a file (used by hooks)
lint-daemon
```

### 2. Setting Up Hooks (Automatic Linting)

Create a hooks configuration file in your project:

```bash
mkdir -p .cursor
```

Create `.cursor/hooks.json`:
```json
{
  "afterFileEdit": {
    "command": "lint-daemon"
  }
}
```

Then configure linting for your project at `~/.cursor/code-lint/<project-hash>/config.json`:

```bash
# Project hash is your project path with slashes replaced by dashes
# Example: /home/caboose/dev/my-project → home-caboose-dev-my-project

mkdir -p ~/.cursor/code-lint/home-caboose-dev-my-project
```

Create `~/.cursor/code-lint/home-caboose-dev-my-project/config.json`:
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
    },
    "go": {
      "enabled": true,
      "linters": [
        {
          "name": "golangci-lint",
          "command": "golangci-lint run {file}",
          "enabled": true
        }
      ]
    }
  }
}
```

### 3. Using Skills in Cursor

Skills are workflow instructions that Cursor's AI can follow. Here's how to use them:

#### Add Skills to Your .cursorrules File

Create or edit `.cursorrules` in your project:

```markdown
# Project Rules

## Available Skills

### Git Workflows

#### /commit-and-push
Stage all changes, commit with semantic commit message, and push to origin.

Usage: When I say "commit and push" or "/commit-and-push":
1. Run `git status` to see changes
2. Stage all files with `git add .`
3. Ask about any files that shouldn't be committed
4. Create a semantic commit message (e.g., "feat:", "fix:", "refactor:")
5. Commit with the message
6. Push to origin

#### /merge-and-cleanup
Merge current PR (squash), delete branch, switch to main, pull, clean up worktree if applicable.

Usage: When I say "merge and cleanup" or "/merge-and-cleanup":
1. Check current branch and PR status
2. Squash merge PR via `gh pr merge --squash --delete-branch`
3. Switch to main/master branch
4. Pull latest changes
5. Clean up worktree if used

### Slack Integration

#### /slack
Check unread messages, read channels, search, set status.

Usage: When I say "/slack" or "check slack":
- Use `slk unread` to show unread messages
- Use `slk messages #channel` to read a channel
- Use `slk search "query"` to search messages

Requires: `gem install slk` and `slk config setup`

#### /setup-slack
Install and configure slk CLI.

### GitHub Copilot Review

#### /gh-copilot-review
Wait for Copilot review, apply all feedback, push fixes.

Usage: When I say "/gh-copilot-review" or "wait for copilot":
1. Run `watch-copilot-reviews` and wait for review
2. Fetch all inline comments
3. Apply suggested changes
4. Reply to each comment
5. Resolve review threads
6. Commit and push fixes

### GitHub Gists

#### /gist-create
Create a gist from a file with auto-generated README.

Usage: When I say "/gist-create <file>" or "create a gist":
1. Read the file
2. Create gist: `gh gist create <file> --desc "..." --public`
3. Generate comprehensive README comment
4. Post README as gist comment

#### /gist-update
Update an existing gist and its README.

Usage: When I say "/gist-update <file>":
1. Find existing gist by filename
2. Update gist: `gh gist edit <gist_id> <file>`
3. Update README comment

### Session Analysis

#### /improve-session
Analyze current Cursor session for inefficiencies.

Usage: When I say "/improve-session":
1. Run `parse-session --current`
2. Analyze output for linter loops, tool failures, token waste
3. Suggest improvements
```

### 4. Integrating with Cursor AI

When chatting with Cursor's AI, you can now:

**Direct Commands:**
```
You: "watch ci for this branch"
AI: *runs watch-ci and reports results*

You: "wait for copilot review"
AI: *runs watch-copilot-reviews, waits, then applies feedback*

You: "commit and push these changes"
AI: *follows /commit-and-push workflow*
```

**Using Skills:**
```
You: "/merge-and-cleanup"
AI: *follows merge-and-cleanup workflow*

You: "/slack"
AI: *checks your unread Slack messages*

You: "/gist-create ./script.sh"
AI: *creates gist with auto-generated README*
```

## Plugin-Specific Setup

### session-improver
```bash
# Analyze current session
parse-session --current

# Analyze specific session
parse-session <session-id>
```

### code-lint
1. Set up hooks (see above)
2. Configure linters per project
3. Edit files in Cursor - linting happens automatically

### watch-ci
```bash
# Watch current branch
watch-ci

# Custom interval (30 seconds)
watch-ci 30
```

### git-utils
- Use via Cursor chat: "commit and push" or "merge and cleanup"
- Add skills to .cursorrules (see above)

### slack
```bash
# First-time setup
gem install slk
slk config setup

# Then use via Cursor: "check slack" or "/slack"
```

### gh-copilot-review
```bash
# Command line
watch-copilot-reviews

# Or via Cursor: "/gh-copilot-review"
```

### gist
- Use via Cursor: "/gist-create ./file.sh" or "/gist-update ./file.sh"
- Requires: `gh` CLI authenticated

## Project Structure for Reference

```
cursor-plugins/
├── bin/                          # Built binaries
│   ├── lint-daemon
│   ├── parse-session
│   ├── watch-ci
│   └── watch-copilot-reviews
├── plugins/                      # Plugin definitions
│   ├── code-lint/
│   ├── gh-copilot-review/
│   ├── gist/
│   ├── git-utils/
│   ├── session-improver/
│   ├── slack/
│   └── watch-ci/
└── cmd/                          # Go source code
    ├── lint-daemon/
    ├── parse-session/
    ├── watch-ci/
    └── watch-copilot-reviews/
```

## Tips

1. **Hooks are automatic** - Once configured, `lint-daemon` runs on every file edit
2. **Skills need .cursorrules** - Copy skill definitions to your project's `.cursorrules`
3. **Binaries work standalone** - All command-line tools work independently
4. **AI integration** - Reference skills in chat: "follow the commit-and-push workflow"
5. **Combine tools** - Chain workflows: "commit and push, then watch ci"

## Troubleshooting

**Hooks not firing?**
- Check `.cursor/hooks.json` exists in project root
- Verify `lint-daemon` is in PATH: `which lint-daemon`
- Check linter config exists at `~/.cursor/code-lint/<project-hash>/config.json`

**Skills not working?**
- Add skill definitions to `.cursorrules`
- Reference skills explicitly in chat
- Check required CLI tools are installed (gh, slk, etc.)

**Binaries not found?**
```bash
# Verify installation
which parse-session
which lint-daemon
which watch-ci
which watch-copilot-reviews

# Reinstall if needed
cd cursor-plugins
make install
```

## Next Steps

1. Configure hooks for automatic linting
2. Add skill definitions to your `.cursorrules`
3. Try a workflow: "watch ci" or "commit and push"
4. Set up Slack integration if needed
5. Explore session analysis with `parse-session --current`
