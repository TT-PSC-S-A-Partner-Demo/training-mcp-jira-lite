#!/usr/bin/env python3
"""
jira-lite CLI - REFERENCE ANSWER for the trainer.

This is what the participants should end up with: the SAME logic as server.py,
exposed as a plain CLI the agent calls on demand instead of a standing MCP socket.

Usage:
    python jira-lite.py list
    python jira-lite.py list --assignee me
    python jira-lite.py show PROJ-101

No SDK needed - just the standard library. That is part of the lesson: a CLI
front-end is lighter than an MCP server both to write and to run.
"""
import argparse
import json
import pathlib

# note: reference/ sits one level below the data, so go up to find fake-jira.json
DATA = json.loads(
    (pathlib.Path(__file__).parent.parent / "fake-jira.json").read_text(encoding="utf-8")
)


def list_issues(assignee=None):
    rows = DATA["issues"]
    if assignee:
        rows = [i for i in rows if i["assignee"] == assignee]
    return [{"key": i["key"], "status": i["status"], "summary": i["summary"]} for i in rows]


def get_issue(key):
    for i in DATA["issues"]:
        if i["key"] == key:
            return i
    return {"error": f"no issue {key}"}


def main():
    p = argparse.ArgumentParser(prog="jira-lite")
    sub = p.add_subparsers(dest="cmd", required=True)

    pl = sub.add_parser("list", help="list issues")
    pl.add_argument("--assignee")

    ps = sub.add_parser("show", help="show one issue by key")
    ps.add_argument("key")

    a = p.parse_args()
    if a.cmd == "list":
        print(json.dumps(list_issues(a.assignee), indent=2))
    elif a.cmd == "show":
        print(json.dumps(get_issue(a.key), indent=2))


if __name__ == "__main__":
    main()
