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
