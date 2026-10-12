package filesystem_test

import (
	"os"
	"path/filepath"
	"testing"

	infrafs "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

func TestResolveHawpProjectUsesCurrentRuntimeLayout(t *testing.T) {
	project := infrafs.ResolveHawpProject("/repo")

	if got, want := project.DB, filepath.Join("/repo", ".hawp", "db"); got != want {
		t.Fatalf("DB = %q, want %q", got, want)
	}
	if got, want := project.Config, filepath.Join("/repo", ".hawp", "config"); got != want {
		t.Fatalf("Config = %q, want %q", got, want)
	}
}

func TestEnsureRuntimeFoldersCreatesDbAndConfig(t *testing.T) {
	project := infrafs.ResolveHawpProject(t.TempDir())

	created, err := project.EnsureRuntimeFolders()
	if err != nil {
		t.Fatalf("EnsureRuntimeFolders() error = %v", err)
	}
	if !created {
		t.Fatalf("EnsureRuntimeFolders() created = false, want true on first run")
	}

	if got, want := project.GetSearchIndexPath(), filepath.Join(project.Root, "db", "index.sqlite"); got != want {
		t.Fatalf("GetSearchIndexPath() = %q, want %q", got, want)
	}
	if got, want := project.GetProjectConfigPath(), filepath.Join(project.Root, "config", "context.json"); got != want {
		t.Fatalf("GetProjectConfigPath() = %q, want %q", got, want)
	}
}

func TestResolveSafeSearchIndexPathAllowsMissingRuntimeDirectories(t *testing.T) {
	root := t.TempDir()
	got, err := infrafs.ResolveSafeSearchIndexPath(root)
	if err != nil {
		t.Fatalf("ResolveSafeSearchIndexPath() error = %v", err)
	}
	want := filepath.Join(root, ".hawp", "db", "index.sqlite")
	if got != want {
		t.Fatalf("ResolveSafeSearchIndexPath() = %q, want %q", got, want)
	}
}

func TestResolveSafeSearchIndexPathRejectsHardLinkedDatabaseFiles(t *testing.T) {
	for _, suffix := range []string{"", "-wal", "-shm", "-journal"} {
		t.Run(suffix, func(t *testing.T) {
			root := t.TempDir()
			outside := filepath.Join(t.TempDir(), "external.sqlite")
			if err := os.MkdirAll(filepath.Join(root, ".hawp", "db"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(outside, []byte("external database"), 0o600); err != nil {
				t.Fatal(err)
			}
			indexPath := filepath.Join(root, ".hawp", "db", "index.sqlite") + suffix
			if err := os.Link(outside, indexPath); err != nil {
				t.Skipf("hard links unavailable: %v", err)
			}
			if _, err := infrafs.ResolveSafeSearchIndexPath(root); err == nil {
				t.Fatal("expected hard-linked search database file to be rejected")
			}
		})
	}
}

func TestResolveSafeSearchIndexPathRejectsSymlinkedDatabaseFiles(t *testing.T) {
	for _, suffix := range []string{"", "-wal", "-shm", "-journal"} {
		t.Run(suffix, func(t *testing.T) {
			root := t.TempDir()
			outside := filepath.Join(t.TempDir(), "external.sqlite")
			if err := os.WriteFile(outside, []byte("external database"), 0o600); err != nil {
				t.Fatal(err)
			}
			dbDir := filepath.Join(root, ".hawp", "db")
			if err := os.MkdirAll(dbDir, 0o755); err != nil {
				t.Fatal(err)
			}
			indexPath := filepath.Join(dbDir, "index.sqlite") + suffix
			if err := os.Symlink(outside, indexPath); err != nil {
				t.Skipf("symlink capability unavailable (enable Windows Developer Mode or SeCreateSymbolicLinkPrivilege): %v", err)
			}
			if _, err := infrafs.ResolveSafeSearchIndexPath(root); err == nil {
				t.Fatal("expected symlinked search database file to be rejected")
			}
		})
	}
}

func TestResolveSafeSearchIndexPathRejectsSymlinkedRuntimeAncestors(t *testing.T) {

	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, root, outside string)
	}{
		{
			name: "hawp root",
			setup: func(t *testing.T, root, outside string) {
				t.Helper()
				if err := os.Symlink(outside, filepath.Join(root, ".hawp")); err != nil {
					t.Skipf("symlink capability unavailable (enable Windows Developer Mode or SeCreateSymbolicLinkPrivilege): %v", err)
				}
			},
		},
		{
			name: "database directory",
			setup: func(t *testing.T, root, outside string) {
				t.Helper()
				if err := os.MkdirAll(filepath.Join(root, ".hawp"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(outside, filepath.Join(root, ".hawp", "db")); err != nil {
					t.Skipf("symlink capability unavailable (enable Windows Developer Mode or SeCreateSymbolicLinkPrivilege): %v", err)
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			outside := t.TempDir()
			tc.setup(t, root, outside)
			if _, err := infrafs.ResolveSafeSearchIndexPath(root); err == nil {
				t.Fatal("expected symlinked search index path to be rejected")
			}
		})
	}
}

func TestEnsureRuntimeFoldersRejectsSymlinkedRuntimeDirectory(t *testing.T) {

	root := t.TempDir()
	outside := t.TempDir()
	projectRoot := infrafs.ResolveHawpProject(root)
	if err := os.MkdirAll(projectRoot.Root, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, projectRoot.DB); err != nil {
		t.Skipf("symlink capability unavailable (enable Windows Developer Mode or SeCreateSymbolicLinkPrivilege): %v", err)
	}
	if _, err := projectRoot.EnsureRuntimeFolders(); err == nil {
		t.Fatal("expected symlinked runtime directory to be rejected")
	}
}
