# TASK-069 Public/Private Standards Audit

Date: 2026-06-04

## Scope

Source material reviewed:

- standards/public/\*\*
- standards/private/\*\*
- docs/\*\*

Destination scope reviewed:

- core/.hawp/kit/standards/\*\*

## Direct Evidence Summary

1. Public docs standards include a docs category with `hawp-install-update-safety.md` that was not present in `core/.hawp/kit/standards/**` before this task.
2. Existing core standards already matched key public guideline files (`guidelines/testing.md`, `guidelines/documentation.md`).
3. Private lane content includes internal auth route inventories, provider strategy cues, product rollout details, and internal runtime notes.
4. A private security guide includes generally reusable practices, but direct publication from private lane is not required because public-safe security guidance already exists in core standards.

## Classification

### Absorb Directly (implemented in TASK-069)

- standards/public/standards/docs/README.md -> core/.hawp/kit/standards/docs/README.md
- standards/public/standards/docs/hawp-install-update-safety.md -> core/.hawp/kit/standards/docs/hawp-install-update-safety.md

### Adapt Required (follow-up)

- standards/private/auth/protected-routes.md
  - Reason: contains concrete route inventory and token model examples tied to internal implementation details.
- standards/private/providers/README.md
  - Reason: includes provider lifecycle strategy language that should be generalized before public reuse.
- standards/private/product/README.md
  - Reason: includes project-origin references and launch-process framing that needs neutralization.

### Do Not Absorb

- standards/private/internal-runtime/README.md
  - Reason: explicitly private internal runtime topology guidance.
- standards/private/context/\*.yml
  - Reason: internal agent scaffolding and organization-specific context configuration.

## Decision

- Promote only public-safe docs standards into `core/.hawp/kit/standards/docs/` in this task.
- Keep private lane unabsorbed.
- Create a follow-up governance item for workflow-only adaptations from private lane.

## Follow-Up Work Item

- TASK-070: Evaluate private workflow-only standards for safe public adaptation.
