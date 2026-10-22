package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

const maxResourceBytes = 1 << 20

type repoResource struct {
	URI         string
	Name        string
	Description string
	Path        string
}

var repoResources = []repoResource{
	{
		URI:         "hawp://repo/kit/start-here",
		Name:        "HAWP operating guide",
		Description: "The selected repository's HAWP operating guide.",
		Path:        ".hawp/kit/start-here.md",
	},
	{
		URI:         "hawp://repo/work/backlog",
		Name:        "HAWP work backlog",
		Description: "The selected repository's HAWP work backlog.",
		Path:        ".hawp/work/BACKLOG.md",
	},
	{
		URI:         "hawp://repo/kit/usage/mcp",
		Name:        "HAWP MCP guide",
		Description: "The selected repository's MCP usage and workflow guide.",
		Path:        ".hawp/kit/usage/mcp/README.md",
	},
}

func listResources() rpcResponse {
	resources := make([]map[string]any, 0, len(repoResources))
	for _, resource := range repoResources {
		resources = append(resources, map[string]any{
			"uri":         resource.URI,
			"name":        resource.Name,
			"description": resource.Description,
			"mimeType":    "text/markdown",
		})
	}
	return rpcResponse{Result: map[string]any{"resources": resources}}
}

func readResource(params json.RawMessage, repoRoot string) rpcResponse {
	var request struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(params, &request); err != nil || request.URI == "" {
		return errResp(-32602, "invalid params: uri is required")
	}

	var selected *repoResource
	for i := range repoResources {
		if repoResources[i].URI == request.URI {
			selected = &repoResources[i]
			break
		}
	}
	if selected == nil {
		return errResp(-32002, "resource not found")
	}

	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return errResp(-32603, "cannot resolve repository root")
	}
	target := filepath.Join(root, filepath.FromSlash(selected.Path))
	if err := filesystem.RejectSymlinkAncestors(root, target); err != nil {
		return errResp(-32603, "resource path is unsafe")
	}
	if err := filesystem.RejectHardLinkedFile(target); err != nil {
		return errResp(-32603, "resource file is unsafe")
	}
	info, err := os.Lstat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return errResp(-32002, "resource not found")
		}
		return errResp(-32603, "cannot inspect resource")
	}
	if !info.Mode().IsRegular() || info.Size() > maxResourceBytes {
		return errResp(-32603, "resource is not a bounded regular file")
	}
	content, err := os.ReadFile(target)
	if err != nil {
		return errResp(-32603, "cannot read resource")
	}
	if len(content) > maxResourceBytes {
		return errResp(-32603, "resource exceeds size limit")
	}

	return rpcResponse{Result: map[string]any{
		"contents": []map[string]any{{
			"uri":      selected.URI,
			"mimeType": "text/markdown",
			"text":     string(content),
		}},
	}}
}
