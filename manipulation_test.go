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
