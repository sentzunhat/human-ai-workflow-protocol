---
work-item: c1d2e3fb
type: fix
title: "Introduce typed index persistence capability contracts"
status: in-progress
created: 2026-08-10
updated: 2026-08-21
parent: b6c4e8a2
depends-on: c1d2e3fa
---

# Fix: Index Persistence Boundary

## Mission

Introduce typed, capability-local persistence contracts for documents, chunks,
and embeddings, moving SQLite construction toward explicit composition.

## Done When

- Application index services depend on typed contracts rather than concrete
  SQLite operations where the audit identifies a seam.
- Existing schema and ingest/embed behavior remain unchanged.
- Focused tests prove document, chunk, and embedding persistence separately.

## Audit Evidence

- The CLI bypasses the shared build service, losing BACKLOG and resolved work
  metadata while assigning all work files `closed`.
- Ingest and embed services construct and manipulate SQLite directly; document
  replacement is not atomic and embedding metadata-read failures are ignored.
- Two document contracts currently require a CLI-only conversion.

## Smallest Safe Slice

- Add `domain/index/store` contracts for atomic enriched-document replacement
  and embedding persistence.
- Implement those contracts in `infrastructure/sqlite/index`, keeping the
  existing SQLite search adapter separate as a read capability.
- Migrate the CLI to the shared build service before removing compatibility
  conversion code.

## Verification

- Prove atomic replacement, metadata-read error propagation, and adapter-backed
  application tests without SQLite construction in the services.
- Prove one CLI indexing flow retains resolved work metadata.
- Preserve schema, FTS triggers, flags, output, and mixed-model protections.

## Progress

- Extended the existing SQLite index owner with `ReplaceDocument`, which
  atomically upserts the document, optional work metadata, clears stale chunks,
  and inserts replacement chunks.
- Migrated ingest to that operation, preserving schema and FTS trigger use.
- Added rollback coverage proving a failed duplicate chunk insert preserves the
  previous document and chunk state.
- Introduced typed capability contracts in `domain/index/store/` (`DocumentStore`
  and `EmbeddingStore` interfaces) using domain types from `domain/index`.
- Implemented the adapter in `infrastructure/sqlite/index/` (`Adapter` struct)
  wrapping `*sqlite.IndexDB`, implementing both store interfaces with type mapping.
- Updated `IngestService` to accept a `store.DocumentStore` via
  `NewIngestServiceFromStore`; the existing `NewIngestService(dbPath)` constructor
  creates the adapter internally (backward compat, no CLI changes needed yet).
- Updated `EmbedService` to accept a `store.EmbeddingStore` via
  `NewEmbedServiceFromStore`; same backward compat approach.
- Fixed embedding metadata-read error propagation: `GetEmbeddingMetadata` errors
  are now returned to the caller rather than silently skipping the mixed-model
  guard (previously `if metaErr == nil && ok` swallowed real read failures).
- Focused tests: `replace_document_test.go` (atomic replacement, rollback,
  metadata upsert, FTS sync), `adapter_test.go` (document/chunk/embedding
  persistence separately, interface conformance), `store_test.go` (mock-based
  service tests proving no SQLite construction in the service path, metadata-read
  error propagation, mixed-model rejection). All 28 non-integration tests pass.

## Next Slice

Migrate CLI composition to inject the `Adapter` explicitly so the application
layer is fully free of infrastructure construction:

- Remove the dbPath-opening backward-compat path from `IngestService.Execute`
  and `EmbedService.Execute`.
- Have `platform/cli/` create `sqlite.IndexDB` and `sqliteindex.NewAdapter(db)`
  directly and pass the adapter to the service constructors.
- Prove the CLI indexing flow retains resolved work metadata end-to-end.

## Outcome

Introduced `domain/index/store` typed contracts (`DocumentStore`, `EmbeddingStore`) and `infrastructure/sqlite/index/adapter.go` implementing them. Application index services now receive typed stores via constructors rather than constructing SQLite directly. Fixed silent `GetEmbeddingMetadata` error suppression in embed-service. Backward-compatible constructors preserved for existing callers.

## Verification

- `go test ./internal/domain/index ./internal/application/index ./internal/infrastructure/sqlite/...` all pass.
- Build: `CGO_ENABLED=0 go build ./...` clean.
- Merged to development as PR #5.

## Close Checklist

- [x] `domain/index/store` contracts defined
- [x] SQLite adapter implements both contracts
- [x] Application services accept contracts via constructors
- [x] Silent error suppression fixed in embed-service
- [x] Backward-compat constructors preserved
- [x] Merged to development
