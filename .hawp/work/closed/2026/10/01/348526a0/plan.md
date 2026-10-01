# Make v0.0.24 release tagging retry-safe

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `348526a0-61bc-466e-8f74-c200586def7e`
**Type:** fix
**Reported:** 2026-10-01

---

## Input (verbatim)

> Move tag creation until release artifacts and notes are ready; allow rerunning the release only when an existing tag points to the same commit.

## Intake Summary

Make the `0.0.24` release dispatch retryable without creating the tag before
the standard binaries and release materials are ready.

## Current Context

PR #42 targets `main` from `development`. Both pull-request checks passed at
`19d69cfd`; approval is still required. `tag-on-merge.yml` dispatches
`release.yml` after a `main` push when the version tag is absent.

## Initial Analysis

**Directly verified:**

- `.github/workflows/release.yml` currently creates and pushes the tag in
  `build-std` before `make dist`, then refuses any existing tag.
- Its `release` job waits for `build-std` and the optional ORT jobs, packages
  the kit, generates checksums, and extracts changelog notes before publishing.
- `.github/workflows/tag-on-merge.yml` skips dispatch when the tag exists.

**Inferred (not yet proven):**

- A failure after tag creation leaves an existing tag that prevents a new
  dispatch and causes a rerun of `build-std` to fail. The actual `0.0.24`
  release has not run, so this is a source-derived failure path.

**Likely scope:**

- `.github/workflows/release.yml` plus this work record. The version, tag name,
  assets, and merge gate need no change.

## Risk + Review Gate

**Risk:** medium: release automation can publish a tag and GitHub release.
**Gate:** The user explicitly agreed to move the tag after builds and to fix
the retry path identified in the review. Do not dispatch or merge while
implementing this source change.

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/active/348526a0/plan.md

## Plan

1. Move tag creation into the `release` job after binary downloads, kit
   packaging, checksums, and release-note extraction.
2. Fetch tags in that job. On dispatch, create a missing tag; if it exists,
   continue only when it resolves to the current checkout commit. Reject a
   same-named tag at a different commit.
3. Inspect the workflow structure and exercise the tag guard against a
   temporary repository for missing, same-commit, and different-commit cases.
   Validate HAWP records and diff hygiene, then push the focused fix to PR #42.

## Verification

- YAML parsing confirmed the release job puts tag creation after checksums
  and changelog extraction, and before publication; `build-std` no longer
  creates a tag.
- The exact tag shell step was exercised in a temporary local/bare Git pair:
  missing tag created and pushed, same-commit retry continued, and a tag at a
  different commit was refused without changing the tag.
- `git diff --check` passed. HAWP work and link validation are recorded in the
  status report after closure.
- No actual tag, release workflow, or merge was run. GitHub Actions execution
  remains unproven until the approved merge to `main`.

## Outcome

`.github/workflows/release.yml` now creates/verifies the tag in the publish
job after standard cross-build success and release-material preparation.
`RELEASE.md` explains how to retry a failed publication. The release tag name,
version, and assets remain unchanged.

## Close Checklist

- [x] Release source change and local tag-guard verification recorded.
- [x] No tag, release dispatch, or merge performed.
- [x] Remaining GitHub Actions proof documented.

Status: [release tag retry gate](../../../../../status/2026/10/01/348526a0/status.md).
