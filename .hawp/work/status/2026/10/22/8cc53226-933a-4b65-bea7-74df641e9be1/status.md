---
uuid: 8cc53226-933a-4b65-bea7-74df641e9be1
title: "PR review to approval-gated remediation work planning"
type: status
date: 2026-10-22
event_date: 2026-09-22
---

# Status Report

#### Intent

Preserve a reusable, review-only coding-agent instruction decision from
2026-09-22, archived on 2026-10-22. Make verified pull-request findings
actionable for a later implementer without authorizing automatic fixes.

#### Current State

**Decided instruction profile, not a software implementation.** A reusable
PR/branch review instruction set was refined into a structured
review → verification → prioritized work plan → explicit approval →
approved remediation workflow. No specific PR was reviewed, no actual
defects were verified, and no remediation work items were created in
this conversation.

#### What Was Inspected

- The conversation's initial current-branch/PR review instruction set and
  its later detailed work-planning/approval revision.
- Existing repository guidance:
  `.hawp/kit/start-here.md`,
  `.hawp/kit/usage/intake-workflow.md`,
  `.hawp/kit/usage/workflow-loop.md`,
  `.hawp/kit/usage/status-report.md`,
  `.hawp/work/README.md`, and
  `.hawp/work/status/README.md`.
- Related historical status records under
  `.hawp/work/status/2026/09/21/` and
  `.hawp/work/status/2026/09/25/`.

#### What Changed (Before → After)

**Before (2026-09-22):** The proposed agent instruction emphasized
reviewing the checked-out branch against its merge base, tracing relevant
callers/contracts, inspecting security and data integrity, running available
checks, verifying findings, and returning a severity-ranked PR verdict.
The output was primarily a review report.

**During (2026-09-22):** The instruction was extended to require a
durable work plan for each confirmed actionable finding. Required work-item
content: ID, priority P0–P3, title, source finding, problem, evidence,
root cause, scope, required changes, **detailed observable outcome**,
binary acceptance criteria, validation plan, likely files, dependencies,
risks/notes, and out-of-scope boundaries. Related findings may share one
item only when the root cause and fix are genuinely shared.

**After (2026-09-22):** The intended sequence is:
1. Inspect current branch, available PR context, and surrounding code.
2. Verify actionable findings; distinguish new defects from pre-existing
   issues and report validation honestly.
3. Prioritize P0 → P1 → P2 → P3, adjusting for genuine dependencies.
4. Create or update executable work plans instead of duplicate records.
5. Present findings, implementation order, and an explicit approval request.
6. **Stop before code changes.** Silence is not approval.
7. Implement only expressly approved items, preferably one at a time;
   add regression tests, verify acceptance criteria, and re-review the diff.
8. Record unexpected out-of-scope issues as proposed work, not unapproved
   changes. Never silently push, merge, or resolve remote review threads.

#### Strategic Impact

This profile separates reviewer, planner, and executor responsibilities,
preserves human control of remediation scope, and makes findings
transferable across agent sessions. The detailed outcome defines what
must be observably true after a fix rather than prescribing only an edit.

This is a **strict PR-remediation approval profile**, not a claim that all
HAWP work universally needs manual approval. The existing HAWP intake
guide still documents its general risk-based gate (including low-risk
auto-implementation). The stricter user-selected policy applies when
using this review instruction profile.

#### What Was Directly Verified

- The repository uses UUID-scoped status reports under
  `.hawp/work/status/YYYY/MM/DD/<uuid>/status.md`.
- Existing HAWP work-item plans live under
  `.hawp/work/active/<uuid>/plan.md`, indexed in
  `.hawp/work/BACKLOG.md`.
- Historical PR-related status reports exist, but describe particular
  remediation sessions rather than this reusable instruction decision.
- The conversation contains an instruction specification and approval
  decision, **not** evidence of a successful PR review or code fix.

#### What Remains Unproven

- Whether this instruction profile has been incorporated into distributed
  agent prompts or tested with a real PR.
- Whether each tool/provider can persist findings directly as HAWP work
  plans without an additional adapter.
- No build, lint, tests, or live remediation validation were run for this
  instruction-only decision.

#### Constraints and Reconciliation

- HAWP remains a shaping protocol, not an agent runtime.
- Do not add fields to the HAWP core shape or create parallel work-item
  directories merely to represent the review template.
- Review-local `WI-001` identifiers are presentation labels; use the
  repository's canonical UUID identifiers when persisting actual HAWP
  work items. Keep a clear mapping if both are displayed.
- Review/planning artifacts may be prepared before approval; actual code
  remediation remains gated. Publishing this historical checkpoint is a
  distinct authorized documentation operation.
- No personal details, non-public business context, or private project
  names are included in this public archival record.

#### Help Wanted

Before distributing this as default agent guidance, reconcile the strict
PR-specific approval contract with existing HAWP risk-based workflow and
test it against a real branch/PR. Update existing instruction sources,
not duplicate generated provider outputs.

#### Suggested Next Step

Select a real PR, run review-only evaluation, persist only confirmed
findings as canonical HAWP plans, present an ordered work plan, and obtain
explicit implementation approval. If adopting this profile as standard
guidance, first inspect the source → generated → consumer propagation
chain and update the authoritative instruction source.
