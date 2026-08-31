# Parked Work

Folder-per-item space for deferred or icebox items.

Parked items are not active and not closed — they are ideas, low-priority
tasks, or blocked items set aside for later.

```text
parked/
  <work-id>/
    plan.md
```

When a parked item becomes active, move it to:

```text
../active/<ID>/plan.md
```

When a parked item is cancelled or will never be done, move it to:

```text
../closed/YYYY/MM/DD/<ID>/plan.md
```

Historical flat parked files may remain as preserved legacy records, but the
preferred forward layout is folder-per-item.
