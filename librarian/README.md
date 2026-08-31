# librarian

HAWP maintenance lives here, with the Go CLI in [src/](src/README.md) as the
single command surface for both user-facing workflow actions and repo
maintenance tasks.

## Common tasks

```bash
cd librarian/src
go test ./...
go run ./cmd/hawp providers sync
go run ./cmd/hawp distribution sync
go run ./cmd/hawp kit validate
go run ./cmd/hawp work validate
go run ./cmd/hawp check
```

## Provider and distribution maintenance

When shared behaviors, provider packs, or distribution fragments change:

1. Edit `core/providers/shared/behaviors/`, `core/providers/.<provider>/`, or `distribution/sources/`.
2. Run `go run ./cmd/hawp distribution sync` from `librarian/src`.
3. Commit the refreshed files under `core/providers/.{claude,cursor,continue,github}/` and `distribution/generated/`.

Provider registration now lives in Go:

- Distribution variants: `librarian/src/internal/domain/distribution/distribution.go`
- Provider materialization targets: `librarian/src/internal/domain/providersync/materialize.go`

## Binary

`make build` in `librarian/src` produces `src/bin/hawp`, and `make dist`
cross-compiles the release binaries. The old Node maintainer workspace and
TypeScript workflow scripts were retired on 2026-08-31 after the Go CLI
reached parity for workflow and distribution/provider maintenance.
