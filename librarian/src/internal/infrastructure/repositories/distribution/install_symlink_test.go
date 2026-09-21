package distribution_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestInstallSurfacesRejectExecutableDestinationSymlinks(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..", "..", ".."))

	makefilePath := filepath.Join(repoRoot, "librarian", "src", "Makefile")
	makefile, err := os.ReadFile(makefilePath)
	if err != nil {
		t.Fatal(err)
	}
	makefileGuard := "[ -L ../../.hawp/bin/$(BINARY) ]"
	if got := strings.Count(string(makefile), makefileGuard); got != 2 {
		t.Fatalf("Makefile executable symlink guard count = %d, want 2", got)
	}

	for _, rel := range []string{
		"distribution/sources/install/script-core.md",
		"distribution/sources/update/script-core.md",
	} {
		path := filepath.Join(repoRoot, filepath.FromSlash(rel))
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		text := string(body)
		if !strings.Contains(text, `if [ -L ".hawp" ]`) {
			t.Fatalf("%s does not reject a symlinked .hawp root", rel)
		}
		if !strings.Contains(text, "find .hawp -type l -print -quit") {
			t.Fatalf("%s no longer scans .hawp for symlinks", rel)
		}
		guardCall := strings.Index(text, "\nreject_hawp_symlinks\n")
		binaryMove := strings.Index(text, `mv -f "$_tmpdir/binary" "$_dest"`)
		if guardCall == -1 {
			t.Fatalf("%s is missing reject_hawp_symlinks preflight", rel)
		}
		if binaryMove == -1 {
			t.Fatalf("%s is missing the binary install/update rename", rel)
		}
		if guardCall > binaryMove {
			t.Fatalf("%s must reject .hawp symlinks before installing the binary", rel)
		}
	}
}

func TestReconciliationParsesCanonicalAndLegacyBacklogs(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..", "..", ".."))

	tests := []struct {
		name       string
		backlog    string
		activePath string
		closedPath string
	}{
		{
			name: "canonical active row with escaped pipe and owner",
			backlog: "## Active Work\n\n" +
				"| UUID | Type | Title | Status | Owner | Plan File | Updated |\n" +
				"| ---- | ---- | ----- | ------ | ----- | --------- | ------- |\n" +
				"| `e3e2b3f0` | fix | Fix A \\| B | done | unassigned | [plan](closed/2026/09/23/e3e2b3f0/plan.md) | 2026-09-23 |\n",
			activePath: "e3e2b3f0/plan.md",
			closedPath: "2026/09/23/e3e2b3f0/plan.md",
		},
		{
			name: "canonical recently closed row",
			backlog: "## Recently Closed\n\n" +
				"| ID | Type | Title | Closed | Detail |\n" +
				"| -- | ---- | ----- | ------ | ------ |\n" +
				"| `7b58eb7d` | bug | Complete | 2026-09-21 | [plan](closed/2026/09/21/7b58eb7d/plan.md) |\n",
			activePath: "7b58eb7d/plan.md",
			closedPath: "2026/09/21/7b58eb7d/plan.md",
		},
		{
			name: "legacy done row",
			backlog: "## Done\n\n" +
				"| ID | Type | Title | Closed | Detail |\n" +
				"| -- | ---- | ----- | ------ | ------ |\n" +
				"| TASK-001 | bug | Complete | 2026-09-21 | [plan](closed/2026/09/21/TASK-001.md) |\n",
			activePath: "TASK-001.md",
			closedPath: "2026/09/21/TASK-001.md",
		},
	}

	for _, source := range []string{
		"distribution/sources/install/script-core.md",
		"distribution/sources/update/script-core.md",
	} {
		function := reconciliationFunction(t, filepath.Join(repoRoot, filepath.FromSlash(source)))
		for _, test := range tests {
			t.Run(filepath.Base(filepath.Dir(source))+"/"+test.name, func(t *testing.T) {
				root := t.TempDir()
				backlogPath := filepath.Join(root, ".hawp", "work", "BACKLOG.md")
				activePath := filepath.Join(root, ".hawp", "work", "active", filepath.FromSlash(test.activePath))
				if err := os.MkdirAll(filepath.Dir(activePath), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(backlogPath, []byte(test.backlog), 0o644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(activePath, []byte("# plan\n"), 0o644); err != nil {
					t.Fatal(err)
				}

				scriptPath := filepath.Join(root, "reconcile.sh")
				script := "set -euo pipefail\n" + function + "\nreconcile_closed_plans_from_backlog\n"
				if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
					t.Fatal(err)
				}
				command := exec.Command("bash", scriptPath)
				command.Dir = root
				output, err := command.CombinedOutput()
				if err != nil {
					t.Fatalf("reconciliation failed: %v\n%s", err, output)
				}

				closedPath := filepath.Join(root, ".hawp", "work", "closed", filepath.FromSlash(test.closedPath))
				if _, err := os.Stat(closedPath); err != nil {
					t.Fatalf("closed plan missing at %s: %v", closedPath, err)
				}
				if _, err := os.Stat(activePath); !os.IsNotExist(err) {
					t.Fatalf("active plan still exists after reconciliation: %v", err)
				}
			})
		}
	}
}

func reconciliationFunction(t *testing.T, sourcePath string) string {
	t.Helper()
	body, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	start := strings.Index(string(body), "reconcile_closed_plans_from_backlog() {")
	end := strings.Index(string(body), "\nMDATE=")
	if start == -1 || end == -1 || end <= start {
		t.Fatalf("could not extract reconciliation function from %s", sourcePath)
	}
	return string(body)[start:end]
}
