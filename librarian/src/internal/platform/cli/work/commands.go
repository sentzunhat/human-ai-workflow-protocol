// Package work exposes the CLI entrypoints for HAWP work commands.
package work

import (
	doccmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/work/doc"
	newcmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/work/new"
	normalizecmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/work/normalize"
	validatecmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/work/validate"
)

// RunNew delegates to the new command adapter.
func RunNew(args []string, cwd string) error {
	return newcmd.Run(args, cwd)
}

// RunNormalize delegates to the normalize command adapter.
func RunNormalize(args []string, cwd string) error {
	return normalizecmd.Run(args, cwd)
}

// RunValidate delegates to the validate command adapter.
func RunValidate(args []string, cwd string) error {
	return validatecmd.Run(args, cwd)
}

// RunDoc delegates to the doc command adapter for secondary document types.
func RunDoc(docType string, args []string, cwd string) error {
	return doccmd.Run(docType, args, cwd)
}
