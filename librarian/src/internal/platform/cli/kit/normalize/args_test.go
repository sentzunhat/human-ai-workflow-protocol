package normalize

import (
	"regexp"
	"testing"
)

var spaceRe = regexp.MustCompile(`\s+`)

func TestParseKitNormalizeArgsValid(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		wantApply bool
		wantPath  string
	}{
		{"no args", nil, false, ""},
		{"apply", []string{"--apply"}, true, ""},
		{"dry-run", []string{"--dry-run"}, false, ""},
		{"kit-path", []string{"--kit-path", "/custom/kit"}, false, "/custom/kit"},
		{"all flags", []string{"--apply", "--kit-path", "/opt/kit"}, true, "/opt/kit"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			opts, err := parseKitNormalizeArgs(c.args)
			if err != nil {
				t.Fatalf("parseKitNormalizeArgs(%v) returned unexpected error: %v", c.args, err)
			}
			if opts.apply != c.wantApply {
				t.Errorf("apply = %v, want %v", opts.apply, c.wantApply)
			}
			if opts.kitPath != c.wantPath {
				t.Errorf("kitPath = %q, want %q", opts.kitPath, c.wantPath)
			}
		})
	}
}

func TestParseKitNormalizeArgsInvalid(t *testing.T) {
	cases := []string{
		"--bogus-flag",
		"--kit-path --apply extra",
		"normalize extra arg",
	}

	for _, tc := range cases {
		t.Run(tc, func(t *testing.T) {
			args := spaceRe.Split(tc, -1)
			if _, err := parseKitNormalizeArgs(args); err == nil {
				t.Errorf("parseKitNormalizeArgs(%v) expected error, got nil", args)
			}
		})
	}
}

func TestKitNormalizeExplicitFalseModes(t *testing.T) {
	for _, tc := range []struct {
		args  []string
		apply bool
	}{
		{[]string{"--apply=false", "--dry-run"}, false},
		{[]string{"--apply", "--dry-run=false"}, true},
		{[]string{"--dry-run=false"}, false},
	} {
		got, err := parseKitNormalizeArgs(tc.args)
		if err != nil || got.apply != tc.apply {
			t.Errorf("%q: %+v, %v", tc.args, got, err)
		}
	}
}
