package clean

import (
	"regexp"
	"testing"
)

var spaceRe = regexp.MustCompile(`\s+`)

func TestParseArgsValid(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		wantApply bool
	}{
		{"no args", nil, false},
		{"apply", []string{"--apply"}, true},
		{"no-update-check", []string{"--no-update-check"}, false},
		{"all flags", []string{"--apply", "--no-update-check"}, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			opts, err := parseArgs(c.args)
			if err != nil {
				t.Fatalf("parseArgs(%v) returned unexpected error: %v", c.args, err)
			}
			if opts.apply != c.wantApply {
				t.Errorf("apply = %v, want %v", opts.apply, c.wantApply)
			}
		})
	}
}

func TestParseArgsInvalid(t *testing.T) {
	cases := []string{
		"--bogus-flag",
		"links clean --apply extra",
		"--bogus --bogus2",
	}

	for _, tc := range cases {
		t.Run(tc, func(t *testing.T) {
			args := spaceRe.Split(tc, -1)
			if _, err := parseArgs(args); err == nil {
				t.Errorf("parseArgs(%v) expected error, got nil", args)
			}
		})
	}
}
