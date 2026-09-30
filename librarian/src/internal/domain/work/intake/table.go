package intake

import (
	"fmt"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/table"
)

// InsertNewItem preserves the Active Work table's column order and optional
// columns. Unsupported tables fail before the application creates a plan.
func InsertNewItem(backlog string, item NewItemInput, date string) (string, error) {
	lines := strings.Split(backlog, "\n")
	active := false
	for n, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "## Active Work" {
			active = true
			continue
		}
		if active && strings.HasPrefix(trimmed, "## ") {
			break
		}
		if !active || !strings.HasPrefix(trimmed, "|") {
			continue
		}
		headers := table.Cells(trimmed)
		if n+1 >= len(lines) || !isTableSeparator(lines[n+1], len(headers)) {
			return "", fmt.Errorf("Active Work table is missing its separator row")
		}
		cells := make([]string, len(headers))
		seen := map[string]bool{}
		hasID, hasTitle, hasStatus := false, false, false
		for j, header := range headers {
			key := strings.ToLower(strings.TrimSpace(header))
			if seen[key] {
				return "", fmt.Errorf("duplicate Active Work column %q", header)
			}
			seen[key] = true
			cells[j] = "-"
			switch key {
			case "uuid", "id", "#":
				cells[j], hasID = "`"+item.shortUUID()+"`", true
			case "legacy id":
				cells[j] = "-"
			case "type":
				cells[j] = item.Type
			case "title":
				cells[j], hasTitle = intakeTableText(item.Title), true
			case "status":
				cells[j], hasStatus = "inbox", true
			case "owner":
				cells[j] = "unassigned"
			case "plan file", "detail":
				cells[j] = "[plan](" + item.PlanRelativePath() + ")"
			case "updated":
				cells[j] = date
			}
		}
		if !hasID || !hasTitle || !hasStatus {
			return "", fmt.Errorf("Active Work table requires an ID/UUID/#, Title, and Status column")
		}
		insertAt := n + 2
		for insertAt < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[insertAt]), "|") {
			insertAt++
		}
		row := "| " + strings.Join(cells, " | ") + " |"
		out := append([]string{}, lines[:insertAt]...)
		out = append(out, row)
		out = append(out, lines[insertAt:]...)
		return strings.Join(out, "\n"), nil
	}
	return "", fmt.Errorf("BACKLOG.md has no supported Active Work table")
}

func isTableSeparator(line string, width int) bool {
	cells := table.Cells(strings.TrimSpace(line))
	if width == 0 || len(cells) != width {
		return false
	}
	for _, cell := range cells {
		if strings.Trim(strings.Trim(cell, ":"), "-") != "" || !strings.Contains(cell, "-") {
			return false
		}
	}
	return true
}

func intakeTableText(value string) string {
	return escapeTableCell(value)
}
