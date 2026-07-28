---
work-item: n3p5r7s9
type: release
title: "Tag v0.0.3: ship context reshaping (ONNX + Ollama backends)"
status: plan-ready
owner: unassigned
created: 2026-07-26
updated: 2026-07-27
---

# Tag v0.0.3

## Mission

Cut and push the actual `v0.0.3` release tag for context reshaping
(ONNX + Ollama backends, `k8i9f5m1` / `p8r0t2u4` / the `v003-ship-audit`
fixes) once the ship-readiness audit's findings are confirmed fixed.

## Status note (2026-07-27)

Confirmed via `git log` that `librarian-go-v0.0.3` points at a 2026-07-21
commit — before any of this work existed. Nothing has actually shipped yet
despite `version.go` already reading `0.0.3`. 61 commits sit on `dev` ahead
of that tag as of 2026-07-27. This is deliberately left open (not tagged)
pending an explicit decision to release — see the same open question for
`m2o4p6q8` about the stale/inconsistent `0.0.1`–`0.0.6` tag history.

## Plan

1. Decide the real version number to tag, accounting for the confusing
   existing tag history (`librarian-go-v0.0.1/2/3/6`, all stale).
2. Confirm the audit's 7 fixes (`v003-ship-audit` closed record) are still
   valid against current `dev`.
3. Tag and push; let the release workflow build + publish binaries.
4. Verify release artifacts (`f3d5a0h6`).

## Outcome (filled at close)

_Pending._

## Verification (filled at close)

_Pending._
