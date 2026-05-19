# TASK-020 Validator CLI Help Output Evidence

Date: 2026-05-11
Owner: agent

## Scope

Implement validator CLI help output only, documenting:

- default local behavior
- --hawp-root
- --work-root
- --debug-closed-task
- exit code behavior
- FAIL vs WARN vs INFO meaning

## Commands Run

1. cd human-ai-workflow-protocol/librarian && npm run validate:workflow -- --help
2. cd human-ai-workflow-protocol/librarian && npm run typecheck
3. cd human-ai-workflow-protocol/librarian && npm run validate:workflow

## Results

- Help command succeeded and printed the full usage and option documentation.
- Typecheck succeeded with no errors.
- Local validator remained unchanged in behavior and returned VALIDATION PASS.

## Notes

No validation rules were changed.
No fix-up/apply/upgrade command work was added.
No UUID migration or SQLite/indexing/search/queueing work was started.
