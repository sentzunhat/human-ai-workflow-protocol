---
work-item: 392313e4-eded-402b-9d5e-20350c86b856
type: migration
title: "Normalize dated UUID work-item folders and repair workflow references"
status: plan-ready
created: 2026-08-13
updated: 2026-08-13
parent: b6c4e8a2
depends-on: c1d2e402
---

# Migration: Work-Item Identity and References

## Mission

Make each newly created work item self-contained in a dated UUID directory
while preserving the compact `BACKLOG.md` index and historical archive.
Provide a dry-run-first repair command for safe, unambiguous workflow drift.

## Target Layout

New work records use one canonical directory per item:

```text
.hawp/work/
  active/YYYY/MM/DD/<uuid>/
    plan.md
    references/
    evidence/
  parked/YYYY/MM/DD/<uuid>/
    plan.md
    references/
    evidence/
  closed/YYYY/MM/DD/<uuid>/
    plan.md
    references/
    evidence/
```

Every work-item-owned reference and evidence file lives in that UUID folder.
Status remains a date-based shared checkpoint unless it is explicitly owned by
one work item. `BACKLOG.md` continues to link to the canonical `plan.md`.

## Scope

- Teach `hawp work validate` to recognize canonical UUIDs, documented legacy
  compact IDs, and non-work checklist rows without treating them as malformed.
- Teach `hawp work normalize --dry-run` to produce a complete move and link
  rewrite plan for unambiguous work items, including existing flat UUID-named
  plan files under `active/`, `parked/`, and `closed/`.
- Add an opt-in apply mode that creates target directories, moves one work
  item folder atomically where possible, repairs its local references and
  evidence links, and revalidates afterward.
- Make `hawp work normalize` produce and apply the work-item folder plan:
  canonical plan path, owned reference/evidence files, `BACKLOG.md` link, and
  repository-owned workflow references are one reviewed operation.
- Upgrade an existing UUID-named work-item file to its dated UUID directory
  instead of generating a second identifier or leaving a duplicate file.
- Repair repository-owned workflow links that point to provider work records;
  provider source and generated packs remain read-only inputs.
- Report ambiguous IDs, duplicate candidates, external links, and references
  escaping `.hawp/work` as blocked operations requiring review.

## Non-Goals

- Do not bulk-convert historical records merely to satisfy a new layout.
- Do not move legacy evidence or reference files unless their owner UUID is
  established by an existing canonical link or unambiguous metadata.
- Do not change provider content or provider-generated workflow guidance in
  this migration; only repository-owned integration references are eligible.
- Do not infer replacement UUIDs for compact or malformed legacy IDs.

## Safety Gates

- Dry-run is the default and emits JSON/text plans without filesystem writes.
- Apply requires an explicit migration flag, clean worktree unless overridden,
  and an approved plan containing no unresolved blocked operations.
- A manifest records each move and rewritten reference so the exact scope is
  reviewable and reversible, including every owned evidence/reference file.
- Validation must pass for every migrated item before the next item is moved.

## Done When

- New `hawp work new` records can be created in the dated UUID layout.
- New work records include owned `references/` and `evidence/` directories.
- Existing flat UUID-named work-item files are reported by dry-run and moved
  by explicit apply into their canonical dated UUID folders without changing
  their UUIDs.
- Validator and normalizer classify legacy compact records correctly and no
  longer propose unsafe blanket ID repairs for checklist rows.
- A fixture proves an active, parked, and closed item folder, including owned
  references and evidence, can be moved with `BACKLOG.md` and local Markdown
  links repaired.
- Broken, ambiguous, and external references remain reported rather than
  changed; focused tests, build, and diff checks pass.
