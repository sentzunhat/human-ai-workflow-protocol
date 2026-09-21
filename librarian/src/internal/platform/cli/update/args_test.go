package update

import "testing"

func TestParseUpdateSyncArgsValid(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		wantProvs []string
	}{
		{name: "no args — kit-only sync", args: nil, wantProvs: nil},
		{name: "single provider", args: []string{"--provider", "cursor"}, wantProvs: []string{"cursor"}},
		{name: "provider all", args: []string{"--provider", "all"}, wantProvs: []string{"all"}},
		{name: "repeated providers", args: []string{"--provider", "claude", "--provider", "codex"}, wantProvs: []string{"claude", "codex"}},
		{name: "no-update-check only", args: []string{"--no-update-check"}, wantProvs: nil},
		{name: "provider and no-update-check", args: []string{"--provider", "claude", "--no-update-check"}, wantProvs: []string{"claude"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseUpdateSyncArgs(c.args)
			if err != nil {
				t.Fatalf("parseUpdateSyncArgs(%v) unexpected error: %v", c.args, err)
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

func TestParseUpdateSyncArgsInvalid(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{name: "unknown flag", args: []string{"--bogus"}},
		{name: "provider flag with no value", args: []string{"--provider"}},
		{name: "extra positional argument", args: []string{"somepath"}},
		{name: "provider then extra positional", args: []string{"--provider", "all", "extra"}},
		{name: "unknown flag mixed in", args: []string{"--provider", "claude", "--no-such-flag"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := parseUpdateSyncArgs(c.args); err == nil {
				t.Errorf("parseUpdateSyncArgs(%v) expected an error, got nil", c.args)
			}
		})
	}
}

func TestParseUpdateFullArgsValid(t *testing.T) {
	cases := []struct {
		name          string
		args          []string
		wantProvs     []string
		wantNoProvs   bool
		wantNoUpdate  bool
	}{
		{name: "no args — defaults", args: nil, wantProvs: nil, wantNoProvs: false, wantNoUpdate: false},
		{name: "single provider", args: []string{"--provider", "claude"}, wantProvs: []string{"claude"}, wantNoProvs: false, wantNoUpdate: false},
		{name: "no-providers flag", args: []string{"--no-providers"}, wantProvs: nil, wantNoProvs: true, wantNoUpdate: false},
		{name: "provider all with no-providers", args: []string{"--provider", "all", "--no-providers"}, wantProvs: []string{"all"}, wantNoProvs: true, wantNoUpdate: false},
		{name: "no-update-check only", args: []string{"--no-update-check"}, wantProvs: nil, wantNoProvs: false, wantNoUpdate: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseUpdateFullArgs(c.args)
			if err != nil {
				t.Fatalf("parseUpdateFullArgs(%v) unexpected error: %v", c.args, err)
			}
			if len(got.providers) != len(c.wantProvs) {
				t.Fatalf("providers = %v, want %v", got.providers, c.wantProvs)
			}
			for i, p := range got.providers {
				if p != c.wantProvs[i] {
					t.Errorf("providers[%d] = %q, want %q", i, p, c.wantProvs[i])
				}
			}
			if got.noProviders != c.wantNoProvs {
				t.Errorf("noProviders = %v, want %v", got.noProviders, c.wantNoProvs)
			}
			if got.noUpdateCheck != c.wantNoUpdate {
				t.Errorf("noUpdateCheck = %v, want %v", got.noUpdateCheck, c.wantNoUpdate)
			}
		})
	}
}

func TestParseUpdateFullArgsInvalid(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{name: "unknown flag", args: []string{"--bogus"}},
		{name: "trailing provider with no value", args: []string{"--provider"}},
		{name: "extra positional argument", args: []string{"somepath"}},
		{name: "provider then extra positional", args: []string{"--provider", "all", "extra"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := parseUpdateFullArgs(c.args); err == nil {
				t.Errorf("parseUpdateFullArgs(%v) expected an error, got nil", c.args)
			}
		})
	}
}
