# Use Lstat for MCP configuration prerequisites

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

**UUID:** `fc6a2c80-c4bd-47be-80cb-3b27aa62a209`
**Type:** fix
**Reported:** 2026-09-23
**Source plan:** `04c18c70`

## Context

MCP configuration already fails closed for non-regular existing configuration files, but its binary and backlog prerequisite loop follows symlinks through `os.Stat`.

## Fix Work

Change prerequisite inspection to `os.Lstat`, retaining the existing regular-file predicate. Add regressions for both symlinkable prerequisite locations and prove no provider config is written when preflight rejects them.

## Verification

Run focused MCP configuration tests, then `go test ./...`, `go vet ./...`, HAWP validation, formatting, and diff-hygiene checks.

## Outcome

Done. `preflightProviderConfigs` uses `os.Lstat` and rejects non-regular
prerequisites before provider configuration writes. Focused MCP configuration
tests prove symlinked prerequisites are rejected and no configuration is
written.
