# Persist work UUID metadata during search index ingestion

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `e752eef1-9d9e-4dfb-95e6-cf5887cccccf`
**Type:** bug
**Reported:** 2026-09-21

---

## Input (verbatim)

> PR #41 review finding: search index derives lifecycle status from canonical work paths but never sets WorkUUID, so documents_metadata and work-aware chunk context are omitted. Derive canonical active/parked/closed IDs and add persistence/search coverage.

## Intake Summary

`hawp search index` recognizes work lifecycle folders but drops the work-item
UUID before handing documents to the ingest service. As a result, work metadata
and work-aware chunk context are never persisted for the normal work corpus.

## Current Context

The affected corpus walker is
`librarian/src/internal/platform/cli/index/ingest/corpus.go`. The persistence
gate is `librarian/src/internal/application/index/ingest-service.go`, which
only writes `documents_metadata` when `EnrichedDocument.WorkUUID` is non-nil.

## Initial Analysis

**Directly verified:**

- `walkWorkFiles` derives `Status` from `active`, `parked`, and `closed`
  paths but does not set `WorkUUID` on the emitted document.
- `IngestService.Execute` persists `documents_metadata` and adds work context
  only when `WorkUUID != nil`.
- The default search configuration indexes `.hawp/work`, so canonical work
  records take this path.

**Inferred (not yet proven):**

- Legacy non-UUID work paths should remain indexable without fabricated
  metadata.

**Likely scope:**

- `librarian/src/internal/platform/cli/index/ingest/corpus.go`
- `librarian/src/internal/platform/cli/index/ingest/corpus_test.go`
- `librarian/src/internal/application/index/ingest-service_test.go`

## Risk + Review Gate

**Risk:** medium
**Gate:** explicit user authorization received to implement the reviewed fix

## Plan

### Root cause

The corpus walker has only enough path parsing to assign a broad lifecycle
status. It does not parse the UUID-bearing segment of canonical active,
parked, or dated closed paths, even though the downstream persistence contract
requires that UUID.

### Options considered

1. Derive and validate UUIDs in the CLI corpus walker, leaving legacy paths
   without metadata.
2. Infer a UUID from document content or fabricate one for legacy paths.

### Recommended fix

Use option 1. Extract the candidate ID only from canonical folder layouts,
validate it with the existing work identity contract, and set status/UUID as a
pair. This retains non-canonical records as ordinary searchable documents.

### Verification target

Add corpus and persistence coverage for `active/<uuid>/plan.md` and
`closed/YYYY/MM/DD/<uuid>/plan.md`. Assert persisted `work_uuid` and `status`
values through the index repository. Run focused index tests, full Go tests,
vet, build, HAWP checks, generated distribution validation, source-layout
tests, and `git diff --check`.

## Backlog + Plan Link

**Status now:** in-progress
**Plan file:** work/active/e752eef1/plan.md

## Next Step

- [x] Investigation recorded above
- [x] Plan written
- [x] User authorized implementation
- [x] Implement and verify

## Outcome

The work corpus walker now extracts UUIDs only from canonical active, parked,
and dated closed work-item paths. It passes that UUID with the derived status
to the ingest service, allowing the index to persist `documents_metadata` and
include work-aware chunk context. Non-canonical legacy paths remain searchable
without fabricated metadata.

## Verification

Directly verified:

- `TestWorkCorpusPersistsCanonicalUUIDMetadata` indexes active and dated-closed
  records and verifies persisted `work_uuid:status` pairs; a non-canonical
  work path produces no metadata row.
- Focused corpus/index tests passed.
- `make check` passed from `librarian/src` (`go vet`, full `go test`, build).
- `go run ./cmd/hawp check` passed.
- `go run ./cmd/hawp work validate` passed with zero issues and one existing
  verification-clarity warning.
- `go run ./cmd/hawp distribution validate` passed.
- `go test ./...` passed from `scripts/source-layout`.
- `git diff --check origin/main` passed after the branch whitespace cleanup.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Plan file moved to closed/YYYY/MM/DD
- [x] BACKLOG.md updated
- [x] Status report not needed; plan carries the bounded evidence
