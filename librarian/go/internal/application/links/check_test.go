package links

import (
	"os"
	"path/filepath"
	"testing"
)

func buildRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
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

func TestCheckFindsBrokenLinks(t *testing.T) {
	root := buildRepo(t, map[string]string{
		"README.md":               "[ok](.hawp/kit/start-here.md)\n[broken](missing.md)\n",
		".hawp/kit/start-here.md": "# start\n",
	})
	result := Check(root)
	if result.FilesChecked != 2 {
		t.Errorf("files checked = %d, want 2", result.FilesChecked)
	}
	if len(result.Failures) != 1 || result.Failures[0] != "README.md -> missing.md" {
		t.Fatalf("failures = %+v, want the one broken link", result.Failures)
	}
}

func TestCheckSkipsExternalAnchorsImagesFencesAndArchives(t *testing.T) {
	root := buildRepo(t, map[string]string{
		"README.md":                           "[ext](https://x.test)\n[anchor](#top)\n[mail](mailto:a@b.c)\n![img](gone.png)\n```\n[fenced](gone.md)\n```\n",
		".hawp/work/closed/2026/01/01/old.md": "[stale](../does-not-exist.md)\n",
		".hawp/work/BACKLOG.md":               "# backlog\n",
	})
	result := Check(root)
	if len(result.Failures) != 0 {
		t.Fatalf("failures = %+v, want none (external/anchor/image/fence/archive all skipped)", result.Failures)
	}
}

func TestCheckRootAbsolutePaths(t *testing.T) {
	root := buildRepo(t, map[string]string{
		"README.md":     "[abs](/docs/guide.md)\n",
		"docs/guide.md": "# guide\n",
	})
	if result := Check(root); len(result.Failures) != 0 {
		t.Fatalf("root-absolute link should resolve from repo root: %+v", result.Failures)
	}
}
