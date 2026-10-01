package goquery

import (
	"strings"
	"testing"
)

func TestClassOperationsPreserveNonASCIIWhitespace(t *testing.T) {
	for _, whitespace := range []string{"\u00a0", "\u2003", "\u202f", "\v"} {
		for _, class := range []string{"a" + whitespace + "b", whitespace + "a", "a" + whitespace, whitespace} {
			for _, op := range []struct {
				name string
				run  func(*Selection, string) *Selection
				want string
			}{
				{"add", func(s *Selection, c string) *Selection { return s.AddClass(c) }, class + " keep"},
				{"remove", func(s *Selection, c string) *Selection { return s.RemoveClass(c) }, "keep"},
				{"toggle", func(s *Selection, c string) *Selection { return s.ToggleClass(c) }, "keep"},
			} {
				t.Run(op.name+"/"+class, func(t *testing.T) {
					doc, err := NewDocumentFromReader(strings.NewReader("<div></div>"))
					if err != nil {
						t.Fatal(err)
					}
					sel := doc.Find("div").SetAttr("class", class+" keep")
					if !sel.HasClass(class) {
						t.Fatal("existing class not recognized")
					}
					if op.run(sel, "\t"+class+"\r\n") != sel {
						t.Fatal("lost chaining")
					}
					if got := sel.AttrOr("class", ""); got != op.want {
						t.Fatalf("class = %q, want %q", got, op.want)
					}
					markup, err := OuterHtml(sel)
					if err != nil {
						t.Fatal(err)
					}
					reopened, err := NewDocumentFromReader(strings.NewReader(markup))
					if err != nil {
						t.Fatal(err)
					}
					if got := reopened.Find("div").AttrOr("class", ""); got != op.want {
						t.Fatalf("round trip = %q, want %q", got, op.want)
					}
				})
			}
		}
	}
}

func TestClassTokensASCIIWhitespace(t *testing.T) {
	for _, sep := range []string{" ", "\t", "\r", "\n", "\f"} {
		got := getClassesSlice(sep + "a\u00a0b" + sep + sep + "keep" + sep)
		if len(got) != 2 || got[0] != "a\u00a0b" || got[1] != "keep" {
			t.Errorf("separator %q: %q", sep, got)
		}
	}
}

func TestClassOperationsKeepUnicodeAtAttributeEdges(t *testing.T) {
	doc, err := NewDocumentFromReader(strings.NewReader("<div></div>"))
	if err != nil {
		t.Fatal(err)
	}
	sel := doc.Find("div")
	sel.AddClass("\u00a0a\u00a0").AddClass("keep").RemoveClass("keep")
	if got := sel.AttrOr("class", ""); got != "\u00a0a\u00a0" {
		t.Fatalf("class = %q", got)
	}
	sel.RemoveClass()
	if _, exists := sel.Attr("class"); exists {
		t.Fatal("RemoveClass() did not clear attribute")
	}
}
