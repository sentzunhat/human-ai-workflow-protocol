# Task: Scripts structure guidelines and CLI boundary alignment

**Backlog ID:** TASK-077

## Context

Follow-up to TASK-076. The librarian scripts had good domain division (backlog-upgrade, backlog-validate, distribution, providers, validate-hawp-workflow, lib) but two domains mixed parsing, logic, and process control in one file, and the structural conventions were only documented implicitly inside `backlog-upgrade/CLI.md`.

## Scope

- Align CLI-shaped domains to the documented boundary pattern (`index.ts` boundary / `cli.ts` adapter / `script.ts` logic).
- Write explicit scripts guidelines covering domain division, single-use script folders, shared `lib/` rules, and how to add new scripts/CLI subcommands.
- Reconcile backlog and status records.

## Outcome

- `backlog-validate/` refactored from one 161-line monolith into `index.ts` + `cli.ts` + `script.ts`; identical observable behavior.
- `validate-hawp-workflow/` gained `script.ts` (`runWorkflowValidation`); `index.ts` is now a 4-line executable boundary, `process.exit` scattering removed.
- New `librarian/scripts/README.md` (guidelines: domain map table, boundary rules, lib constraints, conventions, add-a-script procedure); linked from `librarian/README.md`.
- Pipelines (`distribution/`, `providers/materialize/`) confirmed already conformant as simple single-use build/validate script folders — left unchanged by design.
- Backlog reconciled: Active Work and parked empty and accurate, Recently Closed capped at 10, STATUS.md current.

## Verification

- [x] `npm --prefix librarian run typecheck` exits 0. Evidence: ../evidence/2026/06/11/TASK-077-scripts-structure-verification.md
- [x] `npm --prefix librarian test` 38/38 pass. Evidence: ../evidence/2026/06/11/TASK-077-scripts-structure-verification.md
- [x] `validate:backlog` and `hawp backlog upgrade --dry-run --validate` PASS. Evidence: ../evidence/2026/06/11/TASK-077-scripts-structure-verification.md

## Close Checklist

- [x] Outcome recorded with explicit unchanged-by-design notes.
- [x] Evidence stored under `evidence/2026/06/11/`.
- [x] BACKLOG.md Recently Closed updated and capped at 10 rows.
- [x] STATUS.md refreshed.
