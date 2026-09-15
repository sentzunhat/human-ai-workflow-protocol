// source-layout prepares and applies a reviewed, repository-wide Go layout.
package main

import (
	"errors"
	"flag"
	"fmt"
	"hawp-source-layout/internal/migration"
	"os"
)

func main() {
	if err := migration.Run(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		fmt.Fprintln(os.Stderr, "source-layout:", err)
		os.Exit(1)
	}
}
