package embed

import (
	"reflect"
	"testing"
)

func TestParseEmbedArgs(t *testing.T) {
	for _, tc := range []struct {
		name        string
		args        []string
		texts       []string
		model, onnx string
		wantErr     bool
	}{
		{name: "no text", wantErr: true},
		{name: "only flags", args: []string{"--model", "org/repo"}, wantErr: true},
		{name: "missing value", args: []string{"hello", "--model"}, wantErr: true},
		{name: "unknown flag", args: []string{"hello", "--unknown"}, wantErr: true},
		{name: "interspersed", args: []string{"first", "--model", "org/repo", "second", "--onnx-file", "model.onnx"}, texts: []string{"first", "second"}, model: "org/repo", onnx: "model.onnx"},
		{name: "multiple options", args: []string{"first", "--model=org/repo", "second", "--onnx-file=model.onnx", "third"}, texts: []string{"first", "second", "third"}, model: "org/repo", onnx: "model.onnx"},
		{name: "equals", args: []string{"--model=org/repo", "hello"}, texts: []string{"hello"}, model: "org/repo"},
		{name: "literal flags", args: []string{"--", "--model", "literal"}, texts: []string{"--model", "literal"}},
		{name: "terminator after option", args: []string{"--model", "org/repo", "--", "first", "--literal"}, texts: []string{"first", "--literal"}, model: "org/repo"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseEmbedArgs(tc.args)
			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			if !reflect.DeepEqual(got.texts, tc.texts) || got.modelRepo != tc.model || got.onnxFile != tc.onnx {
				t.Fatalf("got %+v", got)
			}
		})
	}
}
