package normalize

import (
	"path/filepath"
	"testing"
)

func TestResolveKitPath(t *testing.T) {
	root := t.TempDir()
	cwd := filepath.Join(root, "nested")

	if got, want := resolveKitPath(root, cwd, ""), filepath.Join(root, ".hawp", "kit"); got != want {
		t.Fatalf("default path = %q, want %q", got, want)
	}
	if got, want := resolveKitPath(root, cwd, "../kit"), filepath.Join(root, "kit"); got != want {
		t.Fatalf("relative path = %q, want %q", got, want)
	}

	absolute := filepath.Join(root, "custom", "kit")
	if got := resolveKitPath(root, cwd, absolute); got != absolute {
		t.Fatalf("absolute path = %q, want %q", got, absolute)
	}
}
