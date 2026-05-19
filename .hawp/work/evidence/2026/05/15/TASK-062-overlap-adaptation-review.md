# TASK-062 Overlap Adaptation Review

Generated: 2026-05-15

## Scope

Compared local vs shared overlap candidates:

- `.hawp/kit/references/install-update-safety.md`
- `shared_standards/public/standards/docs/hawp-install-update-safety.md`
- `.hawp/kit/templates/adr-template.md`
- `shared_standards/public/templates/ADR.template.md`

## Diff Findings

### 1) Install/Update Safety Reference

- Shared file is much longer and includes operational execution detail.
- Local file is intentionally principle-level and explicitly avoids implementation details.
- Decision: keep local file unchanged to preserve repo boundary model and reference scope.

### 2) ADR Template

- Shared template is example-heavy and includes scenario-specific narrative.
- Local template is structured for repo-agnostic architecture decisions with explicit evidence/inference separation.
- Decision: keep local file unchanged to preserve current HAWP template style and governance intent.

## Adaptation Outcome

- Applied adaptation strategy: **reviewed, no-content-change** for both overlap files.
- Reason: local variants are already stricter for this repo's boundary and governance model.
- No direct absorbs were made from shared files in this pass.

## Follow-Up Split

Zacatl principle-level adaptation candidates were not merged in this docs/template overlap pass and are split into a dedicated follow-up:

- `TASK-063` — Extract generalized standards from Zacatl adaptation candidates.

## Verification Inputs

Commands used:

- `diff -u .hawp/kit/references/install-update-safety.md shared_standards/public/standards/docs/hawp-install-update-safety.md`
- `diff -u .hawp/kit/templates/adr-template.md shared_standards/public/templates/ADR.template.md`
