package goquery

import (
	"strings"
	"testing"
)

func TestIssue114(t *testing.T) {
	d, err := NewDocumentFromReader(strings.NewReader("<html><head><title>some title</title></head><body><h1>H1</h1></body></html>"))
	if err != nil {
		t.Error(err)
	}
	sel := d.Find("~")
	t.Log(sel)
}
