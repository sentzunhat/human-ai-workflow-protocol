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

## Outcome

Created `.hawp/kit/usage/hawp-first-workflow.md` (62 lines) covering: why hawp-first (token cost rationale), session-start pattern, during-work pattern, fall-through criteria, re-index triggers, parallel agent worktree cleanup, token budget and session continuity, and a quick reference table. Added item 6 to `start-here.md` Workflow Guides list. `hawp kit validate` — 3 checks passed, 0 issues.

## Verification

- [x] `.hawp/kit/usage/hawp-first-workflow.md` exists and passes kit validate. Evidence: see Outcome section above.
- [x] `start-here.md` references the new doc. Evidence: see Outcome section above.
- [x] No Go source changes. Evidence: see Outcome section above.
- [x] `hawp kit validate` — 0 issues. Evidence: see Outcome section above.

## Close Checklist

- [x] Doc written and validated
- [x] start-here.md updated
- [x] BACKLOG updated; plan moved to closed

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
