package mcp

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func configureFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, path := range []string{hawpBinaryPath(root), filepath.Join(root, ".hawp", "work", "BACKLOG.md")} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("fixture"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestConfigureFreshAndRepeat(t *testing.T) {
	root := configureFixture(t)
	for i := 0; i < 2; i++ {
		if err := Configure(root, []string{"claude", "codex", "cursor", "github"}); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{".mcp.json", ".codex/config.toml", ".cursor/mcp.json", ".vscode/mcp.json"} {
		data, err := os.ReadFile(filepath.Join(root, name))
		if err != nil || !strings.Contains(string(data), "--repo-root") {
			t.Fatalf("%s: %s %v", name, data, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(root, ".hawp/work/BACKLOG.md"))
	if err != nil || string(data) != "fixture" {
		t.Fatal("modified backlog", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".hawp/kit")); !os.IsNotExist(err) {
		t.Fatal("created kit")
	}
	ignore, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range []string{".mcp.json", ".cursor/mcp.json", ".codex/config.toml", ".vscode/mcp.json"} {
		if strings.Count(string(ignore), entry) != 1 {
			t.Fatalf("gitignore entry %q count = %d, content:\n%s", entry, strings.Count(string(ignore), entry), ignore)
		}
	}
}

func TestConfigurePreservesRemoteCodexBeforeOtherWrites(t *testing.T) {
	root := configureFixture(t)
	path := filepath.Join(root, ".codex/config.toml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	input := "[mcp_servers.hawp]\nurl = \"https://example.invalid/mcp\"\nenabled = false\n"
	if err := os.WriteFile(path, []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Configure(root, []string{"claude", "codex"}); err == nil {
		t.Fatal("replaced custom Codex configuration")
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != input {
		t.Fatal("changed custom config", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".mcp.json")); !os.IsNotExist(err) {
		t.Fatal("partial configuration write")
	}
}

func TestConfigurePreflightsGitHubBeforeOtherWrites(t *testing.T) {
	root := configureFixture(t)
	path := filepath.Join(root, ".vscode", "mcp.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"servers":{"hawp":{"type":"sse"}}}`), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := Configure(root, []string{"claude", "github"}); err == nil {
		t.Fatal("accepted non-stdio GitHub/Copilot configuration")
	}
	if _, err := os.Stat(filepath.Join(root, ".mcp.json")); !os.IsNotExist(err) {
		t.Fatalf("Claude configuration was written before GitHub preflight: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".gitignore")); !os.IsNotExist(err) {
		t.Fatalf("gitignore was written before GitHub preflight: %v", err)
	}
}

func TestConfigureMissingBinaryAndInvalidSelection(t *testing.T) {
	for _, providers := range [][]string{nil, {"codxe"}, {"claude"}} {
		root := t.TempDir()
		if err := Configure(root, providers); err == nil {
			t.Fatal("accepted missing prerequisites")
		}
		entries, err := os.ReadDir(root)
		if err != nil || len(entries) != 0 {
			t.Fatal("wrote configuration", err)
		}
	}
}

func TestConfigureRejectsSymlinkedPrerequisites(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}
	for _, prerequisite := range []string{"binary", "backlog", "hawp root"} {
		t.Run(prerequisite, func(t *testing.T) {
			root := configureFixture(t)
			path := hawpBinaryPath(root)
			if prerequisite == "backlog" {
				path = filepath.Join(root, ".hawp", "work", "BACKLOG.md")
			}
			external := filepath.Join(t.TempDir(), "external")
			if err := os.WriteFile(external, []byte("external"), 0o755); err != nil {
				t.Fatal(err)
			}
			if prerequisite == "hawp root" {
				externalRoot := t.TempDir()
				for _, rel := range []string{"bin/hawp", "work/BACKLOG.md"} {
					candidate := filepath.Join(externalRoot, rel)
					if err := os.MkdirAll(filepath.Dir(candidate), 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(candidate, []byte("external"), 0o755); err != nil {
						t.Fatal(err)
					}
				}
				if err := os.RemoveAll(filepath.Join(root, ".hawp")); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(externalRoot, filepath.Join(root, ".hawp")); err != nil {
					t.Skipf("symlinks unavailable: %v", err)
				}
			} else {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(external, path); err != nil {
					t.Skipf("symlinks unavailable: %v", err)
				}
			}

			if err := Configure(root, []string{"claude"}); err == nil {
				t.Fatal("accepted symlinked prerequisite")
			}
			if _, err := os.Lstat(filepath.Join(root, ".mcp.json")); !os.IsNotExist(err) {
				t.Fatalf("configuration was written despite rejected prerequisite: %v", err)
			}
		})
	}
}

func TestConfigureMigratesCustomCodexPolicy(t *testing.T) {
	root := configureFixture(t)
	path := filepath.Join(root, ".codex/config.toml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("# keep\n[mcp_servers.hawp]\ncommand='old'\nenabled=false\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Configure(root, []string{"codex"}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(data), "enabled=false\n") || !strings.HasPrefix(string(data), "# keep\n") {
		t.Fatalf("lost policy: %s %v", data, err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatal("changed permissions", err)
	}
}

func TestEnsureGitignoreEntryRejectsSymlink(t *testing.T) {
	root := t.TempDir()
	external := filepath.Join(t.TempDir(), "external-gitignore")
	if err := os.WriteFile(external, []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, filepath.Join(root, ".gitignore")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if err := ensureGitignoreEntry(root, ".mcp.json"); err == nil {
		t.Fatal("expected symlinked .gitignore to be rejected")
	}
	content, err := os.ReadFile(external)
	if err != nil || string(content) != "keep\n" {
		t.Fatalf("external gitignore was modified: %v %q", err, content)
	}
}

func TestManagedWritesReplaceHardLinks(t *testing.T) {
	for _, tt := range []struct {
		name    string
		write   func(root, path string) error
		path    string
		initial []byte
		want    string
	}{
		{
			name: "mcp json",
			write: func(root, path string) error {
				return writeMCPJSON(path, claudeServerEntry(root))
			},
			path:    ".mcp.json",
			initial: []byte("{}\n"),
			want:    `"hawp"`,
		},
		{
			name:    "codex toml",
			write:   func(root, path string) error { return writeCodexTOML(path, root) },
			path:    filepath.Join(".codex", "config.toml"),
			initial: []byte("[model]\nname = \"keep\"\n"),
			want:    "[mcp_servers.hawp]",
		},
		{
			name:    "gitignore",
			write:   func(root, path string) error { return ensureGitignoreEntry(root, ".mcp.json") },
			path:    ".gitignore",
			initial: []byte("keep\n"),
			want:    ".mcp.json",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, tt.path)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			shared := filepath.Join(t.TempDir(), "shared")
			original := tt.initial
			if err := os.WriteFile(shared, original, 0o600); err != nil {
				t.Fatal(err)
			}
			if err := os.Link(shared, path); err != nil {
				t.Skipf("hard links unavailable: %v", err)
			}

			if err := tt.write(root, path); err != nil {
				t.Fatal(err)
			}
			gotExternal, err := os.ReadFile(shared)
			if err != nil {
				t.Fatal(err)
			}
			if string(gotExternal) != string(original) {
				t.Fatalf("hard-linked source was modified: %q", gotExternal)
			}
			gotManaged, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(gotManaged), tt.want) {
				t.Fatalf("managed destination was not replaced: %q", gotManaged)
			}
		})
	}
}
