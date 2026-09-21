# Reject external Markdown symlinks during corpus indexing and export

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `225708ae-cbed-4589-9a13-7b947bb28e7c`
**Type:** bug
**Reported:** 2026-09-20

---

## Input (verbatim)

> Reject external Markdown symlinks during corpus indexing and export

## Intake Summary

Repository-side Markdown file symlinks are accepted by both the current search
ingest walkers and the older enriched-corpus builder. `os.ReadFile` follows the
link, so an in-repository `credentials.md` symlink can copy a readable external
file into SQLite search chunks or a JSON corpus export.

## Current Context

This is distinct from configured `../` path traversal: even fixed, legitimate
index roots remain unsafe when a file entry inside them is a symlink. Directory
symlinks encountered by the current walkers are not recursively traversed, but
file symlinks with a `.md` name are read.

## Initial Analysis

**Directly verified:**

- `librarian/src/internal/platform/cli/index/ingest/corpus.go:39-164` filters
  entries by directory status and `.md` suffix, then calls `os.ReadFile`
  without rejecting symlinks or checking the resolved target.
- `librarian/src/internal/infrastructure/markdown/markdown.go:25-47` similarly
  includes `.md` symlink entries in the shared collection helper.
- `librarian/src/internal/application/index/build-service.go:31-52` passes
  collected paths to enrichment readers; `BuildResult.Export` serializes full
  content when requested.
- The resulting search chunks can be returned through CLI, RAG, and MCP search.

**Inferred (not yet proven):**

- Broken or unreadable links fail or are skipped; the leak requires a readable
  target.
- Indexing/export is user-invoked, and no automatic transmission to the
  repository supplier was established.
- Other mutation callers of the shared Markdown collector may need a separate
  follow-up if implementation review proves their security impact; this item
  must not silently broaden into unrelated refactoring.

**Likely scope:**

- `librarian/src/internal/platform/cli/index/ingest/corpus.go`
- `librarian/src/internal/platform/cli/index/ingest/corpus_test.go`
- `librarian/src/internal/infrastructure/markdown/markdown.go`
- `librarian/src/internal/application/index/build-service.go`
- Focused tests for both SQLite ingestion and JSON export.

## Root Cause

The walkers equate "not a directory and named `.md`" with a safe document.
Neither `filepath.Walk` callbacks nor `os.ReadDir` regular-file filters prove
that the entry is not a symlink, and the subsequent `os.ReadFile` follows the
external target while retaining the harmless-looking repository path.

## Options Considered

1. Reject all symlink and non-regular entries at each indexing walker. Small
   and explicit, but duplicates policy.
2. Introduce a shared safe Markdown collector/read boundary that rejects
   symlinks and non-regular files and verifies resolved containment. This also
   protects the legacy export path. **Recommended**, provided mutation callers
   are reviewed before changing shared behavior.
3. Allow symlinks whose canonical targets remain inside the collection root.
   More compatible, but adds platform and race complexity with little proven
   need.

## Recommended Fix

- Default to rejecting symlinks and other non-regular entries before reads.
- Make the accepted root explicit so every returned file has both lexical and
  physical containment evidence.
- Apply the policy to current `search index` walkers and the shared collector
  used by `index build --export`.
- Return a deterministic error or documented skip; do not silently ingest an
  external target.

## Verification Plan

- Add tests for readable external file symlinks, internal file symlinks,
  symlinked directories, broken links, and non-regular entries where portable.
- Assert external target content never reaches SQLite chunks or exported JSON.
- Preserve ordinary kit/work/custom Markdown indexing behavior.
- Run focused index tests, `go test ./...`, `go vet ./...`,
  `git diff --check`, `go run ./cmd/hawp check`, and
  `go run ./cmd/hawp work validate` from `librarian/src`.

## Risk + Review Gate

**Risk:** high — security boundary and shared filesystem-walker behavior
**Gate:** explicit review before implementation

**Overlap check:** no current product-source edits overlap this planned index
scope. Review all callers before changing the shared collector.

## Backlog + Plan Link

**Status now:** plan-ready
**Plan file:** work/active/225708ae/plan.md

## Next Step

- [x] Investigation recorded above
- [x] Plan written with options, recommendation, and verification scope
- [x] Backlog moved to `plan-ready`
- [ ] Obtain explicit approval before implementation
