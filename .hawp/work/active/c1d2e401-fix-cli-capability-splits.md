---
work-item: c1d2e401
type: fix
title: "Split CLI routing by command capability"
status: plan-ready
created: 2026-08-10
updated: 2026-08-10
parent: b6c4e8a2
depends-on: c1d2e400
---

# Fix: CLI Capability Splits

## Mission

Mechanically split CLI routing into capability-local command files, starting
with search/index, then update and work, without changing behavior.

## Done When

- `run.go` no longer owns unrelated command families.
- Command registration and transport output remain stable.
- Capability-local tests cover routing and error paths.
- Targeted tests, build, and diff checks pass.
