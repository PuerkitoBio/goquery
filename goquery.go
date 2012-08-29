package goquery

import (
	"code.google.com/p/cascadia"
	"exp/html"
	//"fmt"
	"net/http"
	"net/url"
)

// TODO : Ensure no node is added more than once in a selection (especially with Add...)
// TODO : Add the following methods:
// - Closest()
// - Parents()
// - Fix ChildrenFiltered, by forking Cascadia and adding a MatchSingle() method?

type Document struct {
	Root *html.Node
	Url  *url.URL
}

func NewDocumentFromNode(root *html.Node) (d *Document) {
	// Create and fill the document
	d = &Document{root, nil}
	return
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

func childrenWithContext(selector string, nodes ...*html.Node) []*html.Node {
	var matches []*html.Node
	//var allChildren bool
	//var sel *cascadia.Selector
	/*
		if selector == "*" || selector == "" {
			// Get all children
			allChildren = true
		} else {
			if sel, e := cascadia.Compile(selector); e != nil {
				// Selector doesn't compile, empty selection
				return nil
			}
		}
	*/
	for _, n := range nodes {
		for _, nchild := range n.Child {
			// TODO : At the moment, given Cascadia's API, cannot call Children with a selector string
			//if allChildren /*|| sel(nchild)*/ {
			matches = append(matches, nchild)
			//}
		}
	}
	return matches
}

// Returns a new Selection object
func (this *Document) Find(selector string) *Selection {
	return &Selection{findWithContext(selector, this.Root), this}
}

// Returns a new Selection object
func (this *Selection) Find(selector string) *Selection {
	return &Selection{findWithContext(selector, this.Nodes...), this.document}
}

// Returns this (same Selection object)
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

// Adds nodes of the specified Selection object to the current selection. Returns this (the same Selection object).
func (this *Selection) AddFromSelection(sel *Selection) *Selection {
	this.Nodes = append(this.Nodes, sel.Nodes...)
	return this
}

// The Attr() method gets the attribute value for only the first element in the Selection.
// To get the value for each element individually, use a looping construct such as Each() or Map() method.
func (this *Selection) Attr(attrName string) (val string, exists bool) {
	if this.Nodes == nil || len(this.Nodes) == 0 {
		return
	}
	for _, a := range this.Nodes[0].Attr {
		if a.Key == attrName {
			val = a.Val
			exists = true
			return
		}
	}
	return
}

// Returns a new Selection object.
func (this *Document) Children() *Selection {
	return this.ChildrenFiltered("")
}

// Returns a new Selection object.
func (this *Selection) Children() *Selection {
	return this.ChildrenFiltered("")
}

// Returns a new Selection object.
func (this *Document) ChildrenFiltered(selector string) *Selection {
	return &Selection{childrenWithContext(selector, this.Root), this}
}

// Returns a new Selection object.
func (this *Selection) ChildrenFiltered(selector string) *Selection {
	return &Selection{childrenWithContext(selector, this.Nodes...), this.document}
}
