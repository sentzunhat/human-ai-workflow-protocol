## Status Report

### Intent

Review and recommend a lightweight markdown-first pattern to prevent parallel human/agent work collisions (ownership confusion, overlapping edits, duplicated execution, and premature implementation before review approval).

### Current State

The repository has a lean intake loop, one-backlog source of truth, a plan-per-item rule, and explicit status transitions. Parallel guardrails are implemented through a pattern doc, coordination fields in templates, and one workflow reminder before implementation.

### What Was Inspected

- `core/.hawp/usage/INTAKE_WORKFLOW.md`
- `core/.hawp/usage/BACKLOG.md`
- `core/.hawp/templates/intake-plan.md`
- `core/.hawp/usage/INIT.md`
- `core/.hawp/usage/STATUS_REPORT.md`
- `.github/instructions/hawp-intake.instructions.md`

### What Changed

- Added `core/.hawp/patterns/parallel-work-guardrails.md`.
- Added a `Work Coordination` block to `core/.hawp/templates/intake-plan.md`.
- Confirmed `core/.hawp/templates/bug-plan.md` includes the `Work Coordination` block.
- Added one concise overlap/approval reminder in `core/.hawp/usage/INTAKE_WORKFLOW.md`.

### What Was Directly Verified

- The status vocabulary already includes: inbox, analyzing, plan-ready, approved, in-progress, done, blocked, wont-fix.
- The intake workflow already requires plan files in `usage/status/` before implementation work.
- The workflow already states: backlog must be checked before starting; one plan file per item; avoid two agents on same ID.
- `BACKLOG.md` currently tracks active and done items and links to plan files.
- The intake and bug plan templates include explicit ownership, overlap, and implementation-readiness coordination fields.

### What Remains Unproven

- How often file-overlap collisions currently happen in real multi-agent usage.
- Whether teams will consistently maintain a backlog-level lock table under pressure without process drift.

### Constraints

- Markdown-first and file-first only.
- No runtime lock service, daemon, database, package install, or required GitHub Issues dependency.
- Keep HAWP lean and adoption-friendly for founder-speed workflows.

### Help Wanted

- Confirm whether this recommendation should be adopted as a pattern-only change first, then optionally template-assisted after one trial cycle.

### Suggested Next Step

Run the guardrails on the next 2-3 parallel items and confirm whether coordination friction stays low while preventing overlap collisions.

### Optional Attached Artifact

# ADR: Parallel Work Guardrails (Markdown-First)

## Title

Parallel Work Guardrails for HAWP (Human + Multiple AI Agents)

## Context

HAWP in this repo is intentionally lean and file-first. Current workflow already supports intake, planning, review gates, implementation, and verification through markdown artifacts.

## Problem

Parallel contributors can still step on each other when ownership is implicit, file overlap is not surfaced early, and one agent starts implementing while another item is still awaiting review approval.

## Current state found in repo

- Existing strengths:
  - Backlog is designated source of truth for active and completed work.
  - One plan file per item is already required.
  - Review gate exists with low/medium/high risk behavior.
  - Done state expects verification.
- Gap:
  - Plan template lacks explicit micro-fields for owner, overlap, and implementation readiness under parallel activity.

## Options

### Option A - Add `patterns/parallel-work-guardrails.md` only

- Pros: minimal change, no template churn, easiest to trial.
- Cons: optional guidance may be skipped during fast execution; inconsistent adoption risk.

### Option B - Add `patterns/parallel-work-guardrails.md` + intake template coordination fields

- Pros: keeps process markdown-only and small while making critical coordination checks explicit at planning time; aligns with existing review gate behavior.
- Cons: adds small extra planning friction (a few lines per non-trivial item).

### Option C - Add full work-lock table to `BACKLOG.md`

- Pros: centralized visibility in one file.
- Cons: highest maintenance burden; creates pseudo-PM overhead; likely to drift and conflict with item-level truth in plan files.

## Recommended option

Option B.

Rationale: It is the smallest reliable guardrail that improves behavior where collisions happen (inside active plan work), without turning backlog into a lock manager or adding runtime systems.

## Proposed markdown fields

Add this section to intake plan template:

```md
### Work Coordination

**Owner:** human / agent / unassigned
**Implementation status:** not-started / in-progress / blocked / done
**Overlapping files:**

- path/to/file.ts

**Parallel work risk:** low / medium / high
**Can implement now:** yes / no / only after approval

**Coordination note:**
Explain whether another active item touches the same files.
```

Usage notes:

- If overlap exists and risk is medium/high, set `Can implement now` to `only after approval`.
- If overlap is exact same file and another item is `in-progress`, default to `no` unless user explicitly approves parallel edits.

## Example usage

```md
### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- src/auth/session.ts
- src/auth/middleware.ts

**Parallel work risk:** high
**Can implement now:** only after approval

**Coordination note:**
TASK-014 is already in-progress and edits src/auth/session.ts.
This item should wait to avoid conflicting auth changes.
```

## Files proposed to add/change

- Add: `core/.hawp/patterns/parallel-work-guardrails.md`
- Change: `core/.hawp/templates/intake-plan.md`

## Files explicitly not to touch

- `core/.hawp/usage/BACKLOG.md` structure beyond normal item row updates
- Runtime/tooling files (no lock service, no scripts, no package config)
- Protocol schema core files (`core/.hawp/SPEC.md`, `core/.hawp/types/shape.ts`) unless separately requested

## Risks

- Mild process overhead may discourage use on tiny tasks.
- False confidence if overlap lists are not updated when scope expands.
- Inconsistent usage unless reinforced by intake workflow language.

## Acceptance criteria

- Guardrails remain pure markdown and optional workflow discipline.
- No runtime lock mechanism is introduced.
- Plan authors can quickly answer: who owns this, what files overlap, and whether implementation can begin now.
- Parallel collision risk is surfaced before edits begin.

## Verification checklist

- [x] Pattern doc exists with simple rules and examples.
- [x] Intake template includes Work Coordination section.
- [x] At least one sample plan demonstrates overlap detection and hold behavior.
- [x] Workflow text still preserves lean HAWP scope (no PM system expansion).

## Implementation Status

Implemented on 2026-04-30 after verification.

Changed files:

- core/.hawp/patterns/parallel-work-guardrails.md
- core/.hawp/templates/intake-plan.md
- core/.hawp/usage/INTAKE_WORKFLOW.md

Confirmed unchanged by this implementation:

- core/.hawp/types/shape.ts
- core/.hawp/SPEC.md
- Runtime/tooling/package-script surfaces

## Implementation plan (do not implement yet)

1. Draft `core/.hawp/patterns/parallel-work-guardrails.md` with 5-8 short rules.
2. Add Work Coordination block to `core/.hawp/templates/intake-plan.md`.
3. Add one sentence in intake workflow reminding agents to check active work + overlap before implementation.
4. Pilot on 2-3 real items, then evaluate friction and collision reduction.
5. Keep, trim, or revert fields based on pilot signal.
