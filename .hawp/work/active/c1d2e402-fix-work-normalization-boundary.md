---
work-item: c1d2e402
type: fix
title: "Extract work normalization scan and mutation boundary"
status: plan-ready
created: 2026-08-13
updated: 2026-08-13
parent: b6c4e8a2
depends-on: c1d2e3f9
---

# Fix: Work Normalization Boundary

## Mission

Move normalization-specific plan scans, path/existence facts, closed-record
moves, and content writes behind a capability-local work normalization
workspace while retaining pure detection and report rules in `domain/work`.

## Constraints

- Begin only after `c1d2e3f9` establishes the read-only work source pattern.
- Preserve dry-run by default, dirty-worktree safety, blocked-operation
  reporting, research-queue export, and all existing normalization semantics.
- Keep ports and adapters together by work-normalization capability.

## Done When

- Domain normalization rules consume typed scan/mutation facts instead of
  directly reading, moving, or writing files.
- Filesystem mutations are adapter-owned and separately tested.
- Existing CLI output, plan exports, and safety behavior remain compatible.
- Focused tests, build, and diff checks pass.
