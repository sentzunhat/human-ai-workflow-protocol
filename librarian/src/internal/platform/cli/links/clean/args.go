package clean

import (
	"flag"
	"fmt"
	"io"
)

type options struct {
	apply bool
}

func parseArgs(args []string) (options, error) {
	opts := options{}

	flags := flag.NewFlagSet("links clean", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.BoolVar(&opts.apply, "apply", false, "apply link repairs")
	flags.Bool("no-update-check", false, "suppress update notice")

	if err := flags.Parse(args); err != nil {
		return options{}, fmt.Errorf("links clean arguments: %w", err)
	}
	if flags.NArg() != 0 {
		return options{}, fmt.Errorf("unexpected links clean argument %q", flags.Arg(0))
	}

	return opts, nil
}
