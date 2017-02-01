package goquery

import (
	"fmt"
	"strconv"
	"testing"

	"golang.org/x/net/html"

	"github.com/stretchr/testify/assert"
)

const testPage = `<!DOCTYPE html>
<html>
  <head>
    <title></title>
    <meta charset="utf-8" />
  </head>
  <body>
    <h1>
      <ul id="resources">
        <li class="resource" order="3">
          <div class="name">Foo</div>
        </li>
        <li class="resource" order="1">
          <div class="name">Bar</div>
        </li>
        <li class="resource" order="4">
          <div class="name">Baz</div>
        </li>
        <li class="resource" order="2">
          <div class="name">Bang</div>
        </li>
        <li class="resource" order="5">
          <div class="name">Zip</div>
        </li>
      </ul>
			<h2 id="anchor-header"><a href="https://foo.com">FOO!!!</a></h2>
    </h1>
		<ul id="structured-list">
		  <li name="foo" val="flip">foo</li>
			<li name="bar" val="flip">bar</li>
			<li name="baz" val="flip">baz</li>
		</ul>
		<ul id="nested-map">
			<ul name="first">
				<li name="foo">foo</li>
				<li name="bar">bar</li>
				<li name="baz">baz</li>
			</ul>
			<ul name="second">
				<li name="bang">bang</li>
				<li name="ring">ring</li>
				<li name="fling">fling</li>
			</ul>
		</ul>
		<div class="foobar">
			<thing foo="yes">1</thing>
			<foo arr="true">true</foo>
			<bar arr="true">false</foo>
			<float>1.2345</float>
			<int>-123</int>
			<uint>100</uint>
		</div>
  </body>
</html>
`

type Page struct {
	Resources []Resource `goquery:"#resources .resource"`
	FooBar    FooBar
}

type Resource struct {
	Name string `goquery:".name"`
}

type Attr struct {
	Key, Value string
}

type FooBar struct {
	Attrs              []Attr
	Val                int
	unmarshalWasCalled bool
}

type AttrSelectorTest struct {
	Header H2 `goquery:"#anchor-header"`
}

type H2 struct {
	Location string `goquery:"a,[href]"`
}

type sliceAttrSelector struct {
	// For arrays/slices, type []primitive can use a source attribute
	Things []bool `goquery:".foobar [arr=\"true\"],[arr]"`
}

// TODO(andrewstuart) maps are unimplemented
type MapTest struct {
	// For map[primitive]primitive we use syntax selector,keySource,valSource
	Names map[string]string `goquery:"#structured-list li,[name],[val]"`
	// For map[primitive]Object we use the same syntax as a []primitive
	Resources map[string]Resource `goquery:"#resources .resource,[order]"`

	Nested map[string]map[string]string `goquery:"#nested-map,[name],[name],text"`
}

func (f *FooBar) UnmarshalHTML(nodes []*html.Node) error {
	f.unmarshalWasCalled = true

	s := &Selection{}
	s = s.AddNodes(nodes...)

	f.Attrs = []Attr{}
	for _, node := range s.Find(".foobar thing").Nodes {
		for _, attr := range node.Attr {
			f.Attrs = append(f.Attrs, Attr{Key: attr.Key, Value: attr.Val})
		}
	}
	thing := s.Find("thing")

	thingText := thing.Text()

	i, err := strconv.Atoi(thingText)
	f.Val = i
	return err
}

type ErrorFooBar struct{}

var errTestUnmarshal = fmt.Errorf("A wild error appeared")

func (e *ErrorFooBar) UnmarshalHTML([]*html.Node) error {
	return errTestUnmarshal
}

var vals = []string{"Foo", "Bar", "Baz", "Bang", "Zip"}

func TestDecoder(t *testing.T) {
	asrt := assert.New(t)

	asrt.Implements((*Unmarshaler)(nil), new(FooBar))

	var p Page

	asrt.NoError(Unmarshal([]byte(testPage), &p))
	asrt.Len(p.Resources, 5)

	for i, val := range vals {
		asrt.Equal(val, p.Resources[i].Name)
	}

	asrt.True(p.FooBar.unmarshalWasCalled, "Unmarshal should have been called.")
	asrt.Equal(1, p.FooBar.Val)
	asrt.Len(p.FooBar.Attrs, 1)
	asrt.Equal("foo", p.FooBar.Attrs[0].Key)
	asrt.Equal("yes", p.FooBar.Attrs[0].Value)
}

func TestArrayUnmarshal(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		Resources [5]Resource `goquery:"#resources .resource"`
	}

	asrt.NoError(Unmarshal([]byte(testPage), &a))
	for i, val := range vals {
		asrt.Equal(val, a.Resources[i].Name)
	}
}

func TestBoolean(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		BoolTest struct {
			Foo bool `goquery:"foo"`
			Bar bool `goquery:"bar"`
		} `goquery:".foobar"`
	}

	asrt.NoError(Unmarshal([]byte(testPage), &a))

	asrt.Equal(true, a.BoolTest.Foo)
	asrt.Equal(false, a.BoolTest.Bar)
}

func TestNumbers(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		BoolTest struct {
			Int   int     `goquery:"int"`
			Float float32 `goquery:"float"`
			Uint  uint16  `goquery:"uint"`
		} `goquery:".foobar"`
	}

	asrt.NoError(Unmarshal([]byte(testPage), &a))

	asrt.Equal(float32(1.2345), a.BoolTest.Float)
	asrt.Equal(-123, a.BoolTest.Int)
	asrt.Equal(uint16(100), a.BoolTest.Uint)
}

func checkErr(asrt *assert.Assertions, err error) *CannotUnmarshalError {
	asrt.Error(err)
	asrt.IsType((*CannotUnmarshalError)(nil), err)
	return err.(*CannotUnmarshalError)
}

func TestUnmarshalError(t *testing.T) {
	asrt := assert.New(t)

	var a []ErrorFooBar

	err := Unmarshal([]byte(testPage), &a)

	asrt.Contains(err.Error(), "an error occurred")

	e := checkErr(asrt, err)
	e2 := checkErr(asrt, e.Err)

	asrt.Equal(errTestUnmarshal, e2.Err)
	asrt.Equal(CustomUnmarshalError, e2.Reason)
}

func TestNilUnmarshal(t *testing.T) {
	asrt := assert.New(t)

	var a *Page

	err := Unmarshal([]byte{}, a)
	e := checkErr(asrt, err)
	asrt.Equal(NilValue, e.Reason)
}

func TestNonPointer(t *testing.T) {
	asrt := assert.New(t)

	var a Page
	err := Unmarshal([]byte{}, a)
	e := checkErr(asrt, err)
	asrt.Equal(NonPointer, e.Reason)
}

func TestWrongArrayLength(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		Resources [1]Resource `goquery:".resource"`
	}

	err := Unmarshal([]byte(testPage), &a)

	e := checkErr(asrt, err)
	asrt.Equal(TypeConversionError, e.Reason)
	e2 := checkErr(asrt, e.Err)
	asrt.Equal(ArrayLengthMismatch, e2.Reason)

	asrt.Contains(e.Error(), "Resource")
	asrt.Contains(e.Error(), "array length")
}

func TestInvalidLiteral(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		Foo int `goquery:"foo"`
	}

	err := Unmarshal([]byte(testPage), &a)

	e := checkErr(asrt, err).unwindReason()

	asrt.Len(e.chain, 2)
	asrt.Error(e.tail)
	asrt.Contains(err.Error(), e.tail.Error())

	asrt.Equal(TypeConversionError, e.chain[0].Reason)
	asrt.Equal(TypeConversionError, e.chain[1].Reason)
}

func TestInvalidArrayEleType(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		Resources [5]int `goquery:".resource"`
	}

	err := Unmarshal([]byte(testPage), &a)
	e := checkErr(asrt, err).unwindReason()
	asrt.Len(e.chain, 3)
}

func TestAttributeSelector(t *testing.T) {
	asrt := assert.New(t)

	var a AttrSelectorTest

	asrt.NoError(Unmarshal([]byte(testPage), &a))
	asrt.Equal("https://foo.com", a.Header.Location)
}

func TestSliceAttrSelector(t *testing.T) {
	asrt := assert.New(t)

	var a sliceAttrSelector

	asrt.NoError(Unmarshal([]byte(testPage), &a))
	asrt.Len(a.Things, 2)
	asrt.True(a.Things[0])
	asrt.True(a.Things[1])
}

func TestMapQuery(t *testing.T) {
	asrt := assert.New(t)

	a := MapTest{}

	asrt.NoError(Unmarshal([]byte(testPage), &a))
	asrt.Len(a.Names, 3)
	asrt.Equal("flip", a.Names["foo"])

	asrt.Len(a.Resources, 5)
	asrt.Len(a.Nested, 2)
	asrt.Len(a.Nested["first"], 3)
	asrt.Len(a.Nested["second"], 3)
}
