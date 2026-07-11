# TASK-030: Implement work-item file tracking v0.1 — File Tracking

**Backlog ID:** TASK-030

**Purpose:** Track files related to TASK-030 using exact repo-relative paths.

**Work Item:** `.hawp/work/active/TASK-030.md`

**Last Updated:** 2026-05-12

---

### Owned Files

Files this task is allowed to create or edit.

- `.hawp/work/active/TASK-030.md` — this work item plan
- `.hawp/work/active/TASK-030-files.md` — this file tracking document
- `core/.hawp/kit/templates/work-item-files.md` — reusable file tracking template
- `core/.hawp/kit/instructions/da-file-tracking.md` — DA instruction reference
- `core/.hawp/kit/references/work-item-file-tracking.md` — general reference for file tracking

---

### Read-Only Context Files

Files this task may read for context but must not edit.

- `.hawp/work/BACKLOG.md` — backlog structure (will update, but only to add this row)
- `core/.hawp/kit/references/backlog-alignment.md` — for understanding backlog conventions
- `core/.hawp/kit/usage/intake-workflow.md` — for understanding work item structure
- `core/.hawp/kit/templates/intake-plan.md` — for template patterns
- `README.md` — project context
- `package.json` — for understanding npm tasks

---

### Do-Not-Touch Files

Files explicitly out of scope for TASK-030.

- `librarian/scripts/` — separate CLI work (TASK-027, TASK-028, TASK-029)
- `core/install.md` — being audited by TASK-026
- `core/update.md` — being audited by TASK-026
- `.github/instructions/` — separate system instruction files
- `shared_standards/` — separate standards system

---

### Locked / Reserved Files

Files currently reserved for this task to avoid parallel-agent conflicts.

- `.hawp/work/BACKLOG.md` — will add one row for TASK-030 in Active Work table

---

### Changed Files

Files actually changed during this task. Updated as work progresses.

**Core files created:**

- `core/.hawp/kit/templates/work-item-files.md` (new)
- `core/.hawp/kit/instructions/da-file-tracking.md` (new)
- `core/.hawp/kit/references/work-item-file-tracking.md` (new)

**Work item files created:**

- `.hawp/work/active/TASK-030.md` (this work item)
- `.hawp/work/active/TASK-030-files.md` (this file)

**Backlog updated:**

- `.hawp/work/BACKLOG.md` (added TASK-030 to Active Work table)

---

### Verification Notes

Commands to verify accuracy before closing:

```bash
# Show new/modified files
git status --short

# Show exact paths changed
git diff --name-status

# Check for trailing whitespace or formatting issues
git diff --check

# Verify no TypeScript errors (if applicable)
npm run typecheck

# Verify HAWP workflow structure is valid
npm run validate:workflow

# Show final commit
git log -1 --oneline
```

Expected output:

- 5 new files created
- All paths repo-relative and full
- No formatting issues
- All verification commands pass

## Outcome

Normalized as historical file-tracking artifact; retained original content and added required close sections for validator compatibility.

## Verification


### Evidence Follow-Up

- [ ] Research evidence for: Required sections added for compatibility\n- [x] Original artifact body preserved
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.

- [x] Required sections added for compatibility\n- [x] Original artifact body preserved

## Close Checklist

- [x] Outcome section filled\n- [x] Verification section filled\n- [x] Historical artifact preserved
