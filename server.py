#!/usr/bin/env python3
"""
jira-lite MCP server - a tiny, self-contained MCP over a fake Jira (fake-jira.json).

Two tools, both READ-ONLY. No network, no auth, no real Jira.
The point of the exercise: this is one front-end over some logic. In jira-lite.py
the SAME logic is exposed as a CLI. Same answers, a lighter interface.

Run it directly to sanity-check:   python server.py
Wire it into Codex (stdio):        codex mcp add jira-lite -- python /abs/path/server.py

Needs the MCP SDK:                 pip install "mcp[cli]"
"""
import json
import pathlib

from mcp.server.fastmcp import FastMCP

DATA = json.loads((pathlib.Path(__file__).parent / "fake-jira.json").read_text(encoding="utf-8"))

mcp = FastMCP("jira-lite")


@mcp.tool()
def list_issues(assignee: str | None = None) -> list[dict]:
    """List issues, optionally filtered by assignee (e.g. 'me').
    Returns a short row per issue: key, status, summary."""
    rows = DATA["issues"]
    if assignee:
        rows = [i for i in rows if i["assignee"] == assignee]
    return [{"key": i["key"], "status": i["status"], "summary": i["summary"]} for i in rows]


@mcp.tool()
def get_issue(key: str) -> dict:
    """Get one issue in full by key, e.g. 'PROJ-101'."""
    for i in DATA["issues"]:
        if i["key"] == key:
            return i
    return {"error": f"no issue {key}"}


if __name__ == "__main__":
    mcp.run()  # stdio transport
