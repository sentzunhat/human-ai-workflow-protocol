# TASK-019 Tekit Validator Compatibility

- Target: `tekit/.hawp`
- Command: `npm run validate:workflow -- --hawp-root tekit/.hawp`
- Exit code: `1`

## Raw Output

```text

> @hawp/librarian@0.0.0 validate:workflow
> tsx scripts/validate-hawp-workflow/index.ts --hawp-root tekit/.hawp

Validating: tekit/.hawp/work

======================================================================
HAWP Workflow Validation Report
======================================================================

1. BACKLOG CONSISTENCY
----------------------------------------------------------------------

Active Work (0 items):
  Found: 0/0

Recently Closed (10 items):
  Found: 10/10

Orphaned Files (in active/ without backlog row):
  (none)

Blocked / Parked (0 items):
  Found: 0/0

Orphaned Files (in parked/ without backlog row):
  (none)

2. CLOSED TASK COMPLETENESS
----------------------------------------------------------------------

Checking 136 plan file(s)  (10 supporting file(s) skipped):
  Outcome: 0/136
  Verification: 1/136
  Close Checklist: 0/136

  [FAIL] Missing sections (2026-05-10 or later — must fix):
    BUG-156: ✗ missing Outcome, Verification, Close Checklist  (2026-05-10)

  [WARN] Legacy untyped closed files (before 2026-05-10 — tolerated, visible):
    2026-04-26-hawp-adr-template-review: legacy file without TASK-/BUG-style ID  (2026-04-26)

  [WARN] Legacy files missing sections (before 2026-05-10 — tolerated):
    BUG-001: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-002: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-003: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-004: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-005: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-007: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-009: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-010: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-013: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-014: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-015: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-016: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-017: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-018: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-019: missing Outcome, Verification, Close Checklist  (2026-04-29)
    BUG-020: missing Outcome, Verification, Close Checklist  (2026-04-30)
    BUG-021: missing Outcome, Verification, Close Checklist  (2026-04-30)
    BUG-022: missing Outcome, Verification, Close Checklist  (2026-04-30)
    BUG-025: missing Outcome, Verification, Close Checklist  (2026-04-30)
    BUG-026: missing Outcome, Verification, Close Checklist  (2026-04-30)
    BUG-027: missing Outcome, Verification, Close Checklist  (2026-04-30)
    BUG-028: missing Outcome, Verification, Close Checklist  (2026-04-30)
    BUG-029: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-030: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-031: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-032: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-033: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-034: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-035: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-036: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-037: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-038: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-039: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-041: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-042: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-043: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-044: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-045: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-046: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-047: missing Outcome, Verification, Close Checklist  (2026-05-01)
    BUG-048: missing Outcome, Verification, Close Checklist  (2026-05-02)
    BUG-049: missing Outcome, Verification, Close Checklist  (2026-05-02)
    BUG-050: missing Outcome, Verification, Close Checklist  (2026-05-02)
    BUG-052: missing Outcome, Verification, Close Checklist  (2026-05-02)
    BUG-053: missing Outcome, Verification, Close Checklist  (2026-05-02)
    BUG-054: missing Outcome, Verification, Close Checklist  (2026-05-02)
    BUG-055: missing Outcome, Verification, Close Checklist  (2026-05-02)
    BUG-056: missing Outcome, Verification, Close Checklist  (2026-05-02)
    BUG-057: missing Outcome, Verification, Close Checklist  (2026-05-02)
    BUG-061: missing Outcome, Verification, Close Checklist  (2026-05-02)
    BUG-058: missing Outcome, Verification, Close Checklist  (2026-05-02)
    BUG-059: missing Outcome, Verification, Close Checklist  (2026-05-02)
    BUG-060: missing Outcome, Verification, Close Checklist  (2026-05-02)
    BUG-062: missing Outcome, Verification, Close Checklist  (2026-05-02)
    BUG-063: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-064: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-065: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-066: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-067: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-068: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-069: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-070: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-071: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-072: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-072: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-073: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-074: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-075: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-076: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-077: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-078: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-079: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-080: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-081: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-082: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-083: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-084: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-085: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-086: missing Outcome, Verification, Close Checklist  (2026-05-03)
    BUG-088: missing Outcome, Verification, Close Checklist  (2026-05-04)
    BUG-089: missing Outcome, Verification, Close Checklist  (2026-05-04)
    BUG-090: missing Outcome, Verification, Close Checklist  (2026-05-04)
    BUG-091: missing Outcome, Verification, Close Checklist  (2026-05-05)
    BUG-092: missing Outcome, Verification, Close Checklist  (2026-05-05)
    BUG-093: missing Outcome, Verification, Close Checklist  (2026-05-05)
    BUG-094: missing Outcome, Verification, Close Checklist  (2026-05-05)
    BUG-096: missing Outcome, Verification, Close Checklist  (2026-05-05)
    BUG-097: missing Outcome, Verification, Close Checklist  (2026-05-05)
    BUG-098: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-099: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-100: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-101: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-102: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-103: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-104: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-105: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-106: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-107: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-108: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-109: missing Outcome, Close Checklist  (2026-05-06)
    BUG-110: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-113: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-115: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-116: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-117: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-118: missing Outcome, Verification, Close Checklist  (2026-05-06)
    BUG-120: missing Outcome, Verification, Close Checklist  (2026-05-07)
    BUG-122: missing Outcome, Verification, Close Checklist  (2026-05-07)
    BUG-125: missing Outcome, Verification, Close Checklist  (2026-05-07)
    BUG-128: missing Outcome, Verification, Close Checklist  (2026-05-07)
    BUG-130: missing Outcome, Verification, Close Checklist  (2026-05-07)
    BUG-131: missing Outcome, Verification, Close Checklist  (2026-05-07)
    BUG-132: missing Outcome, Verification, Close Checklist  (2026-05-08)
    BUG-133: missing Outcome, Verification, Close Checklist  (2026-05-08)
    BUG-134: missing Outcome, Verification, Close Checklist  (2026-05-08)
    BUG-135: missing Outcome, Verification, Close Checklist  (2026-05-08)
    BUG-136: missing Outcome, Verification, Close Checklist  (2026-05-08)
    BUG-137: missing Outcome, Verification, Close Checklist  (2026-05-08)
    BUG-138: missing Outcome, Verification, Close Checklist  (2026-05-08)
    BUG-140: missing Outcome, Verification, Close Checklist  (2026-05-08)
    BUG-141: missing Outcome, Verification, Close Checklist  (2026-05-08)
    BUG-142: missing Outcome, Verification, Close Checklist  (2026-05-08)
    BUG-143: missing Outcome, Verification, Close Checklist  (2026-05-08)
    BUG-144: missing Outcome, Verification, Close Checklist  (2026-05-09)
    BUG-145: missing Outcome, Verification, Close Checklist  (2026-05-09)
    BUG-146: missing Outcome, Verification, Close Checklist  (2026-05-09)
    BUG-147: missing Outcome, Verification, Close Checklist  (2026-05-09)
    BUG-148: missing Outcome, Verification, Close Checklist  (2026-05-09)
    BUG-149: missing Outcome, Verification, Close Checklist  (2026-05-09)
    BUG-150: missing Outcome, Verification, Close Checklist  (2026-05-09)
    BUG-151: missing Outcome, Verification, Close Checklist  (2026-05-09)
    BUG-152: missing Outcome, Verification, Close Checklist  (2026-05-09)
    BUG-153: missing Outcome, Verification, Close Checklist  (2026-05-09)
    BUG-154: missing Outcome, Verification, Close Checklist  (2026-05-09)
    BUG-155: missing Outcome, Verification, Close Checklist  (2026-05-09)

  [INFO] Supporting files skipped by pattern:
    2026-05-02-backlog-archive-BUG-001-to-BUG-042: matches archive supporting-file pattern  (2026-05-02)
    2026-05-02-backlog-archive-BUG-043-to-BUG-044: matches archive supporting-file pattern  (2026-05-02)
    2026-05-02-backlog-archive-BUG-045-to-BUG-046: matches archive supporting-file pattern  (2026-05-02)
    BUG-067-checkpoint: matches supporting suffix pattern (-checkpoint)  (2026-05-03)
    BUG-069-checkpoint: matches supporting suffix pattern (-checkpoint)  (2026-05-03)
    BUG-071-checkpoint: matches supporting suffix pattern (-checkpoint)  (2026-05-03)
    BUG-072-close-gate-checkpoint: matches supporting suffix pattern (-close-gate-checkpoint)  (2026-05-03)
    BUG-072-final-checkpoint: matches supporting suffix pattern (-final-checkpoint)  (2026-05-03)
    BUG-091-status-report: matches supporting suffix pattern (-status-report)  (2026-05-05)
    BUG-139-mictlan-deleted-user-archive-save-missing-id: matches archive supporting-file pattern  (2026-05-08)

3. EVIDENCE INTEGRITY
----------------------------------------------------------------------

  Found 0 evidence links
  ✓ 0 valid links

4. VERIFICATION CLARITY
----------------------------------------------------------------------

  Proven: 0/0

======================================================================
SUMMARY
======================================================================

✓ Checks passed:     3
✗ Issues found:      1
! Warnings:          0

Result: VALIDATION FAIL

======================================================================
```
