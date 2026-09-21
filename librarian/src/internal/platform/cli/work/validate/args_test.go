package validatecmd

import (
	"path/filepath"
	"testing"
)

func TestValidateArgs(t *testing.T) {
	for _, tc := range []struct {
		args []string
		want string
	}{{nil, ""}, {[]string{"--no-update-check"}, ""}, {[]string{"--work-root", "work with spaces"}, "work with spaces"}, {[]string{"--hawp-root=.hawp"}, filepath.Join(".hawp", "work")}} {
		got, err := parseValidateArgs(tc.args)
		if err != nil || got != tc.want {
			t.Fatalf("%v: got %q, %v; want %q", tc.args, got, err, tc.want)
		}
	}
	for _, args := range [][]string{{"--repo-root", "."}, {"--unknown"}, {"extra"}, {"--work-root"}, {"--hawp-root="}, {"--work-root=one", "--hawp-root=two"}} {
		if _, err := parseValidateArgs(args); err == nil {
			t.Errorf("parser accepted %v", args)
		}
		if err := Run(args, t.TempDir()); err == nil {
			t.Errorf("command accepted %v", args)
		}
	}
}
