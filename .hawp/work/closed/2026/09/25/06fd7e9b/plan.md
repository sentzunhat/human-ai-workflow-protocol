# Handle all-zero base SHA on branch creation in quality workflow

## Outcome

The PR-review finding was implemented and reconciled with the current branch.

## Verification

The focused regression coverage and repository-wide Go, HAWP, distribution,
formatting, and diff-hygiene checks passed before close.

## Close Checklist

- [x] Outcome recorded.
- [x] Verification evidence recorded or referenced.
- [x] Backlog row removed from active coordination.
- [x] Plan archived under the close date.

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `06fd7e9b-186d-4ea4-be3d-311ba7a3b4ca`
**Type:** bug
**Reported:** 2026-09-23

## Input (verbatim)

> Handle all-zero base SHA on branch creation in quality workflow

## Intake Summary

The Copilot review identifies a branch-creation failure in the quality workflow. On a
push that creates a branch, `github.event.before` is Git's all-zero object ID, so the
branch-diff hygiene step can fail before the remaining quality checks execute.

## Current Context

This is PR #41's CI contract. The fix must preserve pull-request merge-base behavior
and ordinary push behavior while making a zero-before push deterministic. The user has
authorized execution of this high-severity review item; merge remains unauthorized.

## Initial Analysis

**Directly verified:**

- `.github/workflows/quality.yml` uses `github.event.before` for non-PR pushes.
- The same step invokes `git diff --check "${base}" "${{ github.sha }}"`.
- The supplied review text confirms the all-zero SHA is the triggering input.

**Inferred (not yet proven):**

- Falling back to the first root commit is sufficient for a non-empty repository;
  skipping the diff would reduce hygiene coverage and is less desirable.

**Likely scope:**

- `.github/workflows/quality.yml` and workflow-level inspection of the resulting shell
  branch.

## Root Cause

The non-PR branch selects an event sentinel as if it were a real Git object. Git
cannot resolve that sentinel as a diff base.

## Options Considered

1. Skip branch-diff hygiene when the before SHA is all zero. This is simple but leaves
   the first push without whitespace validation.
2. Replace the all-zero value with `git rev-list --max-parents=0 "${{ github.sha }}" |
   head -n1`. This retains the check and is the recommended fix.

## Recommended Fix

Add an explicit all-zero comparison after selecting `base` and before `git diff`.
Keep the PR merge-base path unchanged.

## Risk + Review Gate

**Risk:** high (CI/security-adjacent boundary)
**Gate:** user authorized implementation in the request; no merge authorization.

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/active/06fd7e9b/plan.md

## Verification

- Inspect the rendered workflow shell for PR, ordinary push, and all-zero push paths.
- Run the repository's workflow/distribution validation and diff hygiene checks.
- Confirm no unrelated workflow behavior changes.

## Next Step

- [x] Investigation recorded above
- [x] Write or update the plan file
- [x] Create the corresponding fix work item and implement sequentially

## Outcome

Done. The quality workflow replaces GitHub's all-zero push predecessor with
the first root commit before running diff hygiene while preserving PR and
ordinary push paths.
