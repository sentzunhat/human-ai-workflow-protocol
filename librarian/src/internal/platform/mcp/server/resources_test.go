package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResourceProtocol(t *testing.T) {
	root := t.TempDir()
	for _, resource := range repoResources {
		path := filepath.Join(root, filepath.FromSlash(resource.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("contents for "+resource.URI), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	listed := dispatch(rpcRequest{Method: "resources/list"}, root, "test")
	resources, ok := listed.Result.(map[string]any)["resources"].([]map[string]any)
	if !ok || len(resources) != len(repoResources) {
		t.Fatalf("resources/list result = %#v, want %d resources", listed.Result, len(repoResources))
	}
	for i, resource := range repoResources {
		if resources[i]["uri"] != resource.URI || resources[i]["mimeType"] != "text/markdown" {
			t.Errorf("resource[%d] = %#v, want URI %q and text/markdown", i, resources[i], resource.URI)
		}
		params, _ := json.Marshal(map[string]string{"uri": resource.URI})
		read := dispatch(rpcRequest{Method: "resources/read", Params: params}, root, "test")
		if read.Error != nil {
			t.Fatalf("resources/read %q error = %v", resource.URI, read.Error)
		}
		contents := read.Result.(map[string]any)["contents"].([]map[string]any)
		if len(contents) != 1 || contents[0]["uri"] != resource.URI || contents[0]["text"] != "contents for "+resource.URI {
			t.Errorf("resources/read %q = %#v", resource.URI, contents)
		}
	}
}

func TestReadResourceRejectsInvalidAndMissingURIs(t *testing.T) {
	for _, test := range []struct {
		name   string
		params string
		code   int
	}{
		{name: "missing uri", params: `{}`, code: -32602},
		{name: "malformed params", params: `[]`, code: -32602},
		{name: "unknown uri", params: `{"uri":"hawp://repo/../../etc/passwd"}`, code: -32002},
	} {
		t.Run(test.name, func(t *testing.T) {
			resp := readResource(json.RawMessage(test.params), t.TempDir())
			if resp.Error == nil || resp.Error.Code != test.code {
				t.Fatalf("response = %+v, want error code %d", resp, test.code)
			}
		})
	}
}

func TestReadResourceRejectsUnsafeFiles(t *testing.T) {
	t.Run("symlink", func(t *testing.T) {
		root := t.TempDir()
		outside := filepath.Join(t.TempDir(), "outside.md")
		if err := os.WriteFile(outside, []byte("secret"), 0o600); err != nil {
			t.Fatal(err)
		}
		resource := repoResources[0]
		path := filepath.Join(root, filepath.FromSlash(resource.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, path); err != nil {
			t.Skipf("symlink unavailable: %v", err)
		}
		params, _ := json.Marshal(map[string]string{"uri": resource.URI})
		resp := readResource(params, root)
		if resp.Error == nil {
			t.Fatal("symlink resource was accepted")
		}
	})

	t.Run("hard link", func(t *testing.T) {
		root := t.TempDir()
		resource := repoResources[0]
		path := filepath.Join(root, filepath.FromSlash(resource.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		outside := filepath.Join(t.TempDir(), "outside.md")
		if err := os.WriteFile(outside, []byte("shared"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Link(outside, path); err != nil {
			t.Skipf("hard link unavailable: %v", err)
		}
		params, _ := json.Marshal(map[string]string{"uri": resource.URI})
		resp := readResource(params, root)
		if resp.Error == nil {
			t.Fatal("hard-linked resource was accepted")
		}
	})

	t.Run("oversized", func(t *testing.T) {
		root := t.TempDir()
		resource := repoResources[0]
		path := filepath.Join(root, filepath.FromSlash(resource.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(strings.Repeat("x", maxResourceBytes+1)), 0o600); err != nil {
			t.Fatal(err)
		}
		params, _ := json.Marshal(map[string]string{"uri": resource.URI})
		resp := readResource(params, root)
		if resp.Error == nil {
			t.Fatal("oversized resource was accepted")
		}
	})
}

func TestReadResourceReportsMissingFile(t *testing.T) {
	params, _ := json.Marshal(map[string]string{"uri": repoResources[0].URI})
	resp := readResource(params, t.TempDir())
	if resp.Error == nil || resp.Error.Code != -32002 {
		t.Fatalf("response = %+v, want resource-not-found error", resp)
	}
}
