# Multi-repo context skill

**When to use:** When working on the HAWP source repo in a way that affects
installed files, or when working in a downstream consumer repo and needing to
understand where the installed files came from and how to update them.

---

## The Three-Lane Model

HAWP spans three lanes. Every change starts in one and propagates forward:

```
┌─────────────┐         ┌──────────────────┐         ┌─────────────────┐
│  SOURCE     │  install │  DOWNSTREAM     │  update │  CONSUMER       │
│  (this repo)│─────────▶│  .hawp/kit/**   │─────────▶│  (user's repo) │
│  core/      │  /update │  .claude/       │  (re-run)│  .hawp/work/** │
│  .hawp/kit/ │          │  .github/       │          │  (project state)│
│  dist/      │          │  .cursor/       │          │                 │
└─────────────┘          │  .continue/     │          └─────────────────┘
                         │  AGENTS.md      │
                         │  CLAUDE.md      │
                         └──────────────────┘
```

| Lane                | Owned by              | Examples                                                           |
| ------------------- | --------------------- | ------------------------------------------------------------------ |
| Source              | HAWP maintainer       | `core/providers/`, `.hawp/kit/`, `distribution/sources/`           |
| Downstream scaffold | Install/update script | `.hawp/kit/**`, `.claude/rules/`, `.github/instructions/`          |
| Consumer state      | User / agent          | `.hawp/work/BACKLOG.md`, `.hawp/work/active/`, decisions, evidence |

**Key invariant:** `.hawp/work/**` is never overwritten by install or update.
Kit and provider overlays refresh on every run.

---

## Propagation Rules

### What changes propagate downstream

| Source change                             | Downstream effect                                      |
| ----------------------------------------- | ------------------------------------------------------ |
| `.hawp/kit/**` edit                       | Full refresh to every consumer's `.hawp/kit/**`        |
| `core/providers/<name>/**` edit           | Refresh to matching provider overlay in consumer       |
| `distribution/sources/**` edit            | Regenerate guides via `hawp distribution sync`         |
| `core/providers/shared/behaviors/**` edit | Sync via `hawp providers sync`; then distribution sync |

### What does NOT propagate

- `.hawp/work/BACKLOG.md` — seed only when missing
- `.hawp/work/active/**`, `parked/**`, `closed/**` — never touched
- `.hawp/decisions/`, `.hawp/evidence/`, `.hawp/status/` — never touched

### Propagation path for a typical change

```
1. Edit source file (e.g. .hawp/kit/start-here.md)
2. (If kit source)  → no sync step needed
3. (If provider)   → go run ./cmd/hawp providers sync
4. (If distribution) → go run ./cmd/hawp distribution sync
5. Commit + push to source repo
6. Consumer runs update guide → kit + provider overlay refresh
7. Consumer's .hawp/work/** is preserved
```

---

## Provider Overlay Map

Source pack → consumer target (from `core/providers/manifest.yaml`):

| Provider    | Source                      | Consumer target                                                                | Install                   | Update                               |
| ----------- | --------------------------- | ------------------------------------------------------------------------------ | ------------------------- | ------------------------------------ |
| Claude Code | `core/providers/.claude/`   | `.claude/rules/hawp-*.md`, `CLAUDE.md`                                         | refresh / seed-if-missing | refresh / skip (CLAUDE.md preserved) |
| Codex       | `core/providers/.codex/`    | `AGENTS.md`                                                                    | seed-if-missing           | seed-if-missing                      |
| GitHub      | `core/providers/.github/`   | `.github/instructions/`, `.github/prompts/`, `.github/copilot-instructions.md` | refresh / seed            | refresh / refresh                    |
| Cursor      | `core/providers/.cursor/`   | `.cursor/rules/*.mdc`, `AGENTS.md`                                             | refresh / seed            | refresh / seed                       |
| Continue    | `core/providers/.continue/` | `.continue/rules/hawp-*.md`                                                    | refresh                   | refresh                              |

Each install/update guide refreshes **only its own provider overlay**, not all
of them. A Cursor guide never writes to `.github/`.

---

## Common Cross-Repo Failure Modes

| Symptom                                        | Root cause                                          | Fix                                                       |
| ---------------------------------------------- | --------------------------------------------------- | --------------------------------------------------------- |
| Consumer kit is stale after source change      | Update guide not re-run                             | Run matching provider's update guide in consumer          |
| Provider overlay missing after install         | Wrong provider guide used                           | Run the guide matching the consumer's provider            |
| `.hawp/work/BACKLOG.md` lost after update      | Manual overwrite or wrong install                   | Re-seed from source: `copy_file_no_clobber` semantics     |
| `distribution/generated/` out of sync          | `distribution sync` not run after source edit       | `cd librarian/src && go run ./cmd/hawp distribution sync` |
| Provider rules out of sync                     | `providers sync` not run after shared behavior edit | `cd librarian/src && go run ./cmd/hawp providers sync`    |
| Consumer shows "already present" without proof | Agent skipped execution                             | Require `Source:`, `Provider:`, `Source mode:` in output  |
| `CLAUDE.md` overwritten on update              | `seed-if-missing` → `skip` in manifest              | Expected behavior; verify manifest.yaml                   |

---

## Working in a Downstream Consumer Repo

When the agent is in a repo that **has** `.hawp/` installed:

1. **Do not edit `.hawp/kit/` directly.** Those files are managed by the
   install/update script. Source changes must go through the HAWP source repo.
2. **Do not edit provider overlay files directly.** `.claude/rules/`,
   `.github/instructions/`, etc. are managed by the install script.
3. **`/memories/`** (if present) is a local, session-scoped store. It is not
   part of the install contract and does not propagate.
4. **`.hawp/work/`** is the user's project state. It is the only part of the
   installed HAWP layer the user owns and edits directly.
5. To apply a source-side change to this consumer, the user must re-run the
   matching update guide (not the install guide) from this repo's root.

---

## Source Repo Change Checklist

Before committing a change in the HAWP source repo, verify propagation:

```bash
# From librarian/src/
go test ./...
go run ./cmd/hawp providers sync
go run ./cmd/hawp distribution sync
go run ./cmd/hawp kit validate
go run ./cmd/hawp kit normalize --apply
go run ./cmd/hawp work validate
go run ./cmd/hawp work normalize --dry-run --validate
go run ./cmd/hawp check
```

All checks must pass before pushing. CI enforces this via
`sync-distribution-generated.yml`.

---

## Reference Files

- `core/providers/manifest.yaml` — provider overlay map
- `distribution/sources/shared/install.md` — install contract
- `distribution/sources/shared/update.md` — update contract
- `distribution/sources/shared/repo-boundaries-kit.md` — kit/work boundaries
- `distribution/sources/install/script-core.md` — install script mechanics
- `.hawp/kit/start-here.md` — kit entry point
- `.hawp/kit/usage/workflow-loop.md` — multi-iteration workflow
