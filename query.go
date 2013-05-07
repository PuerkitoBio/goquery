package goquery

import (
	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
	"regexp"
	"strings"
)

var rxClassTrim = regexp.MustCompile("[\t\r\n]")

// Is() checks the current matched set of elements against a selector and
// returns true if at least one of these elements matches.
func (this *Selection) Is(selector string) bool {
	if len(this.Nodes) > 0 {
		// Attempt a match with the selector
		cs := cascadia.MustCompile(selector)
		if len(this.Nodes) == 1 {
			return cs.Match(this.Nodes[0])
		} else {
			return len(cs.Filter(this.Nodes)) > 0
		}
	}

	return false
}

// IsFunction() checks the current matched set of elements against a predicate and
// returns true if at least one of these elements matches.
func (this *Selection) IsFunction(f func(int, *Selection) bool) bool {
	return this.FilterFunction(f).Length() > 0
}

// IsSelection() checks the current matched set of elements against a Selection object
// and returns true if at least one of these elements matches.
func (this *Selection) IsSelection(s *Selection) bool {
	return this.FilterSelection(s).Length() > 0
}

// IsNodes() checks the current matched set of elements against the specified nodes
// and returns true if at least one of these elements matches.
func (this *Selection) IsNodes(nodes ...*html.Node) bool {
	return this.FilterNodes(nodes...).Length() > 0
}

// HasClass() determines whether any of the matched elements are assigned the
// given class.
func (this *Selection) HasClass(class string) bool {
	class = " " + class + " "
	for _, n := range this.Nodes {
		// Applies only to element nodes
		if n.Type == html.ElementNode {
			if elClass, ok := getAttributeValue("class", n); ok {
				elClass = rxClassTrim.ReplaceAllString(" "+elClass+" ", " ")
				if strings.Index(elClass, class) > -1 {
					return true
				}
			}
		}
	}
	return false
}

// Contains() returns true if the specified Node is within,
// at any depth, one of the nodes in the Selection object.
// It is NOT inclusive, to behave like jQuery's implementation, and
// unlike Javascript's .contains(), so if the contained
// node is itself in the selection, it returns false.
func (this *Selection) Contains(n *html.Node) bool {
	return sliceContains(this.Nodes, n)
}
