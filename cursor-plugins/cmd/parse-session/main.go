package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cxm6467/caboose-ai/cursor-plugins/pkg/session"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: parse-session <--current|path-to-transcript.jsonl>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "For Cursor IDE:")
		fmt.Fprintln(os.Stderr, "  parse-session --current")
		fmt.Fprintln(os.Stderr, "  parse-session ~/.cursor/projects/my-project/agent-transcripts/UUID/UUID.jsonl")
		os.Exit(1)
	}

	arg := os.Args[1]

	var summary *session.SessionSummary
	var err error

	// Determine input type
	switch arg {
	case "--current":
		// Find most recent Cursor transcript
		transcriptPath, err := session.FindRecentCursorSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to find recent transcript: %v\n", err)
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "Make sure you have Cursor transcripts at:")
			fmt.Fprintln(os.Stderr, "  ~/.cursor/projects/*/agent-transcripts/")
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "Analyzing: %s\n", transcriptPath)
		summary, err = session.ParseCursorSession(transcriptPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse transcript: %v\n", err)
			os.Exit(1)
		}

	default:
		// Assume it's a file path
		if _, err := os.Stat(arg); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "File not found: %s\n", arg)
			os.Exit(1)
		}

		summary, err = session.ParseCursorSession(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse transcript: %v\n", err)
			os.Exit(1)
		}
	}

	// Output as JSON
	output, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal output: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
