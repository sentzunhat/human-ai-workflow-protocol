package providersync

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FrontmatterField struct {
	Key   string
	Value any
}

type MaterializationTarget struct {
	Behavior    string
	OutputPath  string
	Frontmatter []FrontmatterField
	Preamble    string
}

type MaterializationResult struct {
	OutputPath string
	Content    string
}

const GeneratedBanner = "<!-- Generated from core/providers/shared/behaviors - edit shared source and run `hawp providers sync` -->\n\n"

var MaterializationTargets = []MaterializationTarget{
	{
		Behavior:   "hawp-intake",
		OutputPath: "core/providers/.github/instructions/hawp-intake.instructions.md",
		Frontmatter: []FrontmatterField{
			{Key: "applyTo", Value: ".hawp/**,**/.hawp/**"},
		},
		Preamble: "Use this as a scoped, drop-in instruction when integrating HAWP into repos that already have Copilot instructions and prompts.\n\n",
	},
	{
		Behavior:   "hawp-backlog-alignment",
		OutputPath: "core/providers/.github/instructions/hawp-backlog-alignment.instructions.md",
		Frontmatter: []FrontmatterField{
			{Key: "applyTo", Value: "**"},
		},
	},
	{
		Behavior:   "hawp-docs-alignment",
		OutputPath: "core/providers/.github/instructions/hawp-docs-alignment.instructions.md",
		Frontmatter: []FrontmatterField{
			{Key: "applyTo", Value: "docs/**,README.md,CHANGELOG.md,.hawp/kit/**,.hawp/work/**"},
		},
	},
	{
		Behavior:   "hawp-core",
		OutputPath: "core/providers/.claude/rules/hawp-core.md",
		Frontmatter: []FrontmatterField{
			{Key: "description", Value: "HAWP workflow core - kit entry, status reports, backlog discipline"},
		},
	},
	{
		Behavior:   "hawp-intake",
		OutputPath: "core/providers/.claude/rules/hawp-intake.md",
		Frontmatter: []FrontmatterField{
			{Key: "paths", Value: []string{".hawp/**"}},
		},
	},
	{
		Behavior:   "hawp-backlog-alignment",
		OutputPath: "core/providers/.claude/rules/hawp-backlog-alignment.md",
		Frontmatter: []FrontmatterField{
			{Key: "description", Value: "HAWP backlog alignment and compaction guardrails"},
		},
	},
	{
		Behavior:   "hawp-docs-alignment",
		OutputPath: "core/providers/.claude/rules/hawp-docs-alignment.md",
		Frontmatter: []FrontmatterField{
			{Key: "paths", Value: []string{".hawp/kit/**", ".hawp/work/**"}},
		},
	},
	{
		Behavior:   "hawp-core",
		OutputPath: "core/providers/.cursor/rules/hawp-core.mdc",
		Frontmatter: []FrontmatterField{
			{Key: "description", Value: "HAWP workflow core - kit entry, status reports, backlog discipline"},
			{Key: "alwaysApply", Value: true},
		},
	},
	{
		Behavior:   "hawp-intake",
		OutputPath: "core/providers/.cursor/rules/hawp-intake.mdc",
		Frontmatter: []FrontmatterField{
			{Key: "description", Value: "HAWP modular intake for .hawp work folders"},
			{Key: "globs", Value: "**/.hawp/**"},
			{Key: "alwaysApply", Value: false},
		},
	},
	{
		Behavior:   "hawp-backlog-alignment",
		OutputPath: "core/providers/.cursor/rules/hawp-backlog-alignment.mdc",
		Frontmatter: []FrontmatterField{
			{Key: "description", Value: "HAWP backlog alignment and compaction guardrails"},
			{Key: "alwaysApply", Value: true},
		},
	},
	{
		Behavior:   "hawp-docs-alignment",
		OutputPath: "core/providers/.cursor/rules/hawp-docs-alignment.mdc",
		Frontmatter: []FrontmatterField{
			{Key: "description", Value: "HAWP docs alignment when editing kit or workflow docs"},
			{Key: "globs", Value: ".hawp/kit/**,.hawp/work/**"},
			{Key: "alwaysApply", Value: false},
		},
	},
	{
		Behavior:   "hawp-core",
		OutputPath: "core/providers/.continue/rules/hawp-01-core.md",
		Frontmatter: []FrontmatterField{
			{Key: "name", Value: "HAWP Core"},
			{Key: "description", Value: "HAWP workflow core - kit entry, status reports, backlog discipline"},
			{Key: "alwaysApply", Value: true},
		},
	},
	{
		Behavior:   "hawp-backlog-alignment",
		OutputPath: "core/providers/.continue/rules/hawp-02-backlog-alignment.md",
		Frontmatter: []FrontmatterField{
			{Key: "name", Value: "HAWP Backlog Alignment"},
			{Key: "description", Value: "HAWP backlog alignment and compaction guardrails"},
			{Key: "alwaysApply", Value: true},
		},
	},
	{
		Behavior:   "hawp-intake",
		OutputPath: "core/providers/.continue/rules/hawp-03-intake.md",
		Frontmatter: []FrontmatterField{
			{Key: "name", Value: "HAWP Intake"},
			{Key: "description", Value: "HAWP modular intake for .hawp work folders"},
			{Key: "globs", Value: []string{"**/.hawp/**"}},
			{Key: "alwaysApply", Value: false},
		},
	},
	{
		Behavior:   "hawp-docs-alignment",
		OutputPath: "core/providers/.continue/rules/hawp-04-docs-alignment.md",
		Frontmatter: []FrontmatterField{
			{Key: "name", Value: "HAWP Docs Alignment"},
			{Key: "description", Value: "HAWP docs alignment when editing kit or workflow docs"},
			{Key: "globs", Value: []string{".hawp/kit/**", ".hawp/work/**"}},
			{Key: "alwaysApply", Value: false},
		},
	},
}

func ComputeOutputs(repoRoot string) ([]MaterializationResult, error) {
	sharedRoot := filepath.Join(repoRoot, "core", "providers", "shared", "behaviors")
	results := make([]MaterializationResult, 0, len(MaterializationTargets))

	for _, target := range MaterializationTargets {
		behaviorPath := filepath.Join(sharedRoot, target.Behavior+".md")
		body, err := os.ReadFile(behaviorPath)
		if err != nil {
			return nil, fmt.Errorf("missing shared behavior %s: %w", behaviorPath, err)
		}

		content := strings.Join([]string{
			serializeFrontmatter(target.Frontmatter),
			"",
			strings.TrimRight(GeneratedBanner, "\n"),
			"",
			target.Preamble + trimBody(string(body)),
			"",
		}, "\n")

		results = append(results, MaterializationResult{
			OutputPath: filepath.Join(repoRoot, filepath.FromSlash(target.OutputPath)),
			Content:    content,
		})
	}

	return results, nil
}

func trimBody(body string) string {
	return strings.TrimRight(normalizeForCompare(body), "\n")
}

func normalizeForCompare(content string) string {
	return strings.ReplaceAll(content, "\r\n", "\n")
}

func serializeFrontmatter(fields []FrontmatterField) string {
	lines := []string{"---"}
	for _, field := range fields {
		switch value := field.Value.(type) {
		case []string:
			lines = append(lines, field.Key+":")
			for _, item := range value {
				lines = append(lines, "  - "+strconv.Quote(item))
			}
		case bool:
			lines = append(lines, fmt.Sprintf("%s: %t", field.Key, value))
		case string:
			lines = append(lines, field.Key+": "+strconv.Quote(value))
		default:
			panic(fmt.Sprintf("unsupported frontmatter type for %s", field.Key))
		}
	}
	lines = append(lines, "---")
	return strings.Join(lines, "\n")
}
