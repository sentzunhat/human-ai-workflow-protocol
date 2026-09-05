# Checkpoint 2026-09-06 — Architecture Audit & v0.0.24 Progress Gate

## Context Before This Session

The HAWP project is in the middle of `v0.0.24` release work on branch `feature/v0.0.24-work-folder-normalization`. Three active work items tracked:

| ID       | Title                                                      | Status      |
| -------- | ---------------------------------------------------------- | ----------- |
| e5fca9c7 | v0.0.24 work-folder normalization and README positioning   | in-progress |
| 47c793d6 | CLI decomposition and architecture audit continuation      | in-progress |
| 0cb0f9b0 | Single-executable installation and configuration migration | in-progress |

The working tree was clean — no uncommitted changes. An existing status report existed for `0cb0f9b0` (migration progress).

## What Changed During This Session

### Architecture Audit Findings (from `47c793d6`)

A critical type duplication issue was identified between two packages:

**Problem:** `domain/index.Chunk` and `sqlite.Chunk` are nearly duplicate definitions with incompatible field differences:

| Type             | `domain/index`                         | `infrastructure/sqlite` |
| ---------------- | -------------------------------------- | ----------------------- |
| Chunk.Text       | string                                 | string                  |
| FolderContext    | `string`                               | `*string` (pointer)     |
| LineStart        | ❌ not present                         | ✅ int (extra field)    |
| LineEnd          | ❌ not present                         | ✅ int (extra field)    |
| DocumentMetadata | ✅ defined with WorkUUID, Status, etc. | ✅ identical duplicate  |

**Scope of impact:** 5 files construct `sqlite.Chunk{}` directly, and `ingest-service.go` imports BOTH packages simultaneously, creating maintenance debt and potential inconsistency.

### Existing Architecture Audit Plan Context (`47c793d6/plan.md`)

The plan explicitly calls for this fix:

> "Continuation: move the search index contract and embedding metadata into `internal/domain/search`, preserve SQLite compatibility via a type alias, and verify injected repositories, error propagation, and resource closure."

### Previous Work Already Completed (Context from Audit Plan)

Per `47c793d6`'s verification section, the following was already committed and verified:

- ✅ Split command owners in CLI package
- ✅ Reject non-finite hybrid ratios with regression tests
- ✅ Strict `flag.NewFlagSet` typed parsers for all command handlers (work_new, index_build, kit_validate, links_clean, etc.)
- ✅ Schema-aware backlog intake supporting UUID-only and legacy-ID headers
- ✅ Corpus file-read failure propagation
- ✅ Six standard target builds (`make dist VERSION=0.0.24`) pass
- ✅ `go test ./...`, `go vet ./...` clean
- ✅ HAWP validation: kit validate PASS, work validate PASS (0 issues), links check PASS (127+ Markdown files)

## What Remains Unproven / Next Steps

### Immediate Blocker: Type Consolidation

Before proceeding with further CLI work or provider integration, the `Chunk`/`DocumentMetadata` duplication must be resolved. Options:

**Option A (preferred):** Promote `LineStart`/`LineEnd` fields into `domain/index.Chunk`, then replace all `sqlite.Chunk{}` usages with `domainindex.Chunk{}`, add type aliases in `sqlite` package for backward compat.

**Option B:** Keep `sqlite.Chunk` but have it embed/alias `domainindex.Chunk` — requires field alignment first.

**Files requiring updates after consolidation:**

1. `librarian/src/internal/application/index/embed_integration_test.go` (line ~65)
2. `librarian/src/internal/application/index/embed-service_test.go` (line ~44)
3. `librarian/src/internal/application/index/ingest-service.go` (lines ~97, ~149)
4. `librarian/src/internal/platform/mcp/tools_e2e_test.go` (line ~51)
5. `librarian/src/tests/application/search/service_test.go` (line ~37)

### Pending CLI Work from Audit Plan

Per `47c793d6`'s "Next compoundable slice":

- [ ] Pure standard-library search argument parsing with typed options and invalid-input tests
- [ ] Schema-aware backlog intake supporting UUID-only and legacy-ID headers (partially done, needs verification)

### Pending Single-Executable Migration (`0cb0f9b0`)

- [ ] Full single-executable upgrade compatibility testing
- [ ] Codex table customization preservation
- [ ] Live provider reload and cross-platform execution
- [ ] Unknown-folder routing (planned, not implemented)
- [ ] UUID enrichment for work items

### Pending Work-Folder Normalization (`e5fca9c7`)

Status unclear — plan file exists but no detailed status report.

## Strategic Impact

1. **Type consolidation is a prerequisite** for cleaner CLI work and provider integration. The current duplication makes it impossible to reason about "the" Chunk type — callers must know which package they're referencing, leading to the awkward dual-import in `ingest-service.go`.

2. **Migration gate remains open.** The single-executable migration (`0cb0f9b0`) cannot be completed until configuration safety is proven and wrapper retirement is planned. This is a v0.1.0 gate item — not blocking v0.0.24, but should be tracked separately.

3. **Architecture debt vs feature work tension.** The audit identified 2+ major remaining slices (type consolidation, search argument parsing) that are architectural improvements rather than new features. These should be treated as enabling work for v0.1.0, not rushed into v0.0.24.

## Constraints

- Do NOT merge or publish to any branch other than `feature/v0.0.24-work-folder-normalization`
- No PRs, releases, or downstream repository edits during this session
- Preserve all existing work copies and sidecars
- HAWP validation must pass (`go run ./cmd/hawp check --no-update-check`) before any commit

## Recommendations

1. **Immediate:** Implement Option A for Chunk/DocumentMetadata consolidation as a focused, verified change. This is the highest-leverage compounding action — it unblocks all future CLI/provider work.

2. **Next iteration:** After consolidation, implement pure standard-library search argument parsing per `47c793d6`'s plan.

3. **Deferral:** Move single-executable migration (`0cb0f9b0`) to v0.1.0 scope — it requires provider coordination and cross-platform testing that is out of scope for a patch release.

4. **Status report update:** Create/update status reports for all three active work items before closing this checkpoint.

---

_Checkpoint created: 2026-09-06_
_Branch: feature/v0.0.24-work-folder-normalization_
_Working tree: CLEAN_
