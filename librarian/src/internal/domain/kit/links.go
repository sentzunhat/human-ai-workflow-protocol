package kit

import (
	"regexp"
	"strings"
)

var linkRe = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)

// Link is one [text](href) occurrence with its byte offset in the content.
type Link struct {
	Text   string
	Href   string
	Offset int
	Image  bool // preceded by "!", i.e. ![alt](src)
}

// ExtractLinks returns all markdown links in content, marking image links.
func ExtractLinks(content string) []Link {
	var links []Link
	for _, idx := range linkRe.FindAllStringSubmatchIndex(content, -1) {
		start := idx[0]
		link := Link{
			Text:   content[idx[2]:idx[3]],
			Href:   content[idx[4]:idx[5]],
			Offset: start,
			Image:  start > 0 && content[start-1] == '!',
		}
		links = append(links, link)
	}
	return links
}

// IsLocalHref reports whether href points at a local file path (not an
// external URL, absolute path, or in-page anchor).
func IsLocalHref(href string) bool {
	if href == "" {
		return false
	}
	if strings.HasPrefix(href, "http") || strings.HasPrefix(href, "/") || strings.HasPrefix(href, "#") {
		return false
	}
	return true
}

// PathPart strips an anchor suffix from a local href ("a.md#x" → "a.md").
func PathPart(href string) string {
	part, _, _ := strings.Cut(href, "#")
	return part
}
