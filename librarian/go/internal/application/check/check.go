// Package check is the composite validation pass: kit validate + work
// validate + links check in one command.
package check

import (
	"fmt"
	"io"
	"path/filepath"

	appkit "github.com/sentzunhat/hawp/librarian/go/internal/application/kit"
	applinks "github.com/sentzunhat/hawp/librarian/go/internal/application/links"
	appwork "github.com/sentzunhat/hawp/librarian/go/internal/application/work"
	domainwork "github.com/sentzunhat/hawp/librarian/go/internal/domain/work"
)

// Run executes all three validations against repoRoot and returns the exit
// code (0 when none fail; work validate WARNs do not fail).
func Run(out, errOut io.Writer, repoRoot string) int {
	fmt.Fprintln(out, "hawp:check")
	fmt.Fprintln(out, "==========")
	fmt.Fprintf(out, "repo: %s\n\n", repoRoot)
	failed := 0

	fmt.Fprintln(out, "[1/3] kit validate")
	kitResult := appkit.Validate(filepath.Join(repoRoot, ".hawp", "kit"))
	if appkit.Render(out, errOut, kitResult) != 0 {
		failed++
	}
	fmt.Fprintln(out)

	fmt.Fprintln(out, "[2/3] work validate")
	workDir := filepath.Join(repoRoot, ".hawp", "work")
	report, err := appwork.Validate(workDir)
	if err != nil {
		fmt.Fprintf(errOut, "work validate error: %v\n", err)
		failed++
	} else {
		fmt.Fprintf(out, "  Result: VALIDATION %s (%d issues, %d warnings)\n",
			report.Overall, report.Failed, report.Warnings)
		if report.Overall == domainwork.StatusFail {
			fmt.Fprintln(errOut, "  Run `hawp work validate` for the detailed report.")
			failed++
		}
	}
	fmt.Fprintln(out)

	fmt.Fprintln(out, "[3/3] links check")
	linksResult := applinks.Check(repoRoot)
	if applinks.Render(out, errOut, linksResult) != 0 {
		failed++
	}
	fmt.Fprintln(out)

	if failed > 0 {
		fmt.Fprintf(errOut, "hawp check: %d of 3 validations failed\n", failed)
		return 1
	}
	fmt.Fprintln(out, "hawp check: all 3 validations passed")
	return 0
}
