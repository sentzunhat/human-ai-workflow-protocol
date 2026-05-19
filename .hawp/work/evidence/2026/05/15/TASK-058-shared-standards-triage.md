# TASK-058 Shared Standards Triage

Generated: 2026-05-15

## Scope

- Source root: `shared_standards/`
- Rubric: `.hawp/kit/guidance/shared-standards-review-rubric.md`
- Action vocabulary: `absorb`, `adapt`, `split-private`

## Coverage Summary

- Total entries triaged: 55
- absorb: 10
- adapt: 25
- split-private: 20
- public-safe: 10
- repo-specific: 25
- private/proprietary: 20

## Per-File Triage Table

| source path | classification | action | destination | rationale |
| --- | --- | --- | --- | --- |
| shared_standards/README.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Non-private artifact requiring selective adaptation. |
| shared_standards/audits/2026-05-11-public-private-boundary-implementation.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Non-private artifact requiring selective adaptation. |
| shared_standards/audits/2026-05-11-standards-boundary-analysis.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Non-private artifact requiring selective adaptation. |
| shared_standards/audits/2026-05-11-standards-candidate-classification.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Non-private artifact requiring selective adaptation. |
| shared_standards/audits/2026-05-11-standards-harvest-review.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Non-private artifact requiring selective adaptation. |
| shared_standards/audits/2026-05-11-standards-verification.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Non-private artifact requiring selective adaptation. |
| shared_standards/audits/BOUNDARY_SUMMARY.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Non-private artifact requiring selective adaptation. |
| shared_standards/audits/source-docs-diff.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Non-private artifact requiring selective adaptation. |
| shared_standards/audits/standards-promotion-log.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Non-private artifact requiring selective adaptation. |
| shared_standards/audits/standards-reorg-audit.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Non-private artifact requiring selective adaptation. |
| shared_standards/private/auth/README.md | private/proprietary | split-private | TASK-059 (.hawp/work/active/TASK-059.md) | Private boundary path and sensitive/internal scope. |
| shared_standards/private/auth/protected-routes.md | private/proprietary | split-private | TASK-059 (.hawp/work/active/TASK-059.md) | Private boundary path and sensitive/internal scope. |
| shared_standards/private/context/guidelines.yml | private/proprietary | split-private | TASK-059 (.hawp/work/active/TASK-059.md) | Private boundary path and sensitive/internal scope. |
| shared_standards/private/context/patterns.yml | private/proprietary | split-private | TASK-059 (.hawp/work/active/TASK-059.md) | Private boundary path and sensitive/internal scope. |
| shared_standards/private/context/security.yml | private/proprietary | split-private | TASK-059 (.hawp/work/active/TASK-059.md) | Private boundary path and sensitive/internal scope. |
| shared_standards/private/internal-runtime/README.md | private/proprietary | split-private | TASK-059 (.hawp/work/active/TASK-059.md) | Private boundary path and sensitive/internal scope. |
| shared_standards/private/product/README.md | private/proprietary | split-private | TASK-059 (.hawp/work/active/TASK-059.md) | Private boundary path and sensitive/internal scope. |
| shared_standards/private/providers/README.md | private/proprietary | split-private | TASK-059 (.hawp/work/active/TASK-059.md) | Private boundary path and sensitive/internal scope. |
| shared_standards/private/security/security.md | private/proprietary | split-private | TASK-059 (.hawp/work/active/TASK-059.md) | Private boundary path and sensitive/internal scope. |
| shared_standards/project-specific/README.md | private/proprietary | split-private | TASK-060 (.hawp/work/active/TASK-060.md) | Project-specific domain references (Tekit/Mictlan/Zacatl) are private by default. |
| shared_standards/project-specific/mictlan/README.md | private/proprietary | split-private | TASK-060 (.hawp/work/active/TASK-060.md) | Project-specific domain references (Tekit/Mictlan/Zacatl) are private by default. |
| shared_standards/project-specific/tekit/README.md | private/proprietary | split-private | TASK-060 (.hawp/work/active/TASK-060.md) | Project-specific domain references (Tekit/Mictlan/Zacatl) are private by default. |
| shared_standards/project-specific/zacatl/README.md | private/proprietary | split-private | TASK-060 (.hawp/work/active/TASK-060.md) | Project-specific domain references (Tekit/Mictlan/Zacatl) are private by default. |
| shared_standards/public/context/README.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Public/shared governance or index artifact; keep shared or adapt selectively. |
| shared_standards/public/exports/hawp-absorbable/README.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Public/shared governance or index artifact; keep shared or adapt selectively. |
| shared_standards/public/exports/hawp-absorbable/manifest.json | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Public/shared governance or index artifact; keep shared or adapt selectively. |
| shared_standards/public/exports/machine-readable/legacy-standards-classification.csv | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Public/shared governance or index artifact; keep shared or adapt selectively. |
| shared_standards/public/exports/machine-readable/legacy-standards-classification.json | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Public/shared governance or index artifact; keep shared or adapt selectively. |
| shared_standards/public/exports/machine-readable/standards-boundary-classification.csv | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Public/shared governance or index artifact; keep shared or adapt selectively. |
| shared_standards/public/exports/machine-readable/standards-boundary-classification.json | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Public/shared governance or index artifact; keep shared or adapt selectively. |
| shared_standards/public/guidelines/README.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Public/shared governance or index artifact; keep shared or adapt selectively. |
| shared_standards/public/guidelines/architecture.md | public-safe | absorb | core/.hawp/kit/standards/guidelines/architecture.md (already absorbed) | Manifest-approved and already present in core standards tree. |
| shared_standards/public/guidelines/code-style.md | public-safe | absorb | core/.hawp/kit/standards/guidelines/code-style.md (already absorbed) | Manifest-approved and already present in core standards tree. |
| shared_standards/public/guidelines/documentation.md | public-safe | absorb | core/.hawp/kit/standards/guidelines/documentation.md (already absorbed) | Manifest-approved and already present in core standards tree. |
| shared_standards/public/guidelines/git-workflow.md | public-safe | absorb | core/.hawp/kit/standards/guidelines/git-workflow.md (already absorbed) | Manifest-approved and already present in core standards tree. |
| shared_standards/public/guidelines/mongodb-schema-design.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Public/shared governance or index artifact; keep shared or adapt selectively. |
| shared_standards/public/guidelines/testing.md | public-safe | absorb | core/.hawp/kit/standards/guidelines/testing.md (already absorbed) | Manifest-approved and already present in core standards tree. |
| shared_standards/public/standards/database/README.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Public/shared governance or index artifact; keep shared or adapt selectively. |
| shared_standards/public/standards/database/mongodb-schema-design.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Public/shared governance or index artifact; keep shared or adapt selectively. |
| shared_standards/public/standards/docs/README.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Useful but overlaps repo-local HAWP docs/templates and needs adaptation review. |
| shared_standards/public/standards/docs/hawp-install-update-safety.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Useful but overlaps repo-local HAWP docs/templates and needs adaptation review. |
| shared_standards/public/standards/nodejs/README.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Public/shared governance or index artifact; keep shared or adapt selectively. |
| shared_standards/public/standards/nodejs/area-composition.md | public-safe | absorb | core/.hawp/kit/standards/nodejs/area-composition.md (already absorbed) | Manifest-approved and already present in core standards tree. |
| shared_standards/public/standards/nodejs/build-and-env.md | public-safe | absorb | core/.hawp/kit/standards/nodejs/build-and-env.md (already absorbed) | Manifest-approved and already present in core standards tree. |
| shared_standards/public/standards/nodejs/code-style.md | public-safe | absorb | core/.hawp/kit/standards/nodejs/code-style.md (already absorbed) | Manifest-approved and already present in core standards tree. |
| shared_standards/public/standards/nodejs/git-workflow.md | public-safe | absorb | core/.hawp/kit/standards/nodejs/git-workflow.md (already absorbed) | Manifest-approved and already present in core standards tree. |
| shared_standards/public/standards/nodejs/project-structure.md | public-safe | absorb | core/.hawp/kit/standards/nodejs/project-structure.md (already absorbed) | Manifest-approved and already present in core standards tree. |
| shared_standards/public/standards/zacatl/README.md | private/proprietary | split-private | TASK-061 (.hawp/work/active/TASK-061.md) | Red-flag internal-only domain reference (Zacatl) requires private-lane review. |
| shared_standards/public/standards/zacatl/contract-testing.md | private/proprietary | split-private | TASK-061 (.hawp/work/active/TASK-061.md) | Red-flag internal-only domain reference (Zacatl) requires private-lane review. |
| shared_standards/public/standards/zacatl/dependency-registration.md | private/proprietary | split-private | TASK-061 (.hawp/work/active/TASK-061.md) | Red-flag internal-only domain reference (Zacatl) requires private-lane review. |
| shared_standards/public/standards/zacatl/evidence-linked-documentation.md | private/proprietary | split-private | TASK-061 (.hawp/work/active/TASK-061.md) | Red-flag internal-only domain reference (Zacatl) requires private-lane review. |
| shared_standards/public/standards/zacatl/handler-responsibilities.md | private/proprietary | split-private | TASK-061 (.hawp/work/active/TASK-061.md) | Red-flag internal-only domain reference (Zacatl) requires private-lane review. |
| shared_standards/public/standards/zacatl/layered-composition.md | private/proprietary | split-private | TASK-061 (.hawp/work/active/TASK-061.md) | Red-flag internal-only domain reference (Zacatl) requires private-lane review. |
| shared_standards/public/standards/zacatl/service-boundaries.md | private/proprietary | split-private | TASK-061 (.hawp/work/active/TASK-061.md) | Red-flag internal-only domain reference (Zacatl) requires private-lane review. |
| shared_standards/public/templates/ADR.template.md | repo-specific | adapt | TASK-062 (.hawp/work/active/TASK-062.md) | Useful but overlaps repo-local HAWP docs/templates and needs adaptation review. |

## Follow-up Work Items Created

- `TASK-059`: private standards lane review and retention policy
- `TASK-060`: project-specific domain boundaries (Tekit/Mictlan/Zacatl)
- `TASK-061`: public `standards/zacatl` red-flag boundary decision
- `TASK-062`: repo-specific adaptation lane for docs/template overlap
