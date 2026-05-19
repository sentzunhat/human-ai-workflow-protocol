# Validator Rerun Evidence (HAWP, Tekit, Mictlan)

Date: 2026-05-11
Owner: agent

## Scope

Rerun validator against:

- local HAWP repo
- Tekit .hawp
- Mictlan .hawp

## Commands

1. cd human-ai-workflow-protocol/librarian && npm run validate:workflow
2. cd human-ai-workflow-protocol/librarian && npm run validate:workflow -- --hawp-root tekit/.hawp
3. cd human-ai-workflow-protocol/librarian && npm run validate:workflow -- --hawp-root mictlan/.hawp

## Results

- Local HAWP:
  - Backlog consistency: Active Work Found 1/1
  - Result: VALIDATION PASS

- Tekit:
  - BUG-157 and BUG-158 are no longer modern FAIL items
  - Result: VALIDATION PASS (warnings only)

- Mictlan:
  - Result: VALIDATION PASS (warnings only)

## Notes

No code changes were made in Tekit or Mictlan.
TASK-020 CLI help output implementation remains pending and has not been started in this step.
