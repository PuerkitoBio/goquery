package goquery

import (
	"bytes"
	"exp/html"
)

// Attr() gets the specified attribute's value for the first element in the
// Selection. To get the value for each element individually, use a looping
// construct such as Each() or Map() method.
func (this *Selection) Attr(attrName string) (val string, exists bool) {
	if len(this.Nodes) == 0 {
		return
	}
	return getAttributeValue(attrName, this.Nodes[0])
}

// Text() gets the combined text contents of each element in the set of matched
// elements, including their descendants.
func (this *Selection) Text() string {
	var buf bytes.Buffer

	for _, n := range this.Nodes {
		buf.WriteString(getNodeText(n))
	}
	return buf.String()
}

// Size() is an alias for Length().
func (this *Selection) Size() int {
	return this.Length()
}

// Length() returns the number of elements in the Selection object.
func (this *Selection) Length() int {
	return len(this.Nodes)
}

// Html() gets the HTML contents of the first element in the set of matched
// elements.
func (this *Selection) Html() (ret string, e error) {
	// Since there is no .innerHtml, the HTML content must be re-created from
	// the nodes usint html.Render().
	var buf bytes.Buffer

	if len(this.Nodes) > 0 {
		for _, c := range this.Nodes[0].Child {
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
		//ret = strings.Trim(node.Data, " \t\r\n")
		// Keep newlines and spaces, like jQuery
		return node.Data
	} else if len(node.Child) > 0 {
		var buf bytes.Buffer
		for _, c := range node.Child {
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
