package main

import (
	"fmt"
	"os"

	"github.com/sentzunhat/hawp/golang/internal/platform/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
