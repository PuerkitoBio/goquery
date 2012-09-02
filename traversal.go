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

// TODO : Tests and doc for contents and children

func (this *Selection) Contents() *Selection {
	var matches []*html.Node

	for _, n := range this.Nodes {
		matches = appendWithoutDuplicates(matches, getChildren(n, false))
	}
	return pushStack(this, matches)
}

func (this *Selection) Children() *Selection {
	var matches []*html.Node

	for _, n := range this.Nodes {
		matches = appendWithoutDuplicates(matches, getChildren(n, true))
	}
	return pushStack(this, matches)
}

// Return the immediate children of the node, filtered on element nodes only
// if requested. The result is necessarily a slice of unique nodes.
func getChildren(n *html.Node, elemOnly bool) (result []*html.Node) {
	if n != nil {
		for _, c := range n.Child {
			if c.Type == html.ElementNode || !elemOnly {
				result = append(result, c)
			}
		}
	}
	return
}
