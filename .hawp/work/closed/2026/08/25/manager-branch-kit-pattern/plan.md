# manager-branch-kit-pattern — Document manager-branch / worktree operating pattern in kit

**Type:** improvement  
**Status:** done  
**Opened:** 2026-08-25  
**Closed:** 2026-08-25
**Target:** v0.0.11

**Closes:** `manager-branch-kit-pattern`

## Goal

Add an optional kit doc (e.g. `.hawp/kit/usage/manager-branch.md` or as a
section in `start-here.md`) describing this pattern clearly enough for a new
team to adopt it. Keep it short and decision-useful.

## Outcome

Added `.hawp/kit/usage/manager-branch.md` as a concise optional operating note
for teams that keep HAWP coordination on a long-lived manager branch while doing
product work in separate git worktrees from an integration branch. The doc
stays explicitly optional, keeps HAWP out of runtime territory, and points back
to existing workflow/backlog guidance instead of inventing new lifecycle fields.

`start-here.md` now links to this doc in the workflow-guide list, and
`usage/hawp-first-workflow.md` now points readers to the manager-branch note
when they need a broader coordination pattern around parallel worktrees.

## Verification

- [x] Added `.hawp/kit/usage/manager-branch.md`. Evidence: see Outcome section above.
- [x] Linked the new guide from `.hawp/kit/start-here.md`. Evidence: see Outcome section above.
- [x] Added cross-reference from `.hawp/kit/usage/hawp-first-workflow.md`. Evidence: see Outcome section above.
- [x] `mise exec node@26.5.0 -- npm run distribution:sync`. Evidence: see Outcome section above.
- [ ] explicitly unproven in that checkout: `mise exec node@26.5.0 -- npm run kit:validate` — attempted, but `../.hawp/bin/hawp kit validate` exited 137 in this checkout
- [ ] explicitly unproven in that checkout: `mise exec node@26.5.0 -- npm run check:markdown-links` — attempted, but `../.hawp/bin/hawp links check` exited 137 in this checkout

## What was done

- Added an optional manager-branch guide covering when to use the pattern, manager ownership, product worktree ownership, and guardrails
- Linked the pattern from the top-level start guide and from the existing parallel-worktree guidance
- Regenerated distribution outputs and confirmed no source-template drift from the new doc links

## Close Checklist

- [x] Outcome recorded
- [x] Verification captures both the successful and unproven checks
- [x] Optional-pattern scope stays documented
- [x] Ready to stay in closed history
