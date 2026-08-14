---
work-item: c1d2e3fa
type: audit
title: "Recursive audit: application index capability"
status: done
created: 2026-08-10
updated: 2026-08-13
parent: b6c4e8a2
follow-up: c1d2e3fb
---

# Audit: Application Index Capability

## Outcome

The index capability has a clear smallest next boundary: application services
need typed persistence contracts, while SQLite remains the capability-local
adapter. The current CLI independently builds its corpus and bypasses shared
metadata resolution.

## Findings

- CLI indexing skips `BACKLOG.md`, assigns all work files `closed`, and does
  not use the shared build service.
- Ingest and embedding services construct SQLite and call persistence details
  directly, making application tests storage-dependent.
- Document upsert, stale-chunk deletion, and replacement inserts are separate
  operations, so a mid-write failure can leave partial state.
- Build and ingest use parallel document contracts with a CLI conversion.
- Embedding metadata-read errors are suppressed before the compatibility guard.

## Follow-up

`c1d2e3fb` will introduce `domain/index/store` contracts for atomic corpus
replacement and embedding persistence, implemented by
`infrastructure/sqlite/index`. It must preserve the current schema, FTS
triggers, CLI behavior, and model guards.

## Verification

- `go test ./internal/application/index ./internal/infrastructure/sqlite ./internal/infrastructure/sqlite/search`
