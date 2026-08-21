# c1d2e400 — Recursive audit: CLI command capabilities

**Type:** audit
**Status:** done
**Updated:** 2026-08-21

## Goal

Map every command/subcommand in `librarian/go/internal/platform/cli/run.go` to its capability group, application services, and infrastructure dependencies. Identify seams that need cleanup before splitting `run.go` into per-group files (`c1d2e401`).

## Scope

- Read-only audit of `librarian/go/internal/platform/cli/`
- Produce command-to-capability map, provider-construction findings, direct-infra findings, split order, and verification plan
- No production code changes

---

## Audit Evidence

### Files audited

| File | Lines | Role |
| ---- | ----- | ---- |
| `run.go` | 1226 | Main router + all command handlers |
| `commands.go` | 71 | Registry renderer + `runCommands` |
| `registry.go` | 49 | `CommandInfo` struct + `Registry` slice |
| `benchmark.go` | 299 | `RunBenchmark` + benchmark helpers; imports `sqlite` directly |

---

### Command-to-capability map

| Command | Group | Application services | Infrastructure (direct) | Approx lines in run.go |
| ------- | ----- | -------------------- | ----------------------- | ---------------------- |
| `uuid [--short]` | util | appuuid | — | ~12 (146–157) |
| `links check` | links | applinks | repo | ~10 (159–168) |
| `links clean [--apply]` | links | applinks | repo | ~22 (181–202) |
| `kit validate` | kit | appkit | repo | ~20 (204–224) |
| `work validate` | work | appwork | repo | ~31 (226–256) |
| `work new` | work | appwork, appuuid | repo | ~58 (263–320) |
| `check` | util | appcheck | repo | ~10 (593–602) |
| `init` | util | appprovision | download (inline) | ~12 (322–333) |
| `version` | util | domainupdate (const) | — | ~3 (88–89) |
| `update verify` | update | appupdate | githubrelease (inline) | ~17 (337–354) |
| `update latest` | update | appupdate | download, githubrelease (inline) | ~26 (358–383) |
| `update sync` | update | appkitsync | download, githubrelease (via doKitSync) | ~4 (388–392) |
| `update` (full) | update | appupdate, appkitsync | download, githubrelease (inline) | ~27 (396–423) |
| `commands [--json]` | util | — (Registry only) | — | ~2 (103–104) |
| `backlog validate` | work | — (alias → runCheck) | — | ~2 (106–108) |
| `backlog upgrade` | work | — (alias → runWorkNormalize) | — | ~2 (109–111) |
| `db init` | util | appdb | filesystem (inline constructor) | ~7 (112–119) |
| `index build` | index | appindex | repo | ~38 (472–509) |
| `model pull` | index | appembed | filesystem via modelsRoot() | ~25 (521–545) |
| `embed` | index | appembed | filesystem via modelsRoot() | ~45 (547–591) |
| `kit normalize` | kit | appkit | repo | ~27 (604–632) |
| `work normalize` | work | appwork | repo | ~37 (634–671) |
| `search index` | search | appindex | repo, os.ReadFile (direct walk) | ~96 (673–777 incl. helpers) |
| `search embed` | search | appindex | sqlite.Open (direct) | ~65 (783–847) |
| `search benchmark` | search | — (RunBenchmark) | sqlite.Open (direct) | ~15 (1117–1134) |
| `search <query>` | search | appsearch, appcontext | sqlite.Open (direct) | ~155 (849–1005) |

Helper functions also in run.go (not commands, but count toward line total):
- `parseProviderFlags`, `doKitSync`, `sortedProviderNames` — update helpers (~45 lines)
- `buildCorpusFromRepo`, `walkKitFiles`, `walkWorkFiles`, `strPtr` — search/index helpers (~75 lines)
- `tryReshapeViaRAGPipeline`, `renderReshapedWithReferences`, `toJSONReferences` — search helpers (~90 lines)
- `getStr`, `getInt`, `getFloat`, `mustGetwd`, `helpText` — util helpers (~60 lines)

**Total confirmed:** 25 command cases in the `Run` switch + ~270 lines of helpers.

---

### Provider construction findings

Every finding below is a place where a command handler constructs its own services inline instead of receiving them from a shared constructor. These are the seams `c1d2e401` must clean up.

| Location (run.go line) | Finding |
| ---------------------- | ------- |
| 113 | `db init`: `appdb.NewInitService(filesystem.NewLayoutService())` — layout service constructed inline |
| 327 | `runInit`: `download.NewHTTPFetcher()` constructed inline |
| 339 | `runUpdateVerify`: `githubrelease.NewClient()` constructed inline |
| 359 | `runUpdateLatest`: `githubrelease.NewClient()` + `download.NewHTTPFetcher()` constructed inline |
| 399–400 | `runUpdateFull`: same — both clients constructed inline again |
| 441 | `doKitSync`: `download.NewHTTPFetcher()` constructed inline (third copy) |
| 495 | `runIndexBuild`: `appindex.NewBuildService(root)` constructed inline |
| 517–519 | `modelsRoot()`: calls `filesystem.ResolveHawpHome(home).Models` directly — filesystem adapter accessed from CLI layer |
| 689–696 | `runSearchIndex`: `appindex.NewIngestService(dbPath)` constructed inline; `dbPath` computed directly in handler |
| 795–806 | `runSearchEmbed`: `sqlite.Open(dbPath)` — SQLite infrastructure opened directly |
| 839 | `runSearchEmbed`: `appindex.NewEmbedService(dbPath)` constructed inline |
| 889–895 | `runSearch`: `sqlite.Open(dbPath)` — SQLite infrastructure opened directly |
| 1125–1133 | `runSearchBenchmark`: `sqlite.Open(dbPath)` — SQLite infrastructure opened directly |

**Top finding:** `download.NewHTTPFetcher()` and `githubrelease.NewClient()` are each constructed 3–4 times across different update handlers. `sqlite.Open(dbPath)` is constructed 3 times across search handlers. Both patterns should be moved to shared constructors or injected.

---

### Direct infrastructure access findings

Infrastructure packages imported directly in the CLI layer (bypassing application services):

| Package | Used in | How |
| ------- | ------- | --- |
| `infrastructure/sqlite` | `run.go` (runSearch, runSearchEmbed, runSearchBenchmark), `benchmark.go` | `sqlite.Open(dbPath)`, then direct DB method calls (`QueryChunksLexical`, `HasVectors`, `ChunksNeedEmbedding`) |
| `infrastructure/filesystem` | `run.go` (modelsRoot, db init) | `filesystem.ResolveHawpHome()`, `filesystem.NewLayoutService()` |
| `infrastructure/download` | `run.go` (runInit, runUpdateLatest, runUpdateFull, doKitSync) | `download.NewHTTPFetcher()` — constructed inline 4 times |
| `infrastructure/githubrelease` | `run.go` (runUpdateVerify, runUpdateLatest, runUpdateFull, runUpdateSync/doKitSync) | `githubrelease.NewClient()` — constructed inline 3 times; also accepted as arg in `doKitSync` |
| `infrastructure/repo` | throughout run.go | `repo.FindBacklogRepoRoot()`, `repo.Exists()` — acceptable in CLI layer (repo discovery is CLI's job) |

Direct `os.ReadFile` and `filepath.Walk` calls also occur inside `walkKitFiles` and `walkWorkFiles` (lines 722–777) — file I/O logic that belongs in an application or infrastructure layer, not in the CLI.

**`benchmark.go` is also an infrastructure violator:** it imports `sqlite` directly and calls `db.QueryChunksLexical` and `db.HasVectors` inside `RunBenchmark`. This file is part of the CLI package but performs infrastructure work directly.

---

## Handoff To c1d2e401

### Target file layout

After `c1d2e401` completes, the `cli/` package should contain:

| File | Contents |
| ---- | -------- |
| `run.go` | `Run` function (router) + `ExitError` type only; all imports trimmed to what the router needs |
| `cmd_links.go` | `runLinksCheck`, `runLinksClean` |
| `cmd_kit.go` | `runKitValidate`, `runKitNormalize` |
| `cmd_work.go` | `runWorkValidate`, `runWorkNormalize`, `runWorkNew` |
| `cmd_update.go` | `runUpdateVerify`, `runUpdateLatest`, `runUpdateSync`, `runUpdateFull`, `parseProviderFlags`, `doKitSync`, `sortedProviderNames` |
| `cmd_util.go` | `runUUID`, `runCheck`, `runInit`, `mustGetwd`, `helpText` |
| `cmd_index.go` | `runIndexBuild`, `runModelPull`, `runEmbed`, `modelsRoot` |
| `cmd_search.go` | `runSearch`, `runSearchIndex`, `runSearchEmbed`, `runSearchBenchmark`, `buildCorpusFromRepo`, `walkKitFiles`, `walkWorkFiles`, `strPtr`, `tryReshapeViaRAGPipeline`, `renderReshapedWithReferences`, `toJSONReferences`, `getStr`, `getInt`, `getFloat` |
| `commands.go` | Unchanged (`runCommands`, `renderCommandsText`, `renderCommandsJSON`) |
| `registry.go` | Unchanged |
| `benchmark.go` | Unchanged for now; the sqlite dependency is a known seam to address in a follow-up |

### Safe mechanical split order

Each step is one file creation + deletion from `run.go`. Each step must leave the build green before proceeding.

**Step 1 — `cmd_links.go`**
Extract `runLinksCheck`, `runLinksClean`. Only depends on `applinks`, `repo`. Simplest: no infrastructure imports beyond `repo`.

**Step 2 — `cmd_kit.go`**
Extract `runKitValidate`, `runKitNormalize`. Only depends on `appkit`, `repo`.

**Step 3 — `cmd_work.go`**
Extract `runWorkValidate`, `runWorkNormalize`, `runWorkNew`. Only depends on `appwork`, `appuuid`, `repo`.

**Step 4 — `cmd_update.go`**
Extract all update functions + helpers. Bring in `download`, `githubrelease`, `appupdate`, `appkitsync`, `domainupdate` imports. Helpers `parseProviderFlags`, `doKitSync`, `sortedProviderNames` move with this file.

**Step 5 — `cmd_util.go`**
Extract `runUUID`, `runCheck`, `runInit`, `mustGetwd`, `helpText`. Brings in `appcheck`, `appuuid`, `appprovision`, `download`, `domainupdate` (for version in helpText indirectly — helpText is pure text so no import needed there).

**Step 6 — `cmd_index.go`**
Extract `runIndexBuild`, `runModelPull`, `runEmbed`, `modelsRoot`. Brings in `appindex`, `appembed`, `domainindex`, `filesystem`.

**Step 7 — `cmd_search.go`**
Extract all search functions + helpers. Largest step: brings in `sqlite`, `appsearch`, `appcontext`, `domainsearch`. This is where the direct sqlite access lives; extract as-is first, then a follow-up can clean up the seams.

**Step 8 — `run.go` cleanup**
Remove all imports no longer needed in `run.go`. Confirm only the router remains.

### Verification plan for c1d2e401

Run after **each step**:

```
# From librarian/go/
go build ./...
go vet ./internal/platform/cli/...
go test ./internal/platform/cli/...
```

All three must pass before proceeding to the next step. If a step breaks the build, revert that step and investigate import cycles or missed function moves before retrying.

After all steps complete, also confirm:
- `run.go` contains only `Run`, `ExitError`, and the `helpText` call (or `helpText` moved to cmd_util.go with an equivalent call)
- Each `cmd_*.go` file compiles independently (confirmed by `go build ./...`)
- No new infrastructure imports appear in `run.go`
- `benchmark.go` sqlite seam is documented as a deferred cleanup item (not addressed in c1d2e401)
