# Slack

> **Port of:** [ericboehs/claude-plugins](https://github.com/ericboehs/claude-plugins) - See [CREDITS.md](../../CREDITS.md) for attribution details.


Slack messaging, status, and search using [`slk`](https://github.com/ericboehs/slk) CLI.

## Skills

### `/slack`

Check unread, read messages, search, set status/DND, and browse activity.

- `/slack` — Show unread messages across all workspaces
- `/slack messages #channel` — Read recent channel messages
- `/slack messages #channel --since 1d` — Read recent messages from the last day
- `/slack search "query"` — Search messages
- `/slack preset focus` — Apply a status preset

### `/setup-slack`

Guided installation and configuration of `slk` CLI, including workspace authentication and cache setup.

## Features

- **Unread:** Check unread across all workspaces, filter by workspace, include muted channels
- **Messages:** Read channels, DMs, and threads by URL; filter by count, duration, or date
- **Search:** Full-text search with filters for channel, user, and date range (requires user token)
- **Activity:** Browse reactions, mentions, and thread replies
- **Status:** Set/clear status with emoji and duration, manage DND, apply presets
- **Workspaces:** Multi-workspace support with `-w` flag or `--all`

## Requirements

- Ruby (for gem install)
- [`slk`](https://github.com/ericboehs/slk) — install via `gem install slk`
- Slack tokens configured via `slk config setup`
- User token (xoxc/xoxs) recommended for full access including search

## Installation

1. Install slk CLI:
   ```bash
   gem install slk
   ```

2. Configure workspaces:
   ```bash
   slk config setup
   ```

3. Build cache:
   ```bash
   slk cache build --all
   ```

## Integration with Cursor

These skills can be invoked via Cursor's command palette or by mentioning them in chat:
- "check slack"
- "setup slack"
- "/slack"
- "/setup-slack"

The AI will use the `slk` CLI to interact with Slack on your behalf.

## References

- [slk CLI](https://github.com/ericboehs/slk)
- [Slack API](https://api.slack.com/)
