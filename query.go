package goquery

import (
	"exp/html"
	"regexp"
	"strings"
)

// Returns true if at least one node in the selection has the given class
func (this *Selection) HasClass(class string) bool {
	var rx = regexp.MustCompile("[\t\r\n]")

	class = " " + class + " "
	for _, n := range this.Nodes {
		// Applies only to element nodes
		if n.Type == html.ElementNode {
			if elClass, ok := getAttributeValue("class", n); ok {
				elClass = rx.ReplaceAllString(" "+elClass+" ", " ")
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

// Contains() returns true if the specified Node is within,
// at any depth, the root node of the Document object.
// It is NOT inclusive, to behave like jQuery's implementation, and
// unlike Javascript's .contains(), so if the contained
// node is itself in the selection, it returns false.
func (this *Document) Contains(n *html.Node) bool {
	return sliceContains([]*html.Node{this.Root}, n)
}
