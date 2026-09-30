# Optional Manager-Branch Pattern

Use this pattern when your repository already has a product integration branch
and you do not want HAWP kit/backlog churn mixed into product pull requests.

## When it helps

- Your team merges product code through an integration branch such as `origin/QA`
- You want one stable checkout to own `.hawp/kit/` refreshes and `.hawp/work/**`
- You want product work to happen in separate worktrees without carrying HAWP
  coordination noise into every product branch

## Pattern

1. Keep one long-lived manager branch for HAWP coordination.
2. Let that manager branch own `.hawp/kit/`, `.hawp/work/**`, and dispatch notes.
3. Cut product worktrees from the product integration branch, not from the manager branch head.
4. Merge completed product work back into the patch or product branch you actually ship.
5. Do not open product PRs from the manager branch.

Example shape:

```text
manager branch: chore/hawp-manager
product integration branch: origin/QA
product work: git worktree add ... -b feature/v0.0.11-<slice> origin/QA
```

## What the manager branch owns

- `hawp init` / `hawp update` for the coordination checkout
- `.hawp/kit/`
- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/`, `closed/`, `status/`, `evidence/`, `decisions/`
- Agent dispatch and session continuity notes

## What product worktrees own

- Product code changes
- Focused verification for each implementation slice
- Branches that are safe to merge into the patch or product line

## Guardrails

- This is optional. HAWP does not require a manager branch.
- Do not treat the manager branch as a runtime control plane.
- If product code also edits `.hawp/`, expect merge friction and tighten ownership.
- Keep machine-local paths out of saved notes; use repo-relative paths in plans and evidence.

## Related guidance

- [hawp-first-workflow.md](hawp-first-workflow.md) — session-first HAWP workflow and parallel worktree notes
- [workflow-loop.md](workflow-loop.md) — iteration and verification flow
- [../references/backlog-alignment.md](../references/backlog-alignment.md) — backlog and closure discipline
