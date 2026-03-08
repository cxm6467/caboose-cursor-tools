---
name: gist-create
description: |
  Create GitHub Gists with AI-generated documentation in one command. Takes any file and:
  1. Reads and analyzes the code
  2. Creates a public gist via gh CLI
  3. Generates a comprehensive README comment (title, features, installation, usage, examples)
  4. Posts the README as the first gist comment
  5. Returns the shareable gist URL

  Use when user says:
  - "/gist-create <file>"
  - "create a gist from this file"
  - "gist this file"
  - "share this as a gist"

  BENEFIT: Automatic documentation generation means your gists are immediately useful to others without manual README writing. Perfect for sharing code snippets, utilities, or examples.

  Optional: Add "--private" for private gists (default is public).
tools: Bash, Read
---

# Create GitHub Gist with AI Documentation

Automatically creates a GitHub Gist with intelligently generated README documentation based on code analysis.

## Steps

1. Read the file to understand what it does
2. Create a gist using: `gh gist create <filepath> --desc "<description>" --public`
   - If the user requests a private gist, omit the `--public` flag (gists default to private)
3. Generate a comprehensive README as a gist comment using the GitHub API:
   ```
   gh api -X POST /gists/<gist_id>/comments -f body='<markdown_readme>'
   ```
4. The README comment should include:
   - Title and description
   - Features list
   - Installation instructions (curl from raw gist URL)
   - Dependencies
   - Usage examples
   - Configuration details (if applicable)
5. Report the gist URL to the user when complete

Note: Use single quotes around the body and escape any single quotes in the content with '\''
