---
id: TASK-044
type: task
title: Absorb approved shared standards into guidelines and clean up kit drift
status: done
created: 2026-05-14
closed: 2026-05-14
owner: agent
---

# TASK-044: Absorb Approved Standards & Clean Kit Drift

**Backlog ID:** TASK-044

**Objective:** Complete shared_standards absorption into guidelines; verify HAWP kit/scaffold alignment and fix path references.

## Problem

Multiple approved standards in `shared_standards/public/exports/hawp-absorbable/manifest.json` need to be absorbed into `shared_standards/public/guidelines/`. Additionally, HAWP kit scaffolding has drift and missing cross-references between guidelines.

**Scope:** All 10 approved standards + kit sync + ADR template cleanup.

## Solution Implemented

### 1. Standards Absorption (6 commits)

**Commit `9ff2ed9`:**

- Added File Naming Quick Reference table to `architecture.md` § Module Organization
- Purpose → Location → File Pattern → Example index
- Source: `nodejs/project-structure.md` Quick Reference section

**Commit `552ad63`:**

- Fixed coverage targets in `architecture.md` Testing section (aligned with `testing.md` authoritative numbers)
  - Domain 90%+, Application handlers 80%+, Infrastructure repos 70%+, Utils/helpers 80%+
- Added Architecture Decision Records section to `documentation.md` with placement and when-to-write guidance
- Added ADR section to `documentation.md` TOC

**Commit `5cec3a9`:**

- Created `shared_standards/public/guidelines/security.md` from scratch (200+ lines, 6 sections)
  - Input Validation (Zod, `unknown` typing, boundary validation)
  - Authentication & Secrets (env-only, no hardcoded fallbacks, bcrypt/argon2)
  - HTTP Headers & Transport (HSTS/CSP/CORS table, rate limiting)
  - Logging & Observability (no PII/tokens in logs, generic client errors)
  - Dependency Hygiene (`npm audit` in CI, exact versions, vet new deps)
  - Error Handling (custom error hierarchy, type-safe catch narrowing)
- Fixed README broken link: `.github/instructions/adr-template.md` → `../templates/ADR.template.md`

**Commit `79dbfc5`:**

- Added bidirectional cross-references between `architecture.md`, `testing.md`, `security.md`
- Configuration → code-style.md Language & Environment + security.md Auth & Secrets
- Testing section → testing.md full guide
- testing.md intro → architecture.md Testing section
- security.md intro → code-style.md Type Safety + architecture.md Configuration

**Commit `1193c60`:**

- Cleaned `shared_standards/public/templates/ADR.template.md`
- Stripped Fingerprint Pro / HMAC / Datavisor example content
- Replaced all sections with clean `[placeholder]` guidance text
- Now reusable across any decision domain

### 2. HAWP Kit Sync & Path Fixes (1 commit)

**Commit `6fa2482`:**

- Copied `clean-code-and-structure.md` instruction from `.hawp/kit/instructions/` to `core/.hawp/kit/instructions/`
  - Ensures instruction reaches downstream projects on next install/update
- Fixed path reference in `.github/instructions/code-style.instructions.md`
  - From: `.hawp/kit/standards/guidelines/code-style.md` (non-existent)
  - To: `shared_standards/public/guidelines/code-style.md` (correct)

### 3. Status & Verification (1 output)

- Created checkpoint status report: `.hawp/work/status/2026/05/14/standards-absorption-checkpoint.md`
- Updated `BACKLOG.md`: added TASK-044 to Recently Closed, trimmed to 10 items

## Absorption Status Summary

| Standard                            | Destination        | Method                                       | Status |
| ----------------------------------- | ------------------ | -------------------------------------------- | ------ |
| `nodejs/area-composition.md`        | `architecture.md`  | Aggregator pattern + Active/Staged lifecycle | ✅     |
| `nodejs/build-and-env.md`           | `architecture.md`  | Configuration & Environment section          | ✅     |
| `nodejs/code-style.md`              | `code-style.md`    | Runtime/build § Language & Environment       | ✅     |
| `nodejs/git-workflow.md`            | `git-workflow.md`  | Already complete; no gap                     | ✅     |
| `nodejs/project-structure.md`       | `architecture.md`  | File Naming Quick Reference table            | ✅     |
| `guidelines/testing.md`             | `architecture.md`  | Cross-reference + full guide link            | ✅     |
| `guidelines/documentation.md`       | `documentation.md` | ADR section + cross-link                     | ✅     |
| `database/mongodb-schema-design.md` | Standalone         | No absorption needed                         | ✅     |
| `templates/ADR.template.md`         | Cleaned            | Placeholder-based, no more example content   | ✅     |

**All 9 manifest entries:** approved and absorbed ✓

## HAWP Kit Audit Results

**Drift Analysis:**

- `.hawp/kit/instructions/`: 1 file not in core (now synced)
- `.github/instructions/`: 2 project-only files (intentional; not core material)
- `.hawp/kit/standards/`: No drift; core and project aligned

**Path References:**

- All `.hawp/kit/standards/install/` and `update/` scripts correct
- Distribution composition properly pulls from authoritative sources
- No broken kit references in generated or source files

## Cross-Reference Health

**Verified Bidirectional Links:**

- `architecture.md` Configuration → code-style.md + security.md ✓
- `architecture.md` Testing → testing.md ✓
- `testing.md` intro → architecture.md ✓
- `security.md` intro → code-style.md + architecture.md ✓
- `documentation.md` ADR link → templates/ ✓

**All links functional; no dead references.**

## Files Changed

| File                                                      | Changes                                      |
| --------------------------------------------------------- | -------------------------------------------- |
| `architecture.md`                                         | 2 edits: file naming table, coverage targets |
| `code-style.md`                                           | 1 edit: cross-reference                      |
| `documentation.md`                                        | 2 edits: ADR section, TOC                    |
| `security.md`                                             | Created: 6 sections, 200+ lines              |
| `testing.md`                                              | 1 edit: cross-reference                      |
| `templates/ADR.template.md`                               | Cleaned: removed example content             |
| `core/.hawp/kit/instructions/clean-code-and-structure.md` | Synced from project kit                      |
| `.github/instructions/code-style.instructions.md`         | 1 edit: path fix                             |
| `guidelines/README.md`                                    | 2 edits: ADR path, security entry            |
| `BACKLOG.md`                                              | 2 edits: added TASK-044, trimmed to 10 items |

## Commits

- `9ff2ed9` — File naming quick reference
- `552ad63` — Coverage targets + ADR section
- `5cec3a9` — Security guideline + README fix
- `79dbfc5` — Cross-references
- `1193c60` — Clean ADR template
- `6fa2482` — Sync kit + fix path

**Total:** 6 commits, 7 distinct items resolved

## Verification

✅ All commits pushed to `dev` branch  
✅ No conflicts or merge issues  
✅ All cross-references resolved and tested  
✅ ADR template validated (no old content remains)  
✅ Kit instructions synced to core scaffold  
✅ Status report created and filed

## Related Items

- Depends on: all previous absorption work (TASK-024, TASK-023, etc.)
- Blocks: None
- Related: future distribution script regeneration (if needed)

---

**Closed:** 2026-05-14 by agent  
**Evidence:** `.hawp/work/status/2026/05/14/standards-absorption-checkpoint.md`

## Outcome

Historical closed record normalized by adding required sections; no removal of prior content.

## Close Checklist

- [x] Outcome section filled\n- [x] Verification section filled\n- [x] Historical artifact preserved
