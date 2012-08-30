package goquery

import (
	"exp/html"
)

// The Attr() method gets the attribute value for only the first element in the Selection.
// To get the value for each element individually, use a looping construct such as Each() or Map() method.
func (this *Selection) Attr(attrName string) (val string, exists bool) {
	if len(this.Nodes) == 0 {
		return
	}
	return getAttributeValue(attrName, this.Nodes[0])
}

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
