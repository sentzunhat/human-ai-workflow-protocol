### Status Report

#### Intent

Continue the runtime adapter boundary work and make the next compoundable wiring slice explicit.

#### Current State

Model adapters are under infrastructure. Bootstrap is at `internal/bootstrap/`, because it imports infrastructure and application packages and therefore cannot belong under `internal/domain/`. Usage ports, values, and report formatting are domain-owned; SQLite persistence and config file I/O are infrastructure-owned. Search and context now receive model factories through bootstrap as well.

#### What Was Inspected

Inspected the active `b61770a6` plan, package-boundaries map, model factory, embedding service, search/context construction paths, CLI entrypoint, and repository status.

#### What Changed

Moved bootstrap back out of domain, split `domain/usage/store.go` into pure domain usage types/policy and `infrastructure/repositories/usage` SQLite/config adapters, moved the usage tests with the adapter, updated platform callers, and injected model factories into index/search/context application services. Removed obsolete no-factory constructors.

#### What Was Directly Verified

Focused tests for usage, repository, index, bootstrap, CLI usage, MCP, search-embed, search, and context pass. Focused vet passes. No production application or domain package imports `infrastructure/models`. Full `go test ./...` and `go vet ./...` pass with live integration/benchmark tests excluded by their explicit build tags.

#### What Remains Unproven

Live `integration` and `benchmark` tests remain opt-in and require their external model services. The next architectural item is the first bounded work-domain separation slice in `742aa60b`.

#### Constraints

No commits, merges, binary restoration, or unrelated dirty-state changes were made. No model downloads or live Ollama checks were run.

#### Help Wanted

Review whether search and context should share one bootstrap-owned model provider or receive separate narrowly scoped factories before the next slice.

#### Suggested Next Step

Begin `742aa60b` with one pure work-domain rule cluster, then run the complete suite and HAWP validation.
