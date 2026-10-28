package goquery

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestSelectorForNode(t *testing.T) {
	const source = `<html><head></head><body><main><section id="content"><article class="card featured"><a>first</a></article><article class="card"><input type="email"></article><aside></aside></section></main></body></html>`
	doc, err := NewDocumentFromReader(strings.NewReader(source))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		node *html.Node
		want string
	}{
		{"id", doc.Find("section").Get(0), "#content"},
		{"classes distinguish siblings", doc.Find("article").First().Get(0), "#content > article.card.featured"},
		{"nth-child distinguishes siblings", doc.Find("article").Last().Get(0), "#content > article:nth-child(2)"},
		{"input type", doc.Find("input").Get(0), "#content > article:nth-child(2) > input[type=email]"},
		{"unique tag", doc.Find("aside").Get(0), "#content > aside"},
		{"structural root", doc.Find("body").Get(0), "body"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := SelectorForNode(test.node)
			if !ok {
				t.Fatal("SelectorForNode returned false")
			}
			if got != test.want {
				t.Errorf("got %q, want %q", got, test.want)
			}
			matches := doc.Find(got)
			if matches.Length() != 1 || matches.Get(0) != test.node {
				t.Errorf("selector %q did not resolve to the target node", got)
			}
		})
	}
}

func TestSelectorForNodeEscapesIdentifiers(t *testing.T) {
	const source = `<html><body><div id="123:x"></div><p class="a.b"></p><p class="plain"></p></body></html>`
	doc, err := NewDocumentFromReader(strings.NewReader(source))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		selector string
		want     string
	}{
		{"div", `#\31 23\:x`},
		{"p.a\\.b", `body > p.a\.b`},
	}

	for _, test := range tests {
		node := doc.Find(test.selector).Get(0)
		got, ok := SelectorForNode(node)
		if !ok {
			t.Fatalf("SelectorForNode(%q) returned false", test.selector)
		}
		if got != test.want {
			t.Errorf("SelectorForNode(%q) = %q, want %q", test.selector, got, test.want)
		}
		matches := doc.Find(got)
		if matches.Length() != 1 || matches.Get(0) != node {
			t.Errorf("selector %q did not resolve to the target node", got)
		}
	}
}

func TestSelectorForNodeIgnoresNonElementSiblings(t *testing.T) {
	const source = `<html><body><div><span>first</span>text<!-- comment --><span>second</span></div></body></html>`
	doc, err := NewDocumentFromReader(strings.NewReader(source))
	if err != nil {
		t.Fatal(err)
	}

	node := doc.Find("span").Last().Get(0)
	got, ok := SelectorForNode(node)
	if !ok {
		t.Fatal("SelectorForNode returned false")
	}
	if want := "body > div > span:nth-child(2)"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSelectorForNodeRejectsNonElements(t *testing.T) {
	for _, node := range []*html.Node{
		nil,
		{Type: html.TextNode, Data: "text"},
		{Type: html.CommentNode, Data: "comment"},
	} {
		if got, ok := SelectorForNode(node); ok || got != "" {
			t.Errorf("SelectorForNode(%v) = %q, %t; want empty, false", node, got, ok)
		}
	}
}

func TestEscapeCSSIdentifier(t *testing.T) {
	tests := map[string]string{
		"":       "",
		"simple": "simple",
		"123":    `\31 23`,
		"-1x":    `-\31 x`,
		"-":      `\-`,
		"a:b":    `a\:b`,
		"a b":    `a\ b`,
		"\x00":   "�",
		"\x1f":   `\1f `,
		"é":      "é",
	}

	for input, want := range tests {
		if got := escapeCSSIdentifier(input); got != want {
			t.Errorf("escapeCSSIdentifier(%q) = %q, want %q", input, got, want)
		}
	}
}
