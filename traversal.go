package goquery

import (
	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
)

type siblingType int

// Sibling type, used internally when iterating over children at the same
// level (siblings) to specify which nodes are requested.
const (
	siblingPrevUntil siblingType = iota - 3
	siblingPrevAll
	siblingPrev
	siblingAll
	siblingNext
	siblingNextAll
	siblingNextUntil
	siblingAllIncludingNonElements
)

// Find() gets the descendants of each element in the current set of matched
// elements, filtered by a selector. It returns a new Selection object
// containing these matched elements.
func (this *Selection) Find(selector string) *Selection {
	return pushStack(this, findWithSelector(this.Nodes, selector))
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
	return pushStack(this, mapNodes(nodes, func(i int, n *html.Node) []*html.Node {
		if sliceContains(this.Nodes, n) {
			return []*html.Node{n}
		}
		return nil
	}))
}

// Contents() gets the children of each element in the Selection,
// including text and comment nodes. It returns a new Selection object
// containing these elements.
func (this *Selection) Contents() *Selection {
	return pushStack(this, getChildrenNodes(this.Nodes, siblingAllIncludingNonElements))
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
	return pushStack(this, getChildrenNodes(this.Nodes, siblingAll))
}

// ChildrenFiltered() gets the child elements of each element in the Selection,
// filtered by the specified selector. It returns a new
// Selection object containing these elements.
func (this *Selection) ChildrenFiltered(selector string) *Selection {
	return filterAndPush(this, getChildrenNodes(this.Nodes, siblingAll), selector)
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

// Closest() gets the first element that matches the selector by testing the
// element itself and traversing up through its ancestors in the DOM tree.
func (this *Selection) Closest(selector string) *Selection {
	cs := cascadia.MustCompile(selector)

	return pushStack(this, mapNodes(this.Nodes, func(i int, n *html.Node) []*html.Node {
		// For each node in the selection, test the node itself, then each parent
		// until a match is found.
		for ; n != nil; n = n.Parent {
			if cs.Match(n) {
				return []*html.Node{n}
			}
		}
		return nil
	}))
}

// ClosestNodes() gets the first element that matches one of the nodes by testing the
// element itself and traversing up through its ancestors in the DOM tree.
func (this *Selection) ClosestNodes(nodes ...*html.Node) *Selection {
	return pushStack(this, mapNodes(this.Nodes, func(i int, n *html.Node) []*html.Node {
		// For each node in the selection, test the node itself, then each parent
		// until a match is found.
		for ; n != nil; n = n.Parent {
			if isInSlice(nodes, n) {
				return []*html.Node{n}
			}
		}
		return nil
	}))
}

// ClosestSelection() gets the first element that matches one of the nodes in the 
// Selection by testing the element itself and traversing up through its ancestors
// in the DOM tree.
func (this *Selection) ClosestSelection(s *Selection) *Selection {
	if s == nil {
		return pushStack(this, nil)
	}
	return this.ClosestNodes(s.Nodes...)
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
	return pushStack(this, getSiblingNodes(this.Nodes, siblingAll, "", nil))
}

// SiblingsFiltered() gets the siblings of each element in the Selection
// filtered by a selector. It returns a new Selection object containing the
// matched elements.
func (this *Selection) SiblingsFiltered(selector string) *Selection {
	return filterAndPush(this, getSiblingNodes(this.Nodes, siblingAll, "", nil), selector)
}

// Next() gets the immediately following sibling of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (this *Selection) Next() *Selection {
	return pushStack(this, getSiblingNodes(this.Nodes, siblingNext, "", nil))
}

// NextFiltered() gets the immediately following sibling of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (this *Selection) NextFiltered(selector string) *Selection {
	return filterAndPush(this, getSiblingNodes(this.Nodes, siblingNext, "", nil), selector)
}

// NextAll() gets all the following siblings of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (this *Selection) NextAll() *Selection {
	return pushStack(this, getSiblingNodes(this.Nodes, siblingNextAll, "", nil))
}

// NextAllFiltered() gets all the following siblings of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (this *Selection) NextAllFiltered(selector string) *Selection {
	return filterAndPush(this, getSiblingNodes(this.Nodes, siblingNextAll, "", nil), selector)
}

// Prev() gets the immediately preceding sibling of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (this *Selection) Prev() *Selection {
	return pushStack(this, getSiblingNodes(this.Nodes, siblingPrev, "", nil))
}

// PrevFiltered() gets the immediately preceding sibling of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (this *Selection) PrevFiltered(selector string) *Selection {
	return filterAndPush(this, getSiblingNodes(this.Nodes, siblingPrev, "", nil), selector)
}

// PrevAll() gets all the preceding siblings of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (this *Selection) PrevAll() *Selection {
	return pushStack(this, getSiblingNodes(this.Nodes, siblingPrevAll, "", nil))
}

// PrevAllFiltered() gets all the preceding siblings of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (this *Selection) PrevAllFiltered(selector string) *Selection {
	return filterAndPush(this, getSiblingNodes(this.Nodes, siblingPrevAll, "", nil), selector)
}

// NextUntil() gets all following siblings of each element up to but not
// including the element matched by the selector. It returns a new Selection
// object containing the matched elements.
func (this *Selection) NextUntil(selector string) *Selection {
	return pushStack(this, getSiblingNodes(this.Nodes, siblingNextUntil,
		selector, nil))
}

// NextUntilSelection() gets all following siblings of each element up to but not
// including the element matched by the Selection. It returns a new Selection
// object containing the matched elements.
func (this *Selection) NextUntilSelection(sel *Selection) *Selection {
	if sel == nil {
		return this.NextAll()
	}
	return this.NextUntilNodes(sel.Nodes...)
}

// NextUntilNodes() gets all following siblings of each element up to but not
// including the element matched by the nodes. It returns a new Selection
// object containing the matched elements.
func (this *Selection) NextUntilNodes(nodes ...*html.Node) *Selection {
	return pushStack(this, getSiblingNodes(this.Nodes, siblingNextUntil,
		"", nodes))
}

// PrevUntil() gets all preceding siblings of each element up to but not
// including the element matched by the selector. It returns a new Selection
// object containing the matched elements.
func (this *Selection) PrevUntil(selector string) *Selection {
	return pushStack(this, getSiblingNodes(this.Nodes, siblingPrevUntil,
		selector, nil))
}

// PrevUntilSelection() gets all preceding siblings of each element up to but not
// including the element matched by the Selection. It returns a new Selection
// object containing the matched elements.
func (this *Selection) PrevUntilSelection(sel *Selection) *Selection {
	if sel == nil {
		return this.PrevAll()
	}
	return this.PrevUntilNodes(sel.Nodes...)
}

// PrevUntilNodes() gets all preceding siblings of each element up to but not
// including the element matched by the nodes. It returns a new Selection
// object containing the matched elements.
func (this *Selection) PrevUntilNodes(nodes ...*html.Node) *Selection {
	return pushStack(this, getSiblingNodes(this.Nodes, siblingPrevUntil,
		"", nodes))
}

// NextFilteredUntil() is like NextUntil(), with the option to filter
// the results based on a selector string.
// It returns a new Selection object containing the matched elements.
func (this *Selection) NextFilteredUntil(filterSelector string, untilSelector string) *Selection {
	return filterAndPush(this, getSiblingNodes(this.Nodes, siblingNextUntil,
		untilSelector, nil), filterSelector)
}

// NextFilteredUntilSelection() is like NextUntilSelection(), with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (this *Selection) NextFilteredUntilSelection(filterSelector string, sel *Selection) *Selection {
	if sel == nil {
		return this.NextFiltered(filterSelector)
	}
	return this.NextFilteredUntilNodes(filterSelector, sel.Nodes...)
}

// NextFilteredUntilNodes() is like NextUntilNodes(), with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (this *Selection) NextFilteredUntilNodes(filterSelector string, nodes ...*html.Node) *Selection {
	return filterAndPush(this, getSiblingNodes(this.Nodes, siblingNextUntil,
		"", nodes), filterSelector)
}

// PrevFilteredUntil() is like PrevUntil(), with the option to filter
// the results based on a selector string.
// It returns a new Selection object containing the matched elements.
func (this *Selection) PrevFilteredUntil(filterSelector string, untilSelector string) *Selection {
	return filterAndPush(this, getSiblingNodes(this.Nodes, siblingPrevUntil,
		untilSelector, nil), filterSelector)
}

// PrevFilteredUntilSelection() is like PrevUntilSelection(), with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (this *Selection) PrevFilteredUntilSelection(filterSelector string, sel *Selection) *Selection {
	if sel == nil {
		return this.PrevFiltered(filterSelector)
	}
	return this.PrevFilteredUntilNodes(filterSelector, sel.Nodes...)
}

// PrevFilteredUntilNodes() is like PrevUntilNodes(), with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (this *Selection) PrevFilteredUntilNodes(filterSelector string, nodes ...*html.Node) *Selection {
	return filterAndPush(this, getSiblingNodes(this.Nodes, siblingPrevUntil,
		"", nodes), filterSelector)
}

// Filter and push filters the nodes based on a selector, and pushes the results
// on the stack, with the srcSel as previous selection.
func filterAndPush(srcSel *Selection, nodes []*html.Node, selector string) *Selection {
	// Create a temporary Selection with the specified nodes to filter using winnow
	sel := &Selection{nodes, srcSel.document, nil}
	// Filter based on selector and push on stack
	return pushStack(srcSel, winnow(sel, selector, true))
}

// Internal implementation of Find that return raw nodes.
func findWithSelector(nodes []*html.Node, selector string) []*html.Node {
	// Compile the selector once
	sel := cascadia.MustCompile(selector)
	// Map nodes to find the matches within the children of each node
	return mapNodes(nodes, func(i int, n *html.Node) (result []*html.Node) {
		// Go down one level, becausejQuery's Find() selects only within descendants
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				result = append(result, sel.MatchAll(c)...)
			}
		}
		return
	})
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
func getSiblingNodes(nodes []*html.Node, st siblingType, untilSelector string, untilNodes []*html.Node) []*html.Node {
	var f func(*html.Node) bool

	// If the requested siblings are ...Until(), create the test function to 
	// determine if the until condition is reached (returns true if it is)
	if st == siblingNextUntil || st == siblingPrevUntil {
		f = func(n *html.Node) bool {
			if untilSelector != "" {
				// Selector-based condition
				sel := newSingleSelection(n, nil)
				return sel.Is(untilSelector)
			} else if len(untilNodes) > 0 {
				// Nodes-based condition
				sel := newSingleSelection(n, nil)
				return sel.IsNodes(untilNodes...)
			}
			return false
		}
	}

	return mapNodes(nodes, func(i int, n *html.Node) []*html.Node {
		return getChildrenWithSiblingType(n.Parent, st, n, f)
	})
}

// Gets the children nodes of each node in the specified slice of nodes,
// based on the sibling type request.
func getChildrenNodes(nodes []*html.Node, st siblingType) []*html.Node {
	return mapNodes(nodes, func(i int, n *html.Node) []*html.Node {
		return getChildrenWithSiblingType(n, st, nil, nil)
	})
}

// Gets the children of the specified parent, based on the requested sibling
// type, skipping a specified node if required.
func getChildrenWithSiblingType(parent *html.Node, st siblingType, skipNode *html.Node,
	untilFunc func(*html.Node) bool) (result []*html.Node) {

	// Create the iterator function
	var iter = func(cur *html.Node) (ret *html.Node) {
		// Based on the sibling type requested, iterate the right way
		for {
			switch st {
			case siblingAll, siblingAllIncludingNonElements:
				if cur == nil {
					// First iteration, start with first child of parent
					// Skip node if required
					if ret = parent.FirstChild; ret == skipNode && skipNode != nil {
						ret = skipNode.NextSibling
					}
				} else {
					// Skip node if required
					if ret = cur.NextSibling; ret == skipNode && skipNode != nil {
						ret = skipNode.NextSibling
					}
				}
			case siblingPrev, siblingPrevAll, siblingPrevUntil:
				if cur == nil {
					// Start with previous sibling of the skip node
					ret = skipNode.PrevSibling
				} else {
					ret = cur.PrevSibling
				}
			case siblingNext, siblingNextAll, siblingNextUntil:
				if cur == nil {
					// Start with next sibling of the skip node
					ret = skipNode.NextSibling
				} else {
					ret = cur.NextSibling
				}
			default:
				panic("Invalid sibling type.")
			}
			if ret == nil || ret.Type == html.ElementNode || st == siblingAllIncludingNonElements {
				return
			} else {
				// Not a valid node, try again from this one
				cur = ret
			}
		}
		panic("Unreachable code reached.")
	}

	for c := iter(nil); c != nil; c = iter(c) {
		// If this is an ...Until() case, test before append (returns true
		// if the until condition is reached)
		if st == siblingNextUntil || st == siblingPrevUntil {
			if untilFunc(c) {
				return
			}
		}
		result = append(result, c)
		if st == siblingNext || st == siblingPrev {
			// Only one node was requested (immediate next or previous), so exit
			return
		}
	}
	return
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
