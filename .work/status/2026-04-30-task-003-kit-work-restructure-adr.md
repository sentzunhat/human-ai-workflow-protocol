# ADR: Restructure `.hawp` into `kit/` and `work/`

- ID: TASK-003
- Date: 2026-04-30
- Status: Approved (user pre-approved Option B in intake prompt)
- Risk: Medium (path churn across docs/install/update; no protocol semantics change)

## Context

`.hawp/` currently mixes reusable HAWP protocol assets (`README.md`, `SPEC.md`, `START_HERE.md`, `AUTHORING_PATTERNS.md`, `templates/`, `patterns/`, `reviews/`, `examples/`, `types/`) with repo-local work state (`usage/` containing `BACKLOG.md`, `INIT.md`, `INTAKE_WORKFLOW.md`, `STATUS_REPORT.md`, ADR docs, and `status/`).

This makes it harder for humans and AI agents to know which files are reusable kit content versus repo-local operating state. The install/update flow already maintains an allowlist that excludes `usage/`, but the visual separation is weak.

## Problem

A single mental model split is missing:

- "kit" — reusable HAWP material that downstream installs copy.
- "work" — the source repository's own operating state (backlog, ADRs, status reports, evidence).

## Current State (Repo Inspection)

```
core/.hawp/
  AUTHORING_PATTERNS.md
  LICENSE
  README.md
  SPEC.md
  START_HERE.md
  examples/
  patterns/
  reviews/
  templates/        (status-report.md, micro-hawp.md, standard-hawp.md, intake-plan.md, work-intake.md, bug-plan.md, audit-report.md)
  types/shape.ts
  usage/
    BACKLOG.md
    INIT.md
    INTAKE_WORKFLOW.md
    STATUS_REPORT.md
    GUARDRAIL_ADR.md
    INSTALL_BOUNDARY_ADR.md
    PUBLICATION_SAFETY_ADR.md
    status/
      .gitkeep
      2026-04-26-hawp-adr-template-review.md
      2026-04-30-task-001-parallel-work-guardrails-adr.md
      2026-04-30-task-002-intake-template-adr.md
```

`.gitignore` ignores `core/.hawp/usage/status/*.md`, so the on-disk status reports are not tracked. Only `.gitkeep` is tracked there.

## Options

- **A. Keep current structure.** No path churn, but kit and work stay mixed.
- **B. Two folders + root LICENSE: `kit/`, `work/`, `LICENSE`.** Clean mental model, modest path churn in install/update/docs. *(User-preferred.)*
- **C. Two folders only; LICENSE moved into `kit/`.** Strictest split, but LICENSE becomes less discoverable.
- **D. `protocol/` and `workspace/` instead of `kit/`/`work/`.** Formal naming, longer, less aligned with user's stated preference.

## Recommendation

**Option B.** Matches the stated mental model, keeps LICENSE in a conventional place, and aligns with the existing install boundary (kit allowlist already excludes `usage/`).

## Proposed Final Tree

```
core/.hawp/
  LICENSE
  kit/
    README.md
    START_HERE.md
    SPEC.md
    AUTHORING_PATTERNS.md
    examples/
    patterns/
    reviews/
    templates/
      audit-report.md
      backlog.md            (NEW — empty reusable backlog template)
      bug-plan.md
      intake-plan.md
      micro-hawp.md
      standard-hawp.md
      status-report.md
      work-intake.md
    types/shape.ts
  work/
    README.md               (NEW — explains BACKLOG/status/adrs/evidence)
    BACKLOG.md
    INIT.md
    INTAKE_WORKFLOW.md
    STATUS_REPORT.md
    adrs/
      README.md             (NEW)
      GUARDRAIL_ADR.md
      INSTALL_BOUNDARY_ADR.md
      PUBLICATION_SAFETY_ADR.md
    status/
      README.md             (NEW — anchor; replaces .gitkeep)
      2026-04-26-hawp-adr-template-review.md
      2026-04-30-task-001-parallel-work-guardrails-adr.md
      2026-04-30-task-002-intake-template-adr.md
      2026-04-30-task-003-kit-work-restructure-adr.md
    evidence/
      README.md             (NEW — anchor)
```

## Move Map

Kit (reusable):

- `core/.hawp/README.md` → `core/.hawp/kit/README.md`
- `core/.hawp/START_HERE.md` → `core/.hawp/kit/START_HERE.md`
- `core/.hawp/SPEC.md` → `core/.hawp/kit/SPEC.md`
- `core/.hawp/AUTHORING_PATTERNS.md` → `core/.hawp/kit/AUTHORING_PATTERNS.md`
- `core/.hawp/examples/` → `core/.hawp/kit/examples/`
- `core/.hawp/patterns/` → `core/.hawp/kit/patterns/`
- `core/.hawp/reviews/` → `core/.hawp/kit/reviews/`
- `core/.hawp/templates/` → `core/.hawp/kit/templates/`
- `core/.hawp/types/` → `core/.hawp/kit/types/`

Work (repo-local):

- `core/.hawp/usage/BACKLOG.md` → `core/.hawp/work/BACKLOG.md`
- `core/.hawp/usage/INIT.md` → `core/.hawp/work/INIT.md`
- `core/.hawp/usage/INTAKE_WORKFLOW.md` → `core/.hawp/work/INTAKE_WORKFLOW.md`
- `core/.hawp/usage/STATUS_REPORT.md` → `core/.hawp/work/STATUS_REPORT.md`
- `core/.hawp/usage/GUARDRAIL_ADR.md` → `core/.hawp/work/adrs/GUARDRAIL_ADR.md`
- `core/.hawp/usage/INSTALL_BOUNDARY_ADR.md` → `core/.hawp/work/adrs/INSTALL_BOUNDARY_ADR.md`
- `core/.hawp/usage/PUBLICATION_SAFETY_ADR.md` → `core/.hawp/work/adrs/PUBLICATION_SAFETY_ADR.md`
- `core/.hawp/usage/status/*.md` → `core/.hawp/work/status/*.md`

LICENSE stays at `core/.hawp/LICENSE`.

## Files Added

- `core/.hawp/kit/templates/backlog.md` (empty reusable backlog template)
- `core/.hawp/work/README.md`
- `core/.hawp/work/adrs/README.md`
- `core/.hawp/work/status/README.md` (replaces `.gitkeep`)
- `core/.hawp/work/evidence/README.md`

## Files Removed

- `core/.hawp/usage/status/.gitkeep` (folder anchored by README.md instead)
- `core/.hawp/usage/.DS_Store`

## Files Explicitly Not Touched (semantics)

- `core/.hawp/kit/types/shape.ts` (only its parent path changes; file content unchanged)
- protocol semantics in `kit/SPEC.md`, `kit/README.md`, `kit/AUTHORING_PATTERNS.md` (content unchanged)
- existing repo-local status reports under `work/status/` (content preserved verbatim — historical record)
- existing ADRs under `work/adrs/` (content preserved verbatim — historical record)

Historical references to `.hawp/usage/...` inside moved status reports and ADRs are intentionally left as-is because they document past state.

## Path Update Surface

External-facing files updated to new paths:

- `README.md`
- `core/install.md`
- `core/update.md`
- `core/.github/copilot-instructions.md`
- `core/.github/instructions/hawp-intake.instructions.md`
- `core/.github/prompts/hawp-status-report.prompt.md`
- `.github/copilot-instructions.md`
- `.github/instructions/hawp-intake.instructions.md`
- `.github/instructions/intake.instructions.md`
- `.github/prompts/hawp-status-report.prompt.md`
- `benchmark/README.md`
- `benchmark/benchmark-prompt.md`
- `.gitignore`
- moved `kit/templates/intake-plan.md`, `kit/templates/work-intake.md`, `kit/patterns/parallel-work-guardrails.md`
- moved `work/INIT.md`, `work/INTAKE_WORKFLOW.md`, `work/STATUS_REPORT.md`, `work/BACKLOG.md`
- moved `kit/README.md`, `kit/START_HERE.md`

## `.gitignore` Strategy

Goal: track folder shape (so the structure is copyable from git) without `.gitkeep` files; ignore work file contents that should remain local.

- Repath ignored prefixes from `core/.hawp/usage/` to `core/.hawp/work/`.
- Keep `*.md` ignore for `core/.hawp/work/status/*.md` and add an unignore for `README.md` so the folder anchor is tracked.
- Same pattern for `core/.hawp/work/evidence/*` with a tracked `README.md` anchor.
- Remove `.gitkeep` files from the repo.

## Risks

- Stale path references in any file not enumerated in the update surface.
- Downstream installs that have already copied `.hawp/` flat will need to refresh via `core/update.md`.
- Search-for-old-path discipline relies on docs distinguishing migration history from active usage.

## Acceptance Criteria

- `core/.hawp/` root contains only `LICENSE`, `kit/`, `work/`.
- All reusable assets under `core/.hawp/kit/`.
- All repo-local work state under `core/.hawp/work/`.
- `core/.hawp/work/BACKLOG.md` exists and is the active backlog.
- `core/.hawp/kit/templates/backlog.md` exists and is empty/generic.
- All status reports moved to `core/.hawp/work/status/`.
- All repo ADRs moved to `core/.hawp/work/adrs/`.
- `core/.hawp/work/evidence/` exists with no fake evidence.
- No `.gitkeep` files remain anywhere under `core/.hawp/`.
- `core/.hawp/kit/types/shape.ts` content unchanged.
- `install.md` and `update.md` reference `kit/` paths.
- No CLI/SQLite/validators/scripts/runtime tooling introduced.

## Verification Checklist

- [ ] `find core/.hawp -maxdepth 1` shows only `LICENSE`, `kit`, `work`.
- [ ] `grep -r "\.hawp/usage" --exclude-dir=.git` returns only matches inside `core/.hawp/work/status/` or `core/.hawp/work/adrs/` (historical records).
- [ ] `grep -r "\.hawp/templates\|\.hawp/patterns\|\.hawp/reviews\|\.hawp/types\|\.hawp/examples" --exclude-dir=.git` returns only historical/migration references.
- [ ] `find core/.hawp -name .gitkeep` returns no results.
- [ ] `core/.hawp/kit/templates/backlog.md` is empty/template-shaped (no real done rows).
- [ ] `core/install.md` and `core/update.md` copy kit files via `kit/` paths.

## Implementation Plan

1. Add TASK-003 row to `BACKLOG.md` (in-progress).
2. Write this ADR.
3. `git mv` files into `kit/` and `work/` per move map.
4. Add new `README.md` anchors and the empty `backlog.md` template.
5. Update `.gitignore` (work/ paths, README anchors, drop `.gitkeep` use).
6. Remove `.gitkeep` and `.DS_Store`.
7. Update `install.md`, `update.md`, top-level `README.md`, copilot/instructions/prompts, benchmark docs.
8. Update intra-repo cross-references in moved files (`work/INIT.md`, `work/INTAKE_WORKFLOW.md`, etc.).
9. Run grep verification.
10. Mark backlog row `done`.
