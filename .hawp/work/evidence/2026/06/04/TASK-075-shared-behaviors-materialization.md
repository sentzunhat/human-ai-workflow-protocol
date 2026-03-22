# Shared provider behavior materialization — 2026-06-04

Implemented TASK-075: canonical behaviors under `core/providers/shared/behaviors/` materialize into 11 provider overlay files.

## Acceptance

- [x] Four shared behavior sources (core, intake, backlog, docs)
- [x] Emit map in `librarian/scripts/providers/materialize/composition.ts`
- [x] `providers:sync` materialize + validate
- [x] `distribution:sync` runs `providers:sync` first
- [x] CI path filters include `core/providers/shared/**` and materialize scripts
- [x] CI drift check includes materialized provider overlays
- [x] `distribution:validate` — 12/12 guides
- [x] Continue-only smoke install — 4 rules with generated banner

## Commands

```bash
npm --prefix librarian run providers:sync
npm --prefix librarian run distribution:sync
```

## Hand-maintained (not materialized)

- GitHub: `copilot-instructions.md`, prompts, `commit-style.instructions.md`, `intake.instructions.md`
- Cursor: `AGENTS.md.seed`
