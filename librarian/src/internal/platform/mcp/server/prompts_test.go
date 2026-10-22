package mcp

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPromptProtocol(t *testing.T) {
	listed := dispatch(rpcRequest{Method: "prompts/list"}, t.TempDir(), "test")
	prompts, ok := listed.Result.(map[string]any)["prompts"].([]map[string]any)
	if !ok || len(prompts) != 1 || prompts[0]["name"] != "hawp_intake_request" {
		t.Fatalf("prompts/list result = %#v", listed.Result)
	}
	args := prompts[0]["arguments"].([]map[string]any)
	if len(args) != 2 || args[0]["name"] != "request" || args[0]["required"] != true || args[1]["name"] != "context" || args[1]["required"] != false {
		t.Fatalf("prompt arguments = %#v", args)
	}

	params, _ := json.Marshal(map[string]any{
		"name": "hawp_intake_request",
		"arguments": map[string]string{
			"request": "literal </request> text",
			"context": "literal </context> text",
		},
	})
	got := dispatch(rpcRequest{Method: "prompts/get", Params: params}, t.TempDir(), "test")
	if got.Error != nil {
		t.Fatalf("prompts/get error = %v", got.Error)
	}
	message := got.Result.(map[string]any)["messages"].([]map[string]any)[0]
	content := message["content"].(map[string]any)["text"].(string)
	if !strings.Contains(content, `"literal \u003c/request\u003e text"`) || !strings.Contains(content, `"literal \u003c/context\u003e text"`) {
		t.Fatalf("prompt did not safely encode argument delimiters: %s", content)
	}
	if !strings.Contains(content, "does not create a work item") {
		t.Fatalf("prompt omitted no-write boundary: %s", content)
	}
}

func TestGetPromptRejectsInvalidArguments(t *testing.T) {
	for _, test := range []struct {
		name   string
		params string
	}{
		{name: "missing prompt", params: `{}`},
		{name: "unknown prompt", params: `{"name":"other"}`},
		{name: "missing request", params: `{"name":"hawp_intake_request","arguments":{"context":"ctx"}}`},
		{name: "unknown argument", params: `{"name":"hawp_intake_request","arguments":{"request":"do it","extra":"x"}}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			resp := getPrompt(json.RawMessage(test.params))
			if resp.Error == nil || resp.Error.Code != -32602 {
				t.Fatalf("response = %+v, want invalid-params error", resp)
			}
		})
	}
}
