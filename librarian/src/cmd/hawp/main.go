package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		var exitErr cli.ExitError
		if errors.As(err, &exitErr) {
			// Findings were already reported by the command.
			os.Exit(exitErr.Code)
		}
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
