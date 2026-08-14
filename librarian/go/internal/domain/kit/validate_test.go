package kit

import (
	"testing"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/kit/source"
)

func completeSnapshot() source.Snapshot {
	files := make([]source.File, 0, len(RequiredFiles))
	entries := make([]source.Entry, 0, len(RequiredFiles))
	for _, rel := range RequiredFiles {
		files = append(files, source.File{Path: "/kit/" + rel, RelPath: rel, Content: "# " + rel})
		entries = append(entries, source.Entry{Path: "/kit/" + rel, RelPath: rel, Name: relName(rel)})
	}
	return source.Snapshot{Entries: entries, Files: files}
}

func relName(rel string) string {
	for i := len(rel) - 1; i >= 0; i-- {
		if rel[i] == '/' {
			return rel[i+1:]
		}
	}
	return rel
}

func TestValidateCleanKit(t *testing.T) {
	issues, checks := Validate(completeSnapshot())
	if checks != 3 {
		t.Errorf("checks = %d, want 3", checks)
	}
	if len(issues) != 0 {
		t.Fatalf("issues on clean kit: %+v", issues)
	}
}

func TestFileNaming(t *testing.T) {
	snapshot := completeSnapshot()
	snapshot.Entries = append(snapshot.Entries,
		source.Entry{Path: "/kit/Bad Name.md", RelPath: "Bad Name.md", Name: "Bad Name.md"},
		source.Entry{Path: "/kit/README.md", RelPath: "README.md", Name: "README.md"},
	)
	issues := CheckFileNaming(snapshot)
	if len(issues) != 1 || issues[0].File != "Bad Name.md" {
		t.Fatalf("naming issues = %+v, want exactly the bad name", issues)
	}
}

func TestRequiredFiles(t *testing.T) {
	snapshot := completeSnapshot()
	snapshot.Files = snapshot.Files[1:]
	issues := CheckRequiredFiles(snapshot)
	if len(issues) != 1 || issues[0].File != "start-here.md" {
		t.Fatalf("required issues = %+v, want missing start-here.md", issues)
	}
}

func TestInternalLinks(t *testing.T) {
	snapshot := completeSnapshot()
	snapshot.Files = append(snapshot.Files, source.File{
		Path: "/kit/usage/guide.md", RelPath: "usage/guide.md",
		Content: "[ok](../start-here.md)\n[broken](missing.md)\n[ext](https://x.test)\n```\n[fenced](nope.md)\n```\n",
		Links:   []source.Link{{Href: "../start-here.md"}, {Href: "missing.md"}, {Href: "https://x.test"}},
	})
	issues := CheckInternalLinks(snapshot)
	if len(issues) != 1 || issues[0].File != "usage/guide.md" || issues[0].Message != "broken link: missing.md" {
		t.Fatalf("link issues = %+v, want only the broken link", issues)
	}
}
