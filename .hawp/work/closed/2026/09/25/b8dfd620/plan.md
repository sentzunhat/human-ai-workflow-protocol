# Reject non-regular existing installation destinations

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

**UUID:** `b8dfd620-e33d-4ce0-a720-715a30b4f117`
**Type:** bug
**Reported:** 2026-09-23

## Intake Summary

Copilot reports that the `librarian/src/Makefile` install target rejects symlinked `.hawp` paths but accepts an existing directory at `.hawp/bin/hawp`. `mv -f` can then place the temporary binary inside that directory while the advertised executable destination remains unusable.

## Current Context

The installer uses a temporary sibling and performs a pre-rename guard, but both guards currently check only `-L`. The destination is a repository-owned executable path and must be either absent or a regular file before replacement.

## Initial Analysis

**Directly verified:** `librarian/src/Makefile` checks symlink status before and after `cp`, then runs `mv -f` without checking destination type.

**Inferred:** An existing directory or other non-regular node can alter `mv` semantics and make installation report success without producing the expected executable.

**Likely scope:** Makefile install recipe and its executable-destination regression coverage.

## Root Cause

The safety invariant is expressed as “not a symlink” rather than “absent or regular file.”

## Options

1. Add shell tests using `[ -e ] && [ ! -f ]` at both destination checks. This is narrow and covers directories and special nodes.
2. Replace the shell install flow with a Go installer helper. This is broader than the finding and increases scope.

## Recommended Fix

Use a small shell predicate that rejects an existing destination unless it is a regular file, while retaining the symlink checks for `.hawp` and `.hawp/bin`. Add a regression test or static recipe assertion for directory destinations.

## Risk + Review Gate

**Risk:** high — installation path mutation.
**Gate:** user explicitly requested implementation and sequential remediation.

## Outcome

Done. The installer rejects an existing executable destination unless it is a
regular file, while retaining explicit symlink checks for `.hawp` and
`.hawp/bin`. The linked Makefile fix is verified by distribution-focused
tests.
