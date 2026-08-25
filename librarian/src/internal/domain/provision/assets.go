// Package provision defines the assets provisioned into ~/.hawp/ (ONNX
// Runtime shared library, embedding model files) and the manifest that
// tracks installed versions.
package provision

import (
	"fmt"
	"runtime"
)

// Asset is one downloadable file: its source, integrity hash, and where it
// lands relative to its category directory.
type Asset struct {
	Name          string // manifest key, e.g. "onnxruntime", "model_quantized.onnx"
	URL           string
	SHA256        string
	Size          int64
	ArchiveMember string // for archives: the internal path to extract; "" for plain files
	DestName      string // filename to write under the category directory
}

// RuntimeVersion is the pinned ONNX Runtime release these assets come from.
const RuntimeVersion = "1.27.1"

// ModelName identifies the pinned embedding model.
const ModelName = "Xenova/all-MiniLM-L6-v2"

// ModelVersion pins the model asset set (bump when swapping model files).
const ModelVersion = "1"

// platformKey identifies the current OS/arch in ONNX Runtime's release
// naming scheme. Returns an error for combinations with no prebuilt asset
// (e.g. darwin/amd64 — upstream ships arm64-only macOS builds as of v1.27.1).
func platformKey() (string, error) {
	switch runtime.GOOS + "/" + runtime.GOARCH {
	case "linux/amd64":
		return "linux-x64", nil
	case "linux/arm64":
		return "linux-aarch64", nil
	case "darwin/arm64":
		return "osx-arm64", nil
	case "windows/amd64":
		return "win-x64", nil
	case "windows/arm64":
		return "win-arm64", nil
	default:
		return "", fmt.Errorf("no ONNX Runtime %s prebuilt asset for %s/%s", RuntimeVersion, runtime.GOOS, runtime.GOARCH)
	}
}

// runtimeAssets maps each supported platform key to its ONNX Runtime
// archive asset. Verified 2026-07-20 against the v1.27.1 GitHub release.
var runtimeAssets = map[string]Asset{
	"linux-x64": {
		Name:          "onnxruntime",
		URL:           "https://github.com/microsoft/onnxruntime/releases/download/v1.27.1/onnxruntime-linux-x64-1.27.1.tgz",
		SHA256:        "25b1ef1fea1acd210d63f8f24dc870ad6e077795ce1f54876252c6d3803c15af",
		Size:          8828892,
		ArchiveMember: "onnxruntime-linux-x64-1.27.1/lib/libonnxruntime.so.1.27.1",
		DestName:      "libonnxruntime.so",
	},
	"linux-aarch64": {
		Name:          "onnxruntime",
		URL:           "https://github.com/microsoft/onnxruntime/releases/download/v1.27.1/onnxruntime-linux-aarch64-1.27.1.tgz",
		SHA256:        "33c67e33d1e25b816878366ea276589a024f71f000e7ff955c4b33224d639edd",
		Size:          7812402,
		ArchiveMember: "onnxruntime-linux-aarch64-1.27.1/lib/libonnxruntime.so.1.27.1",
		DestName:      "libonnxruntime.so",
	},
	"osx-arm64": {
		Name:          "onnxruntime",
		URL:           "https://github.com/microsoft/onnxruntime/releases/download/v1.27.1/onnxruntime-osx-arm64-1.27.1.tgz",
		SHA256:        "e42b77a7281cc6e55141bf44fcfbac2c782b823a491bbb6ac33c781dd991f8a6",
		Size:          31959937,
		ArchiveMember: "onnxruntime-osx-arm64-1.27.1/lib/libonnxruntime.1.27.1.dylib",
		DestName:      "libonnxruntime.dylib",
	},
	"win-x64": {
		Name:          "onnxruntime",
		URL:           "https://github.com/microsoft/onnxruntime/releases/download/v1.27.1/onnxruntime-win-x64-1.27.1.zip",
		SHA256:        "2e00414a63fdef0914cd5a5ede6c707844878e0c08e1b6693842f0451b2df2a1",
		Size:          77242362,
		ArchiveMember: "onnxruntime-win-x64-1.27.1/lib/onnxruntime.dll",
		DestName:      "onnxruntime.dll",
	},
	"win-arm64": {
		Name:          "onnxruntime",
		URL:           "https://github.com/microsoft/onnxruntime/releases/download/v1.27.1/onnxruntime-win-arm64-1.27.1.zip",
		SHA256:        "6e22c2061ba6400b42a59663d700c8694e4e8fe654cf452c4700c24237407ae1",
		Size:          78590093,
		ArchiveMember: "onnxruntime-win-arm64-1.27.1/lib/onnxruntime.dll",
		DestName:      "onnxruntime.dll",
	},
}

// RuntimeAsset returns the ONNX Runtime asset for the current platform.
func RuntimeAsset() (Asset, error) {
	key, err := platformKey()
	if err != nil {
		return Asset{}, err
	}
	asset, ok := runtimeAssets[key]
	if !ok {
		return Asset{}, fmt.Errorf("no ONNX Runtime asset registered for platform key %q", key)
	}
	return asset, nil
}

// ModelAssets are the embedding model files provisioned by hawp init.
// Only models with verified SHA-256 checksums are included here.
// all-MiniLM-L6-v2 files go to ~/.hawp/models/minilm/.
// Verified 2026-07-22 against HuggingFace repos.
//
// BGE-base-en-v1.5 (768-dim) is tracked separately in BGEModelAssets.
// Its checksums were never verified; it is excluded from the default
// provision set to prevent guaranteed download failures blocking init.
var ModelAssets = []Asset{
	// all-MiniLM-L6-v2: 384-dim, verified checksums
	{
		Name:     "minilm_model.onnx",
		URL:      "https://huggingface.co/Xenova/all-MiniLM-L6-v2/resolve/main/onnx/model_quantized.onnx",
		SHA256:   "afdb6f1a0e45b715d0bb9b11772f032c399babd23bfc31fed1c170afc848bdb1",
		Size:     22972370,
		DestName: "minilm/model_quantized.onnx",
	},
	{
		Name:     "minilm_tokenizer.json",
		URL:      "https://huggingface.co/Xenova/all-MiniLM-L6-v2/resolve/main/tokenizer.json",
		SHA256:   "da0e79933b9ed51798a3ae27893d3c5fa4a201126cef75586296df9b4d2c62a0",
		Size:     711661,
		DestName: "minilm/tokenizer.json",
	},
	{
		Name:     "minilm_config.json",
		URL:      "https://huggingface.co/Xenova/all-MiniLM-L6-v2/resolve/main/config.json",
		SHA256:   "7135149f7cffa1a573466c6e4d8423ed73b62fd2332c575bf738a0d033f70df7",
		Size:     650,
		DestName: "minilm/config.json",
	},
}

// BGEModelAssets describes the BGE-base-en-v1.5 (768-dim) model files.
// These are NOT included in the default provision set because their SHA-256
// checksums have never been verified against the upstream HuggingFace repo.
// Add them to a custom Registry once real checksums are obtained.
var BGEModelAssets = []Asset{
	{
		Name:     "bge_model.onnx",
		URL:      "https://huggingface.co/Xenova/bge-base-en-v1.5/resolve/main/onnx/model.onnx",
		SHA256:   "", // TODO: verify from HuggingFace before use
		Size:     130000000,
		DestName: "bge/model.onnx",
	},
	{
		Name:     "bge_tokenizer.json",
		URL:      "https://huggingface.co/Xenova/bge-base-en-v1.5/resolve/main/tokenizer.json",
		SHA256:   "", // TODO: verify from HuggingFace before use
		Size:     800000,
		DestName: "bge/tokenizer.json",
	},
	{
		Name:     "bge_config.json",
		URL:      "https://huggingface.co/Xenova/bge-base-en-v1.5/resolve/main/config.json",
		SHA256:   "", // TODO: verify from HuggingFace before use
		Size:     700,
		DestName: "bge/config.json",
	},
}
