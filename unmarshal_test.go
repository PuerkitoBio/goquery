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

	asrt.Contains(err.Error(), "[]goquery.ErrorFooBar[0]")

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
	e := checkErr(asrt, Unmarshal([]byte{}, a))
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

	e := checkErr(asrt, err).unwind()

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
	e := checkErr(asrt, err).unwind()
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

type MapTest struct {
	// For map[primitive]primitive we use syntax selector,keySource,valSource
	Names map[string]string `goquery:"#structured-list li,[name],[val]"`
	// For map[primitive]Object we use the same syntax as a []primitive
	Resources map[string]Resource `goquery:"#resources .resource,[order]"`

	Nested map[string]map[string]string `goquery:"#nested-map,[name],[name],text"`
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

func TestMapNonStringKey(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		Map map[int]Resource `goquery:".resource,[order]"`
	}

	asrt.NoError(Unmarshal([]byte(testPage), &a))
	asrt.Len(a.Map, 5)
	asrt.Equal(a.Map[1].Name, "Bar")
}

func TestErroringKey(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		Map map[ErrorFooBar]Resource `goquery:".resource,[order]"`
	}
	err := checkErr(asrt, Unmarshal([]byte(testPage), &a))
	asrt.Equal(errTestUnmarshal, err.unwind().tail)
}

func TestDirectInsertion(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		Nodes []*html.Node `goquery:"ul#resources .resource"`
	}

	asrt.NoError(Unmarshal([]byte(testPage), &a))
	asrt.Len(a.Nodes, 5)
}

func TestInnerHtml(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		HTML []string `goquery:"ul#resources .resource,html"`
	}

	asrt.NoError(Unmarshal([]byte(testPage), &a))
	asrt.Len(a.HTML, 5)
	asrt.Equal(a.HTML[0], `<div class="name">Foo</div>`)
}

func TestMapShortTag(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		Names map[string]string `goquery:"#structured-list li,[name]"`
	}

	asrt.NoError(Unmarshal([]byte(testPage), &a))
	asrt.Len(a.Names, 3)
	// Test that we just use inner text when missing a value selector
	asrt.Equal("foo", a.Names["foo"])
	asrt.Equal("bar", a.Names["bar"])
}

func TestNoKeySelector(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		Names map[string]string `goquery:"#structured-list li"`
	}

	err := checkErr(asrt, Unmarshal([]byte(testPage), &a))
	asrt.Equal(MissingValueSelector, err.unwind().last().Reason)
}

func TestMapInnerError(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		Names map[string]ErrorFooBar `goquery:"#structured-list li,[name]"`
	}
	err := checkErr(asrt, Unmarshal([]byte(testPage), &a))
	asrt.Equal(errTestUnmarshal, err.unwind().tail)
}

func TestInterfaceDecode(t *testing.T) {
	asrt := assert.New(t)
	var a struct {
		IF interface{} `goquery:"#structured-list li"`
	}
	asrt.NoError(Unmarshal([]byte(testPage), &a))
	asrt.Equal("foobarbaz", a.IF.(string))
}

const hnPage = `<html op="news"><head><meta name="referrer" content="origin"><meta name="viewport" content="width=device-width, initial-scale=1.0"><link rel="stylesheet" type="text/css" href="news.css?HLnf3vl4tF17hLCHxIT6">
        <link rel="shortcut icon" href="favicon.ico">
          <link rel="alternate" type="application/rss+xml" title="RSS" href="rss">
        <title>Hacker News</title></head><body><center><table id="hnmain" border="0" cellpadding="0" cellspacing="0" width="85%" bgcolor="#f6f6ef">
        <tr><td bgcolor="#ff6600"><table border="0" cellpadding="0" cellspacing="0" width="100%" style="padding:2px"><tr><td style="width:18px;padding-right:4px"><a href="http://www.ycombinator.com"><img src="y18.gif" width="18" height="18" style="border:1px white solid;"></a></td>
                  <td style="line-height:12pt; height:10px;"><span class="pagetop"><b class="hnname"><a href="news">Hacker News</a></b>
              <a href="newest">new</a> | <a href="newcomments">comments</a> | <a href="show">show</a> | <a href="ask">ask</a> | <a href="jobs">jobs</a> | <a href="submit">submit</a>            </span></td><td style="text-align:right;padding-right:4px;"><span class="pagetop">
                              <a href="login?goto=news">login</a>
                          </span></td>
              </tr></table></td></tr>
<tr style="height:10px"></tr><tr><td><table border="0" cellpadding="0" cellspacing="0" class="itemlist">
              <tr class='athing' id='13543927'>
      <td align="right" valign="top" class="title"><span class="rank">1.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13543927' href='vote?id=13543927&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="https://medium.com/airbnb-engineering/introducing-lottie-4ff4a0afac0e" class="storylink">Introducing Lottie: Airbnb's tool for adding animations to native apps</a><span class="sitebit comhead"> (<a href="from?site=medium.com"><span class="sitestr">medium.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13543927">295 points</span> by <a href="user?id=dikaiosune" class="hnuser">dikaiosune</a> <span class="age"><a href="item?id=13543927">4 hours ago</a></span> <span id="unv_13543927"></span> | <a href="hide?id=13543927&amp;goto=news">hide</a> | <a href="item?id=13543927">59&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13541679'>
      <td align="right" valign="top" class="title"><span class="rank">2.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13541679' href='vote?id=13541679&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="item?id=13541679" class="storylink">Ask HN: Who is hiring? (February 2017)</a></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13541679">438 points</span> by <a href="user?id=whoishiring" class="hnuser">whoishiring</a> <span class="age"><a href="item?id=13541679">7 hours ago</a></span> <span id="unv_13541679"></span> | <a href="hide?id=13541679&amp;goto=news">hide</a> | <a href="item?id=13541679">698&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13545330'>
      <td align="right" valign="top" class="title"><span class="rank">3.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13545330' href='vote?id=13545330&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="https://www.bloomberg.com/news/articles/2017-02-01/apple-developing-new-mac-chip-in-test-of-intel-independence" class="storylink">Apple Said to Work on Mac Chip That Would Lessen Intel Role</a><span class="sitebit comhead"> (<a href="from?site=bloomberg.com"><span class="sitestr">bloomberg.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13545330">107 points</span> by <a href="user?id=VeXocide" class="hnuser">VeXocide</a> <span class="age"><a href="item?id=13545330">2 hours ago</a></span> <span id="unv_13545330"></span> | <a href="hide?id=13545330&amp;goto=news">hide</a> | <a href="item?id=13545330">106&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13543233'>
      <td align="right" valign="top" class="title"><span class="rank">4.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13543233' href='vote?id=13543233&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="https://www.opensourcery.co.za/2017/01/05/the-jvm-is-not-that-heavy/" class="storylink">The JVM is not that heavy</a><span class="sitebit comhead"> (<a href="from?site=opensourcery.co.za"><span class="sitestr">opensourcery.co.za</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13543233">235 points</span> by <a href="user?id=khy" class="hnuser">khy</a> <span class="age"><a href="item?id=13543233">5 hours ago</a></span> <span id="unv_13543233"></span> | <a href="hide?id=13543233&amp;goto=news">hide</a> | <a href="item?id=13543233">151&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13545564'>
      <td align="right" valign="top" class="title"><span class="rank">5.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13545564' href='vote?id=13545564&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="http://petrolicious.com/apartment-find-this-ferrari-250-gt-pf-coupe-was-hidden-in-hollywood-for-decades" class="storylink">This Ferrari 250 GT PF Coupe Was Hidden in a Hollywood Apartment for Decades</a><span class="sitebit comhead"> (<a href="from?site=petrolicious.com"><span class="sitestr">petrolicious.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13545564">75 points</span> by <a href="user?id=6stringmerc" class="hnuser">6stringmerc</a> <span class="age"><a href="item?id=13545564">2 hours ago</a></span> <span id="unv_13545564"></span> | <a href="hide?id=13545564&amp;goto=news">hide</a> | <a href="item?id=13545564">26&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13544491'>
      <td align="right" valign="top" class="title"><span class="rank">6.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13544491' href='vote?id=13544491&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="https://blogs.msdn.microsoft.com/dotnet/2017/02/01/the-net-language-strategy/" class="storylink">The .NET Language Strategy</a><span class="sitebit comhead"> (<a href="from?site=microsoft.com"><span class="sitestr">microsoft.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13544491">140 points</span> by <a href="user?id=benaadams" class="hnuser">benaadams</a> <span class="age"><a href="item?id=13544491">3 hours ago</a></span> <span id="unv_13544491"></span> | <a href="hide?id=13544491&amp;goto=news">hide</a> | <a href="item?id=13544491">108&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13542587'>
      <td align="right" valign="top" class="title"><span class="rank">7.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13542587' href='vote?id=13542587&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="http://blog.2ndquadrant.com/dataloss-at-gitlab/" class="storylink">Data Loss at GitLab</a><span class="sitebit comhead"> (<a href="from?site=2ndquadrant.com"><span class="sitestr">2ndquadrant.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13542587">220 points</span> by <a href="user?id=umairshahid" class="hnuser">umairshahid</a> <span class="age"><a href="item?id=13542587">6 hours ago</a></span> <span id="unv_13542587"></span> | <a href="hide?id=13542587&amp;goto=news">hide</a> | <a href="item?id=13542587">68&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13543072'>
      <td align="right" valign="top" class="title"><span class="rank">8.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13543072' href='vote?id=13543072&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="https://play.google.com/store/apps/details?id=us.arrived.arrived" class="storylink">Show HN: Arrived – Stack Overflow for US Immigration</a><span class="sitebit comhead"> (<a href="from?site=play.google.com"><span class="sitestr">play.google.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13543072">159 points</span> by <a href="user?id=wjmclaugh" class="hnuser">wjmclaugh</a> <span class="age"><a href="item?id=13543072">6 hours ago</a></span> <span id="unv_13543072"></span> | <a href="hide?id=13543072&amp;goto=news">hide</a> | <a href="item?id=13543072">62&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13541589'>
      <td align="right" valign="top" class="title"><span class="rank">9.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13541589' href='vote?id=13541589&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="http://debrouwere.org/2017/02/01/unlearning-descriptive-statistics/" class="storylink">Unlearning descriptive statistics</a><span class="sitebit comhead"> (<a href="from?site=debrouwere.org"><span class="sitestr">debrouwere.org</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13541589">219 points</span> by <a href="user?id=stdbrouw" class="hnuser">stdbrouw</a> <span class="age"><a href="item?id=13541589">8 hours ago</a></span> <span id="unv_13541589"></span> | <a href="hide?id=13541589&amp;goto=news">hide</a> | <a href="item?id=13541589">67&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13545310'>
      <td align="right" valign="top" class="title"><span class="rank">10.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13545310' href='vote?id=13545310&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="http://www.economist.com/news/business/21716020-tech-firms-are-last-departing-their-see-no-evil-stance-society-and-politics" class="storylink">Tech firms are departing from their see-no-evil stance on society and politics</a><span class="sitebit comhead"> (<a href="from?site=economist.com"><span class="sitestr">economist.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13545310">28 points</span> by <a href="user?id=martincmartin" class="hnuser">martincmartin</a> <span class="age"><a href="item?id=13545310">2 hours ago</a></span> <span id="unv_13545310"></span> | <a href="hide?id=13545310&amp;goto=news">hide</a> | <a href="item?id=13545310">2&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13544752'>
      <td align="right" valign="top" class="title"><span class="rank">11.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13544752' href='vote?id=13544752&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="https://mgba.io/2014/12/28/classic-nes/" class="storylink">The Tactics That Certain GBA Cartridges Use to Defeat Emulation Software</a><span class="sitebit comhead"> (<a href="from?site=mgba.io"><span class="sitestr">mgba.io</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13544752">61 points</span> by <a href="user?id=kibwen" class="hnuser">kibwen</a> <span class="age"><a href="item?id=13544752">3 hours ago</a></span> <span id="unv_13544752"></span> | <a href="hide?id=13544752&amp;goto=news">hide</a> | <a href="item?id=13544752">25&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13543827'>
      <td align="right" valign="top" class="title"><span class="rank">12.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13543827' href='vote?id=13543827&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="http://blog.runkit.com/2017/02/01/stop-filing-bugs-file-a-container.html" class="storylink">Stop Filing Bugs, File a Container</a><span class="sitebit comhead"> (<a href="from?site=runkit.com"><span class="sitestr">runkit.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13543827">94 points</span> by <a href="user?id=tolmasky" class="hnuser">tolmasky</a> <span class="age"><a href="item?id=13543827">4 hours ago</a></span> <span id="unv_13543827"></span> | <a href="hide?id=13543827&amp;goto=news">hide</a> | <a href="item?id=13543827">23&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13545412'>
      <td align="right" valign="top" class="title"><span class="rank">13.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13545412' href='vote?id=13545412&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="https://github.com/FredKSchott/CoVim" class="storylink">Collaborative Editing for Vim</a><span class="sitebit comhead"> (<a href="from?site=github.com"><span class="sitestr">github.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13545412">22 points</span> by <a href="user?id=JepZ" class="hnuser">JepZ</a> <span class="age"><a href="item?id=13545412">2 hours ago</a></span> <span id="unv_13545412"></span> | <a href="hide?id=13545412&amp;goto=news">hide</a> | <a href="item?id=13545412">4&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13539552'>
      <td align="right" valign="top" class="title"><span class="rank">14.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13539552' href='vote?id=13539552&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="https://garbagecollected.org/2017/01/31/four-column-ascii/" class="storylink">Four Column ASCII</a><span class="sitebit comhead"> (<a href="from?site=garbagecollected.org"><span class="sitestr">garbagecollected.org</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13539552">432 points</span> by <a href="user?id=nishs" class="hnuser">nishs</a> <span class="age"><a href="item?id=13539552">14 hours ago</a></span> <span id="unv_13539552"></span> | <a href="hide?id=13539552&amp;goto=news">hide</a> | <a href="item?id=13539552">54&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13546231'>
      <td align="right" valign="top" class="title"><span class="rank">15.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13546231' href='vote?id=13546231&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="https://github.com/akelleh/causality" class="storylink" rel="nofollow">Causal inference in python</a><span class="sitebit comhead"> (<a href="from?site=github.com"><span class="sitestr">github.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13546231">10 points</span> by <a href="user?id=aleyan" class="hnuser">aleyan</a> <span class="age"><a href="item?id=13546231">1 hour ago</a></span> <span id="unv_13546231"></span> | <a href="hide?id=13546231&amp;goto=news">hide</a> | <a href="item?id=13546231">discuss</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13546090'>
      <td align="right" valign="top" class="title"><span class="rank">16.</span></td>      <td></td><td class="title"><a href="https://boards.greenhouse.io/bright/jobs/549319#.WJJb_bYrJbU" class="storylink" rel="nofollow">Bright (W15 – developing countries solar) is hiring its first front-end eng (SF)</a><span class="sitebit comhead"> (<a href="from?site=greenhouse.io"><span class="sitestr">greenhouse.io</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="age"><a href="item?id=13546090">1 hour ago</a></span> | <a href="hide?id=13546090&amp;goto=news">hide</a>      </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13541331'>
      <td align="right" valign="top" class="title"><span class="rank">17.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13541331' href='vote?id=13541331&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="http://www.nature.com/news/astronaut-twin-study-hints-at-stress-of-space-travel-1.21380" class="storylink">Scott Kelly's DNA shows unexpected telomere lengthening after year in space</a><span class="sitebit comhead"> (<a href="from?site=nature.com"><span class="sitestr">nature.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13541331">132 points</span> by <a href="user?id=ramzyo" class="hnuser">ramzyo</a> <span class="age"><a href="item?id=13541331">8 hours ago</a></span> <span id="unv_13541331"></span> | <a href="hide?id=13541331&amp;goto=news">hide</a> | <a href="item?id=13541331">49&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13539767'>
      <td align="right" valign="top" class="title"><span class="rank">18.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13539767' href='vote?id=13539767&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="https://medium.com/@squeaky_pl/million-requests-per-second-with-python-95c137af319#.3eiek12me" class="storylink">Million requests per second with Python</a><span class="sitebit comhead"> (<a href="from?site=medium.com"><span class="sitestr">medium.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13539767">371 points</span> by <a href="user?id=d_theorist" class="hnuser">d_theorist</a> <span class="age"><a href="item?id=13539767">13 hours ago</a></span> <span id="unv_13539767"></span> | <a href="hide?id=13539767&amp;goto=news">hide</a> | <a href="item?id=13539767">134&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13544871'>
      <td align="right" valign="top" class="title"><span class="rank">19.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13544871' href='vote?id=13544871&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="http://variety.com/2017/digital/news/facebook-oculus-zenimax-decision-1201975648/" class="storylink">Facebook Ordered to Pay $500M in Oculus Lawsuit</a><span class="sitebit comhead"> (<a href="from?site=variety.com"><span class="sitestr">variety.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13544871">135 points</span> by <a href="user?id=_pius" class="hnuser">_pius</a> <span class="age"><a href="item?id=13544871">3 hours ago</a></span> <span id="unv_13544871"></span> | <a href="hide?id=13544871&amp;goto=news">hide</a> | <a href="item?id=13544871">75&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13541038'>
      <td align="right" valign="top" class="title"><span class="rank">20.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13541038' href='vote?id=13541038&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="http://www.newyorker.com/business/currency/becoming-warren-buffett-the-man-not-the-investor?intcid=mod-latest" class="storylink">“Becoming Warren Buffett,” the Man, Not the Investor</a><span class="sitebit comhead"> (<a href="from?site=newyorker.com"><span class="sitestr">newyorker.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13541038">153 points</span> by <a href="user?id=artsandsci" class="hnuser">artsandsci</a> <span class="age"><a href="item?id=13541038">9 hours ago</a></span> <span id="unv_13541038"></span> | <a href="hide?id=13541038&amp;goto=news">hide</a> | <a href="item?id=13541038">74&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13542735'>
      <td align="right" valign="top" class="title"><span class="rank">21.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13542735' href='vote?id=13542735&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="https://motherboard.vice.com/en_us/article/big-data-cambridge-analytica-brexit-trump" class="storylink">Cambridge Analytica: The Data That Turned the World Upside Down</a><span class="sitebit comhead"> (<a href="from?site=vice.com"><span class="sitestr">vice.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13542735">99 points</span> by <a href="user?id=lakeeffect" class="hnuser">lakeeffect</a> <span class="age"><a href="item?id=13542735">6 hours ago</a></span> <span id="unv_13542735"></span> | <a href="hide?id=13542735&amp;goto=news">hide</a> | <a href="item?id=13542735">45&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13542294'>
      <td align="right" valign="top" class="title"><span class="rank">22.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13542294' href='vote?id=13542294&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="http://www.unofficialgoogledatascience.com/2017/01/causality-in-machine-learning.html" class="storylink">Causality in machine learning</a><span class="sitebit comhead"> (<a href="from?site=unofficialgoogledatascience.com"><span class="sitestr">unofficialgoogledatascience.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13542294">61 points</span> by <a href="user?id=maverick_iceman" class="hnuser">maverick_iceman</a> <span class="age"><a href="item?id=13542294">7 hours ago</a></span> <span id="unv_13542294"></span> | <a href="hide?id=13542294&amp;goto=news">hide</a> | <a href="item?id=13542294">17&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13544514'>
      <td align="right" valign="top" class="title"><span class="rank">23.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13544514' href='vote?id=13544514&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="https://dolphin-emu.org/blog/2017/02/01/dolphin-progress-report-january-2017/" class="storylink">Dolphin Progress Report: January 2017</a><span class="sitebit comhead"> (<a href="from?site=dolphin-emu.org"><span class="sitestr">dolphin-emu.org</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13544514">53 points</span> by <a href="user?id=dEnigma" class="hnuser">dEnigma</a> <span class="age"><a href="item?id=13544514">3 hours ago</a></span> <span id="unv_13544514"></span> | <a href="hide?id=13544514&amp;goto=news">hide</a> | <a href="item?id=13544514">7&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13541544'>
      <td align="right" valign="top" class="title"><span class="rank">24.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13541544' href='vote?id=13541544&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="https://techcrunch.com/2017/02/01/tesla-motors-inc-is-now-officially-tesla-inc/" class="storylink">Tesla Motors, Inc. Is Now Officially Tesla, Inc</a><span class="sitebit comhead"> (<a href="from?site=techcrunch.com"><span class="sitestr">techcrunch.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13541544">168 points</span> by <a href="user?id=pearlsteinj" class="hnuser">pearlsteinj</a> <span class="age"><a href="item?id=13541544">8 hours ago</a></span> <span id="unv_13541544"></span> | <a href="hide?id=13541544&amp;goto=news">hide</a> | <a href="item?id=13541544">72&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13539724'>
      <td align="right" valign="top" class="title"><span class="rank">25.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13539724' href='vote?id=13539724&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="http://www.bbc.com/future/story/20170131-the-secret-to-living-a-meaningful-life" class="storylink">Living a Meaningful Life</a><span class="sitebit comhead"> (<a href="from?site=bbc.com"><span class="sitestr">bbc.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13539724">110 points</span> by <a href="user?id=bootload" class="hnuser">bootload</a> <span class="age"><a href="item?id=13539724">12 hours ago</a></span> <span id="unv_13539724"></span> | <a href="hide?id=13539724&amp;goto=news">hide</a> | <a href="item?id=13539724">24&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13540214'>
      <td align="right" valign="top" class="title"><span class="rank">26.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13540214' href='vote?id=13540214&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="http://introtodeeplearning.com/index.html" class="storylink">6.S191: Introduction to Deep Learning</a><span class="sitebit comhead"> (<a href="from?site=introtodeeplearning.com"><span class="sitestr">introtodeeplearning.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13540214">150 points</span> by <a href="user?id=seycombi" class="hnuser">seycombi</a> <span class="age"><a href="item?id=13540214">11 hours ago</a></span> <span id="unv_13540214"></span> | <a href="hide?id=13540214&amp;goto=news">hide</a> | <a href="item?id=13540214">12&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13541162'>
      <td align="right" valign="top" class="title"><span class="rank">27.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13541162' href='vote?id=13541162&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="https://andychase.me/mail/gigster-contract/" class="storylink">An Email Thread Between a Developer and Gigster</a><span class="sitebit comhead"> (<a href="from?site=andychase.me"><span class="sitestr">andychase.me</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13541162">667 points</span> by <a href="user?id=mfts0" class="hnuser">mfts0</a> <span class="age"><a href="item?id=13541162">9 hours ago</a></span> <span id="unv_13541162"></span> | <a href="hide?id=13541162&amp;goto=news">hide</a> | <a href="item?id=13541162">230&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13546201'>
      <td align="right" valign="top" class="title"><span class="rank">28.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13546201' href='vote?id=13546201&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="item?id=13546201" class="storylink">Ask HN: When and how should I release my open source tool?</a></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13546201">20 points</span> by <a href="user?id=flaque" class="hnuser">flaque</a> <span class="age"><a href="item?id=13546201">1 hour ago</a></span> <span id="unv_13546201"></span> | <a href="hide?id=13546201&amp;goto=news">hide</a> | <a href="item?id=13546201">7&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13542428'>
      <td align="right" valign="top" class="title"><span class="rank">29.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13542428' href='vote?id=13542428&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="http://www.fakenewschallenge.org/" class="storylink">Fake News Challenge</a><span class="sitebit comhead"> (<a href="from?site=fakenewschallenge.org"><span class="sitestr">fakenewschallenge.org</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13542428">153 points</span> by <a href="user?id=phreeza" class="hnuser">phreeza</a> <span class="age"><a href="item?id=13542428">6 hours ago</a></span> <span id="unv_13542428"></span> | <a href="hide?id=13542428&amp;goto=news">hide</a> | <a href="item?id=13542428">150&nbsp;comments</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
                <tr class='athing' id='13546354'>
      <td align="right" valign="top" class="title"><span class="rank">30.</span></td>      <td valign="top" class="votelinks"><center><a id='up_13546354' href='vote?id=13546354&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="http://crypto.stackexchange.com/questions/26336/sha512-faster-than-sha256" class="storylink">SHA-512 is 1.5x faster than SHA-256 on 64-bit platforms</a><span class="sitebit comhead"> (<a href="from?site=stackexchange.com"><span class="sitestr">stackexchange.com</span></a>)</span></td></tr><tr><td colspan="2"></td><td class="subtext">
        <span class="score" id="score_13546354">7 points</span> by <a href="user?id=steffenweber" class="hnuser">steffenweber</a> <span class="age"><a href="item?id=13546354">51 minutes ago</a></span> <span id="unv_13546354"></span> | <a href="hide?id=13546354&amp;goto=news">hide</a> | <a href="item?id=13546354">discuss</a>              </td></tr>
      <tr class="spacer" style="height:5px"></tr>
            <tr class="morespace" style="height:10px"></tr><tr><td colspan="2"></td><td class="title"><a href="news?p=2" class="morelink" rel="nofollow">More</a></td></tr>
  </table>
</td></tr>
<tr><td><img src="s.gif" height="10" width="0"><table width="100%" cellspacing="0" cellpadding="1"><tr><td bgcolor="#ff6600"></td></tr></table><br><center><span class="yclinks"><a href="newsguidelines.html">Guidelines</a>
        | <a href="newsfaq.html">FAQ</a>
        | <a href="mailto:hn@ycombinator.com">Support</a>
        | <a href="https://github.com/HackerNews/API">API</a>
        | <a href="security.html">Security</a>
        | <a href="lists">Lists</a>
        | <a href="bookmarklet.html">Bookmarklet</a>
        | <a href="dmca.html">DMCA</a>
        | <a href="http://www.ycombinator.com/apply/">Apply to YC</a>
        | <a href="mailto:hn@ycombinator.com">Contact</a></span><br><br><form method="get" action="//hn.algolia.com/">Search:
          <input type="text" name="q" value="" size="17" autocorrect="off" spellcheck="false" autocapitalize="off" autocomplete="false"></form>
            </center></td></tr>      </table></center></body><script type='text/javascript' src='hn.js?HLnf3vl4tF17hLCHxIT6'></script></html>`

type page struct {
	Items map[int]*item `goquery:".itemlist,[id]"`
}

type item struct {
	Link string `goquery:".title a,[href]"`
	Site string `goquery:".title .sitebit,text"`
}

func TestHNPage(t *testing.T) {
	asrt := assert.New(t)

	var p page

	asrt.NoError(Unmarshal([]byte(hnPage), &p))
	asrt.Len(p.Items, 30)
	asrt.NotNil(p.Items[13546354])
}
