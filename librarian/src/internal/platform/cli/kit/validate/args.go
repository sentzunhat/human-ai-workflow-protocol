package validate

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

type kitValidateOptions struct {
	kitPath string
}

func parseKitValidateArgs(args []string) (kitValidateOptions, error) {
	opts := kitValidateOptions{}

	flags := flag.NewFlagSet("kit validate", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&opts.kitPath, "kit-path", "", "custom kit directory")
	flags.Bool("no-update-check", false, "suppress update notice")

	if err := flags.Parse(args); err != nil {
		return kitValidateOptions{}, fmt.Errorf("kit validate arguments: %w", err)
	}
	if flags.NArg() != 0 {
		return kitValidateOptions{}, fmt.Errorf("unexpected kit validate argument %q", flags.Arg(0))
	}

	var pathSet bool
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "kit-path" {
			pathSet = true
		}
	})
	if pathSet && strings.TrimSpace(opts.kitPath) == "" {
		return kitValidateOptions{}, fmt.Errorf("--kit-path requires a non-empty path")
	}

	return opts, nil
}
