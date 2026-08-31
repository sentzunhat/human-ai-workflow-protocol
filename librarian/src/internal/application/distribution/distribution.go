package distribution

import (
	"os"
	"path/filepath"
	"strings"

	domaindistribution "github.com/sentzunhat/hawp/librarian/src/internal/domain/distribution"
)

type BuildResult struct {
	Updated   []string
	Unchanged []string
}

type ValidateResult struct {
	Missing       []string
	Stale         []string
	LegacyPresent []string
	PathLeaks     []domaindistribution.PathLeak
}

func Build(repoRoot string) (BuildResult, error) {
	outputs, err := domaindistribution.ComputeExpectedOutputs(repoRoot)
	if err != nil {
		return BuildResult{}, err
	}

	result := BuildResult{}
	for _, output := range outputs {
		if err := os.MkdirAll(filepath.Dir(output.OutputPath), 0o755); err != nil {
			return result, err
		}

		next := normalizeForCompare(output.Content)
		current, err := os.ReadFile(output.OutputPath)
		if err == nil && normalizeForCompare(string(current)) == next {
			result.Unchanged = append(result.Unchanged, output.OutputPath)
			continue
		}
		if err != nil && !os.IsNotExist(err) {
			return result, err
		}

		if err := os.WriteFile(output.OutputPath, []byte(output.Content), 0o644); err != nil {
			return result, err
		}
		result.Updated = append(result.Updated, output.OutputPath)
	}

	return result, nil
}

func Validate(repoRoot string) (ValidateResult, error) {
	outputs, err := domaindistribution.ComputeExpectedOutputs(repoRoot)
	if err != nil {
		return ValidateResult{}, err
	}

	leaks, err := domaindistribution.FindDownstreamPathLeaks(repoRoot)
	if err != nil {
		return ValidateResult{}, err
	}

	result := ValidateResult{PathLeaks: leaks}
	for _, legacy := range domaindistribution.LegacyRootGuides {
		legacyPath := filepath.Join(repoRoot, "distribution", "generated", legacy)
		if _, err := os.Stat(legacyPath); err == nil {
			result.LegacyPresent = append(result.LegacyPresent, legacyPath)
		}
	}

	for _, output := range outputs {
		current, err := os.ReadFile(output.OutputPath)
		if os.IsNotExist(err) {
			result.Missing = append(result.Missing, output.OutputPath)
			continue
		}
		if err != nil {
			return result, err
		}
		if normalizeForCompare(string(current)) != normalizeForCompare(output.Content) {
			result.Stale = append(result.Stale, output.OutputPath)
		}
	}

	return result, nil
}

func (r ValidateResult) OK() bool {
	return len(r.Missing) == 0 && len(r.Stale) == 0 && len(r.LegacyPresent) == 0 && len(r.PathLeaks) == 0
}

func normalizeForCompare(content string) string {
	return strings.ReplaceAll(content, "\r\n", "\n")
}
