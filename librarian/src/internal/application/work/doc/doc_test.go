package doc_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sentzunhat/hawp/librarian/src/internal/application/work/doc"
)

func TestCreateWorkDoc_PathShape(t *testing.T) {
	workDir := t.TempDir()

	for _, docType := range doc.ValidTypes {
		r, err := doc.CreateWorkDoc(docType, "my title", workDir, "")
		if err != nil {
			t.Fatalf("type=%s: unexpected error: %v", docType, err)
		}
		if r.UUID == "" {
			t.Errorf("type=%s: empty UUID", docType)
		}
		if _, err := os.Stat(r.FilePath); err != nil {
			t.Errorf("type=%s: file not found at %s: %v", docType, r.FilePath, err)
		}
		if !strings.Contains(r.FilePath, r.UUID) {
			t.Errorf("type=%s: UUID %s not in path %s", docType, r.UUID, r.FilePath)
		}
		if filepath.Base(r.FilePath) != docType+".md" {
			t.Errorf("type=%s: expected filename %s.md, got %s", docType, docType, filepath.Base(r.FilePath))
		}
	}
}

func TestCreateWorkDoc_UUIDUniqueness(t *testing.T) {
	workDir := t.TempDir()
	r1, _ := doc.CreateWorkDoc("status", "first", workDir, "")
	r2, _ := doc.CreateWorkDoc("status", "second", workDir, "")
	if r1.UUID == r2.UUID {
		t.Error("two calls produced identical UUIDs")
	}
}

func TestCreateWorkDoc_WorkItemID_Short(t *testing.T) {
	workDir := t.TempDir()
	r, err := doc.CreateWorkDoc("evidence", "downstream savings", workDir, "288d543c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.UUID != "288d543c" {
		t.Errorf("expected UUID=288d543c, got %s", r.UUID)
	}
	if !strings.Contains(r.FilePath, "288d543c") {
		t.Errorf("work item ID not in path: %s", r.FilePath)
	}
}

func TestCreateWorkDoc_WorkItemID_Full(t *testing.T) {
	workDir := t.TempDir()
	fullUUID := "288d543c-1234-5678-abcd-000000000000"
	r, err := doc.CreateWorkDoc("status", "gate check", workDir, fullUUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Full UUID normalised to 8-char short form for the folder.
	if r.UUID != "288d543c" {
		t.Errorf("expected short UUID=288d543c, got %s", r.UUID)
	}
	if !strings.Contains(r.FilePath, "288d543c") {
		t.Errorf("work item ID not in path: %s", r.FilePath)
	}
}

func TestCreateWorkDoc_TemplateContent(t *testing.T) {
	workDir := t.TempDir()
	r, _ := doc.CreateWorkDoc("evidence", "downstream savings", workDir, "")
	b, err := os.ReadFile(r.FilePath)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	content := string(b)
	for _, want := range []string{"uuid:", "title:", "type: evidence", "date:"} {
		if !strings.Contains(content, want) {
			t.Errorf("template missing %q\ngot:\n%s", want, content)
		}
	}
}

func TestCreateWorkDoc_FolderNames(t *testing.T) {
	workDir := t.TempDir()
	cases := map[string]string{
		"status":   "status",
		"evidence": "evidence",
		"decision": "decisions",
		"note":     "notes",
	}
	for docType, wantDir := range cases {
		r, err := doc.CreateWorkDoc(docType, "t", workDir, "")
		if err != nil {
			t.Fatalf("type=%s: %v", docType, err)
		}
		rel, _ := filepath.Rel(workDir, r.FilePath)
		parts := strings.Split(rel, string(filepath.Separator))
		if parts[0] != wantDir {
			t.Errorf("type=%s: expected folder %q, got %q (full path: %s)", docType, wantDir, parts[0], r.FilePath)
		}
	}
}

func TestCreateWorkDoc_UnknownTypeRejected(t *testing.T) {
	workDir := t.TempDir()
	_, err := doc.CreateWorkDoc("plan", "title", workDir, "")
	if err == nil {
		t.Error("expected error for unknown type, got nil")
	}
}

func TestCreateWorkDoc_EmptyTitleRejected(t *testing.T) {
	workDir := t.TempDir()
	_, err := doc.CreateWorkDoc("status", "", workDir, "")
	if err == nil {
		t.Error("expected error for empty title, got nil")
	}
}
