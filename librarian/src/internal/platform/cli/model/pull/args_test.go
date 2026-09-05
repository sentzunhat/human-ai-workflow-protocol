package pull

import "testing"

func TestParseModelPullArgs(t *testing.T) {
	for _, tc := range []struct {
		name       string
		args       []string
		repo, onnx string
		wantErr    bool
	}{
		{name: "missing repository", wantErr: true},
		{name: "repository only", args: []string{"org/repo"}, repo: "org/repo"},
		{name: "documented trailing option", args: []string{"org/repo", "--onnx-file", "onnx/model.onnx"}, repo: "org/repo", onnx: "onnx/model.onnx"},
		{name: "equals", args: []string{"org/repo", "--onnx-file=onnx/model.onnx"}, repo: "org/repo", onnx: "onnx/model.onnx"},
		{name: "leading option compatibility", args: []string{"--onnx-file", "model.onnx", "org/repo"}, repo: "org/repo", onnx: "model.onnx"},
		{name: "mixed options", args: []string{"--no-update-check", "org/repo", "--onnx-file", "model.onnx"}, repo: "org/repo", onnx: "model.onnx"},
		{name: "trailing bool", args: []string{"org/repo", "--no-update-check"}, repo: "org/repo"},
		{name: "only flags", args: []string{"--onnx-file", "model.onnx"}, wantErr: true},
		{name: "empty repository", args: []string{" "}, wantErr: true},
		{name: "missing trailing value", args: []string{"org/repo", "--onnx-file"}, wantErr: true},
		{name: "missing leading value", args: []string{"--onnx-file"}, wantErr: true},
		{name: "empty value", args: []string{"org/repo", "--onnx-file="}, wantErr: true},
		{name: "blank value", args: []string{"org/repo", "--onnx-file", " "}, wantErr: true},
		{name: "unknown option", args: []string{"org/repo", "--unknown"}, wantErr: true},
		{name: "extra repository", args: []string{"org/repo", "other/repo"}, wantErr: true},
		{name: "extra after option", args: []string{"org/repo", "--onnx-file=model.onnx", "other/repo"}, wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseModelPullArgs(tc.args)
			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tc.wantErr)
			}
			if !tc.wantErr && (got.modelRepo != tc.repo || got.onnxFile != tc.onnx) {
				t.Fatalf("got %+v", got)
			}
		})
	}
}
