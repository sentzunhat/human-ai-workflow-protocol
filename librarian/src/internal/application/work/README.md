# work

`hawp work validate` (backlog/plan/evidence integrity checks — see
`internal/domain/work` for the actual rule implementations) and
`hawp work new` (intake scaffolding).

## Exports

`Validate` / `Render` (validation) · `NewItem` (intake scaffolding) ·
`Normalize` (work-record drift detection/fixing).

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
