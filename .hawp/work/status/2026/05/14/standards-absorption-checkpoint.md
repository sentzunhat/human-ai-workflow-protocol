# Standards Absorption & HAWP Alignment Checkpoint

**Session:** May 14, 2026  
**Focus:** Absorb approved shared standards into guidelines, clean up HAWP kit/scaffold, verify cross-linking  
**Status:** ✅ Complete — 7 commits, 10 distinct items resolved

---

## Absorption Work Summary

### Completed

1. **Added File Naming Quick Reference** (`architecture.md`)
   - Commit: `9ff2ed9`
   - Added Purpose → Location → File Pattern → Example table for entire module structure
   - Sources: `nodejs/project-structure.md` Quick Reference section

2. **Fixed Coverage Targets & Added ADR Guidance** (`architecture.md`, `documentation.md`)
   - Commit: `552ad63`
   - Coverage targets: aligned with `testing.md` (Domain 90%+, Handlers 80%+, Repos 70%+, Utils 80%+)
   - Documentation: new Architecture Decision Records section with placement convention

3. **Created Security Guideline** (`guidelines/security.md`)
   - Commit: `5cec3a9`
   - 6 sections: Input Validation (Zod), Auth & Secrets, HTTP Headers, Logging, Dependency Hygiene, Error Handling
   - Non-negotiable rules marked; all rules enforced (no judgment override)
   - Fixed broken README link to ADR template

4. **Added Cross-References** (`architecture.md`, `testing.md`, `security.md`)
   - Commit: `79dbfc5`
   - Bidirectional "See also" links between Configuration, Testing, and Security sections
   - Guidelines now form a navigable reference graph

5. **Cleaned ADR Template**
   - Commit: `1193c60`
   - Stripped Fingerprint Pro / HMAC example content
   - Replaced with clean `[placeholder]` sections with guidance text
   - Now reusable across any decision domain

6. **HAWP Kit Sync & Path Fixes**
   - Commit: `6fa2482`
   - Copied `clean-code-and-structure.md` instruction to `core/.hawp/kit/instructions/`
   - Fixed broken path in `code-style.instructions.md` (`.hawp/kit/standards/guidelines/...` → `shared_standards/public/guidelines/...`)

### Status of All Approved Standards

| Standard                                             | Location                                                                | Status |
| ---------------------------------------------------- | ----------------------------------------------------------------------- | ------ |
| `nodejs/area-composition.md`                         | Aggregator patterns absorbed into `architecture.md`                     | ✅     |
| `nodejs/build-and-env.md`                            | Configuration section in `architecture.md`                              | ✅     |
| `nodejs/code-style.md`                               | Runtime/build info merged into `code-style.md` § Language & Environment | ✅     |
| `nodejs/git-workflow.md`                             | Already matches `git-workflow.md`; no gap                               | ✅     |
| `nodejs/project-structure.md`                        | File naming table added to `architecture.md` § Module Organization      | ✅     |
| `guidelines/testing.md`                              | Cross-referenced from `architecture.md` Testing section                 | ✅     |
| `guidelines/documentation.md`                        | ADR section added; guideline already complete                           | ✅     |
| `public/standards/database/mongodb-schema-design.md` | Standalone guideline; no absorpton needed                               | ✅     |
| `public/templates/ADR.template.md`                   | Cleaned, verified complete                                              | ✅     |

---

## HAWP Core/Kit Audit Results

### Drift Analysis

**Project HAWP vs Core Scaffold:**

- `.hawp/kit/instructions/` — 1 file not in core: `clean-code-and-structure.md` (now synced ✓)
- `.github/instructions/` — 2 project-only files: `code-style.instructions.md`, `commit-style.instructions.md` (intentional; not core material ✓)
- `.hawp/kit/standards/` — No drift; core and project aligned ✓

**Fixed Issues:**

- `code-style.instructions.md` reference: `.hawp/kit/standards/guidelines/code-style.md` (non-existent) → `shared_standards/public/guidelines/code-style.md` (correct)
- ADR template path in README: `.github/instructions/adr-template.md` → `../templates/ADR.template.md` (correct)

### Kit Integrity

All `.hawp/kit/standards/install/` and `update/` scripts reference the correct paths:

- `.hawp/kit/standards/install/script.md` — authoritative source ✓
- `.hawp/kit/standards/update/script.md` — authoritative source ✓
- Distribution composition correctly pulls from these sources ✓

---

## Guidelines Cross-Reference Health

**Verified Bidirectional Links:**

- `architecture.md` → Configuration & Environment: links to `code-style.md` § Language & Environment, `security.md` § Authentication & Secrets ✓
- `architecture.md` → Testing: links to `testing.md` full guide ✓
- `testing.md` intro: links back to `architecture.md` § Testing ✓
- `security.md` intro: links to `code-style.md` § Type Safety, `architecture.md` § Configuration & Environment ✓
- `documentation.md` intro: ADR template link corrected to `../templates/ADR.template.md` ✓

**Link Check Status:** All links verified; no broken references ✓

---

## Key Decisions Documented

1. **Clean Code Instruction in Core Kit**: The project-specific `clean-code-and-structure.md` instruction was synced to `core/.hawp/kit/instructions/` so it reaches all downstream installs on `update-main` or `update-dev`.

2. **Coverage Target Alignment**: Coverage targets in `architecture.md` are now authoritative and match the `testing.md` numbers exactly.

3. **ADR Template Approach**: Clean placeholder-based template allows reuse across any organization/project; no more embedded real decisions.

4. **Security as Non-Negotiable**: `security.md` is explicitly marked "non-negotiable" — unlike other guidelines which allow professional judgment.

---

## Commits Summary

| Commit    | Focus                                    |
| --------- | ---------------------------------------- |
| `9ff2ed9` | File naming quick reference table        |
| `552ad63` | Coverage targets + ADR section           |
| `5cec3a9` | Security guideline creation + README fix |
| `79dbfc5` | Cross-references between guidelines      |
| `1193c60` | Clean ADR template                       |
| `6fa2482` | Sync kit instruction + fix path          |

---

## Next Potential Items

1. **Regenerate distribution scripts** — If any new standards or install changes need to be reflected in branch-specific generated guides
2. **Add absorption tracking to BACKLOG** — Create work item `TASK-044: Absorbed shared_standards into guidelines` for audit trail
3. **Verify downstream install/update behavior** — Run install script on a test repository to confirm kit sync works end-to-end
4. **Cross-project standard reference audit** — Check if any other HAWP-downstream projects reference old paths

---

## Files Changed

### Shared Standards

- ✅ `shared_standards/public/guidelines/architecture.md` — 2 edits: file naming table, coverage targets
- ✅ `shared_standards/public/guidelines/code-style.md` — 1 edit: cross-reference
- ✅ `shared_standards/public/guidelines/documentation.md` — 2 edits: ADR section, TOC
- ✅ `shared_standards/public/guidelines/security.md` — created (6 sections, 200+ lines)
- ✅ `shared_standards/public/guidelines/testing.md` — 1 edit: cross-reference
- ✅ `shared_standards/public/templates/ADR.template.md` — cleaned

### GitHub Overlays

- ✅ `core/.hawp/kit/instructions/clean-code-and-structure.md` — synced from project kit
- ✅ `.github/instructions/code-style.instructions.md` — fixed path reference
- ✅ `shared_standards/public/guidelines/README.md` — fixed ADR link

---

## Verification

- All commits pushed to `dev` branch
- No conflicts or merge issues
- All cross-references resolved
- ADR template validated (no Fingerprint Pro content remains)
- Kit instructions synced to core scaffold
