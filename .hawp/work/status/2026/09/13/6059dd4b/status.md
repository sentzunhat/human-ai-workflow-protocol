---
uuid: 6059dd4b-8ca5-485d-a158-f664fbd64f11
title: v0.0.24 cutoff and v0.0.25 MCP lane
type: status
date: 2026-09-13
---

# Status Report

## Intent

Record the release lesson from the MCP planning pass: `v0.0.24` should stay a
compact, robust patch focused on the compound intake tool, while broader MCP
surface expansion belongs to a follow-on `v0.0.25` lane.

This report supports deciding what to ship, where to cut scope, and which active
work items should move to the next version.

## Current State

The live backlog has four `plan-ready` active MCP records:

- `a8797f44` - `hawp_work_intake` MCP tool: compound search + reshape in one call
- `c8d9e1cb` - Interactive MCP intake refinement with structured results
- `2eea565c` - MCP capability catalog for HAWP resources and prompts
- `429e075e` - Remote MCP transport and authorization readiness audit

The backlog now shows `a8797f44` as the active `v0.0.24` implementation item,
with the other MCP records held as plan-ready follow-on work. It points to this
report as the release-scope source.

## Lesson Learned

The compounding release move is not to add many MCP surfaces at once. It is to
make the existing agent workflow harder to misuse:

1. Search context should not be optional when shaping work.
2. Reshape/compact should happen in the same MCP call as retrieval.
3. Missing context should be visible as a structured state instead of hidden in
   a confident-looking draft.
4. Token accounting should be returned with the draft so HAWP's token-reduction
   value is visible during normal agent use.

## Proposed v0.0.24 Cutoff

Ship through:

- `a8797f44` as the main implementation slice.
- The smallest compatible subset of `c8d9e1cb`: structured `hawp_work_intake`
  result states such as `ready_for_work_new`, `needs_user_input`, and
  `blocked_missing_index`.

Cut off before:

- MCP resources/prompts.
- Remote MCP transport or authorization design changes.
- Usage metering CLI/MCP work.
- Cloud model, embedding, rate-limit, or cost-tracking work.

## Proposed v0.0.25 Lane

Assign these work items to `v0.0.25` after the `v0.0.24` intake tool ships:

- `2eea565c` - capability catalog, read-only resources, and prompt candidates.
- `429e075e` - remote MCP transport/readiness audit and security boundary.
- Remaining parts of `c8d9e1cb` after the `v0.0.24` structured result subset,
  especially elicitation research and richer agent guidance.
- `usage-tracking` may be revisited only after `a8797f44` demonstrates the
  token-accounting surface in real intake use.

## What Was Inspected

- `.hawp/kit/start-here.md`
- `.hawp/kit/usage/status-report.md`
- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/a8797f44/plan.md`
- `.hawp/work/active/c8d9e1cb/plan.md`
- `.hawp/work/active/2eea565c/plan.md`
- `.hawp/work/active/429e075e/plan.md`

## What Changed

This report records the release cutoff and version lane. The active plans and
backlog point to this status report as the version-scope source.

## What Was Directly Verified

- `hawp work status --title "v0.0.24 cutoff and v0.0.25 MCP lane"` created this
  status artifact under the expected UUID subfolder.
- The four MCP records above exist in `.hawp/work/active/`.
- The backlog marks all four MCP records `plan-ready`.
- The backlog now separates the `v0.0.24` cutoff from the proposed `v0.0.25`
  lane.

## What Remains Unproven

- The `hawp_work_intake` implementation has not been written or tested yet.
- The exact MCP SDK support for resources, prompts, remote transport, and
  elicitation still needs fresh documentation review before any `v0.0.25`
  implementation.
- Token-accounting ergonomics in real agent workflows remain to be proven by
  the `a8797f44` implementation and smoke tests.

## Constraints

- Keep `v0.0.24` focused and robust.
- Do not add remote MCP listeners in the `v0.0.24` lane.
- Do not move parked cloud/cost/model items into the patch release.
- Preserve the current local-first HAWP posture unless a later decision changes
  it with explicit security requirements.

## Help Wanted

Review whether `v0.0.24` should include only `a8797f44`, or whether the small
structured-state subset from `c8d9e1cb` should be bundled into the same
implementation slice.

## Suggested Next Step

Approve `a8797f44` as the next implementation slice for `v0.0.24`, with the
structured-result subset from `c8d9e1cb` included only if it stays small and
directly improves `hawp_work_intake`.
