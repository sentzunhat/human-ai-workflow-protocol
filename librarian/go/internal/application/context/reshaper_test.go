package context

import (
	"context"
	"strings"
	"testing"
)

func TestNewContextReshaper(t *testing.T) {
	// Note: ONNX LLM is blocked on hugot's CGO ORT backend (see
	// llm.ErrGenerativeRequiresORT) — use Ollama for the LLM backend in tests.
	config := ReshapingConfig{
		EmbeddingsBackend: "onnx",
		EmbeddingsModel:   "all-MiniLM-L6-v2",
		LLMBackend:        "ollama",
		LLMModel:          "mistral",
		MaxTokens:         256,
		TopK:              5,
	}

	reshaper, err := NewContextReshaper(config)
	if err != nil {
		// Ollama might not be running - skip this test
		if strings.Contains(err.Error(), "Ollama") || strings.Contains(err.Error(), "connect") {
			t.Skipf("Skipping: Ollama required for test")
		}
		t.Fatalf("NewContextReshaper failed: %v", err)
	}
	defer reshaper.Close()

	if reshaper.config.TopK != 5 {
		t.Errorf("config.TopK should be 5, got %d", reshaper.config.TopK)
	}
	if reshaper.config.MaxTokens != 256 {
		t.Errorf("config.MaxTokens should be 256, got %d", reshaper.config.MaxTokens)
	}
}

func TestNewContextReshaperDefaults(t *testing.T) {
	// Test that defaults are set correctly
	config := ReshapingConfig{
		EmbeddingsBackend: "onnx",
		EmbeddingsModel:   "all-MiniLM-L6-v2",
		LLMBackend:        "ollama",
		LLMModel:          "mistral",
		// TopK and MaxTokens not set
	}

	reshaper, err := NewContextReshaper(config)
	if err != nil {
		if strings.Contains(err.Error(), "Ollama") || strings.Contains(err.Error(), "connect") {
			t.Skipf("Skipping: Ollama required for test")
		}
		t.Fatalf("NewContextReshaper failed: %v", err)
	}
	defer reshaper.Close()

	// Should have defaults
	if reshaper.config.TopK != 5 {
		t.Errorf("default TopK should be 5, got %d", reshaper.config.TopK)
	}
	if reshaper.config.MaxTokens != 512 {
		t.Errorf("default MaxTokens should be 512, got %d", reshaper.config.MaxTokens)
	}
}

func TestSplitIntoSentences(t *testing.T) {
	reshaper := &ContextReshaper{
		config: ReshapingConfig{TopK: 5},
	}

	content := "This is a sentence. Here is another one! And a third question? Short. Not this."
	sentences := reshaper.splitIntoSentences(content)

	// Should have 3 sentences (short ones filtered out)
	if len(sentences) != 3 {
		t.Errorf("expected 3 sentences, got %d: %v", len(sentences), sentences)
	}

	expected := []string{"This is a sentence", "Here is another one", "And a third question"}
	for i, exp := range expected {
		if i < len(sentences) && sentences[i] != exp {
			t.Errorf("sentence %d: expected %q, got %q", i, exp, sentences[i])
		}
	}
}

func TestExtractConceptsFromSentences(t *testing.T) {
	reshaper := &ContextReshaper{
		config: ReshapingConfig{TopK: 5},
	}

	sentences := []string{
		"Python is a programming language",
		"Kubernetes manages containers",
		"Docker provides containerization",
	}

	// Create mock embeddings (all zeros for simplicity)
	embeddings := make([][]float32, len(sentences))
	for i := range embeddings {
		embeddings[i] = make([]float32, 384)
	}

	concepts := reshaper.extractConceptsFromSentences(sentences, embeddings)

	// Should extract capitalized words (Python, Kubernetes, Docker)
	conceptTexts := make(map[string]bool)
	for _, c := range concepts {
		conceptTexts[c.Text] = true
	}

	expected := []string{"Python", "Kubernetes", "Docker"}
	for _, exp := range expected {
		if !conceptTexts[exp] {
			t.Errorf("expected concept %q not found in %v", exp, conceptTexts)
		}
	}
}

func TestDeduplicateConcepts(t *testing.T) {
	reshaper := &ContextReshaper{
		config: ReshapingConfig{TopK: 5},
	}

	concepts := []Concept{
		{Text: "Python", Relevance: 0.9},
		{Text: "python", Relevance: 0.8}, // Should be deduplicated
		{Text: "Kubernetes", Relevance: 0.85},
		{Text: "Docker", Relevance: 0.8},
	}

	unique := reshaper.deduplicateConcepts(concepts)

	// Should have 3 unique concepts (python deduplicated)
	if len(unique) != 3 {
		t.Errorf("expected 3 unique concepts, got %d", len(unique))
	}
}

func TestRankConceptsByRelevance(t *testing.T) {
	reshaper := &ContextReshaper{
		config: ReshapingConfig{TopK: 5},
	}

	concepts := []Concept{
		{Text: "C", Relevance: 0.5},
		{Text: "A", Relevance: 0.9},
		{Text: "B", Relevance: 0.7},
	}

	reshaper.rankConceptsByRelevance(concepts)

	// Should be ranked by relevance (descending)
	expectedOrder := []string{"A", "B", "C"}
	for i, exp := range expectedOrder {
		if concepts[i].Text != exp {
			t.Errorf("position %d: expected %s, got %s", i, exp, concepts[i].Text)
		}
	}
}

func TestBuildReshapingPrompt(t *testing.T) {
	reshaper := &ContextReshaper{
		config: ReshapingConfig{TopK: 5},
	}

	content := "This is test content about Python and Kubernetes."
	concepts := []Concept{
		{Text: "Python", Relevance: 0.9},
		{Text: "Kubernetes", Relevance: 0.85},
	}

	prompt := reshaper.buildReshapingPrompt(content, concepts)

	// Verify prompt contains key parts
	if !strings.Contains(prompt, "Python") {
		t.Error("prompt should contain concept 'Python'")
	}
	if !strings.Contains(prompt, "Kubernetes") {
		t.Error("prompt should contain concept 'Kubernetes'")
	}
	if !strings.Contains(prompt, content) {
		t.Error("prompt should contain original content")
	}
	if !strings.Contains(prompt, "Reshaped content") {
		t.Error("prompt should contain reshaping instruction")
	}
}

func TestReshapeEmptyBlock(t *testing.T) {
	config := ReshapingConfig{
		EmbeddingsBackend: "onnx",
		EmbeddingsModel:   "all-MiniLM-L6-v2",
		LLMBackend:        "ollama",
		LLMModel:          "mistral",
	}

	reshaper, err := NewContextReshaper(config)
	if err != nil {
		if strings.Contains(err.Error(), "Ollama") || strings.Contains(err.Error(), "connect") {
			t.Skipf("Skipping: Ollama required for test")
		}
		t.Fatalf("NewContextReshaper failed: %v", err)
	}
	defer reshaper.Close()

	// Test nil block
	_, err = reshaper.Reshape(context.Background(), nil)
	if err == nil {
		t.Error("should fail with nil block")
	}

	// Test empty content block
	emptyBlock := &ContextBlock{Title: "empty", Query: "", Results: []FormattedResult{}}
	_, err = reshaper.Reshape(context.Background(), emptyBlock)
	if err == nil {
		t.Error("should fail with empty content")
	}
}

// Integration test with real ONNX embeddings backend (LLM stays Ollama —
// ONNX LLM is blocked on hugot's CGO ORT backend, see ErrGenerativeRequiresORT).
func TestReshapeWithONNX(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	config := ReshapingConfig{
		EmbeddingsBackend: "onnx",
		EmbeddingsModel:   "all-MiniLM-L6-v2", // Downloads ~50MB on first run
		LLMBackend:        "ollama",           // Use Ollama instead of ONNX (ONNX not ready)
		LLMModel:          "mistral",
		MaxTokens:         256,
		TopK:              3,
	}

	reshaper, err := NewContextReshaper(config)
	if err != nil {
		// Skip if Ollama not available (expected in CI)
		if strings.Contains(err.Error(), "Ollama") || strings.Contains(err.Error(), "connect") {
			t.Skipf("Skipping: Ollama required for test")
		}
		t.Fatalf("NewContextReshaper with ONNX failed: %v", err)
	}
	defer reshaper.Close()

	testBlock := &ContextBlock{
		Title: "test - Python Kubernetes Docker",
		Query: "test query",
		Results: []FormattedResult{
			{
				Rank:      1,
				Relevance: 0.95,
				Source:    "test.md",
				Title:     "Programming",
				Content:   "Python is a high-level programming language. Kubernetes is an orchestration platform. Docker provides containerization technologies.",
				Tokens:    50,
			},
		},
		Metadata: make(map[string]string),
	}

	reshaped, err := reshaper.Reshape(context.Background(), testBlock)
	if err != nil {
		t.Fatalf("Reshape failed: %v", err)
	}

	// Verify output structure
	if reshaped.Original.Title != testBlock.Title {
		t.Errorf("Title mismatch: expected %s, got %s", testBlock.Title, reshaped.Original.Title)
	}
	if reshaped.Content == "" {
		t.Error("Reshaped content should not be empty")
	}
	if len(reshaped.KeyConcepts) == 0 {
		t.Error("Key concepts should be extracted")
	}
}

// Integration test with Ollama backend
// NOTE: Requires running Ollama locally (`ollama serve`)
// Skipped in CI unless Ollama is available
func TestReshapeWithOllama(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Ollama integration test")
	}

	// Try to create Ollama reshaper
	// If Ollama isn't running, this will fail gracefully
	config := ReshapingConfig{
		EmbeddingsBackend: "ollama",
		EmbeddingsModel:   "nomic-embed-text",
		LLMBackend:        "ollama",
		LLMModel:          "mistral:latest", // Requires: ollama pull mistral
		MaxTokens:         256,
		TopK:              3,
	}

	reshaper, err := NewContextReshaper(config)
	if err != nil {
		t.Skipf("Skipping: Ollama not available (%v)", err)
	}
	defer reshaper.Close()

	testBlock := &ContextBlock{
		Title: "Ollama Test",
		Query: "search query",
		Results: []FormattedResult{
			{
				Rank:      1,
				Relevance: 0.9,
				Source:    "test.md",
				Content:   "Semantic search finds relevant documents. Vector databases store embeddings. Language models generate text.",
				Tokens:    40,
			},
		},
		Metadata: make(map[string]string),
	}

	reshaped, err := reshaper.Reshape(context.Background(), testBlock)
	if err != nil {
		t.Fatalf("Reshape with Ollama failed: %v", err)
	}

	// Verify output
	if reshaped.Content == "" {
		t.Error("Reshaped content should not be empty")
	}
	if reshaped.Pipeline != "ollama-ollama" {
		t.Errorf("Pipeline should be 'ollama-ollama', got %s", reshaped.Pipeline)
	}
}

// Integration test with hybrid backend (ONNX embeddings + Ollama LLM)
func TestReshapeHybridONNXOllama(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping hybrid integration test")
	}

	config := ReshapingConfig{
		EmbeddingsBackend: "onnx",
		EmbeddingsModel:   "all-MiniLM-L6-v2",
		LLMBackend:        "ollama",
		LLMModel:          "mistral:latest",
		MaxTokens:         256,
		TopK:              3,
	}

	reshaper, err := NewContextReshaper(config)
	if err != nil {
		if strings.Contains(err.Error(), "Ollama") {
			t.Skipf("Skipping: Ollama not available")
		}
		t.Fatalf("NewContextReshaper failed: %v", err)
	}
	defer reshaper.Close()

	testBlock := &ContextBlock{
		Title: "Hybrid Test",
		Query: "test",
		Results: []FormattedResult{
			{
				Rank:      1,
				Relevance: 0.9,
				Source:    "hybrid.md",
				Content:   "Context reshaping improves clarity. Key concepts guide the LLM. Embeddings identify themes.",
				Tokens:    30,
			},
		},
		Metadata: make(map[string]string),
	}

	reshaped, err := reshaper.Reshape(context.Background(), testBlock)
	if err != nil {
		t.Fatalf("Reshape failed: %v", err)
	}

	if reshaped.Pipeline != "onnx-ollama" {
		t.Errorf("Pipeline should be 'onnx-ollama', got %s", reshaped.Pipeline)
	}
}

func TestIdentifyKeyConcepts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping concept identification test")
	}

	config := ReshapingConfig{
		EmbeddingsBackend: "onnx",
		EmbeddingsModel:   "all-MiniLM-L6-v2",
		LLMBackend:        "ollama", // ONNX LLM blocked on hugot's ORT backend
		LLMModel:          "mistral",
		TopK:              3,
	}

	reshaper, err := NewContextReshaper(config)
	if err != nil {
		// Skip if Ollama not available (expected in CI)
		if strings.Contains(err.Error(), "Ollama") || strings.Contains(err.Error(), "connect") {
			t.Skipf("Skipping: Ollama required for test")
		}
		t.Fatalf("NewContextReshaper failed: %v", err)
	}
	defer reshaper.Close()

	content := "Kubernetes is a container orchestration platform. Docker packages applications. Go is a systems programming language."

	concepts, err := reshaper.identifyKeyConcepts(context.Background(), content)
	if err != nil {
		t.Fatalf("identifyKeyConcepts failed: %v", err)
	}

	if len(concepts) == 0 {
		t.Error("should identify at least one concept")
	}
	if len(concepts) > config.TopK {
		t.Errorf("should not exceed TopK concepts, got %d", len(concepts))
	}
}

func TestReshapedBlockInheritsReferences(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping reference inheritance test")
	}

	config := ReshapingConfig{
		EmbeddingsBackend: "onnx",
		EmbeddingsModel:   "all-MiniLM-L6-v2",
		LLMBackend:        "ollama",
		LLMModel:          "mistral",
		MaxTokens:         256,
		TopK:              3,
	}

	reshaper, err := NewContextReshaper(config)
	if err != nil {
		if strings.Contains(err.Error(), "Ollama") || strings.Contains(err.Error(), "connect") {
			t.Skipf("Skipping: Ollama required for test")
		}
		t.Fatalf("NewContextReshaper failed: %v", err)
	}
	defer reshaper.Close()

	// Create a ContextBlock with references
	originalBlock := &ContextBlock{
		Title: "Test - References",
		Query: "test query",
		Results: []FormattedResult{
			{
				Rank:      1,
				Relevance: 0.95,
				Source:    "deployment.md",
				Title:     "Deployment",
				Content:   "Kubernetes orchestrates containers. Docker packages applications.",
				Tokens:    30,
			},
			{
				Rank:      2,
				Relevance: 0.82,
				Source:    "architecture.md",
				Title:     "Architecture",
				Content:   "Microservices are distributed services.",
				Tokens:    20,
			},
		},
		References: []DocumentReference{
			{
				Source:    "deployment.md",
				Relevance: 0.95,
			},
			{
				Source:    "architecture.md",
				Relevance: 0.82,
			},
		},
		Metadata: make(map[string]string),
	}

	// Reshape the block
	reshaped, err := reshaper.Reshape(context.Background(), originalBlock)
	if err != nil {
		t.Fatalf("Reshape failed: %v", err)
	}

	// Verify References are inherited in ReshapedBlock
	if len(reshaped.References) != len(originalBlock.References) {
		t.Errorf("ReshapedBlock References length = %d, want %d", len(reshaped.References), len(originalBlock.References))
	}

	// Verify reference sources are preserved
	sourceMap := make(map[string]float32)
	for _, ref := range reshaped.References {
		sourceMap[ref.Source] = ref.Relevance
	}

	if relevance, exists := sourceMap["deployment.md"]; !exists || relevance != 0.95 {
		t.Errorf("deployment.md reference should exist with relevance 0.95")
	}
	if relevance, exists := sourceMap["architecture.md"]; !exists || relevance != 0.82 {
		t.Errorf("architecture.md reference should exist with relevance 0.82")
	}
}
