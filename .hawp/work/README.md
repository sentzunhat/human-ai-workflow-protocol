# `.hawp/work/` — Repo-Local Work State

This folder holds this repository's HAWP operating state: the backlog, decisions, plan files, notes, and verification artifacts.

## Layout

```

.hawp/work/
	BACKLOG.md          — compact active index (not permanent history)
	active/             — open bugs and tasks (flat, easy to find)
	parked/             — deferred or icebox items (not active, not closed)
	closed/YYYY/MM/DD/  — archived closed work, filed by date
	decisions/YYYY/MM/DD/ — ADRs and significant project decisions
	evidence/YYYY/MM/DD/  — verification artifacts (only when real)
	status/YYYY/MM/DD/  — checkpoint summaries and manager reviews
	notes/YYYY/MM/DD/   — session notes, scratch pads, migration records
```

## Core Rule

```
Read from kit.
Write to work.
Archive by date.
Never overwrite project truth.
```

The reusable workflow guides (operating loop, status report shape, repo-local init) live with the kit:

- `../kit/usage/init.md`
- `../kit/usage/intake-workflow.md`
- `../kit/usage/status-report.md`

## Distinction from `kit/`

- `kit/` is reusable HAWP material (protocol, templates, patterns, reviews, examples, types, and reusable workflow usage guides) that downstream installs copy.
- `work/` is this repo's operating memory. Only the folder scaffolding (README files and a starter `BACKLOG.md`) is seeded by the install; all content you create here belongs to your repo and is never overwritten by HAWP updates.
