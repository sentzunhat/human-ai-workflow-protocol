# Harden tar kit-bundle extraction against path traversal

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `c5a48d9a-a553-454a-85ca-2e51cf7d632a`
**Type:** bug
**Reported:** 2026-09-20

---

## Input (verbatim)

> Harden tar kit-bundle extraction against path traversal

## Intake Summary

The kit-bundle update path extracts a downloaded tar.gz before copying kit and
provider files. Archive member names are currently joined directly to the
temporary extraction root.

## Current Context

`librarian/src/internal/application/kitsync/kitsync.go:63` calls
`archive.ExtractAll`. `librarian/src/internal/infrastructure/archive/targz.go`
uses `filepath.Join(destDir, header.Name)` without rejecting absolute or
parent-traversal names. A crafted bundle can therefore write outside the
temporary extraction directory during `hawp update sync` or `hawp update`.

## Initial Analysis

**Directly verified:**

- `archive.ExtractAll` handles regular files and directories from tar headers.
- No member-name containment check or traversal regression test exists.
- The caller downloads the release bundle and invokes extraction before any
  manifest or provider validation.

**Inferred (not yet proven):**

- A malicious release asset could overwrite user-writable files outside the
  temporary extraction root.

**Likely scope:**

- `librarian/src/internal/infrastructure/archive/targz.go`
- `librarian/src/internal/infrastructure/archive/extract_test.go`
- `librarian/src/internal/application/kitsync/kitsync.go` (caller only)

## Risk + Review Gate

**Risk:** high
**Gate:** explicit user authorization to implement security fixes

## Plan

### Root cause

Archive member names are treated as trusted filesystem paths.

### Options considered

1. Reject absolute names and verify `filepath.Rel(cleanDest, target)` stays
   within the extraction root. This preserves the existing tar layout.
2. Extract only an allowlisted `kit/` and `providers/` prefix. This is more
   restrictive but couples the generic archive helper to kitsync packaging.

### Recommended fix

Use option 1 in the generic archive helper, reject unsafe member names before
creating directories/files, and add traversal plus absolute-path regression
tests.

### Verification target

Run archive tests, full Go tests, vet, HAWP checks, work validation, and diff
checks. Confirm traversal attempts return errors and do not create outside-root
files.

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/closed/2026/09/20/c5a48d9a/plan.md

## Outcome

Hardened tar.gz extraction with root containment checks. Absolute and
parent-traversal member names are rejected before directories or files are
created. Added regression tests proving traversal and absolute members do not
write outside the extraction root.

## Verification

Directly verified:

- `go test ./internal/infrastructure/archive` passed, including both new
  unsafe-member tests.
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
