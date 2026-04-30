# Bug: update.md does not migrate old `.hawp/usage` layouts into `.hawp/work`

- ID: BUG-001
- Date: 2026-04-30
- Status: Approved (user-provided intake spec)
- Risk: Medium (touches install/update shell flows used by downstream repos)

## Input

User report: after the kit/work restructure, running update on a project with an older HAWP layout (e.g. `.hawp/usage/BACKLOG.md`, `.hawp/usage/status/*`, `.hawp/usage/*_ADR.md`, root-level `.hawp/templates|patterns|reviews|examples|types`, or non-dot `hawp/`) refreshes `.hawp/kit/` and `.hawp/LICENSE` but does not reliably:

- create `.hawp/work/` when missing
- migrate `.hawp/usage/BACKLOG.md` → `.hawp/work/BACKLOG.md`
- migrate `.hawp/usage/status/*` → `.hawp/work/status/*`
- migrate `.hawp/usage/*_ADR.md` → `.hawp/work/adrs/*`
- preserve previously created `.hawp/work/**` files
- clean up old root-level kit folders (`.hawp/templates`, `.hawp/patterns`, `.hawp/reviews`, `.hawp/examples`, `.hawp/types`, `.hawp/usage`) after the new `.hawp/kit/` is in place
- remove legacy `human-ai-workflow-protocol-*` overlay files left over from earlier installs

## Context

- Source repo layout: `core/.hawp/{LICENSE, kit/, work/}` and `core/.github/{instructions, prompts, copilot-instructions.md}`.
- Install boundary: `.hawp/kit/**` is HAWP-managed (refreshable); `.hawp/work/**` is repo-owned (must never be overwritten).
- Current `core/update.md` already handles legacy `hawp/` and legacy `.hawp/status/`, but it does NOT handle `.hawp/usage/` migration, old root-level kit folders, missing `.hawp/work/` scaffold seeding, or stale overlay file cleanup.

## Analysis

**Root cause:** `core/update.md` was written for the kit/work transition assuming the only legacy shape was a flat `.hawp/` or a non-dot `hawp/`. It missed the realistic intermediate layout where downstream projects had `.hawp/usage/` (the pre-restructure name for what is now `work/`) plus root-level kit folders.

**Directly verified:**

- `core/update.md` migrates `hawp/ → .hawp/` and `.hawp/status/ → .hawp/work/status/` only.
- `core/install.md` already has `.hawp/work/` scaffold seeding logic; `update.md` does not.
- `core/.github/prompts/` and `core/.github/instructions/` no longer ship any `human-ai-workflow-protocol-*` files; downstream copies of those names are stale.

**Inferred:** A target repo updated with the current script ends up with the new `kit/` plus stale `usage/`, `templates/`, etc. siblings, and missing `work/` scaffold.

## Recommended Fix

Rewrite the migration block in both `core/update.md` and `core/install.md` to use guarded helpers and run cases A–F end-to-end. Same migration body in both files; only the surrounding prose differs.

### Migration order (one-shot shell)

1. Download source tarball.
2. If `hawp/` exists and is not a symlink: migrate `hawp/work`, `hawp/usage`, `hawp/status` into `.hawp/work` (no-clobber); preserve `hawp/LICENSE` if `.hawp/LICENSE` missing; `rm -rf hawp`.
3. If `.hawp/usage/` exists: migrate `BACKLOG.md`, `status/*`, `*_ADR.md` into `.hawp/work/{BACKLOG.md,status,adrs}` (no-clobber).
4. If `.hawp/status/` (root, pre-`work/`) exists: migrate to `.hawp/work/status/` (no-clobber).
5. Refresh `.hawp/LICENSE` and `.hawp/kit/**` from source (rm + cp).
6. Remove legacy root-level kit folders only after kit refresh succeeds: `.hawp/templates`, `.hawp/patterns`, `.hawp/reviews`, `.hawp/examples`, `.hawp/types`, and `.hawp/usage` (now empty).
7. Seed `.hawp/work/` scaffold files only when missing: `README.md`, `BACKLOG.md` (from `kit/templates/backlog.md`), `adrs/README.md`, `status/README.md`, `evidence/README.md`.
8. Refresh `.github/instructions/*.instructions.md` and `.github/prompts/*.prompt.md`.
9. Remove stale `.github/instructions/human-ai-workflow-protocol-*.instructions.md` and `.github/prompts/human-ai-workflow-protocol-*.prompt.md` overlay files.
10. Never overwrite `.github/copilot-instructions.md` if present.
11. Print summary.

### Helpers

```sh
copy_dir_no_clobber() { src="$1"; dest="$2"; [ -d "$src" ] && { mkdir -p "$dest"; cp -Rn "$src"/. "$dest"/ 2>/dev/null || true; }; }
copy_file_no_clobber() { src="$1"; dest="$2"; [ -f "$src" ] && [ ! -f "$dest" ] && { mkdir -p "$(dirname "$dest")"; cp "$src" "$dest"; }; }
```

### Files to change

- `core/update.md` — rewrite migration block, doc the new behavior.
- `core/install.md` — align migration block (same body) and confirm scaffold seeding remains.

### What to verify after

- [ ] Cases A–F handled per spec.
- [ ] `.hawp/work/**` never overwritten on re-run.
- [ ] Old `.hawp/{usage,templates,patterns,reviews,examples,types}` removed after refresh.
- [ ] `.github/copilot-instructions.md` preserved.
- [ ] Legacy `human-ai-workflow-protocol-*` overlay files removed.

## Status

- [x] Plan written
- [x] Approved (user-provided spec)
- [x] Implemented
- [x] Verified (Cases C + E migration; Case B re-run preserves `.hawp/work/BACKLOG.md`)
- [x] Closed
