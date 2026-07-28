# Release Playbook: v0.0.1 and Beyond

**Purpose:** One-click releases for HAWP. Tag push → Automated builds (6 platforms) → GitHub Release.

---

## Quick Start

### Option 1: Manual Tag Push (Recommended for First Release)

```bash
# 1. Update version in version.go
sed -i '' 's/var Version = "dev"/var Version = "0.0.1"/' librarian/go/internal/domain/update/version.go

# 2. Update CHANGELOG.md with release notes
# (Add new section: ## [0.0.1] - 2026-07-23)

# 3. Commit
git add -A
git commit -m "Release v0.0.1"

# 4. Create tag and push
git tag v0.0.1
git push origin main
git push origin --tags

# 5. GitHub Actions takes over automatically
# - Workflow triggers on tag push
# - Builds all 6 binaries in parallel
# - Creates GitHub Release with all files
# - Releases appear in https://github.com/sentzunhat/human-ai-workflow-protocol/releases
```

**Timeline:** Push tag → All builds complete in ~10-15 minutes → Release live

---

### Option 2: GitHub Actions UI (Manual Dispatch)

1. Go to **Actions** tab in GitHub
2. Select **"Build & Release"** workflow
3. Click **"Run workflow"** button
4. Enter version: `0.0.1`
5. Click **"Run workflow"**
6. Wait for builds to complete
7. Release automatically created

**Note:** Requires version.go and CHANGELOG.md already updated locally.

---

## Release Process Detail

### Step 1: Update Version Source

**File:** `librarian/go/internal/domain/update/version.go`

```go
// BEFORE
var Version = "dev"

// AFTER
var Version = "0.0.1"
```

**Why:** This version gets embedded in every binary at compile time via ldflags.

### Step 2: Update CHANGELOG

**File:** `librarian/go/CHANGELOG.md`

```markdown
## [0.0.1] - 2026-07-23

### Added
- Cross-platform binary builds (Windows, macOS, Linux)
- Lexical search (FTS5): <1ms queries
- Semantic search (Cosine similarity): 100ms queries
- Hybrid search (blended ranking): 15-20ms queries
- Vector embedding via ONNX (BGE-base-en-v1.5)
- Auto-update mechanism (`hawp update`)

### Fixed
- Transaction persistence for vector embeddings

### Verified
- 18/18 unit tests passing
- 15-query benchmark suite (3 patterns)
- Cross-platform binaries tested
```

### Step 3: Commit and Tag

```bash
# Stage changes
git add librarian/go/internal/domain/update/version.go
git add librarian/go/CHANGELOG.md

# Commit
git commit -m "Release v0.0.1: Search system + release infrastructure"

# Tag (exact format matters — must match v* pattern)
git tag v0.0.1

# Push both commit and tag
git push origin main --follow-tags
# OR explicitly:
git push origin main
git push origin --tags
```

**Critical:** Tag MUST match `v*` pattern (e.g., `v0.0.1`, `v0.1.0`, `v1.0.0`)

---

## What Happens Next (Automated)

```
Tag push: v0.0.1
    ↓
GitHub detects tag
    ↓
Workflow trigger: `.github/workflows/release.yml`
    ↓
Build job matrix starts (6 jobs in parallel):
    • windows-amd64: GOOS=windows GOARCH=amd64 go build
    • windows-arm64: GOOS=windows GOARCH=arm64 go build
    • darwin-amd64:  GOOS=darwin GOARCH=amd64 go build
    • darwin-arm64:  GOOS=darwin GOARCH=arm64 go build
    • linux-amd64:   GOOS=linux GOARCH=amd64 go build
    • linux-arm64:   GOOS=linux GOARCH=arm64 go build
    ↓
Each build:
    • Checkout code
    • Run: INTEGRATION=1 go test ./...  (FAIL BUILD IF TESTS FAIL)
    • Compile for target OS/arch
    • Generate SHA256 checksum
    • Upload artifact
    ↓
All 6 builds must pass
    ↓
Release job starts:
    • Download all 12 artifacts (6 binaries + 6 checksums)
    • Create GitHub Release
    • Attach all 12 files
    • Publish (visible immediately)
    ↓
Release live at:
https://github.com/sentzunhat/human-ai-workflow-protocol/releases/tag/v0.0.1
    ↓
Users can:
    • Download binaries for their platform
    • Verify SHA256 checksums
    • Run `hawp update` to auto-download latest
```

---

## Verification Checklist

### Before Tagging

- [ ] Version bumped in `version.go`
- [ ] CHANGELOG.md updated with release notes
- [ ] Commit message is clear and descriptive
- [ ] All tests pass locally: `INTEGRATION=1 go test ./...`
- [ ] No uncommitted changes: `git status` is clean

### After Tag Push

- [ ] Go to Actions tab → "Build & Release" workflow
- [ ] All 6 build jobs complete successfully (green checkmarks)
- [ ] Release job completes and creates GitHub Release
- [ ] GitHub Release has all 12 files (6 binaries + 6 SHA256s):
  ```
  hawp-windows-amd64.exe
  hawp-windows-amd64.sha256
  hawp-windows-arm64.exe
  hawp-windows-arm64.sha256
  hawp-darwin-amd64
  hawp-darwin-amd64.sha256
  hawp-darwin-arm64
  hawp-darwin-arm64.sha256
  hawp-linux-amd64
  hawp-linux-amd64.sha256
  hawp-linux-arm64
  hawp-linux-arm64.sha256
  ```
- [ ] Release notes populated (auto-pulled from CHANGELOG)
- [ ] Each binary has matching SHA256 file

### Download & Verify

```bash
# Download binary for your platform
# Example: macOS arm64
curl -L -o hawp https://github.com/sentzunhat/human-ai-workflow-protocol/releases/download/v0.0.1/hawp-darwin-arm64
curl -L -o hawp.sha256 https://github.com/sentzunhat/human-ai-workflow-protocol/releases/download/v0.0.1/hawp-darwin-arm64.sha256

# Verify checksum
sha256sum -c hawp.sha256
# Expected output: hawp: OK

# Make executable
chmod +x hawp

# Test
./hawp version
# Expected output: v0.0.1

./hawp search "test query"
# Expected output: Search results from your index
```

---

## Rollback (If Release is Broken)

```bash
# Option 1: Delete the tag locally and remotely
git tag -d v0.0.1
git push origin --delete v0.0.1
# Fix the issue
# Re-tag and push

# Option 2: Re-release from a corrected commit
# (GitHub Release can be edited, but binary rebuild requires new tag)
git tag v0.0.1-fixed
git push origin --tags
```

---

## Release Cadence

- **v0.0.1** — Initial release (Slices 1-3 complete)
- **v0.1.0** — Context packing (Slice 4)
- **v0.2.0** — Provider overlays + agentic loops
- **v1.0.0** — Stable release candidate

Each release follows the same process: update version → update CHANGELOG → tag → push → automated builds.

---

## Automation Details

### Version Detection

The binary embeds its version at build time:

```bash
# Build command (in GitHub Actions)
GOOS=darwin GOARCH=arm64 go build -ldflags="-X github.com/sentzunhat/human-ai-workflow-protocol/librarian/go/internal/domain/update.Version=v0.0.1" -o hawp ./librarian/go/cmd/hawp

# User verification
./hawp version
# Output: v0.0.1
```

### Update Mechanism

Users can self-update to latest release:

```bash
hawp update
# Detects current version (embedded in binary)
# Queries GitHub for latest release
# Downloads platform-specific binary
# Verifies SHA256 checksum
# Atomically replaces binary
# Confirms new version
```

---

## Troubleshooting

### Workflow Doesn't Trigger

**Problem:** Pushed tag but GitHub Actions didn't start

**Solution:**
```bash
# Verify tag format (must start with 'v')
git tag -l

# Check tag exists locally and remotely
git push origin --tags

# Manually trigger from Actions tab as fallback
```

### Build Fails (Tests or Compilation)

**Problem:** Workflow shows red X on one or more build jobs

**Solution:**
1. Click on failed job in GitHub Actions
2. Read error output
3. Fix locally: `INTEGRATION=1 go test ./...` to reproduce
4. Push fix to main
5. Re-tag and push (e.g., `v0.0.1-fix1`)

### Release Has Wrong Files

**Problem:** GitHub Release created but missing binaries

**Solution:**
1. Check if all 6 build jobs passed (one may have failed silently)
2. Delete release from GitHub
3. Delete tag: `git push origin --delete v0.0.1`
4. Fix the issue
5. Re-tag and push

### SHA256 Mismatch

**Problem:** Downloaded binary fails checksum verification

**Solution:**
1. Try re-downloading (network corruption possible)
2. If persistent, delete release and rebuild with latest code
3. Verify using: `sha256sum -c <file>.sha256`

---

## For Next Release (v0.1.0)

1. Update `version.go`: `"0.1.0"`
2. Update `CHANGELOG.md`: Add new `## [0.1.0]` section
3. Commit: `git commit -am "Release v0.1.0: Context packing"`
4. Tag: `git tag v0.1.0`
5. Push: `git push origin main --follow-tags`
6. Wait ~15 min for builds
7. Verify release on GitHub
8. Done! Users auto-update via `hawp update`

---

## Questions?

- **GitHub Actions troubleshooting:** See `.github/workflows/release.yml`
- **Version format:** Must be `vX.Y.Z` (semantic versioning)
- **CHANGELOG format:** Follow existing pattern with `## [X.Y.Z] - YYYY-MM-DD`
- **Checksums:** Generated automatically by workflow; no manual action needed
