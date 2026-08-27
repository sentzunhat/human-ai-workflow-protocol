# Status Report

## Intent

Capture the current release state after merging the ORT release-workflow repair so
the next agent session can continue from a clean checkpoint instead of
reconstructing the release history.

## Current State

- `main` is at `0e14a63` (`chore: update checked-in hawp binary to v0.0.14`)
- Release tag `0.0.14` exists and the pre-release is published on GitHub
- Checked-in `.hawp/bin/hawp` on `main` has been refreshed from the published
  `hawp-darwin-arm64` release asset and now reports `0.0.14`
- The `Release` workflow for `0.0.14` published assets but completed with
  failure because the Linux ORT lane still fails

## What Was Inspected

- `.github/workflows/release.yml`
- `.github/workflows/tag-on-merge.yml`
- `librarian/src/internal/domain/update/version.go`
- `librarian/src/CHANGELOG.md`
- GitHub PR `#27`
- GitHub Actions runs for `Quality`, `Tag and release on merge to main`,
  `Auto-update integration test`, and `Release`
- GitHub release `0.0.14`

## What Changed

- Repaired the ORT release workflow pins for `onnxruntime-genai` and
  `tokenizers`
- Added `pipefail` to the ORT download steps
- Fixed `onnxruntime-genai` archive extraction for the macOS and Linux ORT lanes
- Bumped release version to `0.0.14`
- Refreshed the checked-in `.hawp/bin/hawp` binary on `main` from the published
  `0.0.14` `hawp-darwin-arm64` asset

## What Was Directly Verified

- PR `#27` merged to `main` on 2026-08-27
- `Quality` passed on the PR head and again on `main`
- `Tag and release on merge to main` passed
- `Auto-update integration test` passed
- GitHub pre-release `0.0.14` exists
- Published `0.0.14` assets include:
  - standard binaries for the six non-ORT platforms
  - `hawp-darwin-arm64-ort.tar.gz`
  - `hawp-kit-bundle.tar.gz`
  - `checksums.txt`
- Local `.hawp/bin/hawp version` returns `0.0.14`
- Local HAWP validation passes after the binary refresh

## What Remains Unproven

- No published Linux ORT tarball exists for `0.0.14`
- `Release` is not fully green for `0.0.14`; the workflow completed with failure
- Linux ORT release packaging still needs one more fix cycle to prove the full
  intended ORT release matrix

## Constraints

- Do not rewrite or delete the existing `0.0.13` or `0.0.14` releases
- Preserve the existing published `0.0.14` asset set as historical evidence
- Treat the remaining Linux ORT issue as next-patch work

## Help Wanted

- Confirm whether the next patch should focus only on the Linux ORT lane or also
  tighten any release-status expectations around partial ORT success
- Review the Linux ORT linker-path fix before cutting the next patch

## Suggested Next Step

Create a focused follow-up patch that fixes the Linux ORT build env path in
`.github/workflows/release.yml`, verify the lane locally where possible, then cut
the next patch release so the full release workflow finishes green with both ORT
tarballs.
