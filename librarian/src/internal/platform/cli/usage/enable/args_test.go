package enable

import "testing"

func TestParseArgsValid(t *testing.T) {
	opts, err := parseArgs(nil)
	if err != nil || opts.logBodies || opts.logBodiesSet {
		t.Fatalf("default: %v, logBodies=%v, logBodiesSet=%v", err, opts.logBodies, opts.logBodiesSet)
	}
	opts, err = parseArgs([]string{"--log-bodies"})
	if err != nil || !opts.logBodies || !opts.logBodiesSet {
		t.Fatalf("--log-bodies: %v, logBodies=%v, logBodiesSet=%v", err, opts.logBodies, opts.logBodiesSet)
	}
	opts, err = parseArgs([]string{"--no-update-check"})
	if err != nil || opts.logBodies || opts.logBodiesSet {
		t.Fatalf("--no-update-check: %v, logBodies=%v, logBodiesSet=%v", err, opts.logBodies, opts.logBodiesSet)
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
