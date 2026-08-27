# v0014-token-speed-bench — v0.0.14 token-savings and speed benchmark gate

**Type:** feature
**Status:** plan-ready
**Branch:** `feature/v0.0.13`
**Opened:** 2026-08-26

## Problem

`v0.0.14` needs a stronger proof story than “feature works.” The release should
show that token counting is actually saving context budget and that the relevant
paths remain fast enough to justify the added instrumentation and shaping logic.

## Goal

Make benchmark evidence a merge/release gate for `v0.0.14`, with two explicit
questions answered:

1. Are token-in / token-out measurements showing meaningful savings?
2. What is the speed cost or speed improvement of the associated feature work?

## Benchmark scope

### Token-savings evidence

- Fixed representative query set across kit docs and work items
- Compare raw candidate context vs final shaped context
- Record:
  - estimated tokens in
  - estimated tokens out
  - estimated tokens saved
  - percentage saved
- Include at least one case where savings are low or neutral so the report is
  not cherry-picked

### Speed evidence

- Benchmark end-to-end command latency for the affected workflow
- Benchmark any new indexing, search, or reshape overhead introduced in `v0.0.14`
- If the feature affects MCP responses, capture end-to-end tool latency
- Report median and spread across repeated runs when feasible

## Deliverables

- Benchmark artifact under `.hawp/work/evidence/2026/08/26/`
- Summary table suitable for release notes or PR description
- Explicit recommendation: keep, tune, or rollback if savings do not justify cost

## Acceptance criteria

- [ ] Token-savings report shows concrete saved-token counts and percentages
- [ ] Speed report shows command or tool latency on a fixed scenario set
- [ ] Evidence includes methodology and hardware/runtime notes
- [ ] Release or PR summary states whether the tradeoff is favorable

## Notes

- Prefer the existing benchmark harness where possible so results stay
  comparable across versions.
- If a feature changes both relevance quality and performance, note that the
  benchmark is necessary but not sufficient for the final ship decision.
