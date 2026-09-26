package index

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
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

func TestExecuteRejectsSymlinkedCorpusRoot(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}
	root := buildFixtureRepo(t)
	external := t.TempDir()
	if err := os.WriteFile(filepath.Join(external, "outside.md"), []byte("# outside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	kitRoot := filepath.Join(root, ".hawp", "kit")
	if err := os.RemoveAll(kitRoot); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, kitRoot); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	_, err := NewBuildService(root).Execute(domainindex.ScopeKit)
	if err == nil {
		t.Fatal("index build accepted a symlinked kit corpus root")
	}
	if !strings.Contains(err.Error(), "unsafe kit corpus root") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteRejectsSymlinkedCorpusDescendants(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}

	newExternalFile := func(t *testing.T, contents string) string {
		t.Helper()
		path := filepath.Join(t.TempDir(), "outside.md")
		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
		return path
	}

	tests := []struct {
		name    string
		scope   domainindex.DocumentScope
		link    string
		target  func(t *testing.T) string
		message string
	}{
		{
			name:    "nested kit markdown file",
			scope:   domainindex.ScopeKit,
			link:    ".hawp/kit/usage/outside.md",
			target:  func(t *testing.T) string { return newExternalFile(t, "# external kit\n") },
			message: "unsafe kit corpus root",
		},
		{
			name:    "work backlog",
			scope:   domainindex.ScopeWork,
			link:    ".hawp/work/BACKLOG.md",
			target:  func(t *testing.T) string { return newExternalFile(t, "# external backlog\n") },
			message: "unsafe work corpus root",
		},
		{
			name:  "work role directory",
			scope: domainindex.ScopeWork,
			link:  ".hawp/work/active",
			target: func(t *testing.T) string {
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "outside.md"), []byte("# external work\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				return dir
			},
			message: "unsafe work corpus root",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := buildFixtureRepo(t)
			link := filepath.Join(root, filepath.FromSlash(test.link))
			if err := os.RemoveAll(link); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(test.target(t), link); err != nil {
				t.Skipf("symlinks unavailable: %v", err)
			}

			_, err := NewBuildService(root).Execute(test.scope)
			if err == nil {
				t.Fatal("index build accepted a symlinked corpus descendant")
			}
			if !strings.Contains(err.Error(), test.message) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
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
