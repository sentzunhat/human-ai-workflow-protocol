# PR 41 temporary-file destination preflight

**UUID:** `b3214b76-705f-40b4-9d6f-c31ce8465b35`
**Type:** bug
**Reported:** 2026-09-29

## Input and investigation

Inspect PR 41, address newest Copilot findings, validate, and do not merge.
The three findings in email `1a0ef59e8c6f1138` are already implemented and
archived in `0fdd9893`. Current head `4061c411` contains the overflow guard,
all SQLite sidecar checks, corrected totals, and their regressions.
New review comments `4140079586` and `4140079620` identify missing destination
revalidation between directory creation/fetch and temporary-file creation.

## Plan

Revalidate full destination ancestry immediately before creating temporary
files for member extraction, whole-archive extraction, and verified downloads.
Add regression coverage for a parent replaced during the fetch callback and
for member extraction through a symlinked parent. User authorized fixes and
validation; merge remains prohibited.

## Outcome

Implemented the requested revalidation. The download regression verifies that
substitution during Fetch is rejected before any content reaches the outside
directory. Existing initial and pre-rename checks remain in place.
These are path-based checks; they do not eliminate every concurrent filesystem
race as a directory-handle-based design would.

## Verification

- PASS: Go 1.26.4 `make check` (vet, full tests, static build).
- PASS: `make install`.
- PASS: source-layout module `go test ./...`.
- PASS: `hawp distribution sync`, zero generated drift.
- PASS: `hawp check`, including 143 Markdown link checks.
- PASS: `hawp work validate`, zero issues, one existing verification-clarity warning.
- PASS: formatting and branch diff hygiene.

## Close Checklist

- [x] Findings inspected against current branch, avoiding duplicate fixes.
- [x] Implementation and regression coverage added.
- [x] Local validation completed.
- [x] GitHub Quality and Validate Distribution Generated: PASS on `646529bd`.
- [x] No merge performed or approval claimed.
