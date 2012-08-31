package goquery

import (
	"testing"
)

func TestHasClass(t *testing.T) {
	sel := Doc().Find("div")
	if !sel.HasClass("span12") {
		t.Error("Expected at least one div to have class span12.")
	}
}

func TestHasClassNone(t *testing.T) {
	sel := Doc().Find("h2")
	if sel.HasClass("toto") {
		t.Error("Expected h1 to have no class.")
	}
}

func TestHasClassNotFirst(t *testing.T) {
	sel := Doc().Find(".alert")
	if !sel.HasClass("alert-error") {
		t.Error("Expected .alert to also have class .alert-error.")
	}
}

func TestDocContains(t *testing.T) {
	sel := Doc().Find("h1")
	if !Doc().Contains(sel.Nodes[0]) {
		t.Error("Expected document to contain H1 tag.")
	}
}

func TestSelContains(t *testing.T) {
	sel := Doc().Find(".row-fluid")
	sel2 := Doc().Find("a[ng-click]")
	if !sel.Contains(sel2.Nodes[0]) {
		t.Error("Expected .row-fluid to contain a[ng-click] tag.")
	}
}

func TestSelNotContains(t *testing.T) {
	sel := Doc().Find("a.link")
	sel2 := Doc().Find("span")
	if sel.Contains(sel2.Nodes[0]) {
		t.Error("Expected a.link to NOT contain span tag.")
	}
}
