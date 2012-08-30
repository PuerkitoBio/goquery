package goquery

import (
	"code.google.com/p/cascadia"
	"exp/html"
	"strings"
)

// TODO : Filtered using Node and other Selection object

// Returns a new Selection object.
func (this *Document) Children() *Selection {
	return this.ChildrenFiltered("")
}

// Returns a new Selection object.
func (this *Selection) Children() *Selection {
	return this.ChildrenFiltered("")
}

// Returns a new Selection object.
func (this *Document) ChildrenFiltered(selector string) *Selection {
	return &Selection{childrenWithContext(selector, this.Root), this, nil}
}

// Returns a new Selection object.
func (this *Selection) ChildrenFiltered(selector string) *Selection {
	return &Selection{childrenWithContext(selector, this.Nodes...), this.document, nil}
}

func childrenWithContext(selector string, nodes ...*html.Node) []*html.Node {
	var matches []*html.Node
	var allChildren bool
	var sel cascadia.Selector

	selector = strings.TrimSpace(selector)
	if selector == "*" || selector == "" {
		// Get all children
		allChildren = true
	} else {
		sel = cascadia.MustCompile(selector)
	}

	for _, n := range nodes {
		for _, nchild := range n.Child {
			if allChildren || sel(nchild) {
				matches = append(matches, nchild)
			}
		}
	}
	return matches
}
