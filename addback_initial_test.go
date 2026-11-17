package goquery

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestAddBackWithoutPreviousSelection(t *testing.T) {
	methods := []struct {
		name  string
		apply func(*Selection) *Selection
	}{
		{"filtered", func(s *Selection) *Selection { return s.AddBackFiltered("i") }},
		{"invalid selector", func(s *Selection) *Selection { return s.AddBackFiltered("[") }},
		{"matcher", func(s *Selection) *Selection { return s.AddBackMatcher(compileMatcher("i")) }},
		{"plain control", func(s *Selection) *Selection { return s.AddBack() }},
	}
	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			factories := []struct {
				name string
				make func() *Selection
			}{
				{"document", func() *Selection {
					d, e := NewDocumentFromReader(strings.NewReader("<i>x</i>"))
					if e != nil {
						t.Fatal(e)
					}
					return d.Selection
				}},
				{"exhausted stack", func() *Selection {
					d, e := NewDocumentFromReader(strings.NewReader("<i>x</i>"))
					if e != nil {
						t.Fatal(e)
					}
					return d.Find("i").End().End()
				}},
				{"standalone", func() *Selection { return &Selection{Nodes: []*html.Node{{Type: html.ElementNode, Data: "i"}}} }},
				{"empty", func() *Selection { return &Selection{} }},
			}
			for _, factory := range factories {
				t.Run(factory.name, func(t *testing.T) {
					s := factory.make()
					defer func() {
						if r := recover(); r != nil {
							t.Errorf("valid selection with no previous stack panicked: %v", r)
						}
					}()
					result := method.apply(s)
					if result == s {
						t.Error("result should be a new selection")
					}
					if result.Length() != s.Length() {
						t.Fatalf("length=%d,want%d", result.Length(), s.Length())
					}
					for i, n := range s.Nodes {
						if result.Nodes[i] != n {
							t.Errorf("node%d changed", i)
						}
					}
					if result.End() != s {
						t.Error("End did not restore original selection")
					}
				})
			}
		})
	}
}

func TestAddBackFilteredPreviousSelection(t *testing.T) {
	for _, method := range []struct {
		name  string
		apply func(*Selection) *Selection
	}{
		{"filtered", func(s *Selection) *Selection { return s.AddBackFiltered(".keep") }},
		{"matcher", func(s *Selection) *Selection { return s.AddBackMatcher(compileMatcher(".keep")) }},
	} {
		t.Run(method.name, func(t *testing.T) {
			d, e := NewDocumentFromReader(strings.NewReader(`<div class="keep"><i>child</i></div><p>omit</p>`))
			if e != nil {
				t.Fatal(e)
			}
			original := d.Find("div,p")
			current := original.Find("i")
			result := method.apply(current)
			if result.Length() != 2 {
				t.Fatalf("length=%d,want2", result.Length())
			}
			if result.Nodes[0] != current.Nodes[0] || result.Nodes[1] != original.Nodes[0] {
				t.Error("incorrect current or filtered previous nodes")
			}
			if result.End() != current {
				t.Error("End did not restore current selection")
			}
		})
	}
}
