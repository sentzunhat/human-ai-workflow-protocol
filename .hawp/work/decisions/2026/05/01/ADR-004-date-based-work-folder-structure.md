# ADR: Date-Based Work Folder Structure and Safe Install/Update Migration

## Status

Accepted

## Context

HAWP currently uses a markdown-first structure where `kit/` holds reusable protocol assets and `work/` holds project-specific work records.

As projects grow, `work/` can collect many bugs, tasks, notes, decisions, and evidence files. If these files stay in one flat folder, it becomes harder for humans and agents to traverse safely.

HAWP also needs install/update behavior that preserves existing project work. Updating HAWP should refresh reusable protocol assets without overwriting or destroying repo-owned work records.

## Decision

HAWP will use this folder model:

```txt
.hawp/
  kit/
    instructions/
    templates/
    examples/
    patterns/

  work/
    STATUS.md
    BACKLOG.md

    active/
      BUG-001.md
      TASK-002.md

    parked/
      IDEA-001.md
      TASK-005.md

    closed/
      YYYY/
        MM/
          DD/
            BUG-001.md
            TASK-002.md

    decisions/
      YYYY/
        MM/
          DD/
            ADR-001-example.md

    evidence/
      YYYY/
        MM/
          DD/
            BUG-001-smoke-check.md

    notes/
      YYYY/
        MM/
          DD/
            session-notes.md

  index/
    hawp.sqlite
```

The core rule is:

```txt
kit/  = reusable HAWP rules, templates, examples
work/ = project-specific truth, progress, bugs, decisions, evidence
index/ = generated search index, safe to rebuild
```

Markdown remains the source of truth.

SQLite, if added later, is only a generated search index.

## Active Work Rule

Active work stays flat and easy to find:

```txt
.hawp/work/active/
```

Examples:

```txt
.hawp/work/active/BUG-003.md
.hawp/work/active/TASK-004.md
```

Active work should not be buried under date folders.

## Parked Work Rule

Parked work is deferred or icebox — not active, not closed:

```txt
.hawp/work/parked/
```

Examples:

```txt
.hawp/work/parked/IDEA-001.md
.hawp/work/parked/TASK-005.md
```

When a parked item becomes active, move it to `active/`.
When a parked item is cancelled, move it to `closed/YYYY/MM/DD/`.

## Closed Work Rule

Closed work moves into date-based archive folders:

```txt
.hawp/work/closed/YYYY/MM/DD/
```

Example:

```txt
.hawp/work/closed/2026/05/01/BUG-003.md
```

This keeps active work small while preserving completed history.

## Decisions, Evidence, and Notes Rule

Decisions, evidence, and notes use date folders:

```txt
.hawp/work/decisions/YYYY/MM/DD/
.hawp/work/evidence/YYYY/MM/DD/
.hawp/work/notes/YYYY/MM/DD/
```

This makes the project easier to browse over time and easier for a future CLI/librarian to index.

## Install Rule

On a fresh install, HAWP should create the default folder structure if it does not already exist.

It may create:

```txt
.hawp/kit/
.hawp/work/
.hawp/work/STATUS.md
.hawp/work/BACKLOG.md
.hawp/work/active/
.hawp/work/parked/
.hawp/work/closed/
.hawp/work/decisions/
.hawp/work/evidence/
.hawp/work/notes/
```

The installer must not overwrite existing project-owned files.

## Update Rule

On update, HAWP may refresh:

```txt
.hawp/kit/
```

But it must preserve:

```txt
.hawp/work/
```

The update command may add missing folders, but it must not overwrite existing work records.

Safe behavior:

```txt
Create missing folders.
Create missing seed files only if absent.
Never overwrite existing STATUS.md.
Never overwrite existing BACKLOG.md.
Never delete existing work files.
Never move files without creating a migration record.
```

## Migration Rule

For older HAWP projects, update/install should migrate known legacy folders safely.

Legacy examples:

```txt
.hawp/usage/
.hawp/status/
hawp/
```

Migration should copy or move old project-owned records into the new `work/` layout without losing data.

Suggested migration mapping:

```txt
.hawp/usage/BACKLOG.md  -> .hawp/work/BACKLOG.md
.hawp/status/STATUS.md  -> .hawp/work/STATUS.md
.hawp/usage/*           -> .hawp/work/notes/YYYY/MM/DD/
hawp/*                  -> .hawp/work/notes/YYYY/MM/DD/
```

If the target file already exists, the migration must not overwrite it.

Instead, create a preserved copy with a suffix:

```txt
BACKLOG.migrated-2026-05-01.md
STATUS.migrated-2026-05-01.md
```

## Migration Record

Every migration should create a migration note:

```txt
.hawp/work/notes/YYYY/MM/DD/migration-note.md
```

The note should include:

```txt
- Date of migration
- Source paths found
- Target paths created
- Files skipped because targets already existed
- Any files that require manual review
```

## Non-Goals

This ADR does not add a required CLI.

This ADR does not make SQLite the source of truth.

This ADR does not add vector search.

This ADR does not introduce a project management system.

This ADR only defines the folder structure and safe migration/update rules.

## Consequences

Benefits:

- Active work stays small and easy to load.
- Closed work remains preserved and searchable.
- Date folders make long-running projects easier to traverse.
- Future CLI and SQLite indexing have a clean structure to read.
- Old projects can update safely without losing records.

Tradeoffs:

- The folder structure becomes slightly more formal.
- Migration logic needs to be careful.
- Agents must follow the rule: write work records to `work/`, not `kit/`.

## Final Rule

```txt
Read from kit.
Write to work.
Archive by date.
Index later.
Never overwrite project truth.
```
