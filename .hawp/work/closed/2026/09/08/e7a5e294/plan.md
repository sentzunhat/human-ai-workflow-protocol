# Group CLI commands into feature-owned packages

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `e7a5e294-97a8-44d1-88f4-bbba7c6c38a5`
**Type:** improvement
**Reported:** 2026-09-07
**Closed:** 2026-09-08

---

## Input (verbatim)

> Decide package boundaries across layers; plan folder-within-folder organization that stays simple, repeatable, and avoids duplication. Planning only now; pace usage and keep implementation slices small.

## Investigation

Confirmed: `librarian/src/internal/platform/cli` contains 30 non-test Go files.
Work parsing/handlers share the root with unrelated commands; several root helpers
are private and must not become child-to-parent imports.

This is a distinct implementation subitem of
[47c793d6](../47c793d6/plan.md), following the
[shared package map](../47c793d6/package-boundaries.md).

## Decision And Plan

Current user-approved convention: use `cli/<family>/<operation>/` with simple
folder names, colocating handlers, parsers and tests. Apply retrospectively to
completed CLI slices and to future slices. Keep only forwarding entrypoints at
the family root; shared code needs a specific owner and no child-to-parent import.
This supersedes earlier flat-family checkpoints below.

Current correction scope (all relative to `librarian/src/internal/platform/cli`):
split `index/index_commands.go` into `index/build/command.go` and
`index/ingest/command.go`; move `index/index_build_args.go` to `index/build/args.go`
and `index/index_corpus{,_test}.go` to `index/ingest/corpus{,_test}.go`.
Split `model/model_commands.go` into `model/{pull,embed,search-embed}/command.go`,
moving each corresponding parser and existing parser test to `args{,_test}.go`.
Extract the shared model path resolver to `model/location/location.go`.
Move root `kit_{normalize,validate}_args{,_test}.go` into the corresponding kit
operation as `args{,_test}.go`, replacing duplicate parsers in child handlers.
Move `mutation_boundary_args_test.go` to `kit/mutation_boundary_test.go`, using
command entrypoints for cross-operation checks. Retain root routing tests.

1. Move only work parsers, handlers and their tests to planned
   `librarian/src/internal/platform/cli/work/`.
2. Keep `Run`, routing, registry/help and public command contracts at the CLI root.
   Export only command entrypoints needed by routing. Resolve shared ExitError
   ownership with a minimal platform-only type if required; keep a root alias
   for compatibility. Never import the CLI root from work.
3. Validate the pattern, then migrate kit/link, search/index/model, and remaining
   command owners as separate slices. Do not split individual flags into packages.

Alternative: retain flat packages with renamed files. Rejected for these scoped
owners because it leaves navigation/ownership concerns unresolved. Do not nest
already-cohesive small packages or duplicate rules between layers.

## File Ownership And Coordination

Own `librarian/src/internal/platform/cli/work_commands.go`, `work_*_args.go`,
associated parser tests, `run.go` dispatch changes, and planned `work/`.
Later slices own their matching *_commands.go/*_args.go and tests.
Read application/work contracts; do not alter their policies in this item.
All shorthand directories above are under `librarian/src/internal` unless stated.
Before implementation, enumerate exact source/test/destination paths for the
selected slice. No broad simultaneous moves; preserve the current dirty checkpoint.

## Implemented First Slice

Moved the work command adapters into the feature-owned package:

- Added `librarian/src/internal/platform/cli/work/commands.go` and its
  `new_args.go`, `normalize_args.go`, and `validate_args.go` parser files.
- Moved the matching parser and mutation-boundary tests into
  `librarian/src/internal/platform/cli/work/`.
- Updated `librarian/src/internal/platform/cli/run.go` to retain command
  routing while delegating work commands with the current working directory.
- Added `librarian/src/internal/platform/exitcode/error.go`; the root
  `cli.ExitError` is an alias, preserving the existing executable contract
  without a child-to-root import.
- Removed the superseded root `work_commands.go` and `work_*_args.go` files.

The exact moved files are confined to this first command family. Application
work policy and domain code were not changed.

## Risk And Execution Gate

Risk: medium (Go package/import topology). Status: in-progress. Owner: Codex.
The first work-command slice is implemented. Order remains: validate this
pattern, then migrate kit/link, search/index/model, and runtime adapters as
separate slices; shared entrypoint changes are serialized.

## Verification / Acceptance

Existing parser/error/no-write tests and black-box CLI work commands pass;
root Run signature and command output remain compatible; no import cycle.
Full Go suite/vet and HAWP validation pass at each completed command-family slice.
Focused verification passed from `librarian/src`:
`go test ./internal/platform/cli/... ./tests/platform/cli/...`.
The final HAWP validation verifies record integrity, not package migration correctness.

## Usage Budget

One cohesive slice per turn. Read only owned files, direct dependencies and tests.
Short updates; update this plan rather than creating a report for each move.
Focused tests during edits, full checks once at slice completion. No recurring
usage monitor, reset redemption, model switch, or paid action requested.

## Completed Families

| Family | Nested operations | Verification |
| --- | --- | --- |
| `work/` | `new/`, `normalize/`, `validate/` | full suite + HAWP ✓ |
| `kit/` | `normalize/`, `validate/` | full suite + HAWP ✓ |
| `index/` | `build/`, `ingest/` | full suite + HAWP ✓ |
| `model/` | `pull/`, `embed/`, `search-embed/`, `location/` | full suite + HAWP ✓ |
| `links/` | `check/`, `clean/` | full suite + HAWP ✓ |
| `search/` | `query/`, `benchmark/` | full suite + HAWP ✓ |
| `mcp/` | `configure/` | full suite + HAWP ✓ |
| `usage/` | `enable/` | full suite + HAWP ✓ |
| `update/` | flat `commands.go`; `init` → `cli/init/command.go` | full suite + HAWP ✓ |
| `distribution/` | flat `commands.go`; imports `providers` sibling | full suite + HAWP ✓ |
| `providers/` | flat `commands.go` | full suite + HAWP ✓ |
| root cleanup | `model_commands.go` + `maintenance_commands.go` deleted; all logic inlined into `run.go` | full suite + HAWP ✓ |

## Remaining Slice Queue

All families migrated. Root cleanup complete. — `distribution_commands.go` (144 L).
- All six operations (`providers materialize/validate/sync` and
  `distribution build/validate/sync`) are argless; no dedicated parsers.
- Create `cli/distribution/commands.go` and `cli/providers/commands.go`, each with
  flat handler functions. No operation subfolders needed.
- Update `run.go`; delete root `distribution_commands.go`.
- Risk: low.

**5. Root cleanup** — after all families are moved:
- Inline `model_commands.go` (3 thin forwarders) directly into `run.go`.
- Consolidate or remove `maintenance_commands.go` residuals (`runUUID`, `runCheck`).
- Goal: `run.go` is the sole routing file; no miscellaneous `*_commands.go` at root.

## Next Step

All slices complete. `run.go` is now the sole routing file; no miscellaneous
`*_commands.go` files remain at the CLI root. See
[package-boundaries.md](../47c793d6/package-boundaries.md).

Retrospective verification: full `go test ./...`, `go vet ./...` and
`git diff --check` passed. Existing parser and corpus tests moved alongside their
owners; kit cross-operation no-write tests now call the actual child handlers.
Root command routing tests remain at the root. No model inference or download
smoke was performed for this mechanical reorganization.

## Earlier Family Checkpoints

The next `kit` slice is now implemented as `cli/kit/` with nested `normalize/`
and `validate/` command packages. Root maintenance handlers remain compatibility
forwarders, and existing root parser tests remain untouched for this bounded step.
Focused CLI tests and `git diff --check` passed.

The next compounding slice moves the cohesive index command family to
`cli/index/`: index build parsing, corpus construction, search-index ingestion,
and the matching corpus tests. Root routing continues to own command dispatch.

Implemented on 2026-09-07. `RunBuild` and `RunSearchIndex` receive the working
directory from root routing; the child does not import its parent. Focused CLI
tests, full `go test ./...`, `go vet ./...`, and `git diff --check` passed.
The next bounded command family is model/embed, subject to an import and test
ownership check before moving it.

Implemented the model/embed slice on 2026-09-07. Model pull, embed, search embed,
their argument parsers, and parser tests now live in `cli/model/`; root routing
uses compatibility forwarders. Focused tests, full `go test ./...`, `go vet ./...`,
and `git diff --check` are required before closing this slice.

## Requested Subcommand Nesting

User requested one additional folder level for new, normalize, and validate.
Inspection confirms each handler and parser is independent. Keep the existing
work entrypoints as forwarding functions; child packages never import their parent.

Exact scope under `librarian/src/internal/platform/cli/work/`:
- Split `commands.go` handlers into `new/command.go`, `normalize/command.go`,
  and `validate/command.go`; retain the forwarding entrypoints in `commands.go`.
- Move `new_args.go` and `new_args_test.go` to `new/args.go` and `new/args_test.go`.
- Move `normalize_args.go` and `normalize_args_test.go` to
  `normalize/args.go` and `normalize/args_test.go`.
- Move `validate_args.go` and `validate_args_test.go` to
  `validate/args.go` and `validate/args_test.go`.
- Move `mutation_boundary_test.go` to `new/mutation_boundary_test.go`.

This user-requested nesting supersedes the earlier preference against nesting
these small packages. No additional command families are included in this slice.

Implemented the nesting above on 2026-09-07. Each child exports only `Run`;
parsers and options remain private. Existing tests moved with their owner.
`go test ./...`, `go vet ./...`, and `git diff --check` passed.
HAWP validation before the changes passed all three checks with zero issues.

User naming refinement: renamed the CLI folder and package to `work`, retaining
a `workcmd` import alias only in root routing. Final paths are `cli/work/new/`,
`cli/work/normalize/`, and `cli/work/validate/`. Updated the shared package map;
continued this existing item without creating duplicate work. Keep later slices
bounded to conserve usage. Post-nesting HAWP checks passed with zero issues.

After the path rename, focused CLI/unit and black-box tests, `git diff --check`,
and all three HAWP validations passed. The full suite and vet passed immediately
before the path-only rename.

## Outcome

All CLI command families migrated from flat root files to nested `cli/<family>/` packages. Eleven families complete (work, kit, index, model, links, search, mcp, usage, update, distribution, providers); each with forwarding entrypoints at the family root, handler/parser/tests colocated in operation subfolders. Root routing file `run.go` is the sole routing file; no miscellaneous `*_commands.go` remain. Full `go test ./...`, `go vet ./...`, and HAWP checks passed at each slice boundary.

## Close Checklist

- [x] All families migrated and verified (see Completed Families table above)
- [x] Full test suite + vet passes clean
- [x] HAWP check passes (all 3 checks)
- [x] Plan moved to closed/2026/09/08/e7a5e294/
- [x] BACKLOG.md updated (removed from active, added to Recently Closed)
- [x] Parent plan 47c793d6 updated to note CLI phase complete
