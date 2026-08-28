# mcp-jira-lite - a tiny MCP over a fake Jira

Workshop kit for Session 2 (Exercise 3): build a small MCP server, then let the agent
turn the same logic into a CLI, and compare the token cost. Self-contained - no real
Jira, no auth, no network. Everything reads `fake-jira.json`.

## Contents

```
server.py        # MCP server (FastMCP): list_issues, get_issue - read-only
reference/jira-lite.py   # the CLI answer for the trainer
fake-jira.json   # 5 fake tickets
PROMPTS.md       # full exercise steps
```

## Setup

```bash
pip install "mcp[cli]<2"     # the MCP Python SDK
python server.py           # sanity check (Ctrl+C to stop)
```

Then follow `PROMPTS.md`: connect the MCP, use it, have the agent rewrite it as a CLI,
and compare `/status` before vs after.
