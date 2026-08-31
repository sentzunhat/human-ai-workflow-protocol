# Provider distribution sources

Per-agent install and update guide fragments. Composed into `distribution/generated/<provider>/<install|update>/<main|dev>.md`.

| Provider | Status | Sources |
| -------- | ------ | ------- |
| `claude` | active | `providers/claude/` |
| `codex` | active | `providers/codex/` |
| `github` | active | `providers/github/` |
| `cursor` | active | `providers/cursor/` |
| `continue` | active | `providers/continue/` |

Each provider guide installs `.hawp/kit/**` plus that provider's overlay from `core/providers/.<provider>/`.

## Source layout (distribution)

```text
sources/
├── install/script-core.md          # kit + migrations (shared)
├── install/script-footer.md
├── update/script-core.md
├── update/script-footer.md
└── providers/<name>/
    ├── safety.md                   # provider overlay safety (tone + paths)
    ├── boundaries.md
    ├── install-contract.md
    ├── update-contract.md
    ├── script-install.md           # overlay only for this provider
    └── script-update.md
```

Generated bash = core + provider script + footer (see `librarian/src/internal/domain/distribution/distribution.go`).
