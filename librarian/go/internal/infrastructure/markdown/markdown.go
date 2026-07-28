// Package markdown provides shared markdown scanning helpers: file
// collection, fenced-code blanking, and local link extraction.
package markdown

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	linkRe  = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)
	fenceRe = regexp.MustCompile("(?ms)^```.*?^```")
)

// Link is one [text](href) occurrence with its byte offset in the content.
type Link struct {
	Text   string
	Href   string
	Offset int
	Image  bool // preceded by "!", i.e. ![alt](src)
}

// CollectFiles recursively gathers .md files under dir. When skipReadme is
// true, README.md files are excluded. Missing directories yield nil.
func CollectFiles(dir string, skipReadme bool) []string {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		full := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			files = append(files, CollectFiles(full, skipReadme)...)
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		if skipReadme && entry.Name() == "README.md" {
			continue
		}
		files = append(files, full)
	}
	return files
}

// BlankFences replaces fenced code blocks with spaces (newlines preserved)
// so links inside them are not scanned, without shifting offsets.
func BlankFences(content string) string {
	return fenceRe.ReplaceAllStringFunc(content, func(m string) string {
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
