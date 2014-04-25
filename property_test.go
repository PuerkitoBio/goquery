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

func TestAttrNotExist(t *testing.T) {
	if val, ok := Doc().Find("div.row-fluid").Attr("href"); ok {
		t.Errorf("Expected no value for the href attribute, got %v.", val)
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
