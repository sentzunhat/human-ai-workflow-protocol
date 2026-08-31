package update_test

import (
	"strings"
	"testing"

	domainupdate "github.com/sentzunhat/hawp/librarian/src/internal/domain/update"
)

func TestParseVersion(t *testing.T) {
	cases := []struct {
		in   string
		want [3]int
	}{
		{"v1.2.3", [3]int{1, 2, 3}},
		{"1.2.3", [3]int{1, 2, 3}},
		{"v1.2", [3]int{1, 2, 0}},
		{"v1", [3]int{1, 0, 0}},
		{"dev", [3]int{0, 0, 0}},
		{"", [3]int{0, 0, 0}},
		{"v1.x.3", [3]int{1, 0, 3}},
	}
	for _, c := range cases {
		if got := domainupdate.ParseVersion(c.in); got != c.want {
			t.Errorf("ParseVersion(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestIsNewer(t *testing.T) {
	cases := []struct {
		current, candidate string
		want               bool
	}{
		{"v1.0.0", "v1.0.1", true},
		{"v1.0.0", "v1.1.0", true},
		{"v1.0.0", "v2.0.0", true},
		{"v1.0.1", "v1.0.0", false},
		{"v1.0.0", "v1.0.0", false},
		{"dev", "v0.0.1", true},
		{"dev", "v1.0.0", true},
		{"v0.0.1", "librarian-go-v0.0.2", true},
		{"v1.0.0", "librarian-go-v2.0.0", true},
		{"v2.0.0", "librarian-go-v1.9.9", false},
	}
	for _, c := range cases {
		if got := domainupdate.IsNewer(c.current, c.candidate); got != c.want {
			t.Errorf("IsNewer(%q, %q) = %v, want %v", c.current, c.candidate, got, c.want)
		}
	}
}

func TestCleanVersion(t *testing.T) {
	cases := []struct{ in, want string }{
		{"librarian-go-v0.0.2", "v0.0.2"},
		{"v0.0.2", "v0.0.2"},
		{"0.0.2", "0.0.2"},
		{"dev", "dev"},
	}
	for _, c := range cases {
		if got := domainupdate.CleanVersion(c.in); got != c.want {
			t.Errorf("CleanVersion(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestParseVersionStripsTagPrefix(t *testing.T) {
	if got := domainupdate.ParseVersion("librarian-go-v2.1.3"); got != [3]int{2, 1, 3} {
		t.Errorf("ParseVersion(with TagPrefix) = %v, want [2 1 3]", got)
	}
}

func TestAssetNameMatchesCurrentPlatform(t *testing.T) {
	name, err := domainupdate.AssetName()
	if err != nil {
		t.Skipf("platform not supported: %v", err)
	}
	if !strings.HasPrefix(name, "hawp-") {
		t.Errorf("AssetName() = %q, want hawp-* prefix", name)
	}
}
