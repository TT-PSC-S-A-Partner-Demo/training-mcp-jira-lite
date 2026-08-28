# MCP -> CLI - build a tiny MCP, then let the agent make it a CLI

Self-contained. No real Jira, no auth, no network. Everything reads `fake-jira.json`.

**The idea:** an MCP server and a CLI are two front-ends over the same logic. MCP is a
standing socket whose tool schemas load into context every turn; a CLI is called on
demand. When the job is "let the agent read a few issues", the CLI is usually the
optimal front-end - same answer, fewer tokens, less standing access. This exercise makes
you feel that difference on the meter.

## Setup (once) - pick your language

**Python** (default):
```bash
pip install "mcp[cli]"          # the MCP Python SDK (FastMCP)
python server.py                # sanity check: should start and wait on stdio (Ctrl+C to stop)
```

**Go** (same two tools, same data - in `go/`):
```bash
cd go
go mod tidy                     # fetch github.com/modelcontextprotocol/go-sdk
go run .                        # sanity check (Ctrl+C to stop); reads ../fake-jira.json
```
Reference CLI answer for the Go track: `reference/jira-lite.go` (stdlib only, run `go run reference/jira-lite.go list`).

## Step 1 - run the MCP and use it (this is the BEFORE)

```bash
codex mcp add jira-lite -- python /ABS/PATH/TO/server.py     # Python
# or, Go (run from the go/ folder so ../fake-jira.json resolves):
#   cd go && codex mcp add jira-lite -- go run .
codex mcp list
```
In Codex:
```
/mcp        # confirm jira-lite + its two tools are reachable
/status     # note the token count - this is BEFORE
```
Ask, through the MCP:
```
List the issues assigned to me, then show PROJ-101 in full.
One line per issue: key, status, summary.
```

## Step 2 - let the agent convert it to a CLI

```
Look at server.py. It exposes list_issues and get_issue over MCP, reading fake-jira.json.
Write a plain CLI, jira-lite.py, with the SAME logic and the same data - no MCP, no SDK,
standard library only. Commands:
  python jira-lite.py list [--assignee me]
  python jira-lite.py show PROJ-101
Output JSON. Do not change fake-jira.json.
```
Reference answer for the trainer: `reference/jira-lite.py`.

## Step 3 - use the CLI instead of the MCP (this is the AFTER)

Drop the MCP so it stops loading its schemas every turn, then start a fresh chat:
```bash
codex mcp remove jira-lite
```
```
/new
/status     # AFTER
```
Ask for the SAME thing, but now the agent just runs the CLI:
```
Run jira-lite.py to list the issues assigned to me, then show PROJ-101 in full.
One line per issue: key, status, summary.
```

## Step 4 - compare

Same answer both times? Compare `/status` BEFORE vs AFTER. Which front-end was lighter,
and why. (The MCP pays for its tool schemas every turn even when unused; the CLI pays
only when it is called.)

> When is the MCP still the right call? When you need it mid-reasoning, repeatedly, or
> the tool set is large and dynamic. "Optimal" is not "always CLI" - it is matching the
> front-end to how often and how deeply the agent needs the tool.

## Bonus - the untrusted ticket

`PROJ-105` has a description telling the assistant to delete the repo. Read it with
`show PROJ-105`. Nothing happens - a read-only front-end cannot act on it. That is the
whole point of scope: the defence is the narrow tool, not the agent reading carefully.
