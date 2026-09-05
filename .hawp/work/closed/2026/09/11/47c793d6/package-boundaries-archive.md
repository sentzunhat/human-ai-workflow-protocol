# Package boundaries — target map and CLI implementation convention

## Decision

Keep layers, then group by capability. A Go directory is a package boundary;
move cohesive behavior and its tests together. Do not make one package per file,
mirror every folder across every layer, or introduce generic utils/common trees.
CLI families use operation subfolders, including earlier work, kit, index and
model slices. Handlers, private parsers and their tests live together; family
roots expose forwarding entrypoints. Future CLI slices follow this convention.
Other destinations below remain planned until their implementation is recorded.
The previous whole-source [migration preview](../../../../scripts/source-layout/review/preview.md)
and [tool guide](../../../../scripts/source-layout/README.md) enumerate the current
mechanical baseline; no broad source migration has been applied. The checkpoint
adds context configuration/deduplication folders to the mapping but deliberately
does not regenerate that saved preview before the requested commit.

## Target ownership

| Layer | Owns | Must not own |
| --- | --- | --- |
| `internal/domain` | Values, invariants, pure decisions, and port interfaces in `domain/providers/<feature>/` | Filesystem walks/writes, SQL, HTTP, model sessions, provider factories |
| `internal/application` | Use-case orchestration; port interfaces in `application/providers/<feature>/` when only application needs them | CLI parsing, provider selection, concrete default adapter construction |
| `internal/infrastructure` | Adapter implementations of ports, named after the domain they serve (`repositories/<domain>/`, `models/<provider>/`) | Port interface definitions, work lifecycle decisions, or CLI routing |
| `internal/platform` | CLI/MCP translation, validation, rendering, command routing | Business policy or direct persistence where a use case owns it |
| `internal/bootstrap` | Explicit wiring: imports both port owners and adapters, constructs services | Business rules; importing CLI/MCP packages back into the composition root |

**Port/adapter rule:** domain and application ports live under an explicit `providers/` subfolder
in the layer that defines the need — `domain/providers/<feature>/` or `application/providers/<feature>/`.
Infrastructure-internal reusable contracts (e.g. a shared SQLite port used across multiple
repository adapters) are the exception — they colocate with their adapter in `infrastructure/repositories/<entity>/`.
The adapter (struct) lives in infrastructure under a folder named for the domain concept it serves.
Bootstrap is the only layer that imports both sides.
Adapter folder names mirror domain concept names — `domain/providers/work/` port → `infrastructure/repositories/work/` adapter.

Paths in this table are relative to `librarian/src`. Domain never imports the
other layers. Infrastructure implements contracts owned by consumers. Bootstrap
imports application/infrastructure; platform entrypoints use the resulting
services. No domain compatibility factory may import infrastructure.

## CLI family status

Convention: `cli/<family>/` owns forwarding entrypoints; `cli/<family>/<operation>/`
holds handler, private parser, and tests when the operation has its own parser.
Argless or flag-only operations stay flat in `<family>/commands.go`.
Maximum depth: two levels (`family/operation/`). No deeper nesting.

| Family | Status | Operations nested |
| --- | --- | --- |
| `work/` | ✓ done | `new/`, `normalize/`, `validate/` |
| `kit/` | ✓ done | `normalize/`, `validate/` |
| `index/` | ✓ done | `build/`, `ingest/` |
| `model/` | ✓ done | `pull/`, `embed/`, `search-embed/`, `location/` |
| `links/` | ✓ done | `check/`, `clean/` |
| `search/` | ✓ done | `query/`, `benchmark/` |
| `mcp/` | ✓ done | `configure/` (handler + invalid-arg tests); server mode flat in `commands.go` |
| `usage/` | ✓ done | `enable/` (parser + tests); disable/log/report/clear/totals flat in `commands.go` |
| `update/` | ✓ done | flat `commands.go`; `init` in `cli/init/command.go` |
| `distribution/` | ✓ done | flat `commands.go`; imports `providers` sibling for sync |
| `providers/` | ✓ done | flat `commands.go` |

Source files remaining at root after all migrations:
`run.go`, `registry.go`, `help.go`, `commands.go` (backlog/db/commands handlers),
`commands_test.go`, `init_test.go`, `embed_test.go`, `model_commands.go` (3 thin
forwarders — inline into `run.go` as a final cleanup), `maintenance_commands.go`
(uuid/check stay here or inline; link/kit forwarders already delegate to child pkgs).

## Planned directory map

```text
librarian/src/internal/
  platform/cli/             # router, registry/help, entrypoint only
    work/                  # forwarding entrypoints ✓
      new/                 # handler, parser, tests ✓
      normalize/           # handler, parser, tests ✓
      validate/            # handler, parser, tests ✓
    kit/                   # kit command family ✓
      normalize/           # normalize args + handler ✓
      validate/            # validate args + handler ✓
    links/                 # ✓
      check/               # handler ✓
      clean/               # handler, parser, tests ✓
    search/                # ✓
      query/               # handler, parser, tests ✓
      benchmark/           # handler, parser, benchmark logic ✓
    index/                 # forwarding entrypoints ✓
      build/               # index build handler and parser ✓
      ingest/              # search index handler, corpus and tests ✓
    model/                 # forwarding entrypoints ✓
      pull/                # handler, parser, tests ✓
      embed/               # handler, parser, tests ✓
      search-embed/        # handler and parser ✓
      location/            # shared model-directory resolver ✓
    mcp/                   # completed CLI family
      configure/           # mcp configure: handler, parser, tests
    usage/                 # completed CLI family
      enable/              # usage enable: handler, parser (only op with parser)
    update/                # completed — flat commands.go
    distribution/          # completed — flat commands.go
    providers/             # completed — flat commands.go
  platform/mcp/            # mechanical proposal, not applied
    configure/             # provider configuration and tests
    server/                # JSON-RPC lifecycle, tools, and server tests
  application/context/
    dedup/                 # deeper proposal: result deduplication and its tests
                           # rendering/reshaping retain shared private helpers
  application/work/
    intake/                # creation + draft orchestration
    validation/            # load facts, evaluate, assemble report
    normalize/             # preview/apply orchestration
  application/providers/work/      # port interfaces needed only by application (not domain)
  domain/work/             # shared values/IDs only; no imports of children
    backlog/               # pure table parsing/rendering
    intake/                # shape/template rules
    validation/            # pure checks against supplied facts
    normalize/             # pure change planning, no disk mutation
  domain/providers/work/           # BacklogRepository + ValidationRepository ports (domain-defined)
  domain/providers/embeddings/     # existing model port, proposed location
  domain/providers/llm/            # existing model port, proposed location
  infrastructure/
    clients/               # mechanical proposal, not applied
      download/            # existing HTTP download client
      githubrelease/       # existing release API client
    repositories/          # adapter implementations — folder name = domain concept served
      index/               # existing SQLite index, proposed location
      work/                # planned adapter for domain/providers/work/
      usage/               # implements UsageStore port from domain/usage/
    models/                # model runtime adapters — grouped by provider identity
      ollama/              # embedding + LLM HTTP implementations
      onnx/                # embedding + LLM runtime implementations
      none/                # explicit no-model adapters
  bootstrap/               # wires ports to adapters; only place importing both sides
```

**Infrastructure naming rule:** folder names reflect the *domain they serve*, not
the mechanism. `repositories/work/` is the work persistence adapter (currently
filesystem); the adapter can gain a `sqlite/` sibling later without renaming the
domain folder. `repositories/usage/` replaces `sqlite/usage/` for the same reason.
`models/` is mechanism-grouped because provider identity (ollama, onnx, none) is
the meaningful distinction for consumers.

This is a destination map, not permission to create all directories at once.
The 2026-09-10 whole-source preview includes the existing SQLite index under
`repositories/index` and external clients under `clients`, matching the canonical
adapter naming rules. It preserves their behavior and Go package names.
Context deduplication now has a focused nested-owner proposal; configuration
stays cohesive with its type methods.
Rendering/reshaping stay together; configuration I/O still needs later consumer
ports. This file grouping does not establish a pure domain or complete extraction.

## Three implementation items

1. [CLI grouping](../../closed/2026/09/08/e7a5e294/plan.md): start only with work. Preserve existing
   application APIs while moving parsers/handlers. This establishes the pattern.
2. [Work boundaries](../../closed/2026/09/10/742aa60b/plan.md): pure backlog/input rules first, then
   application subpackages and `infrastructure/repositories/work`. Keep facts and exact status semantics intact.
3. [Runtime adapters](../../closed/2026/09/09/b61770a6/plan.md): Ollama first, then ONNX/none, then usage
   storage. Move factories to bootstrap and update consumers without cycles.

These are owned subitems of architecture item 47c793d6. Pause new reshape runtime
wiring during topology work; preserve its verified draft contract.

## Repeatable implementation gate

- Read only the selected slice, its imports/callers, tests, and this map.
- Record exact source/destination file paths before moving; move tests alongside.
- Child packages must not import a parent facade that imports them.
- Preserve public command syntax/output, errors, build tags, and stored data.
- Share code only when repeated behavior and ownership are established; do not
  merge distinct draft/intake/context types solely because fields look similar.
- Focused tests plus compilation after the first move; full tests/vet/HAWP checks
  once at the completed slice boundary. No model downloads for mechanical moves.
- One slice per turn, short result, one updated plan; no repeated broad audits.
- Preserve current dirty work. Plan-ready is not implementation-complete.

## Evidence

Inspected source imports and file inventory on 2026-09-07. CLI has mixed command
owners; domain/work currently performs filesystem operations; domain/embeddings
and domain/llm contain HTTP/ONNX implementations; domain/usage owns SQL storage.
These are actual ownership issues, not file-count targets. Relevant guidance:
`.hawp/kit/standards/service-design/layered-composition.md`,
`.hawp/kit/standards/service-design/handler-responsibilities.md`, and
`.hawp/kit/instructions/clean-code-and-structure.md`.
