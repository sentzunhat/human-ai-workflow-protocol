---
work-item: e2c4f9g5
type: infrastructure
title: "Update command v2: auto-detect platform, download binary, verify SHA256"
status: done
closed: 2026-07-23
---

# `e2c4f9g5` — Update command v2 already complete

## Outcome (filled at close)

Confirmed already fully implemented: platform detection (all 6 combinations),
GitHub Release asset download, SHA256 verification, atomic binary replacement,
and full CLI wiring (`hawp update`, `hawp update --check`). No new
implementation needed — this was discovered pre-built during the v0.0.1
release-readiness pass.

Full analysis: [evidence](../../../evidence/2026/07/23/e2c4f9g5-update-already-complete.md)

## Verification (filled at close)

- ✅ `AssetName()` covers all 6 platform/arch combinations
- ✅ `Apply()` downloads via GitHub Release API + `VerifiedFile()` SHA256 check
- ✅ `selfreplace.Replace()` provides atomic, corruption-safe binary swap
- ✅ `hawp update` / `hawp update --check` / `hawp version` wired in the CLI
- ✅ `update_test.go` covers platform detection + version comparison

## Close Checklist

- [x] Outcome recorded
- [x] Verification recorded (see linked evidence for full command output)
- [x] No follow-on work — item was already complete
