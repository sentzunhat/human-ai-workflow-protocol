// Package mcp implements a stdio MCP server that exposes hawp tools to
// AI agents (Claude Code, Cursor, Continue, etc.) over JSON-RPC 2.0.
package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type rpcRequest struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Result  any              `json:"result,omitempty"`
	Error   *rpcError        `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Serve reads JSON-RPC 2.0 messages from stdin and writes responses to
// stdout until stdin closes. repoRoot is the HAWP project directory used
// by all tool handlers.
func Serve(repoRoot, version string) error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	enc := json.NewEncoder(os.Stdout)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var req rpcRequest
		if err := json.Unmarshal(line, &req); err != nil {
			continue
		}
		// Notifications carry no ID — no response.
		if req.ID == nil {
			continue
		}
		resp := dispatch(req, repoRoot, version)
		resp.JSONRPC = "2.0"
		resp.ID = req.ID
		_ = enc.Encode(resp)
	}
	return scanner.Err()
}

func dispatch(req rpcRequest, repoRoot, version string) rpcResponse {
	switch req.Method {
	case "initialize":
		return rpcResponse{Result: map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo":      map[string]any{"name": "hawp", "version": version},
		}}
	case "tools/list":
		return rpcResponse{Result: map[string]any{"tools": toolDefs()}}
	case "tools/call":
		return callTool(req.Params, repoRoot)
	default:
		return rpcResponse{Error: &rpcError{
			Code:    -32601,
			Message: fmt.Sprintf("method not found: %s", req.Method),
		}}
	}
}

func errResp(code int, msg string) rpcResponse {
	return rpcResponse{Error: &rpcError{Code: code, Message: msg}}
}
