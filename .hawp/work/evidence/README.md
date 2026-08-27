# Verification Evidence

Optional verification artifacts referenced from plan files in `../active/` before close and `../closed/YYYY/MM/DD/` after close.

Discipline:

- Only real, captured evidence belongs here (logs, command output snapshots, screenshots).
- Never fabricate evidence to make a plan look complete.
- If a plan has no captured artifact, leave the reference out.
- Keep evidence grouped by date: `evidence/YYYY/MM/DD/`.
- Prefer filenames that start with the work item ID, for example
  `<work-id>-verification.md` or `<work-id>-smoke.txt`.
- Per-item evidence subfolders are not required. Flat-by-date evidence files are
  the preferred steady-state layout.
- Historical evidence filenames may use older IDs or descriptive names; preserve
  them unless a dedicated cleanup task says otherwise.
