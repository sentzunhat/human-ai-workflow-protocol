# Context Handling Audit — v0.0.3

**Date:** 2026-07-25  
**Scope:** v0.0.3 only (incremental patch versions, no v0.1.0 planning)  
**Purpose:** Verify what context types are being handled, identify gaps, document decisions

---

## Current Context Pipeline (v0.0.3)

### Input: ContextBlock
```go
type ContextBlock struct {
    Title       string              // e.g., "Search Results for 'vector embedding'"
    Query       string              // Original query
    ResultCount int                 // Number of results included
    TokenCount  int                 // Approximate token usage
    Results     []FormattedResult    // Ordered, deduplicated results
    Metadata    map[string]string    // Query timestamp, filters, etc
}

type FormattedResult struct {
    Rank      int     // Position (1, 2, 3, ...)
    Relevance float32 // Confidence (0.0 - 1.0)
    Source    string  // Document source/path (e.g., "README.md", "api/search.go")
    Title     string  // Chunk title/heading
    Content   string  // Actual text content
    Tokens    int     // Approximate tokens in this result
}
```

### Processing Pipeline (Reshaper)

1. **Identification:** Text only
   - Extracts sentences from Content
   - Embeds sentences via embeddings backend
   - Identifies key concepts (capitalized words)

2. **LLM Input:** Markdown text only
   ```
   Key concepts:
   - Embedding
   - Semantic Search
   - Vector
   
   Original content:
   [full markdown text]
   
   Reshaped content (maintain critical info, improve structure):
   ```

3. **Output:** ReshapedBlock
   ```go
   type ReshapedBlock struct {
       Original    ContextBlock    // Preserved for audit trail
       Content     string          // Reshaped markdown
       KeyConcepts []Concept       // Extracted concepts with embeddings
       Pipeline    string          // Backend combo used
   }
   ```

### Current Limitations ⚠️

| Context Type | Current | Handling | Example |
|---|---|---|---|
| **Text** | ✅ Supported | Fully processed | "Semantic search is a technique..." |
| **Code** | ⚠️ Partial | Text-only (no syntax/structure) | `func NewEmbedder() {}` |
| **Docs** | ✅ Supported | As text via semantic search | README.md chunks |
| **Images** | ❌ Not supported | No image handling | PNG/JPG from docs |
| **URLs** | ⚠️ Partial | Passed as source metadata only | `Source: "https://docs.../api.html"` |
| **Metadata** | ⚠️ Partial | Extracted to metadata map but not passed to LLM | timestamp, query filters |

---

## What v0.0.3 Currently DOES

### ✅ Text Processing (Full)
- Splits into sentences
- Embeds via semantic search
- Extracts concepts
- Passes to LLM for reshaping
- Preserves original content

**Example:**
```
Input:  "The semantic search pipeline uses embeddings to find relevant documents."
Output: "Semantic search pipeline uses vector embeddings to locate relevant documents efficiently."
```

### ✅ Deduplication & Formatting
- Removes duplicate results
- Groups by source/title
- Formats markdown with proper hierarchy
- Respects token budgets

### ✅ Concept Extraction
- Identifies capitalized terms (proper nouns)
- Ranks by relevance
- Embeds concepts alongside content
- Returns top-K (default 5)

---

## What v0.0.3 Does NOT Do

### ❌ Code File Interpretation
**Current:** Code is treated as plain text
```
func NewEmbedder(backend string) (Embedder, error) {
    // ...
}
```

**LLM receives:** Just the raw function text, no syntax awareness

**What's missing:**
- Language detection (Go vs Python vs Rust)
- Syntax highlighting context
- AST-aware importance weighting
- Function signature extraction
- Dependency analysis

**Impact:** LLM can still understand code, but less precisely

---

### ❌ File Artifacts / Binary Content
**Current:** Not supported
- PDFs: Not handled (would need text extraction)
- Images: Not handled (would need OCR/vision)
- Binary: Not handled (would need serialization)
- Archives: Not handled (would need extraction)

**Example gaps:**
- Screenshots (PNG/JPG) → can't pass to LLM
- Diagrams (SVG/PDF) → can't pass to LLM
- Notebooks (IPYNB) → JSON only, not parsed

---

### ❌ URL Fetching
**Current:** URLs are passed as metadata, not fetched
```
Source: "https://docs.example.com/api/search"  // Metadata only
```

**What's missing:**
- Automatic URL fetching
- HTML parsing
- Link following (no recursive crawl)
- Caching of fetched content

**Impact:** User must pre-fetch and paste content; URLs in results stay unresolved

---

### ❌ Rich Context Formats
**Current:** Text-only pipeline

**What's missing:**
- Markdown metadata preservation (frontmatter)
- Link references in markdown
- Code block language hints
- Structured JSON/YAML
- Table parsing

---

### ⚠️ Partial: Query Metadata
**Current:** Extracted but not used by LLM
```go
Metadata: map[string]string{
    "timestamp":  "2026-07-25T12:34:56Z",
    "filter":     "type:code",
    "source_db": "sqlite-fts5",
}
```

**What's missing:**
- Metadata isn't passed to LLM prompt
- LLM can't use query context to improve reshaping
- No temporal awareness

**Fix would be:** Include metadata in reshaping prompt

---

## Audit: Context Flow Diagram

```
Search Query
    ↓
[Lexical + Semantic Search] ← Finds documents
    ↓
ContextBlock {Results}
    ├─ Source path (Source)  ✅ Extracted
    ├─ Text content (Content) ✅ Processed
    ├─ Rank/Relevance ✅ Extracted
    ├─ Query metadata ⚠️ Extracted but unused
    ├─ Timestamps ⚠️ Extracted but unused
    ├─ Binary data ❌ Not supported
    └─ Images/PDFs ❌ Not supported
    ↓
[Deduplication & Formatting]
    ↓
[Context Reshaper Pipeline]
    ├─ Step 1: Identify concepts (embeddings) ✅
    ├─ Step 2: Build prompt ✅
    │   ├─ Content ✅ Passed
    │   ├─ Concepts ✅ Passed
    │   ├─ Metadata ❌ NOT passed
    │   └─ Source info ❌ NOT passed
    └─ Step 3: LLM reshape ✅
    ↓
ReshapedBlock {Content, KeyConcepts, Pipeline}
```

---

## Specific Answers to Your Questions

### Q1: "Are we getting files or artifacts for the LLM to take into consideration?"

**A:** Partially.

- **Text files:** ✅ Yes (via source search, content passed as text)
- **Code files:** ⚠️ Text only, no syntax awareness
- **Binary artifacts (PDF, images, notebooks):** ❌ No, not supported

**Example:**
```
Query: "show me the embedding reshaper"
→ Search finds: reshaper.go, reshaper_test.go
→ Content is plain text → LLM receives as raw text
→ LLM can understand it but has no syntax/structure hints
```

---

### Q2: "Can we pass code files to interpret or docs for interpretation?"

**A:** Yes for docs, partially for code.

**Docs:**
```
✅ README.md → semantic search → content extracted → LLM reshapes
✅ BACKENDS.md → chunks indexed → LLM can interpret
✅ Comments in code → extracted as text → interpreted as prose
```

**Code:**
```
⚠️ Go functions → indexed as plain text → LLM sees raw function
⚠️ Structure unknown (type definitions, method receivers, etc.)
⚠️ Dependencies not tracked
❌ No syntax tree analysis
```

**To improve:** Store language + highlight hints in metadata

---

### Q3: "If we have notes from semantic search, where to find them?"

**A:** In the ContextBlock.Results array.

**Structure:**
```go
FormattedResult{
    Source:    "docs/BACKENDS.md:line42",     // File + location
    Title:     "Embeddings Backends > ONNX",  // Heading hierarchy
    Content:   "Status: Production-ready...", // Actual text
    Relevance: 0.87,                          // Search score
}
```

**Finding them:**
- Via lexical search (FTS5): Exact phrase match
- Via semantic search: Vector similarity
- Deduped across sources
- Ranked by relevance
- Formatted with markdown hierarchy

**LLM access:** Currently only Content is passed; Source + Title could be added

---

### Q4: "What if it's a picture? Do we pass picture object or URL?"

**A:** Neither. v0.0.3 doesn't handle images.

**Options if we added image support:**

1. **URL to image** (NOT RECOMMENDED)
   - ❌ LLM can't fetch
   - ❌ Requires network
   - ❌ May break privacy

2. **Base64-encoded image data** (POSSIBLE)
   - ✅ LLM can interpret with vision
   - ⚠️ Token cost (images are large)
   - ⚠️ Requires vision-enabled LLM (Ollama doesn't support)

3. **Image description/OCR** (BEST)
   - ✅ Text-only pipeline
   - ✅ All LLMs can understand
   - ✅ Lower tokens
   - Requires: OCR tool or manual descriptions

4. **Skip images, extract text**
   - ✅ Simplest for v0.0.3
   - ⚠️ Loses visual structure
   - Works for docs with embedded diagrams

**Recommendation for v0.0.3:** Skip images, focus on text content

---

## Files Involved in Context Handling

### Core Reshaping
- `internal/application/context/reshaper.go` — Main pipeline
- `internal/application/context/reshaper_test.go` — Tests

### Context Structure
- `internal/application/context/format.go` — ContextBlock, FormattedResult
- `internal/application/context/dedup.go` — Deduplication
- `internal/application/context/encryption.go` — Encryption (metadata)

### Search Results
- `internal/domain/search/` — Semantic + lexical search
- Returns results that become ContextBlock.Results

### LLM Integration
- `internal/domain/llm/ollama_client.go` — LLM backend
- `internal/domain/embeddings/` — Embeddings backend

---

## Audit Results Summary

### What's Working ✅
- Text content extraction
- Semantic search on text
- Key concept identification
- Deduplication
- Markdown formatting
- Context token budgeting
- LLM reshaping of text
- Concept embeddings

### What's Partial ⚠️
- Code handling (text only, no syntax)
- URL handling (metadata only, not fetched)
- Query metadata (extracted but unused)

### What's Missing ❌
- Image support (PNG, JPG, SVG)
- PDF extraction
- Binary artifacts
- Notebook parsing
- Syntax-aware code processing
- URL auto-fetching
- Vision-based LLM calls

---

## Patch Version Improvements (v0.0.x)

### v0.0.3 (Current)
✅ Text-only, Ollama + ONNX backends

### v0.0.4 (Example patch)
Could add one of:
- **Metadata in prompts:** Include query context + source info in LLM prompt
- **Code language hints:** Store Go/Python/Rust in metadata, hint to LLM
- **URL fetching:** Optional inline URL resolution in semantic search

### v0.0.5 (Example patch)
Could add:
- **OCR support:** Extract text from images via tesseract
- **PDF extraction:** Use pdftotext on PDFs in search results

### v0.0.6+ (Example patches)
- **HTML parsing:** Fetch and parse URLs
- **Syntax highlighting:** Tree-sitter for code AST
- **Notebook parsing:** Read Jupyter JSON format

---

## References

**Library audit files:**
- `librarian/go/internal/application/context/` — Context pipeline
- `librarian/go/internal/domain/embeddings/` — Embeddings backends
- `librarian/go/internal/domain/llm/` — LLM backends
- `librarian/go/internal/domain/search/` — Search indexing

**Documentation:**
- `librarian/docs/BACKENDS.md` — Backend configuration
- `librarian/docs/CONTEXT_FLOW.md` — [if exists] Pipeline details
- `.hawp/kit/usage/` — HAWP workflow guidance

**Tests:**
- `librarian/go/internal/domain/integration_test.go` — Live backend tests
- `librarian/go/internal/domain/benchmark_test.go` — Performance data
- `librarian/go/internal/application/context/reshaper_test.go` — Reshaping

---

## Conclusion

**v0.0.3 is production-ready for:**
- ✅ Text document context reshaping
- ✅ Code file indexing (as text)
- ✅ Semantic search on all content
- ✅ Concept extraction
- ✅ Local LLM inference (Ollama)
- ✅ Offline embedding generation (ONNX)

**v0.0.3 does NOT support:**
- ❌ Images, PDFs, binary artifacts
- ❌ URL fetching
- ❌ Vision-based LLM reshaping
- ❌ Syntax-aware code processing

**For patch versions (v0.0.4+):** Any of the missing features can be added incrementally without major rewrites, following the same pluggable pipeline architecture.
