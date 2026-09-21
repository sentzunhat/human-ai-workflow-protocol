// Package table contains the small Markdown-table primitives shared by work
// record parsers. It deliberately does not know about backlog semantics.
package table

import "strings"

// Cells splits a pipe-delimited Markdown row and trims each cell.
func Cells(line string) []string {
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

// Cell returns the first mapped cell matching one of the supplied headers.
func Cell(cells []string, headers map[string]int, aliases ...string) string {
	for _, alias := range aliases {
		if index, ok := headers[alias]; ok {
			if index < len(cells) {
				return cells[index]
			}
			return ""
		}
	}
	return ""
}

// StripCodeSpan removes the simple backtick wrapping used in backlog cells.
func StripCodeSpan(value string) string { return strings.Trim(value, "`") }
