package distribution

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComputeExpectedOutputsProducesAllVariants(t *testing.T) {
	root := t.TempDir()
	writeDistributionFixture(t, root)

	outputs, err := ComputeExpectedOutputs(root)
	if err != nil {
		t.Fatal(err)
	}

	want := len(ActiveProviders) * 2 * 2
	if len(outputs) != want {
		t.Fatalf("outputs = %d, want %d", len(outputs), want)
	}

	for _, output := range outputs {
		if !strings.HasSuffix(output.Content, "\n") {
			t.Fatalf("trailing newline missing for %s", output.OutputPath)
		}
		if !strings.Contains(output.Content, "set -euo pipefail") {
			t.Fatalf("script body missing for %s", output.OutputPath)
		}
		if !strings.Contains(output.Content, "Local sync: run `hawp distribution sync`") {
			t.Fatalf("sync guidance missing for %s", output.OutputPath)
		}
	}
}

func TestFindDownstreamPathLeaksReportsLeakedPaths(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "core", ".hawp", "kit", "instructions", "da-file-tracking.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("bad core/.hawp/path\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	leaks, err := FindDownstreamPathLeaks(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(leaks) != 1 {
		t.Fatalf("leaks = %d, want 1", len(leaks))
	}
	if leaks[0].Line != 1 {
		t.Fatalf("line = %d, want 1", leaks[0].Line)
	}
}

func writeDistributionFixture(t *testing.T, root string) {
	t.Helper()

	files := map[string]string{
		"distribution/sources/shared/safety.md":              "# Safety\n",
		"distribution/sources/shared/repo-boundaries-kit.md": "# Boundaries\n",
		"distribution/sources/shared/install.md":             "# Install\n",
		"distribution/sources/shared/update.md":              "# Update\n",
		"distribution/sources/install/script-core.md":        "```bash\nset -euo pipefail\nREF=\"placeholder\"\nPROVIDER=\"placeholder\"\ncore\n```\n",
		"distribution/sources/install/script-footer.md":      "```bash\nfooter\n```\n",
		"distribution/sources/update/script-core.md":         "```bash\nset -euo pipefail\nREF=\"placeholder\"\nPROVIDER=\"placeholder\"\ncore\n```\n",
		"distribution/sources/update/script-footer.md":       "```bash\nfooter\n```\n",
	}

	for _, provider := range ActiveProviders {
		files["distribution/sources/providers/"+provider+"/preamble-install.md"] = "# pre install\n"
		files["distribution/sources/providers/"+provider+"/preamble-update.md"] = "# pre update\n"
		files["distribution/sources/providers/"+provider+"/safety.md"] = "# provider safety\n"
		files["distribution/sources/providers/"+provider+"/boundaries.md"] = "# provider boundaries\n"
		files["distribution/sources/providers/"+provider+"/install-contract.md"] = "# install contract\n"
		files["distribution/sources/providers/"+provider+"/update-contract.md"] = "# update contract\n"
		files["distribution/sources/providers/"+provider+"/install/main.md"] = "# install main\n"
		files["distribution/sources/providers/"+provider+"/install/dev.md"] = "# install dev\n"
		files["distribution/sources/providers/"+provider+"/update/main.md"] = "# update main\n"
		files["distribution/sources/providers/"+provider+"/update/dev.md"] = "# update dev\n"
		files["distribution/sources/providers/"+provider+"/script-install.md"] = "```bash\nprovider install\n```\n"
		files["distribution/sources/providers/"+provider+"/script-update.md"] = "```bash\nprovider update\n```\n"
	}

	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
