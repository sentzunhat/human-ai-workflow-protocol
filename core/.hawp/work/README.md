# `.hawp/work/` — Repo-Local Work State

This folder holds this repository's HAWP operating state: the backlog, decisions, plan files, notes, and verification artifacts.

## Layout

```
work/
  STATUS.md           — current state dashboard
  BACKLOG.md          — single source of truth for all work items
  active/             — open bugs and tasks (flat, easy to find)
  closed/YYYY/MM/DD/  — archived closed work, filed by date
  decisions/YYYY/MM/DD/ — ADRs and significant project decisions
  evidence/YYYY/MM/DD/  — verification artifacts (only when real)
  notes/YYYY/MM/DD/   — session notes, scratch pads, migration records
```

The reusable workflow guides (operating loop, status report shape, repo-local init) live with the kit:

- `../kit/usage/INIT.md`
- `../kit/usage/INTAKE_WORKFLOW.md`
- `../kit/usage/STATUS_REPORT.md`

## Core Rule

```
Read from kit.
Write to work.
Archive by date.
Never overwrite project truth.
```

## Distinction from `kit/`

- `kit/` is reusable HAWP material (protocol, templates, patterns, reviews, examples, types, and reusable workflow usage guides) that downstream installs copy.
- `work/` is this repo's operating memory. Only the folder scaffolding (README files and a starter `BACKLOG.md`) is seeded by the install; all content you create here belongs to your repo and is never overwritten by HAWP updates.
