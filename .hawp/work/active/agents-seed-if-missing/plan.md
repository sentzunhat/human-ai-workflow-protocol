# agents-seed-if-missing — AGENTS.md must not be overwritten on update

**Type:** fix  
**Status:** plan-ready  
**Opened:** 2026-08-25  
**Target:** v0.0.11

## Input

Downstream install of hawp 0.0.9 with `--provider cursor`. On first `hawp init`,
`AGENTS.md` is correctly seeded only when missing. On subsequent `hawp init`
(update / re-init), the kit sync step overwrites `AGENTS.md` with the upstream
seed, discarding any product-specific content the repo had blended in.

Downstream repos using a blended `AGENTS.md` (e.g. product law + HAWP
work-tracking section) lose their additions on every update.

Evidence source: downstream install evidence 2026-08-25.

## Goal

`hawp init` (including `--provider` re-runs and `hawp update --provider`) must
never overwrite an existing `AGENTS.md`. The correct behavior is identical on
install and update: seed only when the file is absent.

## Constraints

- `.hawp/work/**` is project-owned and must not be overwritten — this rule already exists; apply the same protection to `AGENTS.md`.
- Do not add a required interactive prompt; silently skip is correct.
- An opt-out flag (`--force-agents` or `--overwrite-agents`) is acceptable as a
  deliberate escape hatch but must default to skip.
- This fix applies to all providers that seed `AGENTS.md` (cursor, and any future providers).

## Plan

### Step 1 — Locate the `AGENTS.md` write path

Find where `hawp init` and kit sync write or copy `AGENTS.md`. Likely in
`librarian/src/internal/application/kitsync/` or `provision/`.

### Step 2 — Add seed-if-missing guard

Before writing `AGENTS.md`, check if the file exists. If it does, skip the
write and (optionally) print a note: `AGENTS.md: already exists, skipping`.
Apply the guard to every code path that writes `AGENTS.md`.

### Step 3 — Test

- Unit test: `hawp init --provider cursor` on a repo with existing `AGENTS.md`
  → file unchanged.
- Unit test: `hawp init --provider cursor` on a fresh repo → file written.
- Idempotency: running twice on a fresh repo → one file, same content.

### Step 4 — Kit docs

Add a note to `usage/search.md` Cursor section and/or install guide:

> If you have a custom `AGENTS.md`, `hawp init` will not overwrite it.
> To get the latest HAWP kit guidance, check `.hawp/kit/AGENTS.md.seed`
> and merge manually.

## Verification

- Run `hawp init --provider cursor` twice on a repo with a custom `AGENTS.md`
- Confirm the file is unchanged after the second run
- Confirm a fresh repo gets the seed file written
