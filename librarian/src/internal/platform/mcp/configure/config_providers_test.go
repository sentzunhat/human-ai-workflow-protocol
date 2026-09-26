package mcp

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestConfigProviderExpansion(t *testing.T) {
	for _, tc := range []struct {
		input, want []string
	}{
		{nil, nil},
		{[]string{"all"}, []string{"claude", "cursor", "continue", "codex", "github"}},
		{[]string{"codex", "all", "claude", "all", "github"}, []string{"codex", "claude", "cursor", "continue", "github"}},
		{[]string{"github", "github"}, []string{"github"}},
	} {
		got, err := expandConfigProviders(tc.input)
		if err != nil || !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("%v: got %v, %v; want %v", tc.input, got, err, tc.want)
		}
	}
}

func TestInvalidConfigProviderRefusesBeforeAnyWrite(t *testing.T) {
	for _, invalid := range []string{"", "codxe", "Codex", "../codex"} {
		t.Run(invalid, func(t *testing.T) {
			root := t.TempDir()
			if err := WriteProviderConfigs(root, []string{"claude", invalid}); err == nil {
				t.Fatal("invalid provider accepted")
			}
			entries, err := os.ReadDir(root)
			if err != nil || len(entries) != 0 {
				t.Fatalf("wrote files before rejecting selection: %v, %v", entries, err)
			}
		})
	}
}

func TestWriteProviderConfigsPreflightsBeforeWrites(t *testing.T) {
	root := configureFixture(t)
	path := filepath.Join(root, ".vscode", "mcp.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"servers":{"hawp":{"type":"sse"}}}`), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := WriteProviderConfigs(root, []string{"claude", "github"}); err == nil {
		t.Fatal("accepted non-stdio GitHub/Copilot configuration")
	}
	for _, path := range []string{".mcp.json", ".gitignore"} {
		if _, err := os.Lstat(filepath.Join(root, path)); !os.IsNotExist(err) {
			t.Fatalf("wrote %s before preflight completed: %v", path, err)
		}
	}
}

func TestWriteProviderConfigsRejectsSymlinkedPrerequisite(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}
	root := configureFixture(t)
	path := hawpBinaryPath(root)
	external := filepath.Join(t.TempDir(), "external-hawp")
	if err := os.WriteFile(external, []byte("external"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, path); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if err := WriteProviderConfigs(root, []string{"claude"}); err == nil {
		t.Fatal("accepted symlinked prerequisite")
	}
	if _, err := os.Lstat(filepath.Join(root, ".mcp.json")); !os.IsNotExist(err) {
		t.Fatalf("configuration was written despite rejected prerequisite: %v", err)
	}
}

func TestWriteProviderConfigsPreflightsLaterProviderParentBeforeWrites(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}
	root := configureFixture(t)
	external := t.TempDir()
	cursorDir := filepath.Join(root, ".cursor")
	if err := os.Symlink(external, cursorDir); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if err := WriteProviderConfigs(root, []string{"claude", "cursor"}); err == nil {
		t.Fatal("accepted symlinked later provider parent")
	}
	for _, path := range []string{".mcp.json", ".gitignore"} {
		if _, err := os.Lstat(filepath.Join(root, path)); !os.IsNotExist(err) {
			t.Fatalf("wrote %s before all provider paths passed preflight: %v", path, err)
		}
	}
}
