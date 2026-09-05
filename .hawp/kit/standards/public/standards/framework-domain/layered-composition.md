# Layered Composition

## Principle

Composition should preserve layer boundaries and enable adapter replacement without rewriting domain behavior.

## Applies To

- Application/domain/infrastructure composition
- Adapter wiring and replacement
- Service module boundaries
- Runtime portability decisions

## Guidance

- Keep domain logic independent from platform or transport details.
- Isolate framework and infrastructure concerns behind adapters.
- Prefer composition patterns that allow controlled replacement.

## Port and Adapter Placement

**Ports** (interfaces) live with the layer that defines the need:
- Domain decisions need external behavior → `domain/providers/<feature>/`
- Only application orchestration needs it → `application/providers/<feature>/`
- Infrastructure-internal reusable contracts (e.g. a shared DB port used across
  multiple repository adapters) → colocated in `infrastructure/repositories/<entity>/`

The `providers/` subfolder under domain and application makes contracts predictable.
Infrastructure ports are the exception — colocated with their adapter, not under a `providers/` tree.

**Adapters** live in infrastructure, named after the domain they serve:
- Persistence → `infrastructure/repositories/<domain>/`
- Model runtimes → `infrastructure/models/<provider>/`
- External clients → `infrastructure/clients/<service>/`

The adapter folder name mirrors the domain concept name so both the port and its
adapter are always findable from the domain concept alone.

**Bootstrap** is the only layer that imports both sides to wire them together.
Entrypoints (CLI, MCP) call bootstrap; they do not construct adapters directly.

## Does Not Include

- Exact startup sequence requirements
- Framework-specific bootstraps
- Project folder contracts
- Runtime vendor lock-in assumptions

## Related

- [README.md](README.md) — framework-domain mirror index
- [service-boundaries.md](service-boundaries.md) — Contract clarity
- [handler-responsibilities.md](handler-responsibilities.md) — Boundary adapters
- [dependency-registration.md](dependency-registration.md) — Composition setup
