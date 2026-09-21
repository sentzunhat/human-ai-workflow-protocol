package markdown

import "testing"

func TestExtractLinksAndPathRules(t *testing.T) {
	content := "[live](plan.md#check)\n```\n[ignored](missing.md)\n```\n![image](img.png)"
	links := ExtractLinks(BlankFences(content))
	if len(links) != 2 || links[0].Href != "plan.md#check" || !links[1].Image {
		t.Fatalf("ExtractLinks() = %#v", links)
	}
	if !IsLocalHref("plan.md#check") || IsLocalHref("https://example.test") || PathPart(links[0].Href) != "plan.md" {
		t.Fatal("Markdown link rules returned unexpected results")
	}
}

func TestBlankFencesPreservesByteOffsetsForUnicodeContent(t *testing.T) {
	content := "```\n日本語\n```\n[live](plan.md)"
	masked := BlankFences(content)
	links := ExtractLinks(masked)
	if len(links) != 1 {
		t.Fatalf("ExtractLinks() = %#v, want one link", links)
	}

	wantOffset := len("```\n日本語\n```\n")
	if links[0].Offset != wantOffset {
		t.Fatalf("link offset = %d, want %d", links[0].Offset, wantOffset)
	}
	if got := content[links[0].Offset:]; got != "[live](plan.md)" {
		t.Fatalf("original content at link offset = %q, want original link", got)
	}
}
