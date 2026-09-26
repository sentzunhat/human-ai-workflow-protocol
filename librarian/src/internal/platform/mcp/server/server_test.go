package mcp

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
)

func TestServeReturnsParseErrorForMalformedJSON(t *testing.T) {
	var out bytes.Buffer
	if err := serve(bytes.NewBufferString("{\n"), &out, t.TempDir(), "test"); err != nil {
		t.Fatalf("serve returned error: %v", err)
	}

	var resp rpcResponse
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.JSONRPC != "2.0" || resp.ID != nil || resp.Error == nil || resp.Error.Code != -32700 {
		t.Fatalf("response = %+v, want JSON-RPC parse error with null id", resp)
	}
}

func TestServeRejectsInvalidJSONRPCVersion(t *testing.T) {
	var out bytes.Buffer
	in := bytes.NewBufferString(`{"jsonrpc":"1.0","id":1,"method":"initialize"}` + "\n")
	if err := serve(in, &out, t.TempDir(), "test"); err != nil {
		t.Fatalf("serve returned error: %v", err)
	}

	var resp rpcResponse
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error == nil || resp.Error.Code != -32600 {
		t.Fatalf("response = %+v, want invalid request", resp)
	}
	if resp.ID == nil || string(*resp.ID) != "1" {
		t.Fatalf("response id = %v, want 1", resp.ID)
	}
}

func TestServeRejectsValidJSONThatIsNotARequest(t *testing.T) {
	tests := []string{
		`[]`,
		`{"jsonrpc":"2.0"}`,
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			var out bytes.Buffer
			if err := serve(bytes.NewBufferString(input+"\n"), &out, t.TempDir(), "test"); err != nil {
				t.Fatalf("serve returned error: %v", err)
			}

			var resp rpcResponse
			if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if resp.Error == nil || resp.Error.Code != -32600 || resp.ID != nil {
				t.Fatalf("response = %+v, want invalid request with null id", resp)
			}
		})
	}
}

func TestServeReturnsWriteError(t *testing.T) {
	err := serve(bytes.NewBufferString("{\n"), failingWriter{}, t.TempDir(), "test")
	if err == nil || !errors.Is(err, errWriteFailure) {
		t.Fatalf("serve error = %v, want write failure", err)
	}
}

var errWriteFailure = errors.New("write failure")

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errWriteFailure
}
