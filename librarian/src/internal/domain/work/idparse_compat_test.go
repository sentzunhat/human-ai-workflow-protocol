package work

import "testing"

func TestIdentityCompatibilityWrappers(t *testing.T) {
	if ExtractShortUUID("0E1C4AFA") != "0e1c4afa" {
		t.Fatal("short UUID wrapper changed")
	}
	if !IDsMatch("0e1c4afa", "0e1c4afa-9668-4d61-b5b6-1e27be42ca23") {
		t.Fatal("ID matching wrapper changed")
	}
}
