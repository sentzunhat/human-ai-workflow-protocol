# Install HAWP GitHub Overlay

Use this package when the target repository already has a repo-root `.hawp/` folder and you want to add the matching GitHub Copilot integration under `.github/`.

## One-shot install + usage quickstart

If you want a single copy/paste installation flow, run this from the target repository root.

This script is **safe to re-run** and is **safe on a repo that already has `hawp/` or `.hawp/` from any earlier layout**. It will:

- Migrate a legacy `hawp/` directory to `.hawp/`, preserving any `hawp/work/`, `hawp/usage/`, and `hawp/status/` content.
- Migrate a legacy `.hawp/usage/` layout into `.hawp/work/` (BACKLOG → `work/BACKLOG.md`, status → `work/active/`, ADRs → `work/decisions/YYYY/MM/DD/`).
- Migrate a legacy `.hawp/status/` folder (pre-`work/` layout) into `.hawp/work/notes/YYYY/MM/DD/` and promote `STATUS.md` to `.hawp/work/STATUS.md`.
- Migrate an existing `.hawp/work/adrs/` into `.hawp/work/decisions/YYYY/MM/DD/`, then remove legacy `.hawp/work/adrs/`.
- Migrate an existing `.hawp/work/status/` into `.hawp/work/active/`, then remove legacy `.hawp/work/status/`.
- Reconcile closed plan files by using `.hawp/work/BACKLOG.md` Done rows to move matching `.hawp/work/active/*.md` files into `.hawp/work/closed/...` (prefer closed plan links; fall back to ID + Closed date).
- Retire orphan active items — files in `.hawp/work/active/` with no matching BACKLOG ID or plan-link — to `closed/YYYY/MM/DD/` using the filename date prefix. This pass only runs when the backlog has at least one data row.
- Refresh `.hawp/kit/` and `.hawp/LICENSE` from the package.
- Remove legacy root-level kit folders (`.hawp/templates`, `.hawp/patterns`, `.hawp/reviews`, `.hawp/examples`, `.hawp/types`, `.hawp/usage`) and stale top-level kit docs (`.hawp/README.md`, `.hawp/spec.md`, `.hawp/start-here.md`, `.hawp/authoring-patterns.md`) after the new `.hawp/kit/` is in place.
- Remove any `.gitkeep` files under `.hawp/` (the kit no longer ships placeholders).
- Never overwrite existing files under `.hawp/work/` (BACKLOG, active items, parked items, decisions, evidence, notes). During reconciliation, eligible done plan files may be moved from `active/` to `closed/...`.
- Seed missing `.hawp/work/` scaffold files only when they do not exist.
- Refresh `.github/instructions/*.instructions.md` and `.github/prompts/*.prompt.md`.
- Remove stale `.github/instructions/human-ai-workflow-protocol-*.instructions.md` and `.github/prompts/human-ai-workflow-protocol-*.prompt.md` overlay files.
- Seed `.github/copilot-instructions.md` only if it does not already exist.

```bash
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
REF="main"

TMP_DIR="$(mktemp -d)"
curl -fsSL "https://github.com/${OWNER}/${REPO}/archive/refs/heads/${REF}.tar.gz" \
  | tar -xz -C "$TMP_DIR"

SRC="$TMP_DIR/${REPO}-${REF}/core"

# --- Helpers (no-clobber copy; never overwrite repo-owned files) ---
copy_dir_no_clobber() {
  src="$1"; dest="$2"
  if [ -d "$src" ]; then
    mkdir -p "$dest"
    cp -Rn "$src"/. "$dest"/ 2>/dev/null || true
  fi
}
copy_file_no_clobber() {
  src="$1"; dest="$2"
  if [ -f "$src" ] && [ ! -f "$dest" ]; then
    mkdir -p "$(dirname "$dest")"
    cp "$src" "$dest"
  fi
}
reconcile_closed_plans_from_backlog() {
  backlog=".hawp/work/BACKLOG.md"
  [ -f "$backlog" ] || return 0

  awk -F'|' '
    function trim(s) { gsub(/^[[:space:]]+|[[:space:]]+$/, "", s); return s }
    /^## Done/ { in_done=1; next }
    /^## / && in_done { in_done=0 }
    in_done && /^\|/ {
      id = trim($2)
      closed = trim($5)
      plan = trim($6)
      if (id == "" || id ~ /^-+$/) next
      link = ""
      if (match(plan, /\(([^)]+)\)/)) {
        link = substr(plan, RSTART + 1, RLENGTH - 2)
      }
      print id "\t" closed "\t" link
    }
  ' "$backlog" | while IFS=$'\t' read -r id closed link_path; do
    [ -n "$id" ] || continue

    closed_path=""
    case "$link_path" in
      .hawp/work/closed/*) closed_path="${link_path#.hawp/work/}" ;;
      work/closed/*) closed_path="${link_path#work/}" ;;
      closed/*) closed_path="$link_path" ;;
    esac
    closed_path="${closed_path%%#*}"

    if [ -n "$closed_path" ]; then
      plan_name="$(basename "$closed_path")"
      src=".hawp/work/active/$plan_name"
      dest=".hawp/work/$closed_path"
      if [ -f "$src" ] && [ ! -e "$dest" ]; then
        mkdir -p "$(dirname "$dest")"
        mv "$src" "$dest"
      fi
      continue
    fi

    case "$closed" in
      ????-??-??) closed_dir="${closed//-//}" ;;
      *) continue ;;
    esac

    find .hawp/work/active -maxdepth 1 -type f -name "*-${id}-*.md" | while IFS= read -r src; do
      [ -n "$src" ] || continue
      dest=".hawp/work/closed/$closed_dir/$(basename "$src")"
      if [ ! -e "$dest" ]; then
        mkdir -p "$(dirname "$dest")"
        mv "$src" "$dest"
      fi
    done
  done
}
retire_orphan_active_items() {
  # Moves .hawp/work/active/ files that have no matching row in
  # .hawp/work/BACKLOG.md (no ID pattern match and no plan-link match) to
  # .hawp/work/closed/YYYY/MM/DD/. Skips when the backlog is absent or has
  # no data rows so fresh installs are never affected.
  backlog=".hawp/work/BACKLOG.md"
  [ -f "$backlog" ] || return 0

  known_ids="$(awk -F'|' '
    function trim(s) { gsub(/^[[:space:]]+|[[:space:]]+$/, "", s); return s }
    /^\|/ {
      id = trim($2)
      if (id == "" || id ~ /^-+$/ || id == "ID") next
      print id
    }
  ' "$backlog")"
  [ -n "$known_ids" ] || return 0

  known_links="$(awk -F'|' '
    function trim(s) { gsub(/^[[:space:]]+|[[:space:]]+$/, "", s); return s }
    /^\|/ {
      for (i = 2; i <= NF; i++) {
        cell = trim($i)
        if (match(cell, /\(([^)]+\.md[^)]*)\)/)) {
          link = substr(cell, RSTART + 1, RLENGTH - 2)
          n = split(link, parts, "/")
          gsub(/#.*$/, "", parts[n])
          if (parts[n] != "") print parts[n]
        }
      }
    }
  ' "$backlog" 2>/dev/null)"

  find .hawp/work/active -maxdepth 1 -type f -name '*.md' ! -name 'README.md' | while IFS= read -r src; do
    fname="$(basename "$src")"
    matched=0
    while IFS= read -r id; do
      [ -n "$id" ] || continue
      case "$fname" in
        *"-${id}-"*|*"-${id}.md") matched=1; break ;;
      esac
    done <<< "$known_ids"
    if [ "$matched" -eq 0 ] && [ -n "$known_links" ]; then
      while IFS= read -r link_name; do
        [ -n "$link_name" ] || continue
        [ "$fname" = "$link_name" ] && { matched=1; break; }
      done <<< "$known_links"
    fi
    [ "$matched" -eq 1 ] && continue
    case "$fname" in
      [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-*)
        file_date="${fname:0:10}"
        closed_dir="${file_date//-//}"
        ;;
      *)
        closed_dir="$(date +%Y/%m/%d)"
        ;;
    esac
    dest=".hawp/work/closed/$closed_dir/$fname"
    if [ ! -e "$dest" ]; then
      mkdir -p "$(dirname "$dest")"
      mv "$src" "$dest"
      echo "  retired (orphan): $fname"
    fi
  done
}

MDATE="$(date +%Y/%m/%d)"

# --- 1. Migration: legacy hawp/ -> .hawp/ (preserves hawp/work/, hawp/usage/, hawp/status/) ---
# Only acts on a real directory; symlinks named hawp are left untouched.
if [ -d "hawp" ] && [ ! -L "hawp" ]; then
  mkdir -p .hawp/work
  copy_dir_no_clobber "hawp/work" ".hawp/work"
  if [ -d "hawp/usage" ]; then
    copy_file_no_clobber "hawp/usage/BACKLOG.md" ".hawp/work/BACKLOG.md"
    copy_dir_no_clobber  "hawp/usage/status"     ".hawp/work/active"
    mkdir -p ".hawp/work/decisions/$MDATE"
    for f in hawp/usage/*_ADR.md; do
      [ -f "$f" ] && copy_file_no_clobber "$f" ".hawp/work/decisions/$MDATE/$(basename "$f")"
    done
  fi
  mkdir -p ".hawp/work/notes/$MDATE"
  copy_dir_no_clobber "hawp/status" ".hawp/work/notes/$MDATE"
  copy_file_no_clobber "hawp/LICENSE" ".hawp/LICENSE"
  rm -rf hawp
fi

# --- 2. Migration: legacy .hawp/usage/ -> .hawp/work/ ---
if [ -d ".hawp/usage" ]; then
  copy_file_no_clobber ".hawp/usage/BACKLOG.md" ".hawp/work/BACKLOG.md"
  copy_dir_no_clobber  ".hawp/usage/status"     ".hawp/work/active"
  mkdir -p ".hawp/work/decisions/$MDATE"
  for f in .hawp/usage/*_ADR.md; do
    [ -f "$f" ] && copy_file_no_clobber "$f" ".hawp/work/decisions/$MDATE/$(basename "$f")"
  done
fi

# --- 3. Migration: pre-work/ layout (.hawp/status/) -> .hawp/work/notes/YYYY/MM/DD/ ---
if [ -d ".hawp/status" ]; then
  mkdir -p ".hawp/work/notes/$MDATE"
  copy_file_no_clobber ".hawp/status/STATUS.md" ".hawp/work/STATUS.md"
  copy_dir_no_clobber ".hawp/status" ".hawp/work/notes/$MDATE"
fi

# --- 3b. Migration: current .hawp/work/adrs/ -> .hawp/work/decisions/YYYY/MM/DD/ ---
if [ -d ".hawp/work/adrs" ]; then
  mkdir -p ".hawp/work/decisions/$MDATE"
  if cp -Rn .hawp/work/adrs/. ".hawp/work/decisions/$MDATE"/ 2>/dev/null; then
    rm -rf .hawp/work/adrs
  fi
fi

# --- 3c. Migration: current .hawp/work/status/ -> .hawp/work/active/ ---
if [ -d ".hawp/work/status" ]; then
  mkdir -p ".hawp/work/active"
  if cp -Rn .hawp/work/status/. .hawp/work/active/ 2>/dev/null; then
    rm -rf .hawp/work/status
  fi
fi

# --- 3d. Reconcile closed plans from backlog (active -> closed/...) ---
reconcile_closed_plans_from_backlog

# --- 3e. Retire orphan active items (no matching BACKLOG entry) ---
retire_orphan_active_items

# --- 4. Refresh .hawp/LICENSE and .hawp/kit/** (preserves .hawp/work/) ---
rm -rf .hawp/kit
mkdir -p .hawp/kit
cp "$SRC/.hawp/LICENSE" .hawp/
cp "$SRC/.hawp/kit/README.md"             .hawp/kit/
cp "$SRC/.hawp/kit/start-here.md"         .hawp/kit/
cp "$SRC/.hawp/kit/spec.md"               .hawp/kit/
cp "$SRC/.hawp/kit/authoring-patterns.md" .hawp/kit/
cp -R "$SRC/.hawp/kit/templates" .hawp/kit/
cp -R "$SRC/.hawp/kit/patterns"  .hawp/kit/
cp -R "$SRC/.hawp/kit/reviews"   .hawp/kit/
cp -R "$SRC/.hawp/kit/examples"  .hawp/kit/
cp -R "$SRC/.hawp/kit/types"     .hawp/kit/
cp -R "$SRC/.hawp/kit/usage"     .hawp/kit/

# --- 5. Cleanup: remove legacy root-level kit folders and stray docs (now under .hawp/kit/) ---
rm -rf .hawp/templates .hawp/patterns .hawp/reviews .hawp/examples .hawp/types .hawp/usage
rm -f .hawp/README.md .hawp/spec.md .hawp/start-here.md .hawp/authoring-patterns.md

# --- 5b. Cleanup: remove any .gitkeep files under .hawp/ (kit no longer ships them) ---
find .hawp -name .gitkeep -type f -delete 2>/dev/null || true

# --- 6. Seed .hawp/work/ scaffold (only when missing; never overwrites) ---
mkdir -p .hawp/work/active .hawp/work/parked .hawp/work/closed .hawp/work/decisions .hawp/work/evidence .hawp/work/notes
copy_file_no_clobber "$SRC/.hawp/work/README.md"               ".hawp/work/README.md"
copy_file_no_clobber "$SRC/.hawp/work/STATUS.md"               ".hawp/work/STATUS.md"
copy_file_no_clobber "$SRC/.hawp/work/BACKLOG.md"              ".hawp/work/BACKLOG.md"
copy_file_no_clobber "$SRC/.hawp/work/active/README.md"        ".hawp/work/active/README.md"
copy_file_no_clobber "$SRC/.hawp/work/parked/README.md"        ".hawp/work/parked/README.md"
copy_file_no_clobber "$SRC/.hawp/work/closed/README.md"        ".hawp/work/closed/README.md"
copy_file_no_clobber "$SRC/.hawp/work/decisions/README.md"     ".hawp/work/decisions/README.md"
copy_file_no_clobber "$SRC/.hawp/work/evidence/README.md"      ".hawp/work/evidence/README.md"
copy_file_no_clobber "$SRC/.hawp/work/notes/README.md"         ".hawp/work/notes/README.md"

# --- 7. Install Copilot overlay into .github/ ---
mkdir -p .github/instructions .github/prompts
cp "$SRC/.github/instructions/"*.instructions.md .github/instructions/
cp "$SRC/.github/prompts/"*.prompt.md            .github/prompts/

# --- 8. Cleanup: remove stale legacy-named overlay files ---
rm -f .github/instructions/human-ai-workflow-protocol-*.instructions.md
rm -f .github/prompts/human-ai-workflow-protocol-*.prompt.md

# --- 9. Seed copilot-instructions only when missing (never overwrite) ---
if [ ! -f .github/copilot-instructions.md ]; then
  cp "$SRC/.github/copilot-instructions.md" .github/copilot-instructions.md
fi

rm -rf "$TMP_DIR"
```

Quick usage after install:

1. In `.github/copilot-instructions.md`, ensure HAWP references point to `.hawp/kit/start-here.md` and `.hawp/kit/templates/status-report.md`.
2. Start task shaping from `.hawp/kit/start-here.md`.
3. Use `.hawp/LICENSE` as the installed Apache 2.0 license text for the HAWP kit content.
4. Use `.hawp/kit/templates/status-report.md` for context-transfer artifacts.
5. A starter `.hawp/work/` area is scaffolded automatically: `STATUS.md`, `BACKLOG.md`, `README.md`, and `README.md` files in `active/`, `parked/`, `closed/`, `decisions/`, `evidence/`, and `notes/` (all seeded from the source repo's `.hawp/work/` scaffold). These are seeded once and owned by your repo from that point on.

## Scope Clarification

This document covers the full HAWP installation: `.hawp/` at the repo root, including `.hawp/LICENSE`, plus the GitHub Copilot overlay under `.github/`.

Install boundary: the source repository's own operating state lives at root `.work/` (active work, parked work, closed work, decisions, evidence, real BACKLOG) and is never installed downstream. The install seeds a clean `work/` scaffold — `README.md`, `STATUS.md`, a starter `BACKLOG.md`, and `README.md` files in `active/`, `parked/`, `closed/`, `decisions/`, `evidence/`, and `notes/` — sourced from `core/.hawp/work/` in the kit repo. Downstream repos never receive the HAWP source repo's own backlog items or decision records.

The `benchmark/` folder is optional reference material and is not installed by this script.

## What This Package Installs

Install these files into the target repository's `.github/` folder:

- `copilot-instructions.md`
- `instructions/hawp-intake.instructions.md`
- `instructions/hawp-docs-alignment.instructions.md`
- `instructions/intake.instructions.md`
- `prompts/hawp-status-report.prompt.md`
- `prompts/hawp-intent-first-handoff.prompt.md`
- `prompts/hawp-docs-alignment-deterministic.prompt.md`
- `prompts/hawp-docs-alignment-simplicity.prompt.md`
- `prompts/hawp-conservative-docs-drift-cleanup.prompt.md`
- `prompts/intake.prompt.md`

These files assume the target repository resolves HAWP content from repo-root paths such as `.hawp/kit/start-here.md`.

The installed `.hawp/` content also includes:

- `.hawp/LICENSE` — Apache 2.0 text
- `.hawp/work/README.md` — work area overview
- `.hawp/work/STATUS.md` — current state dashboard
- `.hawp/work/BACKLOG.md` — starter backlog (seeded from the source repo's `.hawp/work/BACKLOG.md` scaffold)
- `.hawp/work/active/README.md` — active work folder description
- `.hawp/work/parked/README.md` — parked/deferred work folder description
- `.hawp/work/closed/README.md` — closed work archive description
- `.hawp/work/decisions/README.md` — decisions/ADR folder description
- `.hawp/work/evidence/README.md` — evidence folder description
- `.hawp/work/notes/README.md` — notes folder description

`work/` files are seeded once and never overwritten on subsequent updates. When Done rows mark plans as closed, reconciliation may move matching files from `active/` to `closed/...`.

## Preconditions

This package is **self-bootstrapping**. The one-shot script above creates `.hawp/` from scratch if it does not exist. There is no requirement that `.hawp/` already be present.

If you are running the manual install procedure (below) instead of the one-shot script, you should already have a repo-root `.hawp/kit/` populated with at least:

- `.hawp/LICENSE`
- `.hawp/kit/README.md`
- `.hawp/kit/spec.md`
- `.hawp/kit/authoring-patterns.md`
- `.hawp/kit/start-here.md`
- `.hawp/kit/templates/status-report.md`

## Upgrading An Existing Install

If the target repo already has `hawp/` (legacy, no dot prefix) or `.hawp/` from any earlier layout, the one-shot script handles the upgrade safely. Specifically:

- **Legacy `hawp/` → `.hawp/` migration.** A real `hawp/` directory (not a symlink) is moved to `.hawp/`. Any `hawp/work/`, `hawp/usage/`, and `hawp/status/` content is copied into `.hawp/work/...` first using `cp -Rn` so existing `.hawp/work/` files are never clobbered. The legacy `hawp/` directory is then removed.
- **Legacy `.hawp/usage/` migration.** `BACKLOG.md` → `.hawp/work/BACKLOG.md`, `status/*` → `.hawp/work/active/`, `*_ADR.md` → `.hawp/work/decisions/YYYY/MM/DD/`. Existing files in `.hawp/work/` always win.
- **Pre-`work/` layout migration.** A legacy `.hawp/status/` folder at the kit root is copied to `.hawp/work/notes/YYYY/MM/DD/`; `STATUS.md` is promoted to `.hawp/work/STATUS.md`.
- **Current `.hawp/work/adrs/` migration.** ADR files are copied to `.hawp/work/decisions/YYYY/MM/DD/`, then the legacy `.hawp/work/adrs/` folder is removed.
- **Current `.hawp/work/status/` migration.** Plan files are copied to `.hawp/work/active/`, then the legacy `.hawp/work/status/` folder is removed.
- **Closed-plan reconciliation.** The script reads `.hawp/work/BACKLOG.md` Done rows and moves matching files from `.hawp/work/active/` into `closed/...` when safe (source exists and destination does not). It prefers explicit closed plan links and falls back to ID + Closed date.
- **Orphan active-item retirement.** After reconciliation, any `.hawp/work/active/` file (excluding `README.md`) with no matching BACKLOG ID or plan-link is retired to `closed/YYYY/MM/DD/` using the filename date prefix. This pass only runs when the backlog has at least one data row, so fresh installs are safe.
- **`.hawp/kit/` refresh.** The `.hawp/kit/` directory is removed and rewritten from the package allowlist. Nothing under `.hawp/work/` is touched by this step.
- **Legacy root-level kit cleanup.** After kit refresh, the script removes `.hawp/templates`, `.hawp/patterns`, `.hawp/reviews`, `.hawp/examples`, `.hawp/types`, and `.hawp/usage`, plus stale top-level kit docs `.hawp/README.md`, `.hawp/spec.md`, `.hawp/start-here.md`, and `.hawp/authoring-patterns.md`. Their reusable content now lives under `.hawp/kit/...` and any repo-local items have already been migrated into `.hawp/work/`. Any `.gitkeep` files under `.hawp/` are also removed.
- **`.hawp/work/` preservation.** Existing `.hawp/work/` files are never overwritten. Scaffold seeders only create missing files; reconciliation may move already-closed plans from `active/` to `closed/...`.
- **`.github/` overlay refresh.** Instruction and prompt files are HAWP-managed and are always replaced. Stale legacy-named overlay files matching `human-ai-workflow-protocol-*.instructions.md` and `human-ai-workflow-protocol-*.prompt.md` are removed. `.github/copilot-instructions.md` is seeded only when missing — if it already exists, merge the HAWP block manually.

If you only want to refresh kit content from upstream (no install steps), use [`update.md`](./update.md).

## Install Procedure

1. Check that the target repository has a repo-root `.hawp/` directory.
2. Check whether the target repository already has a `.github/` directory. Create it if missing.
3. Copy all `instructions/*.instructions.md` files into `.github/instructions/`, replacing existing versions. These instruction files are HAWP-managed and safe to overwrite.
4. Copy all `prompts/*.prompt.md` files into `.github/prompts/`, replacing existing versions. These prompt files are HAWP-managed and safe to overwrite during reinstall.
5. Check whether `.github/copilot-instructions.md` already exists.
6. If `.github/copilot-instructions.md` does not exist, copy this package's `copilot-instructions.md` into place.
7. If `.github/copilot-instructions.md` already exists, merge the HAWP block below into the existing file instead of overwriting unrelated repository guidance.

## Reinstall Behavior

Reinstalling is safe. The procedure is the same as a fresh install with the following notes:

- Steps 3 and 4 (intake instructions and status-report prompt): always replace the existing files. They are entirely HAWP-managed. No user content will be lost.
- Step 7 (copilot-instructions.md merge): if the file was previously merged, re-check that the HAWP block still matches the current package version. If the HAWP block content has changed, update it in place without disturbing unrelated repo instructions.
- After reinstalling, re-run the Validation Checklist to confirm the installed state is correct.

For a dedicated upstream refresh flow that pulls the latest content from GitHub `main`, use `update.md`.

## Cleanup Step

After the overlay files are installed, remove packaging leftovers that are not part of the target repository's steady-state HAWP layout.

Clean up these cases:

- If you copied this package as a temporary `core/` folder, remove that temporary folder after its files have been merged into the target layout.
- Remove duplicate overlay files that were copied outside the target `.github/` folder.
- Keep only the repo-root `.hawp/` folder and the target repo's `.github/` integration files. The target steady state should not require a sibling package source folder.
- Do not remove existing non-HAWP `.github/` files that belong to the target repository.

## HAWP Block For Existing Copilot Instructions

Add or merge this block into the target repository's `.github/copilot-instructions.md`:

```md
This repository uses HAWP as a lightweight workflow method.

Follow the repo-local HAWP guidance in:

- .hawp/kit/start-here.md
- .hawp/kit/templates/status-report.md

Use .hawp/kit/start-here.md as the operating guide for how this repo applies HAWP in practice.

Use .hawp/kit/templates/status-report.md when the user asks for a:

- status report
- checkpoint summary
- context transfer summary
- second-brain review artifact

Saved status reports belong in:

- .hawp/work/active/

For bugs/tasks, track in .hawp/work/BACKLOG.md. Active plan files go in .hawp/work/active/. Deferred items can live in .hawp/work/parked/. Close by moving to .hawp/work/closed/YYYY/MM/DD/.

Keep the repo-local HAWP layer lean.

HAWP in this repository is a practical workflow layer, not a runtime engine, compiler, validator, orchestrator, or memory system.

Do not:

- invent per-field folders such as input/, context/, mission/, constraints/, output/, or checkpoint/
- imply a runtime engine, compiler, validator, orchestrator, persistence layer, or memory system
- overstate certainty; keep direct evidence separate from inference

Prefer compact, decision-useful outputs.
```

## Validation Checklist

After installation, confirm all of the following are true:

- `.github/copilot-instructions.md` references `.hawp/kit/start-here.md`
- `.github/copilot-instructions.md` references `.hawp/kit/templates/status-report.md`
- `.github/instructions/hawp-intake.instructions.md` exists
- `.github/instructions/hawp-docs-alignment.instructions.md` exists
- `.github/instructions/intake.instructions.md` exists
- `.github/prompts/hawp-status-report.prompt.md` exists
- `.github/prompts/hawp-intent-first-handoff.prompt.md` exists
- `.github/prompts/hawp-docs-alignment-deterministic.prompt.md` exists
- `.github/prompts/hawp-docs-alignment-simplicity.prompt.md` exists
- `.github/prompts/hawp-conservative-docs-drift-cleanup.prompt.md` exists
- `.github/prompts/intake.prompt.md` exists
- `.hawp/LICENSE` exists and contains the Apache 2.0 text
- `.hawp/work/README.md` exists
- `.hawp/work/STATUS.md` exists
- `.hawp/work/BACKLOG.md` exists
- `.hawp/work/active/README.md` exists
- `.hawp/work/parked/README.md` exists
- `.hawp/work/closed/README.md` exists
- `.hawp/work/decisions/README.md` exists
- `.hawp/work/evidence/README.md` exists
- `.hawp/work/notes/README.md` exists
- legacy `.hawp/work/adrs/` and `.hawp/work/status/` folders are absent after migration
- done items are not left behind in `.hawp/work/active/` when `.hawp/work/BACKLOG.md` Done rows provide reconcilable data (closed link or ID + Closed date)
- no temporary `.github/` overlay folder remains in the target repository unless the user is intentionally keeping the kit source there

## Failure Cases

- Missing `.hawp/` folder when running the manual install procedure: use the one-shot script instead, or copy the HAWP content first.
- Missing `.github/` folder: create it before copying the overlay files.
- Existing `.github/copilot-instructions.md` with unrelated repo rules: merge carefully and do not replace project-specific instructions.
- Legacy `hawp/` exists as a symlink: the migration step intentionally skips symlinks. Resolve the symlink manually before running the script.
- Temporary package files left in the target repository: remove the extra package source folder and keep only the installed `.github/` files.

## Intent

This package is a modular overlay. It adds HAWP-aware Copilot behavior without replacing the target repository's baseline instruction set.

## Path Note

All installed `.github/` files use `.hawp/...` paths that assume `.hawp/` sits at the target repository root. No path adaptation is required when using this install script.
