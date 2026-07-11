# Continue.dev provider pack

Installs HAWP guidance as local workspace rules under `.continue/rules/`.

| Source | Target |
|--------|--------|
| `rules/hawp-*.md` | `.continue/rules/` |

Local rules load automatically in Continue (Agent, Chat, Edit). No Hub `config.yaml` is required for this pack.

See [Continue rules docs](https://docs.continue.dev/customize/deep-dives/rules) and [doc-verification.md](../doc-verification.md).

Rules use lexicographic order (`hawp-01-core.md` … `hawp-04-docs-alignment.md`).
