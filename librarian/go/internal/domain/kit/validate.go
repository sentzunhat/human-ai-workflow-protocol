// Package kit implements pure .hawp/kit structure rules.
package kit

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/kit/source"
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

// CheckFileNaming flags entries that are not lowercase-hyphen named.
func CheckFileNaming(snapshot source.Snapshot) []Issue {
	var issues []Issue
	for _, entry := range snapshot.Entries {
		if entry.Name != "README.md" && !validNameRe.MatchString(entry.Name) {
			issues = append(issues, Issue{File: entry.RelPath, Message: `name should be lowercase-hyphen (got "` + entry.Name + `")`})
		}
	}
	return issues
}

// CheckRequiredFiles flags missing required kit files.
func CheckRequiredFiles(snapshot source.Snapshot) []Issue {
	present := make(map[string]bool, len(snapshot.Files))
	for _, file := range snapshot.Files {
		present[filepath.ToSlash(file.RelPath)] = true
	}
	var issues []Issue
	for _, rel := range RequiredFiles {
		if !present[rel] {
			issues = append(issues, Issue{File: rel, Message: "required kit file is missing"})
		}
	}
	return issues
}

// CheckInternalLinks flags relative links whose targets are absent. The
// adapter supplies links with fenced code blocks already excluded.
func CheckInternalLinks(snapshot source.Snapshot) []Issue {
	present := make(map[string]bool, len(snapshot.Files))
	for _, file := range snapshot.Files {
		present[filepath.ToSlash(file.RelPath)] = true
	}

	var issues []Issue
	for _, file := range snapshot.Files {
		for _, link := range file.Links {
			if !isLocalHref(link.Href) {
				continue
			}
			pathPart := TrimAnchor(link.Href)
			if pathPart == "" {
				continue
			}
			target := filepath.ToSlash(filepath.Clean(filepath.Join(filepath.Dir(file.RelPath), pathPart)))
			if !present[target] {
				issues = append(issues, Issue{File: file.RelPath, Message: "broken link: " + link.Href})
			}
		}
	}
	return issues
}

// Validate runs all three kit checks and returns the combined issues plus the
// number of checks run.
func Validate(snapshot source.Snapshot) (issues []Issue, checks int) {
	issues = append(issues, CheckFileNaming(snapshot)...)
	issues = append(issues, CheckRequiredFiles(snapshot)...)
	issues = append(issues, CheckInternalLinks(snapshot)...)
	return issues, 3
}

// TrimAnchor is exported for tests; strips "#..." from an href.
func TrimAnchor(href string) string {
	head, _, _ := strings.Cut(href, "#")
	return head
}

func isLocalHref(href string) bool {
	return href != "" && !strings.HasPrefix(href, "http") && !strings.HasPrefix(href, "/") && !strings.HasPrefix(href, "#")
}
