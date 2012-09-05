package goquery

import (
	"exp/html"
	"os"
	"testing"
)

// Test helper functions and members
var doc *Document

func Doc() *Document {
	if doc == nil {
		EnsureDocLoaded()
	}
	return doc
}

func AssertLength(t *testing.T, nodes []*html.Node, length int) {
	if len(nodes) != length {
		t.Errorf("Expected %d nodes, found %d.", length, len(nodes))
		for i, n := range nodes {
			t.Logf("Node %d: %+v.", i, n)
		}
	}
}

func AssertClass(t *testing.T, sel *Selection, class string) {
	if !sel.HasClass(class) {
		t.Errorf("Expected node to have class %s.", class)
	}
}

func AssertPanic(t *testing.T) {
	if e := recover(); e == nil {
		t.Error("Expected a panic.")
	}
}

func AssertEqual(t *testing.T, s1 *Selection, s2 *Selection) {
	if s1 != s2 {
		t.Error("Expected selection objects to be the same.")
	}
}

func EnsureDocLoaded() {
	if f, e := os.Open("./testdata/page.html"); e != nil {
		panic(e.Error())
	} else {
		defer f.Close()
		if node, e := html.Parse(f); e != nil {
			panic(e.Error())
		} else {
			doc = NewDocumentFromNode(node)
		}
	}
}

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
