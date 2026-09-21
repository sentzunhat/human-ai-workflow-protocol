package migration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"hawp-source-layout/internal/mapping"
)

func put(t *testing.T, root, name, content string) {
	t.Helper()
	full := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	for name, content := range map[string]string{
		"go.mod":                                  "module " + modulePath + "\n\ngo 1.26\n",
		"internal/application/work/intake.go":     "package work\nfunc Intake() string { return \"intake\" }\n",
		"internal/application/work/validate.go":   "package work\nfunc Validate() string { return \"valid\" }\n",
		"internal/application/work/normalize.go":  "package work\nfunc Normalize() string { return Validate() }\n",
		"internal/application/work/draft_test.go": "//go:build integration\n\npackage work\nimport \"testing\"\nfunc TestIntake(t *testing.T) { if Intake() != \"intake\" { t.Fatal(\"wrong value\") } }\n",
		"cmd/hawp/main.go":                        "package main\nimport work \"" + modulePath + "/internal/application/work\"\n// Keep work.Intake in this comment.\nconst literal = \"work.Intake\"\nfunc main() { println(work.Intake(), work.Normalize()); work := struct { Intake string }{Intake: literal}; println(work.Intake) }\n",
		"assets/schema.json":                      "{\"keep\":true}\n",
	} {
		put(t, root, filepath.Join(sourceRoot, name), content)
	}
	return root
}

func prepared(t *testing.T, root string) (plan, []*file) {
	t.Helper()
	p, files, err := prepare(root)
	if err != nil {
		t.Fatal(err)
	}
	return p, files
}

func requireState(t *testing.T, root string, p plan, want string) {
	t.Helper()
	state, err := currentState(root, p)
	if err != nil || state != want {
		t.Fatalf("state = %q, %v; want %q", state, err, want)
	}
}

func TestPreviewRewritesImportsWithoutTouchingSource(t *testing.T) {
	root := fixture(t)
	p, files := prepared(t, root)
	if len(p.Files) != 7 {
		t.Fatalf("inventory = %d", len(p.Files))
	}
	for _, f := range files {
		switch f.Source {
		case "cmd/hawp/main.go":
			s := string(f.after)
			for _, want := range []string{"/work/intake\"", "/work/normalize\"", "work_layout1.Normalize()", "println(work.Intake)", "// Keep work.Intake in this comment.", "\"work.Intake\""} {
				if !strings.Contains(s, want) {
					t.Errorf("missing %q in rewritten caller:\n%s", want, s)
				}
			}
		case "internal/application/work/normalize.go":
			if !strings.Contains(string(f.after), "work_layout1.Validate()") {
				t.Fatalf("same-package reference was not qualified: %s", f.after)
			}
		case "internal/application/work/draft_test.go":
			if !strings.HasPrefix(string(f.after), "//go:build integration\n") || !strings.Contains(f.Destination, "/intake/") {
				t.Fatal("test ownership or build tag lost")
			}
		case "assets/schema.json":
			if f.Source != f.Destination || f.Before != f.After {
				t.Fatal("asset must be retained unchanged")
			}
		}
	}
	report := previewReport(p)
	for _, e := range p.Files {
		if !strings.Contains(report, e.Source) || !strings.Contains(report, e.Destination) {
			t.Errorf("report omitted %s", e.Source)
		}
	}
	if err := previewDiff(files); err != nil {
		t.Fatal(err)
	}
	requireState(t, root, p, "ready")
	// Rendering is deterministic, and a destination is never nested again.
	p2, _ := prepared(t, root)
	if !samePlan(p, p2) || previewReport(p) != previewReport(p2) {
		t.Fatal("non-deterministic preview")
	}
	for _, e := range p.Files {
		if dest, _ := mapping.Destination(e.Destination); dest != e.Destination {
			t.Errorf("destination is not stable: %s -> %s", e.Destination, dest)
		}
	}
}

func TestApplyAndAlreadyApplied(t *testing.T) {
	root := fixture(t)
	p, files := prepared(t, root)
	if err := apply(root, p, files); err != nil {
		t.Fatal(err)
	}
	requireState(t, root, p, "already-applied")
	for _, e := range p.Files {
		if e.Source != e.Destination {
			if _, err := os.Stat(filepath.Join(root, sourceRoot, e.Source)); !os.IsNotExist(err) {
				t.Fatalf("moved original still exists: %s", e.Source)
			}
		}
	}
}

func TestApplyRollsBackPostWriteMismatch(t *testing.T) {
	root := fixture(t)
	p, files := prepared(t, root)
	original, _ := prepared(t, root)
	// Simulate a transformed output that disagrees with the reviewed fingerprint.
	for i := range p.Files {
		if p.Files[i].Source != p.Files[i].Destination {
			p.Files[i].After = digest([]byte("mismatched expected output"))
			break
		}
	}
	if err := apply(root, p, files); err == nil {
		t.Fatal("expected output-verification failure")
	}
	requireState(t, root, original, "ready")
}

func TestChangedSnapshotIsRejected(t *testing.T) {
	for _, mutate := range []string{"content", "mode", "extra", "missing", "partial"} {
		t.Run(mutate, func(t *testing.T) {
			root := fixture(t)
			p, _ := prepared(t, root)
			name := filepath.Join(root, sourceRoot, "internal/application/work/intake.go")
			switch mutate {
			case "content":
				put(t, root, filepath.Join(sourceRoot, "internal/application/work/intake.go"), "package work\n")
			case "mode":
				if err := os.Chmod(name, 0o755); err != nil {
					t.Fatal(err)
				}
			case "extra":
				put(t, root, filepath.Join(sourceRoot, "extra.txt"), "new source support file")
			case "missing":
				if err := os.Remove(name); err != nil {
					t.Fatal(err)
				}
			case "partial":
				if err := os.MkdirAll(strings.TrimSuffix(name, ".go"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Rename(name, strings.TrimSuffix(name, ".go")+"/intake.go"); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := currentState(root, p); err == nil {
				t.Fatal("changed source must invalidate the plan")
			}
		})
	}
}

func TestUnsafeDestinationsAreRejected(t *testing.T) {
	for _, kind := range []string{"occupied", "symlink", "ignored-go", "private-dependency"} {
		t.Run(kind, func(t *testing.T) {
			root := fixture(t)
			p, _ := prepared(t, root)
			switch kind {
			case "occupied":
				put(t, root, ".gitignore", "/librarian/src/internal/application/work/intake/\n")
				put(t, root, sourceRoot+"/internal/application/work/intake/intake.go", "do not replace")
				if err := checkDestinations(root, p); err == nil {
					t.Fatal("ignored destination was not protected")
				}
			case "symlink":
				if err := os.Symlink(t.TempDir(), filepath.Join(root, sourceRoot, "internal/application/work/intake")); err != nil {
					t.Fatal(err)
				}
				if err := checkDestinations(root, p); err == nil {
					t.Fatal("destination symlink was not protected")
				}
			case "ignored-go":
				put(t, root, ".gitignore", "/librarian/src/ignored.go\n")
				put(t, root, sourceRoot+"/ignored.go", "package main\n")
			case "private-dependency":
				put(t, root, sourceRoot+"/internal/application/work/validate.go", "package work\nfunc validate() string { return \"valid\" }\n")
				put(t, root, sourceRoot+"/internal/application/work/normalize.go", "package work\nfunc Normalize() string { return validate() }\n")
			}
			if _, _, err := prepare(root); err == nil {
				t.Fatal("unsafe migration must be refused")
			}
		})
	}
}

func TestArtifactCannotWriteSourceOrFollowSymlink(t *testing.T) {
	root := fixture(t)
	p, _ := prepared(t, root)
	if err := writeArtifact(root, filepath.Join(root, sourceRoot, "go.mod"), []byte("bad")); err == nil {
		t.Fatal("artifact overwrote source")
	}
	link := filepath.Join(root, "preview.md")
	if err := os.Symlink(filepath.Join(root, sourceRoot, "go.mod"), link); err != nil {
		t.Fatal(err)
	}
	if err := writeArtifact(root, link, []byte("bad")); err == nil {
		t.Fatal("artifact followed symlink")
	}
	requireState(t, root, p, "ready")
}

func TestUnsafePlanMappingsAreRejected(t *testing.T) {
	for _, bad := range []string{"../escape", "/absolute", "cmd/hawp/main.go"} {
		root := fixture(t)
		p, _ := prepared(t, root)
		p.Files[0].Destination = bad
		if _, err := currentState(root, p); err == nil {
			t.Fatalf("accepted %s", bad)
		}
	}
}
