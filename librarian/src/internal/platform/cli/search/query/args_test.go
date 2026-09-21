package query

import "testing"

func TestParseArgsValid(t *testing.T) {
	for _, args := range [][]string{
		{"multi word query"},
		{"query", "--limit=5", "--context", "--format=json", "--max-tokens", "400", "-v", "--hybrid-ratio=0", "--no-update-check"},
		{"query", "--semantic", "--hybrid-ratio", "1"},
	} {
		opts, err := parseArgs(args)
		if err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		if opts.query != args[0] {
			t.Fatalf("query lost: %+v", opts)
		}
		if len(args) == 1 && (opts.limit != 10 || opts.maxTokens != 2000 || opts.hybridRatio != 0.3) {
			t.Fatalf("defaults changed: %+v", opts)
		}
		if len(args) > 5 && (opts.limit != 5 || !opts.context || opts.format != "json" || opts.maxTokens != 400 || !opts.verbose || opts.hybridRatio != 0) {
			t.Fatalf("options not applied: %+v", opts)
		}
	}
}

func TestParseArgsInvalid(t *testing.T) {
	for _, args := range [][]string{
		nil, {""}, {"   "}, {"--limit", "2"},
		{"q", "--limit"}, {"q", "--limit", "bad"}, {"q", "--limit", "0"}, {"q", "--limit", "-1"},
		{"q", "--max-tokens"}, {"q", "--max-tokens", "0"}, {"q", "--max-tokens", "bad"},
		{"q", "--format"}, {"q", "--format", "yaml"}, {"q", "--hybrid-ratio"},
		{"q", "--unknown"}, {"q", "extra"}, {"q", "--context=maybe"},
	} {
		if _, err := parseArgs(args); err == nil {
			t.Errorf("accepted invalid args: %v", args)
		}
	}
}
