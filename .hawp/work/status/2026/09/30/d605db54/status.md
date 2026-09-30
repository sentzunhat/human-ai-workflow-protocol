---
uuid: d605db54
title: "v0.0.24 binary and release preparation"
type: status
date: 2026-09-30
---

## Intent

Prepare `0.0.24` release source and notes without carrying a host binary in
Git. See the [closed work plan](../../../../../closed/2026/09/30/d605db54/plan.md)
for file-level scope and verification.

## Current State

The version constant is `0.0.24`; the latest published release seen during
preparation was `0.0.23`. The current local `.hawp/bin/hawp` remains installed
but is now ignored. Install/update guides fetch the newest published platform
asset, including prereleases, with SHA256 verification.

## What Was Directly Verified

The focused installer test, full Go vet/test/build check, all six standard
cross-builds, generated distribution validation, kit/work/link validation,
and changelog extraction passed. Work validation had zero issues and its
existing warning for 25 ambiguous historical verification records.

## What Remains Unproven

The `0.0.24` tag and release do not exist yet. Read back PR #42 checks after
the push; merge approval is still a human gate. After merge to `main`, confirm
the release job, published assets and checksums, and a fresh downstream
install from the published release.
