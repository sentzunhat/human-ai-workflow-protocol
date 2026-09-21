package embed

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

type embedOptions struct {
	texts     []string
	modelRepo string
	onnxFile  string
}

func parseEmbedArgs(args []string) (embedOptions, error) {
	var texts []string
	opts := embedOptions{}

	flags := flag.NewFlagSet("embed", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&opts.modelRepo, "model", "", "Hugging Face model repo (org/repo)")
	flags.StringVar(&opts.onnxFile, "onnx-file", "", "path to ONNX file in repo")
	flags.Bool("no-update-check", false, "suppress update notice")

	// Parse each option with its value, allowing text between options.
	for len(args) > 0 {
		if args[0] == "--" {
			texts = append(texts, args[1:]...)
			break
		}
		if !strings.HasPrefix(args[0], "-") || args[0] == "-" {
			texts = append(texts, args[0])
			args = args[1:]
			continue
		}
		count := 1
		name, _, hasValue := strings.Cut(strings.TrimLeft(args[0], "-"), "=")
		if option := flags.Lookup(name); option != nil && !hasValue {
			boolean, ok := option.Value.(interface{ IsBoolFlag() bool })
			if (!ok || !boolean.IsBoolFlag()) && len(args) > 1 {
				count = 2
			}
		}
		if err := flags.Parse(args[:count]); err != nil {
			return embedOptions{}, fmt.Errorf("embed arguments: %w", err)
		}
		args = args[count:]
	}

	opts.texts = texts
	if len(opts.texts) == 0 {
		return embedOptions{}, fmt.Errorf("usage: hawp embed <text> [<text>...] [--model <hf-org/hf-repo>] [--onnx-file <path>]")
	}

	return opts, nil
}
