package goquery

import (
	"exp/html"
)

// Get() without parameter is not implemented, its behaviour would be exactly the same as getting selection.Nodes
func (this *Selection) Get(index int) *html.Node {
	var l = len(this.Nodes)

	if index < 0 {
		index += l // Negative index gets from the end
	}
	if index >= 0 && index < l {
		return this.Nodes[index]
	}
	return nil
}
