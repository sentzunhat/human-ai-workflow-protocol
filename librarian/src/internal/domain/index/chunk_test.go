package index

import (
	"strings"
	"testing"
)

func TestChunkBySection(t *testing.T) {
	content := `# Analysis

This is a short paragraph that should stay together.

## Subsection

A longer section with multiple paragraphs. This paragraph contains enough words to potentially trigger chunking if it gets very long. Let me add more content to make this realistic for testing purposes. The chunking algorithm should respect paragraph boundaries and not split in the middle of a paragraph, even if the word count exceeds the target.

This is another paragraph in the same subsection that should be grouped with the previous content if they fit together.

## Another Subsection

Final section with some content.`

	chunks := ChunkBySection(content)

	if len(chunks) == 0 {
		t.Fatal("expected non-empty chunks")
	}

	// Verify no empty chunks
	for i, chunk := range chunks {
		if strings.TrimSpace(chunk) == "" {
			t.Errorf("chunk %d is empty", i)
		}
	}

	// Verify content is preserved
	joined := strings.Join(chunks, "\n\n")
	if !strings.Contains(joined, "Analysis") || !strings.Contains(joined, "Another Subsection") {
		t.Error("content not preserved after chunking")
	}
}

func TestChunkBySectionPreservesBoundaries(t *testing.T) {
	content := `## Section 1

Paragraph one.

Paragraph two.

## Section 2

Paragraph three.`

	chunks := ChunkBySection(content)

	// Should not split at section boundaries inappropriately
	hasSection1 := false
	hasSection2 := false
	for _, chunk := range chunks {
		if strings.Contains(chunk, "Section 1") {
			hasSection1 = true
		}
		if strings.Contains(chunk, "Section 2") {
			hasSection2 = true
		}
	}

	if !hasSection1 || !hasSection2 {
		t.Error("section headers not preserved in chunks")
	}
}

func TestDeterministicUUID(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{".hawp/kit/start-here.md", "kit-start-here"},
		{".hawp/kit/usage/status-report.md", "kit-usage-status-report"},
		{".hawp/work/active/fbf12a93-plan.md", "work-active-fbf12a93-plan"},
	}

	for _, c := range cases {
		got := DeterministicUUID(c.path)
		if got != c.want {
			t.Errorf("DeterministicUUID(%q) = %q, want %q", c.path, got, c.want)
		}
	}
}

func TestBuildFolderContextKitDoc(t *testing.T) {
	doc := Document{
		Path:       ".hawp/kit/start-here.md",
		Type:       "guide",
		FolderRole: "kit/start-here",
	}
	context := BuildFolderContext(doc, nil)

	if !strings.Contains(context, "kit/start-here") {
		t.Error("folder role not in context")
	}
	if !strings.Contains(context, "guide") {
		t.Error("type not in context")
	}
	if !strings.Contains(context, "start-here.md") {
		t.Error("path not in context")
	}
}

func TestBuildFolderContextWorkDoc(t *testing.T) {
	doc := Document{
		Path:       ".hawp/work/active/fbf12a93-plan.md",
		Type:       "plan",
		FolderRole: "work/active",
	}
	status := "in-progress"
	uuid := "fbf12a93"
	metadata := &DocumentMetadata{
		Status:   status,
		WorkUUID: uuid,
	}
	context := BuildFolderContext(doc, metadata)

	if !strings.Contains(context, "in-progress") {
		t.Error("status not in context")
	}
	if !strings.Contains(context, "fbf12a93") {
		t.Error("work_uuid not in context")
	}
}
