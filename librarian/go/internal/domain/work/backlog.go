package work

import (
	"os"
	"regexp"
	"strings"
)

var separatorRowRe = regexp.MustCompile(`^\|\s*-+`)

// ParseBacklog reads BACKLOG.md and extracts rows from the Active Work,
// Recently Closed (or Done), and Blocked / Parked sections.
func ParseBacklog(backlogPath string) (*Backlog, error) {
	content, err := os.ReadFile(backlogPath)
	if err != nil {
		return nil, err
	}

	backlog := &Backlog{}
	section := ""
	var headerMap map[string]int

	for _, line := range strings.Split(string(content), "\n") {
		switch {
		case strings.Contains(line, "## Active Work"):
			section, headerMap = "active", nil
			continue
		case strings.Contains(line, "## Recently Closed") || strings.Contains(line, "## Done"):
			section, headerMap = "closed", nil
			continue
		case strings.Contains(line, "## Blocked / Parked"):
			section, headerMap = "parked", nil
			continue
		}
		// Any ##+ header ends the current section so nested subsections are
		// not parsed as table rows.
		if section != "" && regexp.MustCompile(`^#{2,}\s`).MatchString(line) {
			section, headerMap = "", nil
			continue
		}

		if section == "" || !strings.HasPrefix(line, "|") {
			continue
		}
		if separatorRowRe.MatchString(line) || strings.Contains(line, "| ---") {
			continue
		}

		cells := parseTableCells(line)
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

	return backlog, nil
}

func parseTableCells(line string) []string {
	parts := strings.Split(line, "|")
	if len(parts) < 3 {
		return nil
	}
	cells := parts[1 : len(parts)-1]
	for i, cell := range cells {
		cells[i] = strings.TrimSpace(cell)
	}
	return cells
}

func mappedCell(cells []string, headerMap map[string]int, aliases ...string) string {
	for _, alias := range aliases {
		if index, ok := headerMap[alias]; ok {
			if index < len(cells) {
				return cells[index]
			}
			return ""
		}
	}
	return ""
}

func stripCodeSpan(value string) string {
	return strings.Trim(value, "`")
}

// buildRow resolves the row ID from the Legacy ID / ID / UUID cells,
// preferring the first cell that parses as a known ID format.
func buildRow(cells []string, headerMap map[string]int) (BacklogRow, bool) {
	var candidates []string
	for _, alias := range []string{"legacy id", "id", "uuid"} {
		value := stripCodeSpan(mappedCell(cells, headerMap, alias))
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
	if normalized == "id" || normalized == "legacy id" {
		return BacklogRow{}, false
	}

	id := ExtractIDFromFilename(rawID)
	if id == "" {
		id = ExtractShortUUID(rawID)
	}
	if id == "" {
		id = rawID
	}

	return BacklogRow{
		ID:     id,
		Type:   mappedCell(cells, headerMap, "type"),
		Title:  mappedCell(cells, headerMap, "title"),
		Status: mappedCell(cells, headerMap, "status", "reason", "closed"),
		Detail: mappedCell(cells, headerMap, "plan file", "detail"),
	}, true
}
