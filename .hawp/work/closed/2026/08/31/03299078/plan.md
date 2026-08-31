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

**Status now:** done
**Plan file:** work/closed/2026/08/31/03299078/plan.md

## Next Step

- [x] Investigation recorded above (required before planning)
- [x] Build the TS-to-Go parity matrix and classify each remaining script as
      remove, retain, or migrate later
- [x] Decide whether to track TS deprecation as one item or split by workflow
      surface vs maintainer pipelines
- [x] Move backlog status to plan-ready once the removal order and proof gates
      are explicit

## 2026-08-31 Parity Matrix

| Surface | Go replacement | Current invocation path | Decision |
| ------- | -------------- | ----------------------- | -------- |
| `hawp/kit-validate` | `hawp kit validate` | npm wrapper calls Go CLI | deprecated now; delete in a later cleanup patch |
| `hawp/kit-normalize` | `hawp kit normalize` | npm wrapper calls Go CLI | deprecated now; delete in a later cleanup patch |
| `hawp/work-validate` | `hawp work validate` | npm wrapper calls Go CLI | deprecated now; keep only as transitional source |
| `hawp/work-normalize` | `hawp work normalize` | npm wrapper calls Go CLI | deprecated now; keep only as transitional source |
| `hawp/hawp-check` | `hawp check` | npm wrapper calls Go CLI | deprecated now; delete in a later cleanup patch |
| `librarian/distribution/*` | none | npm TypeScript pipeline | retain; release-generation tooling still lives here |
| `librarian/providers/materialize/*` | none | npm TypeScript pipeline | retain; provider/distribution tooling still lives here |
| `lib/*` | none | shared by retained Node tooling | retain until Node maintainer pipelines are retired |

## Decision

- Treat TypeScript workflow deprecation as complete for `v0.0.23` at the
  documentation and ownership layer: the Go CLI is authoritative, npm wrappers
  already route through it, and `librarian/scripts/README.md` now describes the
  `hawp/` tree as deprecated legacy source instead of the primary runtime path.
- Keep actual TypeScript source deletion out of this release because
  distribution/provider pipelines still keep Node tooling live in the same
  workspace and there is no release need to mix deletion risk into this patch.

## Outcome

Closed 2026-08-31.

The TypeScript workflow deprecation lane is complete as a release-readiness
audit. The remaining TypeScript workflow command surfaces are now explicitly
mapped to their Go replacements, the retained Node-only pipelines are separated
from deprecated workflow code, and the live scripts documentation no longer
implies that `librarian/scripts/hawp/` is the primary path.

## Verification

- [x] `librarian/src/README.md` already marks the workflow command set as ported
      to Go. Evidence: [README.md](/Users/beltrd/Desktop/projects/sentzunhat/human-ai-workflow-protocol/librarian/src/README.md)
      command matrix shows `uuid`, `links check`, `kit/work validate`,
      `kit/work normalize`, and `check` as available via Go
- [x] `librarian/scripts/README.md` now labels `librarian/scripts/hawp/` as a
      deprecated legacy implementation surface. Evidence:
      [README.md](/Users/beltrd/Desktop/projects/sentzunhat/human-ai-workflow-protocol/librarian/scripts/README.md)
      documents Go ownership for user-facing workflow commands
- [x] `librarian/src/CHANGELOG.md` records the v0.0.23 TypeScript deprecation
      checkpoint. Evidence: [CHANGELOG.md](/Users/beltrd/Desktop/projects/sentzunhat/human-ai-workflow-protocol/librarian/src/CHANGELOG.md)
      includes the TS script deprecation note under `0.0.23`

## Close Checklist

- [x] Parity matrix recorded
- [x] Retained Node-only boundaries made explicit
- [x] Live docs reflect the release-ready ownership model
