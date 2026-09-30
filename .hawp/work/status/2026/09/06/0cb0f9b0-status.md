# Status Report

## Intent

Prepare a preservation-safe single-executable migration for v0.0.24.
Primary record: [migration plan](../../../../closed/2026/09/07/0cb0f9b0/plan.md).

## Current State

In progress. JSON configuration safety and local installation improved; the
wrapper/native installation contract is unchanged. Later architecture and
unknown-record work remain queued behind the migration gate.

## What Was Inspected

Canonical install/update sources, MCP configuration writers/tests, CLI update
handlers, Makefile, quality workflow, operating guidance, and active work records.

## What Changed

Extracted JSON config merging into a small module. Preserve custom settings and
numeric values; refuse incompatible shapes before writing. Local/CI installation
uses temporary-file replacement. Root/core guidance prefers available HAWP MCP
tools and same-item context updates without inventing an update tool.

**New findings from architecture audit (2026-09-06):** Type duplication between
`domain/index.Chunk` and `sqlite.Chunk` identified as critical blocker for
future CLI/provider work — though not directly part of migration scope, it
affects the ingest-service.go which is relevant to single-executable concerns.

**Strategic decision:** Single-executable migration (`0cb0f9b0`) deferred to
v0.1.0 scope pending cross-platform testing and provider coordination (see
checkpoint at `2026/09/06/checkpoint-2026-09-06-architecture-audit.md`).

## What Was Directly Verified

Focused and full Go tests, vet, six standard target builds, distribution
build/validation, root/core guide parity, MCP work validation, and native HAWP
check passed. In-place install initially yielded exit 137; replacing the file
by rename restored execution. The install was repeated successfully.

**Cross-referenced with architecture audit (`47c793d6`):** All CLI command
handlers now use strict typed parsers; schema-aware backlog intake verified;
HAWP validation clean across all 3 checks (kit, work, links).

## What Remains Unproven

Full single-executable upgrade compatibility; Codex table customization
preservation; live provider reload and cross-platform execution. Unknown-folder
routing, UUID enrichment, storage extraction, and expanded indexing are plans,
not completed features. No hosted CI or downstream rehearsal ran in this step.

## Constraints

Preserve work copies, sidecars, identities, and unrelated changes. No downstream
edits, push, release publication, or PR merge. Do not claim signature-cache
involvement as the established cause of exit 137.

## Help Wanted

Review supported configuration-migration forms and preservation tests before
retiring launchers.

## Suggested Next Step

Finish Codex configuration preservation and config-only migration; then change
installer/update paths and regenerate guides with temporary-repository proofs.
