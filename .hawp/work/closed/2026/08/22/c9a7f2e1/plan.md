---
work-item: c9a7f2e1
type: infrastructure
title: "Cross-platform GitHub Actions CI/CD pipeline (6-binary matrix)"
status: plan-ready
owner: unassigned
created: 2026-07-23
updated: 2026-07-23
---

# GitHub Actions CI/CD Pipeline for v0.0.1 Release

## Mission

Build and deploy binaries for Windows, macOS, and Linux across amd64 and arm64 architectures automatically on version tag push (`v0.0.1`, `v0.1.0`, etc.). Enable seamless distribution and automatic release creation.

---

## Context

**Current State:**
- Search system complete and tested (18/18 tests passing)
- Ready to release as v0.0.1
- Need automated cross-platform binary building
- No CI/CD pipeline currently exists

**Release Target:** July 30 - August 6, 2026

---

## Constraints

| Constraint | Detail |
|-----------|--------|
| **Platforms** | Windows, macOS (Intel + Apple Silicon), Linux (x86_64, ARM64) |
| **Binary formats** | `.exe` (Windows), bare binary (macOS/Linux) |
| **Verification** | SHA256 checksums for every binary |
| **Automation** | Tag push → automatic build + release, no manual intervention |
| **Rollback** | Failed builds don't release (CI blocks bad tags) |
| **Quality** | All tests must pass before binary upload |

---

## Implementation Plan

### Phase 1: GitHub Actions Workflow Setup

**File:** `.github/workflows/release.yml`

```yaml
name: Build & Release
on:
  push:
    tags:
      - 'v*'  # Trigger on version tags (v0.0.1, v0.1.0, etc.)

jobs:
  # ========== BUILD MATRIX ==========
  build:
    name: Build ${{ matrix.os }}-${{ matrix.arch }}
    runs-on: ${{ matrix.runs-on }}
    strategy:
      matrix:
        include:
          # Windows builds
          - os: windows
            arch: amd64
            runs-on: windows-latest
            ext: .exe
          - os: windows
            arch: arm64
            runs-on: windows-latest
            ext: .exe
          # macOS builds
          - os: darwin
            arch: amd64
            runs-on: macos-13  # Intel Macs
            ext: ''
          - os: darwin
            arch: arm64
            runs-on: macos-14  # Apple Silicon
            ext: ''
          # Linux builds
          - os: linux
            arch: amd64
            runs-on: ubuntu-latest
            ext: ''
          - os: linux
            arch: arm64
            runs-on: ubuntu-latest
            ext: ''
    steps:
      - uses: actions/checkout@v4
        
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      # Run tests first - fail if any fail
      - name: Run Tests
        run: |
          INTEGRATION=1 go test ./...
      
      # Compile binary for target OS/arch
      - name: Build binary
        env:
          GOOS: ${{ matrix.os }}
          GOARCH: ${{ matrix.arch }}
          CGO_ENABLED: 0  # Pure Go, no C dependencies
        run: |
          go build -o hawp${{ matrix.ext }} ./librarian/go/cmd/hawp
      
      # Generate SHA256 for verification
      - name: Generate checksum
        run: |
          sha256sum hawp${{ matrix.ext }} > hawp-${{ matrix.os }}-${{ matrix.arch }}.sha256
      
      # Upload artifacts for release job
      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: hawp-${{ matrix.os }}-${{ matrix.arch }}
          path: |
            hawp${{ matrix.ext }}
            hawp-${{ matrix.os }}-${{ matrix.arch }}.sha256

  # ========== RELEASE CREATION ==========
  release:
    name: Create GitHub Release
    needs: build
    runs-on: ubuntu-latest
    if: startsWith(github.ref, 'refs/tags/')
    steps:
      - uses: actions/checkout@v4
      
      - name: Download all artifacts
        uses: actions/download-artifact@v3
      
      - name: Prepare release notes
        run: |
          # Extract version from tag
          VERSION=${GITHUB_REF#refs/tags/}
          echo "VERSION=$VERSION" >> $GITHUB_ENV
          
          # Generate release notes from CHANGELOG
          if [ -f CHANGELOG.md ]; then
            cp CHANGELOG.md release_notes.md
          else
            echo "# Release $VERSION" > release_notes.md
          fi
      
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            hawp-windows-amd64/hawp.exe
            hawp-windows-amd64/hawp-windows-amd64.sha256
            hawp-windows-arm64/hawp.exe
            hawp-windows-arm64/hawp-windows-arm64.sha256
            hawp-darwin-amd64/hawp
            hawp-darwin-amd64/hawp-darwin-amd64.sha256
            hawp-darwin-arm64/hawp
            hawp-darwin-arm64/hawp-darwin-arm64.sha256
            hawp-linux-amd64/hawp
            hawp-linux-amd64/hawp-linux-amd64.sha256
            hawp-linux-arm64/hawp
            hawp-linux-arm64/hawp-linux-arm64.sha256
          body_path: release_notes.md
          draft: false
          prerelease: false
```

**Deliverables:**
- [ ] `.github/workflows/release.yml` created
- [ ] Workflow triggers on tag push
- [ ] All 6 binaries build in parallel
- [ ] SHA256 checksums generated for each
- [ ] GitHub Release auto-created with binaries

---

### Phase 2: Versioning & Release Tagging

**Task:** Make releases one-click: `git tag v0.0.1` + `git push --tags`

**Required files:**
- [ ] `version.go` — Canonical version source
- [ ] `Makefile` or `scripts/release.sh` — Tag creation helper
- [ ] `.github/workflows/release.yml` already in Phase 1

**Version source (Go):**
```go
package main

const Version = "0.0.1"
```

**CLI integration:**
```bash
hawp version  # Prints current version
```

**Release process:**
```bash
# 1. Update version in version.go
# 2. Update CHANGELOG.md
# 3. Commit
git commit -am "Release v0.0.1"

# 4. Tag and push
git tag v0.0.1
git push origin main
git push origin --tags
```

**Deliverables:**
- [ ] Version source embedded in binary
- [ ] `hawp version` command works
- [ ] CHANGELOG.md has release notes template
- [ ] Tag naming docs in README.md

---

### Phase 3: Binary Verification & Download

**Auto-detection in update command:**
```go
import "runtime"

// DetectPlatform returns {os, arch} for current system
func DetectPlatform() (os, arch string) {
    return runtime.GOOS, runtime.GOARCH  // windows/darwin/linux, amd64/arm64
}

// ConstructBinaryURL builds GitHub release URL
func ConstructBinaryURL(version string) string {
    os, arch := DetectPlatform()
    ext := ""
    if os == "windows" { ext = ".exe" }
    return fmt.Sprintf(
        "https://github.com/anthropics/hawp/releases/download/%s/hawp-%s-%s%s",
        version, os, arch, ext,
    )
}
```

**SHA256 verification:**
```go
// VerifyBinary checks downloaded binary against published SHA256
func VerifyBinary(binaryPath string, checksumURL string) error {
    // 1. Download SHA256 file from release
    // 2. Compute SHA256 of downloaded binary
    // 3. Compare — fail if mismatch
}
```

**Deliverables:**
- [ ] Update command detects OS/arch automatically
- [ ] Download URL constructed correctly
- [ ] SHA256 verification implemented
- [ ] Rollback if verification fails

---

## Testing Checklist

### Local Testing (Before Pushing Tag)

- [ ] `go test ./...` passes with `INTEGRATION=1`
- [ ] Build succeeds for all 6 platforms:
  ```bash
  GOOS=windows GOARCH=amd64 go build -o test.exe ./cmd/hawp
  GOOS=darwin GOARCH=arm64 go build -o test-arm64 ./cmd/hawp
  # ... etc for all 6
  ```
- [ ] Binaries execute: `./hawp version` works on target systems
- [ ] SHA256 matches: `sha256sum hawp | diff - hawp.sha256`

### CI Testing (After Tag Push)

- [ ] GitHub Actions workflow triggers
- [ ] All 6 builds complete without error
- [ ] All tests pass before binary upload
- [ ] GitHub Release created automatically
- [ ] All 12 files present (6 binaries + 6 SHA256s)

### Integration Testing (Real Downloads)

- [ ] Download binary from GitHub release
- [ ] Verify SHA256 matches published checksum
- [ ] Binary executes on target platform
- [ ] `hawp search` works on fresh system

---

## Risk Mitigation

| Risk | Mitigation |
|------|-----------|
| ARM64 build fails on Linux | Use `ubuntu-latest` with emulation or self-hosted arm64 runner |
| Windows ARM64 build unavailable | Build on windows-latest, skip if incompatible |
| GitHub Actions quota exceeded | Monitor usage; escalate if hitting limits |
| Release notes missing | Use commit message or template fallback |
| SHA256 mismatch in wild | Include verification tool in binary |

---

## Acceptance Criteria

### Definition of Done
- [ ] `.github/workflows/release.yml` exists and passes syntax check
- [ ] Tag `v0.0.1` triggers automated build
- [ ] All 6 binaries build successfully
- [ ] GitHub Release contains all 12 files (6 binaries + 6 SHA256s)
- [ ] Release notes populated from CHANGELOG.md
- [ ] Download + verify step works end-to-end on 3 platforms (Windows, macOS, Linux)
- [ ] Update command successfully detects platform and downloads correct binary
- [ ] Rollback works if downloaded binary is corrupted

---

## Dependencies

- ✅ Search system complete (Slice 1-3)
- ✅ All 18 tests passing
- ✅ Version source defined
- ⏳ Repo cleanup (separate work item `d1b3e8f4`)
- ⏳ Update command v2 (separate work item `e2c4f9g5`)

---

## Effort Estimate

| Phase | Work | Effort |
|-------|------|--------|
| **Phase 1** | GitHub Actions workflow | 2-3h |
| **Phase 2** | Versioning + tagging | 1h |
| **Phase 3** | Binary verification | 1-2h |
| **Testing** | Local + CI + integration | 1h |
| **Total** | | **5-7h** |

---

## Success Metrics

- ✅ One command tags release and builds 6 binaries
- ✅ No manual binary compilation on release day
- ✅ GitHub Release auto-populated with all files
- ✅ Update command works on all 6 platform combos
- ✅ Users can download + verify in <2 min

---

## Notes

- CGO_ENABLED=0 ensures pure Go builds (no C dependency hell)
- Parallel builds (6 jobs) complete in ~10-15 min total
- ARM64 Linux may need self-hosted runner if ubuntu-latest doesn't support it (check Actions panel)
- Keep CHANGELOG.md updated before tagging — release notes auto-pull from it


## Outcome

Shipped in the `0.0.1` release (tag `0.0.1`, 2026-08-21). Work complete.

## Verification

Release published at https://github.com/sentzunhat/human-ai-workflow-protocol/releases/tag/0.0.1 with all 7 assets.

## Close Checklist

- [x] Work shipped in 0.0.1 release
- [x] Archived to closed/2026/08/22/v001-shipped-cleanup/
