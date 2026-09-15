package normalizecmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"

	appwork "github.com/sentzunhat/hawp/librarian/src/internal/application/work/normalize"
)

func parseNormalizeArgs(args []string) (appwork.NormalizeOptions, string, error) {
	opts := appwork.NormalizeOptions{}
	var dryRun bool
	var hawpRoot, workRoot, format string
	flags := flag.NewFlagSet("work normalize", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.BoolVar(&opts.Apply, "apply", false, "apply normalization")
	flags.BoolVar(&dryRun, "dry-run", false, "preview normalization")
	flags.BoolVar(&opts.Validate, "validate", false, "validate workflow")
	flags.BoolVar(&opts.MigrateFolders, "migrate-folders", false, "migrate folder layout")
	flags.BoolVar(&opts.ForceDirty, "force-dirty", false, "allow dirty worktree")
	flags.BoolVar(&opts.Verbose, "verbose", false, "verbose report")
	flags.Bool("no-update-check", false, "suppress update notification")
	flags.StringVar(&format, "format", "text", "text or json")
	flags.StringVar(&opts.Output, "output", "", "report path")
	flags.StringVar(&opts.ExportPlan, "export-plan", "", "plan path")
	flags.StringVar(&opts.ExportResearchQueue, "export-research-queue", "", "research queue path")
	flags.StringVar(&hawpRoot, "hawp-root", "", "HAWP root")
	flags.StringVar(&workRoot, "work-root", "", "work root")
	if err := flags.Parse(args); err != nil {
		return opts, "", fmt.Errorf("work normalize arguments: %w", err)
	}
	if flags.NArg() != 0 {
		return opts, "", fmt.Errorf("unexpected work normalize argument %q", flags.Arg(0))
	}
	var emptyFlag string
	flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "hawp-root", "work-root", "output", "export-plan", "export-research-queue":
			if f.Value.String() == "" {
				emptyFlag = f.Name
			}
		}
	})
	if emptyFlag != "" {
		return opts, "", fmt.Errorf("--%s requires a non-empty path", emptyFlag)
	}
	if opts.Apply && dryRun {
		return opts, "", errors.New("--apply and --dry-run are mutually exclusive")
	}
	if hawpRoot != "" && workRoot != "" {
		return opts, "", errors.New("--hawp-root and --work-root are mutually exclusive")
	}
	if format != "text" && format != "json" {
		return opts, "", errors.New("--format must be text or json")
	}
	opts.FormatJSON = format == "json"
	if hawpRoot != "" {
		return opts, filepath.Dir(filepath.Clean(hawpRoot)), nil
	}
	if workRoot != "" {
		return opts, filepath.Dir(filepath.Dir(filepath.Clean(workRoot))), nil
	}
	return opts, "", nil
}
