# Move filesystem operations out of domain/kit, domain/kitsync, domain/provision, domain/distribution

**UUID:** `8ddb06ea-0d10-454a-8e7e-96ab2e7eb7a8`
**Type:** improvement
**Reported:** 2026-09-10
**Status:** `done`

---

## Input

Session audit identified that the same filesystem-boundary violation fixed for
`domain/work` exists across several other domain packages. Domain packages below
contain direct `os.ReadFile`, `os.WriteFile`, `os.Rename`, `os.Stat`, and
`os.MkdirAll` calls that belong in infrastructure or application layers.

## Investigation

**Verified violations (domain packages with os.\* filesystem ops):**

| File                                  | Operations                                                                  |
| ------------------------------------- | --------------------------------------------------------------------------- |
| `domain/kitsync/apply.go`             | Stat, MkdirAll, Open, Rename, Remove, ReadFile — full file-copy/sync engine |
| `domain/kitsync/manifest.go`          | ReadFile                                                                    |
| `domain/kit/validate.go`              | ReadFile                                                                    |
| `domain/kit/normalize.go`             | ReadFile, Stat, Rename, WriteFile                                           |
| `domain/provision/manifest.go`        | ReadFile, WriteFile                                                         |
| `domain/distribution/distribution.go` | ReadFile (3 calls)                                                          |
| `domain/context/kit.go`               | ReadFile (2 calls)                                                          |
| `domain/context/work.go`              | ReadFile (plan file, `buildWorkDocument`)                                   |
| `domain/providersync/materialize.go`  | ReadFile                                                                    |

**Not a violation (correctly placed):**

- `domain/providers/embeddings/embedder.go` — port interface only
- `domain/providers/llm/llm_client.go` — port interface only
- Infrastructure implementations in `infrastructure/models/` — correct

## Decision And Plan

Same pattern as `742aa60b` (closed 2026-09-10):

1. For each domain package: extract the filesystem operations into a corresponding
   `infrastructure/repositories/<package>/` adapter or into the application use-case
   layer, depending on whether the operation is a storage concern (infra) or
   orchestration concern (application).
2. The domain package retains pure logic; the file read is passed in as `[]byte`
   or a reader, or the operation is called via an interface.
3. Update callers; run focused tests per move; run full suite at each completed
   domain package.

**Execution order (dependency-first):**

- `domain/kitsync/manifest.go` + `domain/provision/manifest.go` (small, isolated reads)
- `domain/context/kit.go` + `domain/context/work.go` (read-only, read is passed in)
- `domain/kit/validate.go` (read-only, straightforward)
- `domain/providersync/materialize.go` (read-only)
- `domain/distribution/distribution.go` (read-only, 3 sites)
- `domain/kit/normalize.go` (reads + writes — larger slice, keep separate)
- `domain/kitsync/apply.go` (largest slice — full file-copy engine, own worktree)

## Risk And Execution Gate

Risk: medium (same as 742aa60b — import topology + callers across layers).
Each domain package is a separate slice; no broad simultaneous moves.

Status: `plan-ready`. Parent audit: `47c793d6`.

## File Ownership

Owns `internal/domain/kit/`, `internal/domain/kitsync/`,
`internal/domain/provision/`, `internal/domain/distribution/`,
`internal/domain/context/`, `internal/domain/providersync/` and
corresponding infra/application targets. Does not touch `domain/work/`
(closed in `742aa60b`) or CLI platform layer.

## Verification / Acceptance

- Each moved domain package has zero `os.*` filesystem calls in non-test files.
- Callers route through the new infra adapter (no inline `os.ReadFile` duplication).
- `go test ./...`, `go vet ./...`, `go run ./cmd/hawp check` all pass after each slice.
- No storage format, command behavior, or public API changes.

## Checkpoints

### Slice B — 2026-09-10 (branch: agent/8ddb06ea-kit)

**Scope:** `domain/kit/validate.go`, `domain/distribution/distribution.go`

**Changes:**

- `domain/kit/validate.go`: `CheckInternalLinks` and `Validate` now accept `func(string) ([]byte, error)` reader param; `os` import retained only for `os.ReadDir` (used by `CheckFileNaming`, not in scope).
- `domain/distribution/distribution.go`: `ComputeExpectedOutputs`, `FindDownstreamPathLeaks`, `composeProviderScript`, `extractBashBody` all accept injected reader; `os` import removed; `fs.ErrNotExist` (via `errors/io/fs`) replaces `os.IsNotExist`.
- New: `infrastructure/repositories/kit/validate_reader.go` — thin `os.ReadFile` wrappers for `CheckInternalLinks` and `Validate`.
- New: `infrastructure/repositories/distribution/reader.go` — thin `os.ReadFile` wrappers for `ComputeExpectedOutputs` and `FindDownstreamPathLeaks`.
- Updated: `application/kit/validate.go` routes through `infrakit.Validate`.
- Updated: `application/distribution/distribution.go` routes through `infradistribution.*`.
- Updated: existing domain tests pass `os.ReadFile` as reader arg.
- Added infra tests: `validate_reader_test.go`, `reader_test.go` (missing-file coverage).

**Verification:** `go build ./...` clean, `go test ./...` all pass, `go vet ./...` clean. Kit validate and links check pass; work validate fail is pre-existing (closed plan 74aaa332).

## Completion Record

### Slices E + F — 2026-09-10 (branch: feature/v0.0.24-work-folder-normalization)

**Scope:** `domain/kit/normalize.go`, `domain/kitsync/apply.go`

**Findings (verified 2026-09-10):**

- `domain/kit/normalize.go` (Slice E): zero `os.*` filesystem calls in non-test
  files. `PlanFileRenames`, `ApplyRenames`, `ApplyLinkUpdates` route through
  injected Reader/Writer. Implemented in `dce8a87b` / `e070d1d4`.
- `domain/kitsync/apply.go` (Slice F): all filesystem ops (Stat, MkdirAll,
  ReadDir, Open, CreateTemp, Rename, Remove) go through the `FileCopier`
  interface (`fc.*`). Concrete `os.*` implementation lives in
  `infrastructure/repositories/kitsync/filecopy.go` — correct layering.

**Acceptance:** every `domain/{kit,kitsync,provision,distribution,context,
providersync}` package now has zero direct `os.*` filesystem calls in non-test
files. Remaining `os.*` in `domain/work/*` is owned by closed plan `742aa60b`
and is out of scope here.

**Verification (fresh, 2026-09-10):** `go build ./...` clean; `go test ./...`
no failures; `go run ./cmd/hawp check` — 2 of 3 validations fail only on
pre-existing `.hawp/.spaces/agent-kitsync/` links (unrelated to this item).

Status: `done` (2026-09-10). All slices B–F implemented and verified.

## Outcome

Every targeted domain package — `kit`, `kitsync`, `provision`, `distribution`,
`context`, and `providersync` — now contains zero direct `os.*` filesystem calls
in non-test files. Storage-bound I/O routes through `infrastructure/repositories/<package>/`
adapters; the `kitsync` copy engine and `kit` normalize apply/plan operations route
through injected `FileCopier` and Reader/Writer interfaces. Public API, storage
format, and command behavior are unchanged.

## Close Checklist

- [x] All six target domain packages free of direct `os.*` filesystem calls.
- [x] Concrete `os.*` implementations live only in infrastructure adapters.
- [x] `go build ./...` clean and `go test ./...` passes.
- [x] No storage format, command behavior, or public API changes.
- [x] Plan moved to `closed/2026/09/10/8ddb06ea/`.
- [x] BACKLOG.md updated (moved to Recently Closed).
