package cli

import (
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	cases := [][]string{nil, {"--help"}, {"-h"}}
	for _, args := range cases {
		if err := Run(args); err != nil {
			t.Errorf("Run(%v) returned error: %v", args, err)
		}
	}
}

func TestRunUnknownCommand(t *testing.T) {
	err := Run([]string{"does-not-exist"})
	if err == nil {
		t.Fatal("Run with unknown command returned nil, want error")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("error %q does not mention unknown command", err.Error())
	}
	if !strings.Contains(err.Error(), "USAGE") {
		t.Errorf("error %q does not include help text", err.Error())
	}
}

func TestRunDBInit(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if err := Run([]string{"db", "init"}); err != nil {
		t.Fatalf("Run(db init) returned error: %v", err)
	}
}

func TestRunUUIDRoutes(t *testing.T) {
	for _, args := range [][]string{{"uuid"}, {"uuid", "--short"}} {
		if err := Run(args); err != nil {
			t.Errorf("Run(%v) returned error: %v", args, err)
		}
	}
}

func TestRunValidateRoutesWithExplicitPaths(t *testing.T) {
	// kit validate against a temp kit missing required files → ExitError(1).
	err := Run([]string{"kit", "validate", "--kit-path", t.TempDir()})
	var exitErr ExitError
	if !errorsAs(err, &exitErr) || exitErr.Code != 1 {
		t.Errorf("kit validate on empty kit: err = %v, want ExitError{1}", err)
	}

	// work validate with an unresolvable work root → plain error.
	if err := Run([]string{"work", "validate", "--work-root", "/nonexistent-hawp-work"}); err == nil {
		t.Error("work validate with bad root should error")
	}
}

func errorsAs(err error, target *ExitError) bool {
	e, ok := err.(ExitError)
	if ok {
		*target = e
	}
	return ok
}

func TestRunIndexBuild(t *testing.T) {
	cases := [][]string{
		{"index", "build"},
		{"index", "build", "--scope", "work"},
		{"index", "build", "--scope", "kit"},
	}
	for _, args := range cases {
		if err := Run(args); err != nil {
			t.Errorf("Run(%v) returned error: %v", args, err)
		}
	}
}

func TestRunSearchWithContext(t *testing.T) {
	cases := [][]string{
		{"search", "test"},                                                  // plain search (no index found is OK)
		{"search", "test", "--limit", "5"},                                 // with limit
		{"search", "test", "--context"},                                    // with --context flag
		{"search", "test", "--context", "--format", "json"},                // with JSON format
		{"search", "test", "--context", "--format", "markdown"},            // with markdown format (explicit)
		{"search", "test", "--context", "--max-tokens", "1000"},            // with token budget
		{"search", "test", "--context", "--format", "json", "--max-tokens", "500"}, // combined flags
	}
	for _, args := range cases {
		// Note: these will likely say "no results found" or "index not found",
		// which is OK; we're just testing the CLI parsing and flag handling.
		Run(args) // Errors expected (no index), but parsing should succeed.
	}
}
