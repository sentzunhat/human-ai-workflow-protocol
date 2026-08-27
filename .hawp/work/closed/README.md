# Closed Work Archive

Date-based archive of completed and cancelled work items.

Layout:

```
closed/
  YYYY/
    MM/
      DD/
        <work-id>/
          plan.md
```

Current preferred shape is one folder per work item, with the final detail file at
`closed/YYYY/MM/DD/<work-id>/plan.md`.

Historical flat files such as `TASK-002.md` may remain in the archive. They are
legacy records, not a requirement for new closures.

To close an item: move its plan from `../active/<work-id>/plan.md` into the
appropriate date folder. Preserve historical archive records; do not bulk-rename
old entries unless a dedicated normalization task explicitly covers that work.
