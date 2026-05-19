```bash
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
REF="main"

echo "Source: ${OWNER}/${REPO}@${REF}"
echo "Archive: https://github.com/${OWNER}/${REPO}/archive/refs/heads/${REF}.tar.gz"

TMP_DIR=""
if [ -n "${HAWP_LOCAL_CORE:-}" ]; then
  SRC="${HAWP_LOCAL_CORE}"
  if [ ! -d "$SRC/.hawp/kit" ] || [ ! -d "$SRC/.github" ]; then
    echo "Error: HAWP_LOCAL_CORE must point to a core directory containing .hawp/kit and .github"
    exit 1
  fi
  echo "Source mode: local core (${SRC})"
else
  TMP_DIR="$(mktemp -d)"
  curl -fsSL "https://github.com/${OWNER}/${REPO}/archive/refs/heads/${REF}.tar.gz" \
    | tar -xz -C "$TMP_DIR"
  SRC="$TMP_DIR/${REPO}-${REF}/core"
  echo "Source mode: remote archive"
fi

if [ -d ".hawp" ]; then
  echo "Preflight: detected existing .hawp/."
  echo "Switching to update-compatible refresh mode for this run."
  echo "Tip: use update-${REF}.md next time when .hawp/ already exists."
fi

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
        echo "  reconciled (link): $src -> $dest"
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
        echo "  reconciled (id-fallback): $src -> $dest"
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
      echo "  retired (orphan): $src -> $dest"
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
cp -R "$SRC/.hawp/kit/instructions" .hawp/kit/
cp -R "$SRC/.hawp/kit/templates" .hawp/kit/
cp -R "$SRC/.hawp/kit/patterns"  .hawp/kit/
cp -R "$SRC/.hawp/kit/reviews"   .hawp/kit/
cp -R "$SRC/.hawp/kit/examples"  .hawp/kit/
cp -R "$SRC/.hawp/kit/types"      .hawp/kit/
cp -R "$SRC/.hawp/kit/usage"      .hawp/kit/
cp -R "$SRC/.hawp/kit/references" .hawp/kit/
cp -R "$SRC/.hawp/kit/standards"  .hawp/kit/

# --- 5. Cleanup: remove legacy root-level kit folders and stray docs (now under .hawp/kit/) ---
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

# --- 7. Install Copilot overlay into .github/ ---
mkdir -p .github/instructions .github/prompts
cp "$SRC/.github/instructions/"*.instructions.md .github/instructions/
cp "$SRC/.github/prompts/"*.prompt.md            .github/prompts/

# --- 8. Cleanup: remove stale legacy-named overlay files ---
find .github/instructions -maxdepth 1 -type f -name 'human-ai-workflow-protocol-*.instructions.md' -delete 2>/dev/null || true
find .github/prompts -maxdepth 1 -type f -name 'human-ai-workflow-protocol-*.prompt.md' -delete 2>/dev/null || true

# --- 9. Install copilot-instructions (HAWP-managed) ---
cp "$SRC/.github/copilot-instructions.md" .github/copilot-instructions.md

if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
  rm -rf "$TMP_DIR"
fi

echo "HAWP install complete."
echo "Refreshed: .hawp/LICENSE, .hawp/kit/**, .github/instructions/*, .github/prompts/*, .github/copilot-instructions.md"
echo "Preserved: .hawp/work/** (no-overwrite)"
echo "Reconciled: Done rows + Active-Work 'done'/'wont-fix' rows moved from .hawp/work/active/ when eligible (see 'reconciled (link):' and 'reconciled (id-fallback):' lines above)"
```
