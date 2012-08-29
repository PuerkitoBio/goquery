package goquery

import (
	"code.google.com/p/cascadia"
	"exp/html"
	"net/http"
	"net/url"
)

// TODO : Ensure no node is added more than once in a selection (especially with Add...)

type Document struct {
	Root *html.Node
	Url  *url.URL
}

func NewDocument(url string) (d *Document, e error) {
	// Load the URL
	res, e := http.Get(url)
	if e != nil {
		return
	}
	defer res.Body.Close()

	// Parse the HTML into nodes
	root, e := html.Parse(res.Body)
	if e != nil {
		return
	}

	// Create and fill the document
	d = &Document{root, res.Request.URL}
	return
}

type Selection struct {
	Nodes    []*html.Node
	document *Document
}

// Private internal implementation of the various Find() methods
func findWithContext(selector string, nodes ...*html.Node) []*html.Node {
	var matches []*html.Node

	sel, e := cascadia.Compile(selector)
	if e != nil {
		// Selector doesn't compile, which means empty selection
		return nil
	}

	// Match the selector on each node
	for _, n := range nodes {
		matches = append(matches, sel.MatchAll(n)...)
	}
	return matches
}

func (this *Document) Find(selector string) *Selection {
	return &Selection{findWithContext(selector, this.Root), this}
}

func (this *Selection) Find(selector string) *Selection {
	return &Selection{findWithContext(selector, this.Nodes...), this.document}
}

func (this *Selection) Each(f func(int, *html.Node)) *Selection {
	for i, n := range this.Nodes {
		f(i, n)
	}
	return this
}

// Adds matching nodes to the current selection. Returns the same Selection object.
// The new selector string is run in the context of the document of the Selection object.
func (this *Selection) Add(selector string) *Selection {
	this.Nodes = append(this.Nodes, findWithContext(selector, this.document.Root)...)
	return this
}

func (this *Selection) AddFromSelection(sel *Selection) *Selection {
	this.Nodes = append(this.Nodes, sel.Nodes...)
	return this
}
