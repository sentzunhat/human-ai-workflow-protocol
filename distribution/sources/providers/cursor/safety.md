# Cursor Overlay Safety

This guide installs the Cursor provider pack only.

- Refreshes `.cursor/rules/hawp-*.mdc` from `core/providers/.cursor/rules/` on every install and update.
- Seeds `AGENTS.md` from `AGENTS.md.seed` only when missing, on both install and update.
- Does not modify `.github/` or `.continue/`.
- Non-HAWP rules already in `.cursor/rules/` are left unchanged.
