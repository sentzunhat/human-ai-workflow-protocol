---
work-item: v010-3-2d
type: feature
title: "v0.1.0 Phase 3.2d: Anthropic Embeddings Stub (Future API)"
status: inbox
owner: unassigned
created: 2026-07-25
updated: 2026-07-25
---

# Phase 3.2d: Anthropic Embeddings (Stub)

## Mission

Scaffold Anthropic embeddings interface. Anthropic doesn't yet have a public embeddings API, so this is a placeholder ready for implementation when the API ships.

## Effort Estimate: 3 hours

---

## Implementation

### 1. Create `anthropic_embedder.go` (1 hour)

```go
type AnthropicEmbedder struct {
    client *anthropic.Client
    model  string // claude-embed (future)
}

func NewAnthropicEmbedder(apiKey, model string) (*AnthropicEmbedder, error) {
    // TODO: API not yet available (as of 2026-07)
    // When available, implement similar to OpenAI
    return nil, fmt.Errorf("Anthropic embeddings API not yet available")
}

func (e *AnthropicEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
    // TODO: Implement when Anthropic API ships
    return nil, fmt.Errorf("not yet implemented")
}

// Other interface methods...
```

### 2. Update Factory (15 min)

Add case in `NewEmbedder()`:
```go
case "anthropic":
    return NewAnthropicEmbedder(apiKey, model)
```

Will return error "not yet available" until API ships.

### 3. Write Tests (1.5 hours)

**Unit Tests:**
- `TestAnthropicEmbedder_NotYetAvailable` — verify error message
- `TestAnthropicEmbedder_InterfaceDefined` — interface conformance (compile check)

**Placeholder for Future:**
```go
// TODO: Integration tests once API available
// func TestAnthropicEmbed_RealAPI(t *testing.T) { ... }
```

### 4. Documentation (30 min)

**BACKENDS.md entry:**
```markdown
### Anthropic (Phase 3.2d)
**Status:** 🔮 Stub for future API

- **Models:** claude-embed (when available)
- **Availability:** Not yet released by Anthropic
- **Implementation Status:** Interface defined, ready to implement when API ships
- **ETA:** Monitor Anthropic announcements

**Use when:**
- Anthropic releases embeddings API
- You want to use Anthropic for both embeddings and LLM
```

---

## Acceptance Criteria

- [x] AnthropicEmbedder struct defined
- [x] All interface methods stubbed
- [x] Factory updated to route "anthropic" backend
- [x] Error message clear ("API not yet available")
- [x] Tests pass (compile check + error path)
- [x] Documentation added (with "future" status)
- [x] Ready to implement (code structure already in place)

---

## Notes

- This is low-priority work; mostly documentation + interface structure
- When Anthropic releases API, implementation is straightforward copy of OpenAI pattern
- No dependencies need to be added until API exists
- Can be done in parallel with higher-priority phases

---

## Related Work

- **Phase 3.2c** (OpenAI Embeddings) — parallel
- **Phase 3.3d** (Anthropic LLM) — similar stub pattern

---

## Files to Create/Modify

| File | Status |
|---|---|
| `librarian/go/internal/domain/embeddings/anthropic_embedder.go` | ✨ NEW |
| `librarian/go/internal/domain/embeddings/anthropic_embedder_test.go` | ✨ NEW |
| `librarian/go/internal/domain/embeddings/embedder.go` | 🔧 Update factory |
| `librarian/docs/BACKENDS.md` | 🔧 Add Anthropic (future) section |

---

**Priority:** LOW (nice to have, future-proofing)
**Effort:** Minimal (mostly stub + docs)
**Can start:** After v0.0.3 ships
**Parallel with:** All v0.1.0 phases
