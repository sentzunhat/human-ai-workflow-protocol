# Standards

Absorbed public standards distributed with the HAWP kit.

All files here are sourced from `shared_standards/public/exports/hawp-absorbable/manifest.json` — reviewed, portable, framework-agnostic, and public-safe.

## guidelines/

| File                                            | Topic                                                     |
| ----------------------------------------------- | --------------------------------------------------------- |
| [architecture.md](guidelines/architecture.md)   | Layered architecture, folder structure, module boundaries |
| [code-style.md](guidelines/code-style.md)       | TypeScript naming, formatting, import/export patterns     |
| [documentation.md](guidelines/documentation.md) | Comments, JSDoc, READMEs, changelogs                      |
| [git-workflow.md](guidelines/git-workflow.md)   | Commit format, branch naming, versioning                  |
| [security.md](guidelines/security.md)           | Input validation, secrets, auth, logging (non-negotiable) |
| [testing.md](guidelines/testing.md)             | Test framework setup, file organization, patterns         |

## nodejs/

| File                                                | Topic                                              |
| --------------------------------------------------- | -------------------------------------------------- |
| [area-composition.md](nodejs/area-composition.md)   | Feature-area module composition pattern            |
| [build-and-env.md](nodejs/build-and-env.md)         | Build pipeline and environment variable validation |
| [code-style.md](nodejs/code-style.md)               | Node.js-specific code style conventions            |
| [git-workflow.md](nodejs/git-workflow.md)           | Git workflow conventions for Node.js projects      |
| [project-structure.md](nodejs/project-structure.md) | Canonical project folder layout                    |

## database/

| File                            | Topic                                       |
| ------------------------------- | ------------------------------------------- |
| [README.md](database/README.md) | Database standards index and usage guidance |
| [nosql.md](database/nosql.md)   | NoSQL and document-model guidance           |
| [sql.md](database/sql.md)       | SQL and SQLite schema standards             |

## patterns/

| File                                                                | Topic                                              |
| ------------------------------------------------------------------- | -------------------------------------------------- |
| [evidence-discipline.md](patterns/evidence-discipline.md)           | Distinguish direct evidence from inference in docs |
| [parallel-work-guardrails.md](patterns/parallel-work-guardrails.md) | Coordination patterns for parallel work            |
| [non-findings.md](patterns/non-findings.md)                         | Audit and review patterns for bounded audits       |

## service-design/

| File                                                                                | Topic                                                  |
| ----------------------------------------------------------------------------------- | ------------------------------------------------------ |
| [README.md](service-design/README.md)                                               | Service-design standards index and extraction boundary |
| [service-boundaries.md](service-design/service-boundaries.md)                       | Service contracts, ownership, and boundary discipline  |
| [handler-responsibilities.md](service-design/handler-responsibilities.md)           | Handler orchestration vs domain responsibilities       |
| [dependency-composition.md](service-design/dependency-composition.md)               | Deterministic dependency composition and validation    |
| [layered-composition.md](service-design/layered-composition.md)                     | Layer-preserving composition and adapter replacement   |
| [contract-testing.md](service-design/contract-testing.md)                           | Contract-focused behavior testing                      |
| [evidence-linked-documentation.md](service-design/evidence-linked-documentation.md) | Documentation tied to implementation evidence          |

## Consolidation Map

Use this map to keep standards absorb/sync decisions simple and safe.

### Canonical standards (normative)

- `guidelines/architecture.md`
- `guidelines/code-style.md`
- `guidelines/documentation.md`
- `guidelines/git-workflow.md`
- `guidelines/security.md`
- `guidelines/testing.md`
- `nodejs/area-composition.md`
- `nodejs/build-and-env.md`
- `nodejs/code-style.md`
- `nodejs/git-workflow.md`
- `nodejs/project-structure.md`
- `database/sql.md`
- `database/nosql.md`
- `patterns/evidence-discipline.md`
- `patterns/parallel-work-guardrails.md`
- `patterns/non-findings.md`
- `service-design/service-boundaries.md`
- `service-design/handler-responsibilities.md`
- `service-design/dependency-composition.md`
- `service-design/layered-composition.md`
- `service-design/contract-testing.md`
- `service-design/evidence-linked-documentation.md`

### Supporting indexes and planning aids (non-normative)

- `database/README.md` (index and usage guidance)
- `.hawp/kit/guidance/da-schema-planning.md` (planning checklist and review gates)
- `.hawp/work/evidence/db-decision-template.md` (decision evidence template)

### Absorption rule of thumb

- Public-safe, portable rules go under `core/.hawp/kit/standards/**`.
- Workflow-specific execution checklists stay in `.hawp/kit/guidance/**`.
- Evidence capture templates stay in `.hawp/work/evidence/**`.

## Exclusion Policy (Private and Proprietary Content)

The `core/.hawp/kit/standards/**` tree is public-safe only.

- Do not absorb private or proprietary standards into this tree.
- Do not include secret sauce recipes, internal architecture playbooks, or internal implementation runbooks.
- Treat internal-domain references (for example Tekit, Micltan, Zacatl, or similar internal-only systems) as private by default.
- If any item is unclear, classify it as private/proprietary and split it into a separate work-item plan under `.hawp/work/active/`.
- Keep private or organization-specific standards in work-project planning lanes until explicitly cleared for public-safe export.
