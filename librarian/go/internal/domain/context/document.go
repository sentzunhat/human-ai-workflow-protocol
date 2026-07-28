// Package context builds a folder-context enrichment layer over
// `.hawp/kit/` and `.hawp/work/`: before any chunk is embedded (see
// work item fbf12a93, Slice 2), each document is tagged with the
// context that makes it meaningful on its own — which kit section it
// belongs to, or a work record's type/status/closed-date — so a
// retrieved chunk carries "this is a closed feature from 2026-07-20",
// not just raw prose.
package context

// Corpus identifies which repo-local corpus a document came from.
type Corpus string

const (
	CorpusKit  Corpus = "kit"
	CorpusWork Corpus = "work"
)

// Document is one enriched markdown file, ready for Slice 2's
// chunking/embedding step to consume.
type Document struct {
	// RelPath is repo-relative (e.g. ".hawp/kit/usage/init.md"), the
	// portable identifier used in exported JSON.
	RelPath string `json:"relPath"`
	Corpus  Corpus `json:"corpus"`
	// Role is the folder classification: for kit, the top-level
	// subfolder name (or "root" for .hawp/kit/*.md); for work, the
	// backlog section or archival folder (active, closed, parked,
	// decisions, evidence, notes, status).
	Role string `json:"role"`
	// Type/Status/ClosedDate/ID are populated for work documents when a
	// matching backlog row or path-derived date is found; empty
	// otherwise (documents with no backlog row still get a Role).
	Type       string `json:"type,omitempty"`
	Status     string `json:"status,omitempty"`
	ClosedDate string `json:"closedDate,omitempty"`
	ID         string `json:"id,omitempty"`
	// ContextPrefix is the short tag Slice 2 prepends before chunking,
	// e.g. "[kit/usage]" or "[work/closed] TASK-086 (feature, closed 2026-07-03)".
	ContextPrefix string `json:"contextPrefix"`
	Content       string `json:"content"`
}
