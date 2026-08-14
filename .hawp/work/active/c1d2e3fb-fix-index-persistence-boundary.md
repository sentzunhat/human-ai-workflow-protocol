---
work-item: c1d2e3fb
type: fix
title: "Introduce typed index persistence capability contracts"
status: in-progress
created: 2026-08-10
updated: 2026-08-13
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

## Next Slice

Introduce typed index-store contracts over this existing atomic operation and
propagate embedding metadata-read failures before migrating CLI composition.
