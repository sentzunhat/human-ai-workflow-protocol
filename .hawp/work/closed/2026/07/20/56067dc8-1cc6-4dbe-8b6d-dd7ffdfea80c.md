# establish Go librarian CLI as the primary librarian project

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `56067dc8-1cc6-4dbe-8b6d-dd7ffdfea80c`
**Type:** improvement
**Reported:** 2026-07-20
**Risk Level:** medium

---

### Input (what was reported)

> The librarian CLI will run as a Go project with all the scripts and "magic"
> — a plain Go binary, no embedded local models and no Node.js bundle.
> The `librarian/go/` workspace becomes the primary librarian implementation.

---

### Context

The repo currently has two CLI workspaces: `nodejs/` and `golang/`. The
`librarian/go/` workspace is a Zacatl-shaped scaffold (~1.6 MB binary proof, no
ONNX or database yet). The 2026-07-20 decision is to build the librarian as
a Go project: the binary stays small because models and the ONNX Runtime are
downloaded after `init` rather than embedded (see
`e98de8c4-d6be-43f0-a5b6-57c939ddd195`).

Related items this supersedes in direction:
`1a3b32a4-ab37-4c86-ade7-71e2eb42b440` (transformers.js Node WASM spike) and
the Node SEA portion of `54e68af7-4622-4383-8482-cc4d4e1e21ee`.

---

### Analysis

**Root cause (or most likely cause):**
_Node packaging (SEA, bundled runtimes) fights the small-binary goal; Go
gives a small static binary and simple cross-platform releases, with models
kept out of the binary entirely._

**Directly verified:**
_`librarian/go/` builds a ~1.6 MB `hawp` binary with layered scaffold
(`cmd/hawp`, `internal/application`, `internal/domain`,
`internal/infrastructure`, `internal/platform`). No database or model
integration exists yet._

**Inferred (not yet proven):**
_The existing librarian scripts (kit/work validate, normalize, uuid,
distribution build) can be ported to or wrapped by the Go CLI without
regressions._

**Scope — what else is affected:**
_`nodejs/` workspace role, `librarian/` npm scripts, `.hawp/bin/hawp`
wrapper, install/update docs in `distribution/`, CI workflows._

---

### Recommended Fix

- Define which librarian capabilities the Go CLI owns first (db init, ingest,
  lexical search) vs. what stays in `librarian/` npm scripts short-term.
- ~~Decide the fate of the `nodejs/` workspace~~ — resolved 2026-07-20:
  `nodejs/` removed and Go scaffold consolidated into `librarian/go/`
  (`6a427d8c-5e3d-4d1b-bc3a-3a2e21bc83a7`, closed).
- Keep the binary model-free: runtime assets are provisioned by `init`
  (`e98de8c4`), updates come from published releases (`cdcf9f78`).

**Port plan (2026-07-20):** the TypeScript audit, per-command Go mapping,
phase order (0 foundations → 1 read-only → 2 mutating → 3 composite +
switchover), and the Go unit test standard are in
`.hawp/work/notes/2026/07/20/librarian-ts-to-go-port-plan.md`.

**Progress (2026-07-20):** Phases 0–2 are closed — Phase 0 in this plan's
scope, Phase 1 as `39bc92b6` (uuid, links check, kit validate, work
validate; count parity verified), Phase 2 as `eddd8339` (kit normalize,
work normalize + backlog upgrade alias; mutation with clean-tree guard).

**Phase 3 (2026-07-20, this item):** composite `hawp check` implemented
(kit + work + links, `backlog validate` now aliases it);
`.hawp/bin/hawp` switched to prefer the Go binary with npm fallback and
the downstream uuid-only behavior preserved; Node CLI PoC retired
(`scripts/hawp-cli-poc/`, `cli:poc*` npm scripts, `@oclif/core`,
`.github/workflows/hawp-cli-poc-binaries.yml` all removed). TS validators
stay as the npm-pipeline fallback until CI switches to the Go binary
(noted in Future Improvements).

**What to verify after:**

- [x] Go CLI capability boundary vs. npm scripts is documented
      (Evidence: command mapping table in the 2026-07-20 port plan note)
- [x] `nodejs/` workspace decision is recorded
      (Evidence: `closed/2026/07/20/6a427d8c-5e3d-4d1b-bc3a-3a2e21bc83a7.md`)
- [x] Binary remains small and model-free
      (Evidence: `make build` 2026-07-20 — bin/hawp ≈ 2.7 MB stripped, no embedded
      assets; `make dist` produced six ~1.6–1.9 MB platform binaries)
- [x] Each ported phase has Go unit tests matching TS test coverage before
      the TS equivalent is removed
      (Evidence: 9 Go test files; only the PoC TS was removed — TS
      validators remain as npm fallback, 58 TS tests still pass)

---

## Outcome (filled at close)

Closed 2026-07-20. The Go librarian CLI is the primary implementation:
all workflow commands (uuid, links check, kit/work validate + normalize,
check, backlog aliases) are native, `.hawp/bin/hawp` routes through the
Go binary with npm fallback, and the Node CLI PoC is fully retired. The
intelligence lane (`e98de8c4`, `cdcf9f78`, `fbf12a93`) builds on this
foundation. Follow-up noted in Future Improvements: remove the TS
validator fallback once CI runs the Go binary.

## Verification (filled at close)

- Evidence: `make check` in `librarian/go/` — vet, all tests, build pass
  (2026-07-20).
- Evidence: `./.hawp/bin/hawp` smoke run 2026-07-20 — uuid, kit validate,
  work validate, check, backlog validate, work normalize dry-run all
  routed through the Go binary with expected output; bare-repo copy still
  falls back to the uuid-only guidance.
- Evidence: `npm --prefix librarian run typecheck` and `npm test`
  (58 tests, 0 failures) pass after the PoC removal.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Plan file saved under `closed/2026/07/20/56067dc8-1cc6-4dbe-8b6d-dd7ffdfea80c.md`
- [x] BACKLOG.md updated
