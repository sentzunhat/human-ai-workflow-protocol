package context

import (
	"path/filepath"
	"testing"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/context/source"
)

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
	corpus := source.KitCorpus{Files: []source.File{
		{RelPath: "start-here.md", RepoPath: ".hawp/kit/start-here.md", Content: "# Start\n\nEntry point for HAWP.\n"},
		{RelPath: "usage/init.md", RepoPath: ".hawp/kit/usage/init.md", Content: "# Init\n\nHow to run init.\n"},
		{RelPath: "standards/README.md", RepoPath: ".hawp/kit/standards/README.md", Content: "# Standards\n\nRules to follow in real work.\n"},
		{RelPath: "standards/naming.md", RepoPath: ".hawp/kit/standards/naming.md", Content: "# Naming\n\nUse kebab-case.\n"},
	}}

	docs := EnrichKit(corpus)
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
