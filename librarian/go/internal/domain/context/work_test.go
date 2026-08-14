package context

import (
	"testing"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/context/source"
	domainwork "github.com/sentzunhat/hawp/librarian/go/internal/domain/work"
)

const backlogFixture = `# Backlog

## Active Work

| ID | Type | Title | Status | Plan File | Updated |
| --- | --- | --- | --- | --- | --- |
| TASK-001 | feature | thing in progress | in-progress | [plan](active/TASK-001.md) | 2026-07-21 |

## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |
| TASK-002 | improvement | parked thing | not needed | [plan](parked/TASK-002.md) | 2026-07-10 |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
| TASK-003 | bug | fixed thing | 2026-07-05 | [plan](closed/2026/07/05/TASK-003.md) |
`

func TestEnrichWorkResolvesBacklogMetadata(t *testing.T) {
	backlog := &domainwork.Backlog{
		Active: []domainwork.BacklogRow{{ID: "TASK-001", Type: "feature", Status: "in-progress"}},
		Parked: []domainwork.BacklogRow{{ID: "TASK-002", Type: "improvement", Status: "not needed"}},
		Closed: []domainwork.BacklogRow{{ID: "TASK-003", Type: "bug", Status: "2026-07-05"}},
	}
	corpus := source.WorkCorpus{
		Backlog: backlog,
		Files: []source.File{
			{RelPath: "active/TASK-001.md", RepoPath: ".hawp/work/active/TASK-001.md", Content: "# active plan\n"},
			{RelPath: "parked/TASK-002.md", RepoPath: ".hawp/work/parked/TASK-002.md", Content: "# parked plan\n"},
			{RelPath: "closed/2026/07/05/TASK-003.md", RepoPath: ".hawp/work/closed/2026/07/05/TASK-003.md", Content: "# closed plan\n"},
			{RelPath: "decisions/2026/07/05/adr-001.md", RepoPath: ".hawp/work/decisions/2026/07/05/adr-001.md", Content: "# a decision with no backlog row\n"},
			{RelPath: "future/TASK-004.md", RepoPath: ".hawp/work/future/TASK-004.md", Content: "# ignored\n"},
		},
	}

	docs := EnrichWork(corpus)
	if len(docs) != 4 {
		t.Fatalf("documents = %d, want 4", len(docs))
	}

	byID := map[string]Document{}
	var decision Document
	for _, doc := range docs {
		if doc.Role == "decisions" {
			decision = doc
			continue
		}
		byID[doc.ID] = doc
	}

	active := byID["TASK-001"]
	if active.Role != "active" || active.Type != "feature" || active.Status != "in-progress" {
		t.Errorf("active doc = %+v", active)
	}
	if active.ContextPrefix != "[work/active] TASK-001 (feature, in-progress)" {
		t.Errorf("active ContextPrefix = %q", active.ContextPrefix)
	}

	parked := byID["TASK-002"]
	if parked.Role != "parked" || parked.Type != "improvement" {
		t.Errorf("parked doc = %+v", parked)
	}

	closed := byID["TASK-003"]
	if closed.Role != "closed" || closed.Type != "bug" || closed.ClosedDate != "2026-07-05" {
		t.Errorf("closed doc = %+v", closed)
	}
	if closed.ContextPrefix != "[work/closed] TASK-003 (bug, closed 2026-07-05)" {
		t.Errorf("closed ContextPrefix = %q", closed.ContextPrefix)
	}

	// Documents with no matching backlog row still get a sensible
	// fallback (Role, but no Type/Status) — no error.
	if decision.Role != "decisions" {
		t.Fatalf("decision doc missing/misclassified: %+v", decision)
	}
	if decision.Type != "" || decision.Status != "" {
		t.Errorf("decision doc should have no resolved metadata: %+v", decision)
	}
	if decision.ContextPrefix != "[work/decisions]" {
		t.Errorf("decision ContextPrefix = %q, want plain tag", decision.ContextPrefix)
	}
}

func TestClosedDateFromPath(t *testing.T) {
	if got := closedDateFromPath("/x/.hawp/work/closed/2026/07/05/TASK-003.md"); got != "2026-07-05" {
		t.Errorf("closedDateFromPath = %q", got)
	}
	if got := closedDateFromPath("/x/.hawp/work/active/TASK-001.md"); got != "" {
		t.Errorf("closedDateFromPath(active) = %q, want empty", got)
	}
}

func TestResolveIDHandlesLegacyAndUUIDFilenames(t *testing.T) {
	cases := []struct{ filename, want string }{
		{"TASK-001.md", "TASK-001"},
		{"0e1c4afa-9668-4d61-b5b6-1e27be42ca23.md", "0e1c4afa-9668-4d61-b5b6-1e27be42ca23"},
		{"2026-04-29-BUG-001-title.md", "BUG-001"},
		{"adr-001.md", ""},
	}
	for _, c := range cases {
		if got := resolveID(c.filename); got != c.want {
			t.Errorf("resolveID(%q) = %q, want %q", c.filename, got, c.want)
		}
	}
}
