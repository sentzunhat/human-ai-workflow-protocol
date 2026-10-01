package query

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"strings"
)

type options struct {
	query       string
	limit       int
	context     bool
	semantic    bool
	format      string
	maxTokens   int
	verbose     bool
	hybridRatio float64
}

// Parsing is independent of repository discovery, storage, and model providers.
func parseArgs(args []string) (options, error) {
	var opts options
	if len(args) == 0 || strings.TrimSpace(args[0]) == "" || strings.HasPrefix(args[0], "--") {
		return opts, errors.New("usage: hawp search <query> [options]; query must precede options")
	}
	opts.query = args[0]
	flags := flag.NewFlagSet("search", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.IntVar(&opts.limit, "limit", 10, "maximum results")
	flags.BoolVar(&opts.context, "context", false, "format context")
	flags.BoolVar(&opts.semantic, "semantic", false, "semantic retrieval")
	flags.StringVar(&opts.format, "format", "markdown", "markdown or json")
	flags.IntVar(&opts.maxTokens, "max-tokens", 2000, "context token budget")
	flags.BoolVar(&opts.verbose, "verbose", false, "token accounting")
	flags.BoolVar(&opts.verbose, "v", false, "token accounting")
	flags.Float64Var(&opts.hybridRatio, "hybrid-ratio", 0.3, "lexical weight")
	flags.Bool("no-update-check", false, "suppress update notification")
	if err := flags.Parse(args[1:]); err != nil {
		return opts, fmt.Errorf("search arguments: %w", err)
	}
	if flags.NArg() != 0 {
		return opts, fmt.Errorf("unexpected search argument %q; quote multi-word queries", flags.Arg(0))
	}
	// Retrieval asks for three times the limit; prevent integer overflow.
	if opts.limit <= 0 || opts.limit > int(^uint(0)>>1)/3 {
		return opts, errors.New("--limit must be positive and small enough for retrieval")
	}
	if opts.maxTokens <= 0 {
		return opts, errors.New("--max-tokens must be positive")
	}
	if opts.format != "markdown" && opts.format != "json" {
		return opts, errors.New("--format must be markdown or json")
	}
	if math.IsNaN(opts.hybridRatio) || math.IsInf(opts.hybridRatio, 0) || opts.hybridRatio < 0 || opts.hybridRatio > 1 {
		return opts, errors.New("--hybrid-ratio must be finite and in [0.0, 1.0]")
	}
	return opts, nil
}
