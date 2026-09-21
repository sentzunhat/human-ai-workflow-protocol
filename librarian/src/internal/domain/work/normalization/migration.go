package normalization

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/identity"
	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/markdown"
)

var (
	migrationUUIDFieldRe = regexp.MustCompile("(?i)\\*\\*UUID:\\*\\*\\s*`?([0-9a-f]{8}(?:-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})?)`?")
	markdownLinkRefRe    = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)
	workItemLineRe       = regexp.MustCompile(`(?m)^\*\*Work Item:\*\*\s+.*$`)
	planFileLineRe       = regexp.MustCompile(`(?m)^\*\*Plan file:\*\*\s+.*$`)
)

// MovedPlan records a plan path change for later backlog-link reconciliation.
type MovedPlan struct {
	OldRel string
	NewRel string
}

// CanonicalFolderID derives the stable folder name without touching the filesystem.
// backlogID and filenameID are supplied by the caller because their parsers remain
// compatibility helpers in the parent work package.
func CanonicalFolderID(content, fallback, backlogID, filenameID string) string {
	if m := migrationUUIDFieldRe.FindStringSubmatch(content); m != nil {
		raw := strings.ToLower(strings.TrimSpace(m[1]))
		if identity.IsFullUUID(raw) {
			return raw[:8]
		}
		if identity.ExtractShortUUID(raw) != "" {
			return raw
		}
	}
	if backlogID != "" {
		return backlogID
	}
	if filenameID != "" {
		return filenameID
	}
	base := filepath.Base(fallback)
	if strings.HasSuffix(strings.ToLower(base), ".md") {
		base = strings.TrimSuffix(base, filepath.Ext(base))
	}
	if base == "" {
		return fallback
	}
	return base
}

func rewriteLocalHref(href, oldPath, newPath string) string {
	if !markdown.IsLocalHref(href) {
		return href
	}
	pathPart, anchor, hasAnchor := strings.Cut(href, "#")
	target := filepath.Clean(filepath.Join(filepath.Dir(oldPath), filepath.FromSlash(pathPart)))
	rel, err := filepath.Rel(filepath.Dir(newPath), target)
	if err != nil {
		return href
	}
	rewritten := filepath.ToSlash(rel)
	if hasAnchor {
		rewritten += "#" + anchor
	}
	return rewritten
}

// RewriteLocalHref rewrites a local link relative to a moved markdown file.
func RewriteLocalHref(href, oldPath, newPath string) string {
	return rewriteLocalHref(href, oldPath, newPath)
}

func rewriteMovedMarkdown(content, oldPath, newPath string) string {
	blanked := markdown.BlankFences(content)
	matches := markdownLinkRefRe.FindAllStringSubmatchIndex(blanked, -1)
	if len(matches) == 0 {
		return content
	}
	var out strings.Builder
	last := 0
	for _, idx := range matches {
		hrefStart, hrefEnd := idx[4], idx[5]
		out.WriteString(content[last:hrefStart])
		out.WriteString(rewriteLocalHref(content[hrefStart:hrefEnd], oldPath, newPath))
		last = hrefEnd
	}
	out.WriteString(content[last:])
	return out.String()
}

// RewriteMovedMarkdown rewrites local markdown links while preserving fenced code.
func RewriteMovedMarkdown(content, oldPath, newPath string) string {
	return rewriteMovedMarkdown(content, oldPath, newPath)
}

// NormalizeMovedPlanContent updates links and the plan-file reference after a move.
func NormalizeMovedPlanContent(content, workRel, oldPath, newPath string) string {
	updated := rewriteMovedMarkdown(content, oldPath, newPath)
	return planFileLineRe.ReplaceAllString(updated, "**Plan file:** work/"+filepath.ToSlash(workRel))
}

// NormalizeMovedFilesContent updates links and the work-item reference after a move.
func NormalizeMovedFilesContent(content, planRel, oldPath, newPath string) string {
	updated := rewriteMovedMarkdown(content, oldPath, newPath)
	return workItemLineRe.ReplaceAllString(updated, "**Work Item:** .hawp/work/"+filepath.ToSlash(planRel))
}
