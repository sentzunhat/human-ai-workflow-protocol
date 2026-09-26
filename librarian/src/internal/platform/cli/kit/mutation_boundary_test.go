package kit

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/kit/normalize"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/kit/validate"
)

func TestMutationBoundaryArgs(t *testing.T) {
	for _, args := range [][]string{
		{"--apply", "--dry-run"},
		{"--dry-run", "--apply"},
		{"--apply=true", "--dry-run=true"},
	} {
		if err := normalize.Run(args, t.TempDir()); err == nil {
			t.Errorf("kit normalize accepted conflicting modes %q", args)
		}
	}
	for _, value := range []string{"", " ", "\t\n"} {
		for _, args := range [][]string{{"--kit-path", value}, {"--kit-path=" + value}} {
			if err := normalize.Run(args, t.TempDir()); err == nil {
				t.Errorf("kit normalize accepted %q", args)
			}
			if err := validate.Run(args, t.TempDir()); err == nil {
				t.Errorf("kit validate accepted %q", args)
			}
		}
	}
}

func TestRejectedMutationCommandsPreserveFiles(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	for name, content := range map[string]string{
		".hawp/work/BACKLOG.md":     "# Backlog\n\n## Active Work\n\n| UUID | Type | Title | Status | Owner | Plan File | Updated |\n| --- | --- | --- | --- | --- | --- | --- |\n",
		".hawp/kit/Needs Rename.md": "# Keep this file\n",
		".hawp/kit/guide.md":        "[Keep link](Needs Rename.md)\n",
	} {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	snapshot := func() map[string]string {
		t.Helper()
		entries := map[string]string{}
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			if entry.IsDir() {
				entries[rel+"/"] = ""
				return nil
			}
			data, err := os.ReadFile(path)
			if err == nil {
				entries[rel] = string(data)
			}
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
		return entries
	}
	for _, tc := range []struct {
		name      string
		run       func([]string) error
		args      []string
		errorText string
	}{
		{"conflicting modes", func(args []string) error { return RunNormalize(args, root) }, []string{"--apply", "--dry-run"}, "mutually exclusive"},
		{"empty normalize path", func(args []string) error { return RunNormalize(args, root) }, []string{"--apply", "--kit-path="}, "non-empty path"},
		{"empty validation path", func(args []string) error { return RunValidate(args, root) }, []string{"--kit-path="}, "non-empty path"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			before := snapshot()
			err := tc.run(tc.args)
			if err == nil || !strings.Contains(err.Error(), tc.errorText) {
				t.Errorf("expected parser refusal, got %v", err)
			}
			if after := snapshot(); !reflect.DeepEqual(before, after) {
				t.Fatalf("rejected command changed fixture: before=%v after=%v", before, after)
			}
		})
	}
}
