// Package kit implements the .hawp/kit/ structure validations: file naming,
// required files, and internal links. Ported from
// librarian/scripts/hawp/kit-validate.
package kit

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/markdown"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

// Issue is one validation finding against a kit-relative path.
type Issue struct {
	File    string
	Message string
}

// RequiredFiles are the kit files every install must have.
var RequiredFiles = []string{
	"start-here.md",
	"usage/status-report.md",
	"usage/intake-workflow.md",
	"usage/init.md",
	"references/spec.md",
	"references/backlog-alignment.md",
}

var validNameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*$`)

// CheckFileNaming flags entries that are not lowercase-hyphen named
// (README.md is allowed).
func CheckFileNaming(kitPath string) []Issue {
	var issues []Issue
	var walk func(dir string)
	walk = func(dir string) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, entry := range entries {
			full := filepath.Join(dir, entry.Name())
			rel := repo.ToRepoRelative(kitPath, full)
			if entry.Name() != "README.md" && !validNameRe.MatchString(entry.Name()) {
				issues = append(issues, Issue{File: rel, Message: `name should be lowercase-hyphen (got "` + entry.Name() + `")`})
			}
			if entry.IsDir() {
				walk(full)
			}
		}
	}
	walk(kitPath)
	return issues
}

// CheckRequiredFiles flags missing required kit files.
func CheckRequiredFiles(kitPath string) []Issue {
	var issues []Issue
	for _, rel := range RequiredFiles {
		if !repo.Exists(filepath.Join(kitPath, filepath.FromSlash(rel))) {
			issues = append(issues, Issue{File: rel, Message: "required kit file is missing"})
		}
	}
	return issues
}

// CheckInternalLinks flags relative links in kit markdown (including
// README.md files) whose targets do not exist. Fenced code blocks are
// ignored.
func CheckInternalLinks(kitPath string) []Issue {
	var issues []Issue
	for _, file := range markdown.CollectFiles(kitPath, false) {
		raw, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		content := markdown.BlankFences(string(raw))
		rel := repo.ToRepoRelative(kitPath, file)
		for _, link := range markdown.ExtractLinks(content) {
			if !markdown.IsLocalHref(link.Href) {
				continue
			}
			pathPart := markdown.PathPart(link.Href)
			if pathPart == "" {
				continue
			}
			target := filepath.Join(filepath.Dir(file), pathPart)
			if !repo.Exists(target) {
				issues = append(issues, Issue{File: rel, Message: "broken link: " + link.Href})
			}
		}
	}
	return issues
}

// Validate runs all three kit checks and returns the combined issues plus
// the number of checks run.
func Validate(kitPath string) (issues []Issue, checks int) {
	issues = append(issues, CheckFileNaming(kitPath)...)
	issues = append(issues, CheckRequiredFiles(kitPath)...)
	issues = append(issues, CheckInternalLinks(kitPath)...)
	return issues, 3
}

// TrimAnchor is exported for tests; strips "#..." from an href.
func TrimAnchor(href string) string {
	head, _, _ := strings.Cut(href, "#")
	return head
}
