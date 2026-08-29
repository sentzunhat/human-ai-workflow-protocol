package filesystem

import (
	"path/filepath"
	"testing"
)

func TestResolveHawpProjectUsesCurrentRuntimeLayout(t *testing.T) {
	project := ResolveHawpProject("/repo")

	if got, want := project.DB, filepath.Join("/repo", ".hawp", "db"); got != want {
		t.Fatalf("DB = %q, want %q", got, want)
	}
	if got, want := project.Config, filepath.Join("/repo", ".hawp", "config"); got != want {
		t.Fatalf("Config = %q, want %q", got, want)
	}
}

func TestEnsureRuntimeFoldersCreatesDbAndConfig(t *testing.T) {
	project := ResolveHawpProject(t.TempDir())

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
