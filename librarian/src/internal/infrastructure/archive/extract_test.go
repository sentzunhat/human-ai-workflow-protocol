package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func buildTarGz(t *testing.T, files map[string]string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "archive.tgz")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz := gzip.NewWriter(f)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	defer tw.Close()

	for name, content := range files {
		if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0o644, Size: int64(len(content))}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	return path
}

func buildZip(t *testing.T, files map[string]string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "archive.zip")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	defer zw.Close()

	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	return path
}

func TestExtractMemberFromTarGz(t *testing.T) {
	archivePath := buildTarGz(t, map[string]string{
		"pkg/lib/libonnxruntime.so.1.27.1": "shared library bytes",
		"pkg/lib/other.so":                 "not this one",
	})
	dest := filepath.Join(t.TempDir(), "out", "libonnxruntime.so")

	if err := ExtractMember(archivePath, "pkg/lib/libonnxruntime.so.1.27.1", dest); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "shared library bytes" {
		t.Errorf("content = %q", got)
	}
	info, err := os.Stat(dest)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm()&0o111 == 0 {
		t.Error("extracted library should be executable")
	}
}

func TestExtractMemberFromZip(t *testing.T) {
	archivePath := buildZip(t, map[string]string{
		"pkg/lib/onnxruntime.dll": "windows dll bytes",
	})
	dest := filepath.Join(t.TempDir(), "onnxruntime.dll")

	if err := ExtractMember(archivePath, "pkg/lib/onnxruntime.dll", dest); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "windows dll bytes" {
		t.Errorf("content = %q", got)
	}
}

func TestExtractMemberNotFound(t *testing.T) {
	archivePath := buildTarGz(t, map[string]string{"a": "1"})
	err := ExtractMember(archivePath, "does-not-exist", filepath.Join(t.TempDir(), "out"))
	if err == nil {
		t.Fatal("expected error for missing member")
	}
}

func TestExtractMemberUnsupportedType(t *testing.T) {
	path := filepath.Join(t.TempDir(), "archive.rar")
	os.WriteFile(path, []byte("x"), 0o644)
	if err := ExtractMember(path, "x", filepath.Join(t.TempDir(), "out")); err == nil {
		t.Fatal("expected error for unsupported archive type")
	}
}

func TestExtractAllRejectsTraversalMember(t *testing.T) {
	archivePath := buildTarGz(t, map[string]string{"../../escaped.txt": "must not write"})
	root := t.TempDir()
	dest := filepath.Join(root, "extracted")

	if err := ExtractAll(archivePath, dest); err == nil {
		t.Fatal("expected traversal member to be rejected")
	}
	if _, err := os.Stat(filepath.Join(root, "escaped.txt")); !os.IsNotExist(err) {
		t.Fatalf("traversal target was created: %v", err)
	}
}

func TestExtractAllRejectsAbsoluteMember(t *testing.T) {
	root := t.TempDir()
	absolute := filepath.Join(root, "absolute.txt")
	archivePath := buildTarGz(t, map[string]string{absolute: "must not write"})
	dest := filepath.Join(root, "extracted")

	if err := ExtractAll(archivePath, dest); err == nil {
		t.Fatal("expected absolute member to be rejected")
	}
	if _, err := os.Stat(absolute); !os.IsNotExist(err) {
		t.Fatalf("absolute target was created: %v", err)
	}
}

func TestExtractAllRejectsSymlinkAncestor(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}

	root := t.TempDir()
	dest := filepath.Join(root, "extracted")
	outside := filepath.Join(root, "outside")
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dest, "sub")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	archivePath := buildTarGz(t, map[string]string{"sub/escaped.txt": "must not write"})
	if err := ExtractAll(archivePath, dest); err == nil {
		t.Fatal("expected symlink ancestor to be rejected")
	}
	if _, err := os.Stat(filepath.Join(outside, "escaped.txt")); !os.IsNotExist(err) {
		t.Fatalf("symlink target was created: %v", err)
	}
}
