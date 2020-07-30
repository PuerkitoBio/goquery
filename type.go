package goquery

import (
	"io"
	"net/url"

	"github.com/andybalholm/cascadia"

	"golang.org/x/net/html"
)

// Document represents an HTML document to be manipulated. Unlike jQuery, which
// is loaded as part of a DOM document, and thus acts upon its containing
// document, GoQuery doesn't know which HTML document to act upon. So it needs
// to be told, and that's what the Document class is for. It holds the root
// document node to manipulate, and can make selections on this document.
type Document struct {
	*Selection
	Url      *url.URL
	rootNode *html.Node
}

// NewDocumentFromNode is a Document constructor that takes a root html Node
// as argument.
func NewDocumentFromNode(root *html.Node) *Document {
	return newDocument(root, nil)
}

// NewDocumentFromReader returns a Document from an io.Reader.
// It returns an error as second value if the reader's data cannot be parsed
// as html. It does not check if the reader is also an io.Closer, the
// provided reader is never closed by this call. It is the responsibility
// of the caller to close it if required.
func NewDocumentFromReader(r io.Reader) (*Document, error) {
	root, e := html.Parse(r)
	if e != nil {
		return nil, e
	}
	return newDocument(root, nil), nil
}

// CloneDocument creates a deep-clone of a document.
func CloneDocument(doc *Document) *Document {
	return newDocument(cloneNode(doc.rootNode), doc.Url)
}

// Private constructor, make sure all fields are correctly filled.
func newDocument(root *html.Node, url *url.URL) *Document {
	// Create and fill the document
	d := &Document{nil, url, root}
	d.Selection = newSingleSelection(root, d)
	return d
}

// Selection represents a collection of nodes matching some criteria. The
// initial Selection can be created by using Document.Find, and then
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

// Matcher is an interface that defines the methods to match
// HTML nodes against a compiled selector string. Cascadia's
// Selector implements this interface.
type Matcher interface {
	Match(*html.Node) bool
	MatchAll(*html.Node) []*html.Node
	Filter([]*html.Node) []*html.Node
}

// compileMatcher compiles the selector string s and returns
// the corresponding Matcher. If s is an invalid selector string,
// it returns a Matcher that fails all matches.
func compileMatcher(s string) Matcher {
	cs, err := cascadia.Compile(s)
	if err != nil {
		return invalidMatcher{}
	}
	return cs
}

// invalidMatcher is a Matcher that always fails to match.
type invalidMatcher struct{}

func (invalidMatcher) Match(n *html.Node) bool             { return false }
func (invalidMatcher) MatchAll(n *html.Node) []*html.Node  { return nil }
func (invalidMatcher) Filter(ns []*html.Node) []*html.Node { return nil }
