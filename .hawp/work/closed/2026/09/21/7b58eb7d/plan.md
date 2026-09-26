# Preflight GitHub MCP configuration before multi-provider writes

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `7b58eb7d-816e-4481-a89b-45fc447bc25e`
**Type:** bug
**Reported:** 2026-09-21

---

## Input (verbatim)

> PR #41 review finding: Configure skips the GitHub/Copilot .vscode/mcp.json preflight, allowing earlier selected providers to be written before an invalid GitHub config fails. Add the missing preflight and a regression proving no partial writes.

## Intake Summary

`hawp mcp configure` promises to preflight every selected provider before
writing. GitHub/Copilot configuration is not included in that preflight, so an
invalid `.vscode/mcp.json` can be discovered only after an earlier selected
provider has already been changed.

## Current Context

The affected runtime path is
`librarian/src/internal/platform/mcp/configure/configure.go`. The existing
preflight supports Claude, Cursor, and Codex; GitHub/Copilot writes use
`.vscode/mcp.json` with the `servers` root key.

## Initial Analysis

**Directly verified:**

- `Configure` calls `writeProviderConfigs` after its preflight loop.
- The loop's `default` skips `github`, while `writeProviderConfigs` later
  creates or merges `.vscode/mcp.json` using the `servers` root key.
- With Claude selected first, a malformed or non-regular GitHub/Copilot config
  can fail after `.mcp.json` and `.gitignore` have been changed.

**Inferred (not yet proven):**

- The same missing preflight can affect any provider ordered before GitHub.

**Likely scope:**

- `librarian/src/internal/platform/mcp/configure/configure.go`
- `librarian/src/internal/platform/mcp/configure/configure_test.go`

## Risk + Review Gate

**Risk:** medium
**Gate:** explicit user authorization received to implement the reviewed fix

## Plan

### Root cause

The preflight switch does not map the `github` provider to its workspace
configuration file and JSON root key. The subsequent write phase therefore
performs validation that should have happened before any selected-provider
write.

### Options considered

1. Add GitHub/Copilot to the existing preflight switch and validate with
   `mergeServerJSON(..., "servers")`.
2. Make provider writes transactional with rollback.

### Recommended fix

Use option 1. It restores the stated all-provider preflight contract without
changing the intentionally non-transactional behavior for filesystem failures
that occur after successful validation.

### Verification target

Add a regression test with `claude` followed by an invalid `.vscode/mcp.json`.
Assert that `Configure` fails and leaves `.mcp.json` and `.gitignore` absent.
Run focused MCP tests, full Go tests, vet, build, HAWP checks, generated
distribution validation, source-layout tests, and `git diff --check`.

## Backlog + Plan Link

**Status now:** in-progress
**Plan file:** work/active/7b58eb7d/plan.md

## Next Step

- [x] Investigation recorded above
- [x] Plan written
- [x] User authorized implementation
- [x] Implement and verify

## Outcome

Added GitHub/Copilot `.vscode/mcp.json` to the selected-provider preflight,
using the same `servers` JSON shape as the writer. A malformed non-stdio HAWP
server now stops the command before Claude or any other earlier provider is
written. Added a multi-provider regression test that proves `.mcp.json` and
`.gitignore` remain absent after this refusal.

## Verification

Directly verified:

- Focused MCP configuration tests passed, including
  `TestConfigurePreflightsGitHubBeforeOtherWrites`.
- `make check` passed from `librarian/src` (`go vet`, full `go test`, build).
- `go run ./cmd/hawp check` passed.
- `go run ./cmd/hawp work validate` passed with zero issues and one existing
  verification-clarity warning.
- `go run ./cmd/hawp distribution validate` passed.
- `go test ./...` passed from `scripts/source-layout`.
- `git diff --check origin/main` passed after the branch whitespace cleanup.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Plan file moved to closed/YYYY/MM/DD
- [x] BACKLOG.md updated
- [x] Status report not needed; plan carries the bounded evidence
