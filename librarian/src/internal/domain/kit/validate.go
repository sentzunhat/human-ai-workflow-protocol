// Package kit implements the .hawp/kit/ structure validations: file naming,
// required files, and internal links. Ported from
// librarian/scripts/hawp/kit-validate.
package kit

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
			rel, _ := filepath.Rel(kitPath, full)
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
func (s *KitSource) CheckRequiredFiles(kitPath string) []Issue {
	var issues []Issue
	for _, rel := range RequiredFiles {
		if !s.Exists(filepath.Join(kitPath, filepath.FromSlash(rel))) {
			issues = append(issues, Issue{File: rel, Message: "required kit file is missing"})
		}
	}
	return issues
}

// CheckInternalLinks flags relative links in kit markdown (including
// README.md files) whose targets do not exist. Fenced code blocks are
// ignored. readFile is called to load each file's content.
func (s *KitSource) CheckInternalLinks(kitPath string, readFile func(string) ([]byte, error)) []Issue {
	var issues []Issue
	for _, file := range s.FileLister(kitPath, false) {
		raw, err := readFile(file)
		if err != nil {
			continue
		}
		content := s.BlankFences(string(raw))
		rel, _ := filepath.Rel(filepath.Dir(kitPath), file)
		for _, link := range ExtractLinks(content) {
			if !IsLocalHref(link.Href) {
				continue
			}
			pathPart := PathPart(link.Href)
			if pathPart == "" {
				continue
			}
			target := filepath.Join(filepath.Dir(file), pathPart)
			if !s.Exists(target) {
				issues = append(issues, Issue{File: rel, Message: "broken link: " + link.Href})
			}
		}
	}
	return issues
}

// Validate runs all three kit checks and returns the combined issues plus
// the number of checks run. readFile is used for the internal-links check.
func (s *KitSource) Validate(kitPath string, readFile func(string) ([]byte, error)) (issues []Issue, checks int) {
	issues = append(issues, CheckFileNaming(kitPath)...)
	issues = append(issues, s.CheckRequiredFiles(kitPath)...)
	issues = append(issues, s.CheckInternalLinks(kitPath, readFile)...)
	return issues, 3
}

// TrimAnchor is exported for tests; strips "#..." from an href.
func TrimAnchor(href string) string {
	head, _, _ := strings.Cut(href, "#")
	return head
}
