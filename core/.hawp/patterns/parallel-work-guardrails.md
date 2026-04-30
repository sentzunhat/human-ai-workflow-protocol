# Parallel Work Guardrails

Use this optional pattern when a human and multiple AI agents may run in parallel.
Goal: reduce conflicting edits and unclear ownership while keeping HAWP lean.

## Guardrails

1. Check active work first in `usage/BACKLOG.md` before planning or implementing.
2. Use one backlog ID and one plan file per item.
3. Record owner and overlapping files in the plan before implementation.
4. If overlap touches the same file as another in-progress item, default to hold.
5. Only proceed with overlapping edits when the user explicitly approves.
6. Keep coordination notes in the plan file, not in a separate lock system.

## Minimal Coordination Block

```md
### Work Coordination

**Owner:** human | agent | unassigned
**Implementation status:** not-started | in-progress | blocked | done
**Overlapping files:**

- path/to/file.ts

**Parallel work risk:** low | medium | high
**Can implement now:** yes | no | only after approval

**Coordination note:**
Explain whether another active item touches the same files.
```

## Quick Example

- TASK-014 is `in-progress` and edits `src/auth/session.ts`.
- TASK-015 also needs `src/auth/session.ts`.
- Set `Can implement now: only after approval` in TASK-015 plan and wait.

This pattern is markdown-first and file-first. It does not add runtime locking, tooling, or schema changes.
