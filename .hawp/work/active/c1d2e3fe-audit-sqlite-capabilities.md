---
work-item: c1d2e3fe
type: audit
title: "Recursive audit: SQLite infrastructure capabilities"
status: done
created: 2026-08-10
updated: 2026-08-21
parent: b6c4e8a2
follow-up: c1d2e3ff
---

# Audit: SQLite Infrastructure Capabilities

## Mission

Audit nested SQLite document, chunk, vector, search, benchmark, and transaction
responsibilities for capability cohesion and raw-map leakage.

## Required Output

- Storage responsibility map.
- Confirmed raw-map and cross-capability coupling findings.
- Safe extraction order and verification plan for `c1d2e3ff`.

## Constraints

Preserve schema, migrations, CGO-free build behavior, and existing search
fallbacks. Do not redesign storage in the audit item.

---

## Audit Evidence

Audited 2026-08-21. All files under `librarian/go/internal/infrastructure/sqlite/`.

### Storage Responsibility Map

| File / Type | Capabilities owned |
|---|---|
| `index.go` / `IndexDB` | Connection management (`Open`, `Close`); schema + migrations (`InitSchema`); document storage (`InsertDocument`, `DocumentCount`); document metadata (`InsertMetadata`, `DocumentMetadata`); chunk storage (`InsertChunk`, `DeleteChunksForDocument`, `ChunkCount`); atomic document + chunk + metadata upsert (`ReplaceDocument`, `DocumentRow`); FTS5 lexical search (`QueryChunksLexical`); vector read/write (`UpdateChunkEmbedding`, `GetChunkVector`, `GetAllChunkVectors`, `GetChunkVectors`, `GetChunksNeedingEmbedding`, `ChunksNeedEmbedding`, `HasVectors`, `ChunkData`); embedding metadata (`GetEmbeddingMetadata`, `SetEmbeddingMetadata`, `EmbeddingMetadata`); transaction control (`BeginTx`, `Commit`, `Rollback`) |
| `index/adapter.go` / `sqliteindex.Adapter` | Domain-to-sqlite type translation; implements `store.DocumentStore` (ReplaceDocument) and `store.EmbeddingStore` (embedding CRUD, transaction control) |
| `search/adapter.go` / `search.Adapter` | Implements domain search index port; translates raw map rows from `QueryChunksLexical` to typed `searchindex.Candidate`; wraps `HasVectors`, `GetEmbeddingMetadata`, `GetChunkVectors` |
| `index/adapter_test.go` | Compile-time and behavioral tests for `store.DocumentStore` and `store.EmbeddingStore` contracts |
| `index_test.go` | Vector round-trip and idempotent ingest tests against `IndexDB` directly |
| `embedding_metadata_test.go` | Embedding metadata upsert tests against `IndexDB` directly |
| `fts_sync_test.go` | FTS5 trigger correctness tests (locks in the 2026-07-27 bug fix) |
| `replace_document_test.go` | Atomicity, rollback, metadata, and FTS-sync tests for `ReplaceDocument` |
| `transaction_test.go` | Integration test for batch-commit embedding writes (skipped unless `INTEGRATION=1`) |

### Confirmed Raw-Map Findings

**Finding 1 — `QueryChunksLexical` returns `[]map[string]interface{}`** (`index.go`:268–315)

`QueryChunksLexical` is a public method on `IndexDB`. It assembles each scanned row into a
`map[string]interface{}` with string-keyed fields (`id`, `text`, `path`, `type`, `category`,
`folder_role`, `chunk_idx`, `folder_context`, `work_uuid`, `status`, `closed_at`). This is
the only public boundary in the `sqlite` package that returns an untyped map.

Direct evidence: `index.go` lines 285–313.

**Finding 2 — `search/adapter.go` consumes the raw map via coercion helpers**

`search/adapter.go:LexicalSearch` calls `QueryChunksLexical` and ranges over
`[]map[string]interface{}`. It reads each value via `stringValue(row["..."])` and
`intValue(row["..."])` — private helper functions (lines 61–83) that handle `nil`, `*string`,
and numeric type assertions. These helpers exist solely to paper over the untyped boundary.

Direct evidence: `search/adapter.go` lines 22–43.

No other public functions return `map[string]interface{}` or `map[string]any`.
`GetAllChunkVectors` and `GetChunkVectors` return `map[int64][]float32`, which is typed.

### Cross-Capability Coupling Findings

**Finding 3 — `IndexDB` conflates six distinct capabilities in one type/file**

`index.go` owns: (a) schema/migrations, (b) document + metadata persistence, (c) chunk
persistence, (d) FTS5 lexical search, (e) vector read/write, (f) embedding metadata, and
(g) transaction state. A caller importing `sqlite` to embed vectors also imports FTS search,
document upsert, and schema management with no way to exclude them.

The `tx *sql.Tx` field on `IndexDB` makes the struct stateful; `UpdateChunkEmbedding` inspects
`ix.tx` to choose its execer (`index.go` lines 341–355). Transaction management is entangled
with vector write rather than isolated.

**Finding 4 — `InsertDocument`/`InsertMetadata` are superseded but not removed**

`ReplaceDocument` (line 484) embeds the same SQL as `InsertDocument` and `InsertMetadata`.
Both older individual methods remain public, creating two code paths for document upsert. The
production ingest pipeline appears to use `ReplaceDocument` exclusively; the loose methods
serve only direct unit tests.

**Finding 5 — Both adapters depend on all of `*sqlite.IndexDB`**

`search/adapter.go` and `index/adapter.go` each hold a `*sqlite.IndexDB` and have access to
every method, but each uses only a subset. The search adapter needs only `QueryChunksLexical`,
`HasVectors`, `GetEmbeddingMetadata`, and `GetChunkVectors`. There is no interface or subset
type limiting each adapter to its own surface.

---

## Handoff To c1d2e3ff

### Safe Extraction Order

Each step must leave `go build ./librarian/go/...` and all non-integration tests passing
before proceeding. Steps 2–6 are same-package file splits and carry no API risk.

**Step 1 — Replace the raw-map boundary in `QueryChunksLexical` (highest priority)**

Define a typed `LexicalRow` struct. Change `QueryChunksLexical` to return `[]LexicalRow`
instead of `[]map[string]interface{}`. Update `search/adapter.go:LexicalSearch` to read
struct fields directly; remove `stringValue`/`intValue` helpers once no callers remain.

This is the only change visible across package boundaries. Doing it first means every
subsequent organizational split starts from a clean typed surface.

**Step 2 — Extract vector/embedding methods into `sqlite/embedding.go` (file-level move)**

Move to a new file: `UpdateChunkEmbedding`, `GetChunksNeedingEmbedding`, `ChunksNeedEmbedding`,
`GetChunkVector`, `GetAllChunkVectors`, `GetChunkVectors`, `HasVectors`, `GetEmbeddingMetadata`,
`SetEmbeddingMetadata`, and the `EmbeddingMetadata` and `ChunkData` types. Methods remain on
`IndexDB`. No signature or import change — pure reorganization.

**Step 3 — Extract FTS/search into `sqlite/search.go` (file-level move)**

Move `QueryChunksLexical` and the `LexicalRow` type (from Step 1) to `sqlite/search.go`.

**Step 4 — Extract document/chunk/metadata into `sqlite/document.go` (file-level move)**

Move `InsertDocument`, `InsertMetadata`, `DeleteChunksForDocument`, `ReplaceDocument`,
`DocumentCount`, `ChunkCount`, and types (`DocumentRow`, `DocumentMetadata`, `Chunk`) to
`sqlite/document.go`.

**Step 5 — Extract transaction control into `sqlite/tx.go` (file-level move)**

Move `BeginTx`, `Commit`, `Rollback`. The `tx *sql.Tx` field stays on `IndexDB` (defined in
`index.go`); note in a comment that it creates statefulness and flag for a follow-up if
needed.

**Step 6 — Extract schema into `sqlite/schema.go` (file-level move)**

Move `InitSchema` and the SQL string. `Open`, `Close`, `ensureDir` stay in `index.go` as the
connection-management entry point.

### Verification Plan for c1d2e3ff

After **Step 1** (raw-map removal):

- `go vet ./librarian/go/...` — no type errors
- `go build ./librarian/go/...` — clean build
- `go test ./librarian/go/internal/infrastructure/sqlite/...` — all non-integration tests pass
- `go test ./librarian/go/internal/infrastructure/sqlite/search/...` — search adapter passes
- Confirm `stringValue`/`intValue` helpers are deleted from `search/adapter.go`

After **each of Steps 2–6** (file splits):

- `go build ./librarian/go/...` — clean build
- `go test ./librarian/go/internal/infrastructure/sqlite/...` — all non-integration tests pass

After **all steps**:

- `go test ./librarian/go/...` — full non-integration suite passes
- `go vet ./librarian/go/...` — clean
- Confirm: `grep -r 'map\[string\]interface{}\|map\[string\]any' librarian/go/internal/infrastructure/sqlite/*.go` returns no hits on function signatures
- Confirm both adapters use only typed struct fields from their respective `sqlite` returns
