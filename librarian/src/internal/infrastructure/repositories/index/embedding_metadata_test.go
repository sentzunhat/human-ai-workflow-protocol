package sqlite

import (
	"path/filepath"
	"testing"
)

func openTestDB(t *testing.T) *IndexDB {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema() error = %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestGetEmbeddingMetadataNotSetYet(t *testing.T) {
	db := openTestDB(t)

	_, ok, err := db.GetEmbeddingMetadata()
	if err != nil {
		t.Fatalf("GetEmbeddingMetadata() error = %v", err)
	}
	if ok {
		t.Error("ok should be false when no embedding has happened yet")
	}
}

func TestSetAndGetEmbeddingMetadata(t *testing.T) {
	db := openTestDB(t)

	want := EmbeddingMetadata{Backend: "onnx", Model: "bge-base-en-v1.5", Dim: 768}
	if err := db.SetEmbeddingMetadata(want); err != nil {
		t.Fatalf("SetEmbeddingMetadata() error = %v", err)
	}

	got, ok, err := db.GetEmbeddingMetadata()
	if err != nil {
		t.Fatalf("GetEmbeddingMetadata() error = %v", err)
	}
	if !ok {
		t.Fatal("ok should be true after SetEmbeddingMetadata")
	}
	if got != want {
		t.Errorf("GetEmbeddingMetadata() = %+v, want %+v", got, want)
	}
}

func TestSetEmbeddingMetadataUpserts(t *testing.T) {
	db := openTestDB(t)

	if err := db.SetEmbeddingMetadata(EmbeddingMetadata{Backend: "onnx", Model: "bge-base-en-v1.5", Dim: 768}); err != nil {
		t.Fatalf("first SetEmbeddingMetadata() error = %v", err)
	}
	if err := db.SetEmbeddingMetadata(EmbeddingMetadata{Backend: "ollama", Model: "nomic-embed-text", Dim: 768}); err != nil {
		t.Fatalf("second SetEmbeddingMetadata() error = %v", err)
	}

	got, ok, err := db.GetEmbeddingMetadata()
	if err != nil {
		t.Fatalf("GetEmbeddingMetadata() error = %v", err)
	}
	if !ok {
		t.Fatal("ok should be true")
	}
	want := EmbeddingMetadata{Backend: "ollama", Model: "nomic-embed-text", Dim: 768}
	if got != want {
		t.Errorf("GetEmbeddingMetadata() after second set = %+v, want %+v (upsert, not a second row)", got, want)
	}
}
