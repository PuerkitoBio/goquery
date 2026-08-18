package goquery

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var allNodes = `<!doctype html>
<html>
	<head>
		<meta a="b">
	</head>
	<body>
		<p><!-- this is a comment -->
		This is some text.
		</p>
		<div></div>
		<h1 class="header"></h1>
		<h2 class="header"></h2>
	</body>
</html>`

func TestNodeName(t *testing.T) {
	doc, err := NewDocumentFromReader(strings.NewReader(allNodes))
	if err != nil {
		t.Fatal(err)
	}

	n0 := doc.Nodes[0]
	nDT := n0.FirstChild
	sMeta := doc.Find("meta")
	nMeta := sMeta.Get(0)
	sP := doc.Find("p")
	nP := sP.Get(0)
	nComment := nP.FirstChild
	nText := nComment.NextSibling

	cases := []struct {
		node *html.Node
		typ  html.NodeType
		want string
	}{
		{n0, html.DocumentNode, nodeNames[html.DocumentNode]},
		{nDT, html.DoctypeNode, "html"},
		{nMeta, html.ElementNode, "meta"},
		{nP, html.ElementNode, "p"},
		{nComment, html.CommentNode, nodeNames[html.CommentNode]},
		{nText, html.TextNode, nodeNames[html.TextNode]},
	}
	for i, c := range cases {
		got := NodeName(newSingleSelection(c.node, doc))
		if c.node.Type != c.typ {
			t.Errorf("%d: want type %v, got %v", i, c.typ, c.node.Type)
		}
		if got != c.want {
			t.Errorf("%d: want %q, got %q", i, c.want, got)
		}
	}
}

func TestNodeNameMultiSel(t *testing.T) {
	doc, err := NewDocumentFromReader(strings.NewReader(allNodes))
	if err != nil {
		t.Fatal(err)
	}

	in := []string{"p", "h1", "div"}
	var out []string
	doc.Find(strings.Join(in, ", ")).Each(func(i int, s *Selection) {
		got := NodeName(s)
		out = append(out, got)
	})
	sort.Strings(in)
	sort.Strings(out)
	if !reflect.DeepEqual(in, out) {
		t.Errorf("want %v, got %v", in, out)
	}
}

var textNodes = `<!doctype html>
<html>
	<body>
		<div id="content">
			<h1>  Hello  </h1>
			<p>world</p>
			<!-- a comment -->
			<script>var ignored = 1;</script>
		</div>
	</body>
</html>`

func TestText_NilOptions(t *testing.T) {
	doc, err := NewDocumentFromReader(strings.NewReader(textNodes))
	if err != nil {
		t.Fatal(err)
	}

	sel := doc.Find("#content")
	// A nil options value must behave exactly like the Text method.
	if got, want := Text(sel, nil), sel.Text(); got != want {
		t.Errorf("nil options: want %q, got %q", want, got)
	}
}

func TestText_SeparatorAndTrim(t *testing.T) {
	doc, err := NewDocumentFromReader(strings.NewReader(textNodes))
	if err != nil {
		t.Fatal(err)
	}

	// Trim drops the whitespace-only text nodes coming from indentation, and
	// Separator joins the remaining fragments. The <script> text is included
	// because no Keep filter is provided.
	got := Text(doc.Find("#content"), &TextOptions{Separator: "|", Trim: true})
	want := "Hello|world|var ignored = 1;"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestText_Keep(t *testing.T) {
	doc, err := NewDocumentFromReader(strings.NewReader(textNodes))
	if err != nil {
		t.Fatal(err)
	}

	// Keep skips the text of <script> (and <style>) elements.
	keep := func(n *html.Node) bool {
		if n.Parent != nil && n.Parent.Type == html.ElementNode {
			switch n.Parent.Data {
			case "script", "style":
				return false
			}
		}
		return true
	}
	got := Text(doc.Find("#content"), &TextOptions{Separator: " ", Trim: true, Keep: keep})
	want := "Hello world"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestText_MultiSelection(t *testing.T) {
	doc, err := NewDocumentFromReader(strings.NewReader(textNodes))
	if err != nil {
		t.Fatal(err)
	}

	got := Text(doc.Find("h1, p"), &TextOptions{Separator: ",", Trim: true})
	want := "Hello,world"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestOuterHtml(t *testing.T) {
	doc, err := NewDocumentFromReader(strings.NewReader(allNodes))
	if err != nil {
		t.Fatal(err)
	}

	n0 := doc.Nodes[0]
	nDT := n0.FirstChild
	sMeta := doc.Find("meta")
	sP := doc.Find("p")
	nP := sP.Get(0)
	nComment := nP.FirstChild
	nText := nComment.NextSibling
	sHeaders := doc.Find(".header")

	cases := []struct {
		node *html.Node
		sel  *Selection
		want string
	}{
		{nDT, nil, "<!DOCTYPE html>"}, // render makes DOCTYPE all caps
		{nil, sMeta, `<meta a="b"/>`}, // and auto-closes the meta
		{nil, sP, `<p><!-- this is a comment -->
		This is some text.
		</p>`},
		{nComment, nil, "<!-- this is a comment -->"},
		{nText, nil, `
		This is some text.
		`},
		{nil, sHeaders, `<h1 class="header"></h1>`},
	}
	for i, c := range cases {
		if c.sel == nil {
			c.sel = newSingleSelection(c.node, doc)
		}
		got, err := OuterHtml(c.sel)
		if err != nil {
			t.Fatal(err)
		}

		if got != c.want {
			t.Errorf("%d: want %q, got %q", i, c.want, got)
		}
	}
}
