# TASK-014 Validation Output — Parked-Work Alignment — 2026-05-10

**npm run typecheck:** 0 errors
**npm run validate:workflow:** VALIDATION FAIL (expected — see notes)

---

## Full Validator Output

```
Validating: <user-home>/.../human-ai-workflow-protocol/.hawp/work

======================================================================
HAWP Workflow Validation Report
======================================================================

1. BACKLOG CONSISTENCY
----------------------------------------------------------------------

Active Work (1 items):
  Found: 1/1

Recently Closed (10 items):
  Found: 10/10

Orphaned Files (in active/ without backlog row):
  (none)

Blocked / Parked (1 items):
  Found: 1/1

Orphaned Files (in parked/ without backlog row):
  (none)

2. CLOSED TASK COMPLETENESS
----------------------------------------------------------------------

Checking 17 closed task file(s):
  Outcome: 0/17
  Verification: 0/17
  Close Checklist: 1/17

  Missing sections:
  2026-04-26-hawp-adr-template-review: missing Outcome, Verification, Close Checklist
  2026-04-30-bug-001-update-flow-migrate-old-layouts: missing Outcome, Verification, Close Checklist
  (... 14 more pre-HAWP-structure files ...)
  TASK-012: missing Outcome, Verification

3. EVIDENCE INTEGRITY
----------------------------------------------------------------------

  Found 0 evidence links
  ✓ 0 valid links

4. VERIFICATION CLARITY
----------------------------------------------------------------------

  Proven: 9/9

======================================================================
SUMMARY
======================================================================

✓ Checks passed:     3
✗ Issues found:      1
! Warnings:          0

Result: VALIDATION FAIL

======================================================================
```

---

## Notes on VALIDATION FAIL

The FAIL is expected and correct — two separate findings:

**Finding 1: Pre-HAWP-structure closed tasks (2026-04-26 to 2026-05-01)**
Files in `closed/2026/04/` and `closed/2026/05/01/` use legacy naming (date-prefixed). These predate the Outcome/Verification/Close Checklist template. They are expected to fail the completeness check. This is a real signal, not a false positive. A future compaction task should address these.

**Finding 2: TASK-012 shows "missing Outcome, Verification"**
TASK-012.md uses the template header `## Outcome (filled at close)` while the validator checks for exact match `## Outcome`. The validator is too strict — it should match section headers that start with `## Outcome`, `## Verification`, etc. This is a follow-up item.

---

## What Worked

- ✅ Parked check: TASK-013 found at `parked/TASK-013.md` (1/1)
- ✅ Orphaned parked: none detected
- ✅ TypeScript: 0 compile errors
- ✅ Parked section appears in report output
- ✅ `extractLinkPath()` correctly parses `[plan](parked/TASK-013.md)` → `parked/TASK-013.md`
