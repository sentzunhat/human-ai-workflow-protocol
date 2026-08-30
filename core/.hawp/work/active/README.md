# Active Work

Folder-per-item workspace for all open bugs and tasks.

Each item gets one directory:

```text
active/
  <work-id>/
    plan.md
    files.md   # optional file-tracking sidecar
```

Active work stays here until it is closed. On close, move the plan to:

```text
../closed/YYYY/MM/DD/<work-id>/plan.md
```

Historical flat files may still exist from older migration phases, but new and
actively maintained UUID-native items should use the folder layout above.
