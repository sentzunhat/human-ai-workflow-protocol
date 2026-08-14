---
work-item: c1d2e3ff
type: fix
title: "Group SQLite persistence by capability and remove raw result maps"
status: plan-ready
created: 2026-08-10
updated: 2026-08-10
parent: b6c4e8a2
depends-on: c1d2e3fe
---

# Fix: SQLite Capability Groups

## Mission

Group SQLite persistence by document, chunk, embedding, and search capability
while replacing raw result maps at the audited boundaries.

## Done When

- Capability contracts are typed and colocated with their owning group.
- Benchmark and application callers no longer require untyped storage rows at
  the migrated seams.
- Schema, performance baselines, and targeted tests remain valid.
