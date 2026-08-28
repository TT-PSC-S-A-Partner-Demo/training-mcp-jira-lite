// jira-lite CLI in Go - REFERENCE ANSWER for the trainer (Go track).
//
// Same logic as go/server.go, exposed as a plain CLI the agent calls on demand
// instead of a standing MCP socket. Standard library only - no MCP SDK.
//
// Run (single file, no module needed - stdlib only):
//
//	go run reference/jira-lite.go list
//	go run reference/jira-lite.go list me
//	go run reference/jira-lite.go show PROJ-101
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Issue struct {
	Key         string `json:"key"`
	Status      string `json:"status"`
	Assignee    string `json:"assignee"`
	Summary     string `json:"summary"`
	Updated     string `json:"updated"`
	Description string `json:"description"`
}

type dataset struct {
	Issues []Issue `json:"issues"`
}

func loadIssues() []Issue {
	for _, path := range []string{"fake-jira.json", "../fake-jira.json"} {
		b, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var d dataset
		if err := json.Unmarshal(b, &d); err != nil {
			fmt.Fprintln(os.Stderr, "parse error:", err)
			os.Exit(1)
		}
		return d.Issues
	}
	fmt.Fprintln(os.Stderr, "fake-jira.json not found (run from the demo-kit-mcp folder)")
	os.Exit(1)
	return nil
}

func printJSON(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: jira-lite list [assignee] | show <key>")
		os.Exit(2)
	}
	issues := loadIssues()

	switch os.Args[1] {
	case "list":
		assignee := ""
		if len(os.Args) > 2 {
			assignee = os.Args[2]
		}
		type row struct {
			Key     string `json:"key"`
			Status  string `json:"status"`
			Summary string `json:"summary"`
		}
		rows := []row{}
		for _, i := range issues {
			if assignee != "" && i.Assignee != assignee {
				continue
			}
			rows = append(rows, row{i.Key, i.Status, i.Summary})
		}
		printJSON(rows)

	case "show":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: jira-lite show <key>")
			os.Exit(2)
		}
		for _, i := range issues {
			if i.Key == os.Args[2] {
				printJSON(i)
				return
			}
		}
		printJSON(map[string]string{"error": "no issue " + os.Args[2]})

	default:
		fmt.Fprintln(os.Stderr, "unknown command:", os.Args[1])
		os.Exit(2)
	}
}
