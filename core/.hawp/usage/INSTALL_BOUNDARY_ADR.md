# ADR: Exclude Repo-Local Usage Artifacts from Default HAWP Installs

Date: 2026-04-27
Status: Proposed
Project: Human AI Workflow Protocol

## Context

HAWP is intended to ship a reusable public protocol kit. The `usage/` area in the source repository is useful for maintainers, but it can naturally accumulate repo-local operating notes, ADRs, and status artifacts.

Copying all of `.hawp/` into downstream repositories risks carrying source-repo maintenance context into projects where it does not belong.

## Decision

Adopt an install boundary:

- Keep `core/.hawp/usage/` in the source repository as maintainer/source-repo operating material.
- Exclude `usage/` from default install and update copy flows.
- Ship reusable protocol guidance from installable paths:
  - `README.md`
  - `START_HERE.md`
  - `SPEC.md`
  - `AUTHORING_PATTERNS.md`
  - `LICENSE`
  - `templates/`, `patterns/`, `reviews/`, `examples/`, `types/`
- Keep moving reusable usage concepts into installable folders over time.

## Rationale

This preserves maintainer context in the source repo while reducing leakage and confusion in downstream installs.

Benefits:

- lowers privacy and accidental-sharing risk
- keeps installed HAWP lean and generic
- clarifies reusable protocol files vs. source-repo maintenance artifacts

## Consequences

Required updates:

- install/update instructions must use an allowlist copy strategy
- downstream-facing overlay docs should not require `usage/` paths
- local/private `usage/` subpaths should be protected in `.gitignore`

Potential tradeoff:

- downstream users relying on `usage/` guides need installable equivalents (for example `START_HERE.md` and `templates/status-report.md`)

## Implementation Notes

- Default install and update flows skip `usage/`.
- Source-repo docs may still mention `usage/` if clearly labeled as source-repo maintenance material.
- Publication-safety checks remain active through installable review documents.
