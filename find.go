package goquery

import (
	"code.google.com/p/cascadia"
	"exp/html"
)

// Returns a new Selection object
func (this *Document) Find(selector string) *Selection {
	return &Selection{findWithContext(selector, this.Root), this, nil}
}

// Returns a new Selection object
func (this *Selection) Find(selector string) *Selection {
	return &Selection{findWithContext(selector, this.Nodes...), this.document, nil}
}

// Private internal implementation of the various Find() methods
func findWithContext(selector string, nodes ...*html.Node) []*html.Node {
	var matches []*html.Node

	sel := cascadia.MustCompile(selector)
	// Match the selector on each node
	for _, n := range nodes {
		matches = append(matches, sel.MatchAll(n)...)
	}
	return matches
}
