---
name: session-improver
description: |
  Analyze Cursor IDE session history to identify and fix workflow inefficiencies. Detects token waste from repeated file reads, linter loops, missing automations, hook failures, permission bottlenecks. Use when user says "session improver", "improve session", "analyze session", "optimize session", "/session-improver", or "/improve-session".
tools: Read, Write, Edit, Bash, Glob, Grep
---

# Session Improver

Same workflow as the **improve-session** skill. Analyzes Cursor session history for inefficiencies and recommends fixes. Use `/session-improver` or `/improve-session`.

## Usage

- `/session-improver` or `/improve-session` — Analyze current/recent session
- `/session-improver <session-id>` or `parse-session <session-id>` — Specific session

## Workflow

### Phase 1: Parse the session

```bash
parse-session --current
```

Or `parse-session <session-id>`. The `parse-session` binary should be in PATH (from this plugin’s `make install`). Output includes: `linter_loops`, `tool_failures`, `repeated_sequences`, `large_reads`, `hook_failures`, `permission_events`. If not available, infer from the current chat.

### Phase 2: Read context

Check project config: `.cursorrules`, `.cursor/hooks.json`, `.cursor/settings.json`, linter configs (e.g. `.eslintrc`, `tsconfig.json`).

### Phase 3: Recommend fixes

- **Linter loops** → `.cursorrules` BAD/GOOD examples, hooks for format-on-edit, or relax noisy rules
- **Repeated reads** → Cache key files in rules or narrow globs
- **Repeated sequences** → Rules, snippets, or hooks
- **Tool/hook failures** → Dependencies, permissions, fallbacks
- **Permission events** → Auto-approve safe commands, document boundaries

### Phase 4: Apply (if requested)

Offer to create/update `.cursorrules`, `.cursor/hooks.json`, `.cursor/settings.json`.

## References

- [Cursor Hooks](https://cursor.com/docs/hooks)
- [Cursor Rules](https://cursor.com/docs/cursorrules)
