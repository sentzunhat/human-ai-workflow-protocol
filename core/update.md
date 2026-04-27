# Update HAWP From GitHub Main

Use this when you already installed HAWP in a target repository and want the latest package content from the GitHub `main` branch.

## What This Updates

This refreshes only HAWP-managed files:

- installable `.hawp/` public kit files and folders
- `.github/instructions/*.instructions.md`
- `.github/prompts/*.prompt.md`

Update boundary: source-repo maintenance paths under `.hawp/usage/` (including `.hawp/usage/status/`) are intentionally excluded from this default update flow.

If a legacy `hawp/` directory exists as a real directory (not a symlink), it is removed so `.hawp/` is the only active kit root.

It does not overwrite an existing `.github/copilot-instructions.md` file.

Because `.hawp/` is rebuilt from the public kit allowlist, `.hawp/LICENSE` is refreshed too.

## Prompt You Can Paste Into GitHub Copilot Chat

```text
Update this repository's HAWP kit from the latest github main branch of human-ai-workflow-protocol.

Requirements:
- refresh .hawp/**
- refresh .github/instructions/*.instructions.md
- refresh .github/prompts/*.prompt.md
- remove legacy hawp/ if it exists as a real directory (not symlink)
- do not overwrite .github/copilot-instructions.md if it already exists
- provide a final summary of changed files and any manual merge points

Source:
- owner: sentzunhat
- repo: human-ai-workflow-protocol
- ref: main
```

## One-Shot Update Command (Copy/Paste)

Run this from the target repository root:

```bash
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
REF="main"

TMP_DIR="$(mktemp -d)"
curl -fsSL "https://github.com/${OWNER}/${REPO}/archive/refs/heads/${REF}.tar.gz" \
  | tar -xz -C "$TMP_DIR"

SRC="$TMP_DIR/${REPO}-${REF}/core"

# Refresh HAWP protocol content
rm -rf .hawp
mkdir -p .hawp
cp "$SRC/.hawp/README.md" .hawp/
cp "$SRC/.hawp/START_HERE.md" .hawp/
cp "$SRC/.hawp/SPEC.md" .hawp/
cp "$SRC/.hawp/AUTHORING_PATTERNS.md" .hawp/
cp "$SRC/.hawp/LICENSE" .hawp/
cp -R "$SRC/.hawp/templates" .hawp/
cp -R "$SRC/.hawp/patterns" .hawp/
cp -R "$SRC/.hawp/reviews" .hawp/
cp -R "$SRC/.hawp/examples" .hawp/
cp -R "$SRC/.hawp/types" .hawp/

# Refresh overlay files
mkdir -p .github/instructions .github/prompts
cp "$SRC/.github/instructions/"*.instructions.md .github/instructions/
cp "$SRC/.github/prompts/"*.prompt.md .github/prompts/

# Seed only when missing (do not overwrite local custom instructions)
if [ ! -f .github/copilot-instructions.md ]; then
  cp "$SRC/.github/copilot-instructions.md" .github/copilot-instructions.md
fi

# Safety: remove only a real legacy hawp/ directory so updates are deterministic
# and .hawp/ remains the single active kit root. Never remove symlinks.
if [ -d "hawp" ] && [ ! -L "hawp" ]; then
  rm -rf hawp
fi

rm -rf "$TMP_DIR"
```

## Post-Update Checklist

1. Confirm `.github/copilot-instructions.md` references `.hawp/START_HERE.md` and `.hawp/templates/status-report.md`.
2. Confirm `.hawp/LICENSE` exists and contains the Apache 2.0 text.
3. Confirm expected prompt files exist under `.github/prompts/`.
4. Confirm expected instruction files exist under `.github/instructions/`.
5. If legacy `hawp/` existed as a real directory, confirm it was removed; if `hawp` is a symlink, confirm it was left untouched.
6. Review git diff before committing.
7. Run your repo checks (lint/test/typecheck) if your workflow requires it.

## Notes

- This update flow is conservative and keeps local Copilot instruction customizations.
- `.hawp/` is always installed flat at the repository root.
- Legacy `hawp/` cleanup is guarded: only real directories are removed, and symlinks are preserved.
- The Apache 2.0 kit license now travels with `.hawp/LICENSE` and is refreshed by this update flow.
