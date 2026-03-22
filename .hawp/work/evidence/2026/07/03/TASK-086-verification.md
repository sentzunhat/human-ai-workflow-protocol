# TASK-086 Verification Evidence — 2026-07-03

## Direct evidence

- `npm --prefix librarian run validate` → PASS end-to-end:
  - typecheck: clean (`tsc --noEmit`)
  - tests: 46/46 pass, 0 fail
  - check:markdown-links: 306 files, all local links valid
  - kit:validate: 3 checks passed, 0 issues
  - distribution:sync (providers:materialize → providers:validate → distribution:build → distribution:validate): clean
  - work:validate: 4 checks passed, 0 issues, 1 tolerated legacy warning
- `npm --prefix librarian run hawp:check` → `Result: PASS` with the new spawn paths (`scripts/librarian/distribution/validate`, `scripts/hawp/work-validate`)
- `./.hawp/bin/hawp kit validate` and `./.hawp/bin/hawp work validate` → PASS after restoring the executable bit

## Final layout

```
librarian/scripts/
├── hawp/           kit-validate, kit-normalize, work-validate, work-normalize, hawp-check
├── librarian/      distribution, providers
├── lib/            shared utilities (flat)
├── check-markdown-links.mjs
└── README.md
```

## Drift repaired in the same session (inference: lost in the prior `WIP on dev` commit)

- Stale duplicates `scripts/validate-hawp-workflow/` and `scripts/backlog-upgrade/` deleted (current code lives in `work-validate/`/`work-normalize/`).
- `hawp-check/script.ts` pointed at the deleted `validate-hawp-workflow` path — fixed.
- `distribution:validate` npm script was referenced by `distribution:sync` but undefined — added.
- BACKLOG rows for TASK-082–086 restored; closed TASK-079 (2026-06-29) Outcome/Verification filled.
- `ADR.template.md` duplicate removed (canonical `adr-template.md` kept); export manifests and docs README updated.
