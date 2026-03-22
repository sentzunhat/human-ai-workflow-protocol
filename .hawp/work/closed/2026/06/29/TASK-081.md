## Task: Add Optional Loop-Lane Note to Intake Plan Template

**Backlog ID:** TASK-081
**UUID:** `c9f5a4e1-2d6b-4b83-9c10-7d1a0f8d6f21`
**Type:** improvement
**Reported:** 2026-06-29
**Risk Level:** low
**Status:** plan-ready

---

### Goal

Add an optional `Loop lane` note to the intake-plan template so parallel Workflow Loop assignments can be recorded consistently without changing core workflow behavior.

### Why this compounds

`TASK-080` needs a second eligible lane to validate parallel coordination. A lightweight template note gives future lanes a standard place to record ownership, reducing ambiguity without introducing runtime orchestration.

### Scope

- Update `.hawp/kit/templates/intake-plan.md`
- Keep the change optional and non-breaking
- Do not add CLI or runtime logic

### Work Coordination

**Owner:** unassigned
**Implementation status:** plan-ready
**Overlapping files:**

- `.hawp/kit/templates/intake-plan.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
This file set is disjoint from `TASK-080`, which is currently operating in `.hawp/work/*` coordination files.

---

### Recommended Fix

Add one optional `Loop lane` note under the coordination section in `.hawp/kit/templates/intake-plan.md`, with wording that makes it clear the field is metadata only and can be left blank.

### Workflow Loop

**Loop status:** active
**Loop mode:** autonomous
**Iteration budget:** 3
**Current iteration:** 0
**Executor:** agent
**Reviewer:** agent (separate session)
**Approver:** agent per risk gate
**Auto-approve:** true

### Iteration Log

| Iter | Date | Outcome | Handoff |
| ---- | ---- | ------- | ------- |
| 001 | 2026-06-29 | started | _pending_ |
| 002 | 2026-06-29 | template note added | _pending_ |

---

### Recommended Fix

Add one short optional bullet or note to the template's coordination section describing the loop lane field, if present, as a non-blocking metadata field for parallel trials.

### Verification

- Template still renders as a generic plan template
- New note remains optional
- No existing required fields change

### Close Checklist

### Outcome

`TASK-081` added the optional `Loop lane` note to the intake-plan template so future parallel trials have a standard metadata field without changing workflow semantics.

### Verification

- [x] Template updated
- [x] Loop-lane note stays optional
- [x] Backlog and plan are aligned

**Direct evidence:** `.hawp/kit/templates/intake-plan.md` now includes a `Loop lane` note under coordination.
