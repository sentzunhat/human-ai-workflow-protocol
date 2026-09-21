// Package markdown contains pure Markdown helpers used by work-record rules.
package markdown

import (
	"regexp"
	"strings"
)

// Link is one Markdown link occurrence with its byte offset.
type Link struct {
	Text   string
	Href   string
	Offset int
	Image  bool
}

var (
	linkRe  = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)
	fenceRe = regexp.MustCompile("(?ms)^```.*?^```")
)

func BlankFences(content string) string {
	return fenceRe.ReplaceAllStringFunc(content, func(m string) string {
		out := []byte(m)
		for i, b := range out {
			if b != '\n' {
				out[i] = ' '
			}
		}
		return string(out)
	})
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
	return href != "" && !strings.HasPrefix(href, "http") &&
		!strings.HasPrefix(href, "/") && !strings.HasPrefix(href, "#")
}

func PathPart(href string) string { part, _, _ := strings.Cut(href, "#"); return part }
