# Docs: hawp_work_doc tool description has wrong path pattern

**UUID:** `cop-workdoc-description`
**Type:** fix
**Severity:** Low
**Reported:** 2026-09-13
**Status:** `inbox`
**Source:** Copilot PR #40 review — `librarian/src/internal/platform/mcp/server/tools.go`

---

## Finding

The tool description advertises a uniform `{type}/YYYY/MM/DD/{id}/{type}.md`
path, but `CreateWorkDoc` maps `decision` → `decisions/` and `note` → `notes/`
(pluralised folder names). Callers following the description will look in the
wrong directory for those two types.

Additionally, the `CHANGELOG.md` entry documents `hawp work doc status` as
the command shape, but the CLI dispatcher routes `hawp work status` directly —
the `doc` subcommand does not exist in that form.

## Fix Plan

1. Update the `hawp_work_doc` description in `tools.go` to list the actual
   per-type directories:
   - `status` → `status/YYYY/MM/DD/{id}/status.md`
   - `evidence` → `evidence/YYYY/MM/DD/{id}/evidence.md`
   - `decision` → `decisions/YYYY/MM/DD/{id}/decision.md`
   - `note` → `notes/YYYY/MM/DD/{id}/note.md`
2. Correct the CHANGELOG v0.0.24 command example to match the dispatcher route
   (`hawp work status` / `hawp work evidence` / etc., not `hawp work doc status`).

## Files

- `librarian/src/internal/platform/mcp/server/tools.go`
- `librarian/src/CHANGELOG.md`

## Verification

- `hawp_work_doc` tool description in `tools.go` lists actual per-type directories: `status/`, `evidence/`, `decisions/`, `notes/`.
- CHANGELOG work doc command examples match the current CLI command signatures.
- `go test ./...` passes; `go vet ./...` clean.

## Outcome

Updated `hawp_work_doc` description in `tools.go` to show the actual directory mapping (`status|evidence/...` vs `decisions|notes/...`). Corrected CHANGELOG command examples to match current CLI. `go test ./...` and `go vet ./...` pass.

## Close Checklist

- [x] `hawp_work_doc` description reflects actual `decisions/`/`notes/` directories.
- [x] CHANGELOG examples match current CLI command signatures.
- [x] `go test ./...` passes; `go vet ./...` clean.
- [x] Plan status: done.
- [x] Plan moved to `.hawp/work/closed/2026/09/14/cop-workdoc-description/plan.md`.
- [x] BACKLOG.md updated.
