package goquery

import (
	"code.google.com/p/go.net/html"
	"errors"
	"github.com/djimenez/iconv-go"
	"io"
	"net/http"
	"net/url"
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

// NewConvertedDocument is a Document constructor that takes a string URL and a string charset as argument.
// It loads the specified document, converts the designated charset document to utf-8,
// parses it, and stores the root Document
// node, ready to be manipulated.
func NewConvertedDocument(url string, charset string) (*Document, error) {
	// Load the URL
	res, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	return NewConvertedDocumentFromResponse(res, charset)
}

// NewDocumentFromReader returns a Document from a generic reader.
// It returns an error as second value if the reader's data cannot be parsed
// as html. It does *not* check if the reader is also an io.Closer, so the
// provided reader is never closed by this call, it is the responsibility
// of the caller to close it if required.
func NewDocumentFromReader(r io.Reader) (*Document, error) {
	root, e := html.Parse(r)
	if e != nil {
		return nil, e
	}
	return newDocument(root, nil), nil
}

// NewDocumentFromResponse is another Document constructor that takes an http resonse as argument.
// It loads the specified response's document, parses it, and stores the root Document
// node, ready to be manipulated. The response's body is closed on return.
func NewDocumentFromResponse(res *http.Response) (*Document, error) {
	if res == nil {
		return nil, errors.New("Response is nil pointer")
	}

	defer res.Body.Close()

	// Parse the HTML into nodes
	root, e := html.Parse(res.Body)
	if e != nil {
		return nil, e
	}

	// Create and fill the document
	return newDocument(root, res.Request.URL), nil
}

// NewConvertedDocumentFromResponse is another Document constructor that takes an http resonse as argument.
// It loads the specified response's document, converts the designated charset document to utf-8,
// parses it, and stores the root Document
// node, ready to be manipulated. The response's body is closed on return.
func NewConvertedDocumentFromResponse(res *http.Response, charset string) (*Document, error) {
	if res == nil {
		return nil, errors.New("Response is nil pointer")
	}

	defer res.Body.Close()

	// Convert the designated charset HTML to utf-8 encoded HTML
	converted_body, err := iconv.NewReader(res.Body, charset, "utf-8")
	if err != nil {
		return nil, err
	}

	// Parse the HTML into nodes
	root, e := html.Parse(converted_body)
	if e != nil {
		return nil, e
	}

	// Create and fill the document
	return newDocument(root, res.Request.URL), nil
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
