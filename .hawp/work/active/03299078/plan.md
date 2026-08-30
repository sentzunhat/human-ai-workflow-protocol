# v0.0.23 TypeScript script deprecation plan and Go parity audit

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `03299078-a456-4e80-8445-cf28521eec31`
**Type:** improvement
**Reported:** 2026-08-30

---

## Input (verbatim)

> Plan the deprecation of legacy TypeScript workflow scripts in librarian/scripts by mapping each remaining script surface to its Go replacement or retained necessity, defining the removal order, validation gates, and repo-structure implications without breaking distribution/provider pipelines.

## Intake Summary

The Go CLI now covers the workflow command surface that was previously carried
by TypeScript under `librarian/scripts/hawp/`, but the repository still keeps
those TypeScript implementations and tests checked in beside newer Go ports.
This item will define the deprecation/removal order, what must stay Node-based,
and what proof gates are required before we retire the remaining TS workflow
script tree.

## Current Context

- `librarian/src/README.md` states all read/mutate/composite workflow commands
  are ported to Go: `uuid`, `links check`, `kit validate/normalize`,
  `work validate/normalize`, and `check`.
- `librarian/package.json` routes the user-facing workflow npm scripts through
  the Go CLI wrapper, while distribution/provider materialization remains in
  TypeScript under `librarian/scripts/librarian/`.
- `librarian/scripts/README.md` still documents TypeScript workflow domains
  under `librarian/scripts/hawp/`, which now function more like retained legacy
  implementations than the primary runtime path.
- On 2026-08-30, `go run ./cmd/hawp work normalize --apply --migrate-folders
  --force-dirty --validate` returned "No work-record changes were necessary"
  and `VALIDATION PASS (0 issues, 1 warnings)`, confirming the UUID-folder
  migration lane is stable and idempotent in this workspace.

## Initial Analysis

**Directly verified:**

- `librarian/scripts/hawp/` still contains these TypeScript command domains:
  `hawp-check`, `kit-normalize`, `kit-validate`, `work-normalize`,
  `work-validate`.
- The active Go implementation now lives under `librarian/src/internal/`
  across application/domain/platform layers, including the new
  `normalize_migrate.go` folder-migration support.
- `go test ./internal/domain/work ./internal/application/work ./internal/platform/cli`
  passed on 2026-08-30.
- The live repo run of `go run ./cmd/hawp work normalize --apply --migrate-folders --force-dirty --validate`
  made no further folder changes after the migration fix, which is the main
  proof that the work-item folder layout upgrade is working.
- Some recent closed-plan files were modified by an earlier combined normalize
  run before `--migrate-folders` was isolated from closed-record normalization;
  those accidental archive edits were reverted from the working tree and should
  stay out of the TS deprecation lane.

**Inferred (not yet proven):**

- The TypeScript workflow script tree can probably be deprecated in phases:
  first documentation/status demotion, then wrapper/test retirement, then
  source deletion after parity evidence and downstream surface checks.
- `librarian/scripts/lib/` may still contain reusable helpers that either need
  a final Go replacement audit or a conscious retention decision for
  Node-based distribution/provider pipelines.
- `librarian/scripts/librarian/` should likely remain for now because the
  provider/distribution generation pipeline is still Node/TS-owned and is not
  covered by the current Go CLI port.

**Likely scope:**

- Build a command-by-command parity matrix: TypeScript implementation, Go
  replacement, current invocation path, remaining tests/docs, and removal risk.
- Separate "deprecated workflow scripts" from "retain Node maintainer tooling"
  so we do not accidentally target distribution/provider scripts that still own
  real release-generation behavior.
- Update repo docs to make the Go path authoritative and the legacy TS workflow
  path explicitly transitional.
- Propose the retirement sequence and verification gates before deleting code.

## Risk + Review Gate

**Risk:** medium
**Gate:** review first on medium/high

## Backlog + Plan Link

**Status now:** analyzing
**Plan file:** work/active/03299078/plan.md

## Next Step

- [x] Investigation recorded above (required before planning)
- [ ] Build the TS-to-Go parity matrix and classify each remaining script as
      remove, retain, or migrate later
- [ ] Decide whether to track TS deprecation as one item or split by workflow
      surface vs maintainer pipelines
- [ ] Move backlog status to plan-ready once the removal order and proof gates
      are explicit
