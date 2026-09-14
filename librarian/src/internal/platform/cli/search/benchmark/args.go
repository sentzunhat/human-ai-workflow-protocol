package benchmark

import (
	"flag"
	"fmt"
	"io"
)

type options struct {
	tokenMode         bool
	exportPath        string
	reshapeToken      bool
	reshapeBackend    string
	reshapeModel      string
	reshapeURL        string
	downstreamToken   bool
	downstreamBackend string
	downstreamModel   string
	downstreamURL     string
}

func parseArgs(args []string) (options, error) {
	opts := options{}

	flags := flag.NewFlagSet("search benchmark", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.BoolVar(&opts.tokenMode, "token", false, "run token-savings benchmark")
	flags.StringVar(&opts.exportPath, "export", "", "export path for token benchmark results")
	flags.BoolVar(&opts.reshapeToken, "reshape-token", false, "run reshape token-savings benchmark")
	flags.StringVar(&opts.reshapeBackend, "reshape-backend", "ollama", "reshape backend: ollama or onnx")
	flags.StringVar(&opts.reshapeModel, "reshape-model", "", "model name for reshape backend (default: backend default)")
	flags.StringVar(&opts.reshapeURL, "reshape-url", "", "Ollama server URL (default: http://localhost:11434)")
	flags.BoolVar(&opts.downstreamToken, "downstream-token", false, "run downstream savings benchmark (request + search context vs shaped intake)")
	flags.StringVar(&opts.downstreamBackend, "downstream-backend", "ollama", "reshape backend for downstream benchmark: ollama or onnx")
	flags.StringVar(&opts.downstreamModel, "downstream-model", "", "model name for downstream backend (default: backend default)")
	flags.StringVar(&opts.downstreamURL, "downstream-url", "", "Ollama server URL for downstream benchmark (default: http://localhost:11434)")
	flags.Bool("no-update-check", false, "suppress update notice")

	if err := flags.Parse(args); err != nil {
		return options{}, fmt.Errorf("search benchmark arguments: %w", err)
	}
	if flags.NArg() != 0 {
		return options{}, fmt.Errorf("unexpected search benchmark argument %q", flags.Arg(0))
	}

	return opts, nil
}
