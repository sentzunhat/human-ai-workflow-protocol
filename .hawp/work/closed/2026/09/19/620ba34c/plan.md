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
- Moved normalization backlog/plan scanning into `internal/domain/work/normalization`; parent aliases preserve the existing application API.
- Moved normalization rule evaluation into `internal/domain/work/normalization`; a narrow function-field `Source` seam preserves the existing parent `WorkSource` callers.
- Moved completed-active-row cleanup into `internal/domain/work/normalization`; the parent method now supplies only its repository-relative path adapter.
- Routed closed-record normalization through `internal/domain/work/normalization`; filesystem path reconciliation remains in the parent mutation boundary.
- Moved migration link/reference rewriting into `internal/domain/work/normalization`; directory moves and file I/O remain in the parent migration engine.
- Moved canonical folder-ID derivation into `internal/domain/work/normalization`; the parent supplies legacy backlog and filename IDs through a compatibility seam.
- Moved the migration `MovedPlan` record into `internal/domain/work/normalization`; filesystem execution still owns collection and application of those records.
- Moved pure normalization operation value types into `internal/domain/work/normalization`.
- Moved normalization report models and rendering into `internal/domain/work/normalization`; mutation code retains parent compatibility through aliases while `WorkSource` coupling remains explicit.
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

**Status now:** done
**Plan file:** work/closed/2026/09/19/620ba34c/plan.md

## Next Step

- [x] Investigation recorded above (required before planning)
- [x] Write or update the plan file
- [x] Move backlog status to analyzing
- [x] Build the domain/work symbol dependency inventory
- [x] Add the reviewed mapping revision for the first bounded extraction
- [x] Generate a fresh preview and compile/vet the candidate
- [x] Review preview diff and retain the candidate without mechanical file moves
- [x] Extract the next reviewed shared boundary only after the identity slice remains stable
- [x] Re-audit normalization parser/link helpers before selecting that boundary
- [x] Extract the shared table-parser boundary used by backlog, intake, and normalization scan
- [x] Extract migration orchestration and result aggregation behind an explicit filesystem seam
- [x] Verify parent `WorkSource` migration compatibility and full repository checks
- [x] Reassess validation boundary; retain it in the parent until shared models/helpers are separated
- [x] Extract reusable work model values behind parent aliases
- [x] Extract validation policies behind a reusable source seam and parent adapters
- [x] Extract backlog parsing and intake/draft capabilities behind parent adapters

## Boundary Decision (2026-09-19)

`internal/domain/work/normalization` now owns migration orchestration,
sidecar movement, moved-plan aggregation, and backlog-link reconciliation
through `MigrationIO`. The parent package retains only preview-copy behavior
and compatibility adapters. Validation remains in `domain/work` because its
report/model types and private helper graph do not yet form a cycle-free
subpackage boundary. The follow-up batch established `model`, `validation`,
`backlog`, and `intake` packages. The parent package now primarily contains
compatibility aliases, normalization adapters, and the remaining shared
source boundary.

## Strict Filesystem Boundary Follow-up (2026-09-19)

The package split is structurally better but `domain/work` is not yet
filesystem-free. Direct production imports of `os` remain in:

- `domain/work/validation`: closed-record, backlog, and dead-link scans
- `domain/work/normalization`: plan scans and completed-row cleanup
- `domain/work/normalize_apply.go`: closed-record path reconciliation
- `domain/work/normalize_migrate.go`: preview copy and concrete OS adapter

The `normalization/migration_apply.go` seam is directionally correct, but its
concrete `os` adapter still lives in the parent domain package. The next safe
batch is to move concrete file walking/copy/mutation into application or
infrastructure adapters and leave domain packages with injected readers,
writers, listers, and path policies only.

Adjacent-domain audit found the same pattern in `domain/context`, `domain/kit`,
`domain/kitsync`, `domain/providersync`, `domain/provision`, and
`domain/distribution`; `domain/kitsync` and `domain/kit` are the next likely
high-value candidates after work, not an incidental part of this slice.

## Completion Evidence (2026-09-19)

## Verification

- `go test ./...` passes from `librarian/src`.
- `go vet ./...` passes from `librarian/src`.
- `git diff --check` passes.
- `go run ./cmd/hawp check --no-update-check` passes after closure reconciliation.

## Outcome

Completed the cohesive subpackage split and follow-up filesystem-boundary deepening. Work, kit, and kitsync now express filesystem dependencies through application/infrastructure seams while preserving the required domain APIs.

## Close Checklist

- [x] Identity, table, markdown, model, validation, backlog, intake, and normalization boundaries extracted.
- [x] Remaining normalization and validation filesystem adapters moved behind injected sources.
- [x] Kit and kitsync follow-up boundaries audited and improved.
- [x] Repository tests, vet, diff checks, and HAWP checks verified.
- [x] Plan moved to `closed/2026/09/19/620ba34c/`.
- [x] `BACKLOG.md` updated.

The strict filesystem follow-up completed the remaining capability boundaries:
`domain/work` normalization and validation consume injected sources,
`domain/kit` naming validation consumes an injected directory reader, and
`domain/kitsync` provider detection and file-copy contracts no longer depend on
concrete `os` types. The domain packages now contain no production `os`
imports in the audited work, kit, or kitsync scopes.
