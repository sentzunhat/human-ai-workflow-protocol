# work

`hawp work validate` (backlog/plan/evidence integrity checks — see
`internal/domain/work` for the actual rule implementations) and
`hawp work new` (intake scaffolding).

## Exports

`Validate` / `Render` (validation) · `NewItem` (intake scaffolding) ·
`Normalize` (work-record drift detection/fixing).

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
