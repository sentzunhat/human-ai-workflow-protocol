package kit

import (
	"testing"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/kit/source"
)

func TestNormalizeFileName(t *testing.T) {
	cases := []struct{ in, want string }{
		{"README.md", ""}, {"start-here.md", ""}, {"Bad Name.MD", "bad-name.md"},
		{"Some_File.md", "some-file.md"}, {"UPPER.md", "upper.md"}, {"weird--name.md", "weird-name.md"},
		{"-leading.md", "leading.md"}, {"noext", ""}, {"Mixed Case Doc", "mixed-case-doc"},
	}
	for _, c := range cases {
		if got := NormalizeFileName(c.in); got != c.want {
			t.Errorf("NormalizeFileName(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestPlanNormalizationFromSnapshot(t *testing.T) {
	snapshot := source.Snapshot{
		Entries: []source.Entry{
			{Path: "/kit/usage/Bad Name.md", RelPath: "usage/Bad Name.md", Name: "Bad Name.md"},
			{Path: "/kit/usage/guide.md", RelPath: "usage/guide.md", Name: "guide.md"},
			{Path: "/kit/start-here.md", RelPath: "start-here.md", Name: "start-here.md"},
		},
		Files: []source.File{
			{Path: "/kit/usage/guide.md", RelPath: "usage/guide.md", Content: "see [bad](Bad%20Name.md)", Links: []source.Link{{Href: "Bad%20Name.md", Offset: 4}}},
			{Path: "/kit/start-here.md", RelPath: "start-here.md", Content: "see [bad too](usage/Bad Name.md)", Links: []source.Link{{Href: "usage/Bad Name.md", Offset: 4}}},
		},
	}
	renames := PlanFileRenames(snapshot)
	if len(renames) != 1 || renames[0].To != "/kit/usage/bad-name.md" {
		t.Fatalf("renames = %+v, want one rename", renames)
	}
	updates := PlanLinkUpdates(snapshot, map[string]string{renames[0].From: renames[0].To})
	if len(updates) != 1 || updates[0].To != "usage/bad-name.md" {
		t.Fatalf("updates = %+v, want start-here.md link update", updates)
	}
}
