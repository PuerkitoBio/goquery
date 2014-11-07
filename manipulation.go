package goquery

import (
	"fmt"
	"strings"

	"code.google.com/p/cascadia"
	"golang.org/x/net/html"
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

// After applies the selector from the root document and inserts the matched elements
// after the elements in the set of matched elements.
//
// If one of the matched elements in the selection is not currently in the
// document, it's impossible to insert nodes after it, so it will be ignored.
//
// This follows the same rules as Selection.Append.
func (s *Selection) After(selector string) *Selection {
	return s.AfterMatcher(cascadia.MustCompile(selector))
}

// AfterMatcher applies the matcher from the root document and inserts the matched elements
// after the elements in the set of matched elements.
//
// If one of the matched elements in the selection is not currently in the
// document, it's impossible to insert nodes after it, so it will be ignored.
//
// This follows the same rules as Selection.Append.
func (s *Selection) AfterMatcher(m Matcher) *Selection {
	return s.AfterNodes(m.MatchAll(s.document.rootNode)...)
}

// AfterSelection inserts the elements in the selection after each element in the set of matched
// elements.
// This follows the same rules as Selection.After.
func (s *Selection) AfterSelection(sel *Selection) *Selection {
	return s.AfterNodes(sel.Nodes...)
}

// AfterHtml parses the html and inserts it after the set of matched elements
// This follows the same rules as Selection.After.
func (s *Selection) AfterHtml(html string) *Selection {
	return s.AfterSelection(parseHtml(html))
}

// AfterNodes inserts the nodes after each element in the set of matched elements.
// This follows the same rules as Selection.After.
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
	return s.AppendMatcher(cascadia.MustCompile(selector))
}

// AppendMatcher applies the matcher from the root document, and append those nodes
// to the set of matched elements.
// This follows the same rules as Selection.Append.
func (s *Selection) AppendMatcher(m Matcher) *Selection {
	return s.AppendNodes(m.MatchAll(s.document.rootNode)...)
}

// AppendSelection appends the elements in the selection to the end of each element in the
// set of matched elements.
// This follows the same rules as Selection.Append.
func (s *Selection) AppendSelection(sel *Selection) *Selection {
	return s.AppendNodes(sel.Nodes...)
}

// AppendHtml parses the html and appends it to the set of matched elements.
func (s *Selection) AppendHtml(html string) *Selection {
	return s.AppendSelection(parseHtml(html))
}

// AppendNodes appends the specified nodes to each node in the set of matched elements.
// This follows the same rules as Selection.Append.
func (s *Selection) AppendNodes(ns ...*html.Node) *Selection {
	return s.manipulateNodes(ns, false, func(sn *html.Node, n *html.Node) {
		sn.AppendChild(n)
	})
}

// Before applies the selector from the root document, and inserts the matched elements
// before each element in the set of matched elements.
// This follows the same rules as Selection.After.
func (s *Selection) Before(selector string) *Selection {
	return s.BeforeMatcher(cascadia.MustCompile(selector))
}

// BeforeMatcher applies the matcher from the root document, and inserts the matched
// elements before each element in the set of matched elements.
// This follows the same rules as Selection.After.
func (s *Selection) BeforeMatcher(m Matcher) *Selection {
	return s.BeforeNodes(m.MatchAll(s.document.rootNode)...)
}

// BeforeSelection inserts the elements in the selection before each element in the set of matched
// elements.
// This follows the same rules as Selection.After.
func (s *Selection) BeforeSelection(sel *Selection) *Selection {
	return s.BeforeNodes(sel.Nodes...)
}

// BeforeHtml parses the html and inserts it before the set of matched elements.
// This follows the same rules as Selection.After.
func (s *Selection) BeforeHtml(html string) *Selection {
	return s.BeforeSelection(parseHtml(html))
}

// BeforeNodes inserts the nodes before each element in the set of matched elements.
// This follows the same rules as Selection.After.
func (s *Selection) BeforeNodes(ns ...*html.Node) *Selection {
	return s.manipulateNodes(ns, false, func(sn *html.Node, n *html.Node) {
		if sn.Parent != nil {
			sn.Parent.InsertBefore(n, sn)
		}
	})
}

// Clone creates a deep copy of the set of matched nodes. The new nodes will not be
// attached to the document.
func (s *Selection) Clone() *Selection {
	ns := newEmptySelection(s.document)
	ns.Nodes = cloneNodes(s.Nodes)
	return ns
}

// Empty removes all children nodes from the set of matched elements.
// Returns the children nodes in a new Selection.
func (s *Selection) Empty() *Selection {
	var nodes []*html.Node

	for _, n := range s.Nodes {
		for c := n.FirstChild; c != nil; c = n.FirstChild {
			n.RemoveChild(c)
			nodes = append(nodes, c)
		}
	}

	return pushStack(s, nodes)
}

// Remove removes the set of matched elements from the document.
// Returns the same selection, now consisting of nodes not in the document.
func (s *Selection) Remove() *Selection {
	for _, n := range s.Nodes {
		if n.Parent != nil {
			n.Parent.RemoveChild(n)
		}
	}

	return s
}

// RemoveFiltered removes the set of matched elements by selector.
// Returns the Selection of removed nodes.
func (s *Selection) RemoveFiltered(selector string) *Selection {
	return s.RemoveMatcher(cascadia.MustCompile(selector))
}

// RemoveMatcher removes the set of matched elements.
// Returns the Selection of removed nodes.
func (s *Selection) RemoveMatcher(m Matcher) *Selection {
	return s.FilterMatcher(m).Remove()
}
