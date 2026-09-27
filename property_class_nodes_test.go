package goquery

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestClassOperationsNonElementNodes(t *testing.T) {
	for _, operation := range []struct {
		name  string
		apply func(*Selection) *Selection
	}{
		{"AddClass", func(s *Selection) *Selection { return s.AddClass("marked") }},
		{"ToggleClass", func(s *Selection) *Selection { return s.ToggleClass("marked") }},
	} {
		for _, selection := range []struct {
			name            string
			selectNodes     func(*Document) *Selection
			containsElement bool
		}{
			{"document", func(d *Document) *Selection { return d.Selection }, false},
			{"doctype", func(d *Document) *Selection {
				return d.Contents().FilterFunction(func(_ int, s *Selection) bool { return s.Get(0).Type == html.DoctypeNode })
			}, false},
			{"text", func(d *Document) *Selection { return d.Find("div").Contents().Eq(0) }, false},
			{"comment", func(d *Document) *Selection { return d.Find("div").Contents().Eq(1) }, false},
			{"mixed content", func(d *Document) *Selection { return d.Find("div").Contents() }, true},
			{"element", func(d *Document) *Selection { return d.Find("b") }, true},
			{"empty", func(d *Document) *Selection { return d.Find("missing") }, false},
		} {
			t.Run(operation.name+"/"+selection.name, func(t *testing.T) {
				defer func() {
					if value := recover(); value != nil {
						t.Fatalf("class operation panicked: %v", value)
					}
				}()
				doc, err := NewDocumentFromReader(strings.NewReader(`<!doctype html><div>text<!-- note --><b class="kept">bold</b></div>`))
				if err != nil {
					t.Fatal(err)
				}
				sel := selection.selectNodes(doc)
				text := doc.Text()
				if operation.apply(sel) != sel {
					t.Fatal("class operation did not return the original selection")
				}
				for _, node := range sel.Nodes {
					if node.Type != html.ElementNode && len(node.Attr) != 0 {
						t.Errorf("non-element node gained attributes: %+v", node.Attr)
					}
				}
				want := "kept"
				if selection.containsElement {
					want += " marked"
				}
				if got := doc.Find("b").AttrOr("class", ""); got != want {
					t.Errorf("class: got %q, want %q", got, want)
				}
				if got := doc.Text(); got != text {
					t.Errorf("text changed: got %q, want %q", got, text)
				}
				if operation.name == "ToggleClass" {
					operation.apply(sel)
					if got := doc.Find("b").AttrOr("class", ""); got != "kept" {
						t.Errorf("second toggle: got %q, want kept", got)
					}
				}
			})
		}
	}
}
