package goquery

import (
	"exp/html"
	"net/http"
	"net/url"
)

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
