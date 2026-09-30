package mcp

import "encoding/json"

// MCP content block types.

type toolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type toolResult struct {
	Content []toolContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// text wraps a plain string as a successful MCP text response.
func text(s string) rpcResponse {
	return rpcResponse{Result: toolResult{Content: []toolContent{{Type: "text", Text: s}}}}
}

// jsonResult serialises v as a JSON text response.
func jsonResult(v any) rpcResponse {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return toolErr("json marshal: " + err.Error())
	}
	return text(string(b))
}

// toolErr wraps a plain string as a failed MCP text response.
func toolErr(s string) rpcResponse {
	return rpcResponse{Result: toolResult{
		Content: []toolContent{{Type: "text", Text: s}},
		IsError: true,
	}}
}

// Search response schema.

// SearchResponse is the structured JSON returned by hawp_search.
type SearchResponse struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
}

// SearchResult is one ranked match.
type SearchResult struct {
	Source    string      `json:"source"`
	Relevance float32     `json:"relevance"`
	Content   string      `json:"content"`
	Lines     LineInfo    `json:"lines"`
	Context   ContextInfo `json:"context"`
}

// LineInfo holds the precise position of a chunk in its source file.
type LineInfo struct {
	Range  LineRange `json:"range"`
	Source int       `json:"source"` // line where the query term first matches
}

// LineRange is an inclusive start/end line pair (1-indexed).
type LineRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// ContextInfo carries the suggested read window around the match.
type ContextInfo struct {
	Window LineRange `json:"window"`
}

// ContextSearchResponse is returned by hawp_search when context:true.
// The content field is a pre-shaped markdown block ready for LLM injection.
type ContextSearchResponse struct {
	Query         string `json:"query"`
	Content       string `json:"content"`
	TokenCount    int    `json:"token_count"`
	Budget        int    `json:"budget"`
	ChunksUsed    int    `json:"chunks_used"`
	ChunksDropped int    `json:"chunks_dropped"`
}

// WorkIntakeResponse is returned by hawp_work_intake. It keeps the shaped
// draft separate from retrieval metadata so callers can distinguish a ready
// draft from missing context.
type WorkIntakeResponse struct {
	State           string              `json:"state"`
	Draft           *WorkIntakeDraft    `json:"draft,omitempty"`
	Retrieval       WorkIntakeRetrieval `json:"retrieval"`
	Questions       []string            `json:"questions,omitempty"`
	Warnings        []string            `json:"warnings,omitempty"`
	TokenAccounting TokenAccounting     `json:"token_accounting"`
}

type WorkIntakeDraft struct {
	Input       string `json:"input"`
	Context     string `json:"context"`
	Mission     string `json:"mission"`
	Constraints string `json:"constraints"`
	Output      string `json:"output_spec"`
	Checkpoint  string `json:"done_signal,omitempty"`
}

type WorkIntakeRetrieval struct {
	Query         string `json:"query"`
	ChunksUsed    int    `json:"chunks_used"`
	ChunksDropped int    `json:"chunks_dropped"`
	Budget        int    `json:"budget"`
}

type TokenAccounting struct {
	ContextTokens int `json:"context_tokens"`
	ShapedTokens  int `json:"shaped_tokens"`
	SavingsPct    int `json:"savings_pct"`
}

// WorkReshapeResponse is the structured response from hawp_work_reshape.
// Fields match the WorkIntakeDraft contract so callers can deserialize them.
type WorkReshapeResponse struct {
	Mission     string `json:"mission"`
	Constraints string `json:"constraints"`
	Output      string `json:"output_spec"`
	Checkpoint  string `json:"done_signal,omitempty"`
}
