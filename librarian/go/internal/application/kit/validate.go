// Package kit is the application service for `hawp kit validate`.
package kit

import (
	"fmt"
	"io"

	domainkit "github.com/sentzunhat/hawp/librarian/go/internal/domain/kit"
)

// ValidateResult carries the kit validation outcome for rendering.
type ValidateResult struct {
	KitPath string
	Issues  []domainkit.Issue
	Checks  int
}

// Validate runs the kit checks against kitPath.
func Validate(kitPath string) ValidateResult {
	issues, checks := domainkit.Validate(kitPath)
	return ValidateResult{KitPath: kitPath, Issues: issues, Checks: checks}
}

// Render writes the report in the same shape as the TS kit:validate output
// and returns the exit code.
func Render(out, errOut io.Writer, result ValidateResult) int {
	fmt.Fprintln(out, "kit:validate")
	fmt.Fprintln(out, "============")
	fmt.Fprintf(out, "kit: %s\n\n", result.KitPath)

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
