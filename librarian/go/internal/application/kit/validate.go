// Package kit is the application service for `hawp kit validate`.
package kit

import (
	"fmt"
	"io"

	domainkit "github.com/sentzunhat/hawp/librarian/go/internal/domain/kit"
	"github.com/sentzunhat/hawp/librarian/go/internal/domain/kit/source"
	filesystemkit "github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/filesystem/kit"
)

// ValidateResult carries the kit validation outcome for rendering.
type ValidateResult struct {
	KitPath string
	Issues  []domainkit.Issue
	Checks  int
	Err     error
}

// Validate runs the kit checks against kitPath.
func Validate(kitPath string) ValidateResult {
	return ValidateWithWorkspace(kitPath, filesystemkit.NewAdapter())
}

// ValidateWithWorkspace allows the application workflow to be tested with a
// capability-specific workspace adapter.
func ValidateWithWorkspace(kitPath string, workspace source.Workspace) ValidateResult {
	snapshot, err := workspace.Read(kitPath)
	if err != nil {
		return ValidateResult{KitPath: kitPath, Err: err}
	}
	issues, checks := domainkit.Validate(snapshot)
	return ValidateResult{KitPath: kitPath, Issues: issues, Checks: checks}
}

// Render writes the report in the same shape as the TS kit:validate output
// and returns the exit code.
func Render(out, errOut io.Writer, result ValidateResult) int {
	fmt.Fprintln(out, "kit:validate")
	fmt.Fprintln(out, "============")
	fmt.Fprintf(out, "kit: %s\n\n", result.KitPath)
	if result.Err != nil {
		fmt.Fprintf(errOut, "kit validate error: %v\n", result.Err)
		return 1
	}

	if len(result.Issues) == 0 {
		fmt.Fprintf(out, "✓ %d checks passed, 0 issues\n", result.Checks)
		return 0
	}
	for _, issue := range result.Issues {
		fmt.Fprintf(errOut, "✗ %s: %s\n", issue.File, issue.Message)
	}
	fmt.Fprintf(errOut, "\n%d issue(s) found\n", len(result.Issues))
	return 1
}
