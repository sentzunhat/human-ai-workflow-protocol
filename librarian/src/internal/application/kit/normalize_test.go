package kit

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNormalizeApplyRejectsSymlinkedKitPaths(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}

	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, repoRoot, kitPath, outside string)
	}{
		{
			name: "kit root",
			setup: func(t *testing.T, repoRoot, kitPath, outside string) {
				t.Helper()
				if err := os.MkdirAll(filepath.Dir(kitPath), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(outside, kitPath); err != nil {
					t.Skipf("symlinks unavailable: %v", err)
				}
			},
		},
		{
			name: "nested markdown file",
			setup: func(t *testing.T, repoRoot, kitPath, outside string) {
				t.Helper()
				if err := os.MkdirAll(kitPath, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(filepath.Join(outside, "outside.md"), filepath.Join(kitPath, "linked.md")); err != nil {
					t.Skipf("symlinks unavailable: %v", err)
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repoRoot := t.TempDir()
			outside := t.TempDir()
			outsideFile := filepath.Join(outside, "outside.md")
			if err := os.WriteFile(outsideFile, []byte("[Guide](Old Name.md)"), 0o644); err != nil {
				t.Fatal(err)
			}
			kitPath := filepath.Join(repoRoot, ".hawp", "kit")
			tc.setup(t, repoRoot, kitPath, outside)

			var out, errOut bytes.Buffer
			code := Normalize(&out, &errOut, NormalizeOptions{KitPath: kitPath, RepoRoot: repoRoot, Apply: true})
			if code == 0 {
				t.Fatalf("Normalize apply unexpectedly succeeded; stderr=%q", errOut.String())
			}
			if !strings.Contains(errOut.String(), "unsafe kit path") {
				t.Fatalf("expected unsafe kit path error, got %q", errOut.String())
			}
			got, err := os.ReadFile(outsideFile)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != "[Guide](Old Name.md)" {
				t.Fatalf("outside file was modified: %q", got)
			}
		})
	}
}
