# Separate work domain rules, use cases, and filesystem adapters

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `742aa60b-5ed6-439b-a8d7-5e9a9212adc1`
**Type:** improvement
**Reported:** 2026-09-07

---

## Input (verbatim)

> Decide package boundaries across layers; plan folder-within-folder organization that stays simple, repeatable, and avoids duplication. Planning only now; pace usage and keep implementation slices small.

## Investigation

Confirmed: domain/work mixes pure IDs/templates with os reads/writes and
infrastructure/repo imports. Application/work mixes four use cases in one package.
The new DraftIntake contract is separate from NewItem persistence and must stay so.

This is a distinct implementation subitem of
[47c793d6](../47c793d6/plan.md), following the
[shared package map](../47c793d6/package-boundaries.md).

## Decision And Plan

1. Extract pure backlog parsing from path-based reading and retain shared work
   values/ID helpers in domain/work. Children may import these shared values;
   domain/work must never import its children.
2. Group domain rules into planned domain/work/backlog, intake, validation,
   normalize. Move facts in, decisions out; preserve all current lifecycle rules.
3. Group application orchestration into work/intake (create + draft), validation,
   normalize; place required filesystem contracts with the consuming use cases.
4. Move concrete reads/writes/moves to infrastructure/repositories/work. Update CLI/MCP
   callers. Any temporary application facade must be one-way and removed once
   callers migrate. Do not change storage format, UUID rules, or normalization
   outcomes during this topology pass; metadata fixes remain a following slice.

Alternative: retain flat packages with renamed files. Rejected for these scoped
owners because it leaves navigation/ownership concerns unresolved. Do not nest
already-cohesive small packages or duplicate rules between layers.

## File Ownership And Coordination

Own `librarian/src/internal/domain/work/`, application/work/ and planned
`librarian/src/internal/infrastructure/repositories/work/` with adjacent tests.
CLI workcmd and MCP caller imports are coordinated follow-ups, not parallel edits.
Read other application/domain packages without reorganizing them.
All shorthand directories above are under `librarian/src/internal` unless stated.
Before implementation, enumerate exact source/test/destination paths for the
selected slice. No broad simultaneous moves; preserve the current dirty checkpoint.

## Risk And Execution Gate

Risk: medium (Go package/import topology). Status: in-progress. Owner: Codex.
Planning is complete; implementation remains slice-gated. The reviewed migration
map is available through `scripts/source-layout/run.sh`; the previous generated
snapshot is `scripts/source-layout/review/preview.md`. The old work-only script
has been removed. Regenerate only after the checkpoint commit. The mapping proposes mechanical application/work package splits;
semantic domain/work extraction is explicitly retained, not auto-moved.

## Verification / Acceptance

Pure rule packages import no os/SQL/HTTP/infrastructure. Temporary-fixture
normalization preservation/idempotency tests and intake/draft contract tests pass.
No lost evidence/history, changed UUID/status behavior, or import cycles.
Run focused tests per move and full Go suite/vet/HAWP at a completed slice.
Planning-only: no implementation tests or model runs performed here. The final
HAWP validation verifies record integrity, not package migration correctness.

## Usage Budget

One cohesive slice per turn. Read only owned files, direct dependencies and tests.
Short updates; update this plan rather than creating a report for each move.
Focused tests during edits, full checks once at slice completion. No recurring
usage monitor, reset redemption, model switch, or paid action requested.

## 2026-09-10 pure backlog parser checkpoint

Implemented the first semantic slice after the mechanical layout migration:
`ParseBacklogMarkdown(content string)` now owns the pure table/section parsing
logic, while `ParseBacklog(path)` remains a compatibility file-reading wrapper.
The parser behavior and existing path-based error contract are unchanged.

Direct verification:

- `go test ./internal/domain/work ./internal/application/work/...` — pass
- `go test ./...` — pass
- `go vet ./...` — pass
- `go run ./cmd/hawp check --no-update-check` — all 3 validations pass, 0 issues/warnings

Attribution: this checkpoint records the human-directed architecture decision,
Codex implementation, and model-assisted review/verification as separate
evidence. No contributor identity or co-author trailer is invented without a
real name/email supplied by the contributor.

## 2026-09-10 filesystem boundary checkpoint

Moved `os.ReadFile` out of `domain/work/backlog.go`. `ParseBacklog(path)` is
removed; `ParseBacklogMarkdown(content)` is the only export. Created
`infrastructure/repositories/work.ReadBacklog` as the canonical file-reading
entry point. Updated `application/work/validation.Validate` and the intake
test to call it. Inlined the read in `normalize_duplicates.go` and
`domain/context/work.go` (both already owned other os operations).

Direct verification:
- `go test ./internal/domain/work/... ./internal/infrastructure/repositories/work/... ./internal/application/work/... ./tests/application/work/...` — pass
- `go test ./...` — pass
- `go vet ./...` — pass
- `go run ./cmd/hawp check --no-update-check` — all 3 validations pass, 0 issues/warnings

## 2026-09-10 backlog reader routing checkpoint

Routed the backlog file read in `domain/context/work.go` through
`infrastructure/repositories/work.ReadBacklog`. The inline `os.ReadFile +
ParseBacklogMarkdown` pair at `EnrichWork` is replaced by a single
`reposwork.ReadBacklog(...)` call. The `os` import is retained in that file
because `buildWorkDocument` still uses `os.ReadFile` for plan file content.

`normalize_duplicates.go` cannot route through `ReadBacklog` in this slice:
it lives in `domain/work`, and `infrastructure/repositories/work` imports
`domain/work`, making the dependency cycle `domain/work →
infra/repositories/work → domain/work` illegal in Go. The inline read there
remains. Eliminating it requires relocating `normalize_duplicates.go` (and its
exported functions `ApplyDuplicateLinks` / `PreviewDuplicateLinks`) to the
application or infrastructure layer — that is the next slice.

Added `infrastructure/repositories/work/backlog_reader_test.go` with a
missing-file error case for `ReadBacklog`.

Direct verification:
- `go test ./internal/domain/work/... ./internal/infrastructure/repositories/work/... ./internal/application/work/... ./tests/application/work/... ./internal/domain/context/...` — pass
- `go test ./...` — pass
- `go vet ./...` — pass (0 issues)
- `go run ./cmd/hawp check --no-update-check` — all 3 validations pass, 0 issues/warnings

## 2026-09-10 linkDuplicatePlans relocation checkpoint

Moved `ApplyDuplicateLinks`, `PreviewDuplicateLinks`, `linkDuplicatePlans`, and
`addRelatedRecordLink` out of `domain/work/normalize_duplicates.go` and into
`application/work/normalize/duplicate_links.go`. Deleted
`domain/work/normalize_duplicates.go` (now empty/gone). The backlog read is now
routed through `reposwork.ReadBacklog` instead of the inline
`os.ReadFile + ParseBacklogMarkdown`. Private helpers `idSet`, `matchesAnyID`,
`ensureBlankLine` are inlined in the new file (identical logic; originals remain
in domain/work for other callers there). Callers in `normalize.go` updated to
call the now-local `ApplyDuplicateLinks` / `PreviewDuplicateLinks` without the
`domainwork.` prefix. Tests moved to `duplicate_links_test.go` in the same
application package; `TestApplyDuplicateLinksPreservesUnreferencedWorkingCopy`
removed from `domain/work/normalize_test.go`.

Direct verification:
- `go build ./...` — clean
- `go test ./internal/domain/work/... ./internal/application/work/... ./internal/infrastructure/repositories/work/...` — pass
- `go test ./...` — all packages pass
- `go vet ./...` — 0 issues
- `go run ./cmd/hawp check --no-update-check` — all 3 validations pass, 0 issues/warnings

Status: this closes the circular-import blocker identified in the previous
checkpoint. `domain/work` no longer contains any function that could create a
cycle with `infrastructure/repositories/work`. This item's core boundary work is
complete; any remaining slices (provider composition, persisted work metadata)
are tracked in the 47c793d6 continuation.

## Outcome

Work domain layer boundary established. `domain/work` now contains only pure
parsing and value logic with no filesystem imports. `ParseBacklog(path)` removed;
`ParseBacklogMarkdown(content)` is the sole export. File reading is routed
through `infrastructure/repositories/work.ReadBacklog`. Duplicate-link logic
moved to `application/work/normalize/duplicate_links.go`, resolving the last
circular-import risk. All Go tests pass; `hawp check` passes.

## Close Checklist

- [x] All focused checks pass (recorded in checkpoints above).
- [x] HAWP check passes (`go run ./cmd/hawp check --no-update-check` — all 3 validations pass).
- [x] Plan status: done.
- [x] Plan moved to `.hawp/work/closed/2026/09/10/742aa60b/plan.md`.
- [x] BACKLOG.md updated.
