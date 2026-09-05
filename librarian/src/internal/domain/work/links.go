// Package work implements the HAWP work-record validations: backlog
// consistency, closed-task completeness, evidence integrity, verification
// clarity, and dead links. Ported from librarian/scripts/hawp/work-validate.
package work

import (
	"regexp"
	"strings"
)

// Link is one [text](href) occurrence with its byte offset in the content.
type link struct {
	text   string
	href   string
	offset int
	image  bool // preceded by "!", i.e. ![alt](src)
}

var (
	workLinkRe  = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)
	workFenceRe = regexp.MustCompile("(?ms)^```.*?^```")
)

// BlankFences replaces fenced code blocks with spaces (newlines preserved)
// so links inside them are not scanned, without shifting offsets.
func BlankFences(content string) string {
	return workFenceRe.ReplaceAllStringFunc(content, func(m string) string {
		out := []rune(m)
		for i, r := range out {
			if r != '\n' {
				out[i] = ' '
			}
		}
		return string(out)
	})
}

// ExtractLinks returns all markdown links in content (fences already blanked
// by the caller when desired), marking image links.
func extractLinks(content string) []link {
	var links []link
	for _, idx := range workLinkRe.FindAllStringSubmatchIndex(content, -1) {
		start := idx[0]
		links = append(links, link{
			text:   content[idx[2]:idx[3]],
			href:   content[idx[4]:idx[5]],
			offset: start,
			image:  start > 0 && content[start-1] == '!',
		})
	}
	return links
}

// IsLocalHref reports whether href points at a local file path (not an
// external URL, absolute path, or in-page anchor).
func isLocalHref(href string) bool {
	if href == "" {
		return false
	}
	if strings.HasPrefix(href, "http") || strings.HasPrefix(href, "/") || strings.HasPrefix(href, "#") {
		return false
	}
	return true
}

// PathPart strips an anchor suffix from a local href ("a.md#x" → "a.md").
func pathPart(href string) string {
	part, _, _ := strings.Cut(href, "#")
	return part
}
