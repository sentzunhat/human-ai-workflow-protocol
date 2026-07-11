# Provider doc verification — 2026-06-04

Verified HAWP provider packs against vendor documentation (online sources).

## Sources checked

| Provider | Documentation URL | Date |
|----------|---------------------|------|
| GitHub / Copilot | https://docs.github.com/en/copilot/how-tos/configure-custom-instructions-in-your-ide/add-repository-instructions-in-your-ide | 2026-06-04 |
| GitHub / Copilot | https://code.visualstudio.com/docs/copilot/customization/custom-instructions | 2026-06-04 |
| Cursor | https://cursor.com/docs/context/rules | 2026-06-04 |
| Continue | https://docs.continue.dev/customize/deep-dives/rules | 2026-06-04 |

## Findings

**GitHub — aligned.** Paths match: `.github/copilot-instructions.md`, `.github/instructions/*.instructions.md` with `applyTo`, `.github/prompts/*.prompt.md` with `name`/`description`.

**Cursor — aligned.** `.cursor/rules/*.mdc` with `description`/`globs`/`alwaysApply`; root `AGENTS.md` as always-on alternative per Cursor docs.

**Continue — adjusted.** Renamed rules to `hawp-01-*` … `hawp-04-*` for lexicographic load order. `globs` use YAML array form. Local `.continue/rules/*.md` only; no Hub config required.

## Verification commands

```bash
npm --prefix librarian run distribution:validate
# continue-only smoke: .continue/rules/hawp-*.md × 4, no .github/.cursor
```

Canonical reference: `core/providers/DOC-VERIFICATION.md`
