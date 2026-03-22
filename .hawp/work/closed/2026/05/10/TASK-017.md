# Work Intake — Plan Template

## Bug / Task: Fix librarian TypeScript validation errors

**Backlog ID:** TASK-017
**Type:** task
**Reported:** 2026-05-10
**Risk Level:** low

---

### Input (what was reported)

> i have some ts problems can you fix them please in the librarian project

---

### Context

The `librarian` package fails typecheck due to mismatched `ClosedTaskCheck` usage and missing helper symbols in backlog consistency validation. There is also a TypeScript configuration warning around `baseUrl` usage.

---

### Analysis

**Root cause (or most likely cause):**
Type definitions and reporter output diverged (`missing` no longer exists on `ClosedTaskCheck`), and helper functions referenced in backlog consistency were removed or not defined. `baseUrl` is set even though `paths` does not require it and it triggers a warning in current TS guidance.

**Directly verified:**

- `npm run -s typecheck` fails in `reporter.ts` and `backlog-consistency.ts`.
- `ClosedTaskCheck` defines `failing`/`warnings` fields, while reporter still reads `missing`.
- `findActiveFile` and `collectOrphanedActive` are referenced but undefined.
- `librarian/tsconfig.json` sets `baseUrl`.

**Inferred (not yet proven):**

- Current behavior should report both warnings and failures explicitly in output.

**Scope — what else is affected:**

- `librarian/scripts/validate-hawp-workflow/types.ts`
- `librarian/scripts/validate-hawp-workflow/reporter.ts`
- `librarian/scripts/validate-hawp-workflow/validations/backlog-consistency.ts`
- `librarian/tsconfig.json`

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `librarian/scripts/validate-hawp-workflow/reporter.ts`
- `librarian/scripts/validate-hawp-workflow/validations/backlog-consistency.ts`
- `librarian/tsconfig.json`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No active in-progress backlog item currently overlaps these `librarian` files.

---

### Options

#### Option A — Restore old `missing` model

Add `missing` back to `ClosedTaskCheck` and continue using legacy reporting.

#### Option B — Align reporter with current model

Keep `failing`/`warnings` structure and update reporter to print both.

---

### Recommended Fix

**Option chosen:** B
**Rationale:**
Preserves the newer model that differentiates legacy warnings from hard failures, avoids schema regression, and resolves strict typing cleanly.

**Files to change:**

- `librarian/scripts/validate-hawp-workflow/reporter.ts` — use `failing` and `warnings` in output.
- `librarian/scripts/validate-hawp-workflow/validations/backlog-consistency.ts` — implement missing helper functions.
- `librarian/tsconfig.json` — remove `baseUrl` while keeping `paths` support.

**What to verify after:**

- [x] `npm run -s typecheck` passes in `librarian`
- [x] report output still includes closed-task missing section details

---

### Implementation Notes

Keep the fix scoped to compile and config issues only.

---

## Outcome

Resolved all current `librarian` TypeScript compile errors and removed the `baseUrl` setting from TypeScript config to align with current guidance.

## Verification

- `npm run -s typecheck` in `librarian` completes with no errors.
- `npm run -s validate:workflow` runs successfully and closed-task reporting renders failing and legacy-warning sections correctly.

## Close Checklist

- [x] Backlog row moved to done in Recently Closed
- [x] Plan file moved from `active/` to dated `closed/` path
- [x] TypeScript errors addressed in scoped files

---

### Status

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
