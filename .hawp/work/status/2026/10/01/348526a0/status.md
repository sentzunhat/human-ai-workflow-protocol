---
uuid: 348526a0
title: "v0.0.24 release tag retry gate"
type: status
date: 2026-10-01
---

## Intent

Keep the `0.0.24` release retryable if an external build or publication step
fails. See the [closed work plan](../../../../../closed/2026/10/01/348526a0/plan.md).

## Current State

PR #42 is still open and awaiting an approving review. The tag step in
`.github/workflows/release.yml` runs only after standard binary builds and
release-material preparation. A rerun continues only if the tag resolves to
the workflow's checkout commit.

## Direct Verification

The workflow parsed as YAML and the tag step passed missing-tag,
same-commit, and different-commit cases in a temporary local/bare Git pair.
No GitHub release was started during this verification.

## Remaining Gate

After PR approval and merge to `main`, inspect the actual release dispatch,
builds, tag, assets, and checksums before calling `0.0.24` published.
