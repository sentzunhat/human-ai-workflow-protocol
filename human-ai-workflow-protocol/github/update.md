# Update HAWP From GitHub Main

Use this when you already installed HAWP in a target repository and want the latest package content from the GitHub `main` branch.

## What This Updates

This refreshes only HAWP-managed files:

- `hawp/**`
- `.github/instructions/*.instructions.md`
- `.github/prompts/*.prompt.md`

It does not overwrite an existing `.github/copilot-instructions.md` file.

## Prompt You Can Paste Into GitHub Copilot Chat

```text
Update this repository's HAWP kit from the latest github main branch of human-ai-workflow-protocol.

Requirements:
- refresh hawp/**
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

Run this from the target repository root:

```bash
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
REF="main"

TMP_DIR="$(mktemp -d)"
curl -fsSL "https://github.com/${OWNER}/${REPO}/archive/refs/heads/${REF}.tar.gz" \
  | tar -xz -C "$TMP_DIR"

SRC="$TMP_DIR/${REPO}-${REF}/human-ai-workflow-protocol"

# Refresh HAWP protocol content
rm -rf hawp
cp -R "$SRC/hawp" .
mkdir -p hawp/usage/status

# Refresh overlay files
mkdir -p .github/instructions .github/prompts
cp "$SRC/github/instructions/"*.instructions.md .github/instructions/
cp "$SRC/github/prompts/"*.prompt.md .github/prompts/

# Seed only when missing (do not overwrite local custom instructions)
if [ ! -f .github/copilot-instructions.md ]; then
  cp "$SRC/github/copilot-instructions.md" .github/copilot-instructions.md
fi

rm -rf "$TMP_DIR"
```

## Post-Update Checklist

1. Confirm `.github/copilot-instructions.md` references `hawp/usage/INIT.md` and `hawp/usage/STATUS_REPORT.md`.
2. Confirm expected prompt files exist under `.github/prompts/`.
3. Confirm expected instruction files exist under `.github/instructions/`.
4. Review git diff before committing.
5. Run your repo checks (lint/test/typecheck) if your workflow requires it.

## Notes

- This update flow is conservative and keeps local Copilot instruction customizations.
- `hawp/` is always installed flat at the repository root.
