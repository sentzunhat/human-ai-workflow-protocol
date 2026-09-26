// Package mcp implements a stdio MCP server that exposes hawp tools to
// AI agents (Claude Code, Cursor, Continue, etc.) over JSON-RPC 2.0.
package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
	ID      *json.RawMessage `json:"id"`
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
	return serve(os.Stdin, os.Stdout, repoRoot, version)
}

func serve(in io.Reader, out io.Writer, repoRoot, version string) error {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	enc := json.NewEncoder(out)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		if !json.Valid(line) {
			if err := enc.Encode(rpcResponse{
				JSONRPC: "2.0",
				Error:   &rpcError{Code: -32700, Message: "parse error"},
			}); err != nil {
				return err
			}
			continue
		}
		var req rpcRequest
		if err := json.Unmarshal(line, &req); err != nil {
			if err := enc.Encode(rpcResponse{
				JSONRPC: "2.0",
				Error:   &rpcError{Code: -32600, Message: "invalid request"},
			}); err != nil {
				return err
			}
			continue
		}
		if req.JSONRPC != "2.0" || req.Method == "" {
			if err := enc.Encode(rpcResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &rpcError{Code: -32600, Message: "invalid request"},
			}); err != nil {
				return err
			}
			continue
		}
		// Notifications carry no ID — no response.
		if req.ID == nil {
			continue
		}
		resp := dispatch(req, repoRoot, version)
		resp.JSONRPC = "2.0"
		resp.ID = req.ID
		if err := enc.Encode(resp); err != nil {
			return err
		}
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
