package newcmd

import "testing"

func TestNewRejectsEmptyRootBeforeWriting(t *testing.T) {
	root := t.TempDir()
	if err := Run([]string{"Should not be created", "--hawp-root="}, root); err == nil {
		t.Fatal("expected empty root to be rejected")
	}
}
