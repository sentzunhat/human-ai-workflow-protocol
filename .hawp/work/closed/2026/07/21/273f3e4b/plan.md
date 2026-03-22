# hawp update also syncs .hawp/kit/ and the installed provider overlay (unified version)

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `273f3e4b-4978-4fc1-ae8c-2f8ecc99402a`
**Type:** feature
**Reported:** 2026-07-21
**Risk Level:** medium

---

### Input (what was reported)

> Will `hawp update` also update the binary and the latest kit changes
> too, with the providers? Instead of just the binary, update
> everything.
>
> Decision (via clarifying question): yes, one unified version — each
> CLI release tag also captures the kit/provider state; `hawp update`
> fetches and applies both, atomically.

---

### Context

Today `hawp update` (`cdcf9f78`) only replaces the compiled binary.
Kit/provider content is updated via a completely separate,
provider-specific bash-script mechanism in `distribution/generated/`.
This gives the HAWP protocol (kit + provider overlays) a real version
number for the first time, tied 1:1 to CLI release tags, and lets
`hawp update` apply both together.

---

### Analysis

**Root cause (or most likely cause):**
_Two update mechanisms for one conceptual "install" (binary +
protocol content) means users have to remember and run both; unifying
under one release/version removes that split._

**Directly verified:**
_`core/providers/manifest.yaml` already fully specifies, per provider,
the source→destination file mapping AND explicit `install`/`update`
semantics (`seed-if-missing` vs `refresh`) — e.g. `.claude/rules/` is
always refreshed on update, but `CLAUDE.md` (from `CLAUDE.md.seed`) is
`update: skip`, so a user's customized CLAUDE.md is never clobbered.
This is exactly the data needed to drive a generic sync — no new
mapping format needs inventing. Confirmed the materialized source files
this manifest points to actually exist as expected
(`core/providers/.claude/rules/hawp-*.md`,
`core/providers/.claude/CLAUDE.md.seed`)._

**Inferred (not yet proven):**
_Auto-detecting which provider(s) a downstream repo has installed can
be done from the manifest's own destination markers (e.g. presence of
`.claude/rules/` or `.cursor/rules/`) without inventing a separate
stored marker file — untested until a real downstream-repo simulation
is run._

**Scope — what else is affected:**
_Release workflow (new bundle asset alongside the 6 binaries), a new
`internal/domain/kitsync` package (manifest parsing — first YAML
dependency — + provider detection + apply), `hawp update` (extended
flow), `librarian/go/README.md`, `librarian/go/CHANGELOG.md`._

---

### Recommended Fix

- Add `gopkg.in/yaml.v3` (pure Go, no cgo) to parse
  `core/providers/manifest.yaml` as-is — no new format.
- Release workflow: package `.hawp/kit/`, `core/providers/manifest.yaml`,
  and all `core/providers/.<name>/` source directories into one
  `hawp-kit-vX.Y.Z.tar.gz` release asset alongside the binaries.
- `hawp update`: after the binary swap, if run inside a HAWP repo
  (skip gracefully otherwise — a globally-installed binary with no
  repo context shouldn't error), download+verify+extract the bundle,
  refresh `.hawp/kit/` wholesale, detect installed provider(s) via the
  manifest's destination markers, and apply only the `update: refresh`
  entries for detected providers (never touch `update: skip` entries
  like a customized `CLAUDE.md`/`AGENTS.md`).
- Verify with a real, simulated downstream repo (not this source repo,
  which has different semantics as the HAWP source itself): seed a
  scratch repo with stale kit content + Claude-provider markers, run
  the sync against a real release bundle, confirm kit refreshed,
  Claude rules refreshed, `CLAUDE.md` untouched, and an undetected
  provider (e.g. Cursor) left completely alone.

**What to verify after:**

- [x] Release bundle asset is produced and downloadable from a real
      release
      (Evidence: real dispatch run 29886028798 cut `librarian-go-v0.0.6`
      — `gh release view` confirms `hawp-kit-bundle.tar.gz` present
      alongside all 6 platform binaries)
- [x] `hawp update` in a repo with no `.hawp/` skips kit/provider sync
      gracefully (binary-only, no error)
      (Evidence: real run in a scratch dir with no `.hawp/` — binary
      updated `v0.0.5` → `v0.0.6`, exit 0, no kit/provider files created)
- [x] Real simulated-downstream-repo test: kit refreshed, detected
      provider's `refresh` entries updated, `skip` entries (customized
      files) left untouched, non-installed providers untouched
      (Evidence: real run — seeded stale `.hawp/kit/start-here.md`,
      stale `.claude/rules/hawp-core.md`, and a "customized" `CLAUDE.md`
      in a scratch dir; after `hawp update`: kit refreshed to 106 real
      files including the correct current `start-here.md` content,
      `.claude/rules/hawp-core.md` refreshed to the real, current
      generated content, `CLAUDE.md` byte-for-byte unchanged, and no
      `.cursor`/`.codex`/`.github` directories were ever created since
      only Claude was detected)
- [x] `make dist`/`make check` still pass after adding the YAML dependency
      (Evidence: `make check` — vet + full test suite + build — passes;
      `make dist` produced all 6 real binaries in the actual CI run)

---

## Outcome (filled at close)

Closed 2026-07-21. `hawp update` now refreshes `.hawp/kit/` and the
installed provider overlay from the same release the binary updates
to — the CLI and the HAWP protocol content share one version for the
first time, replacing two separate update paths (binary via `hawp
update`, kit/providers via manual `distribution/generated/` scripts).

The sync is driven entirely by the existing
`core/providers/manifest.yaml` (first use of `gopkg.in/yaml.v3`, pure
Go, no cgo) — no new mapping format was invented. Provider
auto-detection uses the manifest's own file-pattern markers (works for
Claude, Cursor, Continue; Codex and GitHub have no such marker and need
an explicit `--provider <name>`). `update: skip` entries (a customized
`CLAUDE.md`/`AGENTS.md`) are never touched — verified for real, not
assumed. Running outside any HAWP repo skips the sync step gracefully,
so a globally installed binary still updates cleanly. `--skip-kit`
restores the old binary-only behavior.

**Real, not mocked, end-to-end proof**: cut an actual test release
(`librarian-go-v0.0.6`) via the updated workflow, confirmed the new
`hawp-kit-bundle.tar.gz` asset published correctly, then ran the real
compiled binary in two scratch scenarios — a simulated downstream repo
with stale kit content and Claude markers (kit and Claude rules
refreshed to genuine current content, customized `CLAUDE.md` untouched,
no other providers touched), and a directory with no `.hawp/` at all
(binary-only update, clean exit, no errors).

## Verification (filled at close)

- Evidence: `gh release view librarian-go-v0.0.6` — 6 binaries +
  `hawp-kit-bundle.tar.gz`, not draft, prerelease (0.0.x convention).
- Evidence: real scratch-repo update run — `Kit refreshed: 106 file(s).`
  / `Provider claude refreshed: 4 file(s).`; refreshed file contents
  verified byte-for-byte against real repo content; `CLAUDE.md`
  verified unchanged; no untouched-provider directories created.
- Evidence: real no-`.hawp/`-directory run — binary updated, exit 0, no
  error, no kit/provider files created.
- Evidence: `go test ./internal/domain/kitsync/...` — 7 tests covering
  manifest parsing against the real schema, refresh-on-update logic,
  detection (including the false-positive-avoidance case for generic
  `.github/` folders), kit sync, and provider apply with skip semantics.
- Evidence: `make check` (vet + tests + build) passes.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Plan file saved under `closed/2026/07/21/273f3e4b-4978-4fc1-ae8c-2f8ecc99402a.md`
- [x] BACKLOG.md updated
