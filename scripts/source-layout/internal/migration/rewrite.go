package migration

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"
)

func rewrite(f *file, defs map[string]map[string]definition, packages map[string]string) ([]byte, error) {
	var edits []edit
	add := func(start, end token.Pos, s string) {
		edits = append(edits, edit{f.set.Position(start).Offset, f.set.Position(end).Offset, s})
	}
	used := map[string]bool{}
	ast.Inspect(f.tree, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			used[id.Name] = true
		}
		return true
	})
	fresh := func(base string) string {
		for n := 1; ; n++ {
			v := base + "_layout" + strconv.Itoa(n)
			if !used[v] {
				used[v] = true
				return v
			}
		}
	}
	for _, imp := range f.tree.Imports {
		old, _ := strconv.Unquote(imp.Path.Value)
		symbols, ok := defs[old]
		if !ok {
			continue
		}
		alias := packages[old]
		if imp.Name != nil {
			alias = imp.Name.Name
		}
		if alias == "." || alias == "_" {
			return nil, fmt.Errorf("review required for dot/side-effect import %s", old)
		}
		groups := map[string][]*ast.Ident{}
		var rewriteErr error
		ast.Inspect(f.tree, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			id, ok := sel.X.(*ast.Ident)
			if !ok || id.Name != alias || id.Obj != nil {
				return true
			}
			d, found := symbols[sel.Sel.Name]
			if !found {
				rewriteErr = fmt.Errorf("unresolved %s.%s", old, sel.Sel.Name)
				return false
			}
			groups[d.destination] = append(groups[d.destination], id)
			return true
		})
		if rewriteErr != nil {
			return nil, rewriteErr
		}
		if len(groups) == 0 {
			continue
		}
		paths := make([]string, 0, len(groups))
		for p := range groups {
			paths = append(paths, p)
		}
		sort.Strings(paths)
		if len(paths) == 1 && paths[0] == old {
			continue
		}
		var specs []string
		for i, p := range paths {
			name := alias
			if i > 0 {
				name = fresh(alias)
			}
			specs = append(specs, name+" "+strconv.Quote(p))
			if name != alias {
				for _, id := range groups[p] {
					add(id.Pos(), id.End(), name)
				}
			}
		}
		// Parenthesize a single-line import when splitting it into several imports.
		spec := strings.Join(specs, "\n")
		for _, d := range f.tree.Decls {
			if g, ok := d.(*ast.GenDecl); ok && g.Tok == token.IMPORT && !g.Lparen.IsValid() && len(g.Specs) == 1 && g.Specs[0] == imp && len(specs) > 1 {
				spec = "(\n" + spec + "\n)"
			}
		}
		add(imp.Pos(), imp.End(), spec)
	}
	// Resolve cross-file, formerly same-package references after a package split.
	oldPkg := modulePath + "/" + path.Dir(f.Source)
	newPkg := modulePath + "/" + path.Dir(f.Destination)
	imports := map[string]string{}
	for _, id := range f.tree.Unresolved {
		d, ok := defs[oldPkg][id.Name]
		if !ok || d.destination == newPkg {
			continue
		}
		if !ast.IsExported(id.Name) {
			return nil, fmt.Errorf("private cross-package reference %s requires extraction", id.Name)
		}
		name := imports[d.destination]
		if name == "" {
			name = fresh(d.packageName)
			imports[d.destination] = name
		}
		add(id.Pos(), id.End(), name+"."+id.Name)
	}
	if len(imports) > 0 {
		paths := make([]string, 0, len(imports))
		for p := range imports {
			paths = append(paths, p)
		}
		sort.Strings(paths)
		extra := "\n"
		for _, p := range paths {
			extra += "import " + imports[p] + " " + strconv.Quote(p) + "\n"
		}
		add(f.tree.Name.End(), f.tree.Name.End(), extra)
	}
	if f.Source == "internal/domain/integration_test.go" || f.Source == "internal/domain/benchmark_test.go" {
		add(f.tree.Name.Pos(), f.tree.Name.End(), "models_test")
	}
	if len(edits) == 0 {
		return f.before, nil
	}
	sort.Slice(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
	out := append([]byte(nil), f.before...)
	boundary := len(out)
	for _, e := range edits {
		if e.end > boundary {
			return nil, errors.New("overlapping source edits")
		}
		out = append(append(append([]byte{}, out[:e.start]...), []byte(e.text)...), out[e.end:]...)
		boundary = e.start
	}
	return format.Source(out)
}
