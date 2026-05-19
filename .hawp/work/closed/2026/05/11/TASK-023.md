# Bug / Task: Add docs discoverability for generated install/update guides

**Backlog ID:** TASK-023
**Type:** task
**Reported:** 2026-05-11
**Risk Level:** low

---

### Input (what was reported)

> Add a small docs discoverability pass so users and contributors can find the correct generated install/update guides.

---

### Context

The distribution pipeline is complete and verified, but root docs still foreground authoritative script files and do not clearly direct users to generated stable/dev guides.

---

### Analysis

**Root cause (or most likely cause):**
Discoverability lag: README links were not updated after generated distribution docs became the preferred user entry points.

**Directly verified:**
- Root README links install/update to `core/install.md` and `core/update.md`.
- Generated guides exist under `core/distribution/generated/`.
- No separate `core/README.md` was found; smallest touchpoint is root README.

**Inferred (not yet proven):**
- Users may miss generated guides unless explicitly linked in README.

**Scope — what else is affected:**
- `README.md` only.

---

### Work Coordination

**Owner:** agent
**Implementation status:** in-progress
**Overlapping files:**

- `README.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No active item in BACKLOG currently overlaps README install/update discoverability.

---

### Options

#### Option A — README discoverability block (minimal)

Add a compact section with stable/dev generated guide links and contributor refresh notes.

#### Option B — Add new docs landing file

Create a dedicated distribution docs index and link from README.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
Smallest effective change with low risk and no tooling impact.

**Files to change:**

- `README.md` — add generated-guide pointers + contributor refresh instructions

**What to verify after:**

- [ ] Links and paths in README match current repository structure
- [ ] Trusted verification suite still passes

---

### Implementation Notes

Keep wording concise, avoid script duplication, and preserve source-of-truth framing for `core/install.md` and `core/update.md`.

---

## Outcome (filled at close)

- Added a small discoverability pass to root `README.md`.
- Added a generated-guide quick-links section with stable (`install-main.md`, `update-main.md`) and dev/test-only (`install-dev.md`, `update-dev.md`) pointers.
- Added contributor notes clarifying:
	- authoritative script sources: `core/install.md`, `core/update.md`
	- editable fragments path: `core/distribution/sources/`
	- required refresh commands: `distribution:build`, `distribution:validate`
- Kept scope documentation-only; no changes to distribution generation logic.

## Verification (filled at close)

- [x] Discoverability pointers added to root docs.
	**Evidence:** `README.md` now contains "Generated guide quick links" and "For contributors" sections with current paths.
- [x] Required verification suite passes after docs change.
	**Evidence:** Ran `npm run distribution:validate`, `npm run typecheck`, and `npm run validate:workflow -- --work-root ../.hawp/work` from `librarian/`; all passed.
- [x] Docs/link validation command availability checked.
	**Evidence:** `librarian/package.json` scripts include only `distribution:build`, `distribution:validate`, `typecheck`, and `validate:workflow`; no separate docs/link validator exists.

## Close Checklist

**Before marking this task done in BACKLOG.md, verify ALL:**

- [x] Outcome section filled (what was actually implemented)
- [x] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [x] All evidence files referenced exist in `../evidence/YYYY/MM/DD/` or are noted as inline
- [x] Plan file moved to `../closed/YYYY/MM/DD/<ID>.md`
- [x] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [x] Status report written (optional) or explicitly not required for this low-risk docs-only task
- [x] Decision file created (optional) or explicitly not applicable (no design decision changed)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed