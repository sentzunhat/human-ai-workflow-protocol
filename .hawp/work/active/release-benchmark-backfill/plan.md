# release-benchmark-backfill — Benchmark evidence gate for v0.0.11 through v0.0.13

**Type:** planning
**Status:** plan-ready
**Branch:** `feature/v0.0.13`
**Opened:** 2026-08-26

## Problem

The current patch train has meaningful product changes across `v0.0.11`,
`v0.0.12`, and `v0.0.13`, but benchmark evidence is not yet consistently
tracked as a release artifact for each version. The requirement is now that
every version in the train ships with explicit benchmark evidence, not just
code changes and changelog notes.

## Goal

Produce benchmark evidence for each patch version before merge/release:

- `v0.0.11` — install/update/provider workflow evidence where timing or
  operational friction changed
- `v0.0.12` — local MCP usage-log overhead and behavior evidence
- `v0.0.13` — ORT release-build and reshape-path evidence

## Required evidence

### v0.0.11

- Confirm no benchmark gap is hiding behind docs-only framing
- Measure at least one operational path affected by the patch:
  `hawp init`, provider config write, or install/update path
- Record before/after notes or absolute timings if no true baseline exists

### v0.0.12

- Measure overhead of logging enabled vs disabled on representative MCP calls
- Record token-estimate outputs on a small fixed set of calls
- Confirm logging remains non-blocking from the user’s perspective

### v0.0.13

- Record ORT release-build durations by job/platform where possible
- Record reshape latency on supported hardware if available
- Confirm release artifact size changes are understood

## Deliverables

- One benchmark summary artifact per version under `.hawp/work/evidence/`
- A short checkpoint or status note cross-linking the three benchmark artifacts
- PR descriptions updated to mention benchmark evidence once captured

## Acceptance criteria

- [ ] `v0.0.11` has benchmark evidence attached or an explicit documented reason
      why only absolute timing evidence is possible
- [ ] `v0.0.12` has usage-log overhead evidence and token-estimate examples
- [ ] `v0.0.13` has release-build and ORT-path benchmark evidence
- [ ] Each patch PR references its benchmark artifact before merge

## Notes

- If a direct pre-change baseline no longer exists, use clearly labeled
  absolute measurements plus scope notes rather than inventing a comparison.
- Evidence must separate directly measured results from inference.
