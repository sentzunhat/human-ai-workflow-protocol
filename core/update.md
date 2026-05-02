# Update HAWP From GitHub Main

Use this when you already installed HAWP in a target repository and want the latest package content from the GitHub `main` branch.

## What This Updates

This refreshes only HAWP-managed files:

- `.hawp/LICENSE`
- `.hawp/kit/**`
- `.github/instructions/*.instructions.md`
- `.github/prompts/*.prompt.md`

**Update boundary.** The source repository's own operating state lives at repo-root `.hawp/work/` (real BACKLOG, ADRs, status plans, evidence) and is never copied to downstream repositories by this flow. Only `.hawp/LICENSE`, `.hawp/kit/`, the `.github/` overlay, and missing `.hawp/work/` scaffold files are touched in the target repository. Existing `.hawp/work/**` files are never overwritten; as housekeeping, files in `.hawp/work/active/` that are already listed as `Done` in `.hawp/work/BACKLOG.md` can be moved to `closed/...` using Done-row metadata (plan link when usable, otherwise ID + Closed date fallback).

**Migration support.** This flow safely upgrades repos coming from older HAWP layouts:

- A legacy `hawp/` directory (no dot prefix, real directory only — symlinks are skipped) is migrated into `.hawp/`. Its `work/`, `usage/`, and `status/` contents are copied into `.hawp/work/...` with `cp -Rn` so existing `.hawp/work/` files always win.
- A legacy `.hawp/usage/` layout is migrated: `BACKLOG.md` → `.hawp/work/BACKLOG.md`, `status/*` → `.hawp/work/active/`, `*_ADR.md` → `.hawp/work/decisions/YYYY/MM/DD/`.
- A legacy `.hawp/status/` folder (pre-`work/` layout) is copied into `.hawp/work/notes/YYYY/MM/DD/`; `STATUS.md` is promoted to `.hawp/work/STATUS.md`.
- A current `.hawp/work/adrs/` folder is migrated to `.hawp/work/decisions/YYYY/MM/DD/`, then the legacy `.hawp/work/adrs/` folder is removed.
- A current `.hawp/work/status/` folder is migrated to `.hawp/work/active/`, then the legacy `.hawp/work/status/` folder is removed.
- A closed-plan reconciliation pass can move stale files from `.hawp/work/active/` to `.hawp/work/closed/...` using `.hawp/work/BACKLOG.md` Done rows (closed links when present, otherwise ID + Closed date fallback). The pass also covers `## Active Work` rows whose status is `done` or `wont-fix`, so items are reconciled even when the backlog has not yet been reformatted into the Done section. When no closed date is available from the backlog, the date is derived from the plan filename's `YYYY-MM-DD` prefix; as a last resort today's date is used.
- Any existing `.hawp/work/parked/` content is preserved as-is.
- After `.hawp/kit/**` is refreshed, legacy root-level kit folders are removed: `.hawp/templates`, `.hawp/patterns`, `.hawp/reviews`, `.hawp/examples`, `.hawp/types`, `.hawp/usage`. Stale top-level docs that now live under `kit/` are also removed: `.hawp/README.md`, `.hawp/spec.md`, `.hawp/start-here.md`, `.hawp/authoring-patterns.md`.
- Any `.gitkeep` files under `.hawp/` are removed (the kit no longer ships placeholder files).
- Stale legacy overlay files named `.github/instructions/human-ai-workflow-protocol-*.instructions.md` and `.github/prompts/human-ai-workflow-protocol-*.prompt.md` are removed; current overlays use the `hawp-*` and `intake.*` naming.
- `.github/copilot-instructions.md` is never overwritten.
- Missing `.hawp/work/` scaffold files are seeded only when absent (`README.md`, `STATUS.md`, `BACKLOG.md`, `active/README.md`, `parked/README.md`, `closed/README.md`, `decisions/README.md`, `evidence/README.md`, `status/README.md`, `notes/README.md`).
- Files in `.hawp/work/active/` with no matching entry in `.hawp/work/BACKLOG.md` (no ID pattern match and no plan-link match) are automatically retired to `closed/YYYY/MM/DD/` using the filename date prefix. This pass only runs when the backlog has at least one data row, so fresh installs are unaffected.

## Prompt You Can Paste Into GitHub Copilot Chat

```text
Update this repository's HAWP kit from the latest github main branch of human-ai-workflow-protocol.

Requirements:
- migrate legacy hawp/ to .hawp/ if present, preserving hawp/work/ contents
- migrate legacy .hawp/usage/ (BACKLOG.md → work/BACKLOG.md, status/ → work/active/, *_ADR.md → work/decisions/YYYY/MM/DD/) into .hawp/work/
- migrate legacy .hawp/status/ into .hawp/work/notes/YYYY/MM/DD/; promote STATUS.md to .hawp/work/STATUS.md
- migrate .hawp/work/adrs/ into .hawp/work/decisions/YYYY/MM/DD/ and remove legacy .hawp/work/adrs/
- migrate .hawp/work/status/ into .hawp/work/active/ and remove legacy .hawp/work/status/
- reconcile closed plan files by using `.hawp/work/BACKLOG.md` Done rows to move matching `.hawp/work/active/*.md` files into `.hawp/work/closed/...` (prefer closed plan links; fall back to ID + Closed date); also reconcile `## Active Work` rows whose status is `done` or `wont-fix`; derive closed date from filename prefix when not present in backlog; fall back to today as last resort
- retire orphan active items — files in `.hawp/work/active/` with no matching BACKLOG ID or plan-link — to `closed/YYYY/MM/DD/` using filename date prefix; only when the backlog has at least one data row
- preserve any existing .hawp/work/parked/ content as-is
- refresh .hawp/LICENSE and .hawp/kit/** (never overwrite .hawp/work/; allow backlog-based active→closed reconciliation via Done rows)
- remove old root-level kit folders (.hawp/templates, .hawp/patterns, .hawp/reviews, .hawp/examples, .hawp/types, .hawp/usage) and stale top-level kit docs (.hawp/README.md, .hawp/spec.md, .hawp/start-here.md, .hawp/authoring-patterns.md) after kit refresh
- remove any .gitkeep files under .hawp/
- seed .hawp/work/ scaffold files only when missing
- refresh .github/instructions/*.instructions.md and .github/prompts/*.prompt.md
- remove stale .github/instructions/human-ai-workflow-protocol-*.instructions.md and .github/prompts/human-ai-workflow-protocol-*.prompt.md
- do not overwrite .github/copilot-instructions.md if it already exists
- print a final summary of changed files and any manual merge points

Source:
- owner: sentzunhat
- repo: human-ai-workflow-protocol
- ref: main
```

## One-Shot Update Command (Copy/Paste)

Run this from the target repository root. It is **safe to re-run** and **safe on a repo that has `hawp/` or `.hawp/` from any earlier layout**. All `.hawp/work/` content is preserved.

```bash
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
REF="main"

TMP_DIR="$(mktemp -d)"
curl -fsSL "https://github.com/${OWNER}/${REPO}/archive/refs/heads/${REF}.tar.gz" \
  | tar -xz -C "$TMP_DIR"

SRC="$TMP_DIR/${REPO}-${REF}/core"
MDATE="$(date +%Y/%m/%d)"

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

  # Scans both ## Done rows (col5=Closed date) and ## Active Work rows where
  # status signals completion (done/wont-fix), so reconciliation works even
  # when items have not yet been reformatted into the Done section.
  awk -F'|' '
    function trim(s) { gsub(/^[[:space:]]+|[[:space:]]+$/, "", s); return s }
    /^## Done/        { section="done";   next }
    /^## Active/      { section="active"; next }
    /^## /            { section="";       next }
    section != "" && /^\|/ {
      id   = trim($2)
      col5 = trim($5)
      col6 = trim($6)
      if (id == "" || id ~ /^-+$/ || id == "ID") next
      if (section == "done") {
        closed = col5
        plan   = col6
      } else {
        # Active Work: col5=Status — only reconcile when status signals done
        if (col5 !~ /^(done|wont-fix)$/ && col5 !~ /^[0-9]{4}-[0-9]{2}-[0-9]{2}$/) next
        closed = (col5 ~ /^[0-9]{4}-[0-9]{2}-[0-9]{2}$/) ? col5 : ""
        plan   = col6
      }
      link = ""
      if (match(plan, /\(([^)]+)\)/)) {
        link = substr(plan, RSTART + 1, RLENGTH - 2)
      } else {
        link = plan
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
        echo "  reconciled (link): $(basename "$src")"
      fi
      continue
    fi

    # Derive closed_dir: prefer explicit date from backlog, then from filename
    # prefix, then fall back to today so nothing is silently skipped.
    # Uses case-glob matching to stay portable (avoids grep -E / rg aliases).
    case "$closed" in
      [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]) closed_dir="${closed//-//}" ;;
      *) closed_dir="" ;;
    esac

    find .hawp/work/active -maxdepth 1 -type f -name "*-${id}-*.md" | while IFS= read -r src; do
      [ -n "$src" ] || continue
      local_dir="$closed_dir"
      if [ -z "$local_dir" ]; then
        fname="$(basename "$src")"
        case "$fname" in
          [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-*)
            file_date="${fname:0:10}"
            local_dir="${file_date//-//}"
            ;;
          *)
            local_dir="$(date +%Y/%m/%d)"
            ;;
        esac
      fi
      dest=".hawp/work/closed/$local_dir/$(basename "$src")"
      if [ ! -e "$dest" ]; then
        mkdir -p "$(dirname "$dest")"
        mv "$src" "$dest"
        echo "  reconciled (id-fallback): $(basename "$src")"
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
# Safe because their reusable content has been rewritten under .hawp/kit/ and
# any repo-local items have already been migrated into .hawp/work/.
rm -rf .hawp/templates .hawp/patterns .hawp/reviews .hawp/examples .hawp/types .hawp/usage
rm -f .hawp/README.md .hawp/spec.md .hawp/start-here.md .hawp/authoring-patterns.md

# --- 5b. Cleanup: remove any .gitkeep files under .hawp/ (kit no longer ships them) ---
find .hawp -name .gitkeep -type f -delete 2>/dev/null || true

# --- 6. Seed .hawp/work/ scaffold (only when missing; never overwrites) ---
mkdir -p .hawp/work/active .hawp/work/parked .hawp/work/closed .hawp/work/decisions .hawp/work/evidence .hawp/work/status .hawp/work/notes
copy_file_no_clobber "$SRC/.hawp/work/README.md"               ".hawp/work/README.md"
copy_file_no_clobber "$SRC/.hawp/work/STATUS.md"               ".hawp/work/STATUS.md"
copy_file_no_clobber "$SRC/.hawp/work/BACKLOG.md"              ".hawp/work/BACKLOG.md"
copy_file_no_clobber "$SRC/.hawp/work/active/README.md"        ".hawp/work/active/README.md"
copy_file_no_clobber "$SRC/.hawp/work/parked/README.md"        ".hawp/work/parked/README.md"
copy_file_no_clobber "$SRC/.hawp/work/closed/README.md"        ".hawp/work/closed/README.md"
copy_file_no_clobber "$SRC/.hawp/work/decisions/README.md"     ".hawp/work/decisions/README.md"
copy_file_no_clobber "$SRC/.hawp/work/evidence/README.md"      ".hawp/work/evidence/README.md"
copy_file_no_clobber "$SRC/.hawp/work/status/README.md"        ".hawp/work/status/README.md"
copy_file_no_clobber "$SRC/.hawp/work/notes/README.md"         ".hawp/work/notes/README.md"

# --- 7. Refresh overlay files (HAWP-managed; safe to overwrite) ---
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

echo "HAWP update complete."
echo "Refreshed: .hawp/LICENSE, .hawp/kit/**, .github/instructions/*, .github/prompts/*"
echo "Preserved: .hawp/work/** (no-overwrite), .github/copilot-instructions.md (if present)"
echo "Reconciled: Done rows + Active-Work 'done'/'wont-fix' rows moved from .hawp/work/active/ when eligible (see 'reconciled:' lines above)"
```

## Post-Update Checklist

1. Confirm `.github/copilot-instructions.md` references `.hawp/kit/start-here.md` and `.hawp/kit/usage/status-report.md`.
2. Confirm `.hawp/LICENSE` exists and contains the Apache 2.0 text.
3. Confirm expected prompt files exist under `.github/prompts/` — including `intake.prompt.md`.
4. Confirm expected instruction files exist under `.github/instructions/` — including `intake.instructions.md`.
5. Confirm backlog alignment overlay files exist: `.github/instructions/hawp-backlog-alignment.instructions.md` and `.github/prompts/hawp-backlog-alignment.prompt.md`.
6. Confirm no `human-ai-workflow-protocol-*` named files remain under `.github/instructions/` or `.github/prompts/`.
7. Confirm `.hawp/work/parked/README.md` exists when the scaffold is seeded.
8. If a legacy `hawp/`, `.hawp/usage/`, `.hawp/status/`, `.hawp/work/adrs/`, or `.hawp/work/status/` directory existed, confirm it has been migrated and its content preserved in `.hawp/work/active/`, `.hawp/work/decisions/YYYY/MM/DD/`, or `.hawp/work/notes/YYYY/MM/DD/` as appropriate.
9. Confirm legacy `.hawp/work/adrs/` and `.hawp/work/status/` folders are no longer present after migration.
10. Confirm done items were moved out of `.hawp/work/active/` using `.hawp/work/BACKLOG.md` Done rows and Active Work rows with `done`/`wont-fix` status (closed links when present, otherwise ID + Closed date fallback, then filename date prefix, then today), when source files existed and destination files were missing.
11. Confirm any orphan active files (no matching BACKLOG ID or plan-link) were retired to `closed/...` when the backlog had at least one data row (see `retired (orphan):` lines above).
12. Confirm legacy root-level kit folders (`.hawp/templates`, `.hawp/patterns`, `.hawp/reviews`, `.hawp/examples`, `.hawp/types`, `.hawp/usage`) and stale top-level kit docs (`.hawp/README.md`, `.hawp/spec.md`, `.hawp/start-here.md`, `.hawp/authoring-patterns.md`) are gone — they live under `.hawp/kit/` now. Confirm no `.gitkeep` files remain under `.hawp/`.
13. Confirm backlog plan links resolve correctly after reconciliation (especially Done-row links to `closed/...` files and no stale `active/...` link targets for closed items).
14. Review git diff before committing — pay special attention to anything inside `.hawp/work/` (expected changes are migration copies, legacy-folder removals, and optional active→closed reconciliations only).
15. Run your repo checks (lint/test/typecheck) if your workflow requires it.

## Notes

- This update flow is conservative and keeps local Copilot instruction customizations.
- `.hawp/` is laid out flat at the repository root with `.hawp/LICENSE`, `.hawp/kit/`, and `.hawp/work/`.
- `.hawp/work/` is repo-owned operating state and is never overwritten. Scaffold files are seeded only when missing, including `parked/README.md`. Closed-plan reconciliation may move files from `active/` to `closed/...` when Done rows in the backlog indicate the item is closed, and also when Active Work rows carry a `done` or `wont-fix` status. The closed date is resolved in order: backlog Closed column → filename `YYYY-MM-DD` prefix → today. Reconciliation prints each file it moves so results are visible in the update log.
- Orphan active-item retirement moves `.hawp/work/active/` files with no matching BACKLOG entry to `closed/...`. The pass only runs when the backlog has at least one data row, so fresh installs and repos that do not yet use the backlog are unaffected. Each retired file is printed as `retired (orphan): <filename>`.
- Legacy migrations use `cp -Rn` everywhere so any existing `.hawp/work/` file always wins.
- Legacy `.hawp/work/adrs/` and `.hawp/work/status/` are removed only after a successful no-clobber migration copy.
- Legacy root-level kit folders are removed only after `.hawp/kit/**` has been rewritten from source.
- Symlinks named `hawp` are skipped on purpose to keep updates deterministic.
- The Apache 2.0 kit license travels with `.hawp/LICENSE` and is refreshed by this update flow.
