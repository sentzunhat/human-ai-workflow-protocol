package pull

import (
	"fmt"
	"strings"
)

type modelPullOptions struct {
	modelRepo string
	onnxFile  string
}

func parseModelPullArgs(args []string) (modelPullOptions, error) {
	opts := modelPullOptions{}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--onnx-file":
			i++
			if i >= len(args) || strings.TrimSpace(args[i]) == "" {
				return modelPullOptions{}, fmt.Errorf("--onnx-file requires a non-empty path")
			}
			opts.onnxFile = args[i]
		case strings.HasPrefix(arg, "--onnx-file="):
			value := strings.TrimPrefix(arg, "--onnx-file=")
			if strings.TrimSpace(value) == "" {
				return modelPullOptions{}, fmt.Errorf("--onnx-file requires a non-empty path")
			}
			opts.onnxFile = value
		case arg == "--no-update-check":
			// Accepted for consistency with the other CLI commands. Update
			// notification handling is owned by the command dispatcher.
		case strings.HasPrefix(arg, "-"):
			return modelPullOptions{}, fmt.Errorf("model pull arguments: unknown option %q", arg)
		case opts.modelRepo == "":
			opts.modelRepo = arg
		default:
			return modelPullOptions{}, fmt.Errorf("unexpected extra argument %q", arg)
		}
	}

	if opts.modelRepo == "" {
		return modelPullOptions{}, fmt.Errorf("usage: hawp model pull <hf-org/hf-repo> [--onnx-file <path-in-repo>]")
	}
	if strings.TrimSpace(opts.modelRepo) == "" {
		return modelPullOptions{}, fmt.Errorf("model repository must not be empty")
	}

	return opts, nil
}
