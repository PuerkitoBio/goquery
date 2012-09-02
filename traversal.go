package goquery

import (
	"code.google.com/p/cascadia"
	"exp/html"
)

// Find() gets the descendants of each element in the current set of matched
// elements, filtered by a selector. It returns a new Selection object
// containing these matched elements.
func (this *Selection) Find(selector string) *Selection {
	return pushStack(this, findWithContext(selector, this.Nodes...))
}

func (this *Selection) FindSelection(sel *Selection) *Selection {
	return nil
}

func (this *Selection) FindNodes(nodes ...*html.Node) *Selection {
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

// Contents() gets the children of each element in the Selection,
// including text and comment nodes. It returns a new Selection object
// containing these elements.
func (this *Selection) Contents() *Selection {
	return pushStack(this, getSelectionChildren(this, false))
}

// ContentsFiltered() gets the children of each element in the Selection,
// filtered by the specified selector. It returns a new Selection
// object containing these elements. Since selectors only act on Element nodes,
// this function is an alias to ChildrenFiltered() unless the selector is empty,
// in which case it is an alias to Contents().
func (this *Selection) ContentsFiltered(selector string) *Selection {
	if selector != "" {
		return this.ChildrenFiltered(selector)
	}
	return this.Contents()
}

// Children() gets the child elements of each element in the Selection.
// It returns a new Selection object containing these elements.
func (this *Selection) Children() *Selection {
	return pushStack(this, getSelectionChildren(this, true))
}

// ChildrenFiltered() gets the child elements of each element in the Selection,
// filtered by the specified selector. It returns a new
// Selection object containing these elements.
func (this *Selection) ChildrenFiltered(selector string) *Selection {
	// Get the Children() unfiltered
	nodes := getSelectionChildren(this, true)
	// Create a temporary Selection to filter using winnow
	sel := &Selection{nodes, this.document, nil}
	// Filter based on selector
	nodes = winnow(sel, selector, true)
	// Push on the stack and return the "real" Selection
	return pushStack(this, nodes)
}

// Return the child nodes of each node in the Selection object, without
// duplicates.
func getSelectionChildren(s *Selection, elemOnly bool) (result []*html.Node) {
	for _, n := range s.Nodes {
		result = appendWithoutDuplicates(result, getChildren(n, elemOnly))
	}
	return
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
