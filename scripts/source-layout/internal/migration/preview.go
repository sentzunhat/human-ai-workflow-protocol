package migration

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

func previewReport(p plan) string {
	var out strings.Builder
	moves, updates := 0, 0
	for _, e := range p.Files {
		if e.Source != e.Destination {
			moves++
		}
		if e.Before != e.After {
			updates++
		}
	}
	fmt.Fprintf(&out, "# Source layout preview\n\nGenerated from the current source snapshot. **Proposed, not applied.**\n\n%d files inventoried; %d proposed moves; %d content updates (including moved files); %d retained paths. No dead-file deletions.\n\n", len(p.Files), moves, updates, len(p.Files)-moves)
	out.WriteString("This report records destinations, not proof that every layer is already pure. Run `--check` separately for candidate compilation and vet; run `--diff` to inspect exact code changes. All file paths below are relative to `librarian/src`.\n\n## Proposed directory tree\n\nCounts are files directly in each directory, not recursive totals. Empty folders are not created.\n\n```text\nlibrarian/src/\n")
	counts := map[string]int{}
	children := map[string]map[string]bool{}
	for _, e := range p.Files {
		dir := path.Dir(e.Destination)
		counts[dir]++
		for dir != "." {
			parent := path.Dir(dir)
			if children[parent] == nil {
				children[parent] = map[string]bool{}
			}
			children[parent][dir] = true
			dir = parent
		}
	}
	fmt.Fprintf(&out, "  (%d root files)\n", counts["."])
	var tree func(string, int)
	tree = func(parent string, depth int) {
		dirs := []string{}
		for d := range children[parent] {
			dirs = append(dirs, d)
		}
		sort.Strings(dirs)
		for _, d := range dirs {
			fmt.Fprintf(&out, "%s%s/ (%d files)\n", strings.Repeat("  ", depth), path.Base(d), counts[d])
			tree(d, depth+1)
		}
	}
	tree(".", 1)
	out.WriteString("```\n\n## Moves and ownership\n\n| Current file | Proposed file | Reason | Content changes |\n| --- | --- | --- | --- |\n")
	for _, e := range p.Files {
		if e.Source != e.Destination {
			fmt.Fprintf(&out, "| `%s` | `%s` | %s | %t |\n", e.Source, e.Destination, e.Reason, e.Before != e.After)
		}
	}
	out.WriteString("\n## Content updates without moves\n\nThese callers stay in their current folders; only their imports/references change.\n\n")
	for _, e := range p.Files {
		if e.Source == e.Destination && e.Before != e.After {
			fmt.Fprintf(&out, "- `%s`\n", e.Source)
		}
	}
	out.WriteString("\n## Unchanged files\n\nRetained with their existing owner; inclusion here does not claim complete architectural purity. Tests and assets are not treated as dead files.\n\n")
	for _, e := range p.Files {
		if e.Source == e.Destination && e.Before == e.After {
			fmt.Fprintf(&out, "- `%s`\n", e.Source)
		}
	}
	out.WriteString("\n## Still requires semantic work\n\n")
	for _, item := range p.Remaining {
		fmt.Fprintf(&out, "- %s\n", item)
	}
	return out.String()
}

// Compare temporary snapshots. Neither tree is installed in the checkout.
func previewDiff(files []*file) error {
	dir, err := os.MkdirTemp("", "hawp-source-layout-diff-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	before := make([]*file, 0, len(files))
	for _, f := range files {
		before = append(before, &file{entry: entry{Destination: f.Source, Mode: f.Mode}, after: f.before})
	}
	if err := writeFiles(filepath.Join(dir, "before", sourceRoot), before); err != nil {
		return err
	}
	if err := writeFiles(filepath.Join(dir, "after", sourceRoot), files); err != nil {
		return err
	}
	cmd := exec.Command("git", "diff", "--no-index", "--no-ext-diff", "--no-textconv", "--no-color", "--", "before", "after")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if e, ok := err.(*exec.ExitError); err != nil && (!ok || e.ExitCode() != 1) {
		return fmt.Errorf("preview diff: %w\n%s", err, output)
	}
	fmt.Print(string(output))
	return nil
}
