# Claude Code Overlay Safety

This guide installs the Claude Code provider pack only.

- Refreshes `.claude/rules/hawp-*.md` from `core/providers/.claude/rules/` on every install and update.
- Seeds `CLAUDE.md` from `CLAUDE.md.seed` on first install only when missing; never overwrites on update.
- Does not modify `.github/`, `.cursor/`, or `.continue/`.
- Non-HAWP files already in `.claude/rules/` are left unchanged.
