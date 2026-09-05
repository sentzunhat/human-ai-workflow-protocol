package kit

import (
	normalize "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/kit/normalize"
	validate "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/kit/validate"
)

func RunValidate(args []string, cwd string) error  { return validate.Run(args, cwd) }
func RunNormalize(args []string, cwd string) error { return normalize.Run(args, cwd) }
