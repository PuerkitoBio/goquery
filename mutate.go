package goquery

import (
	"exp/html"
	"strings"
)

// Similar to append() in jQuery, except that it can only handle Selections.
// Note: like the jQuery version, the first target being appended to will have
// the node moved to it, and all subsequent targets will have clones inserted
func (this *Selection) Append(newKids *Selection) *Selection {
	for i, node := range this.Nodes {
		if i == 0 { // move the kid(s)
			for _, kid := range newKids.Nodes {
				if kid.Parent != nil {
					kid.Parent.RemoveChild(kid)
				}
				node.AppendChild(kid)
			}
		} else { // clone the kid(s)
			for _, kid := range newKids.Nodes {
				node.AppendChild(cloneNode(kid))
			}
		}
	}
	return this
}

// Appends a clone of each element in the template to each selected parent.
// AppendClones() isn't in the jQuery API, it was just handy.
func (this *Selection) AppendClones(template *html.Node) *Selection {
	for _, parent := range this.Nodes {
		parent.AppendChild(cloneNode(template))
	}
	return this
}

// Clone() returns a deep copy of the set of selected elements.
// This is the same behavior as jQuery's clone() function.
func (this *Selection) Clone() *Selection {
	results := newEmptySelection(this.document)
	this.Each(func(_ int, sel *Selection) {
		results = results.AddNodes(cloneNode(sel.Node()))
	})
	return results
}

func cloneNode(node *html.Node) *html.Node {
	result := &html.Node{
		Attr:     make([]html.Attribute, len(node.Attr)),
		Data:     node.Data,
		DataAtom: node.DataAtom,
		Type:     node.Type,
	}
	copy(result.Attr, node.Attr)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		result.AppendChild(cloneNode(child))
	}
	return result
}

// Remove all children from each selected node, just like empty() in jQuery.
func (this *Selection) Empty() *Selection {
	for _, node := range this.Nodes {
		for node.FirstChild != nil {
			node.RemoveChild(node.FirstChild)
		}
	}
	return this
}

// InsertBefore() inserts "this" before futureNextSib in the DOM.
// This is the same behavior as jQuery's insertBefore function.
func (this *Selection) InsertBefore(futureNextSib *Selection) *Selection {
	parent := futureNextSib.Node().Parent // TODO: this is a nilref bug, come up with a better way of finding the parent
	for _, n := range this.Nodes {
		if futureNextSib == nil || len(futureNextSib.Nodes) == 0 {
			parent.AppendChild(n)
		} else {
			parent.InsertBefore(n, futureNextSib.Node())
		}
	}
	return this
}

// Remove() removes all selected elements (and their decendents) from the DOM.
// This is the same behavior as jQuery's remove() function.
func (this *Selection) Remove() *Selection {
	for _, n := range this.Nodes {
		n.Parent.RemoveChild(n)
	}
	return this
}

// RemoveAttr removes all references to all attributes in the attrs string
// (space separated) in the selection.  Equivalent to jQuery's removeAttr.
func (this *Selection) RemoveAttr(attrs string) *Selection {
	names := strings.Split(attrs, " ")
	for _, node := range this.Nodes {
		for _, name := range names {
			removeAttr(node, name)
		}
	}
	return this
}

// removes an attribute from a node
func removeAttr(node *html.Node, attrName string) {
	for i := 0; i < len(node.Attr); i++ {
		if node.Attr[i].Key == attrName {
			last := len(node.Attr) - 1
			node.Attr[i] = node.Attr[last] // overwrite the target with the last attribute
			node.Attr = node.Attr[:last]   // then slice off the last attribute
			i--
		}
	}
}

// Adds an attribute key=value to each selected node.  Equivalent to jQuery's attr.
func (this *Selection) SetAttr(key, value string) *Selection {
	for _, node := range this.Nodes {
		removeAttr(node, key)
		node.Attr = append(node.Attr, html.Attribute{Key: key, Val: value})
	}
	return this
}

// Just like jQuery's html() setter.
func (this *Selection) SetHtml(s string) *Selection {
	result := newEmptySelection(this.document)
	for _, n := range this.Nodes {
		newNodes, e := html.ParseFragment(strings.NewReader(s), n)
		if e == nil {
			for _, child := range newNodes {
				n.AppendChild(child)
			}
			result.AddNodes(newNodes...)
		}
	}
	return result
}

// SetText() replaces the children of each selected node with the given text 
// (properly escaped of course).
// This is the same behavior as jQuery's .text() function.
func (this *Selection) SetText(s string) *Selection {
	for _, n := range this.Nodes {
		setNodeText(n, s)
	}
	return this
}

// Replace the given node's children with the given string.
func setNodeText(node *html.Node, s string) {
	// remove all existing children
	for node.FirstChild != nil {
		node.RemoveChild(node.FirstChild)
	}
	// add the text
	node.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: s,
	})
}
