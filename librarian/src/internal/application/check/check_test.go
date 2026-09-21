package check

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunRejectsSymlinkedValidationTrees(t *testing.T) {
	tests := []struct {
		name string
		tree string
		want string
	}{
		{name: "kit", tree: filepath.Join(".hawp", "kit"), want: "kit validate error: unsafe validation path"},
		{name: "work", tree: filepath.Join(".hawp", "work"), want: "work validate error: unsafe validation path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, ".hawp"), 0o755); err != nil {
				t.Fatal(err)
			}
			external := t.TempDir()
			if err := os.Symlink(external, filepath.Join(root, tt.tree)); err != nil {
				t.Fatal(err)
			}

			var out, errOut bytes.Buffer
			if code := Run(&out, &errOut, root); code == 0 {
				t.Fatal("composite validation accepted a symlinked validation tree")
			}

			if !strings.Contains(errOut.String(), tt.want) {
				t.Fatalf("stderr = %q, want substring %q", errOut.String(), tt.want)
			}
			if strings.Contains(out.String(), external) || strings.Contains(errOut.String(), external) {
				t.Fatalf("composite validation exposed symlink target %q:\nstdout=%s\nstderr=%s", external, out.String(), errOut.String())
			}
		})
	}
}
