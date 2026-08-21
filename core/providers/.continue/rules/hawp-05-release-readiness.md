---
name: "HAWP Release Readiness"
description: "HAWP release, dependency, CI, and publication safety"
globs:
  - "package.json"
  - "package-lock.json"
  - ".github/**"
  - "docs/**"
  - ".hawp/**"
  - "**/Dockerfile"
alwaysApply: false
---

<!-- Generated from core/providers/shared/behaviors — edit shared source and run npm --prefix librarian run providers:sync -->

# HAWP Release Readiness

Use this behavior when preparing a release, consolidating dependency updates,
or reviewing CI and publication safety.

## Release candidate discipline

- Inspect the current branch, worktree, package manifests, lockfiles, release
  scripts, and CI workflows before editing.
- Consolidate compatible dependency updates in the intended work branch. Do
  not delete, merge, or close automation branches without explicit approval.
- Keep major-runtime or peer-graph migrations separate when the compatibility
  evidence is incomplete; record the deferred update and its reason.
- Treat the changelog, package version, lockfile, and generated package
  metadata as one release surface.

## Verification evidence

Match verification to the release surface. When applicable, run and record:

- clean install from the lockfile;
- production and full dependency audits, distinguishing findings owned by
  optional or fixture-only consumers from the production graph;
- tests, type checks, lint, and ESM/CJS or equivalent builds;
- packed consumer smoke tests for each supported integration;
- database and Docker smoke tests for every supported example;
- prepublish/version guards and package staging checks.

Do not call a check passed when it was skipped, ran under an unsupported
runtime, or only covered a subset of the release matrix. Label those limits
and keep them as release gates.

## Secure publication

- Prefer the package registry's trusted OIDC publishing flow over long-lived
  publish tokens when the registry and CI provider support it.
- For GitHub Actions OIDC publishing, require the narrow `id-token: write`
  permission, configure the exact repository/workflow trust relationship, and
  verify provenance on a controlled release path.
- Do not describe OIDC migration as complete until the workflow and registry
  trust are configured and tested. Revoke superseded long-lived credentials
  only as an explicitly approved operational action.
- A release candidate may be prepared and verified without merging, tagging,
  publishing, or creating an external release. Require owner approval before
  those state-changing actions.

Keep direct verification separate from inference, and store the checkpoint in
the repository's normal HAWP status/evidence paths.
