package work

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// NormalizeSection identifies which backlog section a row came from.
type NormalizeSection string

const (
	SectionActive  NormalizeSection = "active"
	SectionBlocked NormalizeSection = "blocked"
	SectionClosed  NormalizeSection = "recently-closed"
	SectionOther   NormalizeSection = "other"
)

// NormalizeRow is a parsed backlog row with position info for rule checks.
type NormalizeRow struct {
	Section    NormalizeSection
	LineNumber int
	ID         string
	Type       string
	Title      string
	Status     string
	Updated    string
	Reason     string
	PlanPath   string
}

// NormalizeBacklog is the detection-oriented backlog parse result.
type NormalizeBacklog struct {
	BacklogPath     string
	Rows            []NormalizeRow
	SectionPresence map[NormalizeSection]bool
}

var mdLinkTargetRe = regexp.MustCompile(`\(([^)]+)\)`)

func extractMarkdownLinkTarget(value string) string {
	if m := mdLinkTargetRe.FindStringSubmatch(value); m != nil {
		return strings.TrimSpace(m[1])
	}
	return strings.TrimSpace(stripCodeSpan(value))
}

func normalizeSectionHeader(line string) NormalizeSection {
	switch {
	case strings.HasPrefix(line, "## Active Work"):
		return SectionActive
	case strings.HasPrefix(line, "## Blocked / Parked"):
		return SectionBlocked
	case strings.HasPrefix(line, "## Recently Closed"):
		return SectionClosed
	default:
		return SectionOther
	}
}

// ParseNormalizeBacklog reads BACKLOG.md for detection: every table row with
// its line number, section, and mapped cells.
func ParseNormalizeBacklog(backlogPath string) (*NormalizeBacklog, error) {
	content, err := os.ReadFile(backlogPath)
	if err != nil {
		return nil, err
	}

	result := &NormalizeBacklog{
		BacklogPath: backlogPath,
		SectionPresence: map[NormalizeSection]bool{
			SectionActive: false, SectionBlocked: false,
			SectionClosed: false, SectionOther: false,
		},
	}

	section := SectionOther
	var headerMap map[string]int

	lines := regexp.MustCompile(`\r?\n`).Split(string(content), -1)
	for index, line := range lines {
		if strings.HasPrefix(line, "## ") {
			section = normalizeSectionHeader(line)
			headerMap = nil
			result.SectionPresence[section] = true
			continue
		}
		if !strings.HasPrefix(line, "|") {
			continue
		}
		cells := parseTableCells(line)
		if len(cells) < 4 {
			continue
		}
		if headerMap == nil {
			headerMap = map[string]int{}
			for i, cell := range cells {
				headerMap[strings.ToLower(strings.TrimSpace(cell))] = i
			}
			continue
		}
		if separatorRowRe.MatchString(line) || strings.Contains(line, "| --------") {
			continue
		}

		var candidates []string
		for _, alias := range []string{"legacy id", "id", "uuid", "#"} {
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
		lowered := strings.ToLower(strings.TrimSpace(rawID))
		if rawID == "" || lowered == "id" || lowered == "legacy id" || lowered == "#" {
			continue
		}
		id := ExtractIDFromFilename(rawID)
		if id == "" {
			id = ExtractShortUUID(rawID)
		}
		if id == "" {
			id = rawID
		}

		row := NormalizeRow{
			Section:    section,
			LineNumber: index + 1,
			ID:         id,
			Type:       mappedCell(cells, headerMap, "type"),
			Title:      mappedCell(cells, headerMap, "title"),
		}
		switch section {
		case SectionActive:
			row.Status = mappedCell(cells, headerMap, "status")
			row.PlanPath = extractMarkdownLinkTarget(mappedCell(cells, headerMap, "plan file", "detail"))
			row.Updated = mappedCell(cells, headerMap, "updated")
		case SectionBlocked:
			row.Reason = mappedCell(cells, headerMap, "reason", "status")
			row.PlanPath = extractMarkdownLinkTarget(mappedCell(cells, headerMap, "detail", "plan file"))
			row.Updated = mappedCell(cells, headerMap, "updated")
		case SectionClosed:
			row.Updated = mappedCell(cells, headerMap, "closed", "updated")
			row.PlanPath = extractMarkdownLinkTarget(mappedCell(cells, headerMap, "detail", "plan file"))
		default:
			row.Status = mappedCell(cells, headerMap, "status", "reason", "closed")
			row.PlanPath = extractMarkdownLinkTarget(mappedCell(cells, headerMap, "detail", "plan file"))
			row.Updated = mappedCell(cells, headerMap, "updated", "closed")
		}
		result.Rows = append(result.Rows, row)
	}

	return result, nil
}

// PlanFileRecord is one scanned plan file with its resolved ID.
type PlanFileRecord struct {
	ID      string
	Path    string
	Content string
}

// PlanScan holds the scanned plan files and workspace shape.
type PlanScan struct {
	Files             []PlanFileRecord
	ByID              map[string][]string
	DirectoryPresence map[string]bool // active, parked, closed
}

var (
	backlogIDFieldRe = regexp.MustCompile(`(?i)\*\*Backlog ID(?: \(Legacy\))?:\*\*\s*([A-Z]+-\d+)`)
	uuidFieldRe      = regexp.MustCompile("(?i)\\*\\*UUID:\\*\\*\\s*`?([0-9a-f]{8}(?:-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})?)`?")
)

func readBacklogID(content string) string {
	if m := backlogIDFieldRe.FindStringSubmatch(content); m != nil {
		return strings.ToUpper(m[1])
	}
	if m := uuidFieldRe.FindStringSubmatch(content); m != nil {
		if id := ExtractIDFromFilename(m[1]); id != "" {
			return id
		}
		if id := ExtractShortUUID(m[1]); id != "" {
			return id
		}
	}
	return ""
}

// walkPlanMarkdown mirrors lib.walkMarkdownFiles: recursive, skips README.md.
func walkPlanMarkdown(dir string) []string {
	var out []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		full := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			out = append(out, walkPlanMarkdown(full)...)
			continue
		}
		if strings.HasSuffix(entry.Name(), ".md") && entry.Name() != "README.md" {
			out = append(out, full)
		}
	}
	return out
}

// ScanPlanFiles reads all plan files under active/, parked/, and closed/,
// resolving each file's work item ID from its front matter or filename.
func ScanPlanFiles(workRoot string) *PlanScan {
	scan := &PlanScan{ByID: map[string][]string{}, DirectoryPresence: map[string]bool{}}
	dirs := map[string]string{
		"active": filepath.Join(workRoot, "active"),
		"parked": filepath.Join(workRoot, "parked"),
		"closed": filepath.Join(workRoot, "closed"),
	}
	var paths []string
	for _, name := range []string{"active", "parked", "closed"} {
		dir := dirs[name]
		_, err := os.Stat(dir)
		scan.DirectoryPresence[name] = err == nil
		paths = append(paths, walkPlanMarkdown(dir)...)
	}

	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		content := string(raw)
		id := readBacklogID(content)
		if id == "" {
			base := filepath.Base(path)
			id = ExtractIDFromFilename(base)
			if id == "" {
				id = ExtractShortUUID(strings.TrimSuffix(base, ".md"))
			}
		}
		if id == "" {
			continue
		}
		scan.Files = append(scan.Files, PlanFileRecord{ID: id, Path: path, Content: content})
		scan.ByID[id] = append(scan.ByID[id], path)
	}
	return scan
}
