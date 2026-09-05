package doccmd

import (
	"flag"
	"fmt"
	"io"
	"strings"

	appdoc "github.com/sentzunhat/hawp/librarian/src/internal/application/work/doc"
)

type docOptions struct {
	title      string
	workItemID string
	hawpRoot   string
}

func parseDocArgs(docType string, args []string) (docOptions, error) {
	flags := flag.NewFlagSet("work "+docType, flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	var opts docOptions
	flags.StringVar(&opts.title, "title", "", "short label for the document")
	flags.StringVar(&opts.workItemID, "work-item", "", "work item ID to link this document to (short or full UUID)")
	flags.StringVar(&opts.hawpRoot, "hawp-root", "", "HAWP directory")
	flags.Bool("no-update-check", false, "suppress update notice")

	if err := flags.Parse(args); err != nil {
		return docOptions{}, usage(docType)
	}
	if flags.NArg() != 0 {
		return docOptions{}, fmt.Errorf("unexpected argument %q", flags.Arg(0))
	}
	if strings.TrimSpace(opts.title) == "" {
		return docOptions{}, usage(docType)
	}
	return opts, nil
}

func usage(docType string) error {
	return fmt.Errorf(
		"usage: hawp work %s --title \"<label>\" [--work-item <uuid>] [--hawp-root <path>]\n  types: %s",
		docType, strings.Join(appdoc.ValidTypes, "|"),
	)
}
