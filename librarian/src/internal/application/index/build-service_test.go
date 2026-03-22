package index

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	domaincontext "github.com/sentzunhat/hawp/librarian/src/internal/domain/context"
	domainindex "github.com/sentzunhat/hawp/librarian/src/internal/domain/index"
)

func buildFixtureRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		".hawp/kit/start-here.md": "# Start\n\nEntry point.\n",
		".hawp/kit/usage/init.md": "# Init\n\nHow to init.\n",
		".hawp/work/BACKLOG.md": "# Backlog\n\n## Active Work\n\n" +
			"| ID | Type | Title | Status | Plan File | Updated |\n| --- | --- | --- | --- | --- | --- |\n" +
			"| TASK-001 | feature | thing | in-progress | [plan](active/TASK-001.md) | 2026-07-21 |\n" +
			"\n## Blocked / Parked\n\n| ID | Type | Title | Reason | Detail | Updated |\n| --- | --- | --- | --- | --- | --- |\n" +
			"\n## Recently Closed\n\n| ID | Type | Title | Closed | Detail |\n| --- | --- | --- | --- | --- |\n",
		".hawp/work/active/TASK-001.md": "# plan\n\ncontent\n",
	}
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestExecuteScopeFiltering(t *testing.T) {
	root := buildFixtureRepo(t)
	service := NewBuildService(root)

	kitOnly, err := service.Execute(domainindex.ScopeKit)
	if err != nil {
		t.Fatal(err)
	}
	for _, doc := range kitOnly.Documents {
		if doc.Corpus != domaincontext.CorpusKit {
			t.Errorf("--scope kit returned a %s document", doc.Corpus)
		}
	}
	if len(kitOnly.Documents) != 2 {
		t.Fatalf("kit documents = %d, want 2", len(kitOnly.Documents))
	}

	workOnly, err := service.Execute(domainindex.ScopeWork)
	if err != nil {
		t.Fatal(err)
	}
	for _, doc := range workOnly.Documents {
		if doc.Corpus != domaincontext.CorpusWork {
			t.Errorf("--scope work returned a %s document", doc.Corpus)
		}
	}
	if len(workOnly.Documents) != 1 {
		t.Fatalf("work documents = %d, want 1", len(workOnly.Documents))
	}

	all, err := service.Execute(domainindex.ScopeAll)
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Documents) != len(kitOnly.Documents)+len(workOnly.Documents) {
		t.Errorf("scope all documents = %d, want sum of kit+work", len(all.Documents))
	}
}

func TestStringReportsCountsWithoutDumpingContent(t *testing.T) {
	root := buildFixtureRepo(t)
	service := NewBuildService(root)
	result, err := service.Execute(domainindex.ScopeAll)
	if err != nil {
		t.Fatal(err)
	}

	report := result.String()
	if !strings.Contains(report, "kit/root: 1") || !strings.Contains(report, "kit/usage: 1") {
		t.Errorf("report missing expected kit role counts:\n%s", report)
	}
	if !strings.Contains(report, "work/active: 1") {
		t.Errorf("report missing expected work role count:\n%s", report)
	}
	if strings.Contains(report, "content") {
		t.Error("String() must not dump raw document content")
	}
}

func TestExportWritesValidJSON(t *testing.T) {
	root := buildFixtureRepo(t)
	service := NewBuildService(root)
	result, err := service.Execute(domainindex.ScopeAll)
	if err != nil {
		t.Fatal(err)
	}

	exportPath := filepath.Join(t.TempDir(), "export.json")
	if err := result.Export(exportPath); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatal(err)
	}
	var docs []domaincontext.Document
	if err := json.Unmarshal(raw, &docs); err != nil {
		t.Fatalf("export is not valid JSON: %v", err)
	}
	if len(docs) != len(result.Documents) {
		t.Errorf("exported %d documents, want %d", len(docs), len(result.Documents))
	}
}
