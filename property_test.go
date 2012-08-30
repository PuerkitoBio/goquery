package goquery

import (
	"testing"
)

func TestAttrExists(t *testing.T) {
	EnsureDocLoaded()

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
