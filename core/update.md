# Update HAWP From GitHub Main

Use this when you already installed HAWP in a target repository and want the latest package content from the GitHub `main` branch.

## What This Updates

This refreshes only HAWP-managed files:

- installable `.hawp/` public kit files and folders (`.hawp/LICENSE` and `.hawp/kit/**`)
- `.github/instructions/*.instructions.md`
- `.github/prompts/*.prompt.md`

**Update boundary.** The source repository's historical `core/.hawp/work/` content (ADR files, status plan files, evidence artifacts) is intentionally excluded. Only `.hawp/LICENSE` and `.hawp/kit/` are refreshed. Any local `.hawp/work/` in the target repo — including files you have created and the scaffold seeded by install — is left completely untouched.

**Migration support.** This flow safely handles repos coming from older HAWP layouts:

- A legacy `hawp/` directory (no dot prefix) is migrated to `.hawp/`. Any `hawp/work/` content is moved into `.hawp/work/` first using `cp -Rn`, so existing `.hawp/work/` files are never overwritten. Symlinks named `hawp` are left untouched.
- A legacy `.hawp/status/` folder at the kit root (pre-`work/` layout) is moved into `.hawp/work/status/`.

`.github/copilot-instructions.md` is never overwritten by this flow.

## Prompt You Can Paste Into GitHub Copilot Chat

```text
Update this repository's HAWP kit from the latest github main branch of human-ai-workflow-protocol.

Requirements:
- migrate legacy hawp/ to .hawp/ if present, preserving hawp/work/ contents
- migrate legacy .hawp/status/ into .hawp/work/status/ if present
- refresh .hawp/LICENSE and .hawp/kit/** (leave .hawp/work/ untouched)
- refresh .github/instructions/*.instructions.md
- refresh .github/prompts/*.prompt.md
- do not overwrite .github/copilot-instructions.md if it already exists
- provide a final summary of changed files and any manual merge points

Source:
- owner: sentzunhat
- repo: human-ai-workflow-protocol
- ref: main
```

## One-Shot Update Command (Copy/Paste)

Run this from the target repository root. It is **safe to re-run** and **safe on a repo that has `hawp/` or `.hawp/` from an earlier layout**. All `.hawp/work/` content is preserved.

```bash
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
REF="main"

TMP_DIR="$(mktemp -d)"
curl -fsSL "https://github.com/${OWNER}/${REPO}/archive/refs/heads/${REF}.tar.gz" \
  | tar -xz -C "$TMP_DIR"

SRC="$TMP_DIR/${REPO}-${REF}/core"

# --- Migration: legacy hawp/ -> .hawp/ (preserves hawp/work/) ---
# Only acts on a real directory; symlinks are left untouched so updates stay deterministic.
if [ -d "hawp" ] && [ ! -L "hawp" ]; then
  mkdir -p .hawp
  if [ -d "hawp/work" ]; then
    mkdir -p .hawp/work
    # cp -Rn = no-clobber; existing .hawp/work/ files always win.
    cp -Rn hawp/work/. .hawp/work/ 2>/dev/null || true
  fi
  [ -f hawp/LICENSE ] && [ ! -f .hawp/LICENSE ] && cp hawp/LICENSE .hawp/LICENSE
  rm -rf hawp
fi

# --- Migration: pre-work/ layout (.hawp/status/) -> .hawp/work/status/ ---
if [ -d ".hawp/status" ] && [ ! -d ".hawp/work/status" ]; then
  mkdir -p .hawp/work/status
  cp -Rn .hawp/status/. .hawp/work/status/ 2>/dev/null || true
  rm -rf .hawp/status
fi

# --- Refresh HAWP kit content (preserve LICENSE at root and any local work/) ---
rm -rf .hawp/kit
mkdir -p .hawp/kit
cp "$SRC/.hawp/LICENSE" .hawp/
cp "$SRC/.hawp/kit/README.md" .hawp/kit/
cp "$SRC/.hawp/kit/START_HERE.md" .hawp/kit/
cp "$SRC/.hawp/kit/SPEC.md" .hawp/kit/
cp "$SRC/.hawp/kit/AUTHORING_PATTERNS.md" .hawp/kit/
cp -R "$SRC/.hawp/kit/templates" .hawp/kit/
cp -R "$SRC/.hawp/kit/patterns" .hawp/kit/
cp -R "$SRC/.hawp/kit/reviews" .hawp/kit/
cp -R "$SRC/.hawp/kit/examples" .hawp/kit/
cp -R "$SRC/.hawp/kit/types" .hawp/kit/
cp -R "$SRC/.hawp/kit/usage" .hawp/kit/

# --- Refresh overlay files ---
mkdir -p .github/instructions .github/prompts
cp "$SRC/.github/instructions/"*.instructions.md .github/instructions/
cp "$SRC/.github/prompts/"*.prompt.md .github/prompts/

# Seed only when missing (do not overwrite local custom instructions)
if [ ! -f .github/copilot-instructions.md ]; then
  cp "$SRC/.github/copilot-instructions.md" .github/copilot-instructions.md
fi

rm -rf "$TMP_DIR"
```

## Post-Update Checklist

1. Confirm `.github/copilot-instructions.md` references `.hawp/kit/START_HERE.md` and `.hawp/kit/templates/status-report.md`.
2. Confirm `.hawp/LICENSE` exists and contains the Apache 2.0 text.
3. Confirm expected prompt files exist under `.github/prompts/` — including `intake.prompt.md`.
4. Confirm expected instruction files exist under `.github/instructions/` — including `intake.instructions.md`.
5. If a legacy `hawp/` directory existed, confirm it was migrated into `.hawp/` and that `.hawp/work/` retains every file you had in `hawp/work/`. Symlinks named `hawp` should be left untouched.
6. If a legacy `.hawp/status/` folder existed, confirm its contents now live under `.hawp/work/status/`.
7. Review git diff before committing — pay special attention to anything inside `.hawp/work/` (should be a clean migration with no deletions).
8. Run your repo checks (lint/test/typecheck) if your workflow requires it.

## Notes

- This update flow is conservative and keeps local Copilot instruction customizations.
- `.hawp/` is laid out flat at the repository root with `.hawp/LICENSE` and `.hawp/kit/`.
- `.hawp/work/` (if present) is repo-local operating state and is left completely untouched after the legacy migrations run. Scaffold README files and `BACKLOG.md` seeded by install are treated as repo-owned content and are never overwritten.
- Legacy `hawp/` migration uses `cp -Rn` to preserve every existing `.hawp/work/` file. Symlinks named `hawp` are skipped on purpose.
- Legacy `.hawp/status/` (pre-`work/` layout) is auto-moved into `.hawp/work/status/` so older installs land on the current layout without manual steps.
- The Apache 2.0 kit license travels with `.hawp/LICENSE` and is refreshed by this update flow.
