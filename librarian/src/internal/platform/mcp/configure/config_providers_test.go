package mcp

import (
	"os"
	"reflect"
	"testing"
)

func TestConfigProviderExpansion(t *testing.T) {
	for _, tc := range []struct {
		input, want []string
	}{
		{nil, nil},
		{[]string{"all"}, []string{"claude", "cursor", "continue", "codex"}},
		{[]string{"codex", "all", "claude", "all", "github"}, []string{"codex", "claude", "cursor", "continue", "github"}},
		{[]string{"github", "github"}, []string{"github"}},
	} {
		got, err := expandConfigProviders(tc.input)
		if err != nil || !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("%v: got %v, %v; want %v", tc.input, got, err, tc.want)
		}
	}
}

func TestInvalidConfigProviderRefusesBeforeAnyWrite(t *testing.T) {
	for _, invalid := range []string{"", "codxe", "Codex", "../codex"} {
		t.Run(invalid, func(t *testing.T) {
			root := t.TempDir()
			if err := WriteProviderConfigs(root, []string{"claude", invalid}); err == nil {
				t.Fatal("invalid provider accepted")
			}
			entries, err := os.ReadDir(root)
			if err != nil || len(entries) != 0 {
				t.Fatalf("wrote files before rejecting selection: %v, %v", entries, err)
			}
		})
	}
}
