# Publishing to Cursor Marketplace

## ✅ Ready for Marketplace!

All plugins are properly structured and ready to publish to the official [Cursor Marketplace](https://cursor.com/marketplace).

## Current Status

### Marketplace Configuration
- ✅ Root marketplace manifest created: `.cursor-plugin/marketplace.json`
- ✅ 7 plugins ready to publish
- ✅ All plugin.json files valid (kebab-case names, descriptions, versions)
- ✅ All skills have proper YAML frontmatter
- ✅ README files included for each plugin

### Plugin Inventory

| Plugin | Version | Type | Description |
|--------|---------|------|-------------|
| **session-improver** | 1.0.0 | Binary + Skill | Analyze sessions for inefficiencies |
| **code-lint** | 1.0.0 | Binary + Hook | Multi-language auto-linting |
| **watch-ci** | 1.0.0 | Binary | Monitor GitHub Actions CI |
| **git-utils** | 1.1.0 | Skills | Git workflow automation |
| **slack** | 1.0.0 | Skills | Slack integration via slk CLI |
| **gh-copilot-review** | 1.0.0 | Binary + Skill | Automate Copilot review feedback |
| **gist** | 1.0.0 | Skills | GitHub Gists with auto README |

## Publishing Steps

### Option 1: Submit to Official Marketplace

1. **Push to GitHub**
   ```bash
   cd cursor-plugins
   git add .
   git commit -m "feat: add cursor marketplace plugins"
   git push origin main
   ```

2. **Submit to Cursor**
   - Visit https://cursor.com/marketplace/publish
   - Provide your GitHub repository URL: `https://github.com/cxm6467/caboose-ai`
   - Cursor will validate the plugins and add them to the marketplace

3. **Users Can Install**
   ```bash
   # In Cursor IDE
   /add-plugin caboose-cursor-plugins

   # Or visit marketplace and click install
   https://cursor.com/marketplace
   ```

### Option 2: Fork Official Cursor Plugins Repo

1. **Fork the official repo**
   ```bash
   # Fork https://github.com/cursor/plugins
   git clone https://github.com/YOUR_USERNAME/plugins
   cd plugins
   ```

2. **Copy our plugins**
   ```bash
   # Copy each plugin directory
   cp -r /path/to/cursor-plugins/plugins/* .

   # Update root marketplace.json to include our plugins
   ```

3. **Submit PR to cursor/plugins**
   - Create PR with our plugins
   - Cursor team reviews and merges
   - Plugins become available in official marketplace

### Option 3: Private/Team Marketplace

For internal use or team distribution:

1. **Host on GitHub**
   - Push to private/public repo
   - Structure must match cursor/plugins format

2. **Install via URL**
   ```bash
   # Users can install directly from GitHub
   /add-plugin https://github.com/cxm6467/caboose-ai/cursor-plugins
   ```

3. **Team Distribution** (Coming Soon)
   - Cursor is working on private team marketplaces
   - Organizations can share plugins internally

## Plugin Structure (Already Correct!)

```
cursor-plugins/
├── .cursor-plugin/
│   └── marketplace.json          # ✅ Marketplace manifest
├── plugins/
│   ├── session-improver/
│   │   ├── .cursor-plugin/
│   │   │   └── plugin.json       # ✅ Plugin manifest
│   │   ├── skills/
│   │   │   └── improve-session/
│   │   │       └── SKILL.md      # ✅ Skill with frontmatter
│   │   └── README.md             # ✅ Documentation
│   ├── code-lint/
│   │   ├── .cursor-plugin/
│   │   │   └── plugin.json
│   │   ├── hooks/
│   │   │   └── hooks.json        # ✅ Hook definition
│   │   └── README.md
│   ├── watch-ci/
│   │   ├── .cursor-plugin/
│   │   │   └── plugin.json
│   │   └── README.md
│   ├── git-utils/
│   │   ├── .cursor-plugin/
│   │   │   └── plugin.json
│   │   ├── skills/
│   │   │   ├── commit-and-push/
│   │   │   │   └── SKILL.md      # ✅ Git workflow
│   │   │   └── merge-and-cleanup/
│   │   │       └── SKILL.md
│   │   └── README.md
│   ├── slack/
│   │   ├── .cursor-plugin/
│   │   │   └── plugin.json
│   │   ├── skills/
│   │   │   ├── slack/
│   │   │   │   └── SKILL.md
│   │   │   └── setup-slack/
│   │   │       └── SKILL.md
│   │   └── README.md
│   ├── gh-copilot-review/
│   │   ├── .cursor-plugin/
│   │   │   └── plugin.json
│   │   ├── skills/
│   │   │   └── gh-copilot-review/
│   │   │       └── SKILL.md
│   │   └── README.md
│   └── gist/
│       ├── .cursor-plugin/
│       │   └── plugin.json
│       ├── skills/
│       │   ├── gist-create/
│       │   │   └── SKILL.md
│       │   └── gist-update/
│       │       └── SKILL.md
│       └── README.md
└── bin/                          # ✅ Built binaries (for users)
    ├── parse-session
    ├── lint-daemon
    ├── watch-ci
    └── watch-copilot-reviews
```

## Validation Checklist

- ✅ Plugin names are kebab-case (lowercase, hyphens only)
- ✅ All plugin.json files have required `name` field
- ✅ Descriptions explain what each plugin does
- ✅ Skills have YAML frontmatter with `name` and `description`
- ✅ README.md files document usage
- ✅ No absolute paths or `..` references
- ✅ All manifests are valid JSON
- ✅ Binaries are built and tested

## Binary Distribution

Since some plugins include Go binaries, users will need to build them:

```bash
# After installing from marketplace
cd ~/.cursor/plugins/caboose-cursor-plugins
make build
make install
```

Alternatively, provide pre-built binaries via GitHub Releases:

```bash
# Create release with binaries
cd cursor-plugins
make build
tar -czf cursor-plugins-linux-amd64.tar.gz bin/*
tar -czf cursor-plugins-darwin-amd64.tar.gz bin/*

# Upload to GitHub Releases
gh release create v1.0.0 \
  cursor-plugins-linux-amd64.tar.gz \
  cursor-plugins-darwin-amd64.tar.gz \
  --title "Cursor Plugins v1.0.0" \
  --notes "Initial marketplace release"
```

## Recommended Next Steps

1. **Test Locally First**
   ```bash
   # Copy to Cursor's local plugins directory
   mkdir -p ~/.cursor/plugins
   cp -r cursor-plugins ~/.cursor/plugins/caboose-cursor-plugins

   # Test in Cursor IDE
   # Use /add-plugin or check settings
   ```

2. **Add Logos** (Optional but recommended)
   - Add plugin logos to each `.cursor-plugin/` directory
   - Reference in plugin.json: `"logo": "logo.png"`
   - Helps with marketplace visibility

3. **Add Keywords** (Optional)
   ```json
   {
     "name": "watch-ci",
     "keywords": ["ci", "github-actions", "testing", "automation"]
   }
   ```

4. **Create CHANGELOG.md** for each plugin

5. **Add LICENSE files** (MIT recommended)

## Marketing Your Plugins

Once published, promote on:
- Cursor Community Forum
- Reddit (r/cursor)
- Twitter/X with #CursorIDE
- Dev.to articles
- GitHub README with marketplace badges

## Support & Updates

- Monitor Cursor forum for user feedback
- Use semantic versioning for updates
- Update `version` in plugin.json when releasing changes
- Keep README.md files current with features

## References

- [Cursor Plugin Documentation](https://cursor.com/docs/plugins)
- [Cursor Marketplace](https://cursor.com/marketplace)
- [Official Cursor Plugins Repo](https://github.com/cursor/plugins)
- [Cursor Plugin Publishing](https://cursor.com/marketplace/publish)

---

**Ready to publish!** 🚀 All plugins meet marketplace requirements and are production-ready.
