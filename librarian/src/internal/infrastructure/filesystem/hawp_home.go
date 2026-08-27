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
	Config           string // Root/config — global default config files
	Runtime          string // Root/runtime — ONNX Runtime shared library
	Cache            string // Root/cache — transient data (can be deleted safely)
	Downloads        string // Root/cache/downloads — staging area for verified downloads
	UpdateCacheFile  string // Root/cache/update-check.json — last release check result
	UpdateConfigFile string // Root/config/update.json — update preferences
	UsageConfigFile  string // Root/config/usage.json — usage logging preferences
	UsageDB          string // Root/usage.db — call log (separate from search index)
}

// ResolveHawpHome builds the ~/.hawp/ layout under the given home directory.
func ResolveHawpHome(home string) HawpHome {
	root := filepath.Join(home, ".hawp")
	modelsRoot := filepath.Join(root, "models")
	cacheRoot := filepath.Join(root, "cache")
	configRoot := filepath.Join(root, "config")
	return HawpHome{
		Root:             root,
		Index:            filepath.Join(root, "index"),
		Models:           modelsRoot,
		ModelsEmbedding:  filepath.Join(modelsRoot, "embedding"),
		ModelsLLM:        filepath.Join(modelsRoot, "llm"),
		Config:           configRoot,
		Runtime:          filepath.Join(root, "runtime"),
		Cache:            cacheRoot,
		Downloads:        filepath.Join(cacheRoot, "downloads"),
		UpdateCacheFile:  filepath.Join(cacheRoot, "update-check.json"),
		UpdateConfigFile: filepath.Join(configRoot, "update.json"),
		UsageConfigFile:  filepath.Join(configRoot, "usage.json"),
		UsageDB:          filepath.Join(root, "usage.db"),
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
		h.Cache,
		h.Downloads,
	}
}
