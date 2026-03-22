# BUG-006 Validator External Root Diagnostics

Date: 2026-05-11
Owner: agent

## Scope

- Re-run validator from HAWP repo against local, Tekit, and Mictlan roots.
- Identify why external modern files were reported missing required close sections.
- Ensure FAIL/WARN output includes exact source file paths.

## Commands Run

```bash
cd librarian
npm run typecheck
npm run validate:workflow
npm run validate:workflow -- --hawp-root tekit/.hawp
npm run validate:workflow -- --hawp-root mictlan/.hawp
```

Focused diagnostic helper run:

```bash
cd librarian
node <<'NODE'
// inspected files:
// - tekit/.hawp/work/closed/2026/05/10/BUG-156-docs-audit-src-alignment.md
// - tekit/.hawp/work/closed/2026/05/11/BUG-157-email-verification-navigation-clarity.md
// - tekit/.hawp/work/closed/2026/05/11/BUG-158-auth-me-request-burst-on-login-redirect.md
// - mictlan/.hawp/work/closed/2026/05/10/TASK-033.md
NODE
```

## Root Cause

The closed-task-completeness validator only matched headings that started with exact `## ...`.
Downstream files like `BUG-156` and `TASK-033` used `### Outcome`, `### Verification`, and `### Close Checklist`, which should be accepted as valid section headings.

## Fix Implemented

- Section detection updated to match required headings at any markdown heading level (`#` to `######`).
- FAIL/WARN output now includes exact source path (`[source: /abs/path]`) for:
  - current untyped failures
  - missing-section failures
  - legacy untyped warnings
  - legacy missing-section warnings
- Added optional `--debug-closed-task` switch to print focused closed-task diagnostics for flagged files.

## Rerun Results Summary

### Local HAWP Repo

- `npm run validate:workflow` runs successfully.
- Result remains FAIL due existing backlog-consistency issue (`TASK-020` missing plan file).
- Closed-task warnings now include exact source paths.

### Tekit

- `BUG-156` no longer fails after heading-level fix.
- `BUG-157` and `BUG-158` still fail, and now report exact source paths.
- These two are true positives (required sections absent).

### Mictlan

- `TASK-033` no longer fails after heading-level fix.
- Validator overall status is PASS (warnings only).

## Focused Diagnostics (requested fields)

### BUG-156

- resolved absolute path: `tekit/.hawp/work/closed/2026/05/10/BUG-156-docs-audit-src-alignment.md`
- file exists: yes
- first matching heading lines: `### Outcome | ### Verification | ### Close Checklist`
- exact missing sections according to validator: `(none)`

### BUG-157

- resolved absolute path: `tekit/.hawp/work/closed/2026/05/11/BUG-157-email-verification-navigation-clarity.md`
- file exists: yes
- first matching heading lines: `(none)`
- exact missing sections according to validator: `Outcome, Verification, Close Checklist`

### BUG-158

- resolved absolute path: `tekit/.hawp/work/closed/2026/05/11/BUG-158-auth-me-request-burst-on-login-redirect.md`
- file exists: yes
- first matching heading lines: `(none)`
- exact missing sections according to validator: `Outcome, Verification, Close Checklist`

### TASK-033

- resolved absolute path: `mictlan/.hawp/work/closed/2026/05/10/TASK-033.md`
- file exists: yes
- first matching heading lines: `### Outcome | ### Verification | ### Close Checklist`
- exact missing sections according to validator: `(none)`
