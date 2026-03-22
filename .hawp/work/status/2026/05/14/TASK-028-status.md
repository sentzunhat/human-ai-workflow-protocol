# Status Report: TASK-028

**Completed:** May 14, 2026  
**Task:** Detection pipeline for backlog-upgrade workflow — rule evaluators, test fixtures, heuristic refinement

## Summary

TASK-028 expanded the `backlog-upgrade` CLI infrastructure with rule-driven detection and repair operations across historical backlog content.

- **Core deliverable:** Backlog consistency detection API with 10 rule classes (A1–A7, B1–B5)
- **Test coverage:** 15 unit tests spanning CLI parsing, script execution, rule evaluation, and edge cases
- **Heuristic refinement:** A4/A5/A7 operations now skip legacy closed files (pre-2026-05-10) to reduce false positives on historical content
- **Architecture alignment:** Reusable detection script API; thin CLI adapter with proper error boundaries

## Evidence

**Code Changes:**

- `librarian/scripts/backlog-upgrade/detection/rules/evaluate-rules.ts`: Added `isLegacyClosedFile()` and `extractDateFromPath()` to implement date-based tolerance for A4/A5/A7 operations
- `librarian/scripts/backlog-upgrade/__tests__/detection.test.ts`: Added legacy-file tolerance test (passing)
- Test suite: 15/15 tests passing (CLI parsing, script execution, detection rules, legacy tolerance)

**Validation:**

- TypeScript compilation: PASS
- Unit tests: 15/15 passing
- Workflow validation: PASS (2 checks, 0 issues, 2 warnings on legacy files accepted)

## Deliverables

1. **Detection Rule Evaluator** (`detection/rules/evaluate-rules.ts`)
   - 10 rule classes with configurable heuristics
   - Date-based legacy tolerance (cutoff: 2026-05-10)
   - Outputs structured fix/blocked operations

2. **Test Coverage** (15 tests)
   - CLI option parsing (3 tests)
   - Script-level execution (3 tests)
   - Rule evaluation (9 tests including legacy tolerance)

3. **Heuristic Refinement**
   - A4/A5: Skip operations on legacy closed files
   - A7: Explicit stale-reference patterns only
   - B4: Unproven verification marker detection
   - B5: Required work-directory validation

## Verification

- All 15 unit tests pass
- Heuristic tolerance applied: legacy files before 2026-05-10 skip A4/A5/A7 operations
- No regressions in existing functionality
- Clean typecheck (no errors)
- Workflow validation confirms backlog consistency

## Notes

- Heuristic cutoff date (2026-05-10) balances backward compatibility with forward-facing correctness enforcement
- Stale-reference patterns narrowed to high-confidence cases (`core/distribution/generated/`, `core/distribution/sources/`)
- Legacy tolerance test validates that files created before cutoff are not flagged for common formatting issues
- Ready for integration into applied repair workflow in follow-up tasks
