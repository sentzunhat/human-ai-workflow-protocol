package backlog

import (
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/table"
)

var separatorRowRe = regexp.MustCompile(`^\|\s*-+`)
var sectionHeaderRe = regexp.MustCompile(`(?i)^##\s+(active\s+work|blocked\s*/\s*parked|recently\s+closed|done)(?:\s+.*)?$`)
var headingSpaceRe = regexp.MustCompile(`\s+`)

// ParseBacklogMarkdown parses backlog table content without accessing the
// filesystem. Callers that already own the input boundary can use this pure
// domain function directly.
func ParseBacklogMarkdown(content string) *Backlog {

	backlog := &Backlog{}
	section := ""
	var headerMap map[string]int

	for _, line := range strings.Split(string(content), "\n") {
		sectionHeading := ""
		if match := sectionHeaderRe.FindStringSubmatch(strings.TrimSpace(line)); match != nil {
			sectionHeading = strings.ToLower(headingSpaceRe.ReplaceAllString(match[1], " "))
		}
		switch {
		case sectionHeading == "active work":
			section, headerMap = "active", nil
			continue
		case sectionHeading == "recently closed" || sectionHeading == "done":
			section, headerMap = "closed", nil
			continue
		case strings.HasPrefix(sectionHeading, "blocked"):
			section, headerMap = "parked", nil
			continue
		}
		// A new level-2 section switches context; deeper subsection headings keep
		// the current top-level section but reset table headers.
		if strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "##\t") {
			section, headerMap = "", nil
			continue
		}
		if section != "" && regexp.MustCompile(`^#{3,}\s`).MatchString(line) {
			headerMap = nil
			continue
		}

		if section == "" || !strings.HasPrefix(line, "|") {
			continue
		}
		if separatorRowRe.MatchString(line) || strings.Contains(line, "| ---") {
			continue
		}

		cells := table.Cells(line)
		if len(cells) == 0 {
			continue
		}
		if headerMap == nil {
			headerMap = map[string]int{}
			for i, cell := range cells {
				headerMap[strings.ToLower(strings.TrimSpace(cell))] = i
			}
			continue
		}

		row, ok := buildRow(cells, headerMap)
		if !ok {
			continue
		}
		switch section {
		case "active":
			backlog.Active = append(backlog.Active, row)
		case "closed":
			backlog.Closed = append(backlog.Closed, row)
		case "parked":
			backlog.Parked = append(backlog.Parked, row)
		}
	}

	return backlog
}

// buildRow resolves the row ID from the Legacy ID / ID / UUID / # cells,
// preferring the first cell that parses as a known ID format.
func buildRow(cells []string, headerMap map[string]int) (BacklogRow, bool) {
	var candidates []string
	for _, alias := range []string{"legacy id", "id", "uuid", "#"} {
		value := table.StripCodeSpan(table.Cell(cells, headerMap, alias))
		if label, _, ok := table.LinkParts(value); ok {
			value = table.StripCodeSpan(label)
		}
		if value != "" && value != "—" && value != "-" {
			candidates = append(candidates, value)
		}
	}

	rawID := ""
	for _, value := range candidates {
		if ExtractIDFromFilename(value) != "" {
			rawID = value
			break
		}
	}
	if rawID == "" {
		for _, value := range candidates {
			if ExtractShortUUID(value) != "" {
				rawID = value
				break
			}
		}
	}
	if rawID == "" && len(candidates) > 0 {
		rawID = candidates[0]
	}
	if rawID == "" {
		return BacklogRow{}, false
	}
	normalized := strings.ToLower(strings.TrimSpace(rawID))
	if normalized == "id" || normalized == "legacy id" || normalized == "#" {
		return BacklogRow{}, false
	}

	id := ExtractIDFromFilename(rawID)
	if id == "" {
		id = ExtractShortUUID(rawID)
	}
	if id == "" {
		id = rawID
	}

	detail := table.Cell(cells, headerMap, "plan file", "detail")
	if _, target, ok := table.LinkParts(table.Cell(cells, headerMap, "id")); ok && detail == "" {
		detail = "[plan](" + target + ")"
	}
	return BacklogRow{
		ID:     id,
		Type:   table.Cell(cells, headerMap, "type"),
		Title:  table.Cell(cells, headerMap, "title", "work"),
		Status: table.Cell(cells, headerMap, "status", "reason", "closed"),
		Detail: detail,
	}, true
}
