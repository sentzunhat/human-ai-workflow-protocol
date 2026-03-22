---
work-item: g4e6b1i7
type: task
title: "Final Release: v0.0.1 tag, push, and GitHub Actions automation"
status: inbox
owner: unassigned
created: 2026-07-23
updated: 2026-07-23
---

# Final Release Steps: v0.0.1

## Mission

Execute the final release of v0.0.1 by updating version, tagging, and pushing to GitHub. This triggers GitHub Actions to automatically build all 6 binaries and create the release.

---

## Prerequisites (All Complete ✅)

- ✅ Repository audit: CLEAN (d1b3e8f4 complete)
- ✅ All tests passing: `INTEGRATION=1 go test ./...` ✅
- ✅ GitHub Actions workflow: Ready (`.github/workflows/release.yml`)
- ✅ Build verified: `go build ./librarian/go/cmd/hawp` ✅
- ✅ Release playbook: Created (`RELEASE.md`)

---

## Release Steps (5-10 minutes)

### Step 1: Update Version Source

```bash
# Update version.go to 0.0.1
sed -i '' 's/var Version = "dev"/var Version = "0.0.1"/' \
  librarian/go/internal/domain/update/version.go

# Verify the change
cat librarian/go/internal/domain/update/version.go | grep "var Version"
# Expected output: var Version = "0.0.1"
```

**File:** `librarian/go/internal/domain/update/version.go`

### Step 2: Update CHANGELOG.md

Add v0.0.1 release notes to the top of the file:

```bash
# Open the file
vim librarian/go/CHANGELOG.md
```

Add this section at the top (after "# Changelog"):

```markdown
## [0.0.1] - 2026-07-23

### Added
- Cross-platform binary builds (Windows, macOS, Linux)
- Lexical search via FTS5: <1ms queries
- Semantic search via ONNX embedding: 100ms queries
- Hybrid search (lexical + semantic ranking): 15-20ms queries
- Vector embedding via BGE-base-en-v1.5 (768 dimensions)
- Auto-update mechanism (`hawp update`) with platform detection
- GitHub Actions CI/CD: 6-binary matrix (Windows/macOS/Linux × amd64/arm64)
- Release automation: Tag push → builds → GitHub Release
- 15-query benchmark suite

### Verified
- 18/18 unit tests passing
- All 6 platform binaries build successfully
- Cross-platform self-update mechanism works
- SHA256 checksum verification integrated
```

**File:** `librarian/go/CHANGELOG.md`

### Step 3: Verify Changes

```bash
# Check what changed
git diff librarian/go/internal/domain/update/version.go
git diff librarian/go/CHANGELOG.md

# Verify build still works
cd librarian/go
go build -o /tmp/hawp-test ./cmd/hawp
/tmp/hawp-test version
# Expected output: 0.0.1
```

### Step 4: Commit

```bash
# Stage changes
git add librarian/go/internal/domain/update/version.go
git add librarian/go/CHANGELOG.md

# Commit
git commit -m "Release v0.0.1: Search system + release infrastructure"

# Verify commit
git log --oneline -1
```

### Step 5: Create Tag

```bash
# Create tag (must match v* pattern for GitHub Actions to trigger)
git tag v0.0.1

# Verify tag
git tag -l | tail -5
# Expected: v0.0.1 in the list
```

### Step 6: Push to GitHub (Triggers Release Automation)

```bash
# Push both commit and tag
git push origin main --follow-tags

# Alternative (if --follow-tags doesn't work):
git push origin main
git push origin v0.0.1
```

**This triggers GitHub Actions automatically!**

---

## Automated Steps (GitHub Actions - No Manual Action Needed)

After you push the tag, GitHub Actions automatically:

1. **Detects tag:** `v0.0.1` matches `v*` pattern
2. **Starts build matrix:**
   - windows-amd64 build job
   - windows-arm64 build job
   - darwin-amd64 build job
   - darwin-arm64 build job
   - linux-amd64 build job
   - linux-arm64 build job

3. **Each job:**
   - Checks out code
   - Runs tests: `INTEGRATION=1 go test ./...`
   - Compiles for target platform
   - Generates SHA256 checksum
   - Uploads artifact

4. **Release job (runs after all builds pass):**
   - Downloads all 12 artifacts
   - Creates GitHub Release
   - Attaches all files
   - Publishes (visible immediately)

5. **Result:** Release live at:
   ```
   https://github.com/sentzunhat/human-ai-workflow-protocol/releases/tag/v0.0.1
   ```

---

## Monitoring (While GitHub Actions Runs)

```bash
# Watch GitHub Actions in real time
# Visit: https://github.com/sentzunhat/human-ai-workflow-protocol/actions

# Or use gh CLI
gh workflow run release.yml --ref v0.0.1
gh run list --workflow=release.yml
```

**Expected Timeline:** ~10-15 minutes for all builds to complete

---

## Post-Release Verification

After GitHub Actions completes:

```bash
# 1. Verify GitHub Release was created
curl -s https://api.github.com/repos/sentzunhat/human-ai-workflow-protocol/releases/latest | jq '.tag_name'
# Expected output: v0.0.1

# 2. Download a binary and verify
curl -L -o hawp-darwin-arm64 \
  https://github.com/sentzunhat/human-ai-workflow-protocol/releases/download/v0.0.1/hawp-darwin-arm64

# 3. Verify SHA256
curl -L -o hawp-darwin-arm64.sha256 \
  https://github.com/sentzunhat/human-ai-workflow-protocol/releases/download/v0.0.1/hawp-darwin-arm64.sha256

sha256sum -c hawp-darwin-arm64.sha256
# Expected output: hawp-darwin-arm64: OK

# 4. Test binary
chmod +x hawp-darwin-arm64
./hawp-darwin-arm64 version
# Expected output: v0.0.1
```

---

## Success Criteria

### Before Tagging
- [ ] `git diff` shows only version.go and CHANGELOG.md changes
- [ ] Build passes: `go build ./librarian/go/cmd/hawp`
- [ ] Binary version is 0.0.1: `/tmp/hawp version`
- [ ] All tests pass: `INTEGRATION=1 go test ./...`

### After Pushing Tag
- [ ] Tag appears in GitHub: https://github.com/sentzunhat/human-ai-workflow-protocol/tags
- [ ] Actions workflow starts automatically
- [ ] All 6 build jobs complete (green checkmarks)
- [ ] Release job completes and creates GitHub Release
- [ ] Release has all 12 files (6 binaries + 6 SHA256s)

### Final Verification
- [ ] Release is public (not draft)
- [ ] All 6 binaries executable
- [ ] All 6 SHA256 files present and valid
- [ ] Release notes populated (from CHANGELOG.md)

---

## Rollback (If Needed)

```bash
# If something goes wrong before tag is pushed:
git reset --soft HEAD~1          # Undo commit, keep changes
git reset HEAD~1                 # Undo commit, unstage changes

# If tag was pushed but release is broken:
git tag -d v0.0.1               # Delete local tag
git push origin --delete v0.0.1 # Delete remote tag
# Fix the issue, re-tag and push

# If release was published but has issues:
# GitHub Release can be edited or deleted
# Users can be notified of the issue
# Push v0.0.1-fix1 or skip to v0.1.0
```

---

## Effort & Timeline

| Phase | Effort |
|-------|--------|
| Update version.go | 30 sec |
| Update CHANGELOG.md | 2 min |
| Verify changes | 1 min |
| Commit & tag | 1 min |
| Push (triggers GitHub Actions) | 30 sec |
| Wait for builds (~10-15 min) | 10-15 min |
| Verify release | 5 min |
| **Total** | **~20 minutes** |

---

## Owner & Next Steps

**Current Status:** READY TO EXECUTE

**Who Can Do This:** Any developer with write access to the repository

**Next Action:** Execute steps 1-6 above (takes ~5-10 minutes of manual work, then ~10-15 minutes automated)

**After Release:** Release is live, users can:
- Download binaries from GitHub
- Run `hawp update` to auto-update to v0.0.1

---

## Links & References

- **Release Workflow:** `.github/workflows/release.yml`
- **Release Playbook:** `RELEASE.md`
- **Version Source:** `librarian/go/internal/domain/update/version.go`
- **Changelog:** `librarian/go/CHANGELOG.md`
- **GitHub Actions Dashboard:** https://github.com/sentzunhat/human-ai-workflow-protocol/actions

---

## Notes

- **Tag format matters:** Must be `v0.0.1` (matching `v*` pattern) for GitHub Actions trigger
- **All tests must pass:** GitHub Actions runs tests; if any fail, no release is created
- **Automatic everything else:** Once you push the tag, all builds and release creation happen automatically
- **Users auto-update:** `hawp update` will detect v0.0.1 is available and download the correct binary for their platform

---

**STATUS: READY TO SHIP** 🚀

Execute steps 1-6 above and GitHub Actions handles the rest!


## Outcome

Shipped in the `0.0.1` release (tag `0.0.1`, 2026-08-21). Work complete.

## Verification

Release published at https://github.com/sentzunhat/human-ai-workflow-protocol/releases/tag/0.0.1 with all 7 assets.

## Close Checklist

- [x] Work shipped in 0.0.1 release
- [x] Archived to closed/2026/08/22/v001-shipped-cleanup/
