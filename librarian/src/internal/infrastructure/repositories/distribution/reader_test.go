package distribution_test

import (
	"errors"
	"io/fs"
	"testing"

	reposdistribution "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repositories/distribution"
)

func TestComputeExpectedOutputs_MissingSourceFragments(t *testing.T) {
	// A non-existent repoRoot causes os.ReadFile to fail on the first source
	// fragment; the error should propagate as a wrapped fs.ErrNotExist.
	_, err := reposdistribution.ComputeExpectedOutputs("/nonexistent/repo/root")
	if err == nil {
		t.Fatal("expected error for missing source fragments, got nil")
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected fs.ErrNotExist in error chain, got %v", err)
	}
}

func TestFindDownstreamPathLeaks_MissingFiles(t *testing.T) {
	// A non-existent repoRoot means all downstream target files are missing;
	// the function should return an empty slice (missing files are skipped).
	leaks, err := reposdistribution.FindDownstreamPathLeaks("/nonexistent/repo/root")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(leaks) != 0 {
		t.Fatalf("expected no leaks for missing files, got %d", len(leaks))
	}
}
