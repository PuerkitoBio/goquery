package goquery

import (
	"bytes"

	"code.google.com/p/go.net/html"
)

// Attr gets the specified attribute's value for the first element in the
// Selection. To get the value for each element individually, use a looping
// construct such as Each or Map method.
func (s *Selection) Attr(attrName string) (val string, exists bool) {
	if len(s.Nodes) == 0 {
		return
	}
	return getAttributeValue(attrName, s.Nodes[0])
}

// Text gets the combined text contents of each element in the set of matched
// elements, including their descendants.
func (s *Selection) Text() string {
	var buf bytes.Buffer

	// Slightly optimized vs calling Each: no single selection object created
	for _, n := range s.Nodes {
		buf.WriteString(getNodeText(n))
	}
	return buf.String()
}

// Size is an alias for Length.
func (s *Selection) Size() int {
	return s.Length()
}

// Length returns the number of elements in the Selection object.
func (s *Selection) Length() int {
	return len(s.Nodes)
}

// Html gets the HTML contents of the first element in the set of matched
// elements. It includes text and comment nodes.
func (s *Selection) Html() (ret string, e error) {
	// Since there is no .innerHtml, the HTML content must be re-created from
	// the nodes using html.Render.
	var buf bytes.Buffer

	if len(s.Nodes) > 0 {
		for c := s.Nodes[0].FirstChild; c != nil; c = c.NextSibling {
			e = html.Render(&buf, c)
			if e != nil {
				return
			}
		}
		ret = buf.String()
	}

	return
}

// Get the specified node's text content.
func getNodeText(node *html.Node) string {
	if node.Type == html.TextNode {
		// Keep newlines and spaces, like jQuery
		return node.Data
	} else if node.FirstChild != nil {
		var buf bytes.Buffer
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			buf.WriteString(getNodeText(c))
		}
		return buf.String()
	}

	return ""
}

// Private function to get the specified attribute's value from a node.
func getAttributeValue(attrName string, n *html.Node) (val string, exists bool) {
	if n == nil {
		return
	}

	for _, a := range n.Attr {
		if a.Key == attrName {
			val = a.Val
			exists = true
			return
		}
	}
	return
}
