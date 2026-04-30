# `.hawp/work/` — Repo-Local Work State

This folder holds this repository's HAWP operating state: the backlog, decisions, plan files, and verification artifacts.

## Layout

- `BACKLOG.md` — single source of truth for active and completed work in this repo.
- `adrs/` — repo-local ADRs (project decisions about HAWP itself).
- `status/` — plan files and status reports per backlog item.
- `evidence/` — verification artifacts referenced from plans (only when needed; never fabricated).

The reusable workflow guides (operating loop, status report shape, repo-local init) live with the kit:

- `../kit/usage/INIT.md`
- `../kit/usage/INTAKE_WORKFLOW.md`
- `../kit/usage/STATUS_REPORT.md`

## Distinction from `kit/`

- `kit/` is reusable HAWP material (protocol, templates, patterns, reviews, examples, types, and reusable workflow usage guides) that downstream installs copy.
- `work/` is this repo's operating memory. Only the folder scaffolding (README files and a starter `BACKLOG.md`) is seeded by the install; all content you create here belongs to your repo and is never overwritten by HAWP updates.
