# Install HAWP GitHub Overlay

Use this package when the target repository already has a repo-root `hawp/` folder and you want to add the matching GitHub Copilot integration under `.github/`.

## One-shot install + usage quickstart

If you want a single copy/paste installation flow, run this from the target repository root:

```bash
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
REF="main"

TMP_DIR="$(mktemp -d)"
curl -fsSL "https://github.com/${OWNER}/${REPO}/archive/refs/heads/${REF}.tar.gz" \
  | tar -xz -C "$TMP_DIR"

SRC="$TMP_DIR/${REPO}-${REF}/human-ai-workflow-protocol"

# Install protocol content under human-ai-workflow-protocol/
mkdir -p human-ai-workflow-protocol
cp -R "$SRC/hawp" human-ai-workflow-protocol/

# Install Copilot overlay into .github/
mkdir -p .github/instructions .github/prompts
cp "$SRC/github/instructions/"*.instructions.md \
  .github/instructions/
cp "$SRC/github/prompts/"*.prompt.md \
  .github/prompts/

# Seed instructions only when missing. If present, merge instead of overwrite.
if [ ! -f .github/copilot-instructions.md ]; then
  cp "$SRC/github/copilot-instructions.md" .github/copilot-instructions.md
fi

rm -rf "$TMP_DIR"
```

Quick usage after install:

1. In `.github/copilot-instructions.md`, ensure HAWP references point to your actual layout (`human-ai-workflow-protocol/hawp/...` or `hawp/...`).
2. Start task shaping from `human-ai-workflow-protocol/hawp/START_HERE.md` (or `hawp/START_HERE.md` in repo-root layout).
3. Save context-transfer artifacts under `.../hawp/usage/status/` using `.../hawp/usage/STATUS_REPORT.md`.

## Scope Clarification

This document covers only the GitHub Copilot overlay installation (`.github/` integration).

For minimal HAWP installation in the nested project layout, copy only:

- `github/`
- `hawp/`

into `human-ai-workflow-protocol/`.

The `benchmark/` folder is optional reference material and is not part of the minimal install pair unless explicitly needed.

## What This Package Installs

Install these files into the target repository's `.github/` folder:

- `copilot-instructions.md`
- `instructions/human-ai-workflow-protocol-intake.instructions.md`
- `instructions/human-ai-workflow-protocol-docs-alignment.instructions.md`
- `prompts/human-ai-workflow-protocol-status-report.prompt.md`
- `prompts/human-ai-workflow-protocol-intent-first-handoff.prompt.md`
- `prompts/human-ai-workflow-protocol-docs-alignment-deterministic.prompt.md`
- `prompts/human-ai-workflow-protocol-docs-alignment-simplicity.prompt.md`
- `prompts/human-ai-workflow-protocol-conservative-docs-drift-cleanup.prompt.md`

These files assume the target repository resolves HAWP content from repo-root paths such as `hawp/usage/INIT.md`.

## Preconditions

Before installing this overlay, confirm the target repository contains:

- `hawp/README.md`
- `hawp/SPEC.md`
- `hawp/AUTHORING_PATTERNS.md`
- `hawp/usage/INIT.md`
- `hawp/usage/STATUS_REPORT.md`

If `hawp/` is missing, install or copy the HAWP folder first. Do not install this `.github/` overlay against a repository that does not have the expected repo-root `hawp/` paths.

## Install Procedure

1. Check that the target repository has a repo-root `hawp/` directory.
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

- If you copied this package as a temporary `github/` folder, remove that temporary folder after its files have been merged into `.github/`.
- Remove duplicate overlay files that were copied outside `.github/`, such as an extra `github/copilot-instructions.md`, `github/instructions/`, or `github/prompts/` tree left in the target repo.
- Keep only the repo-root `hawp/` folder and the target repo's `.github/` integration files. The target steady state should not require a sibling `github/` package folder.
- Do not remove existing non-HAWP `.github/` files that belong to the target repository.

## HAWP Block For Existing Copilot Instructions

Add or merge this block into the target repository's `.github/copilot-instructions.md`:

```md
This repository uses HAWP as a lightweight workflow method.

Follow the repo-local HAWP guidance in:

- hawp/usage/INIT.md
- hawp/usage/STATUS_REPORT.md

Use hawp/usage/INIT.md as the operating guide for how this repo applies HAWP in practice.

Use hawp/usage/STATUS_REPORT.md when the user asks for a:

- status report
- checkpoint summary
- context transfer summary
- second-brain review artifact

Saved status reports belong in:

- hawp/usage/status/

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

- `.github/copilot-instructions.md` references `hawp/usage/INIT.md`
- `.github/copilot-instructions.md` references `hawp/usage/STATUS_REPORT.md`
- `.github/instructions/human-ai-workflow-protocol-intake.instructions.md` exists
- `.github/instructions/human-ai-workflow-protocol-docs-alignment.instructions.md` exists
- `.github/prompts/human-ai-workflow-protocol-status-report.prompt.md` exists
- `.github/prompts/human-ai-workflow-protocol-intent-first-handoff.prompt.md` exists
- `.github/prompts/human-ai-workflow-protocol-docs-alignment-deterministic.prompt.md` exists
- `.github/prompts/human-ai-workflow-protocol-docs-alignment-simplicity.prompt.md` exists
- `.github/prompts/human-ai-workflow-protocol-conservative-docs-drift-cleanup.prompt.md` exists
- the referenced `hawp/usage/status/` directory exists in the target repository
- no temporary `github/` overlay folder remains in the target repository unless the user is intentionally keeping the kit source there

## Failure Cases

- Missing `hawp/` folder: stop and install the HAWP content first.
- Missing `.github/` folder: create it before copying the overlay files.
- Existing `.github/copilot-instructions.md` with unrelated repo rules: merge carefully and do not replace project-specific instructions.
- Temporary package files left in the target repository: remove the extra `github/` overlay folder and keep only the installed `.github/` files.

## Intent

This package is a modular overlay. It adds HAWP-aware Copilot behavior without replacing the target repository's baseline instruction set.

## Path Adaptation Note

The package source files (`github/copilot-instructions.md`, `github/instructions/`, `github/prompts/`) use portable `hawp/...` paths. These assume `hawp/` sits at the target repository root.

When the target repository stores HAWP content at a nested path — for example, `human-ai-workflow-protocol/hawp/` — the installer must adapt all `hawp/...` path references in the installed `.github/` files to match the actual layout.

The `human-ai-workflow-protocol` repository itself is an example of this: its installed `.github/` files use `human-ai-workflow-protocol/hawp/...` paths rather than the portable `hawp/...` defaults, because HAWP content lives under a subdirectory rather than at the repo root.

If you are comparing the package source files against installed `.github/` files and see path differences, they are intentional path adaptation — not unintended drift.
