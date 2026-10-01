package normalizecmd

import (
	"path/filepath"
	"testing"
)

func TestNormalizeArgs(t *testing.T) {
	for _, args := range [][]string{{"--apply", "--dry-run"}, {"--hawp-root", "a", "--work-root", "b"}, {"--apply", "--format=xml"}, {"--unknown"}, {"extra"}, {"--output"}, {"--hawp-root="}} {
		if _, _, err := parseNormalizeArgs(args); err == nil {
			t.Errorf("accepted %v", args)
		}
	}
	opts, root, err := parseNormalizeArgs(nil)
	if err != nil || opts.Apply || root != "" {
		t.Fatalf("unsafe defaults: %+v %q %v", opts, root, err)
	}
	repo := t.TempDir()
	opts, root, err = parseNormalizeArgs([]string{"--apply", "--validate", "--migrate-folders", "--force-dirty", "--format=json", "--work-root=" + filepath.Join(repo, ".hawp", "work")})
	if err != nil || root != repo || !opts.Apply || !opts.Validate || !opts.MigrateFolders || !opts.ForceDirty || !opts.FormatJSON {
		t.Fatalf("wrong options: %+v %q %v", opts, root, err)
	}
}
