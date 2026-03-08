---
name: gh-copilot-review
description: |
  Fully automated GitHub Copilot code review workflow. Handles the complete cycle:
  1. Auto-requests @copilot as reviewer if not already requested
  2. Waits for Copilot to submit review (polls every 10s)
  3. Fetches all inline comments and suggestions
  4. Applies suggested changes automatically
  5. Replies to each comment confirming the fix
  6. Resolves all review threads
  7. Commits and pushes fixes

  Use when user says "/gh-copilot-review", "copilot review", "wait for copilot", or when they want automated PR review iteration.

  BENEFITS:
  - Eliminates context switching between editor and GitHub
  - Applies Copilot suggestions with perfect accuracy
  - Resolves threads automatically after fixes
  - Provides clear commit message for review fixes
  - Saves 5-10 minutes per review cycle
tools: Bash, Read, Edit, Write, AskUserQuestion
---

# GitHub Copilot Review Automation

Fully automated workflow for receiving, addressing, and resolving GitHub Copilot code reviews without leaving Cursor. Perfect for rapid PR iteration.

## Steps

### 1. Wait for Copilot review

Run `watch-copilot-reviews` (from PATH) and wait for it to exit. It auto-detects the PR from the current branch. If @copilot hasn't been requested yet, it adds @copilot as a reviewer automatically. It polls until the review is submitted, then exits.

### 2. Fetch Copilot's inline comments

Store the repo name and PR number in variables:
```
REPO=$(gh repo view --json nameWithOwner --jq '.nameWithOwner')
PR_NUMBER=$(gh pr view --json number --jq '.number')
```

First, get the latest Copilot review ID so we only address current feedback (not old, already-resolved reviews):
```
LATEST_REVIEW_ID=$(gh api "repos/$REPO/pulls/$PR_NUMBER/reviews" --jq '[.[] | select(.user.login | test("copilot"; "i"))] | last | .id')
```

Then fetch only the comments from that review:
```
gh api "repos/$REPO/pulls/$PR_NUMBER/comments" --jq "[.[] | select(.user.login | test(\"copilot\"; \"i\")) | select(.pull_request_review_id == $LATEST_REVIEW_ID)]"
```

For each comment, extract: `id`, `node_id`, `path`, `line` (fall back to `original_line`), `body`, and any ` ```suggestion ` code blocks.

### 3. Address each comment

For each inline comment:
- Read the file at the indicated path and line
- If the comment contains a `suggestion` code block, apply the exact suggested change
- If it's general feedback (no suggestion block), make the appropriate code fix
- If the feedback is not applicable or incorrect, prepare a brief explanation why

### 4. Reply to each comment

For each comment you addressed, post a reply explaining what was done:
```
gh api "repos/$REPO/pulls/$PR_NUMBER/comments/$COMMENT_ID/replies" -f body="Done — applied the suggested change."
```

Tailor the reply to what you actually did (applied suggestion, made a fix, or explained why no change was needed).

### 5. Resolve each review thread

First, get the review thread IDs using GraphQL. Split `$REPO` into owner and name parts:
```
OWNER=$(echo "$REPO" | cut -d/ -f1)
REPO_NAME=$(echo "$REPO" | cut -d/ -f2)

gh api graphql -f query='
  query {
    repository(owner: "'"$OWNER"'", name: "'"$REPO_NAME"'") {
      pullRequest(number: '"$PR_NUMBER"') {
        reviewThreads(first: 100) {
          nodes {
            id
            isResolved
            comments(first: 1) {
              nodes {
                body
                author { login }
              }
            }
          }
        }
      }
    }
  }
'
```

Filter to only threads where `comments.nodes[0].author.login` matches Copilot (case-insensitive) and `isResolved` is `false`. Then resolve each one:
```
gh api graphql -f query='
  mutation {
    resolveReviewThread(input: { threadId: "'"$THREAD_NODE_ID"'" }) {
      thread { isResolved }
    }
  }
'
```

### 6. Commit and push

Stage all changes, commit with a message like `fix: address Copilot review feedback`, and push to origin.
Do not chain git commands with `&&` — run `git add`, `git commit`, and `git push` as separate Bash calls.

## Important notes

- The user is EXPLICITLY asking you to perform these git and GitHub API tasks.
- Do not ask for confirmation before applying suggestions or making fixes — just do it.
- If Copilot's review is APPROVED with no inline comments, just report that and stop.
- Only run through the review cycle once — do not re-request a review after pushing fixes.
