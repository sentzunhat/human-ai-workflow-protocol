package normalize

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

type kitNormalizeOptions struct {
	kitPath string
	apply   bool
}

func parseKitNormalizeArgs(args []string) (kitNormalizeOptions, error) {
	opts := kitNormalizeOptions{}
	var dryRun bool

	flags := flag.NewFlagSet("kit normalize", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&opts.kitPath, "kit-path", "", "custom kit directory")
	flags.BoolVar(&opts.apply, "apply", false, "apply normalization")
	flags.BoolVar(&dryRun, "dry-run", false, "preview normalization")
	flags.Bool("no-update-check", false, "suppress update notice")

	if err := flags.Parse(args); err != nil {
		return kitNormalizeOptions{}, fmt.Errorf("kit normalize arguments: %w", err)
	}
	if flags.NArg() != 0 {
		return kitNormalizeOptions{}, fmt.Errorf("unexpected kit normalize argument %q", flags.Arg(0))
	}

	if opts.apply && dryRun {
		return kitNormalizeOptions{}, fmt.Errorf("--apply and --dry-run are mutually exclusive")
	}
	var pathSet bool
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "kit-path" {
			pathSet = true
		}
	})
	if pathSet && strings.TrimSpace(opts.kitPath) == "" {
		return kitNormalizeOptions{}, fmt.Errorf("--kit-path requires a non-empty path")
	}

	return opts, nil
}
