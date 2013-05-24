# goquery - a little like that j-thing, only in Go

GoQuery brings a syntax and a set of features similar to [jQuery][] to the [Go language][go]. It is based on the experimental html package and the CSS Selector library [cascadia][]. Since the experimental html parser returns tokens (nodes), and not a full-featured DOM object, jQuery's manipulation and modification functions have been left off (no point in modifying data in the parsed tree of the HTML, it has no effect).

Supported functions are query-oriented features (`hasClass()`, `attr()` and the likes), as well as traversing functions that make sense given what we have to work with. This makes GoQuery a great library for scraping web pages.

Syntax-wise, it is as close as possible to jQuery, with the same function names when possible, and that warm and fuzzy chainable interface. jQuery being the ultra-popular library that it is, I felt that writing a similar HTML-manipulating library was better to follow its API than to start anew (in the same spirit as Go's `fmt` package), even though some of its methods are less than intuitive (looking at you, [index()][index]...).

## Installation

    $ go get github.com/PuerkitoBio/goquery

(optional) To run unit tests:
    
    $ cd $GOPATH/src/github.com/PuerkitoBio/goquery
    $ go test

(optional) To run benchmarks:

    $ cd $GOPATH/src/github.com/PuerkitoBio/goquery
    $ go test -bench=".*"

## Changelog

*    **v0.3.0** : Add `EachWithBreak()` which allows to break out of an `Each()` loop by returning false. This function was added instead of changing the existing `Each()` to avoid breaking compatibility.
*    **v0.2.1** : Make go-getable, now that [go.net/html is Go1.0-compatible][gonet] (thanks to @matrixik for pointing this out).
*    **v0.2.0** : Add support for negative indices in Slice(). **BREAKING CHANGE** `Document.Root` is removed, `Document` is now a `Selection` itself (a selection of one, the root element, just like `Document.Root` was before). Add jQuery's Closest() method.
*    **v0.1.1** : Add benchmarks to use as baseline for refactorings, refactor Next...() and Prev...() methods to use the new html package's linked list features (Next/PrevSibling, FirstChild). Good performance boost (40+% in some cases).
*    **v0.1.0** : Initial release.

## API

GoQuery exposes two classes, `Document` and `Selection`. Unlike jQuery, which is loaded as part of a DOM document, and thus acts on its containing document, GoQuery doesn't know which HTML document to act upon. So it needs to be told, and that's what the `Document` class is for. It holds the root document node as the initial Selection object to manipulate.

jQuery often has many variants for the same function (no argument, a selector string argument, a jQuery object argument, a DOM element argument, ...). Instead of exposing the same features in GoQuery as a single method with variadic empty interface arguments, I use statically-typed signatures following this naming convention:

*    When the jQuery equivalent can be called with no argument, it has the same name as jQuery for the no argument signature (e.g.: `Prev()`), and the version with a selector string argument is called `XxxFiltered()` (e.g.: `PrevFiltered()`)
*    When the jQuery equivalent **requires** one argument, the same name as jQuery is used for the selector string version (e.g.: `Is()`)
*    The signatures accepting a jQuery object as argument are defined in GoQuery as `XxxSelection()` and take a `*Selection` object as argument (e.g.: `FilterSelection()`)
*    The signatures accepting a DOM element as argument in jQuery are defined in GoQuery as `XxxNodes()` and take a variadic argument of type `*html.Node` (e.g.: `FilterNodes()`)
*    Finally, the signatures accepting a function as argument in jQuery are defined in GoQuery as `XxxFunction()` and take a function as argument (e.g.: `FilterFunction()`)

GoQuery's complete [godoc reference documentation can be found here][doc].

Please note that Cascadia's selectors do NOT necessarily match all supported selectors of jQuery (Sizzle). See the [cascadia project][cascadia] for details.

## Examples

Taken from example_test.go:

```Go
import (
  "fmt"
  // In real use, this import would be required (not in this example, since it
  // is part of the goquery package)
  //"github.com/PuerkitoBio/goquery"
  "strconv"
)

// This example scrapes the 10 reviews shown on the home page of MetalReview.com,
// the best metal review site on the web :) (and no, I'm not affiliated to them!)
func ExampleScrape_MetalReview() {
  // Load the HTML document (in real use, the type would be *goquery.Document)
  var doc *Document
  var e error

  if doc, e = NewDocument("http://metalreview.com"); e != nil {
    panic(e.Error())
  }

  // Find the review items (the type of the Selection would be *goquery.Selection)
  doc.Find(".slider-row:nth-child(1) .slider-item").Each(func(i int, s *Selection) {
    var band, title string
    var score float64

    // For each item found, get the band, title and score, and print it
    band = s.Find("strong").Text()
    title = s.Find("em").Text()
    if score, e = strconv.ParseFloat(s.Find(".score").Text(), 64); e != nil {
      // Not a valid float, ignore score
      fmt.Printf("Review %d: %s - %s.\n", i, band, title)
    } else {
      // Print all, including score
      fmt.Printf("Review %d: %s - %s (%2.1f).\n", i, band, title, score)
    }
  })
  // To see the output of the Example while running the test suite (go test), simply
  // remove the leading "x" before Output on the next line. This will cause the
  // example to fail (all the "real" tests should pass).

  // xOutput: voluntarily fail the Example output.
}
```

## License

The [BSD 3-Clause license][bsd], the same as the [Go language][golic]. Cascadia's license is [here][caslic].

[jquery]: http://jquery.com/
[go]: http://golang.org/
[cascadia]: http://code.google.com/p/cascadia/
[bsd]: http://opensource.org/licenses/BSD-3-Clause
[golic]: http://golang.org/LICENSE
[caslic]: http://code.google.com/p/cascadia/source/browse/LICENSE
[doc]: http://godoc.org/github.com/PuerkitoBio/goquery
[index]: http://api.jquery.com/index/
[gonet]: http://code.google.com/p/go/source/detail?r=f7f5159120f51ba0070774d3c5907969b5fe7858&repo=net
