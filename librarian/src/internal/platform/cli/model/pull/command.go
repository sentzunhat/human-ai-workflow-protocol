package pull

import (
	"context"
	"fmt"

	appembed "github.com/sentzunhat/hawp/librarian/src/internal/application/embed"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/model/location"
)

func Run(args []string) error {
	opts, err := parseModelPullArgs(args)
	if err != nil {
		return err
	}

	dir, err := location.Root()
	if err != nil {
		return err
	}
	fmt.Printf("Pulling %s into %s...\n", opts.modelRepo, dir)
	modelPath, err := appembed.PullModel(context.Background(), opts.modelRepo, opts.onnxFile, dir)
	if err != nil {
		return err
	}
	fmt.Printf("Model ready at %s\n", modelPath)
	return nil
}
