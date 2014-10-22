package goquery

import (
	"testing"
)

func TestAppend(t *testing.T) {
	doc := Doc2Clone()
	doc.Find("#main").Append("#nf6")

	AssertLength(t, doc.Find("#foot #nf6").Nodes, 0)
	AssertLength(t, doc.Find("#main #nf6").Nodes, 1)
}

func TestAppendBody(t *testing.T) {
	doc := Doc2Clone()
	doc.Find("body").Append("#nf6")

	AssertLength(t, doc.Find("#foot #nf6").Nodes, 0)
	AssertLength(t, doc.Find("#main #nf6").Nodes, 0)
	AssertLength(t, doc.Find("body > #nf6").Nodes, 1)
}

func TestAppendSelection(t *testing.T) {
	doc := Doc2Clone()
	doc.Find("#main").AppendSelection(doc.Find("#nf1, #nf2"))

	AssertLength(t, doc.Find("#foot #nf1").Nodes, 0)
	AssertLength(t, doc.Find("#foot #nf2").Nodes, 0)
	AssertLength(t, doc.Find("#main #nf1").Nodes, 1)
	AssertLength(t, doc.Find("#main #nf2").Nodes, 1)
}

func TestAppendClone(t *testing.T) {
	doc := Doc2Clone()
	doc.Find("#n1").AppendSelection(doc.Find("#nf1").Clone())

	AssertLength(t, doc.Find("#foot #nf1").Nodes, 1)
	AssertLength(t, doc.Find("#main #nf1").Nodes, 1)
}

func TestAppendHtml(t *testing.T) {
	doc := Doc2Clone()
	doc.Find("div").AppendHtml("<strong>new node</strong>")

	AssertLength(t, doc.Find("strong").Nodes, 14)
}
func TestEmpty(t *testing.T) {
	doc := Doc2Clone()
	s := doc.Find("#main").Empty()

	AssertLength(t, doc.Find("#main").Children().Nodes, 0)
	AssertLength(t, s.Filter("div").Nodes, 6)
}

func TestRemove(t *testing.T) {
	doc := Doc2Clone()
	doc.Find("#nf1").Remove()

	AssertLength(t, doc.Find("#foot #nf1").Nodes, 0)
}

func TestRemoveAll(t *testing.T) {
	doc := Doc2Clone()
	doc.Find("*").Remove()

	AssertLength(t, doc.Find("*").Nodes, 0)
}

func TestRemoveRoot(t *testing.T) {
	doc := Doc2Clone()
	doc.Find("html").Remove()

	AssertLength(t, doc.Find("html").Nodes, 0)
}

func TestRemoveFilter(t *testing.T) {
	doc := Doc2Clone()
	nf6 := doc.Find("#nf6")
	s := doc.Find("div").RemoveFilter("#nf6")

	AssertLength(t, doc.Find("#nf6").Nodes, 0)
	AssertLength(t, s.Nodes, 1)
	if nf6.Nodes[0] != s.Nodes[0] {
		t.Error("Removed node does not match original")
	}
}
