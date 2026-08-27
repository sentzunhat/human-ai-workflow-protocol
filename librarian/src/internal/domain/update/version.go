// Package update implements the self-update client: compare the running
// binary's version against the latest published release, verify the
// platform asset's checksum, and atomically replace the executable.
package update

import (
	"strconv"
	"strings"
)

// Version is the running binary's version, set at build time via
// `-ldflags="-X .../update.Version=vX.Y.Z"` (see librarian/src/Makefile).
// "dev" marks a local, non-release build — update checks still work but
// comparisons against it always report an update available.
var Version = "0.0.14"

// Repo is the GitHub "owner/repo" this binary updates from.
const Repo = "sentzunhat/human-ai-workflow-protocol"

// TagPrefix was prepended to legacy release tags (e.g. "librarian-go-v0.0.1").
// New releases use plain X.Y.Z tags — CleanVersion strips this prefix when
// present and is a no-op when it isn't, so both formats are handled.
const TagPrefix = "librarian-go-"

// CleanVersion strips TagPrefix from a raw release tag for display,
// leaving the "vX.Y.Z" form (matching how Version itself looks) — e.g.
// "librarian-go-v0.0.2" -> "v0.0.2". Real-world bug caught during the
// v0.0.1/v0.0.2 test releases (4c152ee3): ParseVersion alone only
// stripped a leading "v", so the raw tag "librarian-go-v0.0.2" left
// "librarian-go-v0" as its unparseable major-version component — it
// happened to still compare correctly against "v0.0.1" only because the
// prefix junk was confined to the always-zero major slot, and would have
// silently broken on any real major-version bump.
func CleanVersion(tag string) string {
	return strings.TrimPrefix(strings.TrimSpace(tag), TagPrefix)
}

// ParseVersion splits a "vX.Y.Z" (or "X.Y.Z", or a full release tag with
// TagPrefix) into numeric components. Non-numeric or missing parts are
// treated as 0, so "v1.2" < "v1.2.1" and malformed tags degrade to 0.0.0
// rather than erroring.
func ParseVersion(tag string) [3]int {
	clean := strings.TrimPrefix(CleanVersion(tag), "v")
	parts := strings.SplitN(clean, ".", 3)
	var out [3]int
	for i := 0; i < len(parts) && i < 3; i++ {
		n, err := strconv.Atoi(parts[i])
		if err == nil {
			out[i] = n
		}
	}
	return out
}

// IsNewer reports whether candidate is a strictly newer version than
// current. "dev" (unparseable as a release version) is always considered
// older than any real release tag.
func IsNewer(current, candidate string) bool {
	if current == "dev" {
		return true
	}
	c, n := ParseVersion(current), ParseVersion(candidate)
	for i := 0; i < 3; i++ {
		if n[i] != c[i] {
			return n[i] > c[i]
		}
	}
	return false
}
