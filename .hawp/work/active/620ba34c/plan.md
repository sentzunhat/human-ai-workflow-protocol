# Split domain work into cohesive subpackages

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `620ba34c-5be8-4cc8-ac9a-1484188aa10c`
**Type:** infrastructure
**Reported:** 2026-09-19

---

## Input (verbatim)

> Audit internal/domain/work and prepare a reviewed source-layout mapping for cohesive subpackages such as backlog, intake, validation, links, and normalization. Map private symbol dependencies first; do not mechanically move files or split packages until the preview compiles and the plan is reviewed.

## Intake Summary

The existing source-layout tool has already applied its reviewed shallow
mapping and currently reports 287 files, 0 moves, and 0 content updates. The
next useful slice is a deeper package-boundary proposal for
`librarian/src/internal/domain/work`, which currently contains 25 Go files.

## Current Context

`domain/work` owns backlog parsing, work-item intake values, validation
reports, Markdown link checks, evidence/clarity checks, and normalization.
Application callers import the package directly from context, indexing,
application work services, MCP handlers, and repository adapters.

The existing migration engine supports explicit deeper destinations and
syntax-aware import/reference rewrites, but refuses unresolved private
cross-package references. Its current mapping intentionally leaves
`domain/work` cohesive.

## Dependency Inventory (2026-09-19)

| Candidate owner | Files | Direct coupling observed |
| --- | --- | --- |
| Backlog/intake | `backlog.go`, `intake.go`, `intake_table.go`, related tests | `intake_table.go` and `normalize_scan.go` both use private `parseTableCells`, `mappedCell`, and `stripCodeSpan` from `backlog.go`; this needs a shared parser package or exported/internal helpers before splitting. |
| Work identity/model | `types.go`, `draft.go`, `idparse.go`, `constants.go` | `idparse.go` is reused by backlog parsing, consistency, completeness, and normalization; it is a good low-risk foundational package if its public API remains stable. |
| Validation/reporting | `consistency.go`, `completeness.go`, `evidence.go`, `clarity.go`, `deadlinks.go`, `links.go` | Validation code shares `Backlog`, `Report`, ID helpers, and private Markdown link helpers. `deadlinks.go` and normalization migration both depend on `links.go` private functions. |
| Normalization scan/rules | `normalize_scan.go`, `normalize_rules.go`, `normalize_report.go` | Scan/rule/report types are tightly connected; scan also depends on backlog parser helpers and ID parsing. Keep this as one first extraction boundary. |
| Normalization mutation/migration | `normalize_apply.go`, `normalize_active_rows.go`, `normalize_migrate.go` | Mutating paths share `ApplyResult`, `WorkSource`, link rewriting, path safety, and scan/rule types. This should be a later boundary after scan/rule extraction. |

Direct import callers are outside the package in context, index, application
work services, MCP handlers, and repositories. The inventory therefore favors
extracting stable value/helper packages first and postponing policy package
splits that would require exporting many private functions.

## Initial Analysis

**Directly verified:**

- Added `internal/domain/work/identity` as the capability-local owner for ID parsing and matching rules.
- Preserved the `internal/domain/work` API with compatibility wrappers in `idparse.go`.
- Updated consistency and normalization callers to use the identity boundary for UUID and numeric-ID checks.
- Added boundary-local identity tests while retaining the parent package API tests.
- Moved the identity parsing test file into `internal/domain/work/identity` and added a focused compatibility-wrapper test in the parent package.
- Extracted shared Markdown table primitives into `internal/domain/work/table`; backlog parsing, intake-table insertion, and normalization scanning now depend on that boundary instead of private helpers in `backlog.go`.
- Moved work-domain Markdown link scanning into `internal/domain/work/markdown`; dead-link validation and normalization migration now use the explicit link boundary.
- Moved normalization backlog/plan scanning into `internal/domain/work/normalization`; parent aliases preserve the existing application API while rules remain in `domain/work`.
- The scan boundary compiles independently; normalization rules/report and mutation/migration remain coupled to `WorkSource` and `ApplyResult` and are intentionally not moved in this slice.
- Moved pure normalization operation value types into `internal/domain/work/normalization`; rule evaluation retains the parent package until its `WorkSource` dependency is separated.
- `go test ./...`, `go vet ./...`, `git diff --check`, and source-layout candidate checks passed.
- Source-layout preview reports 288 files, 0 moves, and 0 content updates.
- `scripts/source-layout/run.sh --preview --check --diff` passed against the
  current mapping with 287 retained paths and no pending moves.
- The largest current package is `internal/domain/work` with 25 Go files.
- `internal/application/context`, `internal/platform/mcp/server`, and
  `internal/infrastructure/repositories/index` have already been given
  explicit ownership decisions and should not be split incidentally.
- `go test ./...`, `go vet ./...`, provider validation, distribution
  validation, `hawp check`, and `git diff --check` passed before this intake.
- Private helper reuse confirms that a filename-based one-folder-per-item
  mapping would create unresolved cross-package references.

**Inferred (not yet proven):**

- The identity boundary is low risk because existing callers remain on `domain/work`; direct consumers of the new package are limited to its parent package.
- A mechanical file-to-folder split would create private cross-package
  references and could turn the domain package into a set of coupled adapters.
- Normalization is likely a separate bounded subdomain, but its backlog,
  plan-scan, link-rewrite, and migration helpers need a dependency graph before
  any package boundary is chosen.
- A first safe extraction candidate is an identity/value package around
  `idparse.go` and selected shared types; a full `domain/work` split is not yet
  justified by folder size alone.

**Likely scope:**

- Keep backlog, validation, links, and normalization together until their private parser and rewrite helpers have an explicit shared boundary.
- Produce a symbol/dependency inventory for `domain/work`.
- Propose explicit destinations for only cohesive groups: backlog/intake,
  validation/evidence/clarity, links, and normalization.
- Add mapping tests and a fresh preview plan; do not apply the move until the
  candidate compiles, vets, and the reviewed plan is accepted.

## Risk + Review Gate

**Risk:** medium — package splitting changes import paths and can expose or
duplicate domain policy if private dependencies are handled poorly.
**Gate:** review first on medium/high; implementation must stop on unresolved
private references, destination collisions, or candidate test/vet failure.

## Backlog + Plan Link

**Status now:** in-progress
**Plan file:** work/active/620ba34c/plan.md

## Next Step

- [x] Investigation recorded above (required before planning)
- [x] Write or update the plan file
- [x] Move backlog status to analyzing
- [x] Build the domain/work symbol dependency inventory
- [x] Add the reviewed mapping revision for the first bounded extraction
- [x] Generate a fresh preview and compile/vet the candidate
- [x] Review preview diff and retain the candidate without mechanical file moves
- [ ] Extract the next reviewed shared boundary only after the identity slice remains stable
- [ ] Re-audit normalization parser/link helpers before selecting that boundary
- [x] Extract the shared table-parser boundary used by backlog, intake, and normalization scan
