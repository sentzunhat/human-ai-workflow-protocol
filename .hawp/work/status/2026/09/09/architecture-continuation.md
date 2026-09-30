### Architecture Continuation Checkpoint

Historical checkpoint. The work-only helper described below is superseded by the
whole-source tool in `scripts/source-layout/`; current commands and limitations
are in [its guide](../../../../../../scripts/source-layout/README.md).

#### Completed

- CLI decomposition is closed under `e7a5e294`.
- Runtime adapter work is closed under `b61770a6`: model adapters, bootstrap
  composition, usage repository/config adapters, and injected search/context/index
  factories are verified.
- Live integration and benchmark tests are explicit opt-in build-tag tests.

#### Remaining

- `742aa60b` remains the next architecture item: separate work-domain rules,
  application orchestration, and filesystem adapters.
- The first slice should isolate backlog parsing and its pure values without
  changing status, UUID, path, or normalization behavior.

#### Reviewed Ownership Map

| Current source | Intended destination | Boundary | Automation |
| --- | --- | --- | --- |
| `internal/domain/work/backlog.go` | `internal/domain/work/backlog/parser.go` | parser needs parent values/type extraction | preview only |
| `internal/domain/work/normalize_*` | `domain/work/normalize/` and `infrastructure/workfs/` | shared rule types and writes | preview only |
| `internal/application/work/intake.go` | `application/work/intake/` | UUID/persistence orchestration | preview only |
| `internal/application/work/validate.go` | `application/work/validation/` | report assembly and ports | preview only |
| `internal/application/work/normalize.go` | `application/work/normalize/` | orchestration/workfs coupling | preview only |

#### Safe Automation

`scripts/migrate-work-boundaries.sh` performs repository-root and source/destination
precondition checks and prints the explicit mapping. Default preview changes no
files. Apply mode is intentionally refused until a semantic extraction slice has
its own implementation and tests.

#### Verification

`go test ./...`, `go vet ./...`, and `git diff --check` passed for the completed
runtime architecture changes. The migration preview passed; its guarded apply mode
returned the expected refusal without changing files.
