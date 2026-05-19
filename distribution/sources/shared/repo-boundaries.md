# Repository Boundaries: What Gets Installed and Updated

The install and update operations respect clear boundaries to protect your project-specific work.

## HAWP-Managed Files (Refreshed on Every Install or Update)

These files are owned by HAWP and are written or replaced on every run:

- `.hawp/LICENSE` — Apache 2.0 license for the HAWP kit content
- `.hawp/kit/**` — Reusable HAWP protocol docs, templates, patterns, examples, references (replaced in full)
- `.github/instructions/*.instructions.md` — HAWP Copilot instructions
- `.github/prompts/*.prompt.md` — HAWP Copilot prompt templates

## Project-Owned Files (Never Overwritten)

These files are owned by your project and are never overwritten by install or update:

- `.hawp/work/**` — Your BACKLOG, active work, parked work, closed work, decisions, evidence, notes
- `.github/copilot-instructions.md` — Your custom Copilot instructions (seeded once when missing; never overwritten after)

## Scaffold Files (Seeded Once; Never Overwritten After)

These are created on first install if they do not already exist:

- `.hawp/work/README.md`, `.hawp/work/STATUS.md`, `.hawp/work/BACKLOG.md`
- `.hawp/work/active/README.md`, `parked/README.md`, `closed/README.md`
- `.hawp/work/decisions/README.md`, `evidence/README.md`, `status/README.md`, `notes/README.md`
- `.github/copilot-instructions.md`

## What Never Gets Installed Downstream

These files from the HAWP source repository are **never** copied to your target:

- The HAWP source repo's own `.hawp/work/BACKLOG.md`
- The HAWP source repo's own `.hawp/work/active/`, `closed/`, `decisions/`, `evidence/` operating state
- Any content from `benchmark/` (optional reference material, not installed)

## The Boundary Model

```
HAWP Source Repo (core/)
├─ .hawp/kit/             → Reusable HAWP assets (refreshed every install/update)
├─ .hawp/work/            → HAWP source repo's own operating state (NEVER installed)
└─ .github/               → HAWP overlay templates (refreshed every install/update)

Target Repo (your project)
├─ .hawp/kit/             ← Refreshed from HAWP source on every install/update
├─ .hawp/work/            ← YOUR project work (never overwritten; scaffold seeded once)
├─ .github/instructions/  ← HAWP overlays refreshed; copilot-instructions.md preserved
└─ [your files]           ← Your code, docs, and config (never touched)
```
