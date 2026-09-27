package goquery

import (
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestSetAttrNonElementNodes(t *testing.T) {
	for _, nodeType := range []struct {
		name string
		kind html.NodeType
	}{
		{"document", html.DocumentNode},
		{"doctype", html.DoctypeNode},
		{"text", html.TextNode},
		{"comment", html.CommentNode},
	} {
		for _, existing := range []bool{false, true} {
			name := nodeType.name + "/new"
			if existing {
				name = nodeType.name + "/existing"
			}
			t.Run(name, func(t *testing.T) {
				node := &html.Node{Type: nodeType.kind}
				if existing {
					node.Attr = []html.Attribute{{Key: "data-mark", Val: "kept"}}
				}
				before := append([]html.Attribute(nil), node.Attr...)
				sel := &Selection{Nodes: []*html.Node{node}}
				if sel.SetAttr("data-mark", "changed") != sel {
					t.Fatal("SetAttr did not return the original selection")
				}
				if !reflect.DeepEqual(node.Attr, before) {
					t.Errorf("non-element attributes changed: got %v, want %v", node.Attr, before)
				}
			})
		}
	}
}

func TestSetAttrMixedContents(t *testing.T) {
	doc, err := NewDocumentFromReader(strings.NewReader(`<div>text<!-- note --><b data-mark="old">bold</b><i>italic</i></div>`))
	if err != nil {
		t.Fatal(err)
	}
	sel := doc.Find("div").Contents()
	text := sel.Text()
	if sel.SetAttr("data-mark", "changed") != sel {
		t.Fatal("SetAttr did not return the original selection")
	}
	for _, node := range sel.Nodes {
		var want []html.Attribute
		if node.Type == html.ElementNode {
			want = []html.Attribute{{Key: "data-mark", Val: "changed"}}
		}
		if !reflect.DeepEqual(node.Attr, want) {
			t.Errorf("node %q attributes: got %v, want %v", node.Data, node.Attr, want)
		}
	}
	if got := sel.Text(); got != text {
		t.Errorf("text changed: got %q, want %q", got, text)
	}
	empty := doc.Find("missing")
	if empty.SetAttr("data-mark", "changed") != empty || empty.Length() != 0 {
		t.Fatal("SetAttr changed an empty selection")
	}
}
