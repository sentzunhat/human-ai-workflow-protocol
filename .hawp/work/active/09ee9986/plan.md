# Diagnose empty HAWP search and add index regression tests

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `09ee9986-5d82-4f24-a1bd-0a814eb402a0`
**Type:** bug
**Reported:** 2026-10-22

---

## Input (verbatim)

> Investigate and fix HAWP MCP search returning empty results in this VS Code session. Verify the connected repository and local index/corpus scope, document the correct rebuild/recovery path, and add deterministic regression tests proving indexed HAWP kit/work files remain searchable after a fresh build. Preserve the source files as the authority; do not treat zero search results as proof that records are absent.

## Intake Summary

The current MCP server is connected to this HAWP repository, but targeted
`hawp_search` calls returned empty result arrays and `hawp_work_intake`
previously returned `needs_user_input` with no indexed context. Empty results
do not establish that source files or work records are missing. This is a new
recurrence: `.hawp/work/closed/2026/10/16/c147f2ae/plan.md` records a
successful local index rebuild and useful searches on 2026-10-16.

## Current Context

- On 2026-10-22, two targeted `hawp_search` queries returned
  `results: []`; neither returned a tool error.
- `hawp_work_validate` identified the expected HAWP repository. Kit checks
  passed (3 checks, 0 issues); work validation reported 2 issues and 1 warning;
  link checking found backlog links to missing active plans
  `active/2eea565c/plan.md` and `active/429e075e/plan.md`. These link issues
  predate this search investigation.
- Detailed work validation confirmed the missing-plan issues and reported
  verification clarity as 309/392 evidence-backed, 73 explicitly unproven,
  and 10 ambiguous.
- `.hawp/config/search.json` scopes indexing to `.hawp/kit`, `.hawp/work`,
  `core/.hawp`, `core/providers`, `librarian/docs`, and `README.md`.
- `librarian/src/internal/platform/mcp/server/tool_search.go` returns an
  error when the index cannot be opened, but returns an empty `results` array
  when the lexical query yields no rows. Thus an empty response does not
  identify index freshness or corpus coverage.
- `librarian/src/internal/platform/mcp/server/tools_e2e_test.go` covers
  missing-index errors, no-match behavior, and a successful structured search.
  `.hawp/work/active/c8d9e1cb/plan.md` also records missing-index and no-match
  intake tests; it does not own this recurrence or a fresh-build search regression.
- `.hawp/work/closed/2026/10/16/c147f2ae/plan.md` records a user-approved
  rebuild on 2026-10-16 that indexed 597 documents, 5,001 chunks, and 118
  work-item metadata records; subsequent searches found relevant files. That
  historical result does not prove the current local index is fresh or covers
  the current repository.

## Initial Analysis

**Directly verified:**

- `hawp_search` returned empty result arrays for two targeted queries in this
  session; the HAWP MCP server returned no search error.
- `hawp_work_validate` identified this HAWP repository and reported the
  current validation results recorded above.
- The MCP search implementation distinguishes an index-open error from a
  successful query with zero lexical rows.
- Existing tests cover missing-index and no-match responses, and a successful
  structured search response.
- The configured source paths and the prior successful rebuild are recorded
  in repository files listed above.
- The full Go test suite and distribution consistency check passed. The
  aggregate HAWP check remains failed only on the known backlog/link issues
  recorded above; no local search index was rebuilt or changed.

**Inferred (not yet proven):**

- The current index may be stale, incomplete, or scoped to a different corpus;
  an ordinary query mismatch is also possible. The root cause is unknown.
- A deterministic regression should build an isolated index from fixture
  source files and verify retrieval, rather than depend on the user's local DB.

**Likely scope:**

- Inspect current index configuration and metadata without mutating the local
  database; trace indexing and MCP search from source files to stored chunks.
- Add an isolated fresh-build regression covering representative `.hawp/kit`
  and `.hawp/work` source files, plus existing no-match and missing-index
  semantics.
- Do not rebuild or overwrite the user's persistent index until the cause and
  recovery plan are clear and approved.

## Read-only Investigation Update (2026-10-22)

### Direct evidence

- `hawp/hawp_search` returned results for `search index` and `backlog`, but no
  results for `backlog alignment database freshness corpus regression`. This
  confirms the current index is not globally empty; it does not prove each
  earlier query had the same cause.
- `.hawp/db/index.sqlite` passed read-only `PRAGMA quick_check`. It contains
  609 documents (17 custom, 116 kit, 476 work), 5,109 chunks, and 118 work
  metadata rows. The newest stored chunk timestamp is `2026-10-16 14:07:02`.
- The index contains `.hawp/work/active/c147f2ae/plan.md` with status `active`,
  while the source plan is now
  `.hawp/work/closed/2026/10/16/c147f2ae/plan.md`. No indexed document path
  contains `09ee9986` or `b152c56c`.
- `.hawp/config/search.json` includes `.hawp/kit` and `.hawp/work` among its
  configured paths. `buildCorpusFromRepo` walks Markdown files under those
  paths, and `hawp search index` ingests that corpus.
- `IngestService.Execute` upserts documents by path and replaces their chunks,
  but it does not remove documents whose paths are absent from a later corpus.
  Therefore the current command refreshes additions and edits but does not
  fully reconcile moved or deleted files.
- `sanitizeFTSQuery` joins query terms with spaces, which FTS5 treats as
  implicit `AND`. The search service returns no results when that lexical
  query has no rows; it does not fall back to semantic ranking in that case.
- Existing tests separately cover corpus walking and metadata persistence,
  FTS search after inserting a synthetic chunk into a fresh database, and MCP
  search response/no-match behavior. They do not cover building a fresh index
  from both kit and work source files and then searching those files.
- The user reported three earlier empty searches; the initial record captured
  two, and the exact query strings are not retained. Their individual causes
  remain unproven.
- `hawp/hawp_usage` returned aggregate usage totals of 272 calls, about 7,409
  input tokens, and about 67,449 output tokens. A per-tool/query report was not
  requested.

### Assessment

**Confirmed:** The local index is structurally valid, has working matches for
short queries, and is stale relative to current work files. Long multi-term
queries can return empty results because every sanitized term must match the
same indexed chunk. Re-running the current `hawp search index` command would
not remove stale paths.

**Likely, not proven for the historical queries:** Staleness explains why
recent work items are absent, while strict `AND` matching explains why a
longer query can return no rows. The missing original query strings prevent
attributing every earlier empty result to either cause.

### Proposed plan for review

**Option A — Reconcile the configured corpus and add an end-to-end regression
(recommended).**

- Give the full `hawp search index` operation source-sync semantics so moved
  or deleted source paths do not remain searchable. Keep partial ingest
  callers from pruning documents they did not supply.
- Add an isolated temporary-repository regression that builds a fresh index
  from representative `.hawp/kit` and `.hawp/work` files, verifies both are
  searchable, then verifies a moved or deleted work path is no longer
  returned after a full sync.
- Update `.hawp/kit/usage/search.md` to explain that lexical terms are combined
  with `AND`, that short focused queries are more reliable, and what the index
  command does about source changes.
- Preserve current query semantics in this change. Do not add an `OR` fallback
  without a separate precision/recall evaluation.
- Do not mutate the persistent local index; verification uses isolated
  temporary databases.

**Option B — Add an `OR` or broader-query fallback.** This may improve recall
for long searches, but it does not fix stale or moved documents and may add
noisy results. Defer it pending a search-quality test set.

**Implementation scope:**

- `librarian/src/internal/application/index/ingest-service.go` — add or call a
  full-corpus synchronization path while preserving partial `Execute`.
- `librarian/src/internal/application/index/ingest-service_test.go` — verify
  partial ingest does not prune and an empty full sync is refused.
- `librarian/src/internal/infrastructure/repositories/index/index.go` — add
  transactional stale-document cleanup with metadata/chunk removal.
- `librarian/src/internal/infrastructure/repositories/index/fts_sync_test.go`
  — verify pruning removes stale FTS matches.
- `librarian/src/internal/platform/cli/index/ingest/command.go` — use full
  synchronization only for the configured-corpus command.
- `librarian/src/internal/platform/cli/index/ingest/corpus_test.go` — add the
  isolated fresh-build and moved-path regression.
- `librarian/src/internal/platform/mcp/server/tools_e2e_test.go` — verify a
  freshly synced kit/work corpus is searchable through `hawp_search`.
- `.hawp/kit/usage/search.md` — document the verified query and sync behavior.

**Coordination:** No active or parked plan currently lists the proposed search
index source/test files. `.hawp/work/BACKLOG.md` is also being changed by
`b152c56c`; this update changes only the `09ee9986` row and preserves the
other work. Do not overlap with the MCP setup-guide edits from `b152c56c`.

### Approved implementation approach (2026-10-22)

The user approved Option A, and the isolated `feature/mcp-providers` worktree
is prepared and validated. Preserve `Execute` as the partial-ingest operation
and add a separate full-corpus sync operation used only by `hawp search index`.
Reject an empty corpus, ingest the current files, then transactionally remove
stale document metadata, chunks, and document rows; chunk-delete triggers
maintain FTS. Add tests for partial-ingest preservation, empty-corpus safety,
a fresh kit/work corpus, moved-path removal, and the MCP search response. Do
not rebuild the persistent local index.

### Repository-root proof (captured before this record update)

Machine-local prefixes are redacted; all status paths are repo-relative.

```text
$ pwd
<repo-root-abs>
$ git rev-parse --show-toplevel
<repo-root-abs>
$ git rev-parse --show-prefix
<empty output>
$ git status --short
 M .github/agents/docs-alignment.agent.md
 M .github/agents/hawp.agent.md
 M .github/agents/master-builder.agent.md
 M .github/copilot-instructions.md
 M .github/skills/hawp-change-review-and-reference-sync/SKILL.md
 M .hawp/kit/usage/mcp/README.md
 M .hawp/kit/usage/mcp/github-copilot.md
 M .hawp/work/BACKLOG.md
 M .hawp/work/active/b152c56c/plan.md
 M .hawp/work/status/2026/10/21/b152c56c/status.md
 M core/.hawp/kit/usage/mcp/README.md
 M core/.hawp/kit/usage/mcp/github-copilot.md
 M core/providers/.github/prompts/hawp-change-review-and-reference-sync.prompt.md
 M librarian/src/internal/platform/mcp/configure/config.go
 M librarian/src/internal/platform/mcp/configure/config_agent.go
 M librarian/src/internal/platform/mcp/configure/configure_test.go
?? .hawp/work/active/09ee9986/
```

## Risk + Review Gate

**Risk:** medium
**Gate:** Investigate read-only first. Present the confirmed root cause and
minimal regression/fix plan for review before changing search behavior or the
user's persistent index.

**Approval:** The user approved Option A on 2026-10-22. The isolated
`feature/mcp-providers` worktree is prepared and validated; implementation
starts there. Do not rebuild or modify the persistent local index.

## Backlog + Plan Link

**Status now:** in-progress (Option A implementation on
`feature/mcp-providers`)
**Plan file:** work/active/09ee9986/plan.md

## Next Step

- [x] Investigation recorded above (required before planning)
- [x] Write the initial plan from verified context
- [x] Inspect current local index metadata and reproduce search behavior
  without changing the persistent database
- [x] Trace the source-file indexing path and existing regression coverage
- [x] Propose the minimal fix and fresh-build regression tests
- [x] Validate the work-record update; findings are in
  [`.hawp/work/evidence/2026/10/22/09ee9986/evidence.md`](../../evidence/2026/10/22/09ee9986/evidence.md)
- [x] Receive user approval for Option A (2026-10-22)
- [x] Prepare and validate the isolated Librarian branch
- [x] Start implementation after the Librarian branch is prepared
- [x] Implement full-corpus sync, regression tests, and search guidance

## Code-branch repository proof before implementation (2026-10-22)

Machine-local prefixes are redacted; status paths remain repo-relative.

```text
$ pwd
<repo-root-abs>
$ git rev-parse --show-toplevel
<repo-root-abs>
$ git rev-parse --show-prefix
(empty)
$ git status --short
 M .github/agents/docs-alignment.agent.md
 M .github/agents/master-builder.agent.md
 M .github/copilot-instructions.md
 M .gitignore
 M .hawp/config/search.json
 M .hawp/kit/usage/mcp/README.md
 M .hawp/kit/usage/mcp/github-copilot.md
 M .hawp/kit/usage/mcp/lessons.md
 M .hawp/kit/usage/mcp/verification.md
 M .hawp/work/BACKLOG.md
 M core/.hawp/kit/usage/mcp/README.md
 M core/.hawp/kit/usage/mcp/github-copilot.md
 M core/providers/.claude/rules/hawp-core.md
 M core/providers/.codex/AGENTS.md.seed
 M core/providers/.continue/rules/hawp-01-core.md
 M core/providers/.cursor/rules/hawp-core.mdc
 M core/providers/.github/copilot-instructions.md
 M core/providers/.github/prompts/hawp-change-review-and-reference-sync.prompt.md
 M core/providers/doc-verification.md
 M core/providers/shared/behaviors/hawp-core.md
 M distribution/generated/claude/install/development.md
 M distribution/generated/claude/install/main.md
 M distribution/generated/claude/update/development.md
 M distribution/generated/claude/update/main.md
 M distribution/generated/codex/install/development.md
 M distribution/generated/codex/install/main.md
 M distribution/generated/codex/update/development.md
 M distribution/generated/codex/update/main.md
 M distribution/generated/continue/install/development.md
 M distribution/generated/continue/install/main.md
 M distribution/generated/continue/update/development.md
 M distribution/generated/continue/update/main.md
 M distribution/generated/cursor/install/development.md
 M distribution/generated/cursor/install/main.md
 M distribution/generated/cursor/update/development.md
 M distribution/generated/cursor/update/main.md
 M distribution/generated/github/install/development.md
 M distribution/generated/github/install/main.md
 M distribution/generated/github/update/development.md
 M distribution/generated/github/update/main.md
 M distribution/sources/shared/install.md
 M distribution/sources/shared/update.md
 M librarian/src/internal/platform/mcp/configure/config.go
 M librarian/src/internal/platform/mcp/configure/configure.go
 M librarian/src/internal/platform/mcp/configure/configure_test.go
 M librarian/src/internal/platform/mcp/server/server.go
?? .github/agents/hawp.agent.md
?? .github/skills/commit-many-small/
?? .github/skills/commit-one-big/
?? .github/skills/conservative-docs-drift-safe-cleanup/
?? .github/skills/docs-alignment-deterministic/
?? .github/skills/docs-alignment-simplicity-pass/
?? .github/skills/hawp-backlog-alignment/
?? .github/skills/hawp-change-review-and-reference-sync/
?? .github/skills/intake/
?? .github/skills/intent-first-handoff/
?? .github/skills/status-report/
?? .hawp/work/active/09ee9986/
?? .hawp/work/active/b152c56c/
?? .hawp/work/active/f8685e56/
?? .hawp/work/evidence/2026/10/
?? .hawp/work/status/2026/10/21/
?? librarian/src/internal/platform/mcp/configure/config_agent.go
?? librarian/src/internal/platform/mcp/server/prompts.go
?? librarian/src/internal/platform/mcp/server/prompts_test.go
?? librarian/src/internal/platform/mcp/server/resources.go
?? librarian/src/internal/platform/mcp/server/resources_test.go
```

## Implementation and verification update (2026-10-22)

Implemented on `feature/mcp-providers`:

- Added `IngestService.Sync` for complete corpora while retaining
  `IngestService.Execute` as partial ingestion. Full sync rejects nil or empty
  corpora, indexes current documents, then transactionally removes stale
  document chunks, metadata, and document rows. Existing chunk-delete
  triggers keep FTS synchronized.
- Routed `hawp search index` through full sync and updated search guidance for
  configured-corpus reconciliation and FTS5 `AND` query behavior.
- Added regression coverage for partial-ingest preservation, empty-corpus
  rejection, stale FTS cleanup, fresh kit/work corpus search, moved-path
  removal, and MCP tool search results.

**Verification:** The four focused Go packages passed, and `go test ./...`
passed on the code branch. `hawp/hawp_search` returned no results for the
verification query against the existing MCP-configured index; that index was
not rebuilt, so this does not exercise the new branch code. The persistent
`.hawp/db/index.sqlite` remains untouched. On `feature/mcp-providers`,
`go run ./cmd/hawp distribution sync` updated 0/40 files and
`go run ./cmd/hawp check` passed all 3 validations with 0 issues/warnings
across 147 Markdown files.

The read-only `hawp/hawp_work_validate` call against the original checkout
still reports 2 work issues and 1 warning, plus the known missing links to
`active/2eea565c/plan.md` and `active/429e075e/plan.md`. These are separate
from the passing code-branch check. Keep this item in progress until the code
is deployed/restarted and the approved local-index recovery path is run or
explicitly deferred.

The implementation was committed and pushed on `feature/mcp-providers` as
`a35dd925a04c397e0a8ce363b3066e7a1edc085b`. The persistent local index
remains untouched, so live search verification is still pending.
