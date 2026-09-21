# Validate links in nested canonical work plans

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

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `a78f1048-7817-4528-a8e7-0023a1ba8b04`
**Type:** bug
**Reported:** 2026-09-23

## Input (verbatim)

> Validate links in nested canonical work plans

## Intake Summary

The Copilot review identifies a coverage gap in dead-link validation. Canonical new plans
live at `active/{uuid}/plan.md` and `parked/{uuid}/plan.md`, but the validator only
inspects Markdown files directly under each directory. Broken links inside a
folder-per-item plan are therefore omitted.

## Current Context

The validator must scan canonical nested plans while retaining flat legacy files and
the existing exclusion of closed/evidence/status archives. The fix belongs at the
shared validation boundary and needs an in-memory filesystem regression.

## Initial Analysis

**Directly verified:**

- `librarian/src/internal/domain/work/validation/deadlinks.go:18-39` calls `ReadDir`
  once per scope and skips directory entries.
- `librarian/src/internal/domain/work/validations_test.go:383-400` covers flat active
  files but no nested canonical plan.

**Inferred (not yet proven):**

- Recursively walking only `active/` and `parked/` and selecting Markdown files will
  include canonical nested plans and flat legacy files without scanning archives.

**Likely scope:**

- `librarian/src/internal/domain/work/validation/deadlinks.go`, validation tests, and
  any injected `ReadDir` behavior required for recursive traversal.

## Root Cause

The file-discovery loop treats each work scope as flat and does not descend into the
UUID directory that contains `plan.md`.

## Options Considered

1. Add a second hard-coded lookup for `<entry>/plan.md`. This covers only the current
   canonical name and misses future nested supporting Markdown files.
2. Recursively walk the selected work scopes and collect Markdown files while keeping
   archive scopes excluded. This is the recommended fix and preserves legacy files.

## Recommended Fix

Introduce a small recursive collector over `active` and `parked`, append regular
Markdown files, and retain the existing link extraction/containment behavior. Add a
nested-plan broken-link regression and preserve the existing flat/fenced-link test.

## Risk + Review Gate

**Risk:** medium (shared validation traversal)
**Gate:** user authorized implementation in the request; no merge authorization.

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/active/a78f1048/plan.md

## Verification

- Run the focused dead-link regression for `active/<uuid>/plan.md`.
- Confirm flat legacy plans, fenced links, and parked plans remain covered.
- Run all Go tests and HAWP workflow validation.

## Next Step

- [x] Investigation recorded above
- [x] Write or update the plan file
- [x] Create the corresponding fix work item and implement sequentially
## Reconciliation Outcome

Done. Nested canonical work-plan links are validated by the current work
validation path; repository tests and HAWP validation pass.
