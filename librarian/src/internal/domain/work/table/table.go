// Package table contains the small Markdown-table primitives shared by work
// record parsers. It deliberately does not know about backlog semantics.
package table

import "strings"

// Cells splits a pipe-delimited Markdown row and trims each cell.
func Cells(line string) []string {
	line = strings.TrimSpace(line)
	if len(line) < 2 || line[0] != '|' || line[len(line)-1] != '|' {
		return nil
	}

	var cells []string
	var cell strings.Builder
	for i := 1; i < len(line)-1; i++ {
		if line[i] == '\\' && i+1 < len(line)-1 && (line[i+1] == '\\' || line[i+1] == '|') {
			cell.WriteByte(line[i+1])
			i++
			continue
		}
		if line[i] == '|' {
			cells = append(cells, strings.TrimSpace(cell.String()))
			cell.Reset()
			continue
		}
		cell.WriteByte(line[i])
	}
	cells = append(cells, strings.TrimSpace(cell.String()))
	for i := range cells {
		cells[i] = strings.NewReplacer("&#124;", "|", "&amp;", "&").Replace(cells[i])
	}
	if len(cells) < 1 {
		return nil
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
