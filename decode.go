package goquery

import "io"

// A Decoder takes an underlying source and maps it onto structs annotated with
// appropriate tags
type Decoder interface {
	Decode(interface{}) error
}

type ReaderDecoder struct {
	doc         *Document
	newDocError error
}

// NewDecoder returns a goquery.Decoder given any io.Reader
func NewDecoder(r io.Reader) *ReaderDecoder {
	doc, err := NewDocumentFromReader(r)

	return &ReaderDecoder{
		doc:         doc,
		newDocError: err,
	}
}

// Decode will properly unmarshal its underlying Document into the given
// interface.
func (r *ReaderDecoder) Decode(iface interface{}) error {
	if r.newDocError != nil {
		return r.newDocError
	}

	return UnmarshalDocument(r.doc, iface)
}
