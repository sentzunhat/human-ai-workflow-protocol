# add agent-facing command/usage discovery to hawp CLI (ship now, not gated on vector search)

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `b4c8af81-177f-4e37-a572-5b6d8ee988c2`
**Type:** feature
**Reported:** 2026-07-21
**Risk Level:** low

---

### Input (what was reported)

> Provide the agent commands and usage to hawp right of the bat — don't
> wait for vector search/context building to land first.

---

### Context

Today an agent (or human) discovers `hawp`'s surface only via
`hawp --help`, which is free-text and not machine-parseable. This item
adds a structured, agent-consumable command registry so any agent using
`hawp` can introspect its capabilities programmatically instead of
scraping help text.

---

### Analysis

**Root cause (or most likely cause):**
_Free-text `--help` is fine for humans but forces an agent to either
hardcode knowledge of the CLI surface or fragile-parse help output._

**Directly verified:**
_`internal/platform/cli/run.go` already has a complete, hand-maintained
list of every command in `helpText()`; no structured equivalent exists._

**Inferred (not yet proven):**
_A single static command registry (name, usage, description, flags,
exit-code meaning) can back both the human `--help` text and a new
`hawp commands --json` output, so the two never drift apart._

**Scope — what else is affected:**
_`internal/platform/cli/` (new registry + route), `librarian/go/README.md`
(agent usage section)._

---

### Recommended Fix

- Add a static `[]CommandInfo` registry describing every command: name,
  usage line, one-line description, flags, and exit-code semantics.
- Add `hawp commands` (text) and `hawp commands --json` (machine-readable)
  reading from the same registry.
- Add an "Agent Usage" section to `librarian/go/README.md` pointing at
  `hawp commands --json` as the canonical discovery path, with exit-code
  conventions (0 success, 1 issues found, 2 usage/dirty-tree error) and
  example invocations.
- Ship in the `v0.0.2` test release (`4c152ee3`) so the update mechanism
  test has real content to upgrade to.

**What to verify after:**

- [x] `hawp commands --json` output is valid JSON covering every command
      in `helpText()`
      (Evidence: `TestRunCommandsJSONIsValidAndComplete` unmarshals the
      output and checks command count + specific names; real run
      pretty-printed valid JSON with 17 commands)
- [x] `hawp commands` (text) and `--help` stay consistent (same command
      list, no drift)
      (Evidence: `TestRegistryStaysInSyncWithHelpText` fails the build if
      any available registry entry's base verb is missing from
      `helpText()`)
- [x] Go unit tests cover the registry and both output formats
      (Evidence: `commands_test.go` — 5 tests covering text output,
      JSON validity/completeness, the `--json` route, and planned-vs-
      available marking)

---

## Outcome (filled at close)

Closed 2026-07-21. `hawp commands` (text) and `hawp commands --json`
are implemented, backed by a single static `Registry` in
`internal/platform/cli/registry.go` — 17 entries (15 available, 2
planned: `search`, `context`) each with usage, description, flags, and
exit-code semantics. A test enforces the registry and `helpText()` never
drift apart. Documented in `librarian/go/README.md` under a new "Agent
usage" section with the shared exit-code convention and a `jq` example.
Shipped as the real content of the `v0.0.2` test release
(`4c152ee3-09af-40c0-b372-cff849d063cd`).

## Verification (filled at close)

- Evidence: `go test ./internal/platform/cli/...` — all 5 new tests plus
  existing CLI route tests pass.
- Evidence: real `./bin/hawp commands` and `./bin/hawp commands --json`
  runs produced correct, complete, valid-JSON output on this machine.
- Evidence: `make check` (vet + full test suite + build) passes.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Plan file saved under `closed/2026/07/21/b4c8af81-177f-4e37-a572-5b6d8ee988c2.md`
- [x] BACKLOG.md updated
