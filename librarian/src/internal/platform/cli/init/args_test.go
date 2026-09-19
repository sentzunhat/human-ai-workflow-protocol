package init

import "testing"

func TestParseInitArgsValid(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		wantProvs []string
	}{
		{name: "no args", args: nil, wantProvs: nil},
		{name: "single provider", args: []string{"--provider", "claude"}, wantProvs: []string{"claude"}},
		{name: "repeated provider", args: []string{"--provider", "claude", "--provider", "cursor"}, wantProvs: []string{"claude", "cursor"}},
		{name: "provider all", args: []string{"--provider", "all"}, wantProvs: []string{"all"}},
		{name: "no-update-check suppressed", args: []string{"--no-update-check"}, wantProvs: nil},
		{name: "provider and no-update-check", args: []string{"--provider", "codex", "--no-update-check"}, wantProvs: []string{"codex"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseInitArgs(c.args)
			if err != nil {
				t.Fatalf("parseInitArgs(%v) unexpected error: %v", c.args, err)
			}
			if len(got.providers) != len(c.wantProvs) {
				t.Fatalf("providers = %v, want %v", got.providers, c.wantProvs)
			}
			for i, p := range got.providers {
				if p != c.wantProvs[i] {
					t.Errorf("providers[%d] = %q, want %q", i, p, c.wantProvs[i])
				}
			}
		})
	}
}

func TestParseInitArgsInvalid(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{name: "unknown flag", args: []string{"--bogus"}},
		{name: "provider flag with no value", args: []string{"--provider"}},
		{name: "extra positional argument", args: []string{"somepath"}},
		{name: "provider then extra positional", args: []string{"--provider", "claude", "extra"}},
		{name: "unknown flag after provider", args: []string{"--provider", "claude", "--bogus"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := parseInitArgs(c.args); err == nil {
				t.Errorf("parseInitArgs(%v) expected an error, got nil", c.args)
			}
		})
	}
}
