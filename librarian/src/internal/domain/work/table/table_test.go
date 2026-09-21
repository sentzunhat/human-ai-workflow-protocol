package table

import "testing"

func TestCellsAndCell(t *testing.T) {
	cells := Cells("| ID | Title |")
	if len(cells) != 2 || cells[0] != "ID" || cells[1] != "Title" {
		t.Fatalf("Cells() = %#v", cells)
	}
	if got := Cell(cells, map[string]int{"title": 1}, "title"); got != "Title" {
		t.Fatalf("Cell() = %q", got)
	}
	if got := StripCodeSpan("`abc`"); got != "abc" {
		t.Fatalf("StripCodeSpan() = %q", got)
	}
}
