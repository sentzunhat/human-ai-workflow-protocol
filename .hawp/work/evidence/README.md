# Verification Evidence

Optional verification artifacts referenced from plan files in `../active/` before close and `../closed/YYYY/MM/DD/` after close.

Discipline:

- Only real, captured evidence belongs here (logs, command output snapshots, screenshots).
- Never fabricate evidence to make a plan look complete.
- If a plan has no captured artifact, leave the reference out.
- Keep new evidence grouped by date and owning work UUID:
  `evidence/YYYY/MM/DD/{uuid}/`.
- Prefer filenames that start with the work item ID, for example
  `<work-id>-verification.md` or `<work-id>-smoke.txt`.
- Per-item evidence subfolders are the preferred layout for new work. Flat
  date-only evidence remains legacy-compatible until a dedicated migration.
- Historical evidence filenames may use older IDs or descriptive names; preserve
  them unless a dedicated cleanup task says otherwise.
