# `hawp work` document-type subcommands and MCP tool

**UUID:** `86c3911c-37c1-4acf-b7c8-5b3d999f2dc2`
**Type:** feature
**Reported:** 2026-09-12
**Status:** `done`

---

## Input (verbatim)

> the work command can do easy type of document if possible and that can be a tool

## Intake Summary

The `hawp work` family already has `new`, `validate`, `normalize`, and
`reshape`. The user wants lightweight subcommands for creating secondary work
documents (status reports, evidence files, decision records, notes) at the
correct UUID-subfolder path automatically — and the same capability exposed as
an MCP tool so AI agents can create these documents without knowing the path
conventions.

## Current State (verified 2026-09-12)

`hawp work new` creates a full investigation scaffold (UUID, plan.md, backlog
entry). There is no command for creating a secondary document (status, evidence,
decision, note) that: (a) generates a UUID, (b) writes to the correct
`YYYY/MM/DD/{uuid}/` path, (c) scaffolds a minimal front-matter template.

`hawp_work_reshape` is the only work-domain MCP tool beyond the three search
tools. No document-creation MCP tool exists.

## Plan

### Slice A — `hawp work status` / `evidence` / `decision` / `note` subcommands

Each command:
1. Accepts a `--title` flag (short label, kebab-cased for filename)
2. Generates a UUID and resolves the date-stamped path:
   `{type}/YYYY/MM/DD/{uuid}/{filename}.md`
3. Writes a minimal template with front-matter (uuid, title, type, date)
4. Prints the path so the caller or agent can open and fill it

Commands: `hawp work status`, `hawp work evidence`, `hawp work decision`,
`hawp work note`

Implementation: add a `doc` subpackage under `platform/cli/work/`; share UUID
generation and path resolution via `application/work/`. No domain changes
needed — document creation is a pure filesystem use case.

- [x] Implement `application/work/doc/` use case: `CreateWorkDoc(type, title, root) (path, error)`
- [x] Implement `platform/cli/work/doc/command.go` dispatching on doc type
- [x] Wire into `cmd/hawp/main.go` work family
- [x] Tests: path shape, UUID uniqueness, template content, unknown type rejected

### Slice B — `hawp_work_doc` MCP tool

Add a 5th MCP tool `hawp_work_doc` accepting `{"type": "status|evidence|decision|note", "title": "..."}`.
Returns the created file path. Agents can then append content to it without
knowing path conventions.

- [x] Add `hawp_work_doc` to `platform/mcp/server/`
- [x] Unit tests with fake filesystem; no CLI dependency

## Risk

Low: isolated new subpackage, no changes to existing commands or domain model.
Slice B depends on Slice A's use case.

## Verification

- `hawp work status --title "v0.1.0-gate"` creates
  `status/YYYY/MM/DD/{uuid}/status.md` with correct front-matter
- `hawp work evidence --title "downstream-savings"` creates
  `evidence/YYYY/MM/DD/{uuid}/evidence.md`
- `hawp_work_doc` MCP tool returns the path; file exists on disk
- `go test ./...` passes; `go vet ./...` clean

## Outcome

`hawp work doc` CLI subcommands (status, evidence, decision, note) implemented
with `--title` and `--work-item` flags. `hawp_work_doc` MCP tool added as the
5th MCP tool. Both generate canonical UUID subfolder paths and return the created
file path. Unit tests and e2e MCP tests pass. `go vet` clean.

## Close Checklist

- [x] All focused checks pass (recorded in Verification section above).
- [x] HAWP check passes.
- [x] Plan status: done.
- [x] Plan moved to `.hawp/work/closed/2026/09/12/86c3911c/plan.md`.
- [x] BACKLOG.md updated.
