# work

`hawp work validate` (backlog/plan/evidence integrity checks — see
`internal/domain/work` for the actual rule implementations) and
`hawp work new` (intake scaffolding).

## Use-case APIs

`Validate` / `Render` (validation) · `NewItem` (intake scaffolding) ·
`Normalize` (work-record drift detection/fixing) · `DraftIntake` (internal draft service).

The source-layout migration preserves these APIs while grouping their owners:

| Flat source file | Nested owner under `internal/application/work/` |
| --- | --- |
| `intake.go`, `draft.go`, `draft_test.go` | `intake/` |
| `validate.go` | `validation/` |
| `normalize.go` | `normalize/` |

This README stays at the family root as an overview. Callers must import the
owning child after migration; no parent facade is introduced. The migration tool
updates those imports, including normalization's call into validation.
Domain/work rule and filesystem separation remains a distinct semantic task.

## CLI quick reference

Run the command surface from `librarian/src`:

```bash
go run ./cmd/hawp work validate
go run ./cmd/hawp work normalize --dry-run --validate
go run ./cmd/hawp work new "title" --type task
```

For another repo or an older HAWP checkout, point the validation or
normalization pass at that target tree:

```bash
go run ./cmd/hawp work validate --hawp-root /path/to/repo/.hawp
go run ./cmd/hawp work normalize --dry-run --hawp-root /path/to/repo/.hawp
go run ./cmd/hawp work normalize --dry-run --work-root /path/to/repo/.hawp/work
```

Current compatibility expectations:

- `work validate` accepts UUID rows, legacy `TASK-*` / `BUG-*` rows, and older numeric IDs
- backlog rows may come from `UUID`, `Legacy ID`, `ID`, or `#` columns
- plan links may be Markdown links or plain relative paths
- nested `###` subsections under `## Active Work` are valid

`work normalize` is intentionally stricter than `work validate`: it can still
surface manual-review drift on historical non-canonical rows even when the repo
is safe to validate.

## `NewItem` — what it does and doesn't do

Generates a UUID, writes an investigation plan file shaped like
`.hawp/kit/templates/work-intake.md`, and inserts a `status: inbox` row into
`BACKLOG.md`'s Active Work table.

It deliberately does **not** investigate the request or write the plan
itself — HAWP is "a shaping protocol, not a runtime" (`AGENTS.md`); actual
investigation requires reasoning about the specific request, which stays a
human/AI-agent job. This only removes the boilerplate of hand-typing the
backlog row and file skeleton.

## Quick use

```go
result, err := work.NewItem(workDir, "bug", "Fix the reshape flag", "the --llm-reshape flag is broken")
// result.PlanFilePath, result.BacklogPath
```

## Intake drafts (internal API)

`DraftIntake(ctx, request, shaper)` proposes a reviewable HAWP shape through an
injected `RequestShaper`. The service preserves nonblank caller-supplied input
and context byte-for-byte; absent/whitespace context becomes an explicit unknown
label. The shaper can propose only mission, constraints, output, and checkpoint.
The first three must be nonblank; checkpoint is optional.

The draft carries no UUID, owner, status, or approval. Validation establishes
field completeness, not truth or preservation of every implied requirement in
the generated fields. Provider errors and cancellation return no partial draft.
The use case does not access storage or call `NewItem`; an adapter's own side
effects and output fidelity must be reviewed separately.

There is no default shaper, model dependency, CLI command, or MCP reshape tool.
Current tests inject deterministic proposals. A real adapter requires separate
fidelity/error tests before client wiring. The existing work-new workflow stays
unchanged; creating a work record is a separate action after draft review.
