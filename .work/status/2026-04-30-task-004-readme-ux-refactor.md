# TASK-004 — Root README UX refactor

## Input

User: "refactor the root readme for better ux ... simple install and update steps and some context of the project ... dont refactor too much ... as simple as using install or update readmes you can get started as simple as copy and paste".

## Context

- Current `README.md` is long and front-loads repo layout.
- `core/install.md` and `core/update.md` already contain canonical copy/paste scripts.
- Intent: make the root README a fast on-ramp that pulls attention to (1) what HAWP is, (2) one-shot install, (3) one-shot update, with links into the canonical docs for depth.

## Mission

Rewrite `README.md` so a human or AI model can:

1. Understand HAWP in ~5 lines.
2. Install via a single copy/paste block (mirrors `core/install.md`).
3. Update via a single copy/paste block (mirrors `core/update.md`).
4. Find depth via a short links section.

## Constraints

- Keep scope to `README.md` only.
- Do not change `core/install.md` or `core/update.md` content.
- Preserve the Apache 2.0 license note.
- Keep references to existing paths (`core/.hawp/...`).
- Tight, concise style; no over-explanation.

## Output

- New `README.md` with: short pitch, what it is (5-field shape), install block, update block, "Use it" 3-step, links, license.
- Backlog row added under Done.

## Risk

Low. Docs-only change in a single file. No code paths touched. Reversible by git.

## Evidence vs inference

- Evidence: `README.md`, `core/install.md`, `core/update.md` already read in this session.
- Inference: a shorter README improves agent + human discoverability. Not measured.
