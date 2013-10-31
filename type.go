package goquery

import (
	"code.google.com/p/go.net/html"
	"io"
	"net/http"
)

// Document represents an HTML document to be manipulated. Unlike jQuery, which
// is loaded as part of a DOM document, and thus acts upon its containing
// document, GoQuery doesn't know which HTML document to act upon. So it needs
// to be told, and that's what the Document class is for. It holds the root
// document node to manipulate, and can make selections on this document.
type Document struct {
	*Selection
	rootNode *html.Node
}

// NewDocumentFromNode() is a Document constructor that takes a root html Node
// as argument.
func NewDocumentFromNode(root *html.Node) (d *Document) {
	return newDocument(root)
}

// NewDocument() is a Document constructor that takes a string URL as argument.
// It loads the specified document, parses it, and stores the root Document
// node, ready to be manipulated.
func NewDocument(url string) (d *Document, e error) {
	// Load the URL
	res, e := http.Get(url)
	if e != nil {
		return
	}
	return NewDocumentFromResponse(res)
}

// NewDocumentFromReader() returns a Document from a generic reader.
// It returns an error as second value if the reader's data cannot be parsed
// as html. It does *not* check if the reader is also an io.Closer, so the
// provided reader is never closed by this call, it is the responsibility
// of the caller to close it if required.
func NewDocumentFromReader(r io.Reader) (d *Document, e error) {
	root, e := html.Parse(r)
	if e != nil {
		return nil, e
	}
	return newDocument(root), nil
}

// NewDocumentFromResponse() is another Document constructor that takes an http resonse as argument.
// It loads the specified response's document, parses it, and stores the root Document
// node, ready to be manipulated.
func NewDocumentFromResponse(res *http.Response) (d *Document, e error) {
	defer res.Body.Close()

	// Parse the HTML into nodes
	root, e := html.Parse(res.Body)
	if e != nil {
		return
	}

	// Create and fill the document
	return newDocument(root), nil
}

// Private constructor, make sure all fields are correctly filled.
func newDocument(root *html.Node) *Document {
	// Create and fill the document
	d := &Document{nil, root}
	d.Selection = newSingleSelection(root, d)
	return d
}

// Selection represents a collection of nodes matching some criteria. The
// initial Selection can be created by using Document.Find(), and then
// manipulated using the jQuery-like chainable syntax and methods.
type Selection struct {
	Nodes    []*html.Node
	document *Document
	prevSel  *Selection
}

// Helper constructor to create an empty selection
func newEmptySelection(doc *Document) *Selection {
	return &Selection{nil, doc, nil}
}

// Helper constructor to create a selection of only one node
func newSingleSelection(node *html.Node, doc *Document) *Selection {
	return &Selection{[]*html.Node{node}, doc, nil}
}
