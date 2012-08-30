package goquery

// Returns this (same Selection object)
func (this *Selection) Each(f func(int, *Selection)) *Selection {
	for i, n := range this.Nodes {
		f(i, newSingleSelection(n, this.document))
	}
	return this
}

// Map() passes each element in the current matched set through a function,
// producing a slice of string holding the returned values.
func (this *Selection) Map(f func(int, *Selection) string) (result []string) {
	for i, n := range this.Nodes {
		result = append(result, f(i, newSingleSelection(n, this.document)))
	}

	return result
}
