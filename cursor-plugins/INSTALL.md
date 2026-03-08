# Installing Cursor Workflow Automation Plugins

## Quick Install (In Cursor IDE)

### Option 1: Install from GitHub
1. Open Cursor IDE
2. In the chat, type:
   ```
   /add-plugin https://github.com/cxm6467/caboose-cursor-tools/tree/main/cursor-plugins
   ```
3. Cursor will download and install all 7 plugins
4. Restart Cursor to activate

### Option 2: Install via Marketplace (Coming Soon)
Once submitted to the official marketplace:
1. Visit https://cursor.com/marketplace
2. Search for "Cursor Workflow Automation Bundle"
3. Click "Install"
4. Restart Cursor

## Manual Installation (For Development/Testing)

### For Cursor IDE Users

**Copy the entire plugin bundle:**
```bash
# Clone the repository
git clone https://github.com/cxm6467/caboose-cursor-tools.git
cd caboose-cursor-tools/cursor-plugins

# Build binaries (for binary-based plugins)
make build

# Install to Cursor's plugins directory
mkdir -p ~/.cursor/plugins
cp -r /path/to/cursor-plugins ~/.cursor/plugins/caboose-cursor-plugins

# Install binaries to PATH
sudo make install

# Restart Cursor
```

**Or install skills only (no binaries):**
```bash
# Copy individual skills to Cursor's skills directory
cp -r cursor-plugins/plugins/gist/skills/* ~/.cursor/skills-cursor/
cp -r cursor-plugins/plugins/git-utils/skills/* ~/.cursor/skills-cursor/
cp -r cursor-plugins/plugins/slack/skills/* ~/.cursor/skills-cursor/
cp -r cursor-plugins/plugins/gh-copilot-review/skills/* ~/.cursor/skills-cursor/
cp -r cursor-plugins/plugins/session-improver/skills/* ~/.cursor/skills-cursor/

# Restart Cursor
```

## Verifying Installation

After restarting Cursor, type `/` in the chat and you should see:

**Available Skills:**
- `/improve-session` - Analyze session for inefficiencies
- `/session-improver` - (alias for improve-session)
- `/commit-and-push` - Semantic git commit workflow
- `/merge-and-cleanup` - PR merge with cleanup
- `/gist-create` - Create gist with AI docs
- `/gist-update` - Update existing gist
- `/slack` - Slack integration
- `/setup-slack` - Slack CLI setup
- `/gh-copilot-review` - Auto-apply Copilot reviews

**Binary Tools** (if installed with `make install`):
```bash
watch-ci                    # Monitor GitHub Actions
parse-session --current     # Analyze current session
lint-daemon                 # Auto-linting (runs via hooks)
watch-copilot-reviews       # Monitor Copilot reviews
```

## Troubleshooting

### Skills don't appear in autocomplete
1. Ensure you restarted Cursor completely (quit and reopen, not just reload)
2. Check permissions: `chmod -R 755 ~/.cursor/skills-cursor/`
3. Verify SKILL.md files exist: `find ~/.cursor/skills-cursor -name "SKILL.md"`
4. Check Cursor logs for errors

### Binary commands not found
1. Ensure binaries are built: `cd cursor-plugins && make build`
2. Install to PATH: `sudo make install`
3. Or add cursor-plugins/bin to your PATH:
   ```bash
   export PATH="$PATH:/path/to/cursor-plugins/bin"
   ```

### Hooks not working
1. Check hooks.json exists: `~/.cursor/plugins/caboose-cursor-plugins/plugins/code-lint/hooks/hooks.json`
2. Ensure lint-daemon is in PATH
3. Restart Cursor to reload hooks

## Platform-Specific Notes

### Linux/WSL
- Default installation paths work as documented
- Use `make build` to compile binaries

### macOS
- Cursor config at: `~/Library/Application Support/Cursor/User/`
- Use `sudo make install` for system-wide binary installation

### Windows
- Install via WSL recommended
- Or use pre-built binaries from GitHub Releases

## Next Steps

Once installed, try:
1. `/improve-session` to analyze your workflow
2. `/commit-and-push` to commit changes with semantic messages
3. `/gist-create <file>` to share code with auto-generated docs
4. `/slack` to check messages without leaving Cursor

For detailed usage of each skill, see the README files in each plugin directory.
