# TASK-067 Evidence — Closed Work Audit Matrix

Date: 2026-05-17

## Repo-Root Proof (Redacted Prefix)

```bash
$ pwd
<repo-root-abs>

$ git rev-parse --show-toplevel
<repo-root-abs>

$ git rev-parse --show-prefix

$ git status --short
 D .hawp/kit/guidance/da-schema-planning.md
 D .hawp/kit/guidance/shared-standards-review-rubric.md
 M .hawp/kit/standards/README.md
 M .hawp/kit/templates/intake-plan.md
 M .hawp/kit/usage/intake-workflow.md
 M .hawp/work/BACKLOG.md
 M .hawp/work/active/TASK-066.md
?? .hawp/kit/standards/database/
?? .hawp/kit/standards/patterns/
?? .hawp/kit/standards/service-design/
?? .hawp/work/active/TASK-067.md
?? .hawp/work/active/TASK-068.md
?? .hawp/work/evidence/2026/05/17/TASK-067-closed-audit-matrix.md
?? .hawp/work/parked/TASK-066.md
```

## Inventory Summary

- Total closed markdown plans scanned: 73
- Closed files outside date folders: 0
- Allowed non-plan file in closed tree: `.hawp/work/closed/README.md`
- Files requiring alignment action: 20

## Action Counts

```json
{
  "total": 73,
  "actions": 20,
  "missingBacklogId": 11,
  "missingOutcome": 12,
  "missingVerification": 9,
  "missingCloseChecklist": 16,
  "dateMismatch": 0,
  "duplicateIds": 0
}
```

## Action Matrix (Exact Paths)

| Path                                                                                          | Backlog ID | Issues                                                                             |
| --------------------------------------------------------------------------------------------- | ---------- | ---------------------------------------------------------------------------------- |
| .hawp/work/closed/2026/04/26/2026-04-26-hawp-adr-template-review.md                           | (none)     | missing-backlog-id, missing-outcome, missing-verification, missing-close-checklist |
| .hawp/work/closed/2026/04/30/2026-04-30-bug-001-update-flow-migrate-old-layouts.md            | (none)     | missing-backlog-id, missing-outcome, missing-verification, missing-close-checklist |
| .hawp/work/closed/2026/04/30/2026-04-30-task-001-parallel-work-guardrails-adr.md              | (none)     | missing-backlog-id, missing-outcome, missing-close-checklist                       |
| .hawp/work/closed/2026/04/30/2026-04-30-task-002-intake-template-adr.md                       | (none)     | missing-backlog-id, missing-outcome, missing-close-checklist                       |
| .hawp/work/closed/2026/04/30/2026-04-30-task-003-kit-work-restructure-adr.md                  | (none)     | missing-backlog-id, missing-outcome, missing-close-checklist                       |
| .hawp/work/closed/2026/04/30/2026-04-30-task-004-readme-ux-refactor.md                        | (none)     | missing-backlog-id, missing-outcome, missing-verification, missing-close-checklist |
| .hawp/work/closed/2026/05/01/2026-05-01-bug-002-update-flow-clean-legacy-work-folders.md      | BUG-002    | missing-outcome, missing-verification, missing-close-checklist                     |
| .hawp/work/closed/2026/05/01/2026-05-01-bug-003-closed-plan-reconciliation-and-doc-sync.md    | BUG-003    | missing-outcome, missing-verification, missing-close-checklist                     |
| .hawp/work/closed/2026/05/01/2026-05-01-bug-004-reconcile-done-active-plans-by-id-and-date.md | BUG-004    | missing-outcome, missing-verification, missing-close-checklist                     |
| .hawp/work/closed/2026/05/01/2026-05-01-bug-005-install-update-alignment.md                   | BUG-005    | missing-outcome, missing-verification, missing-close-checklist                     |
| .hawp/work/closed/2026/05/02/TASK-005.md                                                      | TASK-005   | missing-close-checklist                                                            |
| .hawp/work/closed/2026/05/02/TASK-006.md                                                      | TASK-006   | missing-outcome, missing-verification, missing-close-checklist                     |
| .hawp/work/closed/2026/05/02/TASK-007.md                                                      | TASK-007   | missing-close-checklist                                                            |
| .hawp/work/closed/2026/05/02/TASK-008.md                                                      | TASK-008   | missing-close-checklist                                                            |
| .hawp/work/closed/2026/05/02/TASK-009.md                                                      | TASK-009   | missing-close-checklist                                                            |
| .hawp/work/closed/2026/05/03/TASK-010.md                                                      | (none)     | missing-backlog-id, missing-outcome, missing-verification, missing-close-checklist |
| .hawp/work/closed/2026/05/12/TASK-030-files.md                                                | (none)     | missing-backlog-id                                                                 |
| .hawp/work/closed/2026/05/13/TASK-030-files.md                                                | (none)     | missing-backlog-id                                                                 |
| .hawp/work/closed/2026/05/13/TASK-038.md                                                      | (none)     | missing-backlog-id                                                                 |
| .hawp/work/closed/2026/05/14/TASK-044.md                                                      | (none)     | missing-backlog-id                                                                 |

## Classification for Next Task Handoff

- Automation-eligible section scaffolding candidates: 16
- Header-only backlog ID insertion candidates: 4
- Folder/date relocation candidates: 0
- Duplicate backlog ID remediation candidates: 0

Notes:

- Matrix uses heading detection across markdown levels `#` through `######` to avoid false missing-section flags.
- Existing workflow validator legacy tolerance remains unchanged in this audit.
