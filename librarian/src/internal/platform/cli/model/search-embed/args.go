package searchembed

import (
	"flag"
	"fmt"
	"io"
)

type searchEmbedOptions struct {
	backend string
	model   string
}

func parseSearchEmbedArgs(args []string) (searchEmbedOptions, error) {
	opts := searchEmbedOptions{}

	flags := flag.NewFlagSet("search embed", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&opts.backend, "backend", "", "embedding backend: onnx|ollama")
	flags.StringVar(&opts.model, "model", "", "model ID (defaults per backend)")
	flags.Bool("no-update-check", false, "suppress update notice")

	if err := flags.Parse(args); err != nil {
		return searchEmbedOptions{}, fmt.Errorf("search embed arguments: %w", err)
	}
	if flags.NArg() != 0 {
		return searchEmbedOptions{}, fmt.Errorf("unexpected search embed argument %q", flags.Arg(0))
	}

	if opts.backend == "" {
		return searchEmbedOptions{}, fmt.Errorf("--backend is required: hawp search embed --backend onnx|ollama [--model <name>]")
	}
	if opts.backend != "onnx" && opts.backend != "ollama" {
		return searchEmbedOptions{}, fmt.Errorf("unknown --backend %q (want onnx|ollama)", opts.backend)
	}

	return opts, nil
}
