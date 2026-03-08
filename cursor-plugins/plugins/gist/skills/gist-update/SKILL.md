---
name: gist-update
description: |
  Update existing GitHub Gists with refreshed code and documentation. Intelligent workflow:
  1. Finds existing gist by filename (or accepts gist ID/URL)
  2. Updates the gist with new file content
  3. Regenerates README comment based on code changes
  4. Updates or creates the README comment
  5. Returns the gist URL

  Use when user says:
  - "/gist-update <file>"
  - "update the gist"
  - "sync this gist"
  - "refresh the gist for <file>"

  SMART FEATURES:
  - Auto-finds existing gist by filename (searches last 100 gists)
  - Preserves gist URL (no new link needed)
  - Updates README to reflect new features/changes
  - Maintains gist history

  Optional: Provide gist ID explicitly: "/gist-update <file> <gist-id>"
tools: Bash, Read, AskUserQuestion
---

# Update GitHub Gist with Refreshed Documentation

Synchronizes local file changes to an existing gist and regenerates documentation based on code updates.

Arguments can be:
- Just a file path (will search for existing gist by filename)
- A file path and gist ID/URL

## Steps

1. If no gist ID provided, try to find an existing gist:
   - `gh gist list --limit 100` and search for the filename
   - Or ask the user for the gist ID/URL
2. Update the gist file: `gh gist edit <gist_id> <filepath>`
3. Update the README comment if it exists:
   - Get the comment ID: `gh api /gists/<gist_id>/comments --jq '.[0].id'`
   - Update it: `gh api -X PATCH /gists/<gist_id>/comments/<comment_id> -f body='<updated_readme>'`
   - If no comment exists, create one with the README
4. The README should reflect any new features or changes in the updated file
5. Report the gist URL when complete

Note: Use single quotes around the body and escape any single quotes in the content with '\''
