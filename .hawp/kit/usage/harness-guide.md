# Harness Guide — Applying the Slice Harness

This guide explains how to use the three standards documents when starting
a new implementation slice. Read it before writing any code or tests.

## Which document to read first

| Slice type | Read |
|------------|------|
| Any slice (always) | `standards/slice-harness.md` |
| New embedding or LLM backend | `standards/provider-harness.md` |
| New CLI command or MCP tool | `standards/tool-harness.md` |

Provider and tool slices read both their specific document and the universal
slice harness.

## Step-by-step

### 1. Open the relevant standards document

Before writing any code, read the applicable standards document. The
document tells you what to verify, what test cases are required, and what
the evidence record must contain.

### 2. Record UUID and owned paths before writing code

Open or create the plan file at `.hawp/work/active/<uuid>/plan.md`. Fill
in the header block from `slice-harness.md`:

```
Work UUID:   <uuid>
Plan file:   .hawp/work/active/<uuid>/plan.md
Branch:      <branch-name>
Owned paths: <list files you will create or modify>
Session:     YYYY-MM-DD
```

Do this before writing a single line of implementation. The UUID anchors
all status checkpoints and evidence under `status/YYYY/MM/DD/<uuid>/` and
`evidence/YYYY/MM/DD/<uuid>/`.

### 3. Run focused tests first

When you have code to test, run only the affected package:

```bash
cd librarian/src && go test ./internal/path/to/package/... -v
```

Do not run `go test ./...` until the focused package tests pass. Record
the actual output, not what you expect it to say.

### 4. Escalate checks in order

Once the focused tests pass, run in sequence:

```bash
cd librarian/src && go build ./...
cd librarian/src && go vet ./...
cd librarian/src && go run ./cmd/hawp check --no-update-check
```

Only run `go test ./...` if you modified Go source and after the above
three steps are clean.

### 5. Record actual test output as evidence

Paste the actual terminal output into the evidence record in the plan file.
Do not paste expected output or fabricate results. Separate verified facts
(things the output proves) from inferences (things that follow from reasoning
but are not directly shown).

### 6. Update the plan status before closing

When the slice is done:

1. Set the plan status to `done`.
2. Move the plan to `.hawp/work/closed/YYYY/MM/DD/<uuid>/plan.md`.
3. Update BACKLOG.md — move the item to Recently Closed, cap that section
   at 5–10 items.

## Quick reference

```
Standards:
  .hawp/kit/standards/slice-harness.md     — universal template
  .hawp/kit/standards/provider-harness.md  — embedding and LLM backends
  .hawp/kit/standards/tool-harness.md      — CLI commands and MCP tools

Artifact paths:
  .hawp/work/active/<uuid>/plan.md          — active plan
  .hawp/work/status/YYYY/MM/DD/<uuid>/      — session checkpoints
  .hawp/work/evidence/YYYY/MM/DD/<uuid>/    — evidence files
  .hawp/work/closed/YYYY/MM/DD/<uuid>/      — closed plan

Check commands (run from librarian/src):
  go test ./internal/path/to/package/... -v     # focused
  go build ./... && go vet ./...                 # build/vet
  go run ./cmd/hawp check --no-update-check      # HAWP validation
  go test ./...                                  # full (only if Go source changed)
```

## What not to do

- Do not run `go test ./...` as the first check — run the focused package first.
- Do not paste expected output into the evidence record — paste actual output.
- Do not invent co-author trailers or model identity labels.
- Do not create per-field folders (`input/`, `context/`, `mission/`,
  `constraints/`, `output/`, `checkpoint/`) — HAWP prohibits them.
- Do not modify files owned by another active work item without coordination.
