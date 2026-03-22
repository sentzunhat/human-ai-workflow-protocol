package context

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFixture(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestKitRole(t *testing.T) {
	if got := kitRole("start-here.md"); got != "root" {
		t.Errorf("kitRole(start-here.md) = %q, want root", got)
	}
	if got := kitRole("usage/init.md"); got != "usage" {
		t.Errorf("kitRole(usage/init.md) = %q, want usage", got)
	}
	if got := kitRole("standards/public/README.md"); got != "standards" {
		t.Errorf("kitRole(nested) = %q, want standards", got)
	}
}

func TestEnrichKitAssignsRolesAndFolderContext(t *testing.T) {
	root := t.TempDir()
	kitPath := filepath.Join(root, ".hawp", "kit")
	writeFixture(t, root, map[string]string{
		".hawp/kit/start-here.md":       "# Start\n\nEntry point for HAWP.\n",
		".hawp/kit/usage/init.md":       "# Init\n\nHow to run init.\n",
		".hawp/kit/standards/README.md": "# Standards\n\nRules to follow in real work.\n",
		".hawp/kit/standards/naming.md": "# Naming\n\nUse kebab-case.\n",
	})

	docs, err := EnrichKit(root, kitPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(docs) != 4 {
		t.Fatalf("documents = %d, want 4 (including the standards README)", len(docs))
	}

	byRelSuffix := map[string]Document{}
	for _, doc := range docs {
		byRelSuffix[filepath.Base(doc.RelPath)] = doc
	}

	start := byRelSuffix["start-here.md"]
	if start.Role != "root" || start.ContextPrefix != "[kit/root]" {
		t.Errorf("start-here.md = %+v", start)
	}

	naming := byRelSuffix["naming.md"]
	if naming.Role != "standards" {
		t.Errorf("naming.md role = %q, want standards", naming.Role)
	}
	if naming.ContextPrefix != "[kit/standards] Rules to follow in real work." {
		t.Errorf("naming.md ContextPrefix = %q, want folder README summary included", naming.ContextPrefix)
	}

	usageDoc := byRelSuffix["init.md"]
	if usageDoc.ContextPrefix != "[kit/usage]" {
		t.Errorf("init.md ContextPrefix = %q, want plain tag (no README in usage/)", usageDoc.ContextPrefix)
	}
}

func TestFirstDescriptiveLineSkipsHeadingsAndBlanks(t *testing.T) {
	content := "# Title\n\n\nActual description here.\nMore text.\n"
	if got := firstDescriptiveLine(content); got != "Actual description here." {
		t.Errorf("firstDescriptiveLine = %q", got)
	}
}
