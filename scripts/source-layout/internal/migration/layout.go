package migration

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"hawp-source-layout/internal/mapping"
)

const sourceRoot = "librarian/src"
const modulePath = "github.com/sentzunhat/hawp/librarian/src"

type entry struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Before      string `json:"before_sha256"`
	After       string `json:"after_sha256"`
	Mode        uint32 `json:"mode"`
	Reason      string `json:"reason"`
}
type plan struct {
	MappingRevision string   `json:"mapping_revision"`
	Version         int      `json:"version"`
	Module          string   `json:"module"`
	Files           []entry  `json:"files"`
	Remaining       []string `json:"remaining"`
}
type file struct {
	entry
	before, after []byte
	tree          *ast.File
	set           *token.FileSet
}
type definition struct{ destination, packageName string }
type edit struct {
	start, end int
	text       string
}

func digest(b []byte) string { return fmt.Sprintf("%x", sha256.Sum256(b)) }
func safeRelative(p string) bool {
	return p != "." && fs.ValidPath(p) && !strings.Contains(p, "\\") && !filepath.IsAbs(p)
}

// Inventory excludes only Git-ignored files; an ignored Go file is rejected so it
// cannot be silently omitted from compilation or from the reviewed mapping.
func inventory(root string) (map[string][]byte, map[string]uint32, error) {
	cmd := exec.Command("git", "-C", root, "ls-files", "-z", "--cached", "--others", "--exclude-standard", "--", sourceRoot)
	out, err := cmd.Output()
	if err != nil {
		return nil, nil, err
	}
	data := map[string][]byte{}
	modes := map[string]uint32{}
	for _, name := range strings.Split(string(out), "\x00") {
		if name == "" {
			continue
		}
		p := strings.TrimPrefix(name, sourceRoot+"/")
		if !safeRelative(p) {
			return nil, nil, fmt.Errorf("unsafe source: %s", name)
		}
		full := filepath.Join(root, filepath.FromSlash(name))
		// An old tracked path may be absent after a completed migration.
		info, err := os.Lstat(full)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, nil, err
		}
		if !info.Mode().IsRegular() {
			return nil, nil, fmt.Errorf("non-regular source: %s", name)
		}
		for dir := filepath.Dir(full); dir != root; dir = filepath.Dir(dir) {
			i, e := os.Lstat(dir)
			if e != nil {
				return nil, nil, e
			}
			if i.Mode()&os.ModeSymlink != 0 {
				return nil, nil, fmt.Errorf("symlink parent: %s", dir)
			}
		}
		b, err := os.ReadFile(full)
		if err != nil {
			return nil, nil, err
		}
		data[p] = b
		modes[p] = uint32(info.Mode().Perm())
	}
	err = filepath.WalkDir(filepath.Join(root, sourceRoot), func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink in source tree: %s", p)
		}
		if strings.HasSuffix(p, ".go") && !d.IsDir() {
			rel, _ := filepath.Rel(filepath.Join(root, sourceRoot), p)
			if _, ok := data[filepath.ToSlash(rel)]; !ok {
				return fmt.Errorf("unreviewed ignored Go source: %s", rel)
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	if !bytes.HasPrefix(data["go.mod"], []byte("module "+modulePath+"\n")) {
		return nil, nil, errors.New("unexpected source module")
	}
	return data, modes, nil
}

func prepare(root string) (plan, []*file, error) {
	data, modes, err := inventory(root)
	if err != nil {
		return plan{}, nil, err
	}
	keys := make([]string, 0, len(data))
	for p := range data {
		keys = append(keys, p)
	}
	sort.Strings(keys)
	files := make([]*file, 0, len(keys))
	destinations := map[string]bool{}
	defs := map[string]map[string]definition{}
	packages := map[string]string{}
	for _, p := range keys {
		dest, reason := mapping.Destination(p)
		// Case-fold as well: the working macOS filesystem can be case insensitive.
		if destinations[strings.ToLower(dest)] {
			return plan{}, nil, fmt.Errorf("destination collision: %s", dest)
		}
		destinations[strings.ToLower(dest)] = true
		f := &file{entry: entry{Source: p, Destination: dest, Before: digest(data[p]), Mode: modes[p], Reason: reason}, before: data[p]}
		if strings.HasSuffix(p, ".go") {
			f.set = token.NewFileSet()
			f.tree, err = parser.ParseFile(f.set, p, f.before, parser.ParseComments)
			if err != nil {
				return plan{}, nil, err
			}
			pkg := modulePath + "/" + path.Dir(p)
			packages[pkg] = strings.TrimSuffix(f.tree.Name.Name, "_test")
			if !strings.HasSuffix(p, "_test.go") {
				if defs[pkg] == nil {
					defs[pkg] = map[string]definition{}
				}
				for name := range f.tree.Scope.Objects {
					d := definition{modulePath + "/" + path.Dir(dest), f.tree.Name.Name}
					if old, ok := defs[pkg][name]; ok && old != d {
						return plan{}, nil, fmt.Errorf("ambiguous declaration: %s.%s", pkg, name)
					}
					defs[pkg][name] = d
				}
			}
		}
		files = append(files, f)
	}
	p := plan{Version: 1, MappingRevision: mapping.Revision, Module: modulePath, Remaining: []string{
		"domain/work normalization still has historical closed-record evidence review; folder normalization is clean",
		"domain/context, kit, kitsync, providersync and distribution still mix policy with I/O; imports alone cannot separate them",
		"application/context configuration and index storage wiring need consumer ports before further layer extraction",
		"MCP server tool handlers share private RPC contracts; keep them together until that contract is extracted",
		"No production file is deleted based on absent filename references; moved sources are removed only after verification",
	}}
	for _, f := range files {
		f.after = f.before
		if f.tree != nil {
			f.after, err = rewrite(f, defs, packages)
			if err != nil {
				return plan{}, nil, fmt.Errorf("%s: %w", f.Source, err)
			}
		}
		f.After = digest(f.after)
		p.Files = append(p.Files, f.entry)
	}
	if err := checkDestinations(root, p); err != nil {
		return plan{}, nil, err
	}
	return p, files, nil
}
