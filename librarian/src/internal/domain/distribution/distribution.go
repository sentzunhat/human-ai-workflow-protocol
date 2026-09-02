package distribution

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Variant struct {
	OutputFile      string
	SectionFiles    []string
	ScriptPartFiles []string
	Ref             string
	Provider        string
	Operation       string
}

type BuildResult struct {
	OutputPath string
	Content    string
}

type PathLeak struct {
	File string
	Line int
	Text string
}

var ActiveProviders = []string{"claude", "codex", "github", "cursor", "continue"}

var LegacyRootGuides = []string{
	"install-main.md",
	"install-dev.md",
	"install-development.md",
	"update-main.md",
	"update-dev.md",
	"update-development.md",
}

var DownstreamTargetFiles = []string{
	"core/.hawp/kit/instructions/da-file-tracking.md",
	"core/.hawp/kit/references/work-item-file-tracking.md",
	"core/.hawp/kit/references/install-update-safety.md",
	"core/.hawp/kit/templates/work-item-files.md",
	"core/.hawp/kit/templates/adr-template.md",
}

var refLinePattern = regexp.MustCompile(`(^|\n)REF="[^"]*"`)
var providerLinePattern = regexp.MustCompile(`(^|\n)PROVIDER="[^"]*"`)

func ComputeExpectedOutputs(repoRoot string) ([]BuildResult, error) {
	sourceRoot := filepath.Join(repoRoot, "distribution")
	outputRoot := filepath.Join(repoRoot, "distribution", "generated")
	variants := distributionPlan()
	outputs := make([]BuildResult, 0, len(variants))

	for _, variant := range variants {
		sections := make([]string, 0, len(variant.SectionFiles)+2)
		for _, sectionFile := range variant.SectionFiles {
			sourcePath := filepath.Join(sourceRoot, filepath.FromSlash(sectionFile))
			content, err := os.ReadFile(sourcePath)
			if err != nil {
				return nil, fmt.Errorf("missing source fragment %s: %w", sourcePath, err)
			}
			sections = append(sections, normalizeForJoin(string(content)))
		}

		scriptSection, err := generateScriptSection(repoRoot, variant)
		if err != nil {
			return nil, err
		}
		refSection := generateSourceReferenceSection(variant)
		joined := strings.Join(append(sections, normalizeForJoin(scriptSection), normalizeForJoin(refSection)), "\n\n")

		outputs = append(outputs, BuildResult{
			OutputPath: filepath.Join(outputRoot, filepath.FromSlash(variant.OutputFile)),
			Content:    ensureTrailingNewline(joined),
		})
	}

	return outputs, nil
}

func FindDownstreamPathLeaks(repoRoot string) ([]PathLeak, error) {
	var leaks []PathLeak
	for _, rel := range DownstreamTargetFiles {
		abs := filepath.Join(repoRoot, filepath.FromSlash(rel))
		body, err := os.ReadFile(abs)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		lines := strings.Split(normalizeForCompare(string(body)), "\n")
		for idx, line := range lines {
			if strings.Contains(line, "core/.hawp/") {
				leaks = append(leaks, PathLeak{
					File: rel,
					Line: idx + 1,
					Text: strings.TrimSpace(line),
				})
			}
		}
	}
	return leaks, nil
}

func normalizeForCompare(content string) string {
	return strings.ReplaceAll(content, "\r\n", "\n")
}

func distributionPlan() []Variant {
	var variants []Variant
	for _, provider := range ActiveProviders {
		for _, operation := range []string{"install", "update"} {
			for _, ref := range []string{"main", "development"} {
				variants = append(variants, providerVariant(provider, operation, ref))
			}
		}
	}
	return variants
}

func providerVariant(provider, operation, ref string) Variant {
	return Variant{
		OutputFile: provider + "/" + operation + "/" + ref + ".md",
		SectionFiles: []string{
			"sources/providers/" + provider + "/preamble-" + operation + ".md",
			"sources/shared/safety.md",
			"sources/providers/" + provider + "/safety.md",
			"sources/shared/repo-boundaries-kit.md",
			"sources/providers/" + provider + "/boundaries.md",
			"sources/shared/" + operation + ".md",
			"sources/providers/" + provider + "/" + operation + "-contract.md",
			"sources/providers/" + provider + "/" + operation + "/" + ref + ".md",
		},
		ScriptPartFiles: []string{
			"distribution/sources/" + operation + "/script-core.md",
			"distribution/sources/providers/" + provider + "/script-" + operation + ".md",
			"distribution/sources/" + operation + "/script-footer.md",
		},
		Ref:       ref,
		Provider:  provider,
		Operation: operation,
	}
}

func generateScriptSection(repoRoot string, variant Variant) (string, error) {
	operation := "Update"
	if variant.Operation == "install" {
		operation = "Install"
	}
	bashBlock, err := composeProviderScript(repoRoot, variant)
	if err != nil {
		return "", err
	}
	return "## " + operation + " Command (Copy/Paste)\n\n" +
		"Run this from the root of your target repository. No edits are required; branch and provider are already configured in the command. Each run fetches the latest commit from that branch.\n\n" +
		bashBlock, nil
}

func composeProviderScript(repoRoot string, variant Variant) (string, error) {
	parts := make([]string, 0, len(variant.ScriptPartFiles))
	for _, part := range variant.ScriptPartFiles {
		body, err := extractBashBody(filepath.Join(repoRoot, filepath.FromSlash(part)))
		if err != nil {
			return "", err
		}
		parts = append(parts, body)
	}
	body := strings.Join(parts, "\n\n")

	if !refLinePattern.MatchString(body) || !providerLinePattern.MatchString(body) {
		return "", fmt.Errorf("composed script missing REF or PROVIDER for %s", variant.OutputFile)
	}

	substituted := refLinePattern.ReplaceAllString(body, "${1}REF=\""+variant.Ref+"\"")
	substituted = providerLinePattern.ReplaceAllString(substituted, "${1}PROVIDER=\""+variant.Provider+"\"")
	return "```bash\n" + substituted + "\n```", nil
}

func extractBashBody(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("script part not found %s: %w", path, err)
	}

	text := normalizeForCompare(string(content))
	start := strings.Index(text, "```bash\n")
	if start == -1 {
		return "", fmt.Errorf("no bash block in %s", path)
	}
	start += len("```bash\n")
	end := strings.Index(text[start:], "\n```")
	if end == -1 {
		return "", fmt.Errorf("no bash block in %s", path)
	}
	return strings.TrimRight(text[start:start+end], "\n"), nil
}

func generateSourceReferenceSection(variant Variant) string {
	fragmentLines := make([]string, 0, len(variant.SectionFiles))
	for _, file := range variant.SectionFiles {
		fragmentLines = append(fragmentLines, "- `distribution/"+file+"`")
	}
	scriptLines := make([]string, 0, len(variant.ScriptPartFiles))
	for _, file := range variant.ScriptPartFiles {
		scriptLines = append(scriptLines, "- `"+file+"`")
	}

	return "## Source Reference\n\n" +
		"This file is generated. Do not edit it directly.\n\n" +
		"- Workflow gate: pushes and pull requests on `main` or `development` fail when generated guides drift from source.\n" +
		"- Local sync: run `hawp distribution sync` after editing `distribution/sources/` or the distribution composition code.\n\n" +
		"Generated output file:\n\n" +
		"- `distribution/generated/" + variant.OutputFile + "`\n\n" +
		"Provider: `" + variant.Provider + "` · Operation: `" + variant.Operation + "` · Branch: `" + variant.Ref + "`\n\n" +
		"Install mapping: `core/providers/." + variant.Provider + "/` -> downstream paths in this guide.\n\n" +
		"This generated guide is built from:\n\n" +
		strings.Join(fragmentLines, "\n") + "\n\n" +
		"Composed shell script (core + provider overlay + footer):\n\n" +
		strings.Join(scriptLines, "\n")
}

func normalizeForJoin(content string) string {
	return strings.TrimRight(normalizeForCompare(content), "\n")
}

func ensureTrailingNewline(content string) string {
	if strings.HasSuffix(content, "\n") {
		return content
	}
	return content + "\n"
}
