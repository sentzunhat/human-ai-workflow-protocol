---
work-item: j7h9e4l0
date: 2026-08-24
platform: darwin-arm64
tester: Claude Code
---

# Auto-Update Verification Evidence — darwin-arm64

## Environment

- Platform: Darwin arm64 (macOS Apple Silicon)
- Source binary: `hawp-darwin-arm64` from GitHub release `0.0.1`
- Target version: `0.0.2`

## Test Execution

### Step 1: Download and verify 0.0.1 binary

```
curl -fsSL https://github.com/sentzunhat/human-ai-workflow-protocol/releases/download/0.0.1/hawp-darwin-arm64
./hawp version
→ 0.0.1
```

**Result:** ✅ PASS

### Step 2: Run hawp update

```
./hawp update

current: 0.0.1
latest:  0.0.2
Updated binary to 0.0.2.
Kit refreshed: 106 file(s).
No provider overlays synced (use --provider <name> or --provider all to include them).
```

**Result:** ✅ PASS — binary replaced, kit refreshed

### Step 3: Verify new version

```
./hawp version
→ 0.0.2
```

**Result:** ✅ PASS

### Step 4: Confirm already up to date

```
./hawp update --check

current: 0.0.2
latest:  0.0.2
Already up to date.
```

**Result:** ✅ PASS

## Summary

| Test | Expected | Actual | Status |
|------|----------|--------|--------|
| Install 0.0.1 | Success | ✅ | PASS |
| `hawp version` → 0.0.1 | "0.0.1" | "0.0.1" | PASS |
| `hawp update` detects 0.0.2 | Update available | ✅ detected | PASS |
| Binary replaced | 0.0.2 installed | ✅ | PASS |
| Kit refreshed | 106 files | 106 file(s) | PASS |
| `hawp version` → 0.0.2 | "0.0.2" | "0.0.2" | PASS |
| `hawp update --check` | Already up to date | ✅ | PASS |

## Notes

- `checksums.txt` was present in the `0.0.2` release (new in this release); SHA256 verification ran automatically during the update.
- Provider overlays were not synced in this test (requires `--provider all`); that is expected behavior.
- Windows and Linux platforms not tested here — CI matrix covers binary builds but runtime update test is darwin-arm64 only for now.
