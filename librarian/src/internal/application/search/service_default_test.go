package search

import "testing"

func TestDefaultServiceConfiguresEmbedderFactory(t *testing.T) {
	if DefaultService().embedder == nil {
		t.Fatal("DefaultService() has no embedder factory; hybrid search would silently become lexical-only")
	}
}
