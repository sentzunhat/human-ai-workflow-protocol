# Historical migration progress

Superseded progress notes, preserved in full on 2026-09-07. Use [the current plan](plan.md) for status and next actions; older open/closed claims below are historical.

# v0.0.24 single-executable installation and configuration migration

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `0cb0f9b0-1df6-4250-8274-1c3bfffbbf38`
**Type:** improvement
**Reported:** 2026-09-06

---

## Input (scope summary)

> Coordinate single-executable installers, updates, guides, and configuration migration. Validate completion before proceeding to remaining CLI argument validation, UUID correctness, storage boundaries, preservation-first unknown work handling, and broader indexing. Prefer available HAWP MCP tools and maintain existing work context. Do not merge PRs or modify downstream repositories.

## Intake Summary

Complete a tested single-executable migration before pulling the subsequent
architecture and work-normalization queue forward. Parent release:
[e5fca9c7](../../../../../active/e5fca9c7/plan.md); architecture context:
[47c793d6](../../../../../active/47c793d6/plan.md).

## Current Context

Investigation began under the architecture parent; this item owns the distinct
installer/configuration compatibility gate. New MCP configs already launch the
native executable with explicit repository selection. Install/update scripts
still install `hawp-bin` plus wrappers; retiring them is not yet implemented.

## Initial Analysis

**Directly verified:**

- Connected HAWP MCP validation targets this repository and passes all checks.
- Canonical install/update scripts still write `.hawp/bin/hawp-bin` and copy
  wrappers from `core/.hawp/bin/`.
- JSON configuration migration replaced the complete HAWP entry, losing custom
  settings. Null top-level JSON could panic; invalid nested shapes could be
  silently replaced. The first implementation step fixes these cases.
- Codex TOML migration still replaces a whole HAWP table. Preservation of custom
  values and nonstandard table layouts requires a separate tested change.

**Inferred (not yet proven):**

- External/user-global configurations may reference old executable names;
  repository-local config inspection cannot establish their absence.

**Likely scope:**

- Source installer/update templates, generated guides, native installation
  paths, provider config migration, local/CI install contracts, and tests.

## Risk + Review Gate

**Risk:** medium (installation and configuration compatibility).
**Gate:** user explicitly authorized this migration and local commits; no merge,
publication, or changes to downstream checkouts. Can implement now: yes.

## Backlog + Plan Link

**Status now:** in-progress
**Plan file:** work/active/0cb0f9b0/plan.md

## Next Step

Installer transition preflight: canonical scripts currently accept a download
when checksums are missing/unusable and share a predictable checksum file in
the system temporary directory. Before executable renaming, require one valid
matching SHA256 entry, stage downloads privately under the destination folder,
and preserve the installed binary on failure. Test extracted install/update
functions in temporary repositories with network mocked; no live installation.
This prerequisite does not itself change the wrapper/native filename contract.

Lossless TOML continuation: replace the line-based table replacement with a
parser-backed, narrowly scoped infrastructure adapter. Pin go-toml v2.4.3
(MIT, no module dependencies); isolate its unstable AST API and cover it with
byte-preservation tests. This is a format parser, not a CLI framework. Preserve
custom policy/environment fields and comments; reject malformed, remote, or
unsupported layouts without writes. Use the same merge for configure preflight
and init's Codex writer. Leave local provider configs, including the untracked
Copilot file, untouched; verify on temporary fixtures.

- [x] Investigation and scope recorded; separate compatibility lane created via MCP.
- [x] Extract pure JSON merge; preserve custom settings and reject incompatible
      shapes before writes. Focused MCP/CLI tests pass.
- [x] Use a same-directory temporary file and rename for local/CI installation;
      avoid in-place executable overwrite and remove temporary files on failure.
- [x] Preserve Codex table customizations; test explicit/quoted tables, nested
      settings, multiline values, byte preservation, and safe refusal.
- [x] Add config-only setup without model downloads or kit mutation. Custom
      explicit Codex tables now migrate; unsupported layouts are preserved by refusal.
- [x] Coordinate canonical native executable path across source installers,
      update paths, Makefile, CI, and provider configuration generation.
- [x] Regenerate guides; cover fresh installs and legacy upgrades, interrupted
      downloads, checksum failures, existing config preservation, and repeat runs.
- [x] Prove migration on temporary downstream copies, plus standard platform
      builds. Record runtime/provider checks separately from compilation.
- [x] Retire wrappers only after compatibility verification; close this gate.

## Subsequent Queue

Do not report these as implemented:

1. Remaining CLI argument validation, followed by persisted UUID/status
   correctness and application-owned corpus/storage boundaries (`47c793d6`).
2. Preservation-first unmatched-record handling (`e5fca9c7`): proposed destination
   `.hawp/work/unknown/`, dry-run manifest with original paths, collision refusal,
   link repair and sidecar preservation, repeat-run and rollback tests. Never
   blindly overwrite records or infer a UUID for an ambiguous match.
3. Index coverage: audit existing kit/work/docs behavior before widening roots.
   Require explicit include rules and exclude secrets, dependencies, generated
   builds, and symlink escapes. Unknown records must retain unknown identity and
   lifecycle rather than being mistaken for active work.
4. Existing-item updates: enrich the same plan and backlog without losing context.
   Current MCP offers creation and validation, not a dedicated update operation.

## Outcome

Latest TOML continuation supersedes the earlier blanket custom-config refusal:
both configure preflight and init use a shared preserving merge. Only HAWP launch
values change; missing defaults are inserted, existing custom fields remain.
An infrastructure adapter confines the pinned go-toml unstable AST usage. Full
semantic parsing occurs before/after edits, with comparison to expected changes.
No personal/provider config was migrated in place; tests use temporary fixtures.
Unsupported inline/dotted-only server definitions, remote transports, malformed
TOML, and comments inside a changed launch value are refused without rewriting.
Installer transition and single-executable retirement remain open.

TOML verification: full Go tests, vet, six standard target builds, local native
rebuild/check, connected MCP validation, distribution validation, kit parity,
and diff checks passed. Regression fixtures cover custom disabled/allowlist/
timeout/env settings, quoted tables, multiline values, CRLF, trailing comments,
argument updates, repeated runs, semantic equivalence, and refusal without writes.
The adapter reparses its result and checks that no unrelated values changed.
Byte preservation excludes launch values intentionally replaced and newly added
defaults; comments inside changed values trigger refusal rather than loss.
Dependency scope is one pinned MIT TOML library, with no module dependencies;
its unstable AST API is isolated and must be retested before version upgrades.
No live provider configuration was edited, and no new provider runtime proof is
claimed. New untracked install-agent guides from parallel work are left unstaged.

2026-09-06 documentation continuation: added canonical and root-mirrored
provider setup guides, shared digital-worker instructions, and lessons under
`.hawp/kit/usage/mcp/`. Linked them from start-here and search guidance.
Codex/Claude init support and GitHub's manual-only MCP setup are distinguished;
init side effects, Codex table replacement, platform paths, and unpublished
branch-version requirements are explicit. Official provider docs were checked.
No provider configurations were changed and no init downloads were run.
Documentation checks passed: focused MCP/CLI tests, connected HAWP MCP
validation, JSON example parsing, complete root/core kit parity, distribution
validation, and diff whitespace checks. Live Claude/Copilot connections remain
unverified. Lessons are repository documentation, not personal memory updates.

Partial: JSON config migration safety and MCP-first operating guidance improved.
Single-executable installation is still open; no wrapper has been removed.

## Verification

Installer download prerequisite: both canonical scripts now require exactly one
valid matching SHA256 entry, restrict asset downloads and redirects to HTTPS,
and stage privately beside the destination with scoped cleanup traps. The 20
network-mocked install/update cases cover success, partial download, missing or
invalid checksums, duplicate entries, digest mismatch, missing/failing hash
utilities, and termination. Failures preserve the prior binary and launcher and
leave no staging directory. Generated all 20 provider guides from these sources.
These fixtures exercise Linux asset selection, not live release downloads or
every platform runtime. Wrapper copying is skipped in the fixture; whole-kit
installation remains nontransactional. The single-executable filename migration
is still pending, and no live/downstream installation was performed.
Verification passed: full Go suite, vet, distribution validation, root/core kit
parity, diff whitespace checks, and connected HAWP MCP validation (three checks;
137 Markdown files; zero work issues or warnings). No native rebuild is needed
for this template-and-test-only change.

MCP-only continuation: added `mcp configure --provider NAME` with explicit
selection, optional repository root, and strict flag validation. Pure config
preflight runs before selected writes. Native file/backlog prerequisites are
checked; existing customized Codex TOML is preserved verbatim with a manual-merge
error rather than invoking the destructive table replacement. JSON configs reuse
the tested merge. Runtime compatibility is still a separate check, and write-time
filesystem failures are not transactional. Existing `init` behavior is unchanged.

Latest provider evidence supersedes earlier Copilot-unavailable notes: user
reported `mcp_hawp2_hawp_work_validate` with all three checks passing. Codex is
directly verified, Claude/Copilot are user-reported. Cursor/Continue live tests
are explicitly deferred by the user. Untracked `.vscode/mcp.json` is preserved
and not staged. No provider config in this checkout was modified by this step.
Verification passed: full Go suite, added focused CLI/MCP tests, vet, six standard
target builds, local binary rebuild, native advice-only configure smoke, connected
MCP validation, distribution validation, root/core kit parity, and diff checks.
Temporary fixtures prove fresh/repeat configuration, no home/kit provisioning,
explicit-root routing, and custom Codex refusal before earlier selected writes.
Cross-platform compilation does not establish live execution on every platform.

2026-09-06 user-supplied consumer evidence: Claude Code reports MCP validation
PASS for this repository (kit/work/links, 136 Markdown files). This session did
not independently observe that call. Copilot's transcript instead runs a shell
`work validate` command and reports five work checks passing; MCP remains
unverified. Updated both kit copies with the distinction and a no-shell-fallback
probe. Its ignored CLI `--repo-root` argument prompted a focused parser fix in
the architecture item. No downstream work records or client configs were changed.

Provider extension continuation: retained five client names without a new
framework. Configuration selection now validates all names before config writes
and expands/deduplicates mixed `all` selections once. Tests cover invalid names
after a valid name with no config files created. This is a config-writer guarantee,
not an init-wide transaction. Saved a provider matrix and connection/extension
checklist in both kit copies. Live HAWP MCP validation passed here; other clients
remain unverified. Installer retirement remains open.

Distribution script artifact continuation: generated provider guides now have
sibling `.sh` outputs beside each `.md` guide. The Markdown remains the
reviewable agent instruction set, while the shell artifact is the same composed
script without fences for review-first download and explicit execution. The
distribution validator now treats the shell scripts as generated outputs, so CI
can catch drift. This does not pipe remote scripts into a shell and does not yet
remove the wrapper/native split.

- Full Go suite, vet, six standard target builds, distribution build/validation,
  and root/core operating-guide comparison pass. Generated guides remain current
  and unchanged because installer templates have not yet migrated.
- Initial in-place local install produced exit 137 when executing the installed
  file; the build output ran and on-disk signature verification passed. Installing
  through a new file plus rename restored successful native execution and HAWP
  check. Signature-cache involvement is an inference, not a proven diagnosis.
- HAWP MCP creation/validation were used against this repository. Work validation
  reports zero issues/warnings. No downstream consumer configuration was edited.

Canonical path continuation (2026-09-06):

- `make install` now writes the native binary to `.hawp/bin/hawp` (changed from
  `hawp-bin`). This aligns the Makefile and CI with what the distribution scripts,
  `hawpBinaryPath` config generation, and binary download tests already expected.
  The shell launcher at `.hawp/bin/hawp` is superseded; `.hawp/bin/hawp-bin`
  remains as a preserved legacy binary for existing consumer configs.
- `hawp distribution sync` regenerated all 40 generated provider guides. The
  install and update scripts now write to `.hawp/bin/hawp${_ext}` consistently.
  Distribution validation passed: generated outputs are current.
- Full Go suite and `go vet ./...` passed with no changes to tests.
- `hawp check --no-update-check` passed: kit validate (0 issues), work validate
  (0 issues, 0 warnings), links check (137 files). `git diff --check` clean.
- `.hawp/bin/hawp` confirmed as Mach-O 64-bit arm64 native binary reporting 0.0.24.
- Mochila-archive-viewer temp copy proof (2026-09-06): binary upgrade from
  0.0.23 → 0.0.24 via temp+rename succeeded. `hawp mcp configure --provider
  claude --provider codex` wrote both configs with absolute canonical paths.
  `hawp work validate` ran against the temp copy and returned pre-existing
  mochila issues (2 closed plans missing sections, 4 sub-task naming warnings)
  — not migration-caused. `hawp-bin` (0.0.23) and `hawp-mcp` were preserved
  unchanged. Temp copy was deleted; no changes committed to mochila.
- Wrapper retirement (2026-09-06): `hawp-mcp` removed from `.hawp/bin/` and
  `core/.hawp/bin/` via `git rm`. Superseded by `hawp mcp --repo-root`.
  `hawp-bin` left as preserved legacy binary for existing consumer configs.
  This item is closed.

## 2026-09-07 closure reconciliation investigation

Canonical downloaders and Makefile use native hawp. Existing download tests
cover legacy hawp/hawp-bin files only; fresh installs, repeat success, and
preservation of hawp-mcp lack direct fixtures. `core/.hawp/bin/hawp` remains a
legacy shell launcher, although canonical installer functions no longer copy it.
The prior closed claim therefore does not establish every retirement gate.

Plan: extend the existing mocked download matrix for fresh/legacy layouts,
repeat success, and old wrapper preservation; retain runtime limitations.
Compact this plan with the full prior progress retained in sibling history.
Correct current kit guidance and architecture-state notes, preserving history.
Risk low: tests and documentation only; no installed binary or wrapper removal.
Can implement now: yes, existing migration UUID retained, user requested cleanup.
Root proof: pwd and top-level `<repo-root-abs>`, prefix empty; prior dirty work
from the architecture/reshape continuations is preserved.
