# TASK-021 Status Report

## Summary

Created a design note for future HAWP librarian fix-up and upgrade command flows with explicit safety and evidence guardrails.

## What Changed

- Added command model covering `validate`, `work-items fix-up`, `work-items upgrade`, `backlog refine`, and `backlog work-items-upgrade`.
- Defined dry-run-first policy, apply preconditions, refusal paths, and reversible-change requirements.
- Split operation classes into mechanical, AI-draft, and manual-required.
- Added explicit AI drafting boundaries for SmolLM2 local WASM/Transformers.js use.
- Included a concrete dry-run output example and recommended V1 scope.

## Verification Evidence

- `.hawp/work/notes/2026/05/11/hawp-librarian-fixup-upgrade-command-model.md`
- `.hawp/work/closed/2026/05/11/TASK-021.md`

## Outcome

The repository now has a reviewable, implementation-oriented baseline spec for safe HAWP cleanup/upgrade command development without introducing non-file system dependencies.
