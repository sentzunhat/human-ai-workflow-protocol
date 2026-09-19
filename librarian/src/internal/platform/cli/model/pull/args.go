package pull

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

type modelPullOptions struct {
	modelRepo string
	onnxFile  string
}

func parseModelPullArgs(args []string) (modelPullOptions, error) {
	opts := modelPullOptions{}

	flags := flag.NewFlagSet("model pull", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&opts.onnxFile, "onnx-file", "", "path to ONNX file in repo")
	flags.Bool("no-update-check", false, "suppress update notice")

	if err := flags.Parse(args); err != nil {
		return modelPullOptions{}, fmt.Errorf("model pull arguments: %w", err)
	}
	if flags.NArg() == 0 {
		return modelPullOptions{}, fmt.Errorf("usage: hawp model pull <hf-org/hf-repo> [--onnx-file <path-in-repo>]")
	}
	opts.modelRepo = flags.Arg(0)
	if strings.TrimSpace(opts.modelRepo) == "" {
		return modelPullOptions{}, fmt.Errorf("model repository must not be empty")
	}
	// Preserve leading options while accepting the documented repo-first form.
	if err := flags.Parse(flags.Args()[1:]); err != nil {
		return modelPullOptions{}, fmt.Errorf("model pull arguments: %w", err)
	}
	if flags.NArg() != 0 {
		return modelPullOptions{}, fmt.Errorf("unexpected extra argument %q", flags.Arg(0))
	}
	var emptyPath bool
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "onnx-file" && strings.TrimSpace(opts.onnxFile) == "" {
			emptyPath = true
		}
	})
	if emptyPath {
		return modelPullOptions{}, fmt.Errorf("--onnx-file requires a non-empty path")
	}

	return opts, nil
}
