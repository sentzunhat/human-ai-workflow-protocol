package enable

import (
	"flag"
	"fmt"
	"io"
)

type options struct {
	logBodies bool
}

func parseArgs(args []string) (options, error) {
	opts := options{}
	flags := flag.NewFlagSet("usage enable", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.BoolVar(&opts.logBodies, "log-bodies", false, "capture raw input/output bodies")
	flags.Bool("no-update-check", false, "suppress update notice")
	if err := flags.Parse(args); err != nil {
		return options{}, err
	}
	if flags.NArg() != 0 {
		return options{}, fmt.Errorf("usage enable accepts no positional arguments, got %d", flags.NArg())
	}
	return opts, nil
}
