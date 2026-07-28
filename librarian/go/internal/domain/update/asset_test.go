package update

import (
	"strings"
	"testing"
)

func TestAssetNameMatchesCurrentPlatform(t *testing.T) {
	name, err := AssetName()
	if err != nil {
		// Only unsupported platforms error; the CI/dev matrix (darwin,
		// linux, windows amd64/arm64) always succeeds.
		t.Skipf("platform not supported: %v", err)
	}
	if !strings.HasPrefix(name, "hawp-") {
		t.Errorf("AssetName() = %q, want hawp-* prefix", name)
	}
}
