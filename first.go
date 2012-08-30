package goquery

// Returns a new Selection object
func (this *Selection) First() *Selection {
	if len(this.Nodes) == 0 {
		return newEmptySelection(this.document)
	}
	return newSingleSelection(this.Nodes[0], this.document)
}
