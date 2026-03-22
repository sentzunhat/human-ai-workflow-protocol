---
work-item: f3d5a0h6
type: test
title: "Release verification: test all 6 binaries cross-platform"
status: plan-ready
owner: unassigned
created: 2026-07-23
updated: 2026-07-23
---

# Release Verification: Cross-Platform Binary Testing

## Mission

Test all 6 release binaries (Windows/macOS/Linux × amd64/arm64) to verify they execute correctly, pass SHA256 checksums, and function as expected before v0.0.1 ships.

---

## Context

**Current State:**
- GitHub Actions workflow compiles all 6 binaries automatically
- Binaries uploaded to GitHub Release
- Need verification that all work on actual target platforms

**Goal:** Confirm v0.0.1 release is production-ready across all supported platforms.

---

## Test Matrix

| OS | Arch | Binary Name | Tester Role | Status |
|----|------|------------|------------|--------|
| Windows | amd64 | hawp-windows-amd64.exe | Dev/CI | pending |
| Windows | arm64 | hawp-windows-arm64.exe | CI or ARM tester | pending |
| macOS | amd64 | hawp-darwin-amd64 | Intel Mac user | pending |
| macOS | arm64 | hawp-darwin-arm64 | Apple Silicon user | pending |
| Linux | amd64 | hawp-linux-amd64 | Linux/CI runner | pending |
| Linux | arm64 | hawp-linux-arm64 | ARM64 Linux user | pending |

---

## Test Protocol

### Phase 1: SHA256 Verification (All Platforms)

**For each platform:**

```bash
# 1. Download release binaries and checksums
VERSION="v0.0.1"
PLATFORM="darwin-arm64"  # Replace with target platform

curl -L -o "hawp-${PLATFORM}" \
  "https://github.com/sentzunhat/human-ai-workflow-protocol/releases/download/${VERSION}/hawp-${PLATFORM}"

curl -L -o "hawp-${PLATFORM}.sha256" \
  "https://github.com/sentzunhat/human-ai-workflow-protocol/releases/download/${VERSION}/hawp-${PLATFORM}.sha256"

# 2. Verify checksum
sha256sum -c hawp-${PLATFORM}.sha256
# Expected: hawp-${PLATFORM}: OK

# 3. Make executable (Unix only)
chmod +x "hawp-${PLATFORM}"
```

**Pass Criteria:**
- [ ] SHA256 checksum matches
- [ ] File is executable
- [ ] No download corruption detected

### Phase 2: Binary Execution Tests

#### Test 2a: Version Check

**All platforms:**

```bash
./hawp version
# Expected output: v0.0.1
```

**Pass Criteria:**
- [ ] Command exits with code 0
- [ ] Output matches "v0.0.1"
- [ ] No error messages

#### Test 2b: Help & Commands

**All platforms:**

```bash
./hawp --help
./hawp commands
./hawp commands --json

# Sample output formats:
# --help: Shows usage
# commands: Lists available commands
# commands --json: Valid JSON output
```

**Pass Criteria:**
- [ ] All commands exit cleanly (code 0)
- [ ] Output is readable and useful
- [ ] JSON is valid (can parse with `jq`)

#### Test 2c: Search System (If Database Available)

**Prerequisite:** Database initialized on test system

```bash
# Check if database exists
./hawp search index --list

# Run test query
./hawp search "test query"
```

**Pass Criteria:**
- [ ] Command executes without crashing
- [ ] Search results display correctly
- [ ] No segfaults or memory errors

---

## Platform-Specific Tests

### Windows (Both amd64 & arm64)

```bash
# Download binary
curl -L -o hawp.exe https://github.com/sentzunhat/human-ai-workflow-protocol/releases/download/v0.0.1/hawp-windows-amd64.exe

# Verify signature (if code-signed)
# (Optional: signtool verify hawp.exe)

# Test execution
hawp.exe version
hawp.exe --help

# Check exit codes
echo %ERRORLEVEL%  # Should be 0
```

**Windows-Specific:**
- [ ] `.exe` extension correct
- [ ] Runs without UAC prompt (unsigned expected for now)
- [ ] Paths with spaces handled correctly
- [ ] Terminal output readable

### macOS (Intel amd64 & Apple Silicon arm64)

```bash
# Download
curl -L -o hawp-darwin-arm64 https://github.com/sentzunhat/human-ai-workflow-protocol/releases/download/v0.0.1/hawp-darwin-arm64
chmod +x hawp-darwin-arm64

# Remove quarantine flag (if present)
xattr -d com.apple.quarantine ./hawp-darwin-arm64 2>/dev/null

# Test
./hawp-darwin-arm64 version

# Check architecture
file ./hawp-darwin-arm64
# Expected: ELF 64-bit LSB executable, ARM aarch64 (or x86-64 for amd64)

# Architecture-specific check
lipo -info ./hawp-darwin-arm64
# Expected: Mach-O 64-bit executable arm64 (or x86_64)
```

**macOS-Specific:**
- [ ] Executes without codesign warnings
- [ ] Architecture matches platform (arm64 on M1+, amd64 on Intel)
- [ ] No "app is damaged" warnings
- [ ] Quarantine flag removable if present

### Linux (Both amd64 & arm64)

```bash
# Download
curl -L -o hawp-linux-amd64 https://github.com/sentzunhat/human-ai-workflow-protocol/releases/download/v0.0.1/hawp-linux-amd64
chmod +x hawp-linux-amd64

# Test execution
./hawp-linux-amd64 version

# Check binary info
file ./hawp-linux-amd64
# Expected: ELF 64-bit LSB executable, x86-64 (or ARM64)

# Check dependencies (should have none for pure Go)
ldd ./hawp-linux-amd64 2>&1 | head -10
# Expected: mostly empty or just libc/glibc

# Verify no CGO dependency
nm ./hawp-linux-amd64 | grep -i "^.*cgo" || echo "No CGO"
```

**Linux-Specific:**
- [ ] Executes cleanly
- [ ] ELF binary format correct
- [ ] No missing dependencies (pure Go, CGO_ENABLED=0)
- [ ] Works across different libc versions (glibc, musl)

---

## Integration Test: Self-Update

**On any platform with v0.0.1 already installed:**

```bash
# Current version
hawp version
# Output: v0.0.1

# Check for update (should report already up-to-date)
hawp update --check
# Expected output:
#   current: v0.0.1
#   latest:  v0.0.1
#   Already up to date.

# (If v0.0.2 was released, verify download works)
hawp update
# Expected output: Successfully updated to v0.0.2
```

**Pass Criteria:**
- [ ] Version detection works
- [ ] Update check succeeds
- [ ] Platform-specific binary correctly identified
- [ ] Self-update mechanism works (if newer version available)

---

## Sign-Off Checklist

### For Each Platform

- [ ] SHA256 verification passed
- [ ] Binary executes: `./hawp version` shows v0.0.1
- [ ] Binary executes: `./hawp --help` shows help text
- [ ] Binary info correct: `file` and `lipo` (macOS) or `file` and `ldd` (Linux)
- [ ] No crashes or segfaults on basic commands
- [ ] Search works (if database available)
- [ ] Platform-specific requirements met

### Overall Release Sign-Off

- [ ] Windows amd64 verified ✅
- [ ] Windows arm64 verified ✅
- [ ] macOS amd64 verified ✅
- [ ] macOS arm64 verified ✅
- [ ] Linux amd64 verified ✅
- [ ] Linux arm64 verified ✅
- [ ] All SHA256 checksums valid ✅
- [ ] All binaries execute correctly ✅
- [ ] No regressions vs Slice 3 baseline ✅

---

## Test Results Template

**Run after each platform test and save to evidence:**

```
Platform: [OS-arch]
Tester: [name]
Date: 2026-07-XX

SHA256 Verification:
  ✅ Downloaded successfully
  ✅ Checksum valid
  ✅ File size reasonable (~50-80 MB expected)

Binary Execution:
  ✅ ./hawp version → v0.0.1
  ✅ ./hawp --help → help text
  ✅ ./hawp commands --json → valid JSON
  ✅ Exit codes: 0 (success)

Platform-Specific:
  ✅ [Windows: No UAC / macOS: No quarantine / Linux: No deps]
  ✅ Architecture correct
  ✅ File format correct

Search System (if tested):
  ✅ Database loads
  ✅ Query executes
  ✅ Results display

Overall: PASS ✅
```

---

## Risk Mitigation

| Risk | Likelihood | Mitigation |
|------|-----------|-----------|
| Binary doesn't exist in release | Low | Check GitHub Release page, re-run build |
| SHA256 mismatch | Low | Re-download, verify network; check Actions logs |
| Platform-specific library missing | Medium | Verify CGO_ENABLED=0; check ldd output |
| Quarantine flag on macOS | Medium | Use `xattr -d` to remove |
| ARM64 Linux untestable (no hardware) | Medium | Use Actions CI runner or skip; document |
| User reports crash after ship | Low | Keep v0.0.1 tag; ability to quickly push v0.0.1-fix1 |

---

## Timeline & Effort

| Phase | Work | Time | Effort |
|-------|------|------|--------|
| Download + verify | Get all 6 binaries, checksums | 10 min | 1 person |
| SHA256 tests | All 6 platforms | 15 min | 1 person |
| Execution tests | All 6 platforms | 20 min | 1 person |
| Platform-specific | OS-specific validation | 15 min | platform-dependent |
| Integration test | Self-update mechanism | 10 min | 1 platform |
| Documentation | Save results, sign off | 10 min | 1 person |
| **Total** | | **~80 min** | — |

**Can be parallelized:** Different people test different platforms in parallel.

---

## Success Criteria

### Definition of Done
- [ ] All 6 binaries verified (SHA256 + execution)
- [ ] All platform-specific checks pass
- [ ] No crashes, segfaults, or unexpected errors
- [ ] Version detection correct on all platforms
- [ ] Help/command output readable on all platforms
- [ ] Test results documented in evidence
- [ ] Sign-off approved by team lead or release manager
- [ ] Ready to announce v0.0.1 publicly

---

## Post-Release

**If issues found:**
- Document in `KNOWN_ISSUES.md`
- Plan v0.0.1-fix1 patch release
- Accelerate v0.1.0 with fixes

**If all pass:**
- Announce v0.0.1 release
- Update installation docs
- Email users: "hawp v0.0.1 is live"

---

## Notes

- **Parallelizable:** Different people can test different platforms simultaneously
- **Blockers:** None (depends on GitHub Actions building successfully)
- **Owner:** Release manager or designated tester
- **Timeline:** Can run same day as repo audit (in parallel)
- **Confidence:** High (builds are automated; testing just validates delivery)


## Outcome

Shipped in the `0.0.1` release (tag `0.0.1`, 2026-08-21). Work complete.

## Verification

Release published at https://github.com/sentzunhat/human-ai-workflow-protocol/releases/tag/0.0.1 with all 7 assets.

## Close Checklist

- [x] Work shipped in 0.0.1 release
- [x] Archived to closed/2026/08/22/v001-shipped-cleanup/
