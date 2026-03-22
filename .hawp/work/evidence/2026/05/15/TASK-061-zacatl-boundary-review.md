# TASK-061 Zacatl Boundary Review

Generated: 2026-05-15

## Scope

- Reviewed paths: `shared_standards/public/standards/zacatl/*.md`
- Rubric: `.hawp/kit/guidance/shared-standards-review-rubric.md`
- Goal: assign per-file boundary decision (`private/proprietary` vs `repo-specific` adaptation)

## Summary Decision

- Total files reviewed: 7
- `private/proprietary` (split-private): 1
- `repo-specific` (adapt): 6
- `public-safe` (direct absorb): 0

Rationale:

- The six standards pages are principle-level and explicitly exclude internal implementation details.
- The Zacatl category `README.md` includes internal-domain/source references and private/internal-only note categories, so it remains private-lane.

## Per-File Decisions

| file                                                                        | evidence snippet                                                                                                            | classification        | action          | destination | rationale                                                                                   |
| --------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- | --------------------- | --------------- | ----------- | ------------------------------------------------------------------------------------------- |
| `shared_standards/public/standards/zacatl/README.md`                        | "Private/Internal-Only Notes"; "derived from ... Mictlan and Tekit implementations"; source pointer `new_docs/zacatl_docs/` | `private/proprietary` | `split-private` | `TASK-060`  | Contains internal-domain references and private/internal notes (red-flag category).         |
| `shared_standards/public/standards/zacatl/service-boundaries.md`            | "Services expose clear contracts and should not leak internal runtime details."                                             | `repo-specific`       | `adapt`         | `TASK-062`  | Principle-level guidance; no private implementation details required.                       |
| `shared_standards/public/standards/zacatl/handler-responsibilities.md`      | "Does Not Include ... Internal routing bootstraps ... Project route inventories"                                            | `repo-specific`       | `adapt`         | `TASK-062`  | Explicitly excludes private/project internals; candidate for framework-agnostic extraction. |
| `shared_standards/public/standards/zacatl/dependency-registration.md`       | "Does Not Include ... Internal DI container mechanics ... Project-specific token maps"                                      | `repo-specific`       | `adapt`         | `TASK-062`  | Principle-level registration guidance with explicit internal exclusions.                    |
| `shared_standards/public/standards/zacatl/layered-composition.md`           | "Does Not Include ... Framework-specific bootstraps ... Runtime vendor lock-in assumptions"                                 | `repo-specific`       | `adapt`         | `TASK-062`  | Architecture principles are portable but still Zacatl-labeled.                              |
| `shared_standards/public/standards/zacatl/contract-testing.md`              | "Does Not Include ... Framework internals testing ... Project-specific endpoint snapshots"                                  | `repo-specific`       | `adapt`         | `TASK-062`  | Behavior-first contract testing guidance can be generalized.                                |
| `shared_standards/public/standards/zacatl/evidence-linked-documentation.md` | "Does Not Include ... Internal-only security details"                                                                       | `repo-specific`       | `adapt`         | `TASK-062`  | Documentation principle appears reusable with naming adaptation.                            |

## Risk Note

No Zacatl file is approved for direct absorb in current form. Any reuse must go through adaptation/extraction lane (`TASK-062`) and preserve exclusion policy boundaries.
