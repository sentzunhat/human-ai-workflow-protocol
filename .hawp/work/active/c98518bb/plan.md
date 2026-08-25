# c98518bb — Hawp-first session workflow: MCP search as default context strategy

**Type:** improvement  
**Status:** in-progress  
**Branch:** feature/v008-hawp-workflow-docs → feature/v0.0.8  
**Updated:** 2026-08-24

## Input

Agent sessions consume tokens reading kit files directly. MCP `hawp_search` returns
structured, relevance-ranked chunks — far cheaper per lookup. We need a canonical
guide telling agents (and humans configuring agents) how to use hawp-first at
session start and during work.

## Scope

- Write `.hawp/kit/usage/hawp-first-workflow.md` — the session start pattern:
  search before reading, pipe MCP results to stay within budget, re-index when content changes.
- Update `start-here.md` to reference the new doc (one line).
- BACKLOG: move c98518bb to in-progress → done on merge.

## Out of scope

- Token counting / measurement (that is 4c88f451).
- Changes to Go source code.

## Acceptance

- New kit doc exists and passes `hawp kit validate`.
- `start-here.md` links to it.
- Commit on `feature/v008-hawp-workflow-docs`, squash-merge → `feature/v0.0.8`.
