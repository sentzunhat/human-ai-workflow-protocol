package sqlite

import (
	"database/sql"
	"path/filepath"
	"testing"

	domainindex "github.com/sentzunhat/hawp/librarian/src/internal/domain/index"
)

func TestDomainChunkStoragePreservesContextAndRanges(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "index.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.InitSchema(); err != nil {
		t.Fatal(err)
	}
	id, err := db.InsertDocument("work", "plan", "work/active/example/plan.md", "work/active")
	if err != nil {
		t.Fatal(err)
	}
	empty, context := "", "Status: active"
	for i, value := range []*string{nil, &empty, &context} {
		if err := db.InsertChunk(domainindex.Chunk{DocumentID: id, ChunkIdx: i, Text: "searchable content", FolderContext: value, LineStart: 3, LineEnd: 7}); err != nil {
			t.Fatal(err)
		}
		var got sql.NullString
		var start, end int
		if err := db.db.QueryRow("SELECT folder_context, line_start, line_end FROM chunks WHERE document_id = ? AND chunk_idx = ?", id, i).Scan(&got, &start, &end); err != nil {
			t.Fatal(err)
		}
		if got.Valid != (value != nil) || (value != nil && got.String != *value) || start != 3 || end != 7 {
			t.Fatalf("case %d: context=%+v range=%d:%d", i, got, start, end)
		}
	}
	if err := db.InsertMetadata(domainindex.DocumentMetadata{DocumentID: id, WorkUUID: "example", Status: "active", Owner: &empty}); err != nil {
		t.Fatal(err)
	}
	var owner, closed sql.NullString
	var uuid, status string
	if err := db.db.QueryRow("SELECT work_uuid, status, owner, closed_at FROM documents_metadata WHERE document_id = ?", id).Scan(&uuid, &status, &owner, &closed); err != nil {
		t.Fatal(err)
	}
	if uuid != "example" || status != "active" || !owner.Valid || owner.String != "" || closed.Valid {
		t.Fatalf("metadata: %q %q owner=%+v closed=%+v", uuid, status, owner, closed)
	}
}
