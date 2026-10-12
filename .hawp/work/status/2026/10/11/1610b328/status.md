---
uuid: 1610b328-dd2f-4909-8c33-59569aa50f93
title: "Mochilada upstream U1-U9 source triage"
type: status
date: 2026-10-11
---

# Status Report

## Intent

Triage Mochilada's U1–U9 upstream report against HAWP `development` and make confirmed source fixes reviewable.

## Current State

Nine items are indexed in `.hawp/work/BACKLOG.md` with investigation and implementation records under `.hawp/work/active/`. Source changes are local on `codex/mochilada-upstream-triage`, based on `origin/development` at `d7c46a4e` (live remote tip checked 2026-10-11). They have not been merged or installed downstream.

## What Was Inspected

The report's referenced parser, normalization, validation, install/update, search, and metadata paths under `librarian/src/`, `distribution/sources/`, and `.hawp/kit/`; current generated guides under `distribution/generated/`; and HAWP intake and propagation guidance under `.hawp/kit/` and `.github/skills/`.

## What Changed

- U1–U3: tolerant backlog headings, linked legacy IDs and plan paths, `Work` title alias, and non-speculative handling of closed tables lacking a Type column.
- U4–U5: explicit inherited-closure and supporting-artifact markers; parked-row status drift warnings, review proposals, and an apply guard.
- U6: `core/.hawp/.gitignore` for runtime binary and database directories, seeded without overwriting downstream content.
- U7–U9: corrected search hint, clarified effective search scope, and recognized numeric folder IDs in metadata enrichment.
- `distribution/generated/` was refreshed from `distribution/sources/`.

## What Was Directly Verified

Focused Go regression tests and `go test ./...` passed. `hawp distribution sync` updated and validated 40 generated guides. `hawp kit validate`, `hawp work validate`, `hawp work normalize --dry-run --validate`, and `hawp check` all passed in the source worktree. The work check found Active 13/13, Closed 12/12, Parked 10/10, and zero issues.

## What Remains Unproven

The Mochilada counts and scratch-copy result were supplied by the downstream report and were not rerun here. No updated source guide or binary has been installed in Mochilada. Existing tracked runtime files there require separate untracking. U4's markers require explicit downstream record updates, and U5 deliberately requires review before moving plans or rewriting links.

## Constraints

The Mochilada repository was read only. No downstream files, binaries, Git index, or workflow records were edited. Source changes are limited to this isolated branch and remain unmerged.

## Help Wanted

Review the compatibility policy for inherited closures and the manual move requirement for parked plans before merging the branch.

## Suggested Next Step

Review the source diff and nine active plans, then integrate the branch through the normal HAWP review process. After release, rerun the matching update guide and validation in Mochilada.
