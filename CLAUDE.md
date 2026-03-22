@AGENTS.md

## Claude Code

This repo uses HAWP as a lightweight task-shaping protocol. Rules for Claude Code are in `.claude/rules/`.

### Commands

```bash
# Full validation (typecheck + tests + links + kit + distribution sync + workflow)
npm --prefix librarian run validate

# By group
npm --prefix librarian run typecheck             # TypeScript check
npm --prefix librarian test                      # unit tests only
npm --prefix librarian run check:markdown-links  # local link check across all .md files

npm --prefix librarian run providers:sync        # materialize shared behaviors → provider packs
npm --prefix librarian run distribution:sync     # providers:sync + build generated guides + validate

npm --prefix librarian run kit:validate           # validate .hawp/kit/ structure (naming, links, required files)
npm --prefix librarian run kit:normalize          # normalize .hawp/kit/ names and internal links
npm --prefix librarian run work:validate         # backlog + plan + evidence integrity
npm --prefix librarian run work:normalize        # normalize work records (apply + validate)
npm --prefix librarian run hawp:check            # combined distribution + work:validate in one command
```

### Project layout

| Path | Purpose |
|------|---------|
| `.hawp/kit/` | Canonical HAWP kit — templates, patterns, standards, usage guides |
| `.hawp/work/` | Active backlog, plans, evidence, status reports |
| `core/providers/` | Provider overlay sources (GitHub Copilot, Cursor, Continue, Claude Code) |
| `distribution/generated/` | Generated install/update guides — do not edit directly |
| `distribution/sources/` | Source fragments composed into generated docs |
| `librarian/scripts/` | Tooling: distribution build, provider materialization (TypeScript) |
| `librarian/src/` | Go source for the `hawp` CLI — build with `cd librarian/src && make build` |
| `benchmark/` | Benchmark runs comparing HAWP vs. no-HAWP |

### Key conventions

- Edit `core/providers/shared/behaviors/` then run `providers:sync` — never edit generated provider overlays directly.
- `distribution/generated/` is fully generated. Edit `distribution/sources/` then run `distribution:sync`.
- Generated files are validated by CI (`sync-distribution-generated.yml`). A passing build requires clean sync.
- Work items tracked in `.hawp/work/BACKLOG.md`; active plans in `.hawp/work/active/`; close to `.hawp/work/closed/YYYY/MM/DD/`.
- Evidence goes in `.hawp/work/evidence/YYYY/MM/DD/`; status snapshots in `.hawp/work/status/YYYY/MM/DD/`.
- Prefer compact, decision-useful outputs. Separate direct evidence from inference.
