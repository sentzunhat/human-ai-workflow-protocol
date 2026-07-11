# Task: Audit provider/distribution sync and fix kit layout drift

**Backlog ID:** TASK-078

## Context

The kit restructure moved `spec.md` and `authoring-patterns.md` from the top level of `core/.hawp/kit/` into `references/`, with matching updates to the shared intake behavior, distribution sources, generated guides, and benchmark docs. A full audit was requested to verify that distribution generated outputs and sources match the provider dot folders at repo root, that root `.hawp/kit` matches `core/.hawp/kit`, that the kit top level holds only `README.md` and `start-here.md` plus category folders, and to fix any remaining drift.

## Scope

- Verify `core/.hawp/kit` ↔ `.hawp/kit` sync and the new `references/` placement of `spec.md` and `authoring-patterns.md`.
- Verify provider packs (`core/providers/.github|.cursor|.continue`) match root dot folders, `AGENTS.md`, and `.github/copilot-instructions.md`.
- Verify `distribution/generated/` is current against `distribution/sources/` and materialized provider files are current against shared behaviors.
- Fix any drift found.

## Outcome

- Kit sync verified: `diff -rq core/.hawp/kit .hawp/kit` clean; kit top level contains only `README.md`, `start-here.md`, and category folders; `spec.md` and `authoring-patterns.md` live in `references/` with all links updated (kit README, start-here, usage/init, intake behavior, benchmark docs, install/update script sources).
- Provider sync verified: cursor rules, continue rules, github instructions/prompts, copilot-instructions, and `AGENTS.md` seed all match their pack sources.
- Drift fixed: `.github/instructions/code-style.instructions.md` was root-only, had broken frontmatter, and referenced nonexistent `shared_standards/public/guidelines/code-style.md`. Rewrote it with valid frontmatter pointing at `.hawp/kit/standards/guidelines/code-style.md`, aligned its key rules to that standard, and added it to `core/providers/.github/instructions/` so pack and root match.
- The only remaining old-path reference is in the closed archive `closed/2026/06/04/TASK-072.md`, preserved intentionally as history.

## Verification

- [x] `providers:validate` — 11 materialized files current. Evidence: ../evidence/2026/06/12/TASK-078-provider-kit-sync-audit-verification.md
- [x] `distribution:validate` — generated outputs current. Evidence: same file.
- [x] `workflow:validate` — PASS (0 issues). Evidence: same file.
- [x] `typecheck` exit 0 and `npm test` 38/38 pass on Node 22. Evidence: same file.

## Close Checklist

- [x] Outcome recorded.
- [x] Evidence stored under `evidence/2026/06/12/`.
- [x] BACKLOG.md Recently Closed updated and capped at 10 rows.
- [x] STATUS.md refreshed.
