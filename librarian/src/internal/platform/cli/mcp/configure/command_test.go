package configure

import (
	"testing"
)

func TestMCPConfigureRejectsInvalidArguments(t *testing.T) {
	cwd := t.TempDir()
	for _, args := range [][]string{
		{"--provider"}, {"--provider="}, {"--repo-root="},
		{"--repo-root"}, {"--unknown"}, {"extra"},
	} {
		if err := Run(args, cwd); err == nil {
			t.Errorf("accepted %v", args)
		}
	}
}
