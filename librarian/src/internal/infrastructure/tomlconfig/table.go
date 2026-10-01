// Package tomlconfig edits values in explicit TOML tables without reformatting
// the surrounding document. Unsupported layouts fail closed.
package tomlconfig

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"

	"github.com/pelletier/go-toml/v2"
)

// Rewrite updates only supplied fields; existing equal values remain byte-for-byte
// unchanged. The caller supplies only missing defaults, not existing policy keys.
func Rewrite(data []byte, path []string, fields map[string]any) ([]byte, error) {
	encodedFields, err := toml.Marshal(fields)
	if err != nil {
		return nil, err
	}
	var normalized map[string]any
	if err := toml.Unmarshal(encodedFields, &normalized); err != nil {
		return nil, err
	}
	fields = normalized
	var before map[string]any
	if err := toml.Unmarshal(data, &before); err != nil {
		return nil, fmt.Errorf("invalid TOML configuration")
	}
	table := before
	for _, key := range path {
		var ok bool
		table, ok = table[key].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected an explicit TOML table")
		}
	}
	values := map[string]string{}
	for key, value := range fields {
		if reflect.DeepEqual(table[key], value) {
			continue
		}
		encoded, err := toml.Marshal(map[string]any{key: value})
		if err != nil {
			return nil, err
		}
		_, raw, ok := bytes.Cut(encoded, []byte("="))
		if !ok {
			return nil, fmt.Errorf("unsupported replacement value")
		}
		values[key] = string(bytes.TrimSpace(raw))
	}
	edits, err := locateEdits(data, path, values)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
	out := append([]byte(nil), data...)
	for _, edit := range edits {
		out = append(append(append([]byte(nil), out[:edit.start]...), edit.text...), out[edit.end:]...)
	}
	var after map[string]any
	if err := toml.Unmarshal(out, &after); err != nil {
		return nil, fmt.Errorf("edited TOML is invalid; no configuration written")
	}
	for key, value := range fields {
		table[key] = value
	}
	if !reflect.DeepEqual(before, after) {
		return nil, fmt.Errorf("TOML edit changed unrelated values; no configuration written")
	}
	return out, nil
}
