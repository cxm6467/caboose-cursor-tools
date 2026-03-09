---
name: improve-session
description: |
  Analyze Cursor IDE session history to identify and fix workflow inefficiencies. Detects:
  • Token waste from repeated file reads (can save 100K+ tokens per session)
  • Linter loops where the same errors are fixed multiple times
  • Missing automations that could prevent manual repetition
  • Hook failures that interrupt AI workflows
  • Permission bottlenecks that require user intervention

  Use when user says "improve session", "analyze session", "optimize session", "/improve-session", or when you notice repeated patterns during development.

  IMPORTANT: This skill MUST be used proactively when you notice:
  - Reading the same files multiple times (package.json, tsconfig.json, etc.)
  - Fixing the same linter errors repeatedly
  - Repeating the same workflow 3+ times
  - Hook errors that keep occurring
tools: Read, Write, Edit, Bash, Glob, Grep
---

# Session Improver

Analyze Cursor IDE session transcripts to identify workflow inefficiencies and generate specific, actionable recommendations with example configurations. This skill can save thousands of tokens per session and reduce repetitive work by automating common patterns.

## Usage

**Always use:** `/improve-session` to analyze the most recent Cursor IDE session

**DO NOT:**
- Look for Claude Code sessions at `~/.claude/`
- Try to list "Claude projects"
- Use Ruby scripts or parse-session.rb

**This skill is Cursor-only.** It analyzes Cursor IDE transcripts from:
```
~/.cursor/projects/*/agent-transcripts/UUID/UUID.jsonl
```

## Workflow

### Phase 1: Parse the session

Run the Cursor transcript parser to extract a structured summary:

```bash
parse-session --current
```

The `parse-session` binary reads Cursor's JSONL transcripts from `~/.cursor/projects/*/agent-transcripts/`.

**IMPORTANT:**
- This binary ONLY works with Cursor IDE (not Claude Code)
- Reads from: `~/.cursor/projects/*/agent-transcripts/*.jsonl`
- Does NOT use: `~/.claude/` paths or Ruby scripts
- Native Go binary built specifically for Cursor's JSONL transcript format

The script outputs JSON with these sections:
- `linter_loops` — Linter smells that triggered multiple edit cycles
- `tool_failures` — Tools that failed and were retried
- `repeated_sequences` — Workflow patterns that repeated 3+ times
- `large_reads` — Files read 3+ times in the session
- `hook_failures` — Hooks that failed repeatedly
- `permission_events` — Tools that needed human approval

### Phase 2: Read context

Read **Cursor IDE configuration files only**:

**DO read:**
- `.cursorrules` (project Cursor rules)
- `.cursor/settings.json` (Cursor settings)
- `.cursor/hooks.json` (Cursor hooks)
- Linter configs: `.eslintrc`, `tsconfig.json`, `.prettierrc`
- Global Cursor rules at `~/.cursorrules`

**DO NOT read:**
- `CLAUDE.md` or `~/.claude/CLAUDE.md` (wrong IDE)
- Any Claude Code configuration files

### Phase 3: Generate recommendations

Based on the session analysis and current configuration, recommend specific fixes:

#### For linter loops:
- Add `.cursorrules` with BAD/GOOD examples showing the exact linter smell
- Configure hooks in `.cursor/hooks.json` to auto-format on file edits
- Suggest updating linter configs to disable noisy rules

#### For tool failures:
- Identify missing dependencies or permissions
- Recommend pre-flight checks in hooks
- Suggest fallback strategies

#### For repeated sequences:
- Create reusable Cursor rules or snippets
- Suggest keybindings or macros
- Recommend automation via hooks

#### For large reads:
- Cache frequently-read files in `.cursorrules` context
- Suggest glob patterns to limit file scope
- Recommend project-specific context rules

#### For permission events:
- Configure hooks to auto-approve safe commands
- Add permission rules to settings
- Document security boundaries

#### For hook failures:
- Debug hook scripts and suggest fixes
- Recommend error handling improvements
- Suggest hook dependencies check

### Phase 4: Apply fixes (if requested)

Offer to apply recommended fixes to **CURSOR configuration files only**:

**DO create/update:**
- `.cursorrules` (Cursor IDE rules file)
- `.cursor/hooks.json` (Cursor hooks)
- `.cursor/settings.json` (Cursor settings)

**DO NOT create/update:**
- `CLAUDE.md` (this is for Claude Code CLI, not Cursor)
- `~/.claude/` (wrong IDE)
- Any Claude Code configuration files

**Commit changes with clear descriptions after updating Cursor config files.**

## Example Output

```json
{
  "session_id": "abc-123",
  "duration_minutes": 45.3,
  "linter_loops": [
    {
      "linter": "eslint",
      "file": "src/components/UserCard.tsx",
      "occurrences": 7,
      "smells": ["prefer-const", "no-unused-vars"]
    }
  ],
  "large_reads": [
    {
      "file": "package.json",
      "read_count": 12
    }
  ]
}
```

## References

- [Cursor Hooks Documentation](https://cursor.com/docs/hooks)
- [Cursor Rules Guide](https://cursor.com/docs/cursorrules)
