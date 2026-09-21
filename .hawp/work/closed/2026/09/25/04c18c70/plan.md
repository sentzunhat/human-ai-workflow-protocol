# Reject symlinked MCP configuration prerequisites

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

**UUID:** `04c18c70-0bcc-42e3-82cb-9ffc65a786f9`
**Type:** bug
**Reported:** 2026-09-23

## Intake Summary

Copilot reports that MCP configuration preflight calls `os.Stat` for `.hawp/bin/hawp` and `.hawp/work/BACKLOG.md`. `Stat` follows symlinks, allowing a redirected executable or backlog prerequisite to pass the regular-file check.

## Current Context

`librarian/src/internal/platform/mcp/configure/configure.go` already uses `os.Lstat` for existing provider configuration files, and the surrounding writer paths use repository-aware symlink validation.

## Initial Analysis

**Directly verified:** prerequisite validation uses `os.Stat`; provider configuration validation in the same function uses `os.Lstat` and rejects non-regular nodes.

**Inferred:** a symlink to a regular file is accepted as a prerequisite despite configuration hardening requiring direct repository-owned regular files.

**Likely scope:** MCP configuration preflight and focused prerequisite-symlink tests.

## Root Cause

The preflight’s regular-file assertion evaluates the target instead of the prerequisite directory entry itself.

## Options

1. Replace `os.Stat` with `os.Lstat`, keeping the current regular-file test and error behavior. This is local and consistent with adjacent checks.
2. Add a separate helper with repository-root traversal checks. This is broader than the reported direct-entry symlink gap and duplicates existing guards.

## Recommended Fix

Use `os.Lstat` for both prerequisites and add a table-driven regression for a symlinked binary and a symlinked `BACKLOG.md`, asserting configuration writes do not occur.

## Risk + Review Gate

**Risk:** high — generated MCP configuration can trust a redirected executable.
**Gate:** user explicitly authorized sequential remediation in this PR loop.

## Outcome

Done. The linked fix uses `os.Lstat` for both prerequisites and rejects
symlinked ancestors before inspection. Focused MCP configuration tests cover
redirected binary, backlog, and `.hawp` root prerequisites.

See also the completed implementation item `fc6a2c80`.
