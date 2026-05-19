# TASK-049 Status Report

## Summary

TASK-049 refactored backlog-upgrade into a script-first architecture. Core execution moved out of CLI control flow and is now reusable programmatically, while CLI remains a thin adapter with documented command contract.

## Delivered

- Added reusable script API: `librarian/scripts/backlog-upgrade/script.ts` (`runBacklogUpgradeScript`).
- Refactored CLI adapter: `librarian/scripts/backlog-upgrade/cli.ts` now parses/validates args and delegates execution.
- Kept process boundary in executable entrypoint: `librarian/scripts/backlog-upgrade/index.ts`.
- Added command contract doc: `librarian/scripts/backlog-upgrade/CLI.md`.
- Added script-focused tests: `librarian/scripts/backlog-upgrade/__tests__/script.test.ts`.

## Verification Evidence

- `npm --prefix librarian run typecheck` passed.
- `cd librarian && npx --yes tsx --test scripts/backlog-upgrade/__tests__/cli.test.ts scripts/backlog-upgrade/__tests__/detection.test.ts scripts/backlog-upgrade/__tests__/script.test.ts` passed (9/9).
- `npm --prefix librarian run validate:workflow` returned `VALIDATION PASS` (warnings only).

## Residual Risk

- `--apply` mode remains intentionally unimplemented in this slice and returns usage error.
- Broader detection heuristics/fixtures remain tracked in TASK-028.
