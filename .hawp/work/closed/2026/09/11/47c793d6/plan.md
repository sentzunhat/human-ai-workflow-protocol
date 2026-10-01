# v0.0.24 CLI decomposition and architecture audit continuation

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `47c793d6-1f02-4fde-b372-3ee20f37fb69`
**Type:** improvement
**Reported:** 2026-09-05

---

## Input (verbatim)

> v0.0.24 CLI decomposition and architecture audit continuation

## Intake Summary

Continue the earlier audit with bounded CLI decomposition on the 0.0.24 branch.

## Current Context

The earlier audit `5957aaf4` closed a queue, not every implementation slice.
Go is the maintained runtime; TypeScript workflow scripts are retired.

## Initial Analysis

**Directly verified:**

- `run.go` mixes routing, indexing, search formatting, updates, and usage storage.

**Inferred (not yet proven):**

- Separating command owners reduces review friction without new package APIs.

**Likely scope:**

- Separate handlers inside the existing CLI package; preserve signatures and
  behavior. Fix non-finite hybrid-ratio validation with a regression test.
- Record remaining findings in [audit.md](audit.md).

## Risk + Review Gate

**Risk:** low for mechanical moves; medium for the boundary validation fix.
**Gate:** User requested audit implementation. Preserve existing patch work,
downstream repository state, and Go package boundaries. No new dependencies,
merges or release publication. The latest user request authorizes local commits.

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/active/47c793d6/plan.md

## Next Step

2026-09-10 checkpoint-first cleanup: consolidate `scripts/` under source-layout,
retire the redundant MCP setup shell script in favor of the existing CLI, and
separate the migration entrypoint, engine, mapping rules, and review artifacts.
Extend only the proposed context deduplication package map with one
additional capability level; do not move runtime source now. Retain and label the
prior preview/plan as an older snapshot, then commit all reviewed changes before
running a new preview or apply. Include the already-recorded core launcher
retirement from `0cb0f9b0`; preserve installed native/legacy compatibility files.

Added `scripts/README.md` to make these standards and lessons reusable for future
maintenance tools. Commit the documentation before generating the new preview.

2026-09-10 continuation: fix and test the whole-source preview, align the package
map and port documentation, and generate a readable directory/file report.
The entrypoint is now `bash scripts/source-layout/run.sh --preview --check`;
`--report review/preview.md` and `--diff` expose destination and content changes.
The [tool guide](../../../../scripts/source-layout/README.md) explains the gates.
This iteration does not apply the broad migration or close semantic work items.

2026-09-09 scope expansion: the user requests an executable migration for the
whole `librarian/src` tree, including MCP, all layers, tests, and import updates.
Replace the limited work-only preview with `scripts/source-layout/`: inventory
every source file, record move/retain decisions, rewrite Go imports using syntax
trees, and validate real apply in an isolated copy. Keep operation nesting shallow.
Prepare the reviewed migration and continuation artifact; do not apply the broad
layout to the working source tree as part of preparing the script. Semantic work
rule extraction stays with `742aa60b`; naming moves alone cannot establish purity.

Wrapper simplification continuation: installers currently pair `hawp` (shell
launcher) with `hawp-bin` (native executable), and Cursor uses `hawp-mcp` to set
cwd. Add `mcp --repo-root` and generate direct native MCP commands for providers.
Keep old wrappers for existing configurations; deleting them requires coordinated
installer/update migration. Verify explicit root selection without changing cwd.

Commit-readiness continuation: strict standard-library `work normalize` parsing
rejects conflicting modes/roots, unknown arguments, and missing/empty path values
before mutation. Restore the tracked shell launcher and install native builds
into the existing `hawp-bin` slot; update CI to use the same layout. Commit
verified implementation, documentation, and binary groups without pushing or merging.

Indexing continuation: propagate corpus file-read failures, use portable relative
paths, and derive lifecycle classification only for active/parked/closed folders.
Supporting documents keep unknown status. Preserve index schema and provider
behavior; verify with temporary directories and deterministic unreadable entries.

Next compoundable slice: pure standard-library search argument parsing with
typed options and invalid-input tests; schema-aware backlog intake supporting
UUID-only and legacy-ID headers; add Go vet to existing CI. No CLI framework
or runtime dependency is needed. Validate before filesystem or provider work.

Continuation: move the search index contract and embedding metadata into
`internal/domain/search`, preserve SQLite compatibility via a type alias, and
verify injected repositories, error propagation, and resource closure. Existing
default adapter construction remains available. Provider implementation moves
and broader CLI parsing remain separate follow-ups.

- [x] Investigation and implementation scope recorded above
- [x] Split command owners and verify focused CLI tests
- [x] Reject non-finite hybrid ratios and add regression coverage
- [x] Run full tests, vet, HAWP checks, and portable release builds
- [x] Refresh the local 0.0.24 binary
- [x] CLI decomposition complete — all families migrated ([e7a5e294](../../closed/2026/09/08/e7a5e294/plan.md) closed 2026-09-08)
- [x] Runtime architecture layer moves: [b61770a6](../../closed/2026/09/09/b61770a6/plan.md) (model/usage adapters)
- [x] Work architecture layer moves: [742aa60b](../../closed/2026/09/10/742aa60b/plan.md) (work domain separation, closed 2026-09-10)

## Verification

2026-09-10 checkpoint-first follow-up: tool compilation, mapping/flag/revision-only
tests, vet, shell syntax, and current MCP/context tests pass after organizing the
tool and extending the map. No new preview/apply or migration fixture execution
ran before the requested checkpoint. Prior artifacts are explicitly stale and
revision-gated; fresh candidate checks are the next action after committing.

2026-09-10 source-layout preview: tool regressions/vet and current-source full
tests/vet pass. Candidate default and integration/benchmark compilation plus vet
pass. Preview leaves source unchanged; fixture apply, rollback, and already-applied
checks pass. All HAWP checks pass after repairing stale closed-plan links.
The [checkpoint](../../status/2026/09/10/source-layout-preview.md) records scope
and remaining runtime risks. The broad migration remains unapplied.

2026-09-06 provider evidence follow-up: Copilot's supplied transcript invoked
`work validate --repo-root` through a terminal, not MCP. Investigation confirmed
that this CLI handler silently ignored unknown flags. Extracted a small typed
parser rejecting unknown/missing/empty/conflicting arguments before validation;
kept supported `--work-root`, `--hawp-root`, and default discovery behavior.
This targeted correction is prompted by observed evidence; the separate
single-executable migration remains open.
Verification: full Go tests, vet, local install, distribution validation, kit
mirror parity, and diff checks passed. The rebuilt binary rejects the supplied
unsupported flag with exit 1; explicit `--hawp-root .hawp` passes five work
checks with zero issues/warnings. Connected HAWP MCP validation also passes.

2026-09-06 CLI decomposition wiring completion: All tracked command handlers now
use strict `flag.NewFlagSet` typed parsers — `work_new`, `index_build`,
`kit_validate`, `links_clean`, `kit_normalize`, `embed`, `model_pull`,
`search_embed`, `search_benchmark`, `usage_enable`. Aggressive import cleanup
that removed `sqlite` and `application/index` dependencies was reversed; all
builds and tests pass clean (`go vet ./...`, `go test ./...`).

Latest commit checkpoint (2026-09-05): implementation committed as `adcb583`;
native binary refreshed separately as `6f5d19c`. Shell launcher remains unchanged
and delegates to `hawp-bin`. Full Go tests, vet, HAWP check, distribution
validation, six-target builds, and diff whitespace checks passed. Local launcher
reports `0.0.24`. No PR was merged or published. Remaining audit work stays open.

Verified 2026-09-05, commands run from `librarian/src` unless stated otherwise:

- `go test ./internal/platform/cli ./tests/platform/cli`: pass after first split.
- `go test ./...`: pass, including invalid hybrid-ratio regression cases.
- `go vet ./...`: pass.
- `go run ./cmd/hawp check --no-update-check`: kit, work, and links pass;
  work has zero issues and zero warnings; 127 Markdown files checked.
- `make dist VERSION=0.0.24`: all six standard targets compile (darwin, linux,
  windows; each amd64 and arm64). Cross-target execution and ORT builds were
  not tested in this pass.
- `make install VERSION=0.0.24`: pass; root `.hawp/bin/hawp version` reports
  `0.0.24`, and a NaN search ratio exits 1 with the expected validation error.
- `git diff --check`: pass.

Remaining findings are explicitly open in [audit.md](audit.md). This pass adds
no performance claim and does not close the full architecture audit.

Repository-port continuation verification: focused search/SQLite/CLI tests,
`go test ./...`, `go vet ./...`, `make dist VERSION=0.0.24` (six standard
targets), `make install VERSION=0.0.24`, root binary `check --no-update-check`,
and `git diff --check` all passed after the contract and error-handling changes.
The full test suite includes environment-dependent integration skips; this is
not proof that every optional model backend executed. No new dependencies added.

Argument/intake/corpus continuation verification (2026-09-05):

- Pure search parser and schema-aware intake regressions pass. The first full
  run caught MCP's minimal ID/Title/Status table; that supported layout now has
  explicit coverage and the subsequent full suite passes.
- Corpus tests cover active/parked/closed, supporting records without inferred
  status, and kit/work/custom read failures with preserved underlying errors.
- Final `go test ./...`, `go vet ./...`, `git diff --check`, six-target
  `make dist VERSION=0.0.24`, and `make install VERSION=0.0.24` passed.
- Root binary rejects `search query --limit=0 --no-update-check` with exit 1.
- Quality CI now includes vet; its commands were tested locally. No hosted CI
  run or PR publication occurred in this continuation.
- Provider wiring, remaining command parsers, and persisted backlog enrichment
  remain open, as detailed in the audit queue.

Native MCP launch continuation (2026-09-05):

- Added a small dedicated MCP command handler with validated `--repo-root`.
  New provider configurations use the native executable and explicit root.
  See [launcher-audit.md](launcher-audit.md) for the remaining migration scope.
- Focused tests, `go test ./...`, `go vet ./...`, `make install VERSION=0.0.24`,
  six-target `make dist VERSION=0.0.24`, HAWP check, distribution validation,
  and `git diff --check` passed.
- The rebuilt native executable handled MCP initialize and work validation
  from a foreign working directory with this repository explicitly selected:
  all three validations passed, with zero work issues or warnings.
- Cross-platform compilation is not cross-platform runtime or live provider
  integration proof. Existing consumer configurations and wrappers are unchanged;
  no downstream repository edits, publishing, or PR merge occurred.

## 2026-09-06 index model consolidation

Investigation: `librarian/src/internal/domain/index/chunk.go` duplicates SQLite's
metadata model and defines an unused chunk shape without source ranges. SQLite
uses nullable folder context; changing it to a string would erase that distinction.
The earlier checkpoint's universal blocker claim is not established. Importing
both domain and adapter packages alone is not proof of a defect.

Plan: retain nullable context and source ranges in the domain model; replace
SQLite definitions with compatibility aliases; construct domain values in ingest.
Alternative: retain separate storage DTOs and explicit conversion, unnecessary
for these matching fields. Risk: low, internal mechanical consolidation without
schema or query changes. Can implement now: yes; source files are clean and this
continuation belongs to the existing architecture item. Preserve the three
pre-existing dirty checkpoint/rule files. Verify database round trips, full Go
tests, vet, build, and HAWP validation.

Baseline HAWP validation fails on orphan `active/c42e04b7/plan.md`; a closed
record with the same UUID exists. Review and preserve that source as supporting
history under its closed item before validating again.

Root proof before edits: `pwd` and `git rev-parse --show-toplevel` both returned
`<repo-root-abs>`; `git rev-parse --show-prefix` returned an empty line.
`git status --short` returned:

```text
 M .hawp/work/status/2026/09/06/0cb0f9b0-status.md
?? .claude/rules/architecture-state.md
?? .hawp/work/status/2026/09/06/checkpoint-2026-09-06-architecture-audit.md
```

## 2026-09-06 review continuation

Reviewed the existing uncommitted domain-model consolidation; preserve its
nullable context, source ranges, compatibility aliases, and database tests.
Baseline full Go tests fail `TestRunEmbedRequiresTextArg`: the embed parser
separates flag values from flags and discards parse errors. Next bounded action:
preserve interspersed text/flags using FlagSet parsing, propagate errors, and
test values, unknown options, missing values, and explicit `--` text. No model
download is needed for parser verification.

Baseline MCP and source CLI checks also report an empty orphan active folder
and unrecognized preserved intake. Remove only the empty directory and rename
the preserved intake to the supported archive filename; retain its bytes and link.

Verification for this continuation: `go test ./...` and `go vet ./...` pass.
The existing domain SQLite round-trip test is included. Source CLI
`go run ./cmd/hawp check --no-update-check` and connected
`hawp_work_validate` pass kit/work/links with zero issues/warnings.
`go run ./cmd/hawp distribution validate`, root/core kit comparison, and
`git diff --check` pass. No binary refresh, model inference, release, or
downstream client verification was performed. Existing uncommitted work is
preserved; this continuation is uncommitted.

Next bounded review: query-first `model pull` parser compatibility. Its current
FlagSet stops at the repository positional argument, so trailing documented
options need focused regression coverage before further parser cleanup.

## 2026-09-07 model pull argument compatibility

Investigation: `librarian/src/internal/platform/cli/model_pull_args.go` calls
FlagSet.Parse on all arguments and rejects the trailing tokens as extra
positionals. The documented repository-first command in
`librarian/src/internal/platform/cli/registry.go` therefore rejects
`org/repo --onnx-file model.onnx`.

Plan: parse leading options, consume one repository, then parse trailing
options; preserve leading options and reject empty repository/explicit empty
ONNX path before model-directory discovery or downloads. Keep this private
parser local; no shared parser framework. Add pure regression cases and align
command registry/help, source README, and audit references. Risk: low.
Can implement now: yes; the parser is clean and the earlier dirty files belong
to the preserved continuation. Apply handler-responsibilities and docs-alignment
guidance from `.hawp/kit/standards/service-design/handler-responsibilities.md`
and `.hawp/kit/references/docs-alignment.md`.
Root proof: pwd/top-level = `<repo-root-abs>`; prefix empty. Git status showed
prior uncommitted work from the 2026-09-06 continuation; no commit or reset.

Verification: the new model-pull table failed before the fix in five cases
(documented trailing option, equals form, mixed flags, trailing boolean, empty
repository). All 16 model-pull cases and existing embed cases pass after the
fix. Full `go test ./...`, `go vet ./...`, distribution validation, kit parity,
and connected `hawp_work_validate` pass. Source README and command registry
now describe the accepted option order. No model downloads/inference, binary
refresh, commit, or release were performed.

Next: compare remaining typed parsers with their documented examples, and
reconcile the migration plan's closure statement with its active backlog row
using acceptance evidence. Request shaping belongs to
[the separate reshape item](../a3df8a9c/plan.md).

## 2026-09-07 mutation-boundary follow-up

Investigation: `kit_normalize_args.go` registers but discards `--dry-run`, so
`--apply --dry-run` sets apply=true. `work_new_args.go` has an empty Visit
callback and a false comment claiming FlagSet rejects empty strings; an
explicit `--hawp-root=` reaches current-repo discovery. Both kit parsers
also accept an empty explicit `--kit-path`. These are confirmed in
`librarian/src/internal/platform/cli/`; source handlers consume the parsed
values before calling the relevant services.

Plan: bind dry-run and reject both enabled modes, matching work normalize;
reject explicitly supplied empty/whitespace roots in work new, kit validate,
and kit normalize while retaining omitted-root defaults. Keep local FlagSet
checks; no framework or filesystem-layout change. Add parser cases and
handler-level temporary-fixture snapshots proving rejection without mutation.
Risk: low, bounded validation before side effects. Can implement now: yes.
Existing UUID and dirty work are preserved; root proof remains
`<repo-root-abs>` with empty git prefix and the prior continuation's dirty paths.
Guidance: `.hawp/kit/standards/service-design/handler-responsibilities.md` and
`.hawp/kit/references/docs-alignment.md`.

Verification: parser regressions failed before the fix for conflicting modes
and explicit empty paths. After the fix, focused parser/handler tests pass;
four handler cases snapshot all fixture file contents and directory entries
before/after rejection. Full `go test ./...` and `go vet ./...` pass.
Connected `hawp_work_validate` passes kit/work/links with zero issues/warnings;
distribution validation, kit parity, and `git diff --check` pass. Source README,
command registry, audit, and work dashboard are aligned. No dependencies,
model downloads, binary refresh, commit, or external repository edits.

The migration plan's contradictory closure record remains unresolved; this
bounded parser correction does not establish migration/release readiness.

## 2026-09-07 guidance cleanup

Replaced stale `.claude/rules/architecture-state.md` blocker/type/deferral
claims with source-confirmed context and current work links. Preserved its full
prior text in [architecture-state-prior.md](architecture-state-prior.md).
Corrected `.github/copilot-instructions.md`: Go imports belong to their declaring
file; package-wide scans apply to symbol changes, not import ownership. This
is a guidance cleanup under the existing architecture item, not a new runtime
policy or a claim that all earlier work is complete. Migration reconciliation
and downloader proof are tracked in the existing `0cb0f9b0` item.

2026-09-07 migration follow-up: the packaging audit resolved the retained source
launcher gate; [0cb0f9b0](../../closed/2026/09/07/0cb0f9b0/plan.md) is now
closed locally. This supersedes earlier unresolved-closure notes above.
Architecture work remains open; local migration completion is not publication.

## 2026-09-07 package-boundary planning

User requested feature subpackages across layers, low duplication, and reduced
usage. [Package map](package-boundaries.md) owns the decision. Three MCP-created
implementation subitems: [CLI](../../closed/2026/09/08/e7a5e294/plan.md) (closed),
[work boundaries](../../closed/2026/09/10/742aa60b/plan.md), [runtime adapters](../../closed/2026/09/09/b61770a6/plan.md).
This parent coordinates; do not duplicate their implementation scope here.
New feature/runtime wiring pauses behind the first structural slices. No code
moved, no tests rerun for documentation-only work, no model/reset changes.

## 2026-09-11 work-normalize legacy slug false positives

Investigation: current active and parked work records are already folder-per-item
with `plan.md` and item-owned sidecar artifacts inside each item folder. The
folder migration dry-run reports no planned changes. Historical closed archive
files before the current folder convention remain flat by policy and are
tolerated by validation.

`work normalize --dry-run --validate` still flagged linked parked and recently
closed legacy slug IDs such as `release-benchmark-backfill` and
`multi-repo-context-d9b2f3a1` as A3 malformed IDs, even though their plan links
resolve and `work validate` accepts them. The rule is now scoped to active rows:
active IDs remain canonicalized, while linked parked/recently closed legacy
slugs are preserved.

Verification: focused `internal/domain/work` tests pass. `work normalize
--dry-run --validate` no longer reports the 17 bogus A3 fixes; only the three
manual-review B7 closed verification blockers remain. `work normalize --dry-run
--migrate-folders --validate` reports clean with 0 planned folder changes.

## 2026-09-11 source-layout launcher alignment

Investigation: `scripts/source-layout/run.sh` was intentionally thin, but it
anchored `--root` by relative path and nearby docs still linked the closed
`742aa60b` work-domain plan as if it were active. The generated preview's
remaining-work note also still described work-domain semantic extraction as
pending under `742aa60b`, which conflicts with the closed work-domain separation
record and the new folder-normalize evidence above.

Plan: keep the launcher logic thin, but resolve the Git repository root
explicitly and refuse to run if the script is not under the expected checkout.
Update source-layout docs and mapping revision so stale review artifacts remain
rejected after the remaining-work note changes.

Verification: `bash -n scripts/source-layout/run.sh`, `go -C
scripts/source-layout test ./...`, and `go -C scripts/source-layout vet ./...`
pass. `bash scripts/source-layout/run.sh --preview --check` passes with 273
retained files, 0 moves, and 0 content updates against the current source tree.
`bash scripts/source-layout/run.sh --preview --plan review/plan.json` rejects the
checked-in review plan as an older mapping, as intended.

**Closed:** 2026-09-11 — work completed in feature/v0.0.24-work-folder-normalization

## Outcome

v0.0.24 CLI decomposition and architecture audit continuation completed. All
command handlers migrated to strict typed parsers. Non-finite hybrid-ratio
validation fixed with regression coverage. Domain boundary moves completed
(model/usage adapters, work domain separation). Source-layout preview script
built and verified. Work-normalize false-positive rule scoped to active rows.
Source-layout launcher alignment completed. Full Go tests, vet, and HAWP checks
pass.

## Close Checklist

- [x] All focused checks pass (recorded in Verification section above).
- [x] HAWP check passes.
- [x] Plan status: done.
- [x] Plan moved to `.hawp/work/closed/2026/09/11/47c793d6/plan.md`.
- [x] BACKLOG.md updated.
