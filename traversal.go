package goquery

import (
	"code.google.com/p/cascadia"
	"exp/html"
	"strings"
)

// TODO : Maybe make the Document's Root return a Selection, Find() on the
// Document is not very useful (would need all the same funcs as Selection)
func (this *Document) Find(selector string) *Selection {
	return &Selection{findWithContext(selector, this.Root), this, nil}
}

// Find() gets the descendants of each element in the current set of matched
// elements, filtered by a selector. It returns a new Selection object
// containing these matched elements.
func (this *Selection) Find(selector string) *Selection {
	return pushStack(this, findWithContext(selector, this.Nodes...))
}

func (this *Selection) FindSelection(sel *Selection) *Selection {
	return nil
}

// Private internal implementation of the Find() methods
func findWithContext(selector string, nodes ...*html.Node) []*html.Node {
	var matches []*html.Node

	sel := cascadia.MustCompile(selector)
	// Match the selector on each node
	for _, n := range nodes {
		// Go down one level, becausejQuery's Find() selects only within descendants
		for _, c := range n.Child {
			if c.Type == html.ElementNode {
				matches = appendWithoutDuplicates(matches, sel.MatchAll(c))
			}
		}
	}
	return matches
}

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
