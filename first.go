package goquery

import (
	"exp/html"
)

// Returns a new Selection object
func (this *Selection) First() *Selection {
	if len(this.Nodes) == 0 {
		return &Selection{nil, this.document}
	}
	return &Selection{[]*html.Node{this.Nodes[0]}, this.document}
}
