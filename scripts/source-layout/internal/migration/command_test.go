package migration

import (
	"errors"
	"flag"
	"strings"
	"testing"
)

// Exercise parsing only: these inputs exit before preparation or application.
func TestCommandFlags(t *testing.T) {
	for _, args := range [][]string{
		{"--preview", "--apply"}, {"--diff", "--apply"},
		{"--report", "preview.md", "--apply"},
		{"--apply"}, {"--unknown"}, {"unexpected"},
	} {
		if err := Run(args); err == nil {
			t.Errorf("accepted invalid flags: %v", args)
		}
	}
	if err := Run([]string{"--help"}); !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("help: %v", err)
	}
}

func TestObsoleteMappingRejected(t *testing.T) {
	// Revision checks must happen before reading the supplied filesystem root.
	_, err := currentState("/does-not-exist", plan{Version: 1, Module: modulePath})
	if err == nil || !strings.Contains(err.Error(), "older mapping") {
		t.Fatalf("old plan accepted: %v", err)
	}
}
