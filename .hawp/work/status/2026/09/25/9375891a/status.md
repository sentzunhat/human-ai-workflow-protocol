---
uuid: 9375891a-0d30-4321-95d0-1701ac51b0fe
title: "PR 41 hard-link-safe managed configuration writes"
type: status
date: 2026-09-25
---

### Status Report

#### Intent

Close Copilot finding `Prevent hard links from overwriting managed configuration files` in PR #41 and inspect the adjacent managed MCP configuration write paths for the same class of issue.

#### Current State

Fixed in the working tree. Managed JSON, Codex TOML, and `.gitignore` writes now use atomic sibling replacement, so a hard-linked destination is replaced rather than truncating the shared inode.

#### What Was Inspected

- `librarian/src/internal/platform/mcp/configure/configure.go`
- `librarian/src/internal/platform/mcp/configure/config.go`
- `librarian/src/internal/platform/mcp/configure/config_json.go`
- `librarian/src/internal/platform/mcp/configure/config_codex.go`
- Existing MCP configuration tests and the shared `filesystem.AtomicWriteFile` implementation.

#### What Changed

Added the shared `managedFilePerm` validation/permission-preservation helper and routed the three production managed-file writers through `filesystem.AtomicWriteFile`. Added a regression table covering hard-linked `.mcp.json`, `.codex/config.toml`, and `.gitignore` destinations.

#### What Was Directly Verified

- The hard-link regression passes: the managed destination receives the new content while the other hard-linked pathname remains unchanged.
- `go test ./...` passed.
- `go vet ./...` passed.
- `go build ./cmd/hawp` passed.
- `GOOS=windows GOARCH=amd64 go build ./...` passed.
- `go test -race ./internal/platform/mcp/configure` passed.
- Focused hard-link, prerequisite, and provider-preflight tests passed uncached.
- `go run ./cmd/hawp check` passed all three validations.
- `git diff --check` passed.

#### What Remains Unproven

Windows runtime behavior was not executed on Windows, but the complete Go tree cross-built successfully for `windows/amd64`. The implementation uses the repository's existing cross-platform `os.CreateTemp` plus `os.Rename` atomic replacement abstraction; runtime tests and race detection passed in the current environment.

#### Constraints

Unrelated worktree changes were preserved. The review was scoped to the PR #41 managed configuration write boundary and directly related callers; no speculative changes were made to unrelated `os.WriteFile` uses.

#### Help Wanted

External PR review and CI should confirm the expected Windows rename semantics and the final Copilot assessment.

#### Suggested Next Step

Review the four changed source/test files, push the patch to PR #41, and request a fresh Copilot review after the new commit is available.
