---
uuid: f8685e56
title: "Website and MCP PR split checkpoint"
type: status
date: 2026-10-23
---

# 2026-10-22 — HAWP website and Librarian changes split into two PRs;
# website PR closed on 2026-10-23

**Event date:** 2026-10-22
**Archive date:** 2026-10-23
**Work item:** `.hawp/work/active/f8685e56/plan.md`

## Status Report

### Intent

Continue the existing split work item, preserve the decision trail, and give
the next maintainer a verified PR, branch, and search-index state without
duplicating the earlier checkpoint protocol record.

### What Was Inspected

- The existing split plan, path-ownership manifest, backlog row, and status
  report conventions.
- Earlier session checkpoints for branch preparation and publication, plus
  the existing public-safe archival-protocol checkpoint.
- The repository README, current branch/worktree state, and GitHub PR
  metadata and checks.

### Before

- A mixed website and Librarian/MCP change set was on `feature/website`.
- The user approved two separate branches based on `main`: website/Pages on
  `feature/website-improvements`, and MCP/provider/Librarian changes on
  `feature/mcp-providers`.
- Search-index Option A was approved: full-corpus synchronization, stale
  document removal, empty-corpus protection, regression tests, and updated
  search guidance. Rebuilding the persistent local index was not approved or
  performed.

### During

- The two branches were prepared from `origin/main` at
  `178fde6ebcb57f3e8100710c72ae71f06e62fb6d`. Their owned path sets were
  disjoint: 39 website paths and 76 code/provider paths.
- The website branch was committed and pushed at
  `f27855b08f0c11d388a22f1700bc124f751a672c`.
- The code branch contains MCP resources/prompts, provider guidance and
  generated outputs, configuration changes, and the approved search-index
  implementation with regression tests. Its implementation commit is
  `a35dd925a04c397e0a8ce363b3066e7a1edc085b`; subsequent record updates were
  also pushed.
- PR #47, **Add the HAWP website and Pages deployment**, was opened from
  `feature/website-improvements` to `main`.
- PR #46, **Expand HAWP MCP tooling and provider setup**, was opened from
  `feature/mcp-providers` to `main`.

### After

- On 2026-10-23, PR #47 was closed without merging, as requested. Its branch
  remains available at `f27855b08f0c11d388a22f1700bc124f751a672c`. PR #46
  remains open.
- The code worktree is aligned with `feature/mcp-providers`; the website
  worktree remains clean on `feature/website-improvements`.
- The separate root checkout's website-social-card changes remain untouched.

### What Was Directly Verified

- PR #47 is closed, not merged. GitHub reported `mergedAt: null` and retained
  its head branch `feature/website-improvements`.
- PR #46 is open against `main`; GitHub reported it mergeable. Its `quality`
  and `validate-generated` checks passed on head
  `b81ad060f9a61c65899f78c14e79313b6514fd7a` before this checkpoint-only
  update.
- This checkpoint and its split-plan, manifest, and backlog updates were
  committed and pushed on `feature/mcp-providers` as
  `59907e715b4624b2a6fa975d605e120370a9d17f`. Both PR checks passed on that
  pushed head.
- PR #47's website build check passed; its deploy check was skipped as
  expected for a pull request.
- Website `npm run check` and `npm run build` passed. The full Go test suite,
  provider sync, distribution sync, and code-branch HAWP validation passed.
- After this checkpoint and plan update, `git diff --check` passed and
  `cd librarian/src && go run ./cmd/hawp check` passed all three checks with
  0 issues, 0 warnings, and 147 valid Markdown files.
- `hawp/hawp_search` returned no indexed results for the targeted split query.
  Repository files and GitHub records were used as the direct source of
  truth; the empty search result is not evidence that records are absent.
- The current `feature/mcp-providers` worktree was switched from `main` to its
  PR branch after confirming it was clean. The root checkout's separate
  website-social-card changes were left untouched. The
  `feature/website-improvements` worktree remains clean at its pushed head.

### What Remains Unproven

- The persistent local SQLite search index was not rebuilt, and live search
  against a restarted server using the new code remains unverified.
- The actual selected-agent Tools list, VS Code approval behavior, and MCP
  migration-preview effect remain unverified.

### Strategic Impact

The split keeps the website/Pages lifecycle independent of MCP/provider code
and gives each branch a clear review boundary. PR #47 is closed without
merging, so the website work remains available on its branch but is not
proposed for `main` through that PR. PR #46 remains the active review lane.
The version and release policy was not changed.

### Decisions and Constraints

- Keep the website and code PRs separate and based on `main`; do not merge
  either PR without a separate user request.
- Do not delete the pushed website branch when closing its PR.
- Keep search-index live recovery and VS Code client evidence separate from
  code-level tests. Do not claim that local search was repaired until live
  verification is completed.
- This public repository checkpoint excludes personal and unrelated private
  project context.

### Continuation State

**Latest milestone:** Both PRs were opened; the website PR was subsequently
closed unmerged; the MCP/provider PR remains open.

**Next actions:**

1. Review PR #46 and address any requested changes.
2. Keep `feature/website-improvements` available; reopen PR #47 only if
   requested.
3. For live search verification, restart the server with the new code, obtain
   approval before rebuilding the persistent index, and run a targeted query.
4. Obtain direct VS Code evidence for the selected-agent Tools list,
   approval behavior, and migration-preview effect.

**Unresolved:** Live search/index recovery and all three client-side MCP
checks remain unverified. The persistent local index was not rebuilt.

**Dependencies:** Live search proof requires the new server code and an
approved local index rebuild. Client behavior requires access to the VS Code
session/UI.

### Help Wanted

No review request is needed for the checkpoint itself. PR #46 review and the
unverified live/client checks remain open.

### Suggested Next Step

Start with PR #46 review; do not merge or re-open PR #47 without a separate
request.
