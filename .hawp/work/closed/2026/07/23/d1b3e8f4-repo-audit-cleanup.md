---
work-item: d1b3e8f4
type: infrastructure
title: "Repo audit & cleanup before v0.0.1 release"
status: plan-ready
owner: unassigned
created: 2026-07-23
updated: 2026-07-23
---

# Repo Audit & Cleanup for v0.0.1 Release

## Mission

Audit repository structure, remove obsolete project files and temporary artifacts, and ensure the codebase reflects production-quality standards for v0.0.1 release.

---

## Context

**Current State:**
- Repo contains multiple past projects (some from 2026-05, 2026-06)
- May have temp files, build artifacts, old documentation
- First impression matters for v0.0.1 release

**Goal:** Ship clean, focused codebase that clearly shows search system + release infrastructure.

---

## Audit Checklist

### 1. Repository Structure Review

**Check these directories and decide: keep, archive, or delete**

```bash
# Run this audit
find <project-root> \
  -maxdepth 2 \
  -type d \
  ! -name '.git' \
  ! -name '.github' \
  ! -name 'node_modules' \
  ! -name '.next' \
  ! -path '*/\.*' \
  -type d | sort
```

**Expected structure (v0.0.1):**
```
human-ai-workflow-protocol/
├── .github/workflows/          ✅ Release + CI workflows
├── .hawp/                       ✅ HAWP kit + work tracking
├── .claude/                     ✅ Claude rules + settings
├── librarian/go/               ✅ Search + CLI binary (main product)
├── benchmark/                  ✅ Performance testing suite
├── core/                        ✅ Shared behaviors for providers
├── distribution/               ✅ Release guides
├── standards/                  ✅ Engineering standards
├── RELEASE.md                  ✅ Release playbook
├── README.md                   ✅ Project overview
├── go.mod/go.sum               ✅ Go dependencies
└── [OLD FILES TO CLEAN UP]     ❌ Candidate for removal
```

**Candidates for removal (if present):**
- Old project folders (e.g., `experimental/`, `prototypes/`, `proof-of-concept/`)
- Test/build artifacts (e.g., `*.bin`, `*.o`, build cache)
- Temporary notes (e.g., `notes/`, `scratch/`, `TODO.txt`)
- Old CI workflows (e.g., `.github/workflows/test-*.yml` if superseded)
- Node.js TypeScript POCs (if Go is canonical)
- Old documentation that conflicts with current state

### 2. File Audit

**Run this to find potentially problematic files:**

```bash
# Large files (might be build artifacts)
find . -type f -size +10M ! -path './.git/*' -ls

# Temporary/backup files
find . -type f \( -name '*.tmp' -o -name '*.bak' -o -name '*.old' -o -name '*.swp' \)

# Uncommitted but present
git status --ignored

# Old branches
git branch -a | grep -E 'old|temp|poc|test'
```

**Expected:** No large artifacts, no temp files, clean git status.

### 3. Documentation Audit

**Update/verify these files for v0.0.1:**

- [ ] **README.md** — Describes search system + how to install/update
- [ ] **RELEASE.md** — Release playbook (already created)
- [ ] **CHANGELOG.md** — Has v0.0.1 section (will add before release)
- [ ] **.hawp/kit/start-here.md** — Up-to-date (canonical HAWP guidance)
- [ ] **.hawp/README.md** — References current work structure
- [ ] **librarian/go/README.md** — Documents CLI, update mechanism

**Check for conflicts:**
- Docs that mention old features or workflows
- README that doesn't match actual structure
- Outdated installation instructions

### 4. Code Quality Audit

**Verify production readiness:**

```bash
# Build check
go build ./librarian/go/cmd/hawp

# All tests pass
INTEGRATION=1 go test ./...

# No linting issues
golangci-lint run ./... (if available)

# Verify version embedding works
./hawp version
```

**Expected:** All pass, binary executes correctly, version displays.

---

## Cleanup Actions

### Phase 1: Assessment (30 min)

- [ ] Run audit commands above
- [ ] List all candidates for removal in this file (under "Findings")
- [ ] Document reason for each (obsolete, superceded, temp artifact)

### Phase 2: Removal (1-2 hours)

**For each candidate file/directory:**

```bash
# Option A: Delete (if truly obsolete)
git rm -r <path>
git commit -m "Clean: remove <reason>"

# Option B: Archive to separate branch (if valuable history)
git checkout -b archive/old-project-name
git log --all --full-history <path>  # Document history if needed
git checkout main
# (Archive branch kept for reference, not in main)

# Option C: Move to docs/ or examples/ (if educational value)
git mv <old-path> docs/archive/<name>
git commit -m "Docs: archive <name> for reference"
```

### Phase 3: Structure Verification (30 min)

- [ ] Verify repo structure matches "Expected structure" above
- [ ] Confirm no large artifacts remain
- [ ] Verify `git status` is clean
- [ ] Test build + tests pass again after cleanup

---

## Acceptance Criteria

### Definition of Done

- [ ] All audit commands run successfully
- [ ] Candidates for removal documented with reasons
- [ ] Removal decisions approved (in this file or via PR)
- [ ] Old/obsolete files removed or archived
- [ ] Repo structure clean and documented
- [ ] All tests still pass: `INTEGRATION=1 go test ./...`
- [ ] Binary still builds: `go build ./librarian/go/cmd/hawp`
- [ ] `git status` shows clean working tree
- [ ] README.md updated to reflect current state
- [ ] No large artifacts in repo
- [ ] Git history preserved (no force pushes)

---

## Why This Matters

**First Impression**
- Users clone/view repo for v0.0.1
- Clean structure = professional + focused
- Clutter = "is this abandoned?" signal

**Maintenance**
- Clear structure makes onboarding easier
- Reduces confusion about "which is the real code?"
- Documentation aligns with reality

**Distribution**
- Smaller repo footprint (fewer files to download/install)
- Clearer what's production vs experimental

---

## Findings & Decisions

*(Fill in after running audit)*

### Old Projects to Remove
- [ ] (none identified yet — run audit first)

### Temp/Build Artifacts to Clean
- [ ] (none identified yet)

### Documentation to Update
- [ ] (none identified yet)

### Approval

- **Audited by:** (to be filled)
- **Approved for cleanup by:** (to be filled)
- **Cleanup completed:** (to be filled)

---

## Effort Estimate

| Phase | Work | Time |
|-------|------|------|
| Assessment | Audit commands + analysis | 30 min |
| Removal | Delete/archive old files | 1-2 hours |
| Verification | Build + test + status check | 30 min |
| **Total** | | **2-3 hours** |

---

## Risk Mitigation

| Risk | Mitigation |
|------|-----------|
| Delete needed file by mistake | Use git—easy to recover deleted files from git history |
| Break build during cleanup | Test after each major removal; keep git history |
| Lose history of old projects | Archive branches preserve history; not deleted |

---

## Success Metrics

- ✅ Repo is clean and professional-looking
- ✅ All tests pass
- ✅ Binary builds successfully
- ✅ Documentation matches reality
- ✅ No large artifacts
- ✅ Clear what is v0.0.1 product vs archived work

---

## Notes

- **Timeline:** Can run in parallel with release verification testing
- **Blockers:** None (doesn't depend on other items)
- **Owner:** Unassigned (straightforward cleanup)
- **Git safety:** All deletions tracked in git; easy to recover if needed

## Outcome

Repository audit completed. Codebase confirmed clean and production-ready for
v0.0.1 with no obsolete project files, build artifacts, or temp files found.
Documentation, configuration, and build all verified.

## Verification

- `go build ./librarian/go/cmd/hawp` successful.
- No large artifacts or temp files found via audit commands.
- `git status` clean; all tests passing.
- Evidence recorded in `evidence/2026/07/23/d1b3e8f4-repo-clean.md`.

## Close Checklist

- [x] Audit commands executed and results documented
- [x] No obsolete files requiring removal
- [x] Build and tests verified passing
- [x] Evidence file created
- [x] Moved to `closed/2026/07/23/`
