package goquery

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var invalidPathNodes = []struct {
	in   string
	path []int
}{
	{"<a>", []int{0, 1, 2}},
	{"<html><head><meta><title></title></head><body><div><p></p><a></a><span></span></div></body></html>", []int{0, 0, 1, 2, 0}},
	{"<html><head><meta><title></title></head><body><div><p></p><a></a><span></span></div></body></html>", []int{1}},
	{"<html><head><meta><title></title></head><body><div><p></p><a></a><span></span></div></body></html>", []int{1, 2}},
	{"<html><head><meta><title></title></head><body><div><p></p><a></a><span></span></div></body></html>", []int{1, 2, 10}},
	{"<a>", []int{}},
	{"<a>", []int{-1}},
}

var validPathNodes = []struct {
	in   string
	el   string
	path []int
}{
	{"<a>", "a", []int{0, 0, 1, 0}},                                                                                                      // root html body(1) a
	{"<html><head><meta></head><body></body></html>", "meta", []int{0, 0, 0, 0}},                                                         // root html head meta
	{"<html><head><meta><title></title></head><body></body></html>", "title", []int{0, 0, 0, 1}},                                         // root html head title
	{"<html><head><meta><title></title></head><body><div><p></p></div></body></html>", "div", []int{0, 0, 1, 0}},                         // root html body(1) div
	{"<html><head><meta><title></title></head><body><div><p></p></div></body></html>", "p", []int{0, 0, 1, 0, 0}},                        // root html body(1) div p
	{"<html><head><meta><title></title></head><body><div><p></p><a></a><span></span></div></body></html>", "a", []int{0, 0, 1, 0, 1}},    // root html body(1) div a(1)
	{"<html><head><meta><title></title></head><body><div><p></p><a></a><span></span></div></body></html>", "span", []int{0, 0, 1, 0, 2}}, // root html body(1) div span(2)
}

func TestPathForNode(t *testing.T) {
	for i, c := range validPathNodes {
		doc, err := NewDocumentFromReader(strings.NewReader(c.in))
		if err != nil {
			t.Errorf("%d: failed to parse: %v", i, err)
			continue
		}

		var n *html.Node
		if sel := doc.Find(c.el); sel.Length() > 0 {
			n = sel.Get(0)
		}

		got := PathForNode(n)
		if !reflect.DeepEqual(c.path, got) {
			h, _ := OuterHtml(doc.Selection)
			t.Errorf("%d: want %v, got %v (html: %s)", i, c.path, got, h)
		}
	}

	// test a nil node
	if got := PathForNode(nil); got != nil {
		t.Errorf("want nil for nil node, got %v", got)
	}
}

func TestNodeAtPath(t *testing.T) {
	// valid cases
	for i, c := range validPathNodes {
		n, err := html.Parse(strings.NewReader(c.in))
		if err != nil {
			t.Errorf("%d: failed to parse: %v", i, err)
			continue
		}

		nn := NodeAtPath(c.path, n)
		if nn == nil || nn.Data != c.el {
			t.Errorf("%d: want element %s, got %v", i, c.el, nn)
		}
	}

	// invalid cases
	for i, c := range invalidPathNodes {
		n, err := html.Parse(strings.NewReader(c.in))
		if err != nil {
			t.Errorf("%d: failed to parse: %v", i, err)
			continue
		}

		if got := NodeAtPath(c.path, n); got != nil {
			t.Errorf("%d: want nil, got %v", i, got)
		}
	}

	// test a nil node
	if got := NodeAtPath([]int{1, 2, 3}, nil); got != nil {
		t.Errorf("want nil for nil node, got %v", got)
	}

	// test empty and nil path
	n, _ := html.Parse(strings.NewReader("<a>"))
	if got := NodeAtPath(nil, n); got != nil {
		t.Errorf("want nil for nil path, got %v", got)
	}
	if got := NodeAtPath([]int{}, n); got != nil {
		t.Errorf("want nil for empty path, got %v", got)
	}
}

func TestSelectionPath(t *testing.T) {
	doc, err := NewDocumentFromReader(strings.NewReader("<html><body><div><p>hello</p></div></body></html>"))
	if err != nil {
		t.Fatal(err)
	}

	p := doc.Find("p")
	path := p.Path()
	want := []int{0, 0, 1, 0, 0}
	if !reflect.DeepEqual(path, want) {
		t.Errorf("want %v, got %v", want, path)
	}

	// empty selection
	empty := doc.Find("nonexistent")
	if got := empty.Path(); got != nil {
		t.Errorf("want nil for empty selection, got %v", got)
	}
}

func TestSelectionNodeAtPath(t *testing.T) {
	doc, err := NewDocumentFromReader(strings.NewReader("<html><body><div><p>hello</p></div></body></html>"))
	if err != nil {
		t.Fatal(err)
	}

	p := doc.Find("p")
	path := p.Path()

	// retrieve from doc
	sel := doc.NodeAtPath(path)
	if sel.Length() != 1 {
		t.Fatalf("want 1 node, got %d", sel.Length())
	}
	if sel.Nodes[0] != p.Nodes[0] {
		t.Errorf("want p node, got %v", sel.Nodes[0])
	}
	if sel.Text() != "hello" {
		t.Errorf("want text 'hello', got %q", sel.Text())
	}

	// retrieve from sub-selection
	div := doc.Find("div")
	sel2 := div.NodeAtPath(path)
	if sel2.Length() != 1 || sel2.Nodes[0] != p.Nodes[0] {
		t.Errorf("want p node from sub-selection, got %v", sel2.Nodes)
	}

	// invalid path
	if bad := doc.NodeAtPath([]int{9, 9, 9}); bad.Length() != 0 {
		t.Errorf("want empty selection for invalid path, got %d nodes", bad.Length())
	}

	// empty selection
	empty := doc.Find("nonexistent")
	if bad := empty.NodeAtPath(path); bad.Length() != 0 {
		t.Errorf("want empty selection from empty selection, got %d nodes", bad.Length())
	}
}

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

func TestText_EmptySelection(t *testing.T) {
	doc, err := NewDocumentFromReader(strings.NewReader(textNodes))
	if err != nil {
		t.Fatal(err)
	}

	// An empty selection has no text nodes to walk, so the result is empty
	// and the call must not panic.
	got := Text(doc.Find(".no-such-thing"), &TextOptions{Separator: " ", Trim: true})
	if got != "" {
		t.Errorf("want empty string, got %q", got)
	}
}

func TestText_KeepFiltersAll(t *testing.T) {
	doc, err := NewDocumentFromReader(strings.NewReader(textNodes))
	if err != nil {
		t.Fatal(err)
	}

	// When Keep rejects every text node, nothing is joined and the result is
	// empty rather than a string of separators.
	keep := func(n *html.Node) bool { return false }
	got := Text(doc.Find("#content"), &TextOptions{Separator: " ", Trim: true, Keep: keep})
	if got != "" {
		t.Errorf("want empty string, got %q", got)
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
