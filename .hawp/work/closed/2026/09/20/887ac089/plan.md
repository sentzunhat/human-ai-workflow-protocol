# Fix model pull trailing --onnx-file parsing

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `887ac089-83f8-4f08-b38e-a0240400c828`
**Type:** bug
**Reported:** 2026-09-20

---

## Input (verbatim)

> Fix model pull trailing --onnx-file parsing

## Intake Summary

The model-pull command documents options after its single repository
positional argument. The current implementation parses the flag set, then
parses a sliced remainder a second time to compensate for Go flag parsing
stopping at the first positional argument.

## Current Context

`librarian/src/internal/platform/cli/model/pull/args.go` calls
`FlagSet.Parse` twice. The existing unit test covers one repo-first trailing
option, but the parser has two independent parse states and is brittle for
interspersed/trailing combinations. The fix should make option ordering
deterministic and keep the documented form covered.

## Initial Analysis

**Directly verified:**

- The documented repo-first `--onnx-file` test currently passes in isolation,
  so the Copilot report is not reproduced by that single case.
- The parser nevertheless relies on a non-standard second parse of the same
  `FlagSet`; this is the maintainability and edge-case risk to remove.

**Inferred (not yet proven):**

- A trailing option combination can be accepted or rejected based on the
  first parse's stopping point rather than one consistent argument contract.

**Likely scope:**

- `librarian/src/internal/platform/cli/model/pull/args.go`
- `librarian/src/internal/platform/cli/model/pull/args_test.go`

## Risk + Review Gate

**Risk:** medium
**Gate:** explicit user authorization to fix the reported parser issue

## Plan

### Root cause

Go's standard `flag` parser stops at the first positional argument, while this
command intentionally permits options before and after the repository. The
current double-parse workaround obscures the contract.

### Options considered

1. Normalize/reorder arguments before one standard `flag` parse.
2. Use one explicit token pass for the two supported options and the one
   positional repository, rejecting unknown options and extra positionals.

### Recommended fix

Use option 2: make the accepted command grammar explicit, preserve both
repo-first and option-first forms, and expand tests for trailing and mixed
options. This avoids repeated `FlagSet` state and keeps validation local.

### Verification target

Run focused parser tests, full Go tests, vet, HAWP checks, work validation, and
diff checks.

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/closed/2026/09/20/887ac089/plan.md

## Outcome

Replaced the double `flag.FlagSet.Parse` workaround with one deterministic
argument pass that accepts options before or after the repository, including
the documented trailing `--onnx-file` form. Added regression coverage for
trailing separate-value and equals-value options.

## Verification

Directly verified:

- Focused model-pull parser tests passed.
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
