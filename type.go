package goquery

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/andybalholm/cascadia"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
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
	charset  *encoding.Decoder
}

// DecodeString decodes given string by this document's charset.
func (doc *Document) DecodeString(src string) string {
	if doc.charset == nil {
		return src
	}
	if dest, e := doc.charset.String(src); e == nil {
		return dest
	}
	return src
}

// NewDocumentFromNode is a Document constructor that takes a root html Node
// as argument.
func NewDocumentFromNode(root *html.Node) *Document {
	return newDocument(root, nil, nil)
}

// NewDocument is a Document constructor that takes a string URL as argument.
// It loads the specified document, parses it, and stores the root Document
// node, ready to be manipulated.
func NewDocument(url string) (*Document, error) {
	// Load the URL
	res, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	return NewDocumentFromResponse(res)
}

// NewDocumentFromReader returns a Document from a generic reader.
// It returns an error as second value if the reader's data cannot be parsed
// as html. It does *not* check if the reader is also an io.Closer, so the
// provided reader is never closed by this call, it is the responsibility
// of the caller to close it if required.
func NewDocumentFromReader(r io.Reader) (*Document, error) {
	return parseReader(r, nil)
}

// NewDocumentFromResponse is another Document constructor that takes an http response as argument.
// It loads the specified response's document, parses it, and stores the root Document
// node, ready to be manipulated. The response's body is closed on return.
func NewDocumentFromResponse(res *http.Response) (*Document, error) {
	if res == nil {
		return nil, errors.New("Response is nil")
	}
	defer res.Body.Close()
	if res.Request == nil {
		return nil, errors.New("Response.Request is nil")
	}

	return parseReader(res.Body, res.Request.URL)
}

func parseReader(r io.Reader, url *url.URL) (*Document, error) {
	b, _ := ioutil.ReadAll(r)
	enc, _, _ := charset.DetermineEncoding(b, "text/html")
	root, e := html.Parse(bytes.NewReader(b))
	if e != nil {
		return nil, e
	}
	return newDocument(root, url, enc.NewDecoder()), nil
}

// CloneDocument creates a deep-clone of a document.
func CloneDocument(doc *Document) *Document {
	return newDocument(cloneNode(doc.rootNode), doc.Url, nil)
}

// Private constructor, make sure all fields are correctly filled.
func newDocument(root *html.Node, url *url.URL, decoder *encoding.Decoder) *Document {
	// Create and fill the document
	d := &Document{Selection: nil, Url: url, rootNode: root, charset: decoder}
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
