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

func getNodeText(node *html.Node) (ret string) {
	if node.Type == html.TextNode {
		//ret = strings.Trim(node.Data, " \t\r\n")
		// Keep newlines and spaces, like jQuery
		ret = node.Data
	} else if len(node.Child) > 0 {
		for _, c := range node.Child {
			ret += getNodeText(c)
		}
	}
	return
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
