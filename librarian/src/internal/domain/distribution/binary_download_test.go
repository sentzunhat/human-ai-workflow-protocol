package distribution

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestBinaryDownloadPreservesInstalledFilesOnFailure(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("requires Bash")
	}
	for _, mode := range []string{"install", "update"} {
		for _, layout := range []string{"fresh", "legacy"} {
			for _, scenario := range []string{"success", "binary-fail", "checksums-fail", "missing-entry", "duplicate-entry", "malformed-entry", "mismatch", "no-hash", "hash-fail", "interrupt"} {
				t.Run(mode+"/"+layout+"/"+scenario, func(t *testing.T) {
					root := t.TempDir()
					binDir := filepath.Join(root, ".hawp/bin")
					expected := map[string]string{}
					if layout == "legacy" {
						if err := os.MkdirAll(binDir, 0755); err != nil {
							t.Fatal(err)
						}
						for _, name := range []string{"hawp", "hawp-bin", "hawp-mcp"} {
							expected[name] = "old " + name
							if err := os.WriteFile(filepath.Join(binDir, name), []byte(expected[name]), 0755); err != nil {
								t.Fatal(err)
							}
						}
					}
					payload := []byte("new platform binary\n")
					if err := os.WriteFile(filepath.Join(root, "payload"), payload, 0600); err != nil {
						t.Fatal(err)
					}
					checksum := fmt.Sprintf("%x", sha256.Sum256(payload))
					runs := 1
					if scenario == "success" {
						runs = 2
						expected["hawp"] = string(payload)
					}
					for run := 0; run < runs; run++ {
						runBinaryDownload(t, root, mode, scenario, checksum)
						entries, err := os.ReadDir(binDir)
						if err != nil || len(entries) != len(expected) {
							t.Fatalf("unexpected binary/staging files: %v %v; want %v", entries, err, expected)
						}
						for name, want := range expected {
							data, err := os.ReadFile(filepath.Join(binDir, name))
							if err != nil || string(data) != want {
								t.Fatalf("%s changed incorrectly: %q %v", name, data, err)
							}
						}
						if scenario == "success" {
							info, err := os.Stat(filepath.Join(binDir, "hawp"))
							if err != nil {
								t.Fatal(err)
							}
							if info.Mode().Perm()&0111 == 0 {
								t.Fatal("installed binary is not executable")
							}
						}
					}
				})
			}
		}
	}
}

func runBinaryDownload(t *testing.T, root, mode, scenario, checksum string) {
	t.Helper()
	script := binaryDownloadMocks + "\n" + binaryDownloadFunction(t, mode) + "\n" + mode + "_hawp_binary\n"
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-c", script)
	cmd.Dir = root
	for _, value := range os.Environ() {
		if !strings.HasPrefix(value, "BASH_ENV=") && !strings.HasPrefix(value, "ENV=") {
			cmd.Env = append(cmd.Env, value)
		}
	}
	cmd.Env = append(cmd.Env, "HAWP_TEST_MODE="+scenario, "HAWP_TEST_SHA="+checksum)
	output, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		t.Fatal("download test timed out")
	}
	if (err == nil) != (scenario == "success") {
		t.Fatalf("unexpected outcome: %v\n%s", err, output)
	}
}

func binaryDownloadFunction(t *testing.T, mode string) string {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(file), "../../../../../distribution/sources", mode, "script-core.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	start := strings.Index(string(data), mode+"_hawp_binary() (")
	if start < 0 {
		t.Fatal("missing binary function")
	}
	end := strings.Index(string(data[start:]), "\n)\n"+mode+"_hawp_binary")
	if end < 0 {
		t.Fatal("missing binary function end")
	}
	return string(data[start : start+end+2])
}
