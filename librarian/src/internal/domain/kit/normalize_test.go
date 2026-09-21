package kit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/markdown"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

func TestNormalizeFileName(t *testing.T) {
	cases := []struct {
		in   string
		want string // "" = already normalized / exempt
	}{
		{"README.md", ""},
		{"start-here.md", ""},
		{"Bad Name.MD", "bad-name.md"},
		{"Some_File.md", "some-file.md"},
		{"UPPER.md", "upper.md"},
		{"weird--name.md", "weird-name.md"},
		{"-leading.md", "leading.md"},
		{"noext", ""},
		{"Mixed Case Doc", "mixed-case-doc"},
	}
	for _, c := range cases {
		if got := NormalizeFileName(c.in); got != c.want {
			t.Errorf("NormalizeFileName(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

var normalizeTestSource = &KitSource{
	FileLister:     func(kitPath string, skipReadme bool) []string { return markdown.CollectFiles(kitPath, skipReadme) },
	BlankFences:    markdown.BlankFences,
	Exists:         repo.Exists,
	ToRepoRelative: repo.ToRepoRelative,
}

func TestPlanAndApplyNormalization(t *testing.T) {
	kitPath := t.TempDir()
	mustWrite := func(rel, content string) {
		full := filepath.Join(kitPath, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mustWrite("usage/Bad Name.md", "# doc\n")
	mustWrite("usage/guide.md", "see [bad](Bad%20Name.md) and [good](../start-here.md)\n")
	mustWrite("start-here.md", "see [bad too](usage/Bad Name.md)\n")

	renames := PlanFileRenames(kitPath, os.ReadDir)
	if len(renames) != 1 || filepath.Base(renames[0].To) != "bad-name.md" {
		t.Fatalf("renames = %+v, want one rename to bad-name.md", renames)
	}

	renameMap := map[string]string{renames[0].From: renames[0].To}
	updates := normalizeTestSource.PlanLinkUpdates(kitPath, renameMap, os.ReadFile)
	// Only the plain-path link resolves to the renamed file (the %20 form
	// does not match on-disk resolution here, matching TS behavior).
	if len(updates) != 1 || updates[0].To != "usage/bad-name.md" {
		t.Fatalf("updates = %+v, want start-here.md link update", updates)
	}

	if from, to, err := ApplyRenames(renames, os.Stat, os.Rename); err != nil || from != "" {
		t.Fatalf("ApplyRenames failed: %v conflict=%s->%s", err, from, to)
	}
	changed, err := ApplyLinkUpdates(updates, os.ReadFile, os.WriteFile)
	if err != nil || changed != 1 {
		t.Fatalf("ApplyLinkUpdates changed=%d err=%v, want 1 file", changed, err)
	}

	content, _ := os.ReadFile(filepath.Join(kitPath, "start-here.md"))
	if !strings.Contains(string(content), "(usage/bad-name.md)") {
		t.Errorf("link not rewritten: %s", content)
	}
	if _, err := os.Stat(filepath.Join(kitPath, "usage", "bad-name.md")); err != nil {
		t.Error("renamed file missing")
	}
}

func TestPlanAndApplyNormalizationPreservesUnicodeFenceOffsets(t *testing.T) {
	kitPath := t.TempDir()
	mustWrite := func(rel, content string) {
		full := filepath.Join(kitPath, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mustWrite("usage/Bad Name.md", "# doc\n")
	mustWrite("start-here.md", "```\n日本語\n```\nSee [guide](usage/Bad Name.md).\n")

	renames := PlanFileRenames(kitPath, os.ReadDir)
	if len(renames) != 1 {
		t.Fatalf("renames = %+v, want one rename", renames)
	}
	updates := normalizeTestSource.PlanLinkUpdates(kitPath, map[string]string{renames[0].From: renames[0].To}, os.ReadFile)
	if len(updates) != 1 || updates[0].From != "usage/Bad Name.md" || updates[0].To != "usage/bad-name.md" {
		t.Fatalf("updates = %+v, want one link rewrite", updates)
	}
	if _, _, err := ApplyRenames(renames, os.Stat, os.Rename); err != nil {
		t.Fatal(err)
	}
	if _, err := ApplyLinkUpdates(updates, os.ReadFile, os.WriteFile); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(kitPath, "start-here.md"))
	if err != nil {
		t.Fatal(err)
	}
	want := "```\n日本語\n```\nSee [guide](usage/bad-name.md).\n"
	if string(content) != want {
		t.Fatalf("normalized content = %q, want %q", content, want)
	}
}

func TestApplyRenamesRefusesOverwrite(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"Doc.md", "doc.md"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	renames := []FileRename{{From: filepath.Join(dir, "Doc.md"), To: filepath.Join(dir, "doc.md")}}
	from, to, err := ApplyRenames(renames, os.Stat, os.Rename)
	if err != nil {
		t.Fatal(err)
	}
	if from == "" || filepath.Base(to) != "doc.md" {
		t.Fatalf("expected conflict, got from=%q to=%q", from, to)
	}
}
