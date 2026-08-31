# agents-seed-if-missing — AGENTS.md must not be overwritten on update

**Type:** fix  
**Status:** done  
**Opened:** 2026-08-25  
**Closed:** 2026-08-25
**Target:** v0.0.11

**Closes:** `agents-seed-if-missing`

## Goal

`hawp init` (including `--provider` re-runs and `hawp update --provider`) must
never overwrite an existing `AGENTS.md`. The correct behavior is identical on
install and update: seed only when the file is absent.

## Outcome

Cursor and Codex provider overlays now treat `AGENTS.md` as seed-if-missing on
both install and update. The kitsync manifest now supports update-time
`seed-if-missing` semantics directly, so the CLI sync path and the generated
shell guides follow the same rule: refresh provider rules, but preserve an
existing repo-local `AGENTS.md`.

This removes the previous split-brain behavior where install was safe but
update/re-init could clobber a blended or product-owned `AGENTS.md`.

## Verification

- [x] `go test ./internal/domain/kitsync/... ./internal/application/kitsync/... ./internal/platform/cli/... ./internal/platform/mcp/...`. Evidence: the command is recorded in this plan's Verification section.
- [x] Added kitsync tests proving update-time seed-if-missing does not overwrite an existing `AGENTS.md`. Evidence: this plan's What was done section names the new seed-if-missing coverage.
- [x] Added kitsync tests proving update-time seed-if-missing creates `AGENTS.md` when it is absent. Evidence: this plan's What was done section names the new seed-if-missing coverage.
- [x] `mise exec node@26.5.0 -- npm run distribution:sync`. Evidence: the command is recorded in this plan's Verification section.

## What was done

- Added `UpdateMode()` support for `seed-if-missing` in `librarian/src/internal/domain/kitsync/manifest.go`
- Updated provider update application in `librarian/src/internal/domain/kitsync/apply.go`
- Changed `core/providers/manifest.yaml` so Cursor and Codex `AGENTS.md` rules use update-time `seed-if-missing`
- Updated source templates and generated Cursor/Codex update guides to use no-clobber AGENTS behavior
- Updated provider safety/boundary/doc-verification text to match the new rule

## Close Checklist

- [x] Outcome recorded
- [x] Verification includes code and generated-output evidence
- [x] Existing `AGENTS.md` preservation rule captured explicitly
- [x] Ready to stay in closed history
