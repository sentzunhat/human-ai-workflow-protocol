package links_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	applinks "github.com/sentzunhat/hawp/librarian/src/internal/application/links"
)

func TestCleanRelinksWhenUniqueMatchFound(t *testing.T) {
	root := buildRepo(t, map[string]string{
		"README.md":                  "See the [setup guide](docs/setup-guide.md) for details.\n",
		"docs/guides/setup-guide.md": "# Setup\n",
	})

	result, err := applinks.Clean(root, true)
	if err != nil {
		t.Fatalf("Clean failed: %v", err)
	}
	if countChanges(result.Changes, "relinked") != 1 {
		t.Fatalf("changes = %#v, want 1 relinked change", result.Changes)
	}

	content, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "(docs/guides/setup-guide.md)") {
		t.Fatalf("README not relinked:\n%s", content)
	}
}

func TestCleanNeutralizesWhenNoUniqueMatchExists(t *testing.T) {
	root := buildRepo(t, map[string]string{
		"README.md": "See the [setup guide](docs/setup-guide.md) for details.\n",
	})

	result, err := applinks.Clean(root, true)
	if err != nil {
		t.Fatalf("Clean failed: %v", err)
	}
	if countChanges(result.Changes, "neutralized") != 1 {
		t.Fatalf("changes = %#v, want 1 neutralized change", result.Changes)
	}

	content, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "[setup guide]") {
		t.Fatalf("README should no longer contain markdown link syntax:\n%s", content)
	}
	if !strings.Contains(string(content), "setup guide") {
		t.Fatalf("README should keep visible text:\n%s", content)
	}
}

func countChanges(changes []applinks.CleanChange, action string) int {
	count := 0
	for _, change := range changes {
		if change.Action == action {
			count++
		}
	}
	return count
}
