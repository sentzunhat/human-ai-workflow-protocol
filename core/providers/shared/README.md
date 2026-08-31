# Shared provider behaviors

Canonical HAWP **agent integration** text - one source per behavior, materialized into each provider pack.

## Layout

```text
shared/
├── behaviors/           # edit these (provider-neutral markdown bodies)
│   ├── hawp-core.md
│   ├── hawp-intake.md
│   ├── hawp-backlog-alignment.md
│   └── hawp-docs-alignment.md
└── README.md
```

Emit maps live in `librarian/src/internal/domain/providersync/materialize.go`.

## What gets materialized

| Behavior | Claude Code | GitHub | Cursor | Continue |
|----------|-------------|--------|--------|----------|
| Core | `hawp-core.md` | - (see `copilot-instructions.md`) | `hawp-core.mdc` | `hawp-01-core.md` |
| Intake | `hawp-intake.md` | `hawp-intake.instructions.md` | `hawp-intake.mdc` | `hawp-03-intake.md` |
| Backlog | `hawp-backlog-alignment.md` | `hawp-backlog-alignment.instructions.md` | `hawp-backlog-alignment.mdc` | `hawp-02-backlog-alignment.md` |
| Docs | `hawp-docs-alignment.md` | `hawp-docs-alignment.instructions.md` | `hawp-docs-alignment.mdc` | `hawp-04-docs-alignment.md` |

**Hand-maintained** (not generated): GitHub `copilot-instructions.md`, prompts, commit-style; Cursor `AGENTS.md.seed`; Codex `AGENTS.md.seed`; Claude `CLAUDE.md.seed`.

## Workflow

Normative workflow policy stays in `.hawp/kit/` and `.hawp/kit/references/`. Shared behaviors are **thin integration overlays** that point agents at the kit.

```bash
cd librarian/src && go run ./cmd/hawp providers sync
cd librarian/src && go run ./cmd/hawp distribution sync
```

`librarian/` is a checked-in maintenance workspace, not a separately installed dependency. Run these commands from the HAWP repo root when you want to refresh the generated provider packs and guides.

CI should fail when materialized provider files drift from `shared/behaviors/`.

## Adding a behavior

1. Add `behaviors/<name>.md`
2. Register targets in `materialize/composition.ts` (frontmatter + output path per provider)
3. Run `hawp providers sync`

## Adding a provider

1. Add the provider source pack under `core/providers/.<name>/`
2. Add emit targets in `librarian/src/internal/domain/providersync/materialize.go` if the provider consumes shared behaviors
3. Add install/update fragments under `distribution/sources/providers/<name>/`
4. Register the provider in `ACTIVE_PROVIDERS`
5. Update provider indexes and `core/providers/manifest.yaml`
6. Run `hawp distribution sync`
