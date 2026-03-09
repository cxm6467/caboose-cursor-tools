# How to Use Cursor Skills (No `/add-plugin` Needed!)

## TL;DR - Skills Are Already Installed!

Your skills are located at `~/.cursor/skills-cursor/`. **There is NO `/add-plugin` command in Cursor** - skills are simply markdown files that Cursor reads automatically.

## ✅ Verify Skills Are Loaded

1. **Open Cursor IDE**
2. **Start a new chat** (or continue existing one)
3. **Type `/`** (just a forward slash)
4. **Look for your skills** in the autocomplete:
   - `/improve-session`
   - `/commit-and-push`
   - `/gist-create`
   - `/slack`
   - etc.

## 🔍 If Skills Don't Appear

### Check 1: Skills Directory
```bash
ls -la ~/.cursor/skills-cursor/
```

You should see these directories:
- `commit-and-push/`
- `gh-copilot-review/`
- `gist-create/`
- `gist-update/`
- `improve-session/`
- `merge-and-cleanup/`
- `session-improver/`
- `setup-slack/`
- `slack/`

### Check 2: Reinstall Skills
```bash
cd /home/caboose/dev/caboose-cursor-tools/cursor-plugins

# Copy all skills to Cursor
cp -r plugins/gist/skills/* ~/.cursor/skills-cursor/
cp -r plugins/git-utils/skills/* ~/.cursor/skills-cursor/
cp -r plugins/slack/skills/* ~/.cursor/skills-cursor/
cp -r plugins/gh-copilot-review/skills/* ~/.cursor/skills-cursor/
cp -r plugins/session-improver/skills/* ~/.cursor/skills-cursor/
```

### Check 3: Restart Cursor
**Completely quit and restart Cursor**:
```bash
# Kill all Cursor processes
pkill -9 cursor

# Then relaunch Cursor from your applications menu
```

## 📚 Available Skills

### Session Analysis
- **`/improve-session`** - Analyze current session for inefficiencies
  - Detects linter loops, repeated file reads, tool failures
  - Recommends hooks, config changes, automations
  - Example: `/improve-session`

- **`/session-improver`** - Alias for `/improve-session`

### Git Workflows
- **`/commit-and-push`** - Semantic commit with auto-push
  - Analyzes changes, generates commit message
  - Follows conventional commits (feat:, fix:, etc.)
  - Example: `/commit-and-push`

- **`/merge-and-cleanup`** - Complete PR merge workflow
  - Merges PR via GitHub
  - Deletes local and remote branches
  - Returns to main/master and pulls
  - Example: `/merge-and-cleanup`

### GitHub Gists
- **`/gist-create <file>`** - Create gist with AI README
  - Creates public gist via gh CLI
  - Generates comprehensive README comment
  - Example: `/gist-create src/utils/helper.ts`

- **`/gist-update <file>`** - Update existing gist
  - Finds gist by filename
  - Updates content and README
  - Example: `/gist-update src/utils/helper.ts`

### Slack Integration
- **`/slack`** - Check messages, search, send DMs
  - Shows unread messages
  - Search conversations
  - Set status
  - Example: `/slack`

- **`/setup-slack`** - First-time slk CLI setup
  - Installs slk gem
  - Configures API tokens
  - Builds cache
  - Example: `/setup-slack`

### GitHub Copilot
- **`/gh-copilot-review`** - Auto-apply Copilot PR reviews
  - Requests @copilot as reviewer
  - Waits for review
  - Applies all suggestions
  - Resolves threads
  - Example: `/gh-copilot-review`

## 🛠️ Binary Tools (Optional)

Some skills use Go binaries for better performance:

```bash
# Build binaries
cd /home/caboose/dev/caboose-cursor-tools/cursor-plugins
make build

# Install to PATH (optional but recommended)
sudo make install
```

**Available binaries:**
- `parse-session --current` - Analyze current Cursor session
- `watch-ci` - Monitor GitHub Actions in real-time
- `lint-daemon` - Auto-linting via hooks
- `watch-copilot-reviews` - Monitor Copilot reviews

## 📍 Where Skills Are Stored

| Type | Location | Purpose |
|------|----------|---------|
| **User skills** | `~/.cursor/skills-cursor/` | Personal skills (used by you) |
| **Project skills** | `.cursor/skills/` | Project-specific skills |
| **Plugin source** | `~/dev/caboose-cursor-tools/cursor-plugins/` | Original source code |

## 🎯 How Skills Actually Work

1. **Cursor scans** `~/.cursor/skills-cursor/` on startup
2. **Finds directories** with `SKILL.md` files inside
3. **Reads YAML frontmatter** to get skill name and description
4. **Makes skills available** via `/skill-name` in chat

**No installation command needed** - it's file-based!

## 🔧 Where to View in Cursor Settings

1. **Open Settings**: `Ctrl+Shift+J` (or `Cmd+Shift+J` on Mac)
2. **Look for**:
   - **Features** section
   - **Tools & MCP Servers** section
   - **Skills** or **Plugins** panel (if available)

Note: Not all Cursor versions show skills in settings UI. The reliable way is typing `/` in chat.

## ❓ Troubleshooting

### Skills don't autocomplete after `/`
1. Ensure SKILL.md files exist: `find ~/.cursor/skills-cursor -name "SKILL.md"`
2. Check permissions: `chmod -R 755 ~/.cursor/skills-cursor/`
3. Restart Cursor completely (quit all processes)

### "Binary not found" errors
1. Build binaries: `cd cursor-plugins && make build`
2. Install to PATH: `sudo make install`
3. Or add to PATH manually: `export PATH="$PATH:/path/to/cursor-plugins/bin"`

### Skills work but no output
- Check that binaries are executable: `chmod +x cursor-plugins/bin/*`
- Ensure dependencies are installed (gh, slk, etc.)

## 🚀 Quick Test

Try this in Cursor chat:

```
/improve-session
```

If it works, you'll see analysis of your current session. If not, follow the "If Skills Don't Appear" section above.

## 📖 Learn More

- Full docs: `cursor-plugins/README.md`
- Installation guide: `cursor-plugins/INSTALL.md`
- Marketplace info: `cursor-plugins/MARKETPLACE.md`
