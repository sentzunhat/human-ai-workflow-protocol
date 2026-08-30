package work

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
	if err := os.MkdirAll(filepath.Join(dir, "active"), 0755); err != nil {
		t.Fatalf("mkdir active: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "BACKLOG.md"), []byte(sampleBacklog), 0644); err != nil {
		t.Fatalf("write BACKLOG.md: %v", err)
	}
	return dir
}

func TestNewItemCreatesPlanFileAndBacklogRow(t *testing.T) {
	workDir := setupWorkDir(t)

	result, err := NewItem(workDir, "bug", "Fix the reshape flag", "the --llm-reshape flag is broken")
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
	_, err := NewItem(workDir, "task", "", "")
	if err == nil {
		t.Error("NewItem should fail when title is empty")
	}
}

func TestNewItemDefaultsTypeToTask(t *testing.T) {
	workDir := setupWorkDir(t)
	result, err := NewItem(workDir, "", "Untyped item", "")
	if err != nil {
		t.Fatalf("NewItem failed: %v", err)
	}
	if result.Type != "task" {
		t.Errorf("result.Type = %q, want task (default)", result.Type)
	}
}

func TestNewItemRejectsUnknownType(t *testing.T) {
	workDir := setupWorkDir(t)
	_, err := NewItem(workDir, "not-a-real-type", "Some item", "")
	if err == nil {
		t.Error("NewItem should fail for an unrecognized Type value")
	}
}

func TestNewItemDefaultsInputToTitle(t *testing.T) {
	workDir := setupWorkDir(t)
	result, err := NewItem(workDir, "task", "Some title as input", "")
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
	dir := t.TempDir() // no BACKLOG.md written
	_, err := NewItem(dir, "task", "Some item", "")
	if err == nil {
		t.Error("NewItem should fail when BACKLOG.md doesn't exist")
	}
}
