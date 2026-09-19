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
