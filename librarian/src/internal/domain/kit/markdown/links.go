package markdown

import (
	"regexp"
	"strings"
)

var (
	linkRe   = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)
	schemeRe = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)
)

// Link is one [text](href) occurrence with its byte offset in the content.
type Link struct {
	Text   string
	Href   string
	Offset int
	Image  bool
}

func ExtractLinks(content string) []Link {
	var links []Link
	for _, idx := range linkRe.FindAllStringSubmatchIndex(content, -1) {
		start := idx[0]
		links = append(links, Link{
			Text: content[idx[2]:idx[3]], Href: content[idx[4]:idx[5]],
			Offset: start, Image: start > 0 && content[start-1] == '!',
		})
	}
	return links
}

func IsLocalHref(href string) bool {
	if href == "" {
		return false
	}
	if strings.HasPrefix(href, "/") || strings.HasPrefix(href, "#") || schemeRe.MatchString(href) {
		return false
	}
	return true
}

func PathPart(href string) string {
	part, _, _ := strings.Cut(href, "#")
	return part
}
