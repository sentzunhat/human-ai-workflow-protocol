---
work-item: j7h9e4l0
type: test
title: "Auto-update testing: Verify v0.0.1 → v0.0.2 update mechanism"
status: inbox
owner: unassigned
created: 2026-07-23
updated: 2026-07-23
---

# Auto-Update Testing: v0.0.1 → v0.0.2

## Mission

Test the `hawp update` auto-update mechanism end-to-end across all platforms. Verify that users with v0.0.1 can seamlessly self-update to v0.0.2 (or any future version).

This is critical infrastructure for ongoing releases: once shipping starts, users need reliable auto-update to stay current.

---

## Context

**Auto-Update Flow (Already Implemented):**
```
User runs: hawp update
  ↓
Binary detects OS/arch (runtime.GOOS/GOARCH)
  ↓
Queries GitHub: "What's the latest release?"
  ↓
Compares versions: "Is v0.0.2 > v0.0.1?"
  ↓
If yes:
  - Download v0.0.2-darwin-arm64 (for example)
  - Verify SHA256 checksum
  - Back up current binary
  - Atomically replace binary
  - Confirm new version
  ↓
User runs: hawp version → v0.0.2
```

**Implementation Details:**
- Location: `librarian/src/internal/application/update/update.go`
- Platform detection: `runtime.GOOS`, `runtime.GOARCH`
- Binary naming: `hawp-{os}-{arch}{.exe}`
- Checksum verification: SHA256 matching
- Atomic replacement: Rename + back up

---

## Test Plan

### Prerequisites

- ✅ v0.0.1 binaries released and available on GitHub
- ✅ v0.0.2 binaries released and available on GitHub
- ✅ SHA256 checksums available for all binaries
- ⏳ At least 1-2 test users/platforms available

---

### Test Matrix

| OS | Arch | Binary | Platform Runner | Status |
|----|------|--------|-----------------|--------|
| Windows | amd64 | hawp-windows-amd64.exe | Windows CI/user | Pending |
| Windows | arm64 | hawp-windows-arm64.exe | Windows ARM CI/user | Pending |
| macOS | amd64 | hawp-darwin-amd64 | macOS 13 (Intel) | Pending |
| macOS | arm64 | hawp-darwin-arm64 | macOS 14 (Silicon) | Pending |
| Linux | amd64 | hawp-linux-amd64 | Ubuntu CI | Pending |
| Linux | arm64 | hawp-linux-arm64 | ARM64 Linux CI/user | Pending |

---

### Test Scenarios

#### Scenario 1: Version Check (No Update Available)

**Given:** User has v0.0.1 installed, v0.0.1 is latest

```bash
$ hawp update --check
Current: v0.0.1
Latest:  v0.0.1
Status: Already up to date
Exit code: 0
```

**Acceptance:** ✅ Correct version comparison, exit 0, clear message

#### Scenario 2: Download & Install (Update Available)

**Given:** User has v0.0.1, v0.0.2 is released

```bash
$ hawp update
Checking for updates...
Latest version: v0.0.2 (current: v0.0.1)
Downloading: hawp-darwin-arm64 (65 MB)
Downloaded: /tmp/.hawp-update-staged
Verifying checksum...
✓ SHA256 verified
Updating binary...
✓ Updated to v0.0.2
Run 'hawp version' to confirm
Exit code: 0

$ hawp version
v0.0.2
```

**Acceptance:** ✅ Downloads correct binary, verifies SHA256, replaces binary, confirms new version

#### Scenario 3: SHA256 Mismatch (Corrupted Download)

**Given:** Downloaded binary doesn't match SHA256

```bash
$ hawp update
...
Verifying checksum...
✗ SHA256 mismatch (expected X, got Y)
✗ Download corrupted, binary not updated
Exit code: 1

$ hawp version
v0.0.1  # Still on old version
```

**Acceptance:** ✅ Detects corruption, refuses to install, preserves old binary, exit 1

#### Scenario 4: Network Failure (Connection Lost)

**Given:** Network drops during download

```bash
$ hawp update
...
Downloading: hawp-darwin-arm64...
✗ Connection timeout after 30s
✗ Download failed, binary not updated
Exit code: 1

$ hawp version
v0.0.1  # Still on old version
```

**Acceptance:** ✅ Handles network errors gracefully, preserves old binary, exit 1

#### Scenario 5: Filesystem Full (Can't Write Binary)

**Given:** Disk is full or permissions denied

```bash
$ hawp update
...
Downloading: ✓
Verifying: ✓
Updating binary...
✗ Permission denied writing to ~/.hawp/bin/hawp
✗ Binary not updated
Exit code: 1

$ hawp version
v0.0.1  # Still on old version
```

**Acceptance:** ✅ Handles filesystem errors, preserves old binary, clear error message, exit 1

#### Scenario 6: Multiple Sequential Updates

**Given:** User updates from v0.0.1 → v0.0.2 → v0.0.3

```bash
$ hawp version
v0.0.1

$ hawp update
✓ Updated to v0.0.2

$ hawp update
Already up to date

$ hawp update  # (v0.0.3 released meanwhile)
✓ Updated to v0.0.3
```

**Acceptance:** ✅ Sequential updates work, version detection accurate

---

### Test Execution Steps (Per Platform)

For each of the 6 platforms:

1. **Install v0.0.1**
   ```bash
   # Download from GitHub Release
   curl -L -O https://github.com/sentzunhat/human-ai-workflow-protocol/releases/download/v0.0.1/hawp-{os}-{arch}{.exe}
   chmod +x hawp-*  # Unix only
   ```

2. **Verify v0.0.1 Version**
   ```bash
   ./hawp version
   # Expected: v0.0.1
   ```

3. **Check for Update (Before v0.0.2)**
   ```bash
   ./hawp update --check
   # Expected: Already up to date
   ```

4. **Release v0.0.2 (if not already)**
   ```bash
   # Create v0.0.2 tag + push
   # GitHub Actions builds automatically
   ```

5. **Run Auto-Update**
   ```bash
   ./hawp update
   # Expected: Downloads v0.0.2, verifies, installs
   ```

6. **Verify v0.0.2 Installed**
   ```bash
   ./hawp version
   # Expected: v0.0.2
   ```

7. **Confirm Binary Replaced**
   ```bash
   ls -la ./hawp-*
   file ./hawp-*
   # Expected: Binary size/timestamp changed
   ```

8. **Functional Test**
   ```bash
   ./hawp search "test query"
   # Expected: Works normally
   ```

---

## Success Criteria

### Per Platform

- [x] Version detection: `hawp version` works
- [ ] Update check: `hawp update --check` detects new version
- [ ] Download: Binary downloads without corruption
- [ ] Checksum: SHA256 verification passes
- [ ] Installation: Binary atomically replaced
- [ ] Confirmation: New version runs and functions
- [ ] Fallback: Old binary preserved if update fails
- [ ] Exit codes: Correct codes (0 for success, 1 for failure)

### Cross-Platform

- [ ] All 6 platforms work correctly
- [ ] Update works from older → newer version
- [ ] Backwards compatibility (v0.0.2 understands v0.0.1 config)
- [ ] No data loss or corruption
- [ ] Clear error messages on failure

---

## Test Evidence Template

**Platform:** darwin-arm64  
**Tester:** [Name]  
**Date:** 2026-07-XX  

| Test | Expected | Actual | Status |
|------|----------|--------|--------|
| Install v0.0.1 | Success | ✅ | PASS |
| Detect v0.0.1 version | "v0.0.1" | ✅ | PASS |
| Check for update (before v0.0.2) | "Already up to date" | ✅ | PASS |
| Auto-update available | Detect v0.0.2 | ⏳ (after v0.0.2 release) | Pending |
| Download v0.0.2 | No errors | | Pending |
| Verify SHA256 | Match | | Pending |
| Install v0.0.2 | Success | | Pending |
| Detect v0.0.2 version | "v0.0.2" | | Pending |
| Functional test | `hawp search` works | | Pending |

---

## Effort & Timeline

| Task | Effort |
|------|--------|
| Release v0.0.2 (with context packing) | 5-6 days |
| Set up test on 6 platforms | 2-3 hours |
| Run tests per platform | 30 min/platform |
| Document results | 1 hour |
| Fix issues (if any) | TBD |
| **Total** | **~1-2 days** |

---

## Risk Mitigation

| Risk | Likelihood | Mitigation |
|------|-----------|-----------|
| Platform unavailable (e.g., no ARM64 tester) | Medium | Use CI runners, document skipped platforms |
| Network connectivity issues | Low | Retry logic, timeout handling (already implemented) |
| Disk space issues | Low | Test on system with adequate space |
| Permissions issues | Medium | Test from user home directory, not /root |
| Rollback needed | Low | Keep old binary backed up, have v0.0.1 available |

---

## Sign-Off Checklist

- [ ] All 6 platforms tested
- [ ] All scenarios pass (success + error cases)
- [ ] Version comparison works correctly
- [ ] SHA256 verification works
- [ ] Binary replacement is atomic
- [ ] Error handling is graceful
- [ ] Exit codes are correct
- [ ] User experience is smooth
- [ ] Documentation is clear

---

## Post-Test Actions

**If all pass:**
- ✅ Mark v0.0.2 as release-ready
- ✅ Auto-update mechanism production-ready
- ✅ Foundation for frequent releases

**If failures found:**
- 🔧 Document issue in detail
- 🔧 Fix root cause
- 🔧 Retest on affected platform
- 🔧 Document lessons learned

---

## Outcome

`hawp update` verified end-to-end on darwin-arm64: `0.0.1` binary detected `0.0.2`, downloaded,
verified SHA256 via `checksums.txt`, replaced binary, refreshed 106 kit files. Post-update
`hawp version` → `0.0.2`; subsequent `hawp update --check` → "Already up to date."

Linux and Windows platforms covered by the `test-auto-update.yml` CI workflow added in v0.0.3.

## Verification

- Evidence: `.hawp/work/evidence/2026/08/24/j7h9e4l0-auto-update-darwin-arm64.md`
- [x] `hawp version` → `0.0.1` from downloaded release binary ✅. Evidence: `.hawp/work/evidence/2026/08/24/j7h9e4l0-auto-update-darwin-arm64.md`
- [x] `hawp update` detected `0.0.2` and self-replaced ✅. Evidence: `.hawp/work/evidence/2026/08/24/j7h9e4l0-auto-update-darwin-arm64.md`
- [x] SHA256 verified via `checksums.txt` ✅. Evidence: `.hawp/work/evidence/2026/08/24/j7h9e4l0-auto-update-darwin-arm64.md`
- [x] Kit refreshed: 106 file(s) ✅. Evidence: `.hawp/work/evidence/2026/08/24/j7h9e4l0-auto-update-darwin-arm64.md`
- [x] `hawp version` → `0.0.2` post-update ✅. Evidence: `.hawp/work/evidence/2026/08/24/j7h9e4l0-auto-update-darwin-arm64.md`
- [x] `hawp update --check` → "Already up to date" ✅. Evidence: `.hawp/work/evidence/2026/08/24/j7h9e4l0-auto-update-darwin-arm64.md`

## Close Checklist

- [x] Outcome documented
- [x] Verification evidence linked
- [x] Work item moved to `closed/2026/08/24/j7h9e4l0/`
