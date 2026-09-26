package search

import (
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/search/benchmark"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/search/query"
)

func RunQuery(args []string, cwd string) error     { return query.Run(args, cwd) }
func RunBenchmark(args []string, cwd string) error { return benchmark.Run(args, cwd) }
