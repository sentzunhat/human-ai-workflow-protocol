package enable

import "testing"

func TestParseArgsValid(t *testing.T) {
	opts, err := parseArgs(nil)
	if err != nil || opts.logBodies {
		t.Fatalf("default: %v, logBodies=%v", err, opts.logBodies)
	}
	opts, err = parseArgs([]string{"--log-bodies"})
	if err != nil || !opts.logBodies {
		t.Fatalf("--log-bodies: %v, logBodies=%v", err, opts.logBodies)
	}
	opts, err = parseArgs([]string{"--no-update-check"})
	if err != nil || opts.logBodies {
		t.Fatalf("--no-update-check: %v, logBodies=%v", err, opts.logBodies)
	}
}

func TestParseArgsInvalid(t *testing.T) {
	for _, args := range [][]string{
		{"positional"},
		{"--unknown"},
		{"--log-bodies", "extra"},
	} {
		if _, err := parseArgs(args); err == nil {
			t.Errorf("accepted %v", args)
		}
	}
}
