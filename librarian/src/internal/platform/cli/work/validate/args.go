package validatecmd

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
)

func parseValidateArgs(args []string) (string, error) {
	flags := flag.NewFlagSet("work validate", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var workRoot, hawpRoot string
	flags.StringVar(&workRoot, "work-root", "", "work directory")
	flags.StringVar(&hawpRoot, "hawp-root", "", "HAWP directory")
	flags.Bool("no-update-check", false, "suppress update notice")
	if err := flags.Parse(args); err != nil {
		return "", err
	}
	if flags.NArg() != 0 {
		return "", fmt.Errorf("unexpected work validate argument %q", flags.Arg(0))
	}
	var roots []string
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "work-root" || f.Name == "hawp-root" {
			roots = append(roots, f.Name)
		}
	})
	if len(roots) > 1 {
		return "", fmt.Errorf("--work-root and --hawp-root are mutually exclusive")
	}
	if len(roots) == 1 {
		if workRoot == "" && hawpRoot == "" {
			return "", fmt.Errorf("--%s requires a non-empty path", roots[0])
		}
		if roots[0] == "hawp-root" {
			return filepath.Join(hawpRoot, "work"), nil
		}
	}
	return workRoot, nil
}
