# Update HAWP From GitHub Main

Use this when you already installed HAWP in a target repository and want the latest package content from the GitHub `main` branch.

## What This Updates

This refreshes only HAWP-managed files:

- `.hawp/LICENSE`
- `.hawp/kit/**`
- `.github/instructions/*.instructions.md`
- `.github/prompts/*.prompt.md`

**Update boundary.** The source repository's historical `core/.hawp/work/` content (ADR files, status plan files, evidence artifacts) is intentionally excluded. Only `.hawp/LICENSE`, `.hawp/kit/`, and the `.github/` overlay are refreshed. Every existing file under `.hawp/work/**` — including BACKLOG, ADRs, status plans, evidence artifacts, and the scaffold seeded by install — is left untouched.

**Migration support.** This flow safely upgrades repos coming from older HAWP layouts:

- A legacy `hawp/` directory (no dot prefix, real directory only — symlinks are skipped) is migrated into `.hawp/`. Its `work/`, `usage/`, and `status/` contents are copied into `.hawp/work/...` with `cp -Rn` so existing `.hawp/work/` files always win.
- A legacy `.hawp/usage/` layout is migrated: `BACKLOG.md` → `.hawp/work/BACKLOG.md`, `status/*` → `.hawp/work/status/*`, `*_ADR.md` → `.hawp/work/adrs/*`.
- A legacy `.hawp/status/` folder (pre-`work/` layout) is moved into `.hawp/work/status/`.
- After `.hawp/kit/**` is refreshed, legacy root-level kit folders are removed: `.hawp/templates`, `.hawp/patterns`, `.hawp/reviews`, `.hawp/examples`, `.hawp/types`, `.hawp/usage`. Stale top-level docs that now live under `kit/` are also removed: `.hawp/README.md`, `.hawp/SPEC.md`, `.hawp/START_HERE.md`, `.hawp/AUTHORING_PATTERNS.md`.
- Any `.gitkeep` files under `.hawp/` are removed (the kit no longer ships placeholder files).
- Stale legacy overlay files named `.github/instructions/human-ai-workflow-protocol-*.instructions.md` and `.github/prompts/human-ai-workflow-protocol-*.prompt.md` are removed; current overlays use the `hawp-*` and `intake.*` naming.
- `.github/copilot-instructions.md` is never overwritten.
- Missing `.hawp/work/` scaffold files are seeded only when absent (`README.md`, `BACKLOG.md`, `adrs/README.md`, `status/README.md`, `evidence/README.md`).

## Prompt You Can Paste Into GitHub Copilot Chat

```text
Update this repository's HAWP kit from the latest github main branch of human-ai-workflow-protocol.

Requirements:
- migrate legacy hawp/ to .hawp/ if present, preserving hawp/work/ contents
- migrate legacy .hawp/usage/ (BACKLOG.md, status/, *_ADR.md) into .hawp/work/
- migrate legacy .hawp/status/ into .hawp/work/status/
- refresh .hawp/LICENSE and .hawp/kit/** (leave .hawp/work/ untouched)
- remove old root-level kit folders (.hawp/templates, .hawp/patterns, .hawp/reviews, .hawp/examples, .hawp/types, .hawp/usage) and stale top-level kit docs (.hawp/README.md, .hawp/SPEC.md, .hawp/START_HERE.md, .hawp/AUTHORING_PATTERNS.md) after kit refresh
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

# --- 1. Migration: legacy hawp/ -> .hawp/ (preserves hawp/work/, hawp/usage/, hawp/status/) ---
# Only acts on a real directory; symlinks named hawp are left untouched.
if [ -d "hawp" ] && [ ! -L "hawp" ]; then
  mkdir -p .hawp/work
  copy_dir_no_clobber "hawp/work" ".hawp/work"
  if [ -d "hawp/usage" ]; then
    copy_file_no_clobber "hawp/usage/BACKLOG.md" ".hawp/work/BACKLOG.md"
    copy_dir_no_clobber  "hawp/usage/status"     ".hawp/work/status"
    for f in hawp/usage/*_ADR.md; do
      [ -f "$f" ] && copy_file_no_clobber "$f" ".hawp/work/adrs/$(basename "$f")"
    done
  fi
  copy_dir_no_clobber "hawp/status" ".hawp/work/status"
  copy_file_no_clobber "hawp/LICENSE" ".hawp/LICENSE"
  rm -rf hawp
fi

# --- 2. Migration: legacy .hawp/usage/ -> .hawp/work/ ---
if [ -d ".hawp/usage" ]; then
  copy_file_no_clobber ".hawp/usage/BACKLOG.md" ".hawp/work/BACKLOG.md"
  copy_dir_no_clobber  ".hawp/usage/status"     ".hawp/work/status"
  for f in .hawp/usage/*_ADR.md; do
    [ -f "$f" ] && copy_file_no_clobber "$f" ".hawp/work/adrs/$(basename "$f")"
  done
fi

# --- 3. Migration: pre-work/ layout (.hawp/status/) -> .hawp/work/status/ ---
if [ -d ".hawp/status" ]; then
  copy_dir_no_clobber ".hawp/status" ".hawp/work/status"
fi

# --- 4. Refresh .hawp/LICENSE and .hawp/kit/** (preserves .hawp/work/) ---
rm -rf .hawp/kit
mkdir -p .hawp/kit
cp "$SRC/.hawp/LICENSE" .hawp/
cp "$SRC/.hawp/kit/README.md"             .hawp/kit/
cp "$SRC/.hawp/kit/START_HERE.md"         .hawp/kit/
cp "$SRC/.hawp/kit/SPEC.md"               .hawp/kit/
cp "$SRC/.hawp/kit/AUTHORING_PATTERNS.md" .hawp/kit/
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
rm -f .hawp/README.md .hawp/SPEC.md .hawp/START_HERE.md .hawp/AUTHORING_PATTERNS.md

# --- 5b. Cleanup: remove any .gitkeep files under .hawp/ (kit no longer ships them) ---
find .hawp -name .gitkeep -type f -delete 2>/dev/null || true

# --- 6. Seed .hawp/work/ scaffold (only when missing; never overwrites) ---
mkdir -p .hawp/work/adrs .hawp/work/status .hawp/work/evidence
copy_file_no_clobber "$SRC/.hawp/work/README.md"           ".hawp/work/README.md"
copy_file_no_clobber "$SRC/.hawp/kit/templates/backlog.md" ".hawp/work/BACKLOG.md"
copy_file_no_clobber "$SRC/.hawp/work/adrs/README.md"      ".hawp/work/adrs/README.md"
copy_file_no_clobber "$SRC/.hawp/work/status/README.md"    ".hawp/work/status/README.md"
copy_file_no_clobber "$SRC/.hawp/work/evidence/README.md"  ".hawp/work/evidence/README.md"

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
echo "Preserved: .hawp/work/** (untouched), .github/copilot-instructions.md (if present)"
```

## Post-Update Checklist

1. Confirm `.github/copilot-instructions.md` references `.hawp/kit/START_HERE.md` and `.hawp/kit/templates/status-report.md`.
2. Confirm `.hawp/LICENSE` exists and contains the Apache 2.0 text.
3. Confirm expected prompt files exist under `.github/prompts/` — including `intake.prompt.md`.
4. Confirm expected instruction files exist under `.github/instructions/` — including `intake.instructions.md`.
5. Confirm no `human-ai-workflow-protocol-*` named files remain under `.github/instructions/` or `.github/prompts/`.
6. If a legacy `hawp/`, `.hawp/usage/`, or `.hawp/status/` directory existed, confirm it has been migrated and removed, and that `.hawp/work/` retains every file you had previously.
7. Confirm legacy root-level kit folders (`.hawp/templates`, `.hawp/patterns`, `.hawp/reviews`, `.hawp/examples`, `.hawp/types`, `.hawp/usage`) and stale top-level kit docs (`.hawp/README.md`, `.hawp/SPEC.md`, `.hawp/START_HERE.md`, `.hawp/AUTHORING_PATTERNS.md`) are gone — they live under `.hawp/kit/` now. Confirm no `.gitkeep` files remain under `.hawp/`.
8. Review git diff before committing — pay special attention to anything inside `.hawp/work/` (should be a clean migration with no deletions).
9. Run your repo checks (lint/test/typecheck) if your workflow requires it.

## Notes

- This update flow is conservative and keeps local Copilot instruction customizations.
- `.hawp/` is laid out flat at the repository root with `.hawp/LICENSE`, `.hawp/kit/`, and `.hawp/work/`.
- `.hawp/work/` is repo-owned operating state and is never overwritten. Scaffold files are seeded only when missing.
- Legacy migrations use `cp -Rn` everywhere so any existing `.hawp/work/` file always wins.
- Legacy root-level kit folders are removed only after `.hawp/kit/**` has been rewritten from source.
- Symlinks named `hawp` are skipped on purpose to keep updates deterministic.
- The Apache 2.0 kit license travels with `.hawp/LICENSE` and is refreshed by this update flow.
