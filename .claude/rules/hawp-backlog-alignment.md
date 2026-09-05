---
description: "HAWP backlog alignment and compaction guardrails"
---

<!-- Generated from core/providers/shared/behaviors - edit shared source and run `hawp providers sync` -->

# HAWP Backlog Alignment

`BACKLOG.md` is an active coordination index, not permanent history.

Rules:

- Keep Active Work short and current.
- On close, move detail files to `.hawp/work/closed/YYYY/MM/DD/`.
- Cap Recently Closed to the last 5–10 items (or the last 14–30 days).
- Store verification evidence in `.hawp/work/evidence/YYYY/MM/DD/{uuid}/evidence.md` (use `hawp work evidence --title "<label>"`).
- Store checkpoint summaries in `.hawp/work/status/YYYY/MM/DD/{uuid}/status.md` (use `hawp work status --title "<label>"`).
- Store decision records in `.hawp/work/decisions/YYYY/MM/DD/{uuid}/decision.md` (use `hawp work decision --title "<label>"`).
- Store notes in `.hawp/work/notes/YYYY/MM/DD/{uuid}/note.md` (use `hawp work note --title "<label>"`).
- Preserve history in archive files; never delete records just to shorten the backlog.
- If `BACKLOG.md` already has many Done rows, create or recommend a work item titled `Compact BACKLOG.md and archive closed work.` before adding more Done rows.

Canonical policy: `.hawp/kit/references/backlog-alignment.md`
