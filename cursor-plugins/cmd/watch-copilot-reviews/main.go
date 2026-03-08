package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[0;31m"
	colorGreen  = "\033[0;32m"
	colorYellow = "\033[1;33m"
	colorBlue   = "\033[0;34m"
	colorCyan   = "\033[0;36m"
	colorDim    = "\033[2m"
	colorBold   = "\033[1m"
)

type RequestedReviewer struct {
	Login string `json:"login"`
}

type RequestedReviewers struct {
	Users []RequestedReviewer `json:"users"`
}

type Review struct {
	ID          int    `json:"id"`
	State       string `json:"state"`
	SubmittedAt string `json:"submitted_at"`
	Body        string `json:"body"`
	User        struct {
		Login string `json:"login"`
	} `json:"user"`
}

type Comment struct {
	Path               string `json:"path"`
	Line               int    `json:"line"`
	OriginalLine       int    `json:"original_line"`
	Body               string `json:"body"`
	PullRequestReviewID int    `json:"pull_request_review_id"`
	User               struct {
		Login string `json:"login"`
	} `json:"user"`
}

type PRInfo struct {
	Number     int    `json:"number"`
	HeadRefOid string `json:"headRefOid"`
}

type CheckRun struct {
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	Name       string `json:"name"`
	App        struct {
		Slug string `json:"slug"`
	} `json:"app"`
}

type CheckRuns struct {
	CheckRuns []CheckRun `json:"check_runs"`
}

func main() {
	prNumber := ""
	interval := 10

	// Parse arguments
	for i, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" {
			printHelp()
			os.Exit(0)
		}

		// Check if it's a number
		if num, err := strconv.Atoi(arg); err == nil {
			if i == 0 {
				prNumber = arg
			} else if i == 1 {
				interval = num
			}
		}
	}

	if interval < 5 {
		fmt.Printf("⚠️  Interval too low (%d), using minimum of 5 seconds\n", interval)
		interval = 5
	}

	// Check prerequisites
	if err := exec.Command("git", "rev-parse", "--git-dir").Run(); err != nil {
		fmt.Println("❌ Not in a git repository")
		os.Exit(1)
	}

	if _, err := exec.LookPath("gh"); err != nil {
		fmt.Println("❌ GitHub CLI (gh) not found. Install with: brew install gh")
		os.Exit(1)
	}

	if _, err := exec.LookPath("jq"); err != nil {
		fmt.Println("❌ jq not found. Install with: brew install jq")
		os.Exit(1)
	}

	// Get PR number if not provided
	if prNumber == "" {
		cmd := exec.Command("git", "branch", "--show-current")
		output, err := cmd.Output()
		if err != nil || len(output) == 0 {
			fmt.Println("❌ Not on a branch (detached HEAD). Pass a PR number instead.")
			os.Exit(1)
		}

		branch := strings.TrimSpace(string(output))
		cmd = exec.Command("gh", "pr", "view", branch, "--json", "number", "--jq", ".number")
		output, err = cmd.Output()
		if err != nil || len(output) == 0 {
			fmt.Printf("❌ No PR found for branch: %s\n", branch)
			fmt.Println("💡 Pass a PR number: watch-copilot-reviews 30")
			os.Exit(1)
		}
		prNumber = strings.TrimSpace(string(output))
	}

	// Get repo info
	cmd := exec.Command("gh", "repo", "view", "--json", "nameWithOwner", "--jq", ".nameWithOwner")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("❌ Failed to get repository info")
		os.Exit(1)
	}
	repo := strings.TrimSpace(string(output))

	fmt.Printf("🤖 Watching Copilot reviews for %s%s#%s%s\n", colorBlue, repo, prNumber, colorReset)
	fmt.Printf("⏱️  Refresh interval: %ds (press Ctrl+C to stop)\n\n", interval)

	// Initial check
	if showReviewStatus(repo, prNumber) == 2 {
		os.Exit(0)
	}

	// Watch loop
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if showReviewStatus(repo, prNumber) == 2 {
			os.Exit(0)
		}
	}
}

// showReviewStatus checks review status and returns: 0 = pending, 1 = error, 2 = complete
func showReviewStatus(repo, prNumber string) int {
	clearScreen()
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	fmt.Printf("🤖 Copilot Review Status for %s%s#%s%s\n", colorBlue, repo, prNumber, colorReset)
	fmt.Printf("⏱️  Last updated: %s\n", timestamp)
	fmt.Println("📝 Press Ctrl+C to stop watching\n")

	// Check for pending review requests
	cmd := exec.Command("gh", "api", fmt.Sprintf("repos/%s/pulls/%s/requested_reviewers", repo, prNumber))
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		var reviewers RequestedReviewers
		if json.Unmarshal(output, &reviewers) == nil {
			for _, user := range reviewers.Users {
				if strings.Contains(strings.ToLower(user.Login), "copilot") {
					fmt.Printf("%s⏳ Copilot review pending...%s\n\n", colorYellow, colorReset)
					fmt.Println("  Copilot is still analyzing the PR. This typically takes 1-3 minutes.")

					// Check copilot agent status
					cmd = exec.Command("gh", "pr", "view", prNumber, "--json", "headRefOid", "--jq", ".headRefOid")
					if headSha, err := cmd.Output(); err == nil && len(headSha) > 0 {
						sha := strings.TrimSpace(string(headSha))
						cmd = exec.Command("gh", "api", fmt.Sprintf("repos/%s/commits/%s/check-runs", repo, sha))
						if checkOutput, err := cmd.Output(); err == nil {
							var checks CheckRuns
							if json.Unmarshal(checkOutput, &checks) == nil {
								for _, check := range checks.CheckRuns {
									if check.App.Slug == "copilot-pull-request-review" || check.Name == "Copilot" {
										conclusion := check.Conclusion
										if conclusion == "" {
											conclusion = "pending"
										}
										fmt.Printf("\n  Agent status: %s%s\t%s%s\n", colorCyan, check.Status, conclusion, colorReset)
									}
								}
							}
						}
					}
					return 0
				}
			}
		}
	}

	// Check for submitted reviews
	cmd = exec.Command("gh", "api", fmt.Sprintf("repos/%s/pulls/%s/reviews", repo, prNumber))
	output, err = cmd.Output()
	if err != nil {
		return 1
	}

	var reviews []Review
	if err := json.Unmarshal(output, &reviews); err != nil {
		return 1
	}

	// Filter to copilot reviews
	var copilotReviews []Review
	for _, review := range reviews {
		if strings.Contains(strings.ToLower(review.User.Login), "copilot") {
			copilotReviews = append(copilotReviews, review)
		}
	}

	if len(copilotReviews) > 0 {
		// Show latest review
		latest := copilotReviews[len(copilotReviews)-1]

		fmt.Printf("%s✅ Copilot review received!%s\n\n", colorGreen, colorReset)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		switch latest.State {
		case "APPROVED":
			fmt.Printf("  %s%s✅ APPROVED%s  %s(%s)%s\n", colorGreen, colorBold, colorReset, colorDim, latest.SubmittedAt, colorReset)
		case "CHANGES_REQUESTED":
			fmt.Printf("  %s%s🔄 CHANGES REQUESTED%s  %s(%s)%s\n", colorRed, colorBold, colorReset, colorDim, latest.SubmittedAt, colorReset)
		case "COMMENTED":
			fmt.Printf("  %s%s💬 COMMENTED%s  %s(%s)%s\n", colorYellow, colorBold, colorReset, colorDim, latest.SubmittedAt, colorReset)
		default:
			fmt.Printf("  %s%s❓ %s%s  %s(%s)%s\n", colorYellow, colorBold, latest.State, colorReset, colorDim, latest.SubmittedAt, colorReset)
		}

		if latest.Body != "" {
			fmt.Println()
			lines := strings.Split(latest.Body, "\n")
			for _, line := range lines {
				fmt.Printf("  %s\n", line)
			}
		}
		fmt.Println()

		if len(copilotReviews) > 1 {
			fmt.Printf("  %s(Showing latest of %d reviews)%s\n\n", colorDim, len(copilotReviews), colorReset)
		}

		// Show inline comments from latest review
		cmd = exec.Command("gh", "api", fmt.Sprintf("repos/%s/pulls/%s/comments", repo, prNumber))
		output, err = cmd.Output()
		if err == nil {
			var comments []Comment
			if json.Unmarshal(output, &comments) == nil {
				var latestComments []Comment
				for _, comment := range comments {
					if strings.Contains(strings.ToLower(comment.User.Login), "copilot") &&
						comment.PullRequestReviewID == latest.ID {
						latestComments = append(latestComments, comment)
					}
				}

				if len(latestComments) > 0 {
					fmt.Printf("%s%s📝 Inline Comments (%d):%s\n\n", colorCyan, colorBold, len(latestComments), colorReset)

					for _, comment := range latestComments {
						line := comment.Line
						if line == 0 {
							line = comment.OriginalLine
						}
						fmt.Printf("  %s%s%s:%s%d%s\n", colorBlue, comment.Path, colorReset, colorYellow, line, colorReset)

						// Strip suggestion blocks for cleaner output
						body := comment.Body
						inSuggestion := false
						var filteredLines []string
						for _, line := range strings.Split(body, "\n") {
							if strings.HasPrefix(line, "```suggestion") {
								inSuggestion = true
								continue
							}
							if inSuggestion && strings.HasPrefix(line, "```") {
								inSuggestion = false
								continue
							}
							if !inSuggestion {
								filteredLines = append(filteredLines, line)
							}
						}
						body = strings.Join(filteredLines, "\n")

						for _, line := range strings.Split(body, "\n") {
							fmt.Printf("    %s\n", line)
						}
						fmt.Println()
					}
				}
			}
		}

		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()
		fmt.Printf("🔗 %shttps://github.com/%s/pull/%s%s\n", colorBlue, repo, prNumber, colorReset)
		fmt.Println()
		fmt.Println("🎉 Exiting since Copilot review is complete")
		return 2
	}

	// No pending and no reviews - request Copilot
	fmt.Printf("%s🔄 No Copilot review request found — requesting one...%s\n", colorYellow, colorReset)
	cmd = exec.Command("gh", "pr", "edit", prNumber, "--add-reviewer", "@copilot")
	if err := cmd.Run(); err == nil {
		fmt.Printf("  %s✅ Added @copilot as reviewer%s\n", colorGreen, colorReset)
	} else {
		fmt.Printf("  %s❌ Failed to add @copilot as reviewer%s\n", colorRed, colorReset)
		fmt.Printf("  💡 Try manually: %sgh pr edit %s --add-reviewer @copilot%s\n", colorCyan, prNumber, colorReset)
	}
	return 0
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func printHelp() {
	fmt.Println("Usage: watch-copilot-reviews [pr_number] [interval_seconds]")
	fmt.Println()
	fmt.Println("Watch for GitHub Copilot code reviews on a pull request.")
	fmt.Println("Both arguments are numeric. A single argument is always the PR number.")
	fmt.Println("If no PR number is given, uses the PR for the current branch.")
	fmt.Println()
	fmt.Println("  pr_number         PR number to watch (default: current branch's PR)")
	fmt.Println("  interval_seconds  Polling interval in seconds (default: 10, minimum: 5)")
}
