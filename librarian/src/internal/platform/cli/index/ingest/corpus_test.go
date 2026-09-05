package ingest

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	appindex "github.com/sentzunhat/hawp/librarian/src/internal/application/index"
)

func TestWorkCorpusLifecycle(t *testing.T) {
	root := t.TempDir()
	want := map[string]string{
		"active/12345678/plan.md":            "active",
		"parked/87654321/plan.md":            "parked",
		"closed/2026/09/05/12345678/plan.md": "closed",
		"evidence/proof.md":                  "",
		"STATUS.md":                          "",
	}
	for rel := range want {
		path := filepath.Join(root, ".hawp/work", filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("# example"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	corpus, err := buildCorpusFromRepo(root, []string{".hawp/work"})
	if err != nil {
		t.Fatal(err)
	}
	if len(corpus.Documents) != len(want) {
		t.Fatalf("wrong document count: %d", len(corpus.Documents))
	}
	for _, doc := range corpus.Documents {
		rel, err := filepath.Rel(".hawp/work", doc.Path)
		if err != nil {
			t.Fatal(err)
		}
		status, ok := want[filepath.ToSlash(rel)]
		if !ok {
			t.Fatalf("unexpected path: %s", doc.Path)
		}
		if status == "" {
			if doc.Status != nil {
				t.Fatalf("invented status for %s: %s", doc.Path, *doc.Status)
			}
		} else if doc.Status == nil || *doc.Status != status {
			t.Fatalf("wrong lifecycle for %s: %v", doc.Path, doc.Status)
		}
		if doc.Content != "# example" {
			t.Fatalf("content lost: %+v", doc)
		}
	}
}

func TestCorpusReadFailuresAbort(t *testing.T) {
	for _, kind := range []string{"kit", "work", "custom"} {
		t.Run(kind, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.Symlink(filepath.Join(dir, "missing"), filepath.Join(dir, "broken.md")); err != nil {
				t.Skipf("symlink unavailable: %v", err)
			}
			corpus := &appindex.EnrichedCorpus{}
			var err error
			switch kind {
			case "kit":
				err = walkKitFiles(dir, corpus)
			case "work":
				err = walkWorkFiles(dir, corpus)
			case "custom":
				err = walkCustomPath(dir, "docs", corpus)
			}
			if !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("expected preserved read error, got %v", err)
			}
			if len(corpus.Documents) != 0 {
				t.Fatal("unreadable document indexed")
			}
		})
	}
}
