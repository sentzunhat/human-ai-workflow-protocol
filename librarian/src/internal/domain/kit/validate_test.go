package kit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func buildKit(t *testing.T, files map[string]string) string {
	t.Helper()
	kitPath := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(kitPath, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return kitPath
}

func completeKit(t *testing.T) map[string]string {
	t.Helper()
	files := map[string]string{}
	for _, rel := range RequiredFiles {
		files[rel] = "# " + rel + "\n"
	}
	return files
}

func TestValidateCleanKit(t *testing.T) {
	kitPath := buildKit(t, completeKit(t))
	issues, checks := Validate(kitPath)
	if checks != 3 {
		t.Errorf("checks = %d, want 3", checks)
	}
	if len(issues) != 0 {
		t.Fatalf("issues on clean kit: %+v", issues)
	}
}

func TestFileNaming(t *testing.T) {
	files := completeKit(t)
	files["Bad Name.md"] = "# bad"
	files["README.md"] = "# allowed uppercase"
	kitPath := buildKit(t, files)

	issues := CheckFileNaming(kitPath)
	if len(issues) != 1 || !strings.Contains(issues[0].Message, "lowercase-hyphen") {
		t.Fatalf("naming issues = %+v, want exactly the bad name", issues)
	}
}

func TestRequiredFiles(t *testing.T) {
	files := completeKit(t)
	delete(files, "start-here.md")
	kitPath := buildKit(t, files)

	issues := CheckRequiredFiles(kitPath)
	if len(issues) != 1 || issues[0].File != "start-here.md" {
		t.Fatalf("required issues = %+v, want missing start-here.md", issues)
	}
}

func TestInternalLinks(t *testing.T) {
	files := completeKit(t)
	files["usage/guide.md"] = "[ok](../start-here.md)\n[broken](missing.md)\n[ext](https://x.test)\n```\n[fenced](nope.md)\n```\n"
	kitPath := buildKit(t, files)

	issues := CheckInternalLinks(kitPath)
	if len(issues) != 1 || !strings.Contains(issues[0].Message, "missing.md") {
		t.Fatalf("link issues = %+v, want only the broken link", issues)
	}
}
