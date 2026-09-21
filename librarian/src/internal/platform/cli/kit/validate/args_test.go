package validate

import (
	"regexp"
	"testing"
)

func TestParseKitValidateArgsValid(t *testing.T) {
	cases := []struct {
		name        string
		args        []string
		wantKitPath string
	}{
		{"no args", nil, ""},
		{"with kit-path", []string{"--kit-path", "/custom/kit"}, "/custom/kit"},
		{"with no-update-check", []string{"--no-update-check"}, ""},
		{"all flags", []string{"--kit-path", "/opt/hawp/kit", "--no-update-check"}, "/opt/hawp/kit"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			opts, err := parseKitValidateArgs(c.args)
			if err != nil {
				t.Fatalf("parseKitValidateArgs(%v) returned unexpected error: %v", c.args, err)
			}
			if opts.kitPath != c.wantKitPath {
				t.Errorf("kitPath = %q, want %q", opts.kitPath, c.wantKitPath)
			}
		})
	}
}

func TestParseKitValidateArgsInvalid(t *testing.T) {
	cases := []string{
		"unrecognized flag",
		"--kit-path /foo --bogus-flag",
		"trailing argument",
		"kit validate extra arg",
	}

	for _, tc := range cases {
		t.Run(tc, func(t *testing.T) {
			args := splitArgs(tc)
			if _, err := parseKitValidateArgs(args); err == nil {
				t.Errorf("parseKitValidateArgs(%v) expected error, got nil", args)
			}
		})
	}
}

func splitArgs(s string) []string {
	var result []string
	for _, part := range spaceSplitRe.FindAllString(s, -1) {
		result = append(result, part)
	}
	return result
}

var spaceSplitRe = regexp.MustCompile(`\s+`)
