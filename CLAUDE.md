@AGENTS.md

## Claude Code

This repo uses HAWP as a lightweight task-shaping protocol. Rules for Claude Code are in `.claude/rules/`.

### Commands

```bash
# Full validation
cd librarian/src && go test ./... && go run ./cmd/hawp distribution sync && go run ./cmd/hawp check

# By group
cd librarian/src && go run ./cmd/hawp links check
cd librarian/src && go run ./cmd/hawp providers sync
cd librarian/src && go run ./cmd/hawp distribution sync
cd librarian/src && go run ./cmd/hawp kit validate
cd librarian/src && go run ./cmd/hawp kit normalize --apply
cd librarian/src && go run ./cmd/hawp work validate
cd librarian/src && go run ./cmd/hawp work normalize --dry-run --validate
cd librarian/src && go run ./cmd/hawp check
```

### Project layout

| Path | Purpose |
|------|---------|
| `.hawp/kit/` | Canonical HAWP kit — templates, patterns, standards, usage guides |
| `.hawp/work/` | Active backlog, plans, evidence, status reports |
| `core/providers/` | Provider overlay sources (GitHub Copilot, Cursor, Continue, Claude Code) |
| `distribution/generated/` | Generated install/update guides — do not edit directly |
| `distribution/sources/` | Source fragments composed into generated docs |
| `librarian/src/` | Go source for the `hawp` CLI — build with `cd librarian/src && make build` |
| `benchmark/` | Benchmark runs comparing HAWP vs. no-HAWP |

### Key conventions

- Edit `core/providers/shared/behaviors/` then run `hawp providers sync` - never edit generated provider overlays directly.
- `distribution/generated/` is fully generated. Edit `distribution/sources/` then run `hawp distribution sync`.
- Generated files are validated by CI (`sync-distribution-generated.yml`). A passing build requires clean sync.
- Work items tracked in `.hawp/work/BACKLOG.md`; active plans in `.hawp/work/active/`; close to `.hawp/work/closed/YYYY/MM/DD/`.
- Evidence goes in `.hawp/work/evidence/YYYY/MM/DD/`; status snapshots in `.hawp/work/status/YYYY/MM/DD/`.
- Prefer compact, decision-useful outputs. Separate direct evidence from inference.
