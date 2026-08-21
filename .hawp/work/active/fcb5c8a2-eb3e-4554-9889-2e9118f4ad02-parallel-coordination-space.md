---
work-item: fcb5c8a2-eb3e-4554-9889-2e9118f4ad02
type: feature
title: "Optional HAWP parallel coordination execution space"
status: plan-ready
created: 2026-08-15
updated: 2026-08-15
parent: b6c4e8a2
depends-on: c1d2e401
---

# Feature: Parallel Coordination Execution Space

## Mission

Document and implement an optional, disposable `.hawp/.space/<branch>/`
convention for isolated Git worktrees, while keeping `.hawp/work/` as durable
shared coordination state.

## Boundaries

- `.hawp/.space/` is gitignored, machine-local, disposable, and never proof
  or durable project memory.
- Preserve the existing HAWP Shape and Markdown-first workflow; do not add a
  runtime, daemon, lock service, database, or permanent manager branch.
- Reuse existing work ownership and overlap guardrails. Provider sources and
  generated overlays remain separate from repository-local coordination.

## Done When

- Canonical guidance defines branch-to-worktree mapping and coordinator role.
- Provider overlays receive aligned instructions without changing the core
  Shape.
- Setup, collision, cleanup, and no-space fallback are documented and tested
  with Git worktree fixtures.
