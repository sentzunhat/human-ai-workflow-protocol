# TASK-079: Development audit and remediation plan

**Backlog ID:** TASK-079
**UUID:** 533dc506-e376-4763-927e-a5ad9e8f25e8
**Status:** done
**Created:** 2026-06-20
**Owner:** Codex

## HAWP Shape

```text
input: |
  can you please review and audit the human ai workflow repo from development and tell me or add a new work item using hawp to review and audit and do the work of ti and writting a report of what needs to be done for

context: |
  This repository uses HAWP as a lightweight workflow method. Development work is centered on the librarian TypeScript tooling, generated distribution guides, provider materialization, and HAWP workflow records.

mission: |
  Review the repository from a development-maintenance perspective, create a HAWP-tracked work item, run available verification, and write a compact report of what needs to be done next.

constraints: |
  Keep the audit bounded to development workflow, validation, tooling, generated artifacts, and repo-local HAWP guidance.
  Do not modify unrelated user changes already present in the worktree.
  Separate direct evidence from inference.
  Do not turn HAWP into a runtime engine, compiler, validator, orchestrator, or memory system.

output: |
  A backlog row, this active work item, and a status/audit report under .hawp/work/status/YYYY/MM/DD/ with verified findings and recommended next work.
```

## Inspected Scope

- Repo root proof: `pwd`, `git rev-parse --show-toplevel`, `git rev-parse --show-prefix`, `git status --short`.
- HAWP operating guide: `.hawp/kit/start-here.md`.
- Status report guide: `.hawp/kit/usage/status-report.md`.
- Backlog: `.hawp/work/BACKLOG.md`.
- Project overview: `README.md`, `librarian/README.md`, `librarian/scripts/README.md`.
- Tooling configuration: `.nvmrc`, `librarian/package.json`, `librarian/tsconfig.json`, `.github/workflows/*.yml`.
- Development scripts and tests under `librarian/scripts/`.
- Review standards: `.hawp/kit/reviews/project-review-checklist.md`, `.hawp/kit/standards/guidelines/testing.md`, `.hawp/kit/standards/nodejs/project-structure.md`.

## Findings

### Confirmed

- The default shell Node runtime was `v18.12.1`, while `.nvmrc` is `26` and `librarian/package.json` requires `node >=26`.
- `npm test` fails under the default Node 18 runtime with `node: bad option: --import`.
- Under Node `v26.3.1`, `npm test` passes: 37 tests, 37 passing.
- Under Node `v26.3.1`, `npm run typecheck`, `npm run distribution:sync`, `npm run provider:sync`, `npm run workflow:validate`, and `npm run validate` pass.
- `npm audit --audit-level=moderate` reports 0 vulnerabilities for `librarian/`.
- `npm outdated --json` reports no outdated dependencies for `librarian/`.
- Workflow validation passes with one tolerated pre-cutoff legacy closed-file warning.
- Generated distribution guides and provider overlays are current; `distribution:sync` updated 0 files.
- README and `librarian/README.md` now tell contributors to activate Node 26 with `nvm use` before running librarian commands.
- The canonical and mirror code-style guides now explain that Node `24.14.0` is a portable standards baseline, while individual repositories may require newer runtimes.

### Likely

- The remaining workflow validation warning is expected legacy-ID drift, not a path-normalization issue.

## Recommended Remediation

1. Keep treating the remaining workflow warning as historical legacy-ID debt unless a maintainer wants a dedicated cleanup pass for old closed records.
2. If future standards work raises the portable runtime floor above Node `24.14.0`, update the canonical and mirror code-style guides together.
3. Keep validating on the repository-declared runtime first whenever local failures appear under older system Node versions.

## Verification

- `npm run typecheck` passed under Node `v26.3.1`.
- `npm test` passed under Node `v26.3.1`: 37/37 tests passing.
- `npm run distribution:sync` passed and reported 0 generated file updates.
- `npm run provider:sync` passed and reported 0 materialized provider file updates.
- `npm run workflow:validate` passed with one tolerated legacy-ID warning and 9 valid evidence links.
- `npm run validate` passed under Node `v26.3.1`.
- `npm audit --audit-level=moderate` reported 0 vulnerabilities.
- `npm outdated --json` returned `{}`.

## Outcome

Development audit completed, remediation applied, and workflow records aligned for closure.

## Close Checklist

- [x] HAWP work item created.
- [x] Audit report written.
- [x] Verification commands run with correct Node runtime.
- [x] Follow-up remediation implemented or split into smaller work items.
