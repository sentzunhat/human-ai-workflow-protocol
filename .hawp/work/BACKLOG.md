# Backlog

Single source of truth for all open bugs, tasks, and improvement work.
Each row links to its plan file when one exists.

---

## Status Key

| Status        | Meaning                             |
| ------------- | ----------------------------------- |
| `inbox`       | Received, not yet analyzed          |
| `analyzing`   | Under investigation                 |
| `plan-ready`  | Plan written, awaiting review       |
| `approved`    | Plan approved, ready to implement   |
| `in-progress` | Being implemented                   |
| `parked`      | Deferred without closing            |
| `done`        | Implemented and verified            |
| `blocked`     | Blocked — reason noted in plan file |
| `wont-fix`    | Decided not to fix — reason noted   |

---

## Active Work

| ID       | Type | Title                                       | Status  | Plan        |
| -------- | ---- | ------------------------------------------- | ------- | ----------- |
| TASK-008 | task | Compact BACKLOG.md and archive closed work. | `inbox` | _(pending)_ |

---

## Done

| ID       | Type        | Title                                                                 | Closed     | Plan                                                                                       |
| -------- | ----------- | --------------------------------------------------------------------- | ---------- | ------------------------------------------------------------------------------------------ |
| BUG-001  | bug         | update migrates old `.hawp/usage` into `.hawp/work`                   | 2026-04-30 | [plan](closed/2026/04/30/2026-04-30-bug-001-update-flow-migrate-old-layouts.md)            |
| BUG-002  | bug         | update flow cleans up legacy work folders after migration             | 2026-05-01 | [plan](closed/2026/05/01/2026-05-01-bug-002-update-flow-clean-legacy-work-folders.md)      |
| BUG-003  | bug         | reconcile closed plans during update/install flow and docs sync       | 2026-05-01 | [plan](closed/2026/05/01/2026-05-01-bug-003-closed-plan-reconciliation-and-doc-sync.md)    |
| BUG-004  | bug         | reconcile done active plans using backlog IDs and closed dates        | 2026-05-01 | [plan](closed/2026/05/01/2026-05-01-bug-004-reconcile-done-active-plans-by-id-and-date.md) |
| BUG-005  | bug         | align install/update flow with validated migration and reconciliation | 2026-05-01 | [plan](closed/2026/05/01/2026-05-01-bug-005-install-update-alignment.md)                   |
| TASK-001 | improvement | ADR: Parallel Work Guardrails                                         | 2026-04-30 | [plan](closed/2026/04/30/2026-04-30-task-001-parallel-work-guardrails-adr.md)              |
| TASK-002 | improvement | ADR: Intake Template Refinement                                       | 2026-04-30 | [plan](closed/2026/04/30/2026-04-30-task-002-intake-template-adr.md)                       |
| TASK-003 | improvement | ADR: Restructure `.hawp` into `kit`/`work`                            | 2026-04-30 | [plan](closed/2026/04/30/2026-04-30-task-003-kit-work-restructure-adr.md)                  |
| TASK-004 | improvement | Root README UX refactor (concise install/update)                      | 2026-04-30 | [plan](closed/2026/04/30/2026-04-30-task-004-readme-ux-refactor.md)                        |
| TASK-005 | improvement | add core backlog alignment guidance for compact active management     | 2026-05-02 | [plan](closed/2026/05/02/TASK-005.md)                                                      |
| TASK-006 | improvement | dogfood repo work state under repo-root .hawp/work                    | 2026-05-02 | [plan](closed/2026/05/02/TASK-006.md)                                                      |
| TASK-007 | improvement | expose backlog alignment as agent-facing prompt and instruction       | 2026-05-02 | [plan](closed/2026/05/02/TASK-007.md)                                                      |
| ADR-004  | decision    | Date-based work folder structure + safe migration                     | 2026-05-01 | [decision](decisions/2026/05/01/ADR-004-date-based-work-folder-structure.md)               |

---

## Notes

- Check this file before starting any new item.
- Active plan files go in `active/`. Close by moving to `closed/YYYY/MM/DD/`.
- Deferred items can live in `parked/` without being closed.
- ADRs and decisions go in `decisions/YYYY/MM/DD/`.
- Each item gets one plan file — no two agents on the same ID.
- Work started outside this loop should still get a row added for visibility.
