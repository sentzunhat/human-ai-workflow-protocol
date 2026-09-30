package newcmd

import (
	"flag"
	"fmt"
	"io"
	"strings"

	domainwork "github.com/sentzunhat/hawp/librarian/src/internal/domain/work"
)

var validItemTypes = []string{
	"task", "bug", "improvement", "feature", "fix",
	"test", "infrastructure", "release", "decision",
}

type newOptions struct {
	title    string
	itemType string
	input    string
	hawpRoot string
}

func parseNewArgs(args []string) (newOptions, error) {
	if len(args) == 0 || strings.HasPrefix(args[0], "--") {
		return newOptions{}, fmt.Errorf("usage: hawp work new \"<title>\" [--type %s] [--input \"<request>\"] [--hawp-root <path>]", strings.Join(validItemTypes, "|"))
	}

	opts := newOptions{title: args[0]}
	flags := flag.NewFlagSet("work new", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&opts.itemType, "type", "task", "work item type")
	flags.StringVar(&opts.input, "input", "", "verbatim request text")
	flags.StringVar(&opts.hawpRoot, "hawp-root", "", "HAWP directory")
	flags.Bool("no-update-check", false, "suppress update notice")
	if err := flags.Parse(args[1:]); err != nil {
		return newOptions{}, fmt.Errorf("work new arguments: %w", err)
	}
	if flags.NArg() != 0 {
		return newOptions{}, fmt.Errorf("unexpected work new argument %q", flags.Arg(0))
	}
	if err := domainwork.ValidateTitle(opts.title); err != nil {
		return newOptions{}, err
	}

	valid := false
	for _, itemType := range validItemTypes {
		if opts.itemType == itemType {
			valid = true
			break
		}
	}
	if !valid {
		return newOptions{}, fmt.Errorf("unknown --type %q (want %s)", opts.itemType, strings.Join(validItemTypes, "|"))
	}

	var rootSet bool
	flags.Visit(func(f *flag.Flag) { rootSet = rootSet || f.Name == "hawp-root" })
	if rootSet && strings.TrimSpace(opts.hawpRoot) == "" {
		return newOptions{}, fmt.Errorf("--hawp-root requires a non-empty path")
	}
	return opts, nil
}
