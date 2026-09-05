package tomlconfig

import (
	"bytes"
	"strings"
	"testing"
)

func TestRewritePreservesAllUnmanagedBytes(t *testing.T) {
	input := []byte("# header\n[target] # table\nvalue = 'old' # explanation\npolicy = 10_000 # formatting\n[next]\nvalue = 'untouched'\n")
	out, err := Rewrite(input, []string{"target"}, map[string]any{"value": "new", "added": 1})
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"# header\n[target] # table\n", "# explanation\npolicy = 10_000 # formatting\n[next]\nvalue = 'untouched'\n"} {
		if !strings.Contains(string(out), expected) {
			t.Fatalf("changed unrelated bytes: %s", out)
		}
	}
	again, err := Rewrite(out, []string{"target"}, map[string]any{"value": "new", "added": 1})
	if err != nil || !bytes.Equal(out, again) {
		t.Fatalf("not idempotent: %v", err)
	}
	if string(input) != "# header\n[target] # table\nvalue = 'old' # explanation\npolicy = 10_000 # formatting\n[next]\nvalue = 'untouched'\n" {
		t.Fatal("mutated input")
	}
}

func TestRewritePreservesEqualValuesIncludingComments(t *testing.T) {
	input := []byte("[target]\nargs = [ # keep\n 'mcp',\n]\n")
	out, err := Rewrite(input, []string{"target"}, map[string]any{"args": []string{"mcp"}})
	if err != nil || !bytes.Equal(input, out) {
		t.Fatalf("changed equal commented array: %s, %v", out, err)
	}
}
