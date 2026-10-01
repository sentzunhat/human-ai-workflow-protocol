package work_test

import (
	"errors"
	"os"
	"testing"

	reposwork "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repositories/work"
)

func TestReadBacklog_MissingFile(t *testing.T) {
	_, err := reposwork.ReadBacklog("/nonexistent/path/BACKLOG.md")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected os.ErrNotExist, got %v", err)
	}
}
