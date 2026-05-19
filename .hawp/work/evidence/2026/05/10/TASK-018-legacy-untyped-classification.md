# TASK-018 Evidence — Legacy Untyped Classification

Date: 2026-05-10

## Commands

- `npm run typecheck` in `librarian` → pass
- `npm run validate:workflow` in `librarian` → pass

## Key Validator Output (Closed Task Completeness)

```text
Checking 20 plan file(s):
  Outcome: 5/20
  Verification: 7/20
  Close Checklist: 5/20

  [WARN] Legacy untyped closed files (before 2026-05-10 — tolerated, visible):
    2026-04-26-hawp-adr-template-review: legacy file without TASK-/BUG-style ID (2026-04-26)

  [WARN] Legacy files missing sections (before 2026-05-10 — tolerated):
    ...
```

## Summary

- Legacy untyped file `2026-04-26-hawp-adr-template-review.md` is now visible as WARN (not silently skipped).
- Supporting-file skips are now explicit-pattern-based and reported under INFO.
- Overall validation remains `VALIDATION PASS` with warnings only.
