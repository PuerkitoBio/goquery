package goquery

import (
	"fmt"
	"strings"

	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
)

func parseHtml(html string) *Selection {
	// Errors are only returned when the io.Reader returns any error besides
	// EOF, but strings.Reader never will
	doc, err := NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		panic(fmt.Sprintf("Could not parse HTML: %s", err))
	}
	return doc.Find("body").Children()
}

// Deep copy a slice of nodes.
func cloneNodes(ns []*html.Node) []*html.Node {
	cns := make([]*html.Node, 0, len(ns))

	for _, n := range ns {
		cns = append(cns, cloneNode(n))
	}

	return cns
}

// Deep copy a node. The new node has clones of all the original node's
// children but none of its parents or siblings.
func cloneNode(n *html.Node) *html.Node {
	nn := &html.Node{
		Type:     n.Type,
		DataAtom: n.DataAtom,
		Data:     n.Data,
		Attr:     make([]html.Attribute, len(n.Attr)),
	}

	copy(nn.Attr, n.Attr)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nn.AppendChild(cloneNode(c))
	}

	return nn
}

func (s *Selection) manipulateNodes(
	ns []*html.Node,
	reverse bool,
	f func(sn *html.Node, n *html.Node)) *Selection {

	lasti := s.Size() - 1

	// net.Html doesn't provide document fragments for insertion, so to get
	// things in the correct order with After() and Prepend(), the callback
	// needs to be called on the reverse of the nodes.
	if reverse {
		for i, j := 0, len(ns)-1; i < j; i, j = i+1, j-1 {
			ns[i], ns[j] = ns[j], ns[i]
		}
	}

	for i, sn := range s.Nodes {
		for _, n := range ns {
			if i != lasti {
				f(sn, cloneNode(n))
			} else {
				if n.Parent != nil {
					n.Parent.RemoveChild(n)
				}
				f(sn, n)
			}
		}
	}

	return s
}

// From the root document, apply the selector, and insert the matched elements
// after element in the set of matched elements.
//
// If one of the matched elements in the selection is not currently in the
// document, it's impossible to insert nodes after it, so it will be ignored.
//
// This follows the same rules as Selection.Append().
func (s *Selection) After(selector string) *Selection {
	return s.AfterSelector(cascadia.MustCompile(selector))
}

// From the root document, apply the cascadia selector, and insert the matched
// elements after each element in the set of matched elements.
// This follows the same rules as Selection.After().
func (s *Selection) AfterSelector(cs cascadia.Selector) *Selection {
	return s.AfterNodes(cs.MatchAll(s.document.rootNode)...)
}

// Insert the elements in the selection after each element in the set of matched
// elements.
// This follows the same rules as Selection.After().
func (s *Selection) AfterSelection(sel *Selection) *Selection {
	return s.AfterNodes(sel.Nodes...)
}

// Parse the html and insert it after the set of matched elements
// This follows the same rules as Selection.After().
func (s *Selection) AfterHtml(html string) *Selection {
	return s.AfterSelection(parseHtml(html))
}

// Insert the nodes after each element in the set of matched elements.
// This follows the same rules as Selection.After().
func (s *Selection) AfterNodes(ns ...*html.Node) *Selection {
	return s.manipulateNodes(ns, true, func(sn *html.Node, n *html.Node) {
		if sn.Parent != nil {
			sn.Parent.InsertBefore(n, sn.NextSibling)
		}
	})
}

// Append the elements, specified by the selector, to the end of each element
// in the set of matched elements.
//
// Take note:
//
// 1) The selector is applied to the root document.
//
// 2) If any elements specified in the parameter are still part of the
// document, they will be moved to the new location.
//
// 3) If there are multiple locations to append to, cloned nodes will be
// appended to all target locations except the last, which will be moved
// as noted in (1).
func (s *Selection) Append(selector string) *Selection {
	return s.AppendSelector(cascadia.MustCompile(selector))
}

// From the root document, apply the cascadia selector, and append those nodes
// to the set of matched elements.
// This follows the same rules as Selection.Append().
func (s *Selection) AppendSelector(cs cascadia.Selector) *Selection {
	return s.AppendNodes(cs.MatchAll(s.document.rootNode)...)
}

// Append the elements in the selection to the end of each element in the
// set of matched elements.
// This follows the same rules as Selection.Append().
func (s *Selection) AppendSelection(sel *Selection) *Selection {
	return s.AppendNodes(sel.Nodes...)
}

// Parse the html and append it to the set of matched elements
func (s *Selection) AppendHtml(html string) *Selection {
	return s.AppendSelection(parseHtml(html))
}

// Append the specified nodes to each node in the set of matched elements.
// This follows the same rules as Selection.Append().
func (s *Selection) AppendNodes(ns ...*html.Node) *Selection {
	return s.manipulateNodes(ns, false, func(sn *html.Node, n *html.Node) {
		sn.AppendChild(n)
	})
}

// Create a deep copy of the set of matched nodes. The new nodes will not be
// attached to the document.
func (s *Selection) Clone() *Selection {
	ns := newEmptySelection(s.document)
	ns.Nodes = cloneNodes(s.Nodes)
	return ns
}

// Remove all children nodes from the set of matched elements.
// Returns the children nodes in a new Selection on the current Selection stack.
func (s *Selection) Empty() *Selection {
	nodes := make([]*html.Node, 0)

	for _, n := range s.Nodes {
		for c := n.FirstChild; c != nil; c = n.FirstChild {
			n.RemoveChild(c)
			nodes = append(nodes, c)
		}
	}

	return pushStack(s, nodes)
}

// Remove the set of matched elements from the document.
// Returns the same selection, now consisting of nodes not in the document.
func (s *Selection) Remove() *Selection {
	for _, n := range s.Nodes {
		if n.Parent != nil {
			n.Parent.RemoveChild(n)
		}
	}

	return s
}

// Filter the set of matched elements by selector before removing.
// Returns the filtered Selection.
func (s *Selection) RemoveFilter(selector string) *Selection {
	return s.RemoveFilterSelector(cascadia.MustCompile(selector))
}

// Filter the set of matched elements by cascadia selector before removing.
// Returns the filtered Selection.
func (s *Selection) RemoveFilterSelector(cs cascadia.Selector) *Selection {
	return s.FilterSelector(cs).Remove()
}
