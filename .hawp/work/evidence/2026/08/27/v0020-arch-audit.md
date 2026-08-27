# Architecture Audit — v0.0.20

Date: 2026-08-27  
Scope: `librarian/src/internal/` — ports-and-adapters alignment vs. mictlan pattern

## Reference pattern (mictlan)

mictlan (`src/backend/`) organizes by **bounded context first**, then by layer:

```
areas/<context>/
  domain/       — types, interfaces (ports), pure business rules
  application/  — use-case orchestration, no infrastructure imports
  infrastructure/ — adapters: DB, HTTP clients, filesystem
main/           — shared/cross-cutting (auth, logging, config)
platform adapters (routes, webhooks) live outside areas/
```

## Current hawp structure

hawp organizes by **layer first**, then by feature:

```
internal/
  domain/{context,embeddings,index,kit,kitsync,llm,provision,search,update,usage,work}/
  application/{check,context,db,embed,index,kit,kitsync,links,provision,search,update,uuidgen,work}/
  infrastructure/{archive,download,filesystem,githubrelease,markdown,repo,selfreplace,sqlite}/
  platform/{cli,mcp}/
```

## Gap analysis

### 1. Organization: layer-first vs. area-first

**Severity: Low (cosmetic / navigational)**

| hawp (layer-first) | Equivalent mictlan (area-first) |
|--------------------|---------------------------------|
| `domain/search/` + `application/search/` + (uses `infrastructure/sqlite`) | `areas/search/domain/` + `areas/search/application/` + `areas/search/infrastructure/` |
| `domain/usage/` (no application/, uses sqlite directly) | `areas/usage/domain/` + `areas/usage/infrastructure/` |
| `domain/work/` + `application/work/` | `areas/work/domain/` + `areas/work/application/` + `areas/work/infrastructure/` |

A directory reorganization would improve cognitive grouping (all search code together rather than split across three top-level directories) but carries high refactor risk (all import paths change). Deferred — not blocking the v0.1.0 gate.

### 2. Missing port interfaces — HIGH PRIORITY

Interfaces define the boundary between domain and infrastructure. Without them, callers depend on concrete types and the storage backend cannot be swapped or stubbed.

| Package | Violation | Fix |
|---------|-----------|-----|
| `domain/usage` | `Store` struct was concrete; callers depended on `*Store` | **Fixed in v0.0.20**: `Store` interface extracted; `sqliteStore` is the concrete impl |
| `domain/kit` | File I/O called directly via `os.ReadFile`, no `FileReader` port | Future patch |
| `domain/work` | 8 files use `os.*` / `filepath.*` directly; no repository interface | Future patch |
| `domain/context` | `kit.go`, `work.go` read from filesystem directly | Future patch |
| `domain/embeddings` | `Embedder` interface exists ✅ | Already compliant |
| `domain/llm` | `LLMClient` interface exists ✅ | Already compliant |
| `domain/search` | `Result` type only; no repository interface | Future patch |

### 3. Constructor placement — MEDIUM PRIORITY

`domain/usage.Open()` creates a sqlite connection inside the domain package. In a strict ports-and-adapters model, the constructor belongs in `infrastructure/sqlite/` and is injected into the domain.

**Current state:** `domain/usage.Open()` imports `database/sql` and `modernc.org/sqlite`.  
**Ideal state:** `infrastructure/sqlite.NewUsageStore(path) usage.Store` — domain stays pure.  
**v0.0.20 fix:** `Open()` now returns `Store` (interface), isolating callers from the concrete type. Full constructor migration is a separate patch (requires updating 5 call sites in cli and mcp).

### 4. Application layer gaps — LOW PRIORITY

Several domain packages have no corresponding application layer:
- `domain/usage` — use cases (enable, disable, report) live directly in `platform/cli/run.go`
- `domain/embeddings` — embedding is called from multiple places with no application-level service
- `domain/update` — update logic split between domain and CLI

These aren't blocking but mean use-case logic accumulates in `run.go`, making it grow. Extracting to `application/usage/`, `application/embed/` would be a future patch.

## What changed in v0.0.20

### `domain/usage/store.go`

- Added `Store` interface (port):
  ```go
  type Store interface {
      Write(tool string, inputJSON, outputJSON []byte, logBodies bool) error
      Recent(n int) ([]Entry, error)
      GetTotals() (Totals, error)
      GetReport() (Report, error)
      Clear() error
      Close()
  }
  ```
- Renamed concrete struct `Store` → `sqliteStore` (unexported)
- `Open()` now returns `(Store, error)` — callers no longer depend on `*Store`
- `store_test.go`: `openTemp` return type updated from `*Store` → `Store`

### Callers unchanged

All 5 call sites (`platform/cli/run.go` × 4, `platform/mcp/tools.go` × 1) already used the returned value only through method calls, so no changes were needed.

## Recommended next steps

| Patch | Action | Risk |
|-------|--------|------|
| v0.0.21 | Move `usage.Open()` to `infrastructure/sqlite/NewUsageStore()` | Low — 5 call sites |
| v0.0.22 | Extract `FileReader` port for `domain/kit` and `domain/work` | Medium — many callers |
| Future | Directory reorganization to area-first | High — all import paths |

## Compliance after v0.0.20

| Criterion | Status |
|-----------|--------|
| `domain/usage` has a Store interface (port) | ✅ |
| Callers depend on interface, not concrete type | ✅ |
| `domain/embeddings` has Embedder interface | ✅ |
| `domain/llm` has LLMClient interface | ✅ |
| `domain/kit`, `domain/work` have repository interfaces | ❌ (future) |
| `usage.Open()` lives in infrastructure | ❌ (partial — returns interface) |
| Directory structure is area-first | ❌ (deferred) |
