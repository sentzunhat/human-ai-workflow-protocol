# HAWP Install and Update Safety (Reference)

## Purpose

This reference defines safety principles for HAWP install and update work.

It is principle-level guidance only. It does not provide operational scripts,
copy-paste commands, or distribution generation procedures.

## Scope

Applies to changes that introduce, refresh, or align HAWP guidance assets while
preserving project-owned work.

This reference covers:

- safety boundaries
- ownership boundaries
- change-scoping expectations
- verification expectations

This reference does not cover:

- install/update command implementation details
- distribution source generation details
- branch-specific rollout mechanics

## Safety Principles

### 0) Treat install/update as explicit execution work items

Each install or update run should be treated as a single execution work item for
the target repository.

Preflight mode checks are required:

- install mode is for repos where `.hawp/` is not present
- update mode is for repos where `.hawp/` already exists
- if `.hawp/` exists, prefer update guidance for the selected branch
- if `.hawp/` is missing, run install first before update

This keeps one-shot execution predictable for both human and software agents.

### 1) Project-owned work must not be clobbered

Project records under work are project-owned and must be preserved.

Never overwrite active, closed, evidence, status, or decisions records as part
of install/update-aligned changes.

### 2) Reusable guidance and project records are different ownership classes

Treat reusable guidance assets and project records as separate domains.

- reusable guidance belongs to kit-level protocol assets
- project records belong to work-level execution history

Changes in one domain must not implicitly rewrite the other domain.

### 3) Path-scoped changes are mandatory

Every install/update-related change must be explicitly path-scoped.

Only modify files that are in the approved task scope. If a needed edit falls
outside the approved scope, stop and request explicit approval.

### 4) Distribution updates are a separate lane

Reference-level safety guidance and distribution implementation are distinct
work lanes.

Do not mix reference updates with distribution-source or generated-output edits
in the same chunk unless the task is explicitly distribution-scoped.

### 5) Preserve user and project files by default

Prefer additive or no-clobber behavior for files that are project-owned.

When uncertainty exists about ownership, treat the file as protected until the
owner confirms otherwise.

### 6) Verify before close

Before closing an install/update-related task, verify that:

- only approved paths were changed
- protected project records were preserved
- no hidden spillover edits occurred in adjacent lanes
- stated safety claims are backed by direct evidence

### 7) Destructive actions require explicit human approval

Any action that deletes, overwrites, or relocates project-owned records requires
explicit human approval in the current task context.

Absence of objection is not approval.

## Reference Boundaries

The following paths are commonly treated as high-sensitivity boundaries and
should not be modified in reference-only safety work unless explicitly approved:

- core/install.md
- core/update.md
- core/distribution/sources/\*\*
- core/distribution/generated/\*\*
- librarian/scripts/distribution/\*\*

## Review Checklist

Use this checklist for reference-only install/update safety changes:

- [ ] Scope is limited to approved reference paths
- [ ] No project-owned work records were overwritten
- [ ] No distribution implementation paths were changed
- [ ] No operational script instructions were added
- [ ] Safety claims separate direct evidence from inference

## Related References

- .hawp/kit/references/docs-alignment.md
- .hawp/kit/references/backlog-alignment.md
- .hawp/kit/references/public-safety-checklist.md
- .hawp/kit/references/publication-safety-guidelines.md
