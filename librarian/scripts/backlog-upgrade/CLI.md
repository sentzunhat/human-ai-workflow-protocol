# Backlog Upgrade CLI Contract

This file documents the CLI adapter contract for backlog-upgrade.

## Command Shape

```bash
hawp backlog upgrade [OPTIONS]
```

## Script Shortcut

From the `librarian/` directory:

```bash
npm run workflow:normalize
```

## Modes

- `--dry-run` (default): run detection and report findings without modifying project files.
- `--apply`: normalize closed records in place, scaffold missing close sections, and preserve ambiguous legacy files for manual follow-up.

## Options

- `--format <text|json>`: render report as text (default) or JSON.
- `--output <path>`: write rendered report to file.
- `--export-plan <path>`: write generated plan JSON to file.
- `--validate`: run workflow validation summary after dry-run or apply.
- `--force-dirty`: skip the apply-mode dirty-tree guard.
- `--verbose`: prints parsed options and script execution notices.
- `--help`, `-h`: print help text.
- `--version`, `-v`: print version text.

## Architecture Boundary

- Script logic lives in `script.ts` and is reusable without CLI assumptions.
- `cli.ts` is an adapter that parses argv and delegates to `script.ts`.
- `index.ts` is the executable boundary that sets process exit code.

## Future Extension Pattern

When adding new commands:

1. Add behavior to script-level modules first.
2. Extend parser/adapter behavior in `cli.ts`.
3. Keep process-level exit behavior only in `index.ts`.
