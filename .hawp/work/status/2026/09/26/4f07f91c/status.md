---

## Intent

Review the current PR #41 head, fix confirmed correctness/security issues, commit and push the fixes, and refresh PR metadata when authenticated GitHub write access is available.

## Current State

The branch is clean and pushed at `f6d519e6`, two commits beyond the prior remote PR head: `309c7c47` (MCP/UTF-8 correctness) and `f6d519e6` (remaining symlinked write paths). PR #41 remains open and its remote title/body have not yet been updated.

## What Was Inspected

- `.hawp/kit/start-here.md` and `.hawp/kit/usage/status-report.md`.
- The immutable PR range `origin/main...309c7c47`, including its 285 changed source/review items and high-risk filesystem, archive, download, distribution, CLI, MCP, and migration paths.
- Current PR #41 metadata and review conversation through the public GitHub page.
- The current local diff and post-fix commit `f6d519e6`.

## What Changed

- Added fail-closed symlink checks before init-time HAWP home provisioning and runtime-folder creation.
- Guarded verified downloads and single-member archive extraction before and after writes.
- Made README generation preflight both roots, use `Lstat`, and write through the guarded atomic writer.
- Guarded generated distribution writes with the shared symlink-safe path and atomic-write helpers.
- Added regression coverage for symlinked provisioning/runtime/project paths.

## What Was Directly Verified

- Focused tests passed for download, archive, filesystem, provisioning, distribution, and new symlink regressions.
- `make check` passed: full Go tests, `go vet`, and static build.
- Source-layout tests passed.
- `hawp check`, `hawp work validate`, and `hawp distribution validate` passed; work validation reported one existing warning and no issues.
- Generated shell syntax checks and `git diff --check` passed.
- GitKraken push succeeded; `git ls-remote` confirmed `origin/feature/v0.0.24` at `f6d519e6`.

## What Remains Unproven

- Fresh remote CI and review readback for `f6d519e6` have not yet been observed.
- The Codex Security scan identified two reportable pre-fix CWE-59 findings, but its final completion call failed because the canonical scan manifest was unavailable after a schema-validation error; no completed sealed report is claimed.
- PR #41 title and description remain at their prior values because the available GitHub CLI/API credential is invalid, the browser is signed out, and GitKraken’s available PR surface does not expose metadata editing.
- A full manual file-by-file review of the 573-file remote PR remains distinct from targeted inspection and local verification.

## Constraints

No merge, approval, review resolution, or force-push was performed. The existing clean branch history was preserved; only the two local commits were pushed normally.

## Help Wanted

Use an authenticated GitHub write session to update PR #41 metadata and request fresh CI/review readback. Treat the prior Copilot `Findings: None` summaries as time-specific and not equivalent to a complete manual review.

## Suggested Next Step

Authenticate GitHub CLI or browser access, update PR #41 title/body to describe the MCP/UTF-8 and symlink-boundary hardening plus exact verification, then read back the remote head, checks, review state, and merge state before any human approval decision.
uuid: 4f07f91c-11a7-4978-92ac-2092deb1f957
title: "PR 41 deeper review and symlink follow-up"
type: status
date: 2026-09-26
---
