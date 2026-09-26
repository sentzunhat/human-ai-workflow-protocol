package embed

import (
	"context"
	"fmt"

	appembed "github.com/sentzunhat/hawp/librarian/src/internal/application/embed"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/model/location"
)

func Run(args []string) error {
	opts, err := parseEmbedArgs(args)
	if err != nil {
		return err
	}

	dir, err := location.Root()
	if err != nil {
		return err
	}

	ctx := context.Background()
	var modelPath string
	if opts.modelRepo == "" {
		modelPath, err = appembed.PullDefaultModel(ctx, dir)
	} else {
		modelPath, err = appembed.PullModel(ctx, opts.modelRepo, opts.onnxFile, dir)
	}
	if err != nil {
		return err
	}

	vectors, err := appembed.Embed(ctx, modelPath, opts.texts)
	if err != nil {
		return err
	}
	for i, vector := range vectors {
		fmt.Printf("%q: %d dims, first 4: %v\n", opts.texts[i], len(vector), vector[:min(4, len(vector))])
	}
	return nil
}
