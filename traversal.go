package goquery

import (
	"code.google.com/p/cascadia"
	"exp/html"
)

type siblingType int

const (
	siblintPrevAll siblingType = iota - 2
	siblingPrev
	siblingAll
	siblingNext
	siblingNextAll
)

// Find() gets the descendants of each element in the current set of matched
// elements, filtered by a selector. It returns a new Selection object
// containing these matched elements.
func (this *Selection) Find(selector string) *Selection {
	return pushStack(this, findWithContext(selector, this.Nodes...))
}

// FindSelection() gets the descendants of each element in the current
// Selection, filtered by a Selection. It returns a new Selection object
// containing these matched elements.
func (this *Selection) FindSelection(sel *Selection) *Selection {
	if sel == nil {
		return pushStack(this, nil)
	}
	return this.FindNodes(sel.Nodes...)
}

// FindNodes() gets the descendants of each element in the current
// Selection, filtered by some nodes. It returns a new Selection object
// containing these matched elements.
func (this *Selection) FindNodes(nodes ...*html.Node) *Selection {
	var matches []*html.Node

	for _, n := range nodes {
		if sliceContains(this.Nodes, n) {
			matches = appendWithoutDuplicates(matches, []*html.Node{n})
		}
	}
	return pushStack(this, matches)
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
	// TODO : Refactor using siblings?
	return pushStack(this, getSelectionChildren(this, true))
}

// ChildrenFiltered() gets the child elements of each element in the Selection,
// filtered by the specified selector. It returns a new
// Selection object containing these elements.
func (this *Selection) ChildrenFiltered(selector string) *Selection {
	return filterAndPush(this, getSelectionChildren(this, true), selector)
}

// Parent() gets the parent of each element in the Selection. It returns a 
// new Selection object containing the matched elements.
func (this *Selection) Parent() *Selection {
	return pushStack(this, getParentNodes(this.Nodes))
}

// ParentFiltered() gets the parent of each element in the Selection filtered by a
// selector. It returns a new Selection object containing the matched elements.
func (this *Selection) ParentFiltered(selector string) *Selection {
	return filterAndPush(this, getParentNodes(this.Nodes), selector)
}

// Parents() gets the ancestors of each element in the current Selection. It
// returns a new Selection object with the matched elements.
func (this *Selection) Parents() *Selection {
	return pushStack(this, getParentsNodes(this.Nodes, "", nil))
}

// ParentsFiltered() gets the ancestors of each element in the current
// Selection. It returns a new Selection object with the matched elements.
func (this *Selection) ParentsFiltered(selector string) *Selection {
	return filterAndPush(this, getParentsNodes(this.Nodes, "", nil), selector)
}

// ParentsUntil() gets the ancestors of each element in the Selection, up to but
// not including the element matched by the selector. It returns a new Selection
// object containing the matched elements.
func (this *Selection) ParentsUntil(selector string) *Selection {
	return pushStack(this, getParentsNodes(this.Nodes, selector, nil))
}

// ParentsUntilSelection() gets the ancestors of each element in the Selection,
// up to but not including the elements in the specified Selection. It returns a
// new Selection object containing the matched elements.
func (this *Selection) ParentsUntilSelection(sel *Selection) *Selection {
	if sel == nil {
		return this.Parents()
	}
	return this.ParentsUntilNodes(sel.Nodes...)
}

// ParentsUntilNodes() gets the ancestors of each element in the Selection,
// up to but not including the specified nodes. It returns a
// new Selection object containing the matched elements.
func (this *Selection) ParentsUntilNodes(nodes ...*html.Node) *Selection {
	return pushStack(this, getParentsNodes(this.Nodes, "", nodes))
}

// ParentsFilteredUntil() is like ParentsUntil(), with the option to filter the
// results based on a selector string. It returns a new Selection
// object containing the matched elements.
func (this *Selection) ParentsFilteredUntil(filterSelector string, untilSelector string) *Selection {
	return filterAndPush(this, getParentsNodes(this.Nodes, untilSelector, nil), filterSelector)
}

// ParentsFilteredUntilSelection() is like ParentsUntilSelection(), with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (this *Selection) ParentsFilteredUntilSelection(filterSelector string, sel *Selection) *Selection {
	if sel == nil {
		return this.ParentsFiltered(filterSelector)
	}
	return this.ParentsFilteredUntilNodes(filterSelector, sel.Nodes...)
}

// ParentsFilteredUntilNodes() is like ParentsUntilNodes(), with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (this *Selection) ParentsFilteredUntilNodes(filterSelector string, nodes ...*html.Node) *Selection {
	return filterAndPush(this, getParentsNodes(this.Nodes, "", nodes), filterSelector)
}

// Siblings() gets the siblings of each element in the Selection. It returns
// a new Selection object containing the matched elements.
func (this *Selection) Siblings() *Selection {
	return pushStack(this, getSiblingNodes(this.Nodes, siblingAll))
}

// SiblingsFiltered() gets the siblings of each element in the Selection
// filtered by a selector. It returns a new Selection object containing the
// matched elements.
func (this *Selection) SiblingsFiltered(selector string) *Selection {
	return filterAndPush(this, getSiblingNodes(this.Nodes, siblingAll), selector)
}

// Next() gets the immediately following sibling of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (this *Selection) Next() *Selection {
	return pushStack(this, getSiblingNodes(this.Nodes, siblingNext))
}

// NextFiltered() gets the immediately following sibling of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (this *Selection) NextFiltered(selector string) *Selection {
	return filterAndPush(this, getSiblingNodes(this.Nodes, siblingNext), selector)
}

// NextAll() gets all the following siblings of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (this *Selection) NextAll() *Selection {
	return pushStack(this, getSiblingNodes(this.Nodes, siblingNextAll))
}

// NextAllFiltered() gets all the following siblings of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (this *Selection) NextAllFiltered(selector string) *Selection {
	return filterAndPush(this, getSiblingNodes(this.Nodes, siblingNextAll), selector)
}

// Prev() gets the immediately preceding sibling of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (this *Selection) Prev() *Selection {
	return pushStack(this, getSiblingNodes(this.Nodes, siblingPrev))
}

// PrevFiltered() gets the immediately preceding sibling of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (this *Selection) PrevFiltered(selector string) *Selection {
	return filterAndPush(this, getSiblingNodes(this.Nodes, siblingPrev), selector)
}

// Filter and push filters the nodes based on a selector, and pushes the results
// on the stack, with the srcSel as previous selection.
func filterAndPush(srcSel *Selection, nodes []*html.Node, selector string) *Selection {
	// Create a temporary Selection with the specified nodes to filter using winnow
	sel := &Selection{nodes, srcSel.document, nil}
	// Filter based on selector and push on stack
	return pushStack(srcSel, winnow(sel, selector, true))
}

// Internal implementation to get all parent nodes, stopping at the specified 
// node (or nil if no stop).
func getParentsNodes(nodes []*html.Node, stopSelector string, stopNodes []*html.Node) []*html.Node {
	return mapNodes(nodes, func(i int, n *html.Node) (result []*html.Node) {
		for p := n.Parent; p != nil; p = p.Parent {
			sel := newSingleSelection(p, nil)
			if stopSelector != "" {
				if sel.Is(stopSelector) {
					break
				}
			} else if len(stopNodes) > 0 {
				if sel.IsNodes(stopNodes...) {
					break
				}
			}
			if p.Type == html.ElementNode {
				result = append(result, p)
			}
		}
		return
	})
}

// Internal implementation of sibling nodes that return a raw slice of matches.
func getSiblingNodes(nodes []*html.Node, st siblingType) []*html.Node {
	return mapNodes(nodes, func(i int, n *html.Node) (result []*html.Node) {
		var prev *html.Node
		var nFound bool

		// Get the parent and loop through all children
		if p := n.Parent; p != nil {
			for _, c := range p.Child {
				// Care only about elements
				if c.Type == html.ElementNode {
					// Is it the existing node?
					if c == n {
						// Found the current node
						nFound = true
						if st == siblingPrev {
							// We want the previous node only, so append it and return
							if prev != nil {
								result = append(result, prev)
							}
							return
						}
					} else if prev == n && st == siblingNext {
						// We want only the next node and this is it, so append it and return
						result = append(result, c)
						return
					}
					// Keep child as previous
					prev = c

					// If child is not the current node, check if sibling type requires
					// to add it to the result.
					if c != n && (st == siblingAll || (st == siblintPrevAll && !nFound) ||
						(st == siblingNextAll && nFound)) {
						result = append(result, c)
					}
				}
			}
		}
		return
	})
}

// Internal implementation of parent nodes that return a raw slice of Nodes.
func getParentNodes(nodes []*html.Node) []*html.Node {
	return mapNodes(nodes, func(i int, n *html.Node) []*html.Node {
		if n.Parent != nil && n.Parent.Type == html.ElementNode {
			return []*html.Node{n.Parent}
		}
		return nil
	})
}

// Internal map function used by many traversing methods. Takes the source nodes
// to iterate on and the mapping function that returns an array of nodes.
// Returns an array of nodes mapped by calling the callback function once for
// each node in the source nodes.
func mapNodes(nodes []*html.Node, f func(int, *html.Node) []*html.Node) (result []*html.Node) {

	for i, n := range nodes {
		if vals := f(i, n); len(vals) > 0 {
			result = appendWithoutDuplicates(result, vals)
		}
	}

	return
}

// Private internal implementation of the Find() methods
func findWithContext(selector string, nodes ...*html.Node) []*html.Node {
	var matches []*html.Node

	// TODO : Refactor to use mapNodes?
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

// Return the child nodes of each node in the Selection object, without
// duplicates.
func getSelectionChildren(s *Selection, elemOnly bool) (result []*html.Node) {
	// TODO : Refactor to use mapNodes?
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
