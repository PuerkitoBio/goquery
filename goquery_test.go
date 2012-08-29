package goquery

import (
	"exp/html"
	"testing"
)

var doc *Document

func TestNewDocument(t *testing.T) {
	var e error
	doc, e = NewDocument("http://provok.in")
	if e != nil {
		t.Error(e.Error())
	}
}

func TestFind(t *testing.T) {
	sel := doc.Find("div.row-fluid")
	if sel.Nodes == nil || len(sel.Nodes) != 9 {
		t.Errorf("Expected 9 matching nodes, found %v.", len(sel.Nodes))
	}
}

func TestFindInvalidSelector(t *testing.T) {
	var cnt int

	sel := doc.Find(":+ ^")
	if sel.Nodes != nil {
		t.Error("Expected a Selection object with Nodes == nil.")
	}
	sel.Each(func(i int, n *html.Node) {
		cnt++
	})
	if cnt > 0 {
		t.Error("Expected Each() to not be called on empty Selection.")
	}
	sel2 := sel.Find("div")
	if sel2.Nodes != nil {
		t.Error("Expected Find() on empty Selection to return an empty Selection.")
	}
}

func TestChainedFind(t *testing.T) {
	sel := doc.Find("div.hero-unit").Find(".row-fluid")
	if sel.Nodes == nil || len(sel.Nodes) != 4 {
		t.Errorf("Expected 4 matching nodes, found %v.", len(sel.Nodes))
	}
}

func TestEach(t *testing.T) {
	var cnt int

	sel := doc.Find(".hero-unit .row-fluid").Each(func(i int, n *html.Node) {
		cnt++
		t.Logf("At index %v, node %v", i, n.Data)
	}).Find("a")

	if cnt != 4 {
		t.Errorf("Expected Each() to call function 4 times, got %v times.", cnt)
	}
	if len(sel.Nodes) != 6 {
		t.Errorf("Expected 6 matching nodes, found %v.", len(sel.Nodes))
	}
}

func TestAdd(t *testing.T) {
	var cnt int

	sel := doc.Find("div.row-fluid")
	cnt = len(sel.Nodes)
	sel2 := sel.Add("a")
	if sel != sel2 {
		t.Error("Expected returned Selection from Add() to be same as 'this'.")
	}
	if len(sel.Nodes) <= cnt {
		t.Error("Expected nodes to be added to Selection.")
	}
	t.Logf("%v nodes after div.row-fluid and a were added.", len(sel.Nodes))
}

func TestAttrExists(t *testing.T) {
	if val, ok := doc.Find("a").Attr("href"); !ok {
		t.Error("Expected a value for the href attribute.")
	} else {
		t.Logf("Href of first anchor: %v.", val)
	}
}

func TestAttrNotExist(t *testing.T) {
	if val, ok := doc.Find("div.row-fluid").Attr("href"); ok {
		t.Errorf("Expected no value for the href attribute, got %v.", val)
	}
}
