---
name: commit-and-push
description: |
  Intelligent git commit workflow with semantic commit messages. Stages all changes, asks for confirmation on sensitive files, groups related changes, generates conventional commit messages (feat:, fix:, docs:, refactor:, etc.), and pushes to origin.

  Use when user says:
  - "/commit-and-push"
  - "commit and push"
  - "push this up"
  - "commit these changes"
  - "push my code"

  BEHAVIOR:
  - Auto-stages all modified and new files
  - Prompts for confirmation if sensitive files detected (.env, credentials, secrets)
  - Suggests splitting into multiple commits if changes are unrelated
  - Generates semantic commit messages following conventional commit format
  - Pushes to origin automatically
  - User is EXPLICITLY consenting to all git operations

  IMPORTANT: Do NOT chain git commands with && or ;. Run each command separately (git add, git commit, git push) as individual Bash tool calls so that prompts don't cascade.
tools: Bash, AskUserQuestion
---

# Commit and Push with Semantic Messages

Automates the commit-and-push workflow with intelligent file staging, semantic commit message generation, and safety confirmations. Follows git best practices and conventional commit standards.
