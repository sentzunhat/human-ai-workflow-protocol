package db_test

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	appdb "github.com/sentzunhat/hawp/librarian/src/internal/application/db"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

func TestInitServicePlansHawpHomeLayout(t *testing.T) {
	home := t.TempDir()
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", home)
	} else {
		t.Setenv("HOME", home)
	}

	service := appdb.NewInitService(filesystem.NewLayoutService())
	result, err := service.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	root := filepath.Join(home, ".hawp")
	want := appdb.InitResult{
		DBPath:        filepath.Join(root, "index", "librarian.db"),
		ModelsPath:    filepath.Join(root, "models"),
		DownloadsPath: filepath.Join(root, "cache", "downloads"),
	}
	if result != want {
		t.Errorf("Execute() = %+v, want %+v", result, want)
	}
}

func TestInitResultStringListsAllPaths(t *testing.T) {
	result := appdb.InitResult{
		DBPath:        "/x/.hawp/index/librarian.db",
		ModelsPath:    "/x/.hawp/models",
		DownloadsPath: "/x/.hawp/cache/downloads",
	}
	out := result.String()
	for _, part := range []string{result.DBPath, result.ModelsPath, result.DownloadsPath} {
		if !strings.Contains(out, part) {
			t.Errorf("String() output missing %q:\n%s", part, out)
		}
	}
}
