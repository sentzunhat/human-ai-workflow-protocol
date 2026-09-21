# Status Report — Intake Draft Contract

## Intent
Advance request-to-intake reshaping without selecting a model or coupling draft
creation to persistence. [Work plan](../../../../active/a3df8a9c/plan.md).

## Current State
Internal DraftIntake API and domain Draft value implemented. Existing work-new
behavior is unchanged. No public command/tool or default shaper is wired.

## What Was Inspected
Work application/domain packages, the existing summarizer contract, HAWP field
specification, current dirty checkpoint, and work coordination records.

## What Changed
Added a typed proposal boundary that cannot replace input/context. The service
retains source strings, labels absent context, validates required fields, and
propagates provider errors/cancellation without a partial draft. Added tests and
package/source docs. Saved a complete pre-draft review-group inventory under
`.hawp/work/evidence/2026/09/07/e5fca9c7-review-groups.md`.

## What Was Directly Verified
Focused/full Go suites, vet, distribution validation, and connected HAWP MCP
validation pass. Tests cover source/wire-shape preservation, missing context,
incomplete proposals, invalid requests, provider failure, and cancellation.
The service has no persistence dependency and does not call work creation.

## What Remains Unproven
A real adapter's semantic fidelity, model execution, and side effects. Nonblank
output is not proof of factual accuracy or approval. No public user flow exists
for automatic generation yet; injected test shapers are deterministic fixtures.

## Constraints
No new dependency, model download, work-record creation from a draft, staging,
commit, binary refresh, or publication. Preserve the prior dirty checkpoint.

## Help Wanted
No immediate user input. Compare adapter value with worker-guided shaping before
adding runtime/client integration.

## Suggested Next Step
Evaluate a concrete adapter with intent-preservation examples, or prioritize
persisted work UUID/status correctness while that product decision is open.
