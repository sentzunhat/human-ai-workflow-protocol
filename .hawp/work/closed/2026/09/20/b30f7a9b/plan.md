# Validate kitsync manifest paths stay within repository and bundle roots

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `b30f7a9b-102d-4089-9cce-6ed978760a7a`
**Type:** bug
**Reported:** 2026-09-20

---

## Input (verbatim)

> Validate kitsync manifest paths stay within repository and bundle roots

## Intake Summary

The kitsync domain applies source and destination paths from a bundled YAML
manifest. Those paths are currently joined to the bundle and repository roots
without containment validation.

## Current Context

`librarian/src/internal/domain/kitsync/apply.go` constructs source paths from
`provider.Source`/`rule.From` and destination paths from `rule.Dest`. A
malicious manifest can use `..` or an absolute path to escape the intended
bundle or repository root.

## Initial Analysis

**Directly verified:**

- `ApplyProviderUpdate` and `ApplyProviderInstall` use `filepath.Join` for
  manifest-controlled paths.
- No source-root or destination-root containment helper currently guards these
  joins.
- The manifest is loaded from the downloaded release bundle.

**Inferred (not yet proven):**

- A malicious bundled manifest could overwrite files outside the target repo
  or read source files outside the extracted bundle.

**Likely scope:**

- `librarian/src/internal/domain/kitsync/apply.go`
- `librarian/src/internal/domain/kitsync/kitsync_test.go`
- `librarian/src/internal/application/kitsync/kitsync.go` (manifest loader)

## Risk + Review Gate

**Risk:** high
**Gate:** explicit user authorization to implement security fixes

## Plan

### Root cause

Manifest path fields are treated as trusted relative paths without checking
that their cleaned results remain under the bundle or repository roots.

### Options considered

1. Add a shared root-containment helper and validate provider sources, rule
   sources, and rule destinations before any filesystem operation.
2. Validate only the checked-in manifest during release creation. This does not
   protect consumers from a tampered or compromised downloaded bundle.

### Recommended fix

Use option 1. Require relative, root-contained paths for all manifest-controlled
source and destination paths, return descriptive errors, and add malicious
manifest tests for traversal and absolute paths.

### Verification target

Run kitsync tests, full Go tests, vet, HAWP checks, work validation, and diff
checks. Confirm rejected paths cause no writes outside the intended root.

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/closed/2026/09/20/b30f7a9b/plan.md

## Outcome

Added manifest path validation that requires provider sources, rule sources,
and rule destinations to remain within their bundle or repository roots. All
rules are validated before application, preventing partial writes from a later
invalid rule. Added regression tests for source and destination traversal.

## Verification

Directly verified:

- `go test ./internal/domain/kitsync` passed, including both new malicious
  manifest tests.
- `go test ./...` passed from `librarian/src`.
- `go vet ./...` passed from `librarian/src`.
- `go run ./cmd/hawp check` passed with zero issues.
- `go run ./cmd/hawp work validate` passed with zero issues and one existing
  verification-clarity warning.
- `go test ./...` passed from `scripts/source-layout`.
- `git diff --check` passed.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Plan file moved to closed/YYYY/MM/DD/
- [x] BACKLOG.md updated
- [x] Status report written if needed
