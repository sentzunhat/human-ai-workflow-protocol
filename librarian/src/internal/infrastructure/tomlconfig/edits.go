package tomlconfig

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2/unstable"
)

type replacement struct {
	start, end int
	text       string
}

// The version-pinned unstable AST API is confined here. Offsets refer to the
// original document; edits are applied in reverse order by Rewrite.
func locateEdits(data []byte, path []string, values map[string]string) ([]replacement, error) {
	parser := unstable.Parser{KeepComments: true}
	parser.Reset(data)
	var edits []replacement
	var current []string
	insertAt := -1
	for parser.NextExpression() {
		expr := parser.Expression()
		switch expr.Kind {
		case unstable.Table, unstable.ArrayTable:
			current = nodeKeys(expr)
			if reflect.DeepEqual(current, path) {
				if expr.Kind != unstable.Table {
					return nil, fmt.Errorf("array table requires manual merge")
				}
				keys := expr.Key()
				end := 0
				for keys.Next() {
					r := keys.Node().Raw
					end = int(r.Offset + r.Length)
				}
				insertAt = len(data)
				if n := bytes.IndexByte(data[end:], '\n'); n >= 0 {
					insertAt = end + n + 1
				}
			}
		case unstable.KeyValue:
			if !reflect.DeepEqual(current, path) {
				continue
			}
			keys := nodeKeys(expr)
			if len(keys) != 1 {
				continue
			}
			value, ok := values[keys[0]]
			if !ok {
				continue
			}
			// Refuse to discard comments embedded in a replaced multiline value.
			if containsComment(expr.Value()) {
				return nil, fmt.Errorf("comments inside changed value require manual merge")
			}
			key := expr.Key()
			key.Next()
			r := key.Node().Raw
			start := int(r.Offset + r.Length)
			end := int(expr.Raw.Offset + expr.Raw.Length)
			for start < end && strings.ContainsRune(" \t=", rune(data[start])) {
				start++
			}
			if start >= end {
				return nil, fmt.Errorf("invalid TOML value range")
			}
			edits = append(edits, replacement{start, end, value})
			delete(values, keys[0])
		}
	}
	if parser.Error() != nil {
		return nil, fmt.Errorf("cannot locate TOML values")
	}
	if insertAt < 0 {
		return nil, fmt.Errorf("inline or dotted HAWP table requires manual merge")
	}
	var keys []string
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var addition strings.Builder
	if len(keys) > 0 && insertAt > 0 && data[insertAt-1] != '\n' {
		addition.WriteByte('\n')
	}
	newline := "\n"
	if bytes.Contains(data, []byte("\r\n")) {
		newline = "\r\n"
	}
	for _, key := range keys {
		fmt.Fprintf(&addition, "%s = %s%s", key, values[key], newline)
	}
	if addition.Len() > 0 {
		edits = append(edits, replacement{insertAt, insertAt, addition.String()})
	}
	return edits, nil
}

func nodeKeys(node *unstable.Node) []string {
	var keys []string
	it := node.Key()
	for it.Next() {
		keys = append(keys, string(it.Node().Data))
	}
	return keys
}

func containsComment(node *unstable.Node) bool {
	if node.Kind == unstable.Comment {
		return true
	}
	it := node.Children()
	for it.Next() {
		if containsComment(it.Node()) {
			return true
		}
	}
	return false
}
