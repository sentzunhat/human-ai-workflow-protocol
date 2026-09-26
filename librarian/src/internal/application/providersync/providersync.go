package providersync

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	domainprovidersync "github.com/sentzunhat/hawp/librarian/src/internal/domain/providersync"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

type MaterializeResult struct {
	Updated   []string
	Unchanged []string
}

type ValidateResult struct {
	Missing []string
	Stale   []string
}

func Materialize(repoRoot string) (MaterializeResult, error) {
	outputs, err := domainprovidersync.ComputeOutputs(repoRoot, os.ReadFile)
	if err != nil {
		return MaterializeResult{}, err
	}

	result := MaterializeResult{}
	for _, output := range outputs {
		if err := validateOutputPath(repoRoot, output.OutputPath); err != nil {
			return result, err
		}
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

		if err := filesystem.AtomicWriteFile(repoRoot, output.OutputPath, []byte(output.Content), 0o644); err != nil {
			return result, err
		}
		result.Updated = append(result.Updated, output.OutputPath)
	}

	return result, nil
}

func validateOutputPath(repoRoot, outputPath string) error {
	if err := filesystem.RejectSymlinkAncestors(repoRoot, outputPath); err != nil {
		return fmt.Errorf("refusing provider materialization path %s: %w", outputPath, err)
	}
	info, err := os.Lstat(outputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat provider materialization path %s: %w", outputPath, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("refusing non-regular provider materialization path %s", outputPath)
	}
	return nil
}

func Validate(repoRoot string) (ValidateResult, error) {
	outputs, err := domainprovidersync.ComputeOutputs(repoRoot, os.ReadFile)
	if err != nil {
		return ValidateResult{}, err
	}

	result := ValidateResult{}
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
	return len(r.Missing) == 0 && len(r.Stale) == 0
}

func normalizeForCompare(content string) string {
	return strings.ReplaceAll(content, "\r\n", "\n")
}
