```bash
set -euo pipefail

OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
REF="main"
PROVIDER="github"

echo "Source: ${OWNER}/${REPO}@${REF}"
echo "Provider: ${PROVIDER}"
echo "Archive: https://github.com/${OWNER}/${REPO}/archive/refs/heads/${REF}.tar.gz"

TMP_DIR=""
if [ -n "${HAWP_LOCAL_CORE:-}" ]; then
  SRC="${HAWP_LOCAL_CORE}"
  if [ ! -d "$SRC/.hawp/kit" ]; then
    echo "Error: HAWP_LOCAL_CORE must point to a core directory containing .hawp/kit"
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
  echo "Tip: use distribution/generated/${PROVIDER}/update/${REF}.md next time when .hawp/ already exists."
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
MDATE="$(date +%Y/%m/%d)"

# --- 1. Migration: legacy hawp/ -> .hawp/ ---
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

# --- 3. Migration: .hawp/status/ -> .hawp/work/notes/ ---
if [ -d ".hawp/status" ]; then
  mkdir -p ".hawp/work/notes/$MDATE"
  copy_file_no_clobber ".hawp/status/STATUS.md" ".hawp/work/STATUS.md"
  copy_dir_no_clobber ".hawp/status" ".hawp/work/notes/$MDATE"
fi

# --- 3b. Migration: .hawp/work/adrs/ -> decisions/ ---
if [ -d ".hawp/work/adrs" ]; then
  mkdir -p ".hawp/work/decisions/$MDATE"
  if cp -Rn .hawp/work/adrs/. ".hawp/work/decisions/$MDATE"/ 2>/dev/null; then
    rm -rf .hawp/work/adrs
  fi
fi

reconcile_closed_plans_from_backlog

# --- 4. Refresh .hawp/LICENSE and .hawp/kit/** ---
rm -rf .hawp/kit
mkdir -p .hawp/kit
cp "$SRC/.hawp/LICENSE" .hawp/
cp "$SRC/.hawp/kit/README.md"             .hawp/kit/
cp "$SRC/.hawp/kit/start-here.md"         .hawp/kit/
cp -R "$SRC/.hawp/kit/instructions" .hawp/kit/
cp -R "$SRC/.hawp/kit/templates" .hawp/kit/
cp -R "$SRC/.hawp/kit/patterns"  .hawp/kit/
cp -R "$SRC/.hawp/kit/reviews"   .hawp/kit/
cp -R "$SRC/.hawp/kit/examples"  .hawp/kit/
cp -R "$SRC/.hawp/kit/types"      .hawp/kit/
cp -R "$SRC/.hawp/kit/usage"      .hawp/kit/
cp -R "$SRC/.hawp/kit/references" .hawp/kit/
cp -R "$SRC/.hawp/kit/standards"  .hawp/kit/

# --- 4b. Install hawp CLI binary (platform-detected from GitHub release) ---
install_hawp_binary() {
  local _os _arch _asset _ext _dest _url _checksum_url _expected _actual
  local _tag

  _os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  _arch="$(uname -m)"
  _ext=""

  case "$_os" in
    linux)  _os="linux"  ;;
    darwin) _os="darwin" ;;
    mingw*|msys*|cygwin*|windows_nt*) _os="windows"; _ext=".exe" ;;
    *)
      echo "hawp install: unsupported OS '$_os' — skipping binary install."
      return 0
      ;;
  esac

  case "$_arch" in
    x86_64|amd64)  _arch="amd64" ;;
    aarch64|arm64) _arch="arm64" ;;
    *)
      echo "hawp install: unsupported arch '$_arch' — skipping binary install."
      return 0
      ;;
  esac

  _asset="hawp-${_os}-${_arch}${_ext}"
  _dest=".hawp/bin/hawp-bin${_ext}"

  # Resolve latest release tag from GitHub API.
  _tag="$(curl -fsSL "https://api.github.com/repos/${OWNER}/${REPO}/releases/latest" 2>/dev/null \
    | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/' || true)"
  if [ -z "$_tag" ]; then
    _tag="$(curl -fsSL "https://api.github.com/repos/${OWNER}/${REPO}/releases?per_page=1" 2>/dev/null \
      | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/' || true)"
    [ -n "$_tag" ] && echo "hawp install: /releases/latest unavailable; using releases list fallback."
  fi

  if [ -z "$_tag" ]; then
    echo "hawp install: could not resolve latest release tag — skipping binary install."
    return 0
  fi

  _url="https://github.com/${OWNER}/${REPO}/releases/download/${_tag}/${_asset}"
  _checksum_url="https://github.com/${OWNER}/${REPO}/releases/download/${_tag}/checksums.txt"

  echo "hawp binary: ${_asset} (release ${_tag})"
  mkdir -p .hawp/bin
  curl -fsSL -o "${_dest}.tmp" "$_url" || {
    echo "hawp install: download failed — skipping binary install."
    return 0
  }

  # Verify SHA256 when checksums.txt is available
  if curl -fsSL -o /tmp/hawp-checksums.txt "$_checksum_url" 2>/dev/null; then
    _expected="$(grep " ${_asset}$" /tmp/hawp-checksums.txt | awk '{print $1}')"
    if [ -n "$_expected" ]; then
      if command -v sha256sum >/dev/null 2>&1; then
        _actual="$(sha256sum "${_dest}.tmp" | awk '{print $1}')"
      elif command -v shasum >/dev/null 2>&1; then
        _actual="$(shasum -a 256 "${_dest}.tmp" | awk '{print $1}')"
      else
        _actual=""
      fi
      if [ -n "$_actual" ] && [ "$_actual" != "$_expected" ]; then
        echo "hawp install: SHA256 mismatch — aborting binary install."
        rm -f "${_dest}.tmp"
        return 1
      fi
      [ -n "$_actual" ] && echo "hawp binary: SHA256 verified."
    fi
    rm -f /tmp/hawp-checksums.txt
  fi

  mv "${_dest}.tmp" "$_dest"
  chmod +x "$_dest"

  # Install the shell wrapper at .hawp/bin/hawp so it delegates to hawp-bin.
  if [ -f "$SRC/.hawp/bin/hawp" ]; then
    cp "$SRC/.hawp/bin/hawp" .hawp/bin/hawp
    chmod +x .hawp/bin/hawp
  fi
  if [ -f "$SRC/.hawp/bin/hawp-mcp" ]; then
    cp "$SRC/.hawp/bin/hawp-mcp" .hawp/bin/hawp-mcp
    chmod +x .hawp/bin/hawp-mcp
  fi

  echo "hawp binary: installed to ${_dest}"
}
install_hawp_binary

rm -rf .hawp/templates .hawp/patterns .hawp/reviews .hawp/examples .hawp/types .hawp/usage
rm -f .hawp/README.md .hawp/spec.md .hawp/start-here.md .hawp/authoring-patterns.md
find .hawp -name .gitkeep -type f -delete 2>/dev/null || true

# --- 5. Seed .hawp/work/ scaffold (only when missing) ---
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
```
