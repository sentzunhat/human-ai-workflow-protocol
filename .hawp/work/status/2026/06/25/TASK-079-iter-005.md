# Workflow Loop — Iteration Handoff

---

## Header

**Backlog ID:** TASK-079 / `a7f3c2e1-9b4d-4a6f-8e2c-1d5f9a3b7c4e`
**Iteration:** 005 / **Budget:** 8
**Loop mode:** autonomous
**Date:** 2026-06-26
**Plan:** `.hawp/work/active/TASK-079.md`
**Executor:** agent
**Reviewer:** agent (review-only hat)
**Risk level:** low
**Auto-approve:** true

---

## Iteration Scope

Add **Autonomous vs Gated Mode** section and **standardized Continue prompt** with mandatory output path in `workflow-loop.md`.

---

## What Changed

**Files touched:**

- `.hawp/kit/usage/workflow-loop.md` — Autonomous vs Gated table; auto-advance rules; standardized Continue prompt block; mandatory handoff path; updated Transition and Copy-Paste Prompts; Quick Reference rows

---

## Verification

### Proven

- Continue prompt specifies same input/output each iteration
- Gated behavior preserved (checkpoint stops for human approve)

### Unproven

- None

---

## Review

**Outcome:** pass

**Findings:**

- Auto-advance conditions match constraints (risk + auto-approve + budget)
- Quick Reference answers "must user say loop again?"

**Blockers:** none

---

## Transition

**Decision:** auto-advance

**Approver:** agent-auto

**Notes:** Proceed to iter 6 — plan section snippet template.

---

## Suggested Next Step

Create `workflow-loop-plan-section.md` template for copy-paste into plans.

---

## Links

- Plan: `.hawp/work/active/TASK-079.md`
- Prior iteration: `.hawp/work/status/2026/06/25/TASK-079-iter-004.md`
- Evidence: none
