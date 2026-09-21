package ingest

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	appindex "github.com/sentzunhat/hawp/librarian/src/internal/application/index"
	_ "modernc.org/sqlite"
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

func TestCorpusWalkersSkipBrokenSymlinks(t *testing.T) {
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
			if err != nil {
				t.Fatalf("expected broken symlink to be skipped, got %v", err)
			}
			if len(corpus.Documents) != 0 {
				t.Fatal("unreadable document indexed")
			}
		})
	}
}

func TestBuildCorpusRejectsPathsOutsideRepository(t *testing.T) {
	root := t.TempDir()
	for _, configuredPath := range []string{"../outside.md", "/tmp/outside.md", "nested/../../outside.md"} {
		t.Run(configuredPath, func(t *testing.T) {
			if _, err := buildCorpusFromRepo(root, []string{configuredPath}); err == nil {
				t.Fatalf("buildCorpusFromRepo(%q) unexpectedly succeeded", configuredPath)
			}
		})
	}
}

func TestBuildCorpusRejectsSymlinkedConfiguredPaths(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}

	for _, configuredPath := range []string{"docs-link", "docs-link/subdir", "linked.md"} {
		t.Run(configuredPath, func(t *testing.T) {
			root := t.TempDir()
			outside := t.TempDir()
			if err := os.MkdirAll(filepath.Join(outside, "subdir"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(outside, "outside.md"), []byte("external"), 0o644); err != nil {
				t.Fatal(err)
			}

			linkTarget := outside
			linkPath := filepath.Join(root, "docs-link")
			if configuredPath == "linked.md" {
				linkTarget = filepath.Join(outside, "outside.md")
				linkPath = filepath.Join(root, "linked.md")
			}
			if err := os.Symlink(linkTarget, linkPath); err != nil {
				t.Skipf("symlinks unavailable: %v", err)
			}

			if _, err := buildCorpusFromRepo(root, []string{configuredPath}); err == nil {
				t.Fatalf("buildCorpusFromRepo(%q) accepted symlinked path", configuredPath)
			}
		})
	}
}

func TestCorpusWalkersSkipSymlinkedMarkdown(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.md")
	if err := os.WriteFile(outside, []byte("external"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "linked.md")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	corpus, err := buildCorpusFromRepo(root, []string{"."})
	if err != nil {
		t.Fatal(err)
	}
	if len(corpus.Documents) != 0 {
		t.Fatalf("symlinked document indexed: %+v", corpus.Documents)
	}
}

func TestWorkCorpusPersistsCanonicalUUIDMetadata(t *testing.T) {
	root := t.TempDir()
	for rel, content := range map[string]string{
		".hawp/work/active/abcdef12/plan.md":             "# active plan",
		".hawp/work/closed/2026/09/21/12345678/plan.md":  "# closed plan",
		".hawp/work/active/not-a-work-id/legacy-plan.md": "# legacy plan",
	} {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	corpus, err := buildCorpusFromRepo(root, []string{".hawp/work"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := appindex.NewIngestService(filepath.Join(root, ".hawp", "db", "index.sqlite")).Execute(corpus); err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite", filepath.Join(root, ".hawp", "db", "index.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(`SELECT work_uuid, status FROM documents_metadata ORDER BY work_uuid`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var got []string
	for rows.Next() {
		var uuid, status string
		if err := rows.Scan(&uuid, &status); err != nil {
			t.Fatal(err)
		}
		got = append(got, uuid+":"+status)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	want := []string{"12345678:closed", "abcdef12:active"}
	if len(got) != len(want) {
		t.Fatalf("metadata rows = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("metadata rows = %v, want %v", got, want)
		}
	}
}
