# GitHub / Copilot Overlay Safety

This guide installs the GitHub provider pack only.

- Refreshes `.github/instructions/` and `.github/prompts/` from `core/providers/.github/` on every install and update.
- Seeds `.github/copilot-instructions.md` on first install if missing; refreshes it on update.
- Does not modify `.cursor/`, `AGENTS.md`, or `.continue/`.
