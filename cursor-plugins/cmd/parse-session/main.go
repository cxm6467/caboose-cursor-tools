package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cxm6467/caboose-ai/cursor-plugins/pkg/session"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: parse-session <session-id|--current|path-to-session>")
		os.Exit(1)
	}

	arg := os.Args[1]

	parser, err := session.NewParser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create parser: %v\n", err)
		os.Exit(1)
	}

	var sessionID string

	// Determine session ID
	switch arg {
	case "--current":
		sessionID, err = parser.FindRecentSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to find recent session: %v\n", err)
			os.Exit(1)
		}
	default:
		// Assume it's a session ID
		sessionID = arg
	}

	// Parse the session
	summary, err := parser.ParseSession(sessionID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse session: %v\n", err)
		os.Exit(1)
	}

	// Output as JSON
	output, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal output: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
