package goquery

import (
	"exp/html"
	"os"
	"testing"
)

var doc *Document

func Doc() *Document {
	if doc == nil {
		EnsureDocLoaded()
	}
	return doc
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

func printNode(n *html.Node, t *testing.T) {
	t.Logf("Type: %v, Data: %v\n", n.Type, n.Data)
	for _, c := range n.Child {
		printNode(c, t)
	}
}

func TestPrintAll(t *testing.T) {
	//printNode(Doc().Root, t)
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
