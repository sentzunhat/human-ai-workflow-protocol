package work

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleBacklog = `# Backlog

## Active Work

### Primary Lane

| UUID | Legacy ID | Type | Title | Status | Owner | Plan File | Updated |
| ---- | --------- | ---- | ----- | ------ | ----- | --------- | ------- |
| ` + "`39bc92b6`" + ` | — | feature | port phase 1 | in-progress | agent | [plan](active/39bc92b6.md) | 2026-07-20 |
| ` + "`eddd8339`" + ` | — | feature | port phase 2 | approved | agent | active/eddd8339.md | 2026-07-20 |

### Nested Notes

| Not | A | Backlog | Row |
| --- | --- | --- | --- |
| x | y | z | w |

## Blocked / Parked

| ID  | Type | Title | Reason | Detail | Updated |
| --- | ---- | ----- | ------ | ------ | ------- |
| ` + "`bee15107`" + ` | improvement | defer adapters | not needed | [plan](parked/bee15107.md) | 2026-07-06 |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| -- | ---- | ----- | ------ | ------ |
| TASK-086 | improvement | restructure scripts | 2026-07-03 | [plan](closed/2026/07/03/TASK-086.md) |
`

func writeBacklog(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "BACKLOG.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParseBacklogSections(t *testing.T) {
	backlog, err := ParseBacklog(writeBacklog(t, sampleBacklog))
	if err != nil {
		t.Fatal(err)
	}
	if len(backlog.Active) != 2 {
		t.Fatalf("active rows = %d, want 2 (nested subsection must not leak)", len(backlog.Active))
	}
	if backlog.Active[0].ID != "39bc92b6" {
		t.Errorf("active[0].ID = %q", backlog.Active[0].ID)
	}
	if backlog.Active[0].Status != "in-progress" {
		t.Errorf("active[0].Status = %q", backlog.Active[0].Status)
	}
	if backlog.Active[1].Detail != "active/eddd8339.md" {
		t.Errorf("active[1].Detail = %q, want plain plan path", backlog.Active[1].Detail)
	}
	if len(backlog.Parked) != 1 || backlog.Parked[0].Detail != "[plan](parked/bee15107.md)" {
		t.Errorf("parked rows = %+v", backlog.Parked)
	}
	if len(backlog.Closed) != 1 || backlog.Closed[0].ID != "TASK-086" {
		t.Errorf("closed rows = %+v", backlog.Closed)
	}
	// Closed rows map the "closed" column into Status.
	if backlog.Closed[0].Status != "2026-07-03" {
		t.Errorf("closed[0].Status = %q, want date from Closed column", backlog.Closed[0].Status)
	}
}

func TestParseBacklogMissingFile(t *testing.T) {
	if _, err := ParseBacklog(filepath.Join(t.TempDir(), "nope.md")); err == nil {
		t.Fatal("expected error for missing backlog")
	}
}
