---
work-item: c42d1443-fbd0-4a98-8379-6058ace02dc4
type: fix
title: "Unify search index through shared build and ingest flow"
status: plan-ready
created: 2026-08-15
updated: 2026-08-15
parent: b6c4e8a2
depends-on: c1d2e3fb
---

# Fix: Shared Search Index Flow

## Mission

Route `hawp search index` through the application index build and ingest flow
so work UUID, status, and folder metadata reach persistence consistently.

## Constraints

- Preserve CLI flags, lexical indexing behavior, schema, FTS triggers, and
  existing output where possible.
- Do not reintroduce CLI file walkers or direct SQLite orchestration.
- Keep conversion compatibility-only until external JSON consumers are proven.

## Done When

- CLI indexing uses one application reindex use case built from shared corpus
  acquisition and typed persistence contracts.
- A fixture proves work metadata reaches indexed search results.
- Focused index, SQLite, CLI tests, build, and diff checks pass.
