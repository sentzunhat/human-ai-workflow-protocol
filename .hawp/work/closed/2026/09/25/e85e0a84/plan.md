# Reconcile UUID-folder work items during close

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
**UUID:** `e85e0a84-1c71-4b1b-87b8-47543621d835`
**Type:** bug
**Reported:** 2026-09-23

## Input (verbatim)

> Reconcile UUID-folder work items during close

## Intake Summary

The Copilot review identifies a mismatch between the canonical UUID-folder work layout
and the generated install/update reconciliation helper. A closed link such as
`closed/YYYY/MM/DD/<uuid>/plan.md` is reduced to the basename `plan.md`, so the
helper looks for `active/plan.md` instead of `active/<uuid>/plan.md`.

## Current Context

The source of truth is `distribution/sources/update/script-core.md`; generated
provider variants must be regenerated after the source fix. The change must retain
flat legacy-file and ID-fallback behavior and must not broaden accepted closed paths
beyond existing containment checks.

## Initial Analysis

**Directly verified:**

- `distribution/sources/update/script-core.md:80-136` derives `plan_name` with
  `basename` and constructs `.hawp/work/active/$plan_name`.
- Generated update variants contain the same implementation.
- Canonical active work uses `.hawp/work/active/{uuid}/plan.md`.

**Inferred (not yet proven):**

- The source path's parent directory is the UUID folder for canonical plans; retaining
  that folder can preserve canonical reconciliation while flat paths remain supported.

**Likely scope:**

- `distribution/sources/update/script-core.md`, generated install/update artifacts,
  and focused shell/regeneration validation.

## Root Cause

Reconciliation models a closed plan as a flat file and discards the folder component
that identifies a UUID-native work item.

## Options Considered

1. Add a special case for literal `plan.md` and infer the UUID from the destination.
   This is ambiguous and can select the wrong active item.
2. Derive the source relative to the closed path: for UUID-folder paths, retain the
   final directory and `plan.md`; for flat legacy paths, retain the basename. This is
   the recommended compatible fix.

## Recommended Fix

Teach the shared reconciliation function to recognize the UUID-folder closed layout,
map it to `.hawp/work/active/<uuid>/plan.md`, and preserve the existing flat-file
fallback. Regenerate all provider variants and add fixture coverage where available.

## Risk + Review Gate

**Risk:** high (generated install/update behavior and filesystem mutation)
**Gate:** user authorized implementation in the request; no merge authorization.

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/active/e85e0a84/plan.md

## Verification

- Exercise UUID-folder and flat legacy closed links against the shared shell logic.
- Run `hawp distribution sync` and `hawp distribution validate`.
- Confirm all generated provider variants agree with the source.

## Next Step

- [x] Investigation recorded above
- [x] Write or update the plan file
- [x] Create the corresponding fix work item and implement sequentially
## Reconciliation Outcome

Done. Close/reconciliation supports canonical UUID-folder work records and
the repository-wide Go gates pass.
