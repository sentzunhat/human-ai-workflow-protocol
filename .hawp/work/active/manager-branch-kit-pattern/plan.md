# manager-branch-kit-pattern — Document manager-branch / worktree operating pattern in kit

**Type:** improvement  
**Status:** plan-ready  
**Opened:** 2026-08-25  
**Target:** v0.0.11

## Input

Downstream install of hawp 0.0.9 in a monorepo with a product integration
branch (`origin/QA`). The team uses a long-lived manager branch
(`chore/hawp-manager`) to hold HAWP kit, backlog, and coordinator dispatch
without polluting product PRs. Product work happens in git worktrees from
`origin/QA`. The manager branch commits but never merges into QA.

This is a useful, repeatable operating pattern for teams with existing
git integration workflows. It is not a HAWP protocol field — it is an
optional project-topology choice.

Evidence source: downstream install evidence 2026-08-25; confirmed working.

## Goal

Add an optional kit doc (e.g. `.hawp/kit/usage/manager-branch.md` or as a
section in `start-here.md`) describing this pattern clearly enough for a new
team to adopt it. Keep it short and decision-useful.

## Constraints

- HAWP must not become a runtime or require this pattern — it is optional.
- Do not invent new HAWP fields, folders, or lifecycle events.
- Keep it one concise doc, not an ADR (this is an operating note, not a design decision for HAWP core).
- Reference the upstream evidence but redact machine-local paths.

## Plan

### Step 1 — Draft `.hawp/kit/usage/manager-branch.md`

Contents:
- When to use: monorepos with an integration branch, where HAWP kit/backlog noise in product PRs is undesirable
- Manager branch setup (one long-lived branch, commits OK, never PR)
- Product worktrees: cut from integration branch, not manager HEAD
- What the manager branch owns: `.hawp/kit/`, `.hawp/work/**`, HAWP init/update
- What it must not touch: product code in worktrees
- One-paragraph example: `chore/hawp-manager` + `origin/QA` + worktrees pattern
- Caveats: merge conflicts on `.hawp/` if product code ever touches it; keep product code and kit on separate paths

### Step 2 — Link from `start-here.md`

Add a one-line reference in `.hawp/kit/start-here.md` under a "Project topologies" or "Advanced setup" section:

> For monorepos with a product integration branch, see
> [`usage/manager-branch.md`](usage/manager-branch.md) for an optional
> manager-branch / worktree pattern.

### Step 3 — Validate links

Run `npm run kit:validate` and `npm run check:markdown-links`.

## Verification

- `npm run kit:validate` passes
- `npm run check:markdown-links` passes
- The doc is readable without prior HAWP knowledge
