package work

import (
	"strings"
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		title string
		want  string
	}{
		{"Fix the --llm-reshape flag", "fix-the-llm-reshape-flag"},
		{"Add ONNX text2text model", "add-onnx-text2text-model"},
		{"  leading and trailing  ", "leading-and-trailing"},
		{"", "item"},
		{"!!!", "item"},
		{"Already-slugged-title", "already-slugged-title"},
	}
	for _, tt := range tests {
		got := Slugify(tt.title)
		if got != tt.want {
			t.Errorf("Slugify(%q) = %q, want %q", tt.title, got, tt.want)
		}
	}
}

func TestSlugifyCapsLength(t *testing.T) {
	long := strings.Repeat("word ", 30) // 150 chars
	got := Slugify(long)
	if len(got) > 60 {
		t.Errorf("Slugify should cap at 60 chars, got %d: %q", len(got), got)
	}
	if strings.HasSuffix(got, "-") {
		t.Errorf("Slugify should not leave a trailing hyphen after truncation, got %q", got)
	}
}

func TestNewItemInputPlanFileName(t *testing.T) {
	item := NewItemInput{
		UUID:  "abcd1234-5678-90ab-cdef-1234567890ab",
		Type:  "bug",
		Title: "Fix the reshape flag",
		Slug:  "fix-the-reshape-flag",
	}
	want := "abcd1234-fix-the-reshape-flag.md"
	if got := item.PlanFileName(); got != want {
		t.Errorf("PlanFileName() = %q, want %q", got, want)
	}
}

func TestNewItemInputBacklogRow(t *testing.T) {
	item := NewItemInput{
		UUID:  "abcd1234-5678-90ab-cdef-1234567890ab",
		Type:  "bug",
		Title: "Fix the reshape flag",
		Slug:  "fix-the-reshape-flag",
	}
	row := item.BacklogRow("2026-07-26")

	for _, want := range []string{"`abcd1234`", "bug", "Fix the reshape flag", "inbox", "unassigned",
		"[plan](active/abcd1234-fix-the-reshape-flag.md)", "2026-07-26"} {
		if !strings.Contains(row, want) {
			t.Errorf("BacklogRow() = %q, missing %q", row, want)
		}
	}

	// Must match the table's 8-column pipe format.
	cols := strings.Count(row, "|")
	if cols != 9 { // 9 pipes = 8 columns
		t.Errorf("BacklogRow() has %d pipes, want 9 (8 columns): %q", cols, row)
	}
}

func TestNewItemInputPlanFileContent(t *testing.T) {
	item := NewItemInput{
		UUID:  "abcd1234-5678-90ab-cdef-1234567890ab",
		Type:  "bug",
		Title: "Fix the reshape flag",
		Slug:  "fix-the-reshape-flag",
	}
	content := item.PlanFileContent("the --llm-reshape flag is broken", "2026-07-26")

	for _, want := range []string{
		"# Fix the reshape flag",
		"**UUID:** `abcd1234-5678-90ab-cdef-1234567890ab`",
		"**Type:** bug",
		"**Reported:** 2026-07-26",
		"> the --llm-reshape flag is broken",
		"## Intake Summary",
		"## Initial Analysis",
		"## Risk + Review Gate",
		"**Status now:** inbox",
		"work/active/abcd1234-fix-the-reshape-flag.md",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("PlanFileContent() missing %q\ngot:\n%s", want, content)
		}
	}
}

func TestInsertActiveRowAppendsBeforeNextHeading(t *testing.T) {
	backlog := `# Backlog

## Active Work

| UUID | Legacy ID | Type | Title | Status | Owner | Plan File | Updated |
| ---- | --------- | ---- | ----- | ------ | ----- | --------- | ------- |
| ` + "`existing1`" + ` | — | task | Existing item | done | — | [plan](active/existing1.md) | 2026-07-20 |

## Blocked / Parked

nothing here
`

	newRow := "| `newitem1` | — | bug | New item | inbox | unassigned | [plan](active/newitem1-new-item.md) | 2026-07-26 |"

	updated, err := InsertActiveRow(backlog, newRow)
	if err != nil {
		t.Fatalf("InsertActiveRow failed: %v", err)
	}

	if !strings.Contains(updated, newRow) {
		t.Errorf("updated backlog should contain the new row")
	}

	// New row must land inside "## Active Work", before "## Blocked / Parked".
	activeIdx := strings.Index(updated, "## Active Work")
	blockedIdx := strings.Index(updated, "## Blocked / Parked")
	newRowIdx := strings.Index(updated, newRow)
	if !(activeIdx < newRowIdx && newRowIdx < blockedIdx) {
		t.Errorf("new row should be inserted between Active Work and Blocked/Parked headings")
	}

	// Existing row must still be present and precede the new row (appended, not prepended).
	existingIdx := strings.Index(updated, "existing1")
	if existingIdx >= newRowIdx {
		t.Errorf("new row should be appended after existing rows, not before")
	}
}

func TestInsertActiveRowMissingSection(t *testing.T) {
	backlog := "# Backlog\n\nNo active work section here.\n"
	_, err := InsertActiveRow(backlog, "| row |")
	if err == nil {
		t.Error("InsertActiveRow should error when \"## Active Work\" section is missing")
	}
}

func TestInsertActiveRowNoTrailingSection(t *testing.T) {
	// Active Work is the last section in the file (no heading after it).
	backlog := `# Backlog

## Active Work

| UUID | Legacy ID | Type | Title | Status | Owner | Plan File | Updated |
| ---- | --------- | ---- | ----- | ------ | ----- | --------- | ------- |
| ` + "`existing1`" + ` | — | task | Existing item | done | — | [plan](active/existing1.md) | 2026-07-20 |
`

	newRow := "| `newitem1` | — | bug | New item | inbox | unassigned | [plan](active/newitem1-new-item.md) | 2026-07-26 |"
	updated, err := InsertActiveRow(backlog, newRow)
	if err != nil {
		t.Fatalf("InsertActiveRow failed: %v", err)
	}
	if !strings.Contains(updated, newRow) {
		t.Errorf("updated backlog should contain the new row")
	}
	if strings.Index(updated, "existing1") >= strings.Index(updated, newRow) {
		t.Errorf("new row should be appended after existing rows")
	}
}
