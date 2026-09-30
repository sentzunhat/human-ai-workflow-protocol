# Preflight MCP prerequisites in the public provider-config wrapper

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

**UUID:** `37703223-232a-4d1d-813b-7905b65f8335`
**Type:** bug
**Reported:** 2026-09-23

## Intake Summary

Copilot reports that `hawp init --provider ...` calls the exported `WriteProviderConfigs` wrapper directly. That wrapper validates provider names but bypasses `Configure`'s prerequisite and existing-config preflight, so a symlinked `.hawp/bin/hawp` or `.hawp/work/BACKLOG.md` can be accepted and multi-provider writes can begin before a later configuration failure is found.

## Current Context

`librarian/src/internal/platform/cli/init/command.go` calls `WriteProviderConfigs`. `librarian/src/internal/platform/mcp/configure/configure.go` performs the required prerequisite, symlink-ancestor, and selected-destination validation before calling the private `writeProviderConfigs`. The private writer is intentionally non-transactional after preflight.

## Initial Analysis

**Directly verified:** `WriteProviderConfigs` only expands provider names before calling the private writer. `Configure` contains the prerequisite and all selected-provider destination preflight. `hawp init` uses the public wrapper, not `Configure`.

**Inferred:** the direct init route can write `.mcp.json` or another earlier provider configuration before encountering a later unsafe configuration, and it does not reject redirected prerequisites.

**Likely scope:** `librarian/src/internal/platform/mcp/configure/config.go`, `librarian/src/internal/platform/mcp/configure/configure.go`, their focused tests, and legitimate direct-wrapper tests in `librarian/src/internal/platform/cli/init_test.go`.

## Root Cause

The safety preflight is attached to one public entry point instead of the shared exported wrapper that reaches the non-transactional provider-file writer.

## Options

1. Make `WriteProviderConfigs` expand, preflight, then invoke the existing private writer; make `Configure` route through that public boundary. This covers both command paths while preserving `writeProviderConfigs` for internally preflighted calls.
2. Change only `hawp init` to call `Configure`. This leaves the exported wrapper unsafe for other callers and duplicates public configuration semantics.

## Recommended Fix

Choose option 1. Extract the existing `Configure` preflight into a private helper, call it from `WriteProviderConfigs` after provider expansion, and make `Configure` delegate to that wrapper. Add direct-wrapper tests proving rejected symlinked prerequisites and invalid later provider configurations leave no provider config or `.gitignore` writes; preserve a valid direct-wrapper configuration path with regular prerequisites.

## Risk + Review Gate

**Risk:** high — provider configuration can trust redirected executable/backlog prerequisites and produce partial repository writes.
**Gate:** user explicitly authorized sequential remediation in this PR loop.

## Outcome

Done. The exported provider wrapper now performs the shared preflight before
the non-transactional writer, covering direct wrapper callers and the
`hawp init` route. Focused direct-wrapper tests prove later-provider failures
leave earlier provider files and `.gitignore` unwritten.

See also the completed implementation item `73bd501c`.
