package providersync

import (
	"os"
	"path/filepath"
	"testing"

	domainprovidersync "github.com/sentzunhat/hawp/librarian/src/internal/domain/providersync"
)

func TestMaterializeRejectsSymlinkedProviderOutputAncestor(t *testing.T) {
	root := t.TempDir()
	shared := filepath.Join(root, "core", "providers", "shared", "behaviors")
	if err := os.MkdirAll(shared, 0o755); err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, target := range domainprovidersync.MaterializationTargets {
		if seen[target.Behavior] {
			continue
		}
		seen[target.Behavior] = true
		if err := os.WriteFile(filepath.Join(shared, target.Behavior+".md"), []byte("# "+target.Behavior+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	outside := t.TempDir()
	providerRoot := filepath.Join(root, "core", "providers")
	if err := os.Symlink(outside, filepath.Join(providerRoot, ".github")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if _, err := Materialize(root); err == nil {
		t.Fatal("Materialize accepted symlinked provider output ancestor")
	}
	entries, err := os.ReadDir(outside)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("Materialize wrote outside repository: %v", entries)
	}
}
