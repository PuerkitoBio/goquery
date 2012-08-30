package goquery

import (
	"code.google.com/p/cascadia"
	"exp/html"
	"strings"
)

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
	return &Selection{childrenWithContext(selector, this.Root), this}
}

// Returns a new Selection object.
func (this *Selection) ChildrenFiltered(selector string) *Selection {
	return &Selection{childrenWithContext(selector, this.Nodes...), this.document}
}

func childrenWithContext(selector string, nodes ...*html.Node) []*html.Node {
	var matches []*html.Node
	var allChildren bool
	var sel cascadia.Selector
	var e error

	selector = strings.TrimSpace(selector)
	if selector == "*" || selector == "" {
		// Get all children
		allChildren = true
	} else {
		if sel, e = cascadia.Compile(selector); e != nil {
			// Selector doesn't compile, empty selection
			return nil
		}
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
