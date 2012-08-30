package goquery

// Returns this (same Selection object)
func (this *Selection) Each(f func(int, *Selection)) *Selection {
	for i, n := range this.Nodes {
		f(i, newSingleSelection(n, this.document))
	}
	return this
}
