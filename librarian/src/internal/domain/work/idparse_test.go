package work

import "testing"

func TestExtractIDFromFilename(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"361fb08e-6457-4ed5-80bd-76337b6f0e89", "361fb08e-6457-4ed5-80bd-76337b6f0e89"},
		{"361FB08E-6457-4ED5-80BD-76337B6F0E89-title", "361fb08e-6457-4ed5-80bd-76337b6f0e89"},
		{"TASK-012", "TASK-012"},
		{"BUG-063-some-title", "BUG-063"},
		{"2026-04-29-BUG-001-title", "BUG-001"},
		{"2026-04-29-bug-001-title", "BUG-001"},
		{"random-name", ""},
		{"", ""},
	}
	for _, c := range cases {
		if got := ExtractIDFromFilename(c.in); got != c.want {
			t.Errorf("ExtractIDFromFilename(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestExtractShortUUID(t *testing.T) {
	if got := ExtractShortUUID("0E1C4AFA"); got != "0e1c4afa" {
		t.Errorf("whole-cell short UUID: got %q", got)
	}
	for _, bad := range []string{"0e1c4afax", "0e1c4af", "0e1c4afa-extra", "deadbeef9"} {
		if got := ExtractShortUUID(bad); got != "" {
			t.Errorf("ExtractShortUUID(%q) = %q, want empty", bad, got)
		}
	}
}

func TestIDsMatch(t *testing.T) {
	full := "0e1c4afa-9668-4d61-b5b6-1e27be42ca23"
	cases := []struct {
		a, b string
		want bool
	}{
		{"TASK-012", "task-012", true},
		{"0e1c4afa", full, true},
		{full, "0e1c4afa", true},
		{"deadbeef", full, false},
		{"TASK-012", "TASK-013", false},
		{"0e1c4afa", "0e1c4afa", true},
	}
	for _, c := range cases {
		if got := IDsMatch(c.a, c.b); got != c.want {
			t.Errorf("IDsMatch(%q, %q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}
