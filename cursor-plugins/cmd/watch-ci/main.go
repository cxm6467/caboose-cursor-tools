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
)

type CheckStatus struct {
	Name        string `json:"name"`
	State       string `json:"state"`
	CompletedAt string `json:"completedAt"`
	DetailsURL  string `json:"detailsUrl"`
}

type RunInfo struct {
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	Jobs       []Job  `json:"jobs"`
}

type Job struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
}

func main() {
	interval := 10
	if len(os.Args) > 1 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			printHelp()
			os.Exit(0)
		}

		var err error
		interval, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: interval must be a positive integer, got '%s'\n", os.Args[1])
			os.Exit(1)
		}

		if interval < 5 {
			fmt.Printf("Warning: Interval too low (%d), using minimum of 5 seconds\n", interval)
			interval = 5
		}
	}

	// Check if we're in a git repo
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	if err := cmd.Run(); err != nil {
		fmt.Println("Not in a git repository")
		os.Exit(1)
	}

	// Check for gh CLI
	if _, err := exec.LookPath("gh"); err != nil {
		fmt.Println("GitHub CLI (gh) not found. Install with: brew install gh")
		os.Exit(1)
	}

	// Check for jq
	if _, err := exec.LookPath("jq"); err != nil {
		fmt.Println("jq not found. Install with: brew install jq")
		os.Exit(1)
	}

	// Get current branch
	cmd = exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		fmt.Println("Not on a branch (detached HEAD)")
		os.Exit(1)
	}
	branch := strings.TrimSpace(string(output))

	fmt.Printf("Watching CI for branch: %s%s%s\n", colorBlue, branch, colorReset)
	fmt.Printf("Refresh interval: %ds (press Ctrl+C to stop)\n\n", interval)

	// Initial check
	status := showCIStatus(branch)
	if status == 2 {
		os.Exit(0)
	}
	if status == 1 {
		os.Exit(1)
	}

	// Watch loop
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		status := showCIStatus(branch)
		if status == 2 {
			os.Exit(0)
		}
		if status == 1 {
			os.Exit(1)
		}
	}
}

// showCIStatus checks CI status and returns: 0 = running, 1 = failed, 2 = passed
func showCIStatus(branch string) int {
	clearScreen()
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	fmt.Printf("CI Status for branch: %s%s%s\n", colorBlue, branch, colorReset)
	fmt.Printf("Last updated: %s\n", timestamp)
	fmt.Println("Press Ctrl+C to stop watching\n")

	// Try PR checks first
	cmd := exec.Command("gh", "pr", "checks", "--json", "name,state,completedAt,detailsUrl")
	output, err := cmd.Output()

	if err == nil && len(output) > 0 && string(output) != "[]" {
		var checks []CheckStatus
		if json.Unmarshal(output, &checks) == nil && len(checks) > 0 {
			return handlePRChecks(checks, branch)
		}
	}

	// No PR found, check workflow runs
	return handleWorkflowRuns(branch)
}

func handlePRChecks(checks []CheckStatus, branch string) int {
	fmt.Println("CI Check Results:")
	fmt.Println("====================")

	hasFail := 0
	hasPending := 0

	for _, check := range checks {
		switch check.State {
		case "SUCCESS":
			fmt.Printf("%sPASS%s %s\n", colorGreen, colorReset, check.Name)
		case "FAILURE", "ERROR":
			fmt.Printf("%sFAIL%s %s\n", colorRed, colorReset, check.Name)
			if check.DetailsURL != "" {
				fmt.Printf("     %s%s%s\n", colorBlue, check.DetailsURL, colorReset)
			}
			hasFail++
		case "PENDING", "QUEUED", "IN_PROGRESS", "WAITING", "REQUESTED", "STARTUP_FAILURE":
			fmt.Printf("%sWAIT%s %s - running...\n", colorYellow, colorReset, check.Name)
			hasPending++
		case "NEUTRAL", "SKIPPED", "STALE":
			fmt.Printf("%sSKIP%s %s - %s\n", colorYellow, colorReset, check.Name, check.State)
		default:
			fmt.Printf("%s????%s %s - %s\n", colorYellow, colorReset, check.Name, check.State)
		}
	}

	fmt.Println()

	if hasFail > 0 && hasPending == 0 {
		fmt.Printf("%sSome checks failed%s\n\n", colorRed, colorReset)
		showFailureLogs(branch)
		return 1
	} else if hasPending > 0 {
		fmt.Printf("%sChecks still running...%s\n", colorYellow, colorReset)
		return 0
	} else {
		fmt.Printf("%sAll checks passed!%s\n", colorGreen, colorReset)
		return 2
	}
}

func handleWorkflowRuns(branch string) int {
	fmt.Printf("CI Workflow Status (latest run on %s):\n", branch)
	fmt.Println("==============================================")

	// Get latest run ID
	cmd := exec.Command("gh", "run", "list", "--branch", branch, "--limit", "1", "--json", "databaseId", "--jq", ".[0].databaseId")
	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		fmt.Printf("No CI runs found for branch %s\n", branch)
		fmt.Println("Make sure you've pushed your branch")
		return 0
	}

	runID := strings.TrimSpace(string(output))

	// Get run info
	cmd = exec.Command("gh", "run", "view", runID, "--json", "status,conclusion,jobs")
	output, err = cmd.Output()
	if err != nil {
		fmt.Println("Could not fetch run info")
		return 0
	}

	var runInfo RunInfo
	if err := json.Unmarshal(output, &runInfo); err != nil {
		fmt.Println("Could not parse run info")
		return 0
	}

	// Display jobs
	for _, job := range runInfo.Jobs {
		if job.Status == "completed" {
			switch job.Conclusion {
			case "success":
				fmt.Printf("%sPASS%s %s\n", colorGreen, colorReset, job.Name)
			case "failure":
				fmt.Printf("%sFAIL%s %s\n", colorRed, colorReset, job.Name)
			case "cancelled":
				fmt.Printf("%sSKIP%s %s - cancelled\n", colorYellow, colorReset, job.Name)
			default:
				fmt.Printf("%s????%s %s - %s\n", colorYellow, colorReset, job.Name, job.Conclusion)
			}
		} else {
			fmt.Printf("%sWAIT%s %s - %s...\n", colorYellow, colorReset, job.Name, job.Status)
		}
	}

	fmt.Println()

	if runInfo.Status == "completed" {
		if runInfo.Conclusion == "success" {
			fmt.Printf("%sAll jobs passed!%s\n", colorGreen, colorReset)
			return 2
		} else {
			fmt.Printf("%sWorkflow %s%s\n\n", colorRed, runInfo.Conclusion, colorReset)
			showFailureLogs(runID)
			return 1
		}
	} else {
		fmt.Printf("%sWorkflow still running...%s\n", colorYellow, colorReset)
		return 0
	}
}

func showFailureLogs(branchOrRunID string) {
	fmt.Println("Failure logs:")
	fmt.Println("===============")

	var cmd *exec.Cmd
	if strings.Contains(branchOrRunID, "/") {
		// It's a branch, get latest run ID first
		runCmd := exec.Command("gh", "run", "list", "--branch", branchOrRunID, "--limit", "1", "--json", "databaseId", "--jq", ".[0].databaseId")
		output, err := runCmd.Output()
		if err != nil || len(output) == 0 {
			return
		}
		runID := strings.TrimSpace(string(output))
		cmd = exec.Command("gh", "run", "view", runID, "--log-failed")
	} else {
		cmd = exec.Command("gh", "run", "view", branchOrRunID, "--log-failed")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func printHelp() {
	fmt.Println("Usage: watch-ci [interval_seconds]")
	fmt.Println()
	fmt.Println("Watch GitHub Actions CI status for the current branch.")
	fmt.Println("Polls until all checks pass or fail, then exits.")
	fmt.Println()
	fmt.Println("  interval_seconds  Polling interval in seconds (default: 10, minimum: 5)")
}
