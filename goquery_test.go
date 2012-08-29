package goquery

import (
	"exp/html"
	"os"
	"testing"
)

var doc *Document

func TestNewDocument(t *testing.T) {
	if f, e := os.Open("./testdata/page.html"); e != nil {
		t.Error(e.Error())
	} else {
		defer f.Close()
		if node, e := html.Parse(f); e != nil {
			t.Error(e.Error())
		} else {
			doc = NewDocumentFromNode(node)
		}
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
	sel.Each(func(i int, n *Selection) {
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

	sel := doc.Find(".hero-unit .row-fluid").Each(func(i int, n *Selection) {
		cnt++
		t.Logf("At index %v, node %v", i, n.Nodes[0].Data)
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

func TestChildren(t *testing.T) {
	sel := doc.Find(".pvk-content").Children()
	if len(sel.Nodes) != 13 {
		t.Errorf("Expected 13 child nodes, got %v.", len(sel.Nodes))
		for _, n := range sel.Nodes {
			t.Logf("%+v", n)
		}
	}
}

func TestChildrenFiltered(t *testing.T) {
	sel := doc.Find(".pvk-content").ChildrenFiltered(".hero-unit")
	if len(sel.Nodes) != 1 {
		t.Errorf("Expected 1 child nodes, got %v.", len(sel.Nodes))
		for _, n := range sel.Nodes {
			t.Logf("%+v", n)
		}
	}
}

func TestChildrenFilteredNone(t *testing.T) {
	sel := doc.Find(".pvk-content").ChildrenFiltered("a.btn")
	if len(sel.Nodes) != 0 {
		t.Errorf("Expected 0 child node, got %v.", len(sel.Nodes))
		for _, n := range sel.Nodes {
			t.Logf("%+v", n)
		}
	}
}

func TestFilter(t *testing.T) {
	sel := doc.Find(".span12").Filter(".alert")
	if len(sel.Nodes) != 1 {
		t.Errorf("Expected 1 node, found %v.", len(sel.Nodes))
	}
}

func TestFilterNone(t *testing.T) {
	sel := doc.Find(".span12").Filter(".zzalert")
	if sel.Nodes != nil {
		t.Error("Expected no node (nil), found some.")
	}
}

func TestFilterFunction(t *testing.T) {
	sel := doc.Find(".pvk-content").FilterFunction(func(i int, s *Selection) bool {
		return i > 0
	})
	if len(sel.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, found %v.", len(sel.Nodes))
	}
}

func TestFilterNode(t *testing.T) {
	sel := doc.Find(".pvk-content")
	sel2 := sel.FilterNode(sel.Nodes[2])
	if len(sel2.Nodes) != 1 {
		t.Errorf("Expected 1 node, found %v.", len(sel2.Nodes))
	}
}

func TestFilterSelection(t *testing.T) {
	sel := doc.Find(".link")
	sel2 := doc.Find("a[ng-click]")
	sel3 := sel.FilterSelection(sel2)
	if len(sel3.Nodes) != 1 {
		t.Errorf("Expected 1 node, found %v.", len(sel3.Nodes))
	}
}

func TestFirst(t *testing.T) {
	sel := doc.Find(".pvk-content").First()
	if len(sel.Nodes) != 1 {
		t.Errorf("Expected 1 node, found %v.", len(sel.Nodes))
	}
}

func TestFirstEmpty(t *testing.T) {
	sel := doc.Find(".pvk-zzcontentzz").First()
	if len(sel.Nodes) != 0 {
		t.Errorf("Expected 0 node, found %v.", len(sel.Nodes))
	}
}

func TestGet(t *testing.T) {
	sel := doc.Find(".pvk-content")
	node := sel.Get(1)
	if sel.Nodes[1] != node {
		t.Errorf("Expected node %v to be %v.", node, sel.Nodes[1])
	}
	node = sel.Get(-3)
	if sel.Nodes[0] != node {
		t.Errorf("Expected node %v to be %v.", node, sel.Nodes[0])
	}
	node = sel.Get(129)
	if node != nil {
		t.Error("Expected node to be nil.")
	}
}

func TestHasClass(t *testing.T) {
	sel := doc.Find("div")
	if !sel.HasClass("span12") {
		t.Error("Expected at least one div to have class span12.")
	}
}
