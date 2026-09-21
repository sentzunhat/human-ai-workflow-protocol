# HAWP docs/work cleanup and next-item review

**UUID:** `e5fca9c7-ba0a-424e-b0a5-ef725786a99e`
**Type:** release
**Reported:** 2026-09-02
**Status:** `done`

---

## Input

> review and organize docs and files and review commits if we deleted and synced the hawp folder with core folder and so on and so forth open a new patch branch update for these cleanup and for what is next on our items work items

## Intent

Create a `v0.0.24` patch branch that verifies the `v0.0.23` closeout state,
improves README positioning, tightens UUID-folder validation, extends safe
work-normalization cleanup for older repositories, and refreshes the visible
next-work dashboard.

## Scope

- Inspect recent release/closeout commits on `main`
- Compare root `.hawp/kit` with `core/.hawp/kit`
- Compare root `.hawp/work` with `core/.hawp/work` only to classify expected
  scaffold-vs-project-state differences
- Improve the README's open-source builder/founder positioning while keeping
  measured claims tied to repo evidence
- Record practical older-repo compatibility targets for the parked UUID
  canonicalization lane
- Update `.hawp/work` bookkeeping for the current cleanup lane
- Preserve release artifacts and generated distribution outputs unless a real
  drift issue is found

## Direct Evidence So Far

- `main` currently points at `94a66001` (`Align release docs and provider
  distribution guides (#39)`), with `05aa28ce` immediately before it for the
  checked-in binary refresh.
- `diff -qr .hawp/kit core/.hawp/kit` produced no differences.
- `diff -qr .hawp/work core/.hawp/work` showed expected differences: root
  `.hawp/work` contains this repo's parked/closed/evidence/status history,
  while `core/.hawp/work` is the downstream scaffold.
- `.hawp/work/STATUS.md` was stale before this cleanup and still referenced an
  old June closeout.
- The README improvement note supplied on 2026-09-05 recommended positioning
  HAWP around project-owned intent, vendor independence, and the line "Your
  model can change. Your intent shouldn't."
- Candidate older-repo proof targets for `89cf7a85` are `mochila-archive`,
  `tekit`, and `local-print-farm`.
- The cleanup item was initially created in a descriptive active folder, which
  violated the UUID-folder convention even though validation passed.
- The branch was promoted from cleanup to `v0.0.24` patch release scope on
  2026-09-05 because it now includes CLI behavior, changelog, and binary
  updates.

## Work Plan

2026-09-05 downstream report: the user supplied Mochila's successful real-repo
normalization report (5/5 checks, zero issues/warnings, no commits, unrelated
dirty files preserved). This is user-reported evidence, not a fresh local audit.
The 35 unidentified legacy files remain intentionally untouched. Follow-up:
preserve all duplicate copies and artifacts, cross-link unique archived matches,
and expose ambiguous review cases in dry-run output before further rollout.

- [x] Open cleanup branch from current `main`
- [x] Record cleanup lane in `.hawp/work/active/`
- [x] Refresh `.hawp/work/STATUS.md`
- [x] Run HAWP work validation from `librarian/src`
- [x] Summarize next queued work from `.hawp/work/BACKLOG.md`
- [x] Add README positioning improvement to this cleanup scope
- [x] Re-run validation after README edits
- [x] Move this cleanup item into UUID-folder layout:
      `.hawp/work/active/e5fca9c7/plan.md`
- [x] Tighten validation so new active non-UUID identifiers fail while legacy
      `TASK-NNN` and numeric rows remain tolerated
- [x] Test UUID-folder validation against `mochila-archive` in read-only mode
- [x] Advance CLI version metadata to `0.0.24`
- [x] Add changelog entry for `0.0.24`
- [x] Build checked-in `.hawp/bin/hawp` as `0.0.24`
- [x] Prove Mochila and Tekit normalization on temporary copies without
      committing to either repo
- [x] Add Mochila digital-agent instruction artifact
- [x] Decide whether to close this cleanup item in the same branch

## Next Work Queue

2026-09-07 scope update: the
[single-executable compatibility gate](../../closed/2026/09/07/0cb0f9b0/plan.md) is closed locally. Continue
CLI validation, persisted UUID/status correctness and storage boundaries. The
user authorized preservation-first routing of unmatched work to a proposed
`.hawp/work/unknown/` folder, not deletion or blind overwrites. Collision-safe
moves, provenance, link repair, indexing rules, and temporary-copy proof remain
planned; no unmatched records have been moved by this continuation.

Current backlog status after `v0.0.23`:

- `89cf7a85` remains the main near-term follow-up: legacy work-item UUID
  canonicalization, older-repo safety, and archive-policy boundaries. Candidate
  targets: `mochila-archive`, `tekit`, or `local-print-farm`.
- `release-benchmark-backfill` is a lower-priority evidence backfill for older
  patch-train releases.
- Cloud provider, usage metering, and cost/rate-limiting work remain parked
  until a fresh product or release lane explicitly pulls them forward.

## Verification

- [x] `go run ./cmd/hawp work validate` — pass, 0 issues, 0 warnings
- [x] `go run ./cmd/hawp kit validate` — pass, 0 issues
- [x] Post-README rerun: `go run ./cmd/hawp work validate` — pass, 0
      issues, 0 warnings
- [x] Post-README rerun: `go run ./cmd/hawp kit validate` — pass, 0 issues
- [x] Focused tests after validator guardrail:
      `go test ./internal/domain/work ./internal/application/work` — pass
- [x] Current-repo validation after UUID folder move:
      `go run ./cmd/hawp work validate` — pass, 0 issues, 0 warnings
- [x] `mochila-archive-viewer` read-only validation:
      `go run ./cmd/hawp work validate --hawp-root /path/to/projects/mochila-archive-viewer/.hawp`
      — fail as expected for pre-migration issues: active rows `048`, `050`,
      `051`, `052`, and `053` still listed active while their plans are closed;
      the corresponding closed plans lack modern Outcome/Verification/Close
      Checklist sections.
- [x] Full Go suite: `go test ./...` — pass.
- [x] Checked-in binary version: `./.hawp/bin/hawp version` — `0.0.24`.
- [x] Mochila temp-copy apply proof:
      default normalize changed 6 files, folder migration changed 11 files,
      final validation passed with 0 issues and 0 warnings.
- [x] Tekit temp-copy apply proof:
      default normalize changed 194 files, folder migration changed 35 files,
      final validation passed with 0 issues and 2 tolerated warnings.

## 2026-09-07 progress report and sequencing

The [current report](../../status/2026/09/07/e5fca9c7-status.md) separates
implemented behavior, verification, architecture gaps, and the next milestones.
Start by organizing the verified uncommitted checkpoint; then define the intake
reshape contract and prove persisted UUID/status correctness before broader
boundary extraction. Release scope is an explicit later gate, not dependent on
finishing every future architecture idea. Existing UUIDs own these milestones;
no duplicate work items were created for this reporting pass.

2026-09-07 checkpoint consolidation: saved the complete tracked/untracked
[review-group inventory](../../evidence/2026/09/07/e5fca9c7-review-groups.md)
before the draft-service continuation. Code groups include their untracked
tests; moved plans include retained archives. The inventory is a review map,
not a claim of independent commit buildability. No staging or commits performed.

**Closed:** 2026-09-11 — work completed in feature/v0.0.24-work-folder-normalization

## Outcome

v0.0.24 patch branch completed. README positioning improved. UUID-folder
validation tightened. Work-normalize false-positive fix applied (active IDs
canonicalized; linked parked/recently closed legacy slugs preserved). Source-
layout launcher aligned. CLI version advanced to 0.0.24, changelog updated,
binary built and checked in. Mochila and Tekit normalization proven on temporary
copies. All active plan items checked off.

## Verification

- [x] `go run ./cmd/hawp work validate` — pass, 0 issues, 0 warnings
- [x] `go run ./cmd/hawp kit validate` — pass, 0 issues
- [x] `go test ./...` — pass
- [x] `./.hawp/bin/hawp version` — `0.0.24`

## Close Checklist

- [x] All focused checks pass (recorded in Verification above).
- [x] HAWP check passes.
- [x] Plan status: done.
- [x] Plan moved to `.hawp/work/closed/2026/09/11/e5fca9c7/plan.md`.
- [x] BACKLOG.md updated.
