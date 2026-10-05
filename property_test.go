package goquery

import (
	"regexp"
	"strings"
	"testing"
)

func TestAttrExists(t *testing.T) {
	if val, ok := Doc().Find("a").Attr("href"); !ok {
		t.Error("Expected a value for the href attribute.")
	} else {
		t.Logf("Href of first anchor: %v.", val)
	}
}

func TestAttrOr(t *testing.T) {
	if val := Doc().Find("a").AttrOr("fake-attribute", "alternative"); val != "alternative" {
		t.Error("Expected an alternative value for 'fake-attribute' attribute.")
	} else {
		t.Logf("Value returned for not existing attribute: %v.", val)
	}
	if val := Doc().Find("zz").AttrOr("fake-attribute", "alternative"); val != "alternative" {
		t.Error("Expected an alternative value for 'fake-attribute' on an empty selection.")
	} else {
		t.Logf("Value returned for empty selection: %v.", val)
	}
}

func TestAttrNotExist(t *testing.T) {
	if val, ok := Doc().Find("div.row-fluid").Attr("href"); ok {
		t.Errorf("Expected no value for the href attribute, got %v.", val)
	}
}

func TestRemoveAttr(t *testing.T) {
	sel := Doc2Clone().Find("div")

	sel.RemoveAttr("id")

	_, ok := sel.Attr("id")
	if ok {
		t.Error("Expected there to be no id attributes set")
	}
}

func TestSetAttr(t *testing.T) {
	sel := Doc2Clone().Find("#main")

	sel.SetAttr("id", "not-main")

	val, ok := sel.Attr("id")
	if !ok {
		t.Error("Expected an id attribute on main")
	}

	if val != "not-main" {
		t.Errorf("Expected an attribute id to be not-main, got %s", val)
	}
}

func TestSetAttr2(t *testing.T) {
	sel := Doc2Clone().Find("#main")

	sel.SetAttr("foo", "bar")

	val, ok := sel.Attr("foo")
	if !ok {
		t.Error("Expected an 'foo' attribute on main")
	}

	if val != "bar" {
		t.Errorf("Expected an attribute 'foo' to be 'bar', got '%s'", val)
	}
}

func TestText(t *testing.T) {
	txt := Doc().Find("h1").Text()
	if strings.Trim(txt, " \n\r\t") != "Provok.in" {
		t.Errorf("Expected text to be Provok.in, found %s.", txt)
	}
}

func TestText2(t *testing.T) {
	txt := Doc().Find(".hero-unit .container-fluid .row-fluid:nth-child(1)").Text()
	if ok, e := regexp.MatchString(`^\s+Provok\.in\s+Prove your point.\s+$`, txt); !ok || e != nil {
		t.Errorf("Expected text to be Provok.in Prove your point., found %s.", txt)
		if e != nil {
			t.Logf("Error: %s.", e.Error())
		}
	}
}

func TestText3(t *testing.T) {
	txt := Doc().Find(".pvk-gutter").First().Text()
	// There's an &nbsp; character in there...
	if ok, e := regexp.MatchString(`^[\s\x{00A0}]+$`, txt); !ok || e != nil {
		t.Errorf("Expected spaces, found <%v>.", txt)
		if e != nil {
			t.Logf("Error: %s.", e.Error())
		}
	}
}

func TestHtml(t *testing.T) {
	txt, e := Doc().Find("h1").Html()
	if e != nil {
		t.Errorf("Error: %s.", e)
	}

	if ok, e := regexp.MatchString(`^\s*<a href="/">Provok<span class="green">\.</span><span class="red">i</span>n</a>\s*$`, txt); !ok || e != nil {
		t.Errorf("Unexpected HTML content, found %s.", txt)
		if e != nil {
			t.Logf("Error: %s.", e.Error())
		}
	}
}

func TestNbsp(t *testing.T) {
	src := `<p>Some&nbsp;text</p>`
	d, err := NewDocumentFromReader(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	txt := d.Find("p").Text()
	ix := strings.Index(txt, "\u00a0")
	if ix != 4 {
		t.Errorf("Text: expected a non-breaking space at index 4, got %d", ix)
	}

	h, err := d.Find("p").Html()
	if err != nil {
		t.Fatal(err)
	}
	ix = strings.Index(h, "\u00a0")
	if ix != 4 {
		t.Errorf("Html: expected a non-breaking space at index 4, got %d", ix)
	}
}

func TestAddClass(t *testing.T) {
	sel := Doc2Clone().Find("#main")
	sel.AddClass("main main main")

	// Make sure that class was only added once
	if a, ok := sel.Attr("class"); !ok || a != "main" {
		t.Error("Expected #main to have class main")
	}
}

func TestAddClassSimilar(t *testing.T) {
	sel := Doc2Clone().Find("#nf5")
	sel.AddClass("odd")

	assertClass(t, sel, "odd")
	assertClass(t, sel, "odder")
	printSel(t, sel.Parent())
}

func TestAddEmptyClass(t *testing.T) {
	sel := Doc2Clone().Find("#main")
	sel.AddClass("")

	// Make sure that class was only added once
	if a, ok := sel.Attr("class"); ok {
		t.Errorf("Expected #main to not to have a class, have: %s", a)
	}
}

func TestAddClasses(t *testing.T) {
	sel := Doc2Clone().Find("#main")
	sel.AddClass("a b")

	// Make sure that class was only added once
	if !sel.HasClass("a") || !sel.HasClass("b") {
		t.Errorf("#main does not have classes")
	}
}

func TestHasClass(t *testing.T) {
	sel := Doc().Find("div")
	if !sel.HasClass("span12") {
		t.Error("Expected at least one div to have class span12.")
	}
}

func TestHasClassNone(t *testing.T) {
	sel := Doc().Find("h2")
	if sel.HasClass("toto") {
		t.Error("Expected h1 to have no class.")
	}
}

func TestHasClassNotFirst(t *testing.T) {
	sel := Doc().Find(".alert")
	if !sel.HasClass("alert-error") {
		t.Error("Expected .alert to also have class .alert-error.")
	}
}

func TestClassOperationsASCIIWhitespace(t *testing.T) {
	for _, tc := range []struct {
		name, separator string
	}{
		{"space", " "},
		{"tab", "\t"},
		{"line feed", "\n"},
		{"carriage return", "\r"},
		{"form feed", "\f"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			markup := "<div class=\"a" + tc.separator + "b\"></div>"
			for _, op := range []struct {
				name string
				run  func(*Selection)
				want string
			}{
				{"add", func(s *Selection) { s.AddClass("a") }, "a b"},
				{"remove", func(s *Selection) { s.RemoveClass("a") }, "b"},
				{"toggle", func(s *Selection) { s.ToggleClass("a") }, "b"},
			} {
				t.Run(op.name, func(t *testing.T) {
					doc, err := NewDocumentFromReader(strings.NewReader(markup))
					if err != nil {
						t.Fatal(err)
					}
					sel := doc.Find("div")
					for _, class := range []string{"a", "b"} {
						if !sel.HasClass(class) {
							t.Errorf("HasClass(%q) = false for %q", class, markup)
						}
					}
					op.run(sel)
					if got := sel.AttrOr("class", ""); got != op.want {
						t.Errorf("class = %q, want %q", got, op.want)
					}
				})
			}
		})
	}
}

func TestClassOperationsRepeatedFormFeed(t *testing.T) {
	for _, class := range []string{"\fa\fb\f", "a\f\fb", "a\fa\fb"} {
		t.Run(class, func(t *testing.T) {
			for _, toggle := range []bool{false, true} {
				doc, err := NewDocumentFromReader(strings.NewReader("<div class=\"" + class + "\"></div>"))
				if err != nil {
					t.Fatal(err)
				}
				sel := doc.Find("div")
				if toggle {
					sel.ToggleClass("a")
				} else {
					sel.RemoveClass("a")
				}
				if got := sel.AttrOr("class", ""); got != "b" {
					t.Errorf("toggle=%t: class = %q, want %q", toggle, got, "b")
				}
			}
		})
	}
}

func TestRemoveClass(t *testing.T) {
	sel := Doc2Clone().Find("#nf1")
	sel.RemoveClass("one row")

	if !sel.HasClass("even") || sel.HasClass("one") || sel.HasClass("row") {
		classes, _ := sel.Attr("class")
		t.Error("Expected #nf1 to have class even, has ", classes)
	}
}

func TestRemoveClassSimilar(t *testing.T) {
	sel := Doc2Clone().Find("#nf5, #nf6")
	assertLength(t, sel.Nodes, 2)

	sel.RemoveClass("odd")
	assertClass(t, sel.Eq(0), "odder")
	printSel(t, sel)
}

func TestRemoveAllClasses(t *testing.T) {
	sel := Doc2Clone().Find("#nf1")
	sel.RemoveClass()

	if a, ok := sel.Attr("class"); ok {
		t.Error("All classes were not removed, has ", a)
	}

	sel = Doc2Clone().Find("#main")
	sel.RemoveClass()
	if a, ok := sel.Attr("class"); ok {
		t.Error("All classes were not removed, has ", a)
	}
}

func TestToggleClass(t *testing.T) {
	sel := Doc2Clone().Find("#nf1")

	sel.ToggleClass("one")
	if sel.HasClass("one") {
		t.Error("Expected #nf1 to not have class one")
	}

	sel.ToggleClass("one")
	if !sel.HasClass("one") {
		t.Error("Expected #nf1 to have class one")
	}

	sel.ToggleClass("one even row")
	if a, ok := sel.Attr("class"); ok {
		t.Errorf("Expected #nf1 to have no classes, have %q", a)
	}
}

func TestRemoveClassRepeatedInSource(t *testing.T) {
	cases := []struct {
		class    string
		remove   string
		expected string
	}{
		{"a a", "a", ""},
		{"a a a", "a", ""},
		{"a b a", "a", "b"},
		{"a b a b", "a", "b b"},
		{"a b a b", "a b", ""},
		{"a a b", "a", "b"},
		{"b a a", "a", "b"},
	}

	for _, c := range cases {
		doc, err := NewDocumentFromReader(strings.NewReader(`<div id="t" class="` + c.class + `"></div>`))
		if err != nil {
			t.Fatal(err)
		}

		sel := doc.Find("#t")
		sel.RemoveClass(c.remove)

		got, _ := sel.Attr("class")
		if got != c.expected {
			t.Errorf("class=%q RemoveClass(%q): got class %q, want %q", c.class, c.remove, got, c.expected)
		}

		for _, removed := range strings.Fields(c.remove) {
			if sel.HasClass(removed) {
				t.Errorf("class=%q RemoveClass(%q): still has class %q", c.class, c.remove, removed)
			}
		}
	}
}

func TestToggleClassRepeatedInSource(t *testing.T) {
	doc, err := NewDocumentFromReader(strings.NewReader(`<div id="t" class="a a b"></div>`))
	if err != nil {
		t.Fatal(err)
	}

	sel := doc.Find("#t")

	sel.ToggleClass("a")
	if sel.HasClass("a") {
		t.Error("expected #t to not have class a after toggling it off")
	}
	if !sel.HasClass("b") {
		t.Error("expected #t to keep class b")
	}

	sel.ToggleClass("a")
	if !sel.HasClass("a") {
		t.Error("expected #t to have class a after toggling it back on")
	}
}
