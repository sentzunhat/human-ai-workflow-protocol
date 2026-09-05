package work_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	appwork "github.com/sentzunhat/hawp/librarian/src/internal/application/work/intake"
	reposwork "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repositories/work"
)

const sampleBacklog = `# Backlog

## Active Work

| UUID | Type | Title | Status | Owner | Plan File | Updated |
| ---- | ---- | ----- | ------ | ----- | --------- | ------- |
| ` + "`existing1`" + ` | task | Existing item | done | — | [plan](active/existing1/plan.md) | 2026-07-20 |

## Blocked / Parked

nothing here
`

func setupWorkDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "active"), 0o755); err != nil {
		t.Fatalf("mkdir active: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "BACKLOG.md"), []byte(sampleBacklog), 0o644); err != nil {
		t.Fatalf("write BACKLOG.md: %v", err)
	}
	return dir
}

func TestNewItemCreatesPlanFileAndBacklogRow(t *testing.T) {
	workDir := setupWorkDir(t)

	result, err := appwork.NewItem(workDir, "bug", "Fix the reshape flag", "the --llm-reshape flag is broken")
	if err != nil {
		t.Fatalf("NewItem failed: %v", err)
	}

	if result.UUID == "" {
		t.Error("result.UUID should not be empty")
	}
	if result.Type != "bug" {
		t.Errorf("result.Type = %q, want bug", result.Type)
	}
	if !strings.HasSuffix(result.PlanFilePath, filepath.Join("active", filepath.Base(filepath.Dir(result.PlanFilePath)), "plan.md")) {
		t.Errorf("PlanFilePath should point to an active/<id>/plan.md file, got %q", result.PlanFilePath)
	}

	planBytes, err := os.ReadFile(result.PlanFilePath)
	if err != nil {
		t.Fatalf("plan file should exist: %v", err)
	}
	plan := string(planBytes)
	if !strings.Contains(plan, "Fix the reshape flag") {
		t.Error("plan file should contain the title")
	}
	if !strings.Contains(plan, "the --llm-reshape flag is broken") {
		t.Error("plan file should contain the verbatim input")
	}

	backlogBytes, err := os.ReadFile(filepath.Join(workDir, "BACKLOG.md"))
	if err != nil {
		t.Fatalf("read BACKLOG.md: %v", err)
	}
	backlog := string(backlogBytes)
	if !strings.Contains(backlog, "Fix the reshape flag") {
		t.Error("BACKLOG.md should contain a row for the new item")
	}
	if !strings.Contains(backlog, "existing1") {
		t.Error("BACKLOG.md should still contain the pre-existing row")
	}
	if !strings.Contains(backlog, "inbox") {
		t.Error("new row should have status inbox")
	}
	if !strings.Contains(backlog, "/plan.md") {
		t.Error("new row should link to the folder-based plan path")
	}
}

func TestNewItemRequiresTitle(t *testing.T) {
	workDir := setupWorkDir(t)
	_, err := appwork.NewItem(workDir, "task", "", "")
	if err == nil {
		t.Error("NewItem should fail when title is empty")
	}
}

func TestNewItemDefaultsTypeToTask(t *testing.T) {
	workDir := setupWorkDir(t)
	result, err := appwork.NewItem(workDir, "", "Untyped item", "")
	if err != nil {
		t.Fatalf("NewItem failed: %v", err)
	}
	if result.Type != "task" {
		t.Errorf("result.Type = %q, want task (default)", result.Type)
	}
}

func TestNewItemRejectsUnknownType(t *testing.T) {
	workDir := setupWorkDir(t)
	_, err := appwork.NewItem(workDir, "not-a-real-type", "Some item", "")
	if err == nil {
		t.Error("NewItem should fail for an unrecognized Type value")
	}
}

func TestNewItemDefaultsInputToTitle(t *testing.T) {
	workDir := setupWorkDir(t)
	result, err := appwork.NewItem(workDir, "task", "Some title as input", "")
	if err != nil {
		t.Fatalf("NewItem failed: %v", err)
	}
	plan, err := os.ReadFile(result.PlanFilePath)
	if err != nil {
		t.Fatalf("read plan file: %v", err)
	}
	if !strings.Contains(string(plan), "> Some title as input") {
		t.Error("plan file's Input section should default to the title when --input is omitted")
	}
}

func TestNewItemFailsOnMissingBacklog(t *testing.T) {
	dir := t.TempDir()
	_, err := appwork.NewItem(dir, "task", "Some item", "")
	if err == nil {
		t.Error("NewItem should fail when BACKLOG.md doesn't exist")
	}
}

func TestNewItemRespectsBacklogColumns(t *testing.T) {
	for _, header := range []string{
		"UUID | Type | Title | Status | Owner | Plan File | Updated",
		"UUID | Legacy ID | Type | Title | Status | Owner | Plan File | Updated",
		"Title | Status | ID | Detail | Type | Priority",
		"# | Status | Title | Plan File | Next action",
		"ID | Title | Status",
	} {
		t.Run(header, func(t *testing.T) {
			dir := t.TempDir()
			columns := strings.Split(header, " | ")
			separator := "|" + strings.Repeat(" --- |", len(columns))
			original := "# Backlog\n\n## Active Work\n\n| " + header + " |\n" + separator + "\n\nKeep this note.\n\n## Blocked / Parked\n"
			path := filepath.Join(dir, "BACKLOG.md")
			if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
				t.Fatal(err)
			}
			result, err := appwork.NewItem(dir, "bug", "Fix pipe | and\nnewline", "")
			if err != nil {
				t.Fatal(err)
			}
			backlog, err := reposwork.ReadBacklog(path)
			if err != nil {
				t.Fatal(err)
			}
			if len(backlog.Active) != 1 {
				t.Fatalf("expected one active row, got %+v", backlog.Active)
			}
			row := backlog.Active[0]
			if row.ID != result.UUID[:8] || row.Status != "inbox" || row.Title != "Fix pipe &#124; and<br>newline" {
				t.Fatalf("shifted or malformed row: %+v", row)
			}
			if strings.Contains(header, "Plan File") || strings.Contains(header, "Detail") {
				if !strings.Contains(row.Detail, "active/"+result.UUID[:8]+"/plan.md") {
					t.Fatalf("missing plan link: %+v", row)
				}
			}
			if _, err := os.Stat(result.PlanFilePath); err != nil {
				t.Fatalf("missing UUID plan: %v", err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			text := string(data)
			if strings.Index(text, result.UUID[:8]) > strings.Index(text, "Keep this note.") {
				t.Fatal("row inserted after trailing prose")
			}
			for _, line := range strings.Split(text, "\n") {
				if strings.Contains(line, result.UUID[:8]) && strings.Count(line, "|") != len(columns)+1 {
					t.Fatalf("wrong table width: %s", line)
				}
			}
		})
	}
}

func TestNewItemInvalidTableDoesNotWrite(t *testing.T) {
	for _, table := range []string{
		"",
		"| UUID | Title |\n| --- | --- |\n",
		"| UUID | Title | Status | Plan File |\n",
		"| UUID | Title | Status | Status | Plan File |\n| --- | --- | --- | --- | --- |\n",
	} {
		dir := t.TempDir()
		path := filepath.Join(dir, "BACKLOG.md")
		original := "# Backlog\n\n## Active Work\n\n" + table
		if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := appwork.NewItem(dir, "task", "Example", ""); err == nil {
			t.Fatalf("accepted invalid table: %q", table)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != original {
			t.Fatal("backlog changed on refusal")
		}
		if _, err := os.Stat(filepath.Join(dir, "active")); !os.IsNotExist(err) {
			t.Fatalf("created artifacts on refusal: %v", err)
		}
	}
}
