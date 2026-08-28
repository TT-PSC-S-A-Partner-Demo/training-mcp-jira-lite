// jira-lite MCP server in Go - the same two read-only tools as server.py,
// over the same fake-jira.json. No network, no auth.
//
// The point of the exercise is unchanged: this is one front-end over some logic.
// A CLI (see reference/jira-lite.go) exposes the SAME logic, a lighter interface.
//
// Run (from this go/ folder so ../fake-jira.json resolves):
//
//	go mod tidy
//	go run .
//
// Wire into Codex (stdio):
//
//	cd go && codex mcp add jira-lite -- go run .
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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

var issues []Issue

// load reads fake-jira.json. It looks in JIRA_DATA, then ../fake-jira.json
// (this file lives in go/, the data one level up), then ./fake-jira.json.
func load() error {
	candidates := []string{
		os.Getenv("JIRA_DATA"),
		filepath.Join("..", "fake-jira.json"),
		"fake-jira.json",
	}
	for _, path := range candidates {
		if path == "" {
			continue
		}
		b, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var d dataset
		if err := json.Unmarshal(b, &d); err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}
		issues = d.Issues
		return nil
	}
	return fmt.Errorf("fake-jira.json not found (set JIRA_DATA or run from the go/ folder)")
}

type ListInput struct {
	Assignee string `json:"assignee" jsonschema:"filter by assignee, e.g. 'me'; empty for all"`
}

type Row struct {
	Key     string `json:"key"`
	Status  string `json:"status"`
	Summary string `json:"summary"`
}

type ListOutput struct {
	Issues []Row `json:"issues"`
}

func listIssues(ctx context.Context, req *mcp.CallToolRequest, in ListInput) (*mcp.CallToolResult, ListOutput, error) {
	rows := []Row{}
	for _, i := range issues {
		if in.Assignee != "" && i.Assignee != in.Assignee {
			continue
		}
		rows = append(rows, Row{Key: i.Key, Status: i.Status, Summary: i.Summary})
	}
	return nil, ListOutput{Issues: rows}, nil
}

type GetInput struct {
	Key string `json:"key" jsonschema:"the issue key, e.g. 'PROJ-101'"`
}

func getIssue(ctx context.Context, req *mcp.CallToolRequest, in GetInput) (*mcp.CallToolResult, Issue, error) {
	for _, i := range issues {
		if i.Key == in.Key {
			return nil, i, nil
		}
	}
	return nil, Issue{}, fmt.Errorf("no issue %s", in.Key)
}

func main() {
	if err := load(); err != nil {
		log.Fatal(err)
	}
	server := mcp.NewServer(&mcp.Implementation{Name: "jira-lite", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_issues",
		Description: "List issues, optionally filtered by assignee (e.g. 'me'). Returns key, status, summary.",
	}, listIssues)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_issue",
		Description: "Get one issue in full by key, e.g. 'PROJ-101'.",
	}, getIssue)
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
