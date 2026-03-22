package selfreplace

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestReplaceSwapsFileAtomically(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Replace is unsupported on Windows by design")
	}
	dir := t.TempDir()
	target := filepath.Join(dir, "hawp")
	source := filepath.Join(dir, "staged")

	if err := os.WriteFile(target, []byte("old binary"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(source, []byte("new binary"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Replace(source, target); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "new binary" {
		t.Errorf("target content = %q, want new binary", content)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		t.Error("replaced binary should be executable")
	}
	if _, err := os.Stat(source); !os.IsNotExist(err) {
		t.Error("source should be consumed (renamed away) after Replace")
	}
}

func TestReplaceOnWindowsReturnsClearError(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-only behavior")
	}
	if err := Replace("a", "b"); err != ErrUnsupportedOnWindows {
		t.Errorf("err = %v, want ErrUnsupportedOnWindows", err)
	}
}
