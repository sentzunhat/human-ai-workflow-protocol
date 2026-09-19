package identity

import "testing"

func TestExtractIDFromFilename(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{"361FB08E-6457-4ED5-80BD-76337B6F0E89-title", "361fb08e-6457-4ed5-80bd-76337b6f0e89"},
		{"TASK-012", "TASK-012"}, {"2026-04-29-bug-001-title", "BUG-001"},
		{"042", "042"}, {"025-analysis", ""}, {"random-name", ""},
	} {
		if got := ExtractIDFromFilename(tc.input); got != tc.want {
			t.Errorf("ExtractIDFromFilename(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestIDsMatchAndClassifiers(t *testing.T) {
	full := "0e1c4afa-9668-4d61-b5b6-1e27be42ca23"
	if !IDsMatch("0e1c4afa", full) || IDsMatch("deadbeef", full) {
		t.Fatal("short and full UUID matching is inconsistent")
	}
	if ExtractShortUUID("0E1C4AFA") != "0e1c4afa" || !IsFullUUID(full) || !IsNumericID("042") {
		t.Fatal("identity classifiers returned unexpected results")
	}
}
