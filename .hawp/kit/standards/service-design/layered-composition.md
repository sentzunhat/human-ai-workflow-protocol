# Layered Composition

## Principle

Composition preserves layer boundaries and enables adapter replacement without rewriting domain behavior.

## Applies To

- application/domain/infrastructure composition
- adapter wiring and replacement
- service module boundaries
- runtime portability decisions

## Guidance

- Keep domain logic independent from platform and transport details.
- Isolate framework and infrastructure concerns behind adapters.
- Prefer composition patterns that allow controlled replacement.

## Port and Adapter Placement (Go)

**Ports** (Go interfaces) live with the layer that *defines the need*:
- Domain decisions need external behavior → `domain/providers/<feature>/` (explicit `providers/` subfolder)
- Only application orchestration needs it → `application/providers/<feature>/`
- Infrastructure-internal reusable contracts (e.g. a shared SQLite transaction port used across
  multiple repository adapters) → colocated in `infrastructure/repositories/<entity>/` alongside the adapter

The `providers/` subfolder under domain and application makes contracts predictable: given a feature
name you always know where to look, regardless of which layer owns the need. Infrastructure ports are
the exception — they stay colocated with their adapter, not under a `providers/` tree.

**Adapters** (struct implementations) live in infrastructure, named after the domain they serve:
- Persistence → `infrastructure/repositories/<domain>/` (e.g. `repositories/work/`)
- Model runtimes → `infrastructure/models/<provider>/` (e.g. `models/ollama/`)
- External clients → `infrastructure/clients/<service>/`

**Naming consistency rule:** the adapter folder name mirrors the domain concept name.
`domain/providers/work/` defines the port; `infrastructure/repositories/work/` provides the adapter.

**Bootstrap** (`internal/bootstrap/`) is the only place that imports both a port's owner
and its adapter to wire them together. CLI/MCP entrypoints call bootstrap; they do not
construct adapters directly.

```
domain/providers/work/           # BacklogRepository interface (port)
  └── port.go

infrastructure/repositories/work/ # FilesystemBacklogRepository (adapter)
  ├── port.go                     # optional: shared SQLite port reused by sibling adapters
  └── repository.go

bootstrap/
  └── wire.go                    # var repo work.BacklogRepository = &workfs.FilesystemBacklogRepository{...}
```

## Does Not Include

- exact startup sequence requirements
- framework bootstraps
- project folder contracts
- runtime vendor lock-in assumptions

## Related

- [service-boundaries.md](service-boundaries.md)
- [handler-responsibilities.md](handler-responsibilities.md)
- [dependency-composition.md](dependency-composition.md)
