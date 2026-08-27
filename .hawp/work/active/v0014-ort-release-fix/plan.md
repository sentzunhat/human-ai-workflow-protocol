# v0014-ort-release-fix — repair ORT release lanes after `v0.0.13` failure

**Type:** fix
**Status:** plan-ready
**Branch:** `feature/v0.0.14`
**Opened:** 2026-08-27

## Problem

`v0.0.13` merged cleanly and passed `Quality`, but the post-merge `Release`
workflow failed before a publishable `0.0.13` release could be created.

Observed failures from GitHub Actions on 2026-08-27:

- `build-ort-linux-amd64` failed in `Download ORT native libs (linux/amd64)`
  with exit code `2`
- `build-ort-darwin-arm64` failed in `Build ORT binary (darwin/arm64)` with
  exit code `1`

Until those ORT release lanes are repaired, `main` cannot publish a verified
`0.0.13` release artifact and the checked-in `.hawp/bin/hawp` binary cannot be
refreshed past `0.0.12`.

## Goal

Make the ORT release lanes reliable enough to cut the next patch release and
refresh the checked-in binary from a published release asset.

## Deliverables

- Root-cause explanation for the linux and darwin ORT release failures
- Workflow or build fixes on `feature/v0.0.14`
- Green `Quality`, `Tag and release on merge to main`, `Auto-update integration test`,
  and `Release` workflows for the next patch
- Published release assets that include the intended ORT outputs
- Follow-up binary refresh on `main` from the published release download

## Acceptance criteria

- [ ] Linux ORT native-lib download step succeeds in GitHub Actions
- [ ] Darwin ORT build step succeeds in GitHub Actions
- [ ] `Release` workflow completes successfully for the next patch
- [ ] Published release exists and is downloadable from GitHub Releases
- [ ] Checked-in `.hawp/bin/hawp` is refreshed from that published asset

## Notes

- Treat this as next-version work, not a retroactive rewrite of the already
  merged `v0.0.13` branch train.
- Preserve the exact failing workflow evidence from 2026-08-27 when diagnosing
  the regression so the fix targets the real release-path breakage.
