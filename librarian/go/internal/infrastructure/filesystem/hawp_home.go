package filesystem

import "path/filepath"

// HawpHome is the resolved ~/.hawp/ layout used by init-time provisioning.
// One location shared across repos on the same machine.
type HawpHome struct {
	Root             string
	Index            string // Root/index — the librarian.db (fbf12a93)
	Models           string // Root/models — embedding and LLM model files
	ModelsEmbedding  string // Root/models/embedding — embedding models (BGE, MiniLM, etc.)
	ModelsLLM        string // Root/models/llm — LLM models (TinyLlama, Mistral, etc.)
	Config           string // Root/config — global default config
	Runtime          string // Root/runtime — ONNX Runtime shared library
	Downloads        string // Root/cache/downloads — staging area for verified downloads
}

// ResolveHawpHome builds the ~/.hawp/ layout under the given home directory.
func ResolveHawpHome(home string) HawpHome {
	root := filepath.Join(home, ".hawp")
	modelsRoot := filepath.Join(root, "models")
	return HawpHome{
		Root:             root,
		Index:            filepath.Join(root, "index"),
		Models:           modelsRoot,
		ModelsEmbedding:  filepath.Join(modelsRoot, "embedding"),
		ModelsLLM:        filepath.Join(modelsRoot, "llm"),
		Config:           filepath.Join(root, "config"),
		Runtime:          filepath.Join(root, "runtime"),
		Downloads:        filepath.Join(root, "cache", "downloads"),
	}
}

// Dirs returns every directory that must exist.
func (h HawpHome) Dirs() []string {
	return []string{
		h.Root,
		h.Index,
		h.Models,
		h.ModelsEmbedding,
		h.ModelsLLM,
		h.Config,
		h.Runtime,
		h.Downloads,
	}
}
