# TASK-030 v0.2 File-Tracking Enforcement Rules - File Tracking

**Backlog ID:** TASK-030

**Purpose:** Enforce plain-text path evidence and path-locked reporting for this lane.

**Work Item:** `.hawp/work/active/TASK-030.md`

**Last Updated:** 2026-05-12

---

### Lane Scope (plain text exact repo-relative paths)

```txt
.hawp/work/active/TASK-030-files.md
core/.hawp/kit/instructions/da-file-tracking.md
core/.hawp/kit/references/work-item-file-tracking.md
core/.hawp/kit/templates/work-item-files.md
```

### Exclusions (plain text exact repo-relative paths)

```txt
.hawp/work/active/0008-CLARIFICATION-exact-paths.md
.hawp/work/active/0008-install-update-distribution-review.md
shared_standards/
```

### Owned Files

- `.hawp/work/active/TASK-030-files.md`
- `core/.hawp/kit/instructions/da-file-tracking.md`
- `core/.hawp/kit/references/work-item-file-tracking.md`
- `core/.hawp/kit/templates/work-item-files.md`

### Read-Only Context Files

- `.hawp/work/active/TASK-030.md`
- `.hawp/work/BACKLOG.md`

### Do-Not-Touch Files

- `.hawp/work/active/0008-CLARIFICATION-exact-paths.md`
- `.hawp/work/active/0008-install-update-distribution-review.md`
- `shared_standards/`

### Locked / Reserved Files

- `core/.hawp/kit/instructions/da-file-tracking.md`
- `core/.hawp/kit/references/work-item-file-tracking.md`
- `core/.hawp/kit/templates/work-item-files.md`
- `.hawp/work/active/TASK-030-files.md`

### Changed Files

- `.hawp/work/active/TASK-030-files.md`
- `core/.hawp/kit/instructions/da-file-tracking.md`
- `core/.hawp/kit/references/work-item-file-tracking.md`
- `core/.hawp/kit/templates/work-item-files.md`

### Verification Notes

```bash
git status --short
git diff --name-status
git diff --name-only
git diff --check
npm --prefix librarian run typecheck
npm --prefix librarian run validate:workflow
```

### Path Evidence Rule

Clickable editor, IDE, or GitHub Copilot file references do not count as path evidence.
Path evidence must be plain text exact repo-relative paths from repository root.

If any visible path appears as basename-only, treat it as unsafe evidence and classify as:

- `INVALID_REPO_RELATIVE_PATH`
- `BASENAME_ONLY_REFERENCE`
- `SELF_VALIDATION_FAILURE`

## Outcome

Normalized as historical file-tracking artifact; retained original content and added required close sections for validator compatibility.

## Verification


### Evidence Follow-Up

- [ ] Research evidence for: Required sections added for compatibility\n- [x] Original artifact body preserved
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.

- [x] Required sections added for compatibility\n- [x] Original artifact body preserved

## Close Checklist

- [x] Outcome section filled\n- [x] Verification section filled\n- [x] Historical artifact preserved
