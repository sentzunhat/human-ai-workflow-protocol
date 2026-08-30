package work

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	domainwork "github.com/sentzunhat/hawp/librarian/src/internal/domain/work"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

// NormalizeOptions configures a work normalize run.
type NormalizeOptions struct {
	RepoRoot            string
	Apply               bool
	MigrateFolders      bool
	Validate            bool
	FormatJSON          bool
	Output              string
	ExportPlan          string
	ExportResearchQueue string
	ForceDirty          bool
	Verbose             bool
}

// Normalize runs work-record normalization: dry-run detection by default,
// closed-record normalization with --apply. Returns the exit code.
func Normalize(out, errOut io.Writer, opts NormalizeOptions) int {
	var notices []string
	if opts.Verbose {
		mode := "dry-run"
		if opts.Apply {
			mode = "apply"
		}
		notices = append(notices, fmt.Sprintf("Script options: mode=%s, validate=%v, forceDirty=%v", mode, opts.Validate, opts.ForceDirty))
	}

	workRoot := filepath.Join(opts.RepoRoot, ".hawp", "work")
	backlogPath := filepath.Join(workRoot, "BACKLOG.md")
	backlogRel := repo.ToRepoRelative(opts.RepoRoot, backlogPath)

	if opts.Apply {
		if !opts.ForceDirty && repo.HasDirtyWorktree(opts.RepoRoot) {
			fmt.Fprintln(errOut, "Error: apply mode requires a clean working tree. Re-run with --force-dirty to override.")
			return 2
		}
		result := domainwork.ApplyResult{}
		if opts.MigrateFolders {
			migration, err := domainwork.ApplyWorkItemFolderMigration(opts.RepoRoot)
			if err != nil {
				fmt.Fprintf(errOut, "Script error: %v\n", err)
				return 1
			}
			result.ChangedFiles = append(result.ChangedFiles, migration.ChangedFiles...)
			notices = append(notices, fmt.Sprintf("Applied work-item folder migration to %d file(s).", len(migration.ChangedFiles)))
		} else {
			closed, err := domainwork.ApplyClosedRecordNormalization(opts.RepoRoot)
			if err != nil {
				fmt.Fprintf(errOut, "Script error: %v\n", err)
				return 1
			}
			result = closed
			notices = append(notices, fmt.Sprintf("Applied closed-record normalization to %d file(s).", len(result.ChangedFiles)))
		}
		if n := len(result.SkippedFiles); n > 0 {
			notices = append(notices, fmt.Sprintf("Skipped %d ambiguous legacy file(s) without inferable Backlog ID.", n))
		}
		if n := len(result.ResearchQueue); n > 0 {
			notices = append(notices, fmt.Sprintf("Added %d verification evidence follow-up item(s) inside Verification sections for agent-friendly research handoff.", n))
		}

		stdoutText := "No work-record changes were necessary."
		if len(result.ChangedFiles) > 0 {
			stdoutText = fmt.Sprintf("%d work-record file(s) normalized.", len(result.ChangedFiles))
		}
		if opts.Validate {
			notices = append(notices, validationSummary(workRoot))
		}
		if opts.Output != "" {
			if err := writeOptionalFile(opts.Output, stdoutText); err != nil {
				fmt.Fprintf(errOut, "Script error: %v\n", err)
				return 1
			}
			notices = append(notices, "Report written: "+opts.Output)
			stdoutText = ""
		}
		if opts.ExportResearchQueue != "" {
			if err := exportJSON(opts.ExportResearchQueue, result.ResearchQueue); err != nil {
				fmt.Fprintf(errOut, "Script error: %v\n", err)
				return 1
			}
			notices = append(notices, "Research queue exported: "+opts.ExportResearchQueue)
		}
		if stdoutText != "" {
			fmt.Fprintln(out, stdoutText)
		}
		printNotices(out, notices)
		return 0
	}

	if opts.ForceDirty {
		notices = append(notices, "--force-dirty has no effect in dry-run mode.")
	}

	backlog, err := domainwork.ParseNormalizeBacklog(backlogPath)
	if err != nil {
		fmt.Fprintf(errOut, "Script error: %v\n", err)
		return 1
	}
	scan := domainwork.ScanPlanFiles(workRoot)
	operations := domainwork.EvaluateRules(opts.RepoRoot, workRoot, backlogRel, backlog, scan)
	report := domainwork.BuildDetectionReport(
		time.Now().UTC().Format(time.RFC3339), backlogRel,
		len(scan.Files), len(backlog.Rows), operations)
	report.ResearchQueue = domainwork.BuildResearchQueue(opts.RepoRoot)

	var rendered string
	if opts.FormatJSON {
		rendered, err = domainwork.RenderJSONReport(report)
		if err != nil {
			fmt.Fprintf(errOut, "Script error: %v\n", err)
			return 1
		}
	} else {
		rendered = domainwork.RenderTextReport(report) + "\n"
	}

	if opts.Validate {
		notices = append(notices, validationSummary(workRoot))
	}
	if opts.Output != "" {
		if err := writeOptionalFile(opts.Output, rendered); err != nil {
			fmt.Fprintf(errOut, "Script error: %v\n", err)
			return 1
		}
		notices = append(notices, "Report written: "+opts.Output)
		rendered = ""
	}
	if opts.ExportResearchQueue != "" {
		if err := exportJSON(opts.ExportResearchQueue, report.ResearchQueue); err != nil {
			fmt.Fprintf(errOut, "Script error: %v\n", err)
			return 1
		}
		notices = append(notices, "Research queue exported: "+opts.ExportResearchQueue)
	}
	if opts.ExportPlan != "" {
		if err := exportJSON(opts.ExportPlan, report.Plan); err != nil {
			fmt.Fprintf(errOut, "Script error: %v\n", err)
			return 1
		}
		notices = append(notices, "Plan exported: "+opts.ExportPlan)
	}

	if rendered != "" {
		fmt.Fprint(out, rendered)
	}
	printNotices(out, notices)
	return 0
}

func printNotices(out io.Writer, notices []string) {
	for _, notice := range notices {
		fmt.Fprintln(out, notice)
	}
}

func validationSummary(workRoot string) string {
	report, err := Validate(workRoot)
	if err != nil {
		return "Validation warning: could not parse BACKLOG.md for workflow validation."
	}
	summary := fmt.Sprintf("Validation summary: VALIDATION %s (%d issues, %d warnings).",
		report.Overall, report.Failed, report.Warnings)
	if report.Backlog.Status == domainwork.StatusFail {
		summary += "\nValidation warning: working files are out of sync with BACKLOG.md; reconcile active/closed plan files before kit sync."
	}
	return summary
}

func writeOptionalFile(path, content string) error {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		return err
	}
	return os.WriteFile(absolute, []byte(content), 0o644)
}

func exportJSON(path string, value any) error {
	rendered, err := domainwork.RenderJSONValue(value)
	if err != nil {
		return err
	}
	return writeOptionalFile(path, rendered)
}
