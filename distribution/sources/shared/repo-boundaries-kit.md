# Kit and Work Boundaries (All Providers)

Every provider install/update guide refreshes the agent-neutral kit and preserves project work.

## Kit (always installed)

| Source | Target | Behavior |
|--------|--------|----------|
| `core/.hawp/kit/**` | `.hawp/kit/**` | Full refresh every install/update |
| `core/.hawp/LICENSE` | `.hawp/LICENSE` | Refreshed |
| `core/.hawp/work/` scaffold | `.hawp/work/` READMEs, `BACKLOG.md` seed | Seed only when missing |

## Project-Owned (never overwritten)

- `.hawp/work/**` — backlog, active/parked/closed work, decisions, evidence, notes

## Never Installed Downstream

- HAWP source repo's `.hawp/work/**` operating state
- `benchmark/` reference material

Provider-specific overlay boundaries are documented in the next section for this guide.
