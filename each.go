package goquery

import (
	"exp/html"
)

// Returns this (same Selection object)
func (this *Selection) Each(f func(int, *Selection)) *Selection {
	for i, n := range this.Nodes {
		f(i, &Selection{[]*html.Node{n}, this.document})
	}
	return this
}
