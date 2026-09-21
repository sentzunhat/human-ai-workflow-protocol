# Superseded architecture note

Preserved on 2026-09-07 for provenance. Its blocker, type, and deferral claims are not current guidance.

# Project Architecture State — HAWP CLI / v0.0.24

**As of:** 2026-09-06

## Active Branch
`feature/v0.0.24-work-folder-normalization`

## Known Type Duplication Issue (CRITICAL - Blocks Future Work)

### Problem: `domain/index.Chunk` vs `sqlite.Chunk` are incompatible duplicates

| Field         | `domain/index.Chunk`              | `infrastructure/sqlite.Chunk`     |
|---------------|-----------------------------------|-----------------------------------|
| Text          | `string`                          | `string`                          |
| FolderContext | `string`                          | `*string` (pointer)               |
| LineStart     | ❌ NOT present                   | ✅ `int` (extra field)            |
| LineEnd       | ❌ NOT present                   | ✅ `int` (extra field)            |

### Problem: `domain/index.DocumentMetadata` vs `sqlite.DocumentMetadata` are identical duplicates
Both packages define same fields: DocumentID, WorkUUID, Status, Owner, RiskLevel, ReportedAt, ClosedAt.

### Impact: 5 files affected
1. `librarian/src/internal/application/index/embed_integration_test.go` (line ~65)
2. `librarian/src/internal/application/index/embed-service_test.go` (line ~44)
3. `librarian/src/internal/application/index/ingest-service.go` (lines ~97, ~149) - imports BOTH packages simultaneously
4. `librarian/src/internal/platform/mcp/tools_e2e_test.go` (line ~51)
5. `librarian/src/tests/application/search/service_test.go` (line ~37)

### Resolution Path (Option A - Preferred)
1. Add `LineStart`, `LineEnd` fields to `domain/index.Chunk` in `chunk.go`
2. Change `FolderContext` from `*string` to `string` (or keep pointer for compat if needed)
3. Replace all `sqlite.Chunk{}` usages with `domainindex.Chunk{}` across the 5 callers
4. Add type alias: `type Chunk = domainindex.Chunk` in sqlite package for backward compat

### Resolution Path (Option B - Alternative)
Keep sqlite types but have them embed/alias domain/index types — requires field alignment first.

### Plan Reference
This consolidation is explicitly called out in `47c793d6/plan.md`:
> "Continuation: move the search index contract and embedding metadata into `internal/domain/search`, preserve SQLite compatibility via a type alias"

## Verification Commands (Must Pass Before Commit)
```bash
cd librarian/src && go build ./...
go test ./...
go vet ./...
go run ./cmd/hawp check --no-update-check  # kit, work, links validations
git diff --check  # no whitespace errors
```

## Deferred Work (Not v0.0.24 Scope)
- Single-executable migration (`0cb0f9b0`) → Move to v0.1.0 scope
  - Requires provider coordination and cross-platform testing
  - Configuration safety proven, but wrapper retirement needs integration work
- Unknown-folder routing, UUID enrichment for work items → planned, not implemented

## Completed Prior Work (from `47c793d6` verification section)
- ✅ Split command owners in CLI package
- ✅ Strict typed parsers for all handlers (work_new, index_build, kit_validate, etc.)
- ✅ Schema-aware backlog intake (partially done, needs verification)
- ✅ Six standard target builds pass
- ✅ HAWP validation clean
