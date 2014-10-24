package goquery

import (
	"testing"
)

func TestAfter(t *testing.T) {
	doc := Doc2Clone()
	doc.Find("#main").After("#nf6")

	AssertLength(t, doc.Find("#main #nf6").Nodes, 0)
	AssertLength(t, doc.Find("#foot #nf6").Nodes, 0)
	AssertLength(t, doc.Find("#main + #nf6").Nodes, 1)
}

func TestAfterWithRemoved(t *testing.T) {
	doc := Doc2Clone()
	s := doc.Find("#main").Remove()
	s.After("#nf6")

	AssertLength(t, s.Find("#nf6").Nodes, 0)
	AssertLength(t, doc.Find("#nf6").Nodes, 0)
}

func TestAfterSelection(t *testing.T) {
	doc := Doc2Clone()
	doc.Find("#main").AfterSelection(doc.Find("#nf1, #nf2"))

	AssertLength(t, doc.Find("#main #nf1, #main #nf2").Nodes, 0)
	AssertLength(t, doc.Find("#foot #nf1, #foot #nf2").Nodes, 0)
	AssertLength(t, doc.Find("#main + #nf1, #nf1 + #nf2").Nodes, 2)
}

func TestAfterHtml(t *testing.T) {
	doc := Doc2Clone()
	doc.Find("#main").AfterHtml("<strong>new node</strong>")

	AssertLength(t, doc.Find("#main + strong").Nodes, 1)
}

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

func TestBefore(t *testing.T) {
	doc := Doc2Clone()
	doc.Find("#main").Before("#nf6")

	AssertLength(t, doc.Find("#main #nf6").Nodes, 0)
	AssertLength(t, doc.Find("#foot #nf6").Nodes, 0)
	AssertLength(t, doc.Find("body > #nf6:first-child").Nodes, 1)
}

func TestBeforeWithRemoved(t *testing.T) {
	doc := Doc2Clone()
	s := doc.Find("#main").Remove()
	s.Before("#nf6")

	AssertLength(t, s.Find("#nf6").Nodes, 0)
	AssertLength(t, doc.Find("#nf6").Nodes, 0)
}

func TestBeforeSelection(t *testing.T) {
	doc := Doc2Clone()
	doc.Find("#main").BeforeSelection(doc.Find("#nf1, #nf2"))

	AssertLength(t, doc.Find("#main #nf1, #main #nf2").Nodes, 0)
	AssertLength(t, doc.Find("#foot #nf1, #foot #nf2").Nodes, 0)
	AssertLength(t, doc.Find("body > #nf1:first-child, #nf1 + #nf2").Nodes, 2)
}

func TestBeforeHtml(t *testing.T) {
	doc := Doc2Clone()
	doc.Find("#main").BeforeHtml("<strong>new node</strong>")

	AssertLength(t, doc.Find("body > strong:first-child").Nodes, 1)
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
