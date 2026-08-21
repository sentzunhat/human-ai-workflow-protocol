---
work-item: u7d4c8a1
type: task
title: "Audit docs drift and checkpoint alignment"
status: done
owner: unassigned
created: 2026-08-10
updated: 2026-08-10
---

# Docs Drift Audit and Checkpoint Alignment

## Mission

Audit workflow and documentation drift against the live HAWP ledger, status-report guidance, and repo-local docs. Update only the smallest confirmed set of files needed to align the docs and checkpoint trail.

## Progress

- Verified that `.hawp/work/STATUS.md` remains an intentional scaffold and migration target because the install/update contract still seeds it when missing and preserves it when present.
- Restored the safety guide entry so it matches that live contract.
- `npm run distribution:validate` passed: generated provider guides are current.
- Current workflow guidance and distribution output are aligned for the inspected scope.
- Historical notes still contain legacy mentions, but those are left untouched for now because they are not current operating guidance.

## Context

- The repo already has HAWP guidance, an active backlog, and dated status-report slots.
- Recent memory notes flagged docs drift away from bootstrap paths and a pattern of splitting broad audits into smaller follow-ups.
- No source-code change is intended unless the audit uncovers a direct docs-to-behavior mismatch.

## Constraints

- Keep the audit bounded to docs, ledger, and workflow surfaces.
- Do not duplicate unrelated GPU, infrastructure, or financial checkpoints unless material evidence changes them.
- Separate direct evidence from inference.
- Preserve unrelated dirty work if it appears later.

## Outcome

- Confirmed the live install/update contract still seeds and preserves `.hawp/work/STATUS.md`.
- Confirmed `npm run distribution:validate` passes with generated outputs current.
- Confirmed no further live docs correction is indicated within the bounded workflow/distribution scope.

## Expected Output

- A short mismatch summary.
- A conservative edit plan.
- A minimal set of doc/workflow fixes, if any.
- A dated status report that records before/after and what remains unproven.

## Verification

- `npm run distribution:validate` passed: generated provider guides are current.
- Confirmed `.hawp/work/STATUS.md` seeding and preservation behavior matches
  live install/update contract.
- No source-code changes made; only minimal doc alignment applied.

## Close Checklist

- [x] Docs drift audit completed within bounded scope
- [x] `distribution:validate` passed
- [x] STATUS.md seeding contract confirmed
- [x] Historical notes with legacy mentions left intentionally untouched
- [x] Moved to `closed/2026/08/10/`
