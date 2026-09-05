package models_test

import (
	"testing"

	models "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/models"
)

func TestNewEmbedderRejectsUnknownBackend(t *testing.T) {
	_, err := models.NewEmbedder("unsupported-backend", "some-model")
	if err == nil {
		t.Error("unsupported backend should return error")
	}
}

func TestNewLLMClientRejectsUnknownBackend(t *testing.T) {
	_, err := models.NewLLMClient("unsupported-backend", "some-model")
	if err == nil {
		t.Error("unsupported backend should return error")
	}
}

func TestNewEmbedderNoneBackend(t *testing.T) {
	e, err := models.NewEmbedder("none", "")
	if err != nil {
		t.Fatalf("none backend should not error: %v", err)
	}
	defer e.Close()
	if e.Backend() != "none" {
		t.Errorf("expected backend=none, got %s", e.Backend())
	}
}

func TestNewLLMClientNoneBackend(t *testing.T) {
	c, err := models.NewLLMClient("none", "")
	if err != nil {
		t.Fatalf("none backend should not error: %v", err)
	}
	defer c.Close()
	if c.Backend() != "none" {
		t.Errorf("expected backend=none, got %s", c.Backend())
	}
}
