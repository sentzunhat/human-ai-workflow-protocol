package kit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/kit/source"
)

func writeKitFile(t *testing.T, root, rel, content string) string {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestAdapterReadsEntriesAndLinks(t *testing.T) {
	root := t.TempDir()
	writeKitFile(t, root, "usage/guide.md", "[ok](../start-here.md)\n```\n[ignored](nope.md)\n```")
	writeKitFile(t, root, "start-here.md", "# Start")

	snapshot, err := NewAdapter().Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Entries) != 3 || len(snapshot.Files) != 2 {
		t.Fatalf("entries=%d files=%d, want 3/2", len(snapshot.Entries), len(snapshot.Files))
	}
	if len(snapshot.Files[0].Links) == 0 && len(snapshot.Files[1].Links) == 0 {
		t.Fatal("expected parsed non-fenced link")
	}
}

func TestAdapterAppliesMutationsAndRefusesConflict(t *testing.T) {
	root := t.TempDir()
	oldPath := writeKitFile(t, root, "Old Doc.md", "see [guide](Guide.md)")
	guidePath := writeKitFile(t, root, "Guide.md", "# guide")
	adapter := NewAdapter()

	conflictFrom, conflictTo, err := adapter.ApplyRenames([]source.Rename{{From: oldPath, To: filepath.Join(root, "old-doc.md")}})
	if err != nil || conflictFrom != "" || conflictTo != "" {
		t.Fatalf("rename failed: %v conflict=%s->%s", err, conflictFrom, conflictTo)
	}
	changed, err := adapter.ApplyLinkUpdates([]source.LinkUpdate{{File: filepath.Join(root, "old-doc.md"), From: "Guide.md", To: "guide.md", Start: 12, End: 20}})
	if err != nil || changed != 1 {
		t.Fatalf("link update changed=%d err=%v", changed, err)
	}
	raw, err := os.ReadFile(filepath.Join(root, "old-doc.md"))
	if err != nil || !strings.Contains(string(raw), "(guide.md)") {
		t.Fatalf("updated content=%q err=%v", raw, err)
	}

	conflictFrom, conflictTo, err = adapter.ApplyRenames([]source.Rename{{From: guidePath, To: filepath.Join(root, "old-doc.md")}})
	if err != nil || conflictFrom == "" || conflictTo == "" {
		t.Fatalf("expected conflict, got err=%v conflict=%s->%s", err, conflictFrom, conflictTo)
	}
}
