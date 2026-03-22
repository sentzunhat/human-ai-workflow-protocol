## Task: Plan `core/` multi-provider layout (GitHub, Cursor, Continue) + distribution

**Backlog ID:** TASK-072
**Type:** task
**Reported:** 2026-06-04 (initial); **revised:** 2026-06-04 (rev.4 — `distribution/sources/providers/` layout)
**Risk Level:** medium
**Status:** done
**Closed:** 2026-06-04

---

### Stakeholder decisions (approved 2026-06-04)

| # | Topic | TASK-072 assumption | **Approved decision** |
|---|-------|---------------------|------------------------|
| 1 | Default install | kit + GitHub only | **No default** — per-provider generated guides under `distribution/generated/{github,cursor,continue}/` by branch |
| 2 | `copilot-instructions.md` | seed-only | **Seed on install if missing; refresh on update** |
| 3 | `AGENTS.md` | seed, no overwrite | **Always refresh** (HAWP-managed) — Phase 2+ |
| 4 | Cursor rules | `hawp-*.mdc` only | **All files** in `providers/.cursor/rules/` |
| 5 | Install guides | one guide + sections | **Split guides per provider** |
| 6 | Source repo `.github/` vs install pack | keep both | **Confirmed** — install pack is `core/providers/.github/` only |

Phase 1 implementation tracked in TASK-074. Split-guide composition is Phase 2.

---

### Stakeholder direction (latest — canonical layout confirmed)

**Keep `core/` as the install source root.** Only two top-level areas under `core/`:

```text
core/
├── .hawp/                    # kit + LICENSE (agent-neutral) — unchanged
└── providers/                # all agent integrations + shared sources
    ├── manifest.yaml         # install routing: source subtree → target repo paths
    ├── README.md             # provider index, boundaries, env flags
    ├── shared/               # canonical prompts / instructions / scripts (.md)
    │   ├── instructions/
    │   ├── prompts/
    │   └── scripts/          # optional bash fragments for install materialization
    ├── .github/              # GitHub / Copilot provider (move from `core/.github/`)
    │   ├── copilot-instructions.md
    │   ├── instructions/
    │   └── prompts/
    ├── .cursor/              # Cursor provider (NEW)
    │   ├── README.md
    │   ├── AGENTS.md.seed
    │   └── rules/hawp-*.mdc
    └── .continue/            # Continue.dev provider (NEW — content in TASK-073)
        ├── README.md
        ├── config.yaml.seed
        └── rules/hawp-*.md
```

**Install mapping (source → target repo):**

| Source under `core/` | Installs to (downstream) |
|----------------------|---------------------------|
| `core/.hawp/` | `.hawp/LICENSE`, `.hawp/kit/**` |
| `core/providers/.github/` | `.github/instructions/`, `.github/prompts/`, `.github/copilot-instructions.md` |
| `core/providers/.cursor/` | `.cursor/rules/hawp-*.mdc`, `AGENTS.md` (seed) |
| `core/providers/.continue/` | `.continue/rules/hawp-*.md`, optional `.continue/config.yaml` (seed) |

Provider folder names **use dot-prefix** (`.github`, `.cursor`, `.continue`) to mirror downstream tool conventions while living under `core/providers/`.

**Migration (TASK-074):** Move `core/.github/**` to `core/providers/.github/**`. Point install and update scripts at `cp "$SRC/providers/.github/..."` instead of `cp "$SRC/.github/..."`. Keep a one-release fallback: if `providers/.github` is missing, read legacy `core/.github/` and print a warning.

**Withdrawn layouts:** `core/agents/` (rev.1); dot-folder providers at `core/.github` + `core/.cursor` siblings (rev.2).

**Distribution must change** — mirror provider categories under `distribution/sources/` (same pattern as existing `install/`, `update/`, `shared/`):

```text
distribution/sources/
├── install/                    # existing: main.md, dev.md, script.md, preamble.md
├── update/                     # existing: main.md, dev.md, script.md, preamble.md
├── shared/                     # existing: install.md, update.md, safety.md, repo-boundaries.md
│   └── providers-overview.md   # NEW: cross-provider index (optional rename from providers.md)
└── providers/                # NEW: per-provider user-facing guide fragments
    ├── README.md               # how provider docs compose into generated guides
    ├── shared/                 # distribution-only shared fragments (not core pack files)
    │   ├── boundaries.md       # provider boundary summary for all guides
    │   └── env-and-flags.md    # HAWP_PROVIDERS, defaults, verification snippets
    ├── github/
    │   ├── install.md          # what GitHub/Copilot install adds; verify steps
    │   └── update.md
    ├── cursor/
    │   ├── install.md
    │   └── update.md
    └── continue/
        ├── install.md          # stub until TASK-073
        └── update.md
```

**Two-layer model (do not conflate):**

| Layer | Path | Purpose |
|-------|------|---------|
| **Install pack** (copied to user repos) | `core/providers/.github/`, `.cursor/`, … | Actual instructions, prompts, rules |
| **Distribution docs** (composed into guides) | `distribution/sources/providers/github/`, … | Human-readable install/update sections in `distribution/generated/*.md` |

Generated guides (`install-main.md`, `update-dev.md`, etc.) gain provider sections by extending `librarian/scripts/distribution/shared/composition.ts` `sectionFiles` to include, after shared install/update body:

```text
sources/providers/shared/boundaries.md
sources/providers/github/install.md    # or update.md per variant
sources/providers/cursor/install.md    # when HAWP_PROVIDERS includes cursor
sources/providers/continue/install.md  # stub / when enabled
```

Install and update **bash scripts** remain in `distribution/sources/install/script.md` and `update/script.md` (one script each, manifest-driven). Optional provider-specific bash notes may later live under `distribution/sources/providers/*/`.

**TASK-073:** Implement `core/providers/.continue/` + `distribution/sources/providers/continue/*` after TASK-074.

---

### Input (full conversation capture)

#### Message 1 (2026-06-04)

User asked to:

- Review the HAWP repository and **plan a new work item** for getting the library to work in **Cursor** automatically.
- Research how Cursor loads project context (rules, `AGENTS.md`, hooks, skills).
- Defer a **second work item for Continue** (another extension/agent library) for later.
- Outcome: planning work item with rich context, not necessarily implementation yet.

#### Assistant findings (message 1) — retained as baseline evidence

- HAWP core is **portable markdown** + optional kit; install today copies `core/.hawp/kit/**` and **only** `core/.github/**` for agents.
- Cursor does **not** read `.github/instructions/*.instructions.md` or `.github/prompts/*.prompt.md` natively.
- Cursor mechanisms: `AGENTS.md` (auto), `.cursor/rules/*.mdc` (`alwaysApply` / `globs` / description-triggered), hooks/skills optional.
- Proposed interim architecture used `core/agents/{shared,github,cursor}` — **superseded** by stakeholder direction below.
- Noted **drift**: docs say `copilot-instructions.md` is seed-only; install/update scripts **overwrite** it.
- Created **TASK-072** plan + backlog row; placeholder **TASK-073** for Continue.

#### Message 2 (2026-06-04)

User asked to:

- **Review latest changes** (3–4 commits on 2026-06-04).
- Create/update work item for **agent type providers** with **two providers now** (GitHub + Cursor), room for more.
- Structure: **`core/` stays `core`** with `.github`, `.hawp`, **`.cursor`**, **`.continue`** — or `core/providers/` containing each provider.
- **Shared** prompts/scripts as `.md` sources materialized into provider-specific folders by install.
- **Distribution changes** to include providers; install guides help users per provider (GitHub, Cursor, Continue, others later).
- Each install script **bases paths from `core/`** provider trees.
- **Grab all conversation context**; surface **questions and assumptions**.

#### Commits reviewed (message 2)

| Commit | Date | Summary |
|--------|------|---------|
| `8060529` | 2026-06-04 | Backlog; closed TASK-069/071; standards audit evidence |
| `22038c6` | 2026-06-04 | `core/.hawp/kit/standards/docs/` + `hawp-install-update-safety.md` |
| `dd97c0c` | 2026-06-04 | Mirror `core/.hawp/kit/standards/public/**` (32 files) |
| `89bf269` | 2026-05-19 | Distribution pipeline, librarian, `.hawp/bin/hawp`, validators (foundation) |

None of these commits added Cursor/Continue under `core/`; agent integration remains GitHub-only in install scripts.

#### Message 3 (2026-06-04) — layout confirmation

User confirmed structure:

```text
core/
├── .hawp/
└── providers/
    ├── .github/
    ├── .cursor/
    ├── .continue/
    └── shared/
```

All agent provider sources live under `core/providers/`; kit stays at `core/.hawp/` only.

#### Message 4 (2026-06-04) — distribution provider sources

User confirmed **`distribution/sources/`** should also have **provider-category folders** (github, cursor, continue, etc.) — parallel to how install/update/shared are organized today — so generated install/update guides can include **per-provider documentation sections**.

---

### Context

#### What HAWP is

- **Protocol:** five-field `Shape` (`input`, `context`, `mission`, optional `checkpoint`, `constraints`, `output`) — see `core/.hawp/kit/spec.md`.
- **Not:** runtime, validator, orchestrator, memory system.
- **Kit:** `.hawp/kit/**` — templates, standards, references, `start-here.md` (agent-neutral).
- **Work tracking:** `.hawp/work/**` — project-owned; never installed from source repo.

#### Current `core/` layout (verified)

```text
core/
├── .hawp/
│   ├── LICENSE
│   └── kit/**              # installed → target .hawp/kit/** (full refresh)
├── .hawp/work/**           # source-repo only; NEVER installed downstream
└── .github/
    ├── copilot-instructions.md
    ├── instructions/*.instructions.md   (5 files)
    └── prompts/*.prompt.md            (11 files)
```

**Missing today:** `core/providers/` entirely; GitHub files still at legacy `core/.github/` (to be moved).

#### Current distribution layout (verified — 12 source files)

```text
distribution/
├── sources/
│   ├── install/     # main.md, dev.md, script.md, preamble.md
│   ├── update/      # main.md, dev.md, script.md, preamble.md
│   └── shared/      # install.md, update.md, repo-boundaries.md, safety.md
└── generated/       # composed by librarian composition.ts
    ├── install-main.md, install-dev.md
    └── update-main.md, update-dev.md   # user has update-dev.md open (683 lines)
```

**Missing today:** `distribution/sources/providers/` (no github/cursor/continue doc fragments).

**Composition today** (`librarian/scripts/distribution/shared/composition.ts`): each variant stitches `preamble` → `shared/safety` → `shared/repo-boundaries` → `shared/install|update` → `install|update/main|dev` → embedded bash from `install|update/script.md`.

Install script (`distribution/sources/install/script.md`) today:

- Copies `core/.hawp/kit/**` (selective `cp` list + standards tree).
- Copies `core/.github/...` today (legacy path) → target `.github/...` (**refresh** / overwrite on `copilot-instructions.md`).
- **After migration:** must copy from `core/providers/.github/...` instead.

`HAWP_LOCAL_CORE` allows local `core/` testing without archive download.

#### GitHub provider inventory (today: `core/.github/` — 16 files; target: `core/providers/.github/`)

**Instructions:**

- `hawp-intake.instructions.md` — `applyTo: ".hawp/**,**/.hawp/**"`
- `hawp-backlog-alignment.instructions.md`
- `hawp-docs-alignment.instructions.md`
- `commit-style.instructions.md`
- `intake.instructions.md`

**Prompts:**

- `hawp-status-report.prompt.md`
- `hawp-backlog-alignment.prompt.md`
- `hawp-intent-first-handoff.prompt.md`
- `hawp-change-review-and-reference-sync.prompt.md`
- `hawp-commit-one-big.prompt.md`, `hawp-commit-many-small.prompt.md`
- `hawp-docs-alignment-deterministic.prompt.md`, `hawp-docs-alignment-simplicity.prompt.md`
- `hawp-conservative-docs-drift-cleanup.prompt.md`
- `intake.prompt.md`

**Entry:** `copilot-instructions.md` → `.hawp/kit/start-here.md`, status reports, backlog compaction, references to instruction/prompt paths.

#### Cursor provider (research — target artifacts under `core/providers/.cursor/`)

| Target repo path | Source (proposed) | Load behavior |
|------------------|-------------------|---------------|
| `AGENTS.md` | `core/providers/.cursor/AGENTS.md.seed` → seed on install | Cursor Agent auto-loads |
| `.cursor/rules/hawp-*.mdc` | `core/providers/.cursor/rules/` → refresh on install | globs / alwaysApply / description |
| Optional hooks | `core/providers/.cursor/hooks/` | Phase 3+; not required for v1 |

Reference: [Cursor Rules](https://cursor.com/docs/context/rules).

**Copilot → Cursor mapping (initial port):**

| Concern | GitHub (`core/providers/.github/`) | Cursor (`core/providers/.cursor/`) |
|---------|-----------------------------------|-------------------------------------|
| Global entry | `copilot-instructions.md` | `AGENTS.md` (from seed) |
| Scoped intake | `hawp-intake.instructions.md` | `rules/hawp-intake.mdc` (`globs: "**/.hawp/**"`) |
| Backlog alignment | `hawp-backlog-alignment.instructions.md` | `rules/hawp-backlog-alignment.mdc` |
| Docs alignment | `hawp-docs-alignment.instructions.md` | `rules/hawp-docs-alignment.mdc` |
| Status / handoff / review | `*.prompt.md` | thin `.mdc` + `@.hawp/kit/...` pointers |
| Repo intake | `intake.instructions.md`, `intake.prompt.md` | optional rules if project uses intake loop |

#### Continue provider (research — target artifacts under `core/providers/.continue/`)

| Target repo path | Source (proposed) | Notes |
|------------------|-------------------|-------|
| `.continue/rules/*.md` | `core/providers/.continue/rules/` | YAML frontmatter: `alwaysApply`, globs |
| `config.yaml` fragments | `core/providers/.continue/config.yaml.seed` or docs-only | Continue reads workspace + `~/.continue/`; seed carefully |

Reference: [Continue Rules](https://docs.continue.dev/customize/deep-dives/rules), [config.yaml](https://docs.continue.dev/reference).

**TASK-073:** implement Continue slice after Cursor + manifest exist; do not block TASK-072 planning on Continue content.

#### Related backlog / standards context

| ID | Status | Note |
|----|--------|------|
| TASK-070 | inbox | Private workflow standards adaptation |
| TASK-069 | closed | Public docs standards; private lane audit |
| TASK-071 | closed | Full `kit/standards/public/**` mirror |
| TASK-073 | inbox | Continue implementation — child of this architecture |

TASK-069 evidence: avoid confusing HAWP **agent providers** with `standards/private/providers/README.md` (domain “providers” in architecture standards). Use clear docs: “agent provider” = Cursor/GitHub/Continue integration lane.

#### Known drift (must fix during implementation)

| Documented | Script actual |
|------------|---------------|
| `copilot-instructions.md` seed-only | always `cp` overwrite in install + update |
| `repo-boundaries.md` lists only GitHub agent paths | no Cursor/Continue boundaries yet |

`core/.hawp/kit/standards/docs/hawp-install-update-safety.md` should be updated when provider boundaries land.

---

### Analysis

**Root cause:** Agent integration is implicit “GitHub-only” in install scripts and distribution guides. There is no `core/.cursor/` or `core/.continue/`, no provider manifest, and no user-facing “install HAWP for Cursor” path.

**Goal:** Users run **one HAWP install/update** (or provider-aware variants) and get:

1. Kit (always).
2. Selected agent provider overlays materialized from matching `core/.<provider>/` trees.
3. Shared canonical `.md` sources (under `core/providers/shared/`) driving long-term deduplication—not duplicate prose across GitHub/Cursor/Continue.

**Directly verified:**

- `core/` has `.hawp/` + legacy `core/.github/` only; no `core/providers/` yet.
- Generated guides (`distribution/generated/install-main.md`, `update-dev.md`) describe GitHub outcomes only in “What Was Added”.
- Root `.hawp/bin/hawp` delegates to librarian; no provider CLI yet.

**Inferred:**

- Install script should loop **`$SRC/providers/manifest.yaml`** rather than hard-coded `cp "$SRC/.github/..."` blocks.
- **TASK-074 must move** `core/.github/` → `core/providers/.github/` before scripts point at new paths only.
- Distribution `librarian/scripts/distribution/build` may need new fragments: `shared/providers.md`, per-provider sections in `install/main.md` / `update/main.md`.
- Default install remains **kit + github** for backward compatibility; Cursor/Continue documented as same script with env flag or explicit “provider pack” sections.

---

### Target architecture (stakeholder-confirmed)

```text
core/
├── .hawp/                         # kit + LICENSE only
│   ├── LICENSE
│   └── kit/**
└── providers/
    ├── manifest.yaml
    ├── README.md
    ├── shared/
    │   ├── instructions/
    │   ├── prompts/
    │   └── scripts/
    ├── .github/                   # from legacy core/.github/ (move)
    ├── .cursor/                   # new
    └── .continue/                 # new (TASK-073)
```

**Downstream after full install (`HAWP_PROVIDERS=hawp,github,cursor,continue`):**

```text
target-repo/
├── .hawp/kit/**                              # from core/.hawp/
├── .github/instructions/, prompts/           # from core/providers/.github/
├── .github/copilot-instructions.md           # seed
├── .cursor/rules/hawp-*.mdc                  # from core/providers/.cursor/
├── AGENTS.md                                   # seed from providers/.cursor/
├── .continue/rules/hawp-*.md                   # from core/providers/.continue/
└── .hawp/work/**                               # never overwrite
```

#### `core/providers/manifest.yaml` (install routing)

```yaml
version: 1
defaults:
  providers: [hawp, github]

providers:
  hawp:
    source_root: .hawp
    targets:
      - { src: LICENSE, dest: .hawp/LICENSE, mode: refresh }
      - { src: kit, dest: .hawp/kit, mode: refresh }

  github:
    source_root: providers/.github
    targets:
      - { src: instructions/*.instructions.md, dest: .github/instructions/, mode: refresh }
      - { src: prompts/*.prompt.md, dest: .github/prompts/, mode: refresh }
      - { src: copilot-instructions.md, dest: .github/copilot-instructions.md, mode: seed }

  cursor:
    source_root: providers/.cursor
    targets:
      - { src: rules/hawp-*.mdc, dest: .cursor/rules/, mode: refresh }
      - { src: AGENTS.md.seed, dest: AGENTS.md, mode: seed }

  continue:
    source_root: providers/.continue
    enabled: false
    targets:
      - { src: rules/hawp-*.md, dest: .continue/rules/, mode: refresh }
      - { src: config.yaml.seed, dest: .continue/config.yaml, mode: seed }
```

**Environment override (proposed):** `HAWP_PROVIDERS=hawp,github,cursor` for install/update scripts.

**`mode` semantics:**

- `refresh` — HAWP-managed; replace on every install/update
- `seed` — `copy_file_no_clobber` only
- `preserve` — never touch (documented for `.hawp/work/**`, user custom rules)

---

### Distribution changes (specification)

#### Target `distribution/sources/` tree (rev.4)

```text
distribution/sources/
├── install/                          # unchanged role + script.md (manifest-driven)
├── update/
├── shared/
│   ├── install.md, update.md, safety.md, repo-boundaries.md  # extend boundaries for cursor/continue
│   └── providers-overview.md         # short pointer → sources/providers/README.md
└── providers/
    ├── README.md
    ├── shared/
    │   ├── boundaries.md             # all providers: managed vs project-owned paths
    │   └── env-and-flags.md          # HAWP_PROVIDERS, defaults, proof commands
    ├── github/
    │   ├── install.md                # “What Was Added” for Copilot overlay
    │   └── update.md                 # what refresh changes; copilot-instructions seed note
    ├── cursor/
    │   ├── install.md                # AGENTS.md + hawp-*.mdc; verification in Cursor Agent
    │   └── update.md
    └── continue/
        ├── install.md                # stub → TASK-073
        └── update.md
```

#### Composition plan update (`composition.ts`)

Extend each variant’s `sectionFiles` (after `shared/install.md` or `shared/update.md`, before branch `main|dev`):

**Install-main / install-dev example order:**

1. `sources/install/preamble.md`
2. `sources/shared/safety.md`
3. `sources/shared/repo-boundaries.md`
4. `sources/shared/install.md`
5. **`sources/providers/shared/boundaries.md`** _(new)_
6. **`sources/providers/shared/env-and-flags.md`** _(new)_
7. **`sources/providers/github/install.md`** _(new — default provider docs)_
8. **`sources/providers/cursor/install.md`** _(new — labeled “optional / Cursor”)_
9. **`sources/providers/continue/install.md`** _(new — stub)_
10. `sources/install/main.md` or `dev.md`
11. bash block from `install/script.md`

**Update variants:** same pattern with `providers/*/update.md` and `shared/update.md`.

Optional: branch-specific overrides `providers/github/install-dev.md` only if dev differs materially (default: single `install.md` per provider).

#### Generated output expectations

After `npm run distribution:sync`, `distribution/generated/update-dev.md` (and siblings) include:

- Core HAWP install/update flow (unchanged preamble)
- **Provider overview** + env flags
- **GitHub section** — files under `.github/`, seed behavior for `copilot-instructions.md`
- **Cursor section** — files under `.cursor/`, `AGENTS.md`, how to enable via `HAWP_PROVIDERS`
- **Continue section** — “coming soon” or minimal stub until TASK-073

#### Script sources (unchanged location)

| File | Role |
|------|------|
| `distribution/sources/install/script.md` | Single install bash; reads `core/providers/manifest.yaml` at runtime |
| `distribution/sources/update/script.md` | Single update bash |

Provider-specific bash snippets (if ever needed) could live in `distribution/sources/providers/*/script-notes.md` for maintainers only—not duplicated into generated guides unless extracted like main script.

#### Librarian validation (TASK-074)

- Assert required `distribution/sources/providers/{github,cursor,continue}/{install,update}.md` exist (continue may be stub with `status: planned` header).
- Assert `composition.ts` lists all new section files.
- Assert `core/providers/manifest.yaml` and `core/providers/.github/` exist post-migration.
- `distribution:validate` fails if generated guides drift from composition plan.

---

### Options (architecture)

| Option | Description | Status |
|--------|-------------|--------|
| **A (chosen)** | `core/.hawp/` + `core/providers/{shared,.github,.cursor,.continue}` + `manifest.yaml` | **Stakeholder confirmed (message 3)** |
| **B** | Dot-folders at `core/` root (`.github`, `.cursor` siblings of `.hawp`) | Withdrawn (rev.2) |
| **C** | `core/agents/{github,cursor}` without dot-prefix | Withdrawn (rev.1) |

**Recommended fix:** Option A.

---

### Implementation phases (follow-up after plan approval)

| Phase | ID | Deliverable |
|-------|-----|-------------|
| 0 | TASK-072 | Approve this plan; optional ADR |
| 1 | TASK-074 | Create `core/providers/`; **move** `core/.github/` → `core/providers/.github/` |
| 2 | TASK-074 | Add `providers/shared/`, `manifest.yaml`, `providers/.cursor/` scaffold |
| 3 | TASK-074 | Refactor install/update scripts to `$SRC/providers/*`; legacy fallback one release |
| 4 | TASK-074 | Add `distribution/sources/providers/{shared,github,cursor,continue}/`; extend `composition.ts` |
| 4b | TASK-074 | `npm run distribution:sync`; verify `generated/update-dev.md` includes provider sections |
| 5 | TASK-074 | Fix `copilot-instructions.md` seed; update install-update-safety standard |
| 6 | TASK-073 | `providers/.continue/` + guide sections |
| 7 | TASK-074+ | Pilot `providers/shared/` → materialize into `.github` / `.cursor` outputs |

---

### Verification plan (implementation)

| Scenario | Expected |
|----------|----------|
| Default `HAWP_PROVIDERS` unset | Same as today: `.hawp/kit` + `.github/*` only |
| `HAWP_PROVIDERS=hawp,github,cursor` | + `.cursor/rules/hawp-*.mdc`, seeded `AGENTS.md` |
| Re-run update | Refreshed hawp-managed; seeded files preserved if customized |
| Custom `AGENTS.md` / `copilot-instructions.md` | Not overwritten |
| User `.cursor/rules/custom.mdc` | Not deleted |
| Cursor Agent: “status report” | Uses `.hawp/kit/usage/status-report.md` without manual @ |
| `npm run distribution:validate` | Passes with provider manifest checks |

```bash
# proof after install
git status --short .hawp/kit .github .cursor/rules AGENTS.md
find .cursor/rules -maxdepth 1 -name 'hawp-*.mdc' 2>/dev/null | sort
```

---

### Questions for stakeholder (remaining — layout decided)

1. **Default install providers:** Fresh install = **kit + GitHub only**; Cursor/Continue via `HAWP_PROVIDERS`? **Assumption: yes.**

2. **`AGENTS.md`:** Seed at repo root from `providers/.cursor/AGENTS.md.seed`? **Assumption: yes.**

3. **Managed Cursor rules:** Only `hawp-*.mdc` refreshed? **Assumption: yes.**

4. **Continue v1:** Rules only under `.continue/rules/`; defer `config.yaml` seed to TASK-073? **Assumption: yes.**

5. **Install guides:** One generated guide with provider sections vs split files? **Assumption: one guide.**

6. **Root `.github/` in source repo:** This HAWP repo also has `.github/` at repo root (workflows, prompts for *this* repo). Only **`core/providers/.github/`** is the install pack—confirm no confusion. **Assumption: keep both; install pack is only under `core/providers/.github/`.**

---

### Assumptions

- **`core/.github/` is removed after move** to `core/providers/.github/` (with one-release script fallback).
- Layout under `core/providers/` uses **dot-prefixed** folder names: `.github`, `.cursor`, `.continue`.
- `core/.continue/` path is **`core/providers/.continue/`**; implemented in TASK-073.
- Install scripts always source from **`$SRC` = `core/`** root (existing archive layout unchanged).
- Provider materialization is **copy/cp**, not a runtime engine—consistent with HAWP non-runtime principle.
- Continue and Cursor are **guidance overlays** pointing at `.hawp/kit/**`, not replacements for kit.
- Implementation work gets a **new ID (TASK-074)**; TASK-072 closes when planning is approved.

---

### Work Coordination

**Owner:** unassigned
**Implementation status:** not-started (planning only)
**DA assignment:** **No.** This plan is not handed off to a digital agent yet. Keep owner `unassigned` until a human approves the plan and explicitly starts TASK-074 with file-tracking per `.hawp/kit/instructions/da-file-tracking.md`.
**Overlapping files:** `core/**`, `distribution/sources/**`, `distribution/generated/**`, `librarian/scripts/distribution/**`, `README.md`, `core/.hawp/kit/standards/docs/hawp-install-update-safety.md`

**Parallel work risk:** medium (TASK-070 standards text may affect shared instruction wording)

**Can implement now:** no — planning only; implementation is TASK-074 after approval

---

### References

| Path | Role |
|------|------|
| `distribution/sources/install/script.md` | Current provider copy (GitHub only) |
| `distribution/generated/update-dev.md` | User-facing update guide (683 lines) |
| `distribution/sources/shared/repo-boundaries.md` | Boundary model to extend |
| `core/.github/` (legacy) | Move to `core/providers/.github/` in TASK-074 |
| `core/providers/.github/copilot-instructions.md` | GitHub entry template (post-move) |
| `core/providers/.github/instructions/hawp-intake.instructions.md` | Cursor port source |
| `.hawp/work/active/TASK-072-agent-provider-architecture-cursor-github.md` | This plan |
| `.hawp/work/closed/2026/06/04/TASK-069.md`, `TASK-071.md` | Recent standards commits context |

---

## Outcome (filled at close)

Planning approved with stakeholder decisions above (2026-06-04). Architecture: `core/.hawp/` + `core/providers/{shared,.github,.cursor,.continue}` + `manifest.yaml`. Distribution moves to per-provider generated guides (Phase 2). TASK-074 Phase 1 started: GitHub provider move + script paths.

## Verification (filled at close)

- Stakeholder answered open questions via structured intake (2026-06-04).
- Plan approved; TASK-074 opened for implementation.
- Active plan moved to `.hawp/work/closed/2026/06/04/TASK-072.md`.

## Close Checklist

- [x] Stakeholder answers open questions (or accepts assumptions)
- [x] Plan approved
- [x] TASK-074 (implementation) created
- [ ] ADR optional
- [x] Plan moved to `closed/YYYY/MM/DD/` when planning complete

**Status:**

- [x] Plan written (revision 4 — `distribution/sources/providers/{github,cursor,continue,...}`)
- [x] Approved
- [ ] Implemented (delegated to TASK-074)
- [x] Verified (planning complete)
- [x] Closed
