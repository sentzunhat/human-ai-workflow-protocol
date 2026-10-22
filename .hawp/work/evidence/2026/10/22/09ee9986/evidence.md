---
uuid: 09ee9986
title: "Read-only search index findings"
type: evidence
date: 2026-10-22
---

# Read-only search/index findings

## Search behavior

- `hawp/hawp_search` query `search index` returned matches.
- Query `backlog` returned a match from the older indexed path
  `.hawp/work/active/387a37b7-files.md`.
- Query `backlog alignment database freshness corpus regression` returned no
  results.
- In `librarian/src/internal/infrastructure/repositories/index/index.go`,
  `sanitizeFTSQuery` joins terms with spaces (FTS5 implicit `AND`).
  `librarian/src/internal/application/search/service.go` returns an empty
  result when lexical search finds no rows; it does not then fall back to
  semantic ranking.

## Local index state

Read-only SQLite inspection of `.hawp/db/index.sqlite` reported:

- `PRAGMA quick_check`: `ok`
- Documents: 609 total (17 custom, 116 kit, 476 work)
- Chunks: 5,109
- Work metadata rows: 118
- Latest chunk update: `2026-10-16 14:07:02`
- `09ee9986` and `b152c56c`: no indexed document paths
- `c147f2ae`: indexed at `.hawp/work/active/c147f2ae/plan.md` with status
  `active`; the source plan is at
  `.hawp/work/closed/2026/10/16/c147f2ae/plan.md`

`.hawp/config/search.json` includes `.hawp/kit` and `.hawp/work`. The CLI
builds a corpus from configured files, but `IngestService.Execute` upserts
documents and chunks without pruning paths absent from the new corpus. This
accounts for stale moved/deleted paths remaining in an existing database.

## Regression coverage gap

`librarian/src/internal/platform/cli/index/ingest/corpus_test.go` verifies
corpus walking and work metadata. `librarian/src/internal/infrastructure/repositories/index/fts_sync_test.go`
proves a synthetic chunk is searchable after fresh schema creation.
`librarian/src/internal/platform/mcp/server/tools_e2e_test.go` covers a
manually seeded search result and the no-match response. There is no test
covering a fresh configured corpus from both kit and work files through to
search, or verifying stale paths disappear after a full sync.

## HAWP MCP usage

- `hawp/hawp_search`: used for short and longer query comparisons.
- `hawp/hawp_usage`: returned aggregate totals of 272 calls, about 7,409 input
  tokens, and about 67,449 output tokens; no per-tool/query report was
  requested.
- `hawp/hawp_work_doc`: created this evidence document for existing work item
  `09ee9986`.

The full read-only findings and repository-root proof are recorded in
[`.hawp/work/active/09ee9986/plan.md`](../../../../../active/09ee9986/plan.md).
At the time of this read-only investigation, no search code had been changed;
the persistent local index was not rebuilt or modified.

## Validation

- `hawp/hawp_work_validate` reported 3 kit checks passed with 0 issues;
  work validation failed with 2 issues and 1 warning.
- Detailed `go run ./cmd/hawp work validate` output: the two existing missing
  plan paths are `active/2eea565c/plan.md` and `active/429e075e/plan.md`;
  closed-task completeness passed; evidence integrity passed (9/9); verification
  clarity reported 309/392 evidence-backed, 73 explicitly unproven, and 10
  ambiguous. No new issue was reported for `09ee9986`.
- The HAWP link check found the same two missing active-plan links in
  `.hawp/work/BACKLOG.md`.
- `git diff --check` passed at this investigation checkpoint; branch commits
  and pushes were not yet authorized.

## Search-index implementation and verification (2026-10-22)

- `IngestService.Execute` remains partial. The new full-corpus `Sync` rejects
  nil/empty corpora, ingests current files, then removes stale chunks,
  metadata, and documents in a transaction. Chunk-delete triggers maintain
  the FTS table.
- The configured `hawp search index` command now uses `Sync`. Search guidance
  documents stale-path cleanup and FTS5's implicit `AND` behavior.
- Regression tests cover partial-ingest preservation, empty-corpus rejection,
  FTS deletion, fresh configured kit/work indexing, moved work paths, and
  MCP tool-search results.
- Focused Go tests passed for the application index, SQLite index, CLI ingest,
  and MCP server packages. The full `go test ./...` suite passed on
  `feature/mcp-providers`.
- A later `hawp/hawp_search` query for `empty search configured corpus moved
  paths index freshness` returned `results: []` from the existing
  MCP-configured index. The persistent `.hawp/db/index.sqlite` was not
  rebuilt, so this is not a test of the new code-branch implementation.
- `hawp/hawp_work_doc` resolved the canonical evidence path for existing work
  item `09ee9986`; no work item was created.
- `hawp/hawp_work_validate` against the original checkout passed 3 kit checks
  with 0 issues, while work validation reported 2 issues and 1 warning. The
  link check found the known missing `active/2eea565c/plan.md` and
  `active/429e075e/plan.md` targets in `.hawp/work/BACKLOG.md`.
- After syncing the updated work records to `feature/mcp-providers`, the
  canonical code-branch validation passed: `go test ./...`, distribution
  sync (0/40 files changed), and `hawp check` (3/3 checks, 0 work issues or
  warnings, 147 Markdown links valid).
- The implementation was committed and pushed on `feature/mcp-providers` as
  `a35dd925a04c397e0a8ce363b3066e7a1edc085b`. The website branch was also
  committed and pushed; details are in
  [`.hawp/work/active/f8685e56/plan.md`](../../../../../active/f8685e56/plan.md).
- No persistent local index rebuild was performed. Live search against a
  restarted MCP server and rebuilt index remains unverified and deferred.
