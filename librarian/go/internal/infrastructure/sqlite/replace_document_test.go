package sqlite

import (
	"path/filepath"
	"testing"
)

func TestReplaceDocumentRollsBackOnChunkFailure(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "index.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.InitSchema(); err != nil {
		t.Fatal(err)
	}

	_, err = db.ReplaceDocument(DocumentReplacement{Category: "kit", Type: "guide", Path: "guide.md", FolderRole: "kit", Chunks: []Chunk{{ChunkIdx: 0, Text: "original"}}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ReplaceDocument(DocumentReplacement{Category: "kit", Type: "changed", Path: "guide.md", FolderRole: "kit", Chunks: []Chunk{{ChunkIdx: 0, Text: "new"}, {ChunkIdx: 0, Text: "duplicate"}}})
	if err == nil {
		t.Fatal("ReplaceDocument() error = nil, want duplicate chunk failure")
	}

	var typ, text string
	if err := db.db.QueryRow(`SELECT d.type, c.text FROM documents d JOIN chunks c ON c.document_id = d.id WHERE d.path = ?`, "guide.md").Scan(&typ, &text); err != nil {
		t.Fatal(err)
	}
	if typ != "guide" || text != "original" {
		t.Fatalf("rolled back state = %q/%q, want guide/original", typ, text)
	}
}
