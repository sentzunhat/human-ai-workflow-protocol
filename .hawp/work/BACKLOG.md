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
| `1c2db85b` | — | improvement | research librarian CLI, binary packaging, and local model runtime feasibility | in-progress | agent | [plan](active/1c2db85b-d57b-47a1-8c4b-b28dc55f2d5c.md) | 2026-07-08 |
| `8dedc4e2` | — | improvement | design installable librarian CLI surface and bin contract | in-progress | agent | [plan](active/8dedc4e2-69c5-42ca-aaad-93b62d7fb899.md) | 2026-07-08 |
| `1a3b32a4` | — | improvement | spike transformers.js Node WASM and shared model cache for librarian | analyzing | unassigned | [plan](active/1a3b32a4-ab37-4c86-ade7-71e2eb42b440.md) | 2026-07-08 |
| `54e68af7` | — | improvement | spike single-binary packaging and shared model asset strategy for librarian | analyzing | unassigned | [plan](active/54e68af7-4622-4383-8482-cc4d4e1e21ee.md) | 2026-07-08 |

---

## Blocked / Parked

| ID  | Type | Title | Reason | Detail | Updated |
| --- | ---- | ----- | ------ | ------ | ------- |
| `bee15107` | improvement | defer CLI participant adapters for Codex, Claude, and GitHub | Not needed yet; provider packs are enough for current use | [plan](parked/bee15107-e4d1-41b3-b67e-04c87328fef3.md) | 2026-07-06 |

---

## Recently Closed

Limited to the last 10 items.

| ID       | Type             | Title                                                                     | Closed     | Detail                                |
| -------- | ---------------- | ------------------------------------------------------------------------- | ---------- | ------------------------------------- |
| `3030ebff` | improvement | final audit of core HAWP kit and provider install instructions | 2026-07-08 | [plan](closed/2026/07/08/3030ebff-55b0-4d34-ad76-0909c0fa552c.md) |
| `a114f1ed` | improvement | add Codex provider pack before broader librarian CLI work | 2026-07-06 | [plan](closed/2026/07/06/a114f1ed-7565-4c7d-b0fd-b22a387ca99f.md) |
| `ddeb9eb3` | improvement | work normalize usability and verification clarity follow-up | 2026-07-03 | [plan](closed/2026/07/03/ddeb9eb3-6dba-4a14-8c24-46b92a225bd3.md) |
| `2bcbf995` | improvement | Distribute hawp bin via core/.hawp; portable wrapper; filename audit | 2026-07-03 | [plan](closed/2026/07/03/2bcbf995-983a-4c23-ab2b-2703ef50d477.md) |
| `582904f4` | improvement | hawp uuid command — easy UUID generation | 2026-07-03 | [plan](closed/2026/07/03/582904f4-1c60-41ba-a357-15f22d7017c7.md) |
| `0e1c4afa` | improvement | UUID IDs Phase 2 — short-UUID display and prefix matching | 2026-07-03 | [plan](closed/2026/07/03/0e1c4afa-9668-4d61-b5b6-1e27be42ca23.md) |
| `361fb08e` | improvement | UUID work item IDs — Phase 1 dual format (was TASK-013) | 2026-07-03 | [plan](closed/2026/07/03/361fb08e-6457-4ed5-80bd-76337b6f0e89.md) |
| TASK-086 | improvement      | Restructure librarian/scripts into hawp/ and librarian/ subfolders        | 2026-07-03 | [plan](closed/2026/07/03/TASK-086.md) |
| TASK-085 | bug              | Fix 7 kit issues found by kit:validate                                    | 2026-07-02 | [plan](closed/2026/07/02/TASK-085.md) |
| TASK-084 | feature          | Implement kit:normalize script                                            | 2026-07-02 | [plan](closed/2026/07/02/TASK-084.md) |
Older closed work (TASK-080 and earlier) is archived under `closed/YYYY/MM/DD/`.

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
- Provider rollout (TASK-072–075), standards adaptation (TASK-070), audit remediation (TASK-076), and scripts structure alignment (TASK-077) are closed.
- TASK-080 and TASK-081 closed 2026-06-29 after a paired-lane workflow-loop trial; see closed plans under `closed/2026/06/29/`.

## Future Improvements

- **UUID migration Phase 3 (optional)**: retroactive migration of closed records to UUID paths — only if librarian/indexing needs it. Phases 1–2 closed 2026-07-03; see `notes/2026/07/03/migration-sequential-to-uuid.md`.
- **Verification clarity cleanup**: `work:validate` now surfaces `86` ambiguous checklist claims plus `1` explicitly unproven claim across closed records. Natural next pass: tighten those historical verification lines with `Evidence:` labels or explicit unproven wording, starting with the 7 most recent `B7` items shown by `./.hawp/bin/hawp work normalize --dry-run --validate`.
- **Optional MCP server overlay for librarian tooling**: a thin MCP wrapper exposing kit/work validate + normalize + uuid as structured tools. CLI stays canonical (works with every agent, no server or per-client config); MCP is an additive overlay if/when a client benefits from structured tool discovery. Decision context: 2026-07-03 discussion, `closed/2026/07/03/2bcbf995-983a-4c23-ab2b-2703ef50d477.md`.
