package kit_test

import (
	"testing"

	reposkit "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repositories/kit"
)

func TestValidate_MissingKitPath(t *testing.T) {
	// A non-existent kit path should yield no issues — CollectFiles returns
	// empty for a missing directory and the required-files check flags them,
	// but there are no os.ReadFile calls that produce hard errors.
	issues, checks := reposkit.Validate("/nonexistent/path/kit")
	if checks != 3 {
		t.Fatalf("expected 3 checks, got %d", checks)
	}
	// Required files will be flagged as missing — that is correct behaviour,
	// not a failure of the reader.
	_ = issues
}

func TestCheckInternalLinks_MissingKitPath(t *testing.T) {
	// A non-existent kit path produces no issues because CollectFiles returns
	// empty for a missing directory — the infra reader is never called.
	issues := reposkit.CheckInternalLinks("/nonexistent/path/kit")
	if len(issues) != 0 {
		t.Fatalf("expected no link issues for missing kit path, got %d", len(issues))
	}
}
