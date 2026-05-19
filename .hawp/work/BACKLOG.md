# Backlog

Active coordination index for open work. Closed history is archived under `.hawp/work/closed/`.

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

| UUID | Legacy ID | Type | Title | Status | Owner | Plan File | Updated |
| ---- | --------- | ---- | ----- | ------ | ----- | --------- | ------- |

_(empty)_

---

## Blocked / Parked

| ID  | Type | Title | Reason | Detail | Updated |
| --- | ---- | ----- | ------ | ------ | ------- |

_(empty)_

---

## Recently Closed

Limited to the last 10 items.

| ID       | Type             | Title                                                                     | Closed     | Detail                                |
| -------- | ---------------- | ------------------------------------------------------------------------- | ---------- | ------------------------------------- |
| TASK-068 | improvement      | Automate closed-record section scaffolding and date-folder reconciliation | 2026-05-17 | [plan](closed/2026/05/17/TASK-068.md) |
| TASK-067 | task             | Audit and align all closed work items to current close-record standards   | 2026-05-17 | [plan](closed/2026/05/17/TASK-067.md) |
| TASK-066 | task             | Normalize remaining legacy closed records flagged by validator warnings   | 2026-05-17 | [plan](closed/2026/05/17/TASK-066.md) |
| TASK-065 | task             | Verify GitHub-hosted distribution auto-sync run (TASK-036 evidence gap)   | 2026-05-15 | [plan](closed/2026/05/15/TASK-065.md) |
| TASK-064 | task             | Compact BACKLOG.md and archive closed work                                | 2026-05-15 | [plan](closed/2026/05/15/TASK-064.md) |
| TASK-054 | task             | Review shared standards and absorb public-safe guidance                   | 2026-05-15 | [plan](closed/2026/05/15/TASK-054.md) |
| TASK-053 | task             | Audit and consolidate opinionated standards folder layout                 | 2026-05-15 | [plan](closed/2026/05/15/TASK-053.md) |
| TASK-060 | governance       | Handle project-specific standards boundaries (Tekit/Mictlan/Zacatl)       | 2026-05-15 | [plan](closed/2026/05/15/TASK-060.md) |
| TASK-059 | governance       | Review and retain private standards lane boundaries                       | 2026-05-15 | [plan](closed/2026/05/15/TASK-059.md) |
| TASK-063 | standards-update | Extract generalized standards from Zacatl adaptation candidates           | 2026-05-15 | [plan](closed/2026/05/15/TASK-063.md) |

---

## Archive

- Closed work: `closed/`
- Status reports: `status/`
- Evidence: `evidence/`
- Decisions: `decisions/`

## Notes

- Check this file before starting any new item.
- Active plan files go in `active/`. Close by moving to `closed/YYYY/MM/DD/`.
- Deferred items can live in `parked/` without being closed.
- ADRs and decisions go in `decisions/YYYY/MM/DD/`.
- Each item gets one plan file — no two agents on the same ID.
- Work started outside this loop should still get a row added for visibility.
- Keep `Recently Closed` capped; archive history lives in `closed/`.

## Future Improvements

- **UUID-based work item IDs** (TASK-013): Move from sequential type-prefixed IDs (TASK-001, BUG-001) to UUIDs for parallel-safe work item creation. Depends on validation script (TASK-012) to surface all ID references in workflow.
