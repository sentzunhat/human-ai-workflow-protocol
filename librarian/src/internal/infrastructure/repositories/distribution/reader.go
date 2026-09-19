// Package distribution provides infrastructure-layer adapters for distribution domain operations.
package distribution

import (
	"os"

	domaindistribution "github.com/sentzunhat/hawp/librarian/src/internal/domain/distribution"
)

// ComputeExpectedOutputs reads source fragments with os.ReadFile and delegates
// output composition to the domain.
func ComputeExpectedOutputs(repoRoot string) ([]domaindistribution.BuildResult, error) {
	return domaindistribution.ComputeExpectedOutputs(repoRoot, os.ReadFile)
}

// FindDownstreamPathLeaks reads downstream target files with os.ReadFile and
// delegates path-leak detection to the domain.
func FindDownstreamPathLeaks(repoRoot string) ([]domaindistribution.PathLeak, error) {
	return domaindistribution.FindDownstreamPathLeaks(repoRoot, os.ReadFile)
}
