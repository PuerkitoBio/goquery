package goquery

import (
	"code.google.com/p/cascadia"
	"exp/html"
	//"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// TODO : Benchmarks
// TODO : Ensure no node is added more than once in a selection (especially with Add...)
// TODO : Add the following methods:
// - Closest()
// - Parents()
// - Contains() (static function?)
// - Contents() (similar to Children(), but includes text and comment nodes, so Children() should filter them out)
// - End()
// - Eq()
// - Find() : Complete with Selection object and Node object as selectors
// - Has()

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
	var allChildren bool
	var sel cascadia.Selector
	var e error

	selector = strings.TrimSpace(selector)
	if selector == "*" || selector == "" {
		// Get all children
		allChildren = true
	} else {
		if sel, e = cascadia.Compile(selector); e != nil {
			// Selector doesn't compile, empty selection
			return nil
		}
	}

	for _, n := range nodes {
		for _, nchild := range n.Child {
			if allChildren || sel(nchild) {
				matches = append(matches, nchild)
			}
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
func (this *Selection) Each(f func(int, *Selection)) *Selection {
	for i, n := range this.Nodes {
		f(i, &Selection{[]*html.Node{n}, this.document})
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

func getAttributeValue(attrName string, n *html.Node) (val string, exists bool) {
	if n == nil {
		return
	}

	for _, a := range n.Attr {
		if a.Key == attrName {
			val = a.Val
			exists = true
			return
		}
	}
	return
}

// The Attr() method gets the attribute value for only the first element in the Selection.
// To get the value for each element individually, use a looping construct such as Each() or Map() method.
func (this *Selection) Attr(attrName string) (val string, exists bool) {
	if len(this.Nodes) == 0 {
		return
	}
	return getAttributeValue(attrName, this.Nodes[0])
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

func (this *Selection) Filter(selector string) *Selection {
	sel, e := cascadia.Compile(selector)
	if e != nil {
		// Selector doesn't compile, which means empty selection
		return &Selection{nil, this.document}
	}

	return &Selection{sel.Filter(this.Nodes), this.document}
}

func (this *Selection) FilterFunction(f func(int, *Selection) bool) *Selection {
	var matches []*html.Node

	// Check for a match for each current selection
	for i, n := range this.Nodes {
		if f(i, &Selection{[]*html.Node{n}, this.document}) {
			matches = append(matches, n)
		}
	}
	return &Selection{matches, this.document}
}

func (this *Selection) FilterNode(node *html.Node) *Selection {
	// TODO : Use Contains() on the this.Nodes array, if it contains, return node Selection, otherwise empty
	for _, n := range this.Nodes {
		if n == node {
			return &Selection{[]*html.Node{n}, this.document}
		}
	}
	return &Selection{nil, this.document}
}

func (this *Selection) FilterSelection(s *Selection) *Selection {
	// TODO : Exactly an Union() of two Selections, so maybe call it Union(), or have it as synonymous
	var matches []*html.Node

	if s == nil {
		return &Selection{nil, this.document}
	}

	// Check for a match for each current selection
	for _, n1 := range this.Nodes {
		for _, n2 := range s.Nodes {
			if n1 == n2 {
				matches = append(matches, n1)
				break
			}
		}
	}
	return &Selection{matches, this.document}
}

// Returns a new Selection object
func (this *Selection) First() *Selection {
	if len(this.Nodes) == 0 {
		return &Selection{nil, this.document}
	}
	return &Selection{[]*html.Node{this.Nodes[0]}, this.document}
}

// Get() without parameter is not implemented, its behaviour would be exactly the same as getting selection.Nodes
func (this *Selection) Get(index int) *html.Node {
	var l = len(this.Nodes)

	if index < 0 {
		index += l // Negative index gets from the end
	}
	if index >= 0 && index < l {
		return this.Nodes[index]
	}
	return nil
}

// Returns true if at least one node in the selection has the given class
func (this *Selection) HasClass(class string) bool {
	var rx = regexp.MustCompile("[\t\r\n]")

	class = " " + class + " "
	for _, n := range this.Nodes {
		// Applies only to element nodes
		if n.Type == html.ElementNode {
			if elClass, ok := getAttributeValue("class", n); ok {
				elClass = rx.ReplaceAllString(" "+elClass+" ", " ")
				if strings.Index(elClass, class) > -1 {
					return true
				}
			}
		}
	}
	return false
}
