---
work-item: c1d2e3fe
type: audit
title: "Recursive audit: SQLite infrastructure capabilities"
status: plan-ready
created: 2026-08-10
updated: 2026-08-10
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
