# Credits & Attribution

## Original Work

These plugins are **direct ports** of the excellent [claude-plugins](https://github.com/ericboehs/claude-plugins) repository by **Eric Boehs**.

### Original Author
- **Name:** Eric Boehs
- **GitHub:** [@ericboehs](https://github.com/ericboehs)
- **Original Repository:** https://github.com/ericboehs/claude-plugins
- **License:** See original repository

## What Changed in This Port

This repository adapts Eric Boehs' claude-plugins for use with **Cursor IDE** instead of Claude Code CLI:

### Structural Changes
- Ruby scripts → Go binaries (for performance and cross-platform compatibility)
- Bash scripts → Go binaries where applicable
- `.claude-plugin` → `.cursor-plugin` format
- `CLAUDE.md` references → `.cursorrules`
- Claude Code CLI paths → Cursor IDE paths

### Specific Adaptations

| Original (claude-plugins) | Port (cursor-plugins) | Changes |
|---------------------------|----------------------|---------|
| **session-improver** | **session-improver** | Ruby JSONL parser → Go SQLite parser for Cursor's session DB |
| **code-lint** | **code-lint** | Bash hook script → Go lint-daemon binary |
| **watch-ci** | **watch-ci** | Bash script → Go binary with color output |
| **git-utils** | **git-utils** | 1:1 port of skills, no code changes |
| **slack** | **slack** | 1:1 port of skills, wraps slk CLI |
| **gh-copilot-review** | **gh-copilot-review** | Bash script → Go binary + skill port |
| **gist** | **gist** | 1:1 port of skills |

### What Stayed the Same
- ✅ All workflow logic and behavior
- ✅ Skill definitions and instructions
- ✅ Command-line interfaces and user experience
- ✅ Documentation and usage patterns
- ✅ Tool requirements (gh CLI, jq, etc.)

## Why Port to Go?

1. **Cross-platform compatibility** - Single binary works on Linux, macOS, Windows
2. **No Ruby dependency** - Users don't need Ruby/gems installed
3. **Better performance** - Compiled binaries are faster than interpreted scripts
4. **Easier distribution** - Single executable vs managing script dependencies
5. **Cursor integration** - Native support for Cursor's SQLite session storage

## Credit Where Credit Is Due

**All original ideas, workflows, and designs** belong to Eric Boehs. This port simply:
- Translates implementation from Ruby/Bash to Go
- Adapts paths and formats for Cursor IDE
- Maintains the original's excellent UX and functionality

If you find these plugins useful, please ⭐ star the original repository:
**https://github.com/ericboehs/claude-plugins**

## License

This port maintains compatibility with the original repository's licensing. Please refer to:
- Original repository: https://github.com/ericboehs/claude-plugins
- Individual plugin licenses as specified by the original author

## Contributing

If you want to improve these plugins:
1. Consider contributing to the **original claude-plugins** repository first
2. We'll gladly port improvements back to this Cursor version
3. PRs welcome for Cursor-specific enhancements

## Thank You

**Eric Boehs** - For creating the original claude-plugins and pioneering AI coding workflows

**Claude Code Team** - For the plugin architecture that inspired this work

**Cursor Team** - For building an amazing AI IDE and plugin marketplace

---

*This port was created with permission and respect for the original author's work.*
