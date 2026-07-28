package cli

import "testing"

// These only cover argument-parsing/usage-error paths — real model
// pull/embed inference is verified manually against the live network
// (see work item 748609a8), not in the unit-test suite.

func TestRunModelPullRequiresRepoArg(t *testing.T) {
	if err := Run([]string{"model", "pull"}); err == nil {
		t.Fatal("expected usage error when no repo is given")
	}
}

func TestRunEmbedRequiresTextArg(t *testing.T) {
	if err := Run([]string{"embed"}); err == nil {
		t.Fatal("expected usage error when no text is given")
	}
	if err := Run([]string{"embed", "--model", "org/repo"}); err == nil {
		t.Fatal("expected usage error when only flags are given, no text")
	}
}
