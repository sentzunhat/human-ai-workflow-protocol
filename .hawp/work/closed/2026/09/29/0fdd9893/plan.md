# Fix PR 41 context and SQLite review findings

**UUID:** `0fdd9893-4756-4154-a2f0-ae6c0beb1001`
**Type:** bug
**Reported:** 2026-09-29
**Risk Level:** medium

### Input

> Deep dive and fix the new Copilot findings: prevent maxTokens multiplication overflow in context truncation; validate SQLite rollback journal path (including symlink and hard-link safety); exclude the failed benchmark query from saved-token totals; add regressions and resolve the PR comments.

### Context and Analysis

The review confirms an integer-overflow panic path, an omitted SQLite rollback-journal sidecar check, and a benchmark total that includes a failed query.

**Directly verified:**

- `librarian/src/internal/application/context/format.go` multiplied `maxTokens * 4` before checking whether the source string fits.
- `librarian/src/internal/infrastructure/filesystem/hawp_project.go` checked hard links on the main index, WAL, and SHM, but omitted the rollback journal and did not apply sidecar symlink checks.
- `librarian/src/tests/infrastructure/filesystem/hawp_project_test.go` lacked journal and symlinked-sidecar regressions.
- `benchmark/runs/2026-09-12-downstream-token-savings-ollama.md` included the failed query's 1,998 tokens; successful rows total 16,098 raw tokens, 831 shaped tokens, and 15,267 saved tokens, rounded to 95%.
- GitHub PR #41 was at `5bd9eb92`; all three inline threads were unresolved before the fix.

**Inferred:** SQLite rollback mode may create or mutate `index.sqlite-journal`, so that sibling requires the same preflight as WAL and SHM.

**Root causes:**

- `truncateToTokens` multiplies a caller-provided positive token count by four before comparing it with the source length.
- `ResolveSafeSearchIndexPath` omits `index.sqlite-journal` and does not apply sidecar symlink checks.
- The benchmark total sums all raw-token rows despite its 9/10 successful denominator.

**Coordination:** Owner Codex. Only this work item is known to touch these paths. Parallel risk low. Repo-root proof captured before edits: `pwd` showed repository root; `git rev-parse --show-toplevel` matched; `git rev-parse --show-prefix` was empty; `git status --short` was clean.

### Options

**A — Minimal guarded fixes (chosen).** Compare the token budget to a safely calculated byte-length threshold before multiplying; preflight all SQLite sidecars for symlink ancestry and hard links; correct the benchmark totals. Add focused regressions. This preserves existing APIs and behavior for normal budgets.

**B — Add a global CLI max-token cap and redesign SQLite opening.** This would constrain callers beyond the actual overflow site and expand the patch into connection/open architecture. It is broader than required and risks changing accepted CLI behavior.

### Recommended Fix

**Option chosen:** A. Keep the boundary safe for direct callers as well as CLI inputs; reject sidecars before SQLite receives the path; make the evidence table arithmetically consistent.

**Files to change:**

- `librarian/src/internal/application/context/format.go` and `format_test.go` — guard before multiplication and test `maxInt` plus non-positive input.
- `librarian/src/internal/infrastructure/filesystem/hawp_project.go` and `librarian/src/tests/infrastructure/filesystem/hawp_project_test.go` — cover the rollback journal and check symlink ancestry and hard links for every SQLite sidecar.
- `benchmark/runs/2026-09-12-downstream-token-savings-ollama.md` — exclude failed query from totals.
- This work item’s outcome and backlog row — record verification and close after success.

**Verify:** focused Go tests, full `make check` from `librarian/src`, source-layout checks if applicable, `git diff --check`, HAWP work validation, then push and read back PR head/checks/threads before resolving comments.

### Outcome

Implemented and pushed in `c6efee6b9ddf0d6583035a43416bbac503ea1dd4`. Token truncation now returns the original short input before any potentially overflowing multiplication and returns empty content for non-positive budgets. The SQLite preflight applies symlink-ancestor and hard-link checks to the main database, WAL, SHM, and rollback journal. The benchmark excludes the failed query from all success totals. The three corresponding PR review threads are resolved.

### Verification

- [x] Context package tests, including max-int/zero/negative budgets: **PASS** (`go test ./internal/application/context`).
- [x] SQLite filesystem tests, including hard-linked and symlinked main/WAL/SHM/journal paths: **PASS** (`go test ./tests/infrastructure/filesystem`).
- [x] Full librarian vet, test suite, and build: **PASS** (`make check` from `librarian/src`).
- [x] Source-layout module tests: **PASS** (`go test ./...` from `scripts/source-layout`).
- [x] HAWP work validation: **PASS**, with one existing repository-wide verification-clarity warning and zero issues.
- [x] Formatting and diff hygiene: **PASS** (`gofmt -l` empty; `git diff --check` clean).
- [x] GitHub `quality` and `validate-generated`: **PASS** on `c6efee6b`.
- [x] PR review threads: all three findings resolved; 78 total threads, zero unresolved.
- [x] PR remains open and unmerged; no human approval is recorded.

### Close Checklist

- [x] Outcome and verification recorded.
- [x] Investigation findings and plan archived under `closed/2026/09/29/0fdd9893/plan.md`.
- [x] Backlog row moved to Recently Closed.
- [x] Staged-path and diff hygiene checks passed before commit.

**Status:** implemented, verified, pushed, review threads resolved, closed.
