package provision

import "testing"

func TestLoadManifestMissingFileReturnsEmpty(t *testing.T) {
	manifest, err := LoadManifest(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Assets) != 0 {
		t.Errorf("expected empty assets, got %+v", manifest.Assets)
	}
}

func TestManifestSaveAndLoadRoundTrip(t *testing.T) {
	root := t.TempDir()
	manifest, _ := LoadManifest(root)
	manifest.RuntimeVersion = "1.27.1"
	manifest.Assets["onnxruntime"] = AssetRecord{SHA256: "abc123", Size: 42, InstalledAt: "2026-07-20T00:00:00Z"}

	if err := manifest.Save(root); err != nil {
		t.Fatal(err)
	}

	reloaded, err := LoadManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.RuntimeVersion != "1.27.1" {
		t.Errorf("RuntimeVersion = %q", reloaded.RuntimeVersion)
	}
	record, ok := reloaded.Assets["onnxruntime"]
	if !ok || record.SHA256 != "abc123" || record.Size != 42 {
		t.Errorf("asset record = %+v", record)
	}
}

func TestManifestSatisfies(t *testing.T) {
	manifest := &Manifest{Assets: map[string]AssetRecord{
		"model_quantized.onnx": {SHA256: "deadbeef"},
	}}
	asset := Asset{Name: "model_quantized.onnx", SHA256: "deadbeef"}
	if !manifest.Satisfies(asset) {
		t.Error("expected Satisfies to match on name+hash")
	}
	if manifest.Satisfies(Asset{Name: "model_quantized.onnx", SHA256: "different"}) {
		t.Error("Satisfies must not match on hash mismatch (stale asset)")
	}
	if manifest.Satisfies(Asset{Name: "unknown", SHA256: "deadbeef"}) {
		t.Error("Satisfies must not match unknown asset name")
	}
}

func TestRuntimeAssetKnownPlatformsHaveConsistentFields(t *testing.T) {
	for key, asset := range runtimeAssets {
		if asset.URL == "" || asset.SHA256 == "" || asset.ArchiveMember == "" || asset.DestName == "" {
			t.Errorf("platform %s has an incomplete asset: %+v", key, asset)
		}
		if len(asset.SHA256) != 64 {
			t.Errorf("platform %s sha256 length = %d, want 64", key, len(asset.SHA256))
		}
	}
}

func TestModelAssetsHaveValidHashes(t *testing.T) {
	for _, asset := range ModelAssets {
		if len(asset.SHA256) != 64 {
			t.Errorf("asset %s sha256 length = %d, want 64", asset.Name, len(asset.SHA256))
		}
		if asset.Size <= 0 {
			t.Errorf("asset %s has non-positive size", asset.Name)
		}
	}
}
