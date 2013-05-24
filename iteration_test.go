package goquery

import (
	"code.google.com/p/go.net/html"
	"testing"
)

func TestEach(t *testing.T) {
	var cnt int

	sel := Doc().Find(".hero-unit .row-fluid").Each(func(i int, n *Selection) {
		cnt++
		t.Logf("At index %v, node %v", i, n.Nodes[0].Data)
	}).Find("a")

	if cnt != 4 {
		t.Errorf("Expected Each() to call function 4 times, got %v times.", cnt)
	}
	AssertLength(t, sel.Nodes, 6)
}

func TestEachWithBreak(t *testing.T) {
	var cnt int

	sel := Doc().Find(".hero-unit .row-fluid").EachWithBreak(func(i int, n *Selection) bool {
		cnt++
		t.Logf("At index %v, node %v", i, n.Nodes[0].Data)
		return false
	}).Find("a")

	if cnt != 1 {
		t.Errorf("Expected Each() to call function 1 time, got %v times.", cnt)
	}
	AssertLength(t, sel.Nodes, 6)
}

func TestEachEmptySelection(t *testing.T) {
	var cnt int

	sel := Doc().Find("zzzz")
	sel.Each(func(i int, n *Selection) {
		cnt++
	})
	if cnt > 0 {
		t.Error("Expected Each() to not be called on empty Selection.")
	}
	sel2 := sel.Find("div")
	AssertLength(t, sel2.Nodes, 0)
}

func TestMap(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	vals := sel.Map(func(i int, s *Selection) string {
		n := s.Get(0)
		if n.Type == html.ElementNode {
			return n.Data
		}
		return ""
	})
	for _, v := range vals {
		if v != "div" {
			t.Error("Expected Map array result to be all 'div's.")
		}
	}
	if len(vals) != 3 {
		t.Errorf("Expected Map array result to have a length of 3, found %v.", len(vals))
	}
}
