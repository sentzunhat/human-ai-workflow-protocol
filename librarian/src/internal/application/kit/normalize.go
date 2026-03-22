package kit

import (
	"fmt"
	"io"

	domainkit "github.com/sentzunhat/hawp/librarian/src/internal/domain/kit"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

// NormalizeOptions configures a kit normalize run.
type NormalizeOptions struct {
	KitPath  string
	RepoRoot string
	Apply    bool
}

// Normalize plans (and in apply mode performs) kit file renames and link
// rewrites, writing a dry-run/apply report. Returns the exit code.
func Normalize(out, errOut io.Writer, opts NormalizeOptions) int {
	fmt.Fprintln(out, "kit:normalize")
	fmt.Fprintln(out, "=============")
	fmt.Fprintf(out, "kit: %s\n", opts.KitPath)
	mode := "dry-run"
	if opts.Apply {
		mode = "apply"
	}
	fmt.Fprintf(out, "mode: %s\n\n", mode)

	renames := domainkit.PlanFileRenames(opts.KitPath)
	renameMap := make(map[string]string, len(renames))
	for _, rename := range renames {
		renameMap[rename.From] = rename.To
	}
	linkUpdates := domainkit.PlanLinkUpdates(opts.KitPath, renameMap)

	if !opts.Apply {
		if len(renames) == 0 && len(linkUpdates) == 0 {
			fmt.Fprintln(out, "No kit normalization needed.")
			return 0
		}
		if len(renames) > 0 {
			fmt.Fprintln(out, "Planned file renames:")
			for _, rename := range renames {
				fmt.Fprintf(out, "- %s -> %s\n",
					repo.ToRepoRelative(opts.RepoRoot, rename.From),
					repo.ToRepoRelative(opts.RepoRoot, rename.To))
			}
		}
		if len(linkUpdates) > 0 {
			fmt.Fprintln(out, "\nPlanned link updates:")
			for _, update := range linkUpdates {
				fmt.Fprintf(out, "- %s: %s -> %s\n",
					repo.ToRepoRelative(opts.RepoRoot, update.File), update.From, update.To)
			}
		}
		return 0
	}

	if repo.HasDirtyWorktree(opts.RepoRoot) {
		fmt.Fprintln(errOut, "Error: apply mode requires a clean working tree. Re-run from a clean tree.")
		return 1
	}

	conflictFrom, conflictTo, err := domainkit.ApplyRenames(renames)
	if err != nil {
		fmt.Fprintf(errOut, "kit normalize error: %v\n", err)
		return 1
	}
	if conflictFrom != "" {
		fmt.Fprintf(errOut, "Error: cannot rename %s because %s already exists.\n",
			repo.ToRepoRelative(opts.RepoRoot, conflictFrom),
			repo.ToRepoRelative(opts.RepoRoot, conflictTo))
		return 1
	}

	changedFiles, err := domainkit.ApplyLinkUpdates(linkUpdates)
	if err != nil {
		fmt.Fprintf(errOut, "kit normalize error: %v\n", err)
		return 1
	}
	fmt.Fprintf(out, "Applied %d rename(s) and %d link update(s) across %d file(s).\n",
		len(renames), len(linkUpdates), changedFiles)
	return 0
}
