# TASK-086 — Restructure librarian/scripts into hawp/ and librarian/ subfolders

**Backlog ID:** TASK-086

**Type:** improvement  
**Status:** inbox  
**Created:** 2026-07-02

## Goal

The current `librarian/scripts/` folder mixes HAWP-workflow scripts (kit:validate, kit:normalize, work:validate, work:normalize, hawp:check) with librarian-specific scripts (distribution, providers). Separate them so domain ownership is clear.

## Proposed layout

```
librarian/scripts/
├── hawp/           ← HAWP workflow scripts
│   ├── kit-validate/
│   ├── kit-normalize/
│   ├── work-validate/
│   ├── work-normalize/
│   └── hawp-check/
├── librarian/      ← Librarian tooling scripts
│   ├── distribution/
│   └── providers/
├── lib/            ← Shared utilities (stays flat)
└── check-markdown-links.mjs
```

## Steps

- [ ] Investigate: confirm no cross-folder imports would break (lib/ is the only shared dependency)
- [ ] Move `hawp-check/`, `kit-validate/`, `kit-normalize/`, `work-validate/`, `work-normalize/` → `hawp/`
- [ ] Move `distribution/`, `providers/` → `librarian/` subfolder
- [ ] Update all `package.json` script paths
- [ ] Update `librarian/scripts/README.md`
- [ ] Update `librarian/README.md` and `CLAUDE.md`
- [ ] Run `npm --prefix librarian run validate` — must pass
- [ ] Close TASK-086

## Constraints

- Do NOT change any logic, only move files
- Go slowly — this touches many paths; run typecheck after each group of moves
- `lib/` stays flat (shared by both groups)

## Outcome

Restructured `librarian/scripts/` into two owner groups plus shared code:

- `scripts/hawp/` — kit-validate, kit-normalize, work-validate, work-normalize, hawp-check
- `scripts/librarian/` — distribution, providers
- `scripts/lib/` — stays flat, shared by both groups
- `check-markdown-links.mjs` and `README.md` stay at the top level

No logic changes — only moves, import-depth bumps for `lib/`, and path updates in `package.json`, `hawp-check` spawn paths, CI workflow filters, and docs.

Also repaired drift found on the way (state lost in the `WIP on dev` commit from the previous session):

- Deleted stale duplicate folders `validate-hawp-workflow/` and `backlog-upgrade/` (old pre-rename copies).
- Fixed `hawp-check` spawning the deleted `validate-hawp-workflow` path — `hawp:check` would have broken.
- Added missing `distribution:validate` npm script (referenced by `distribution:sync` but never defined).
- Restored BACKLOG rows for TASK-082–085 and the TASK-086 active row.
- Filled missing Outcome/Verification in closed TASK-079 (2026-06-29).
- Re-applied the `ADR.template.md` → `adr-template.md` rename (old file had reappeared as a duplicate) and updated export manifests.
- Restored the executable bit on `.hawp/bin/hawp`.

## Verification

- `npm --prefix librarian run validate` — full suite PASS (typecheck, 46 tests, markdown links across 306 files, kit:validate 3/3, distribution:sync clean, work:validate 4 checks / 0 issues).
- `npm --prefix librarian run hawp:check` — PASS with new spawn paths.
- `./.hawp/bin/hawp kit validate` and `work validate` — PASS via wrapper.
- Evidence: `evidence/2026/07/03/TASK-086-verification.md`.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Plan file moved to `closed/YYYY/MM/DD/`
- [x] BACKLOG.md updated
