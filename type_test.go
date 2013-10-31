package goquery

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// Test helper functions and members
var doc *Document
var doc2 *Document
var docB *Document
var docW *Document

func Doc() *Document {
	if doc == nil {
		doc = LoadDoc("page.html")
	}
	return doc
}
func Doc2() *Document {
	if doc2 == nil {
		doc2 = LoadDoc("page2.html")
	}
	return doc2
}
func DocB() *Document {
	if docB == nil {
		docB = LoadDoc("gotesting.html")
	}
	return docB
}
func DocW() *Document {
	if docW == nil {
		docW = LoadDoc("gowiki.html")
	}
	return docW
}

func AssertLength(t *testing.T, nodes []*html.Node, length int) {
	if len(nodes) != length {
		t.Errorf("Expected %d nodes, found %d.", length, len(nodes))
		for i, n := range nodes {
			t.Logf("Node %d: %+v.", i, n)
		}
	}
}

func AssertClass(t *testing.T, sel *Selection, class string) {
	if !sel.HasClass(class) {
		t.Errorf("Expected node to have class %s, found %+v.", class, sel.Get(0))
	}
}

func AssertPanic(t *testing.T) {
	if e := recover(); e == nil {
		t.Error("Expected a panic.")
	}
}

func AssertEqual(t *testing.T, s1 *Selection, s2 *Selection) {
	if s1 != s2 {
		t.Error("Expected selection objects to be the same.")
	}
}

func AssertSelectionIs(t *testing.T, sel *Selection, is ...string) {
	for i := 0; i < sel.Length(); i++ {
		if !sel.Eq(i).Is(is[i]) {
			t.Errorf("Expected node %d to be %s, found %+v", i, is[i], sel.Get(i))
		}
	}
}

func LoadDoc(page string) *Document {
	if f, e := os.Open(fmt.Sprintf("./testdata/%s", page)); e != nil {
		panic(e.Error())
	} else {
		defer f.Close()
		if node, e := html.Parse(f); e != nil {
			panic(e.Error())
		} else {
			return NewDocumentFromNode(node)
		}
	}
	return nil
}

func TestNewDocument(t *testing.T) {
	if f, e := os.Open("./testdata/page.html"); e != nil {
		t.Error(e.Error())
	} else {
		defer f.Close()
		if node, e := html.Parse(f); e != nil {
			t.Error(e.Error())
		} else {
			doc = NewDocumentFromNode(node)
		}
	}
}

func TestNewDocumentFromReader(t *testing.T) {
	cases := []struct {
		src string
		err bool
		sel string
		cnt int
	}{
		0: {
			src: `
<html>
<head>
<title>Test</title>
<body>
<h1>Hi</h1>
</body>
</html>`,
			sel: "h1",
			cnt: 1,
		},
		1: {
			// Actually pretty hard to make html.Parse return an error
			// based on content...
			src: `<html><body><aef<eqf>>>qq></body></ht>`,
		},
	}
	buf := bytes.NewBuffer(nil)

	for i, c := range cases {
		buf.Reset()
		buf.WriteString(c.src)

		d, e := NewDocumentFromReader(buf)
		if (e != nil) != c.err {
			if c.err {
				t.Errorf("[%d] - expected error, got none", i)
			} else {
				t.Errorf("[%d] - expected no error, got %s", i, e)
			}
		}
		if c.sel != "" {
			s := d.Find(c.sel)
			if s.Length() != c.cnt {
				t.Errorf("[%d] - expected %d nodes, found %d", i, c.cnt, s.Length())
			}
		}
	}
}

func TestNewDocumentFromText(t *testing.T) {
	// simple version
	doc, e := NewDocumentFromText("<html><head><title>Test</title><body><h1>Hi</h1></body></html>")
	if e != nil {
		t.Error(e.Error())
	}
	if doc.Find("title").Text() != "Test" {
		t.Error("Error parsing simple plain text version")
	}

	// real file
	page1 := readFile(t, "./testdata/page.html")

	doc, e = NewDocumentFromText(page1)
	if e != nil {
		t.Error(e.Error())
	}
	if doc.Find("div.span12 strong").Text() != "Beta Version." {
		t.Error("Error parsing file 'testdata/page.html'")
	}
}

func readFile(t *testing.T, file string) string {
	fileBytes, e := ioutil.ReadFile(file)
	if e != nil {
		t.Error(e.Error())
	}
	return string(fileBytes)
}
