# goquery - a little like that j-thing, only in Go

GoQuery brings a syntax and features similar to [jQuery][] to the [Go language][go]. It is based on the [experimental html package][exphtml] and the CSS Selector library [cascadia][]. Since the html parser returns tokens (nodes), and not a full-featured DOM object, jQuery's manipulation and modification functions have been left off (no point in modifying data in the parsed tree of the HTML, it has no effect).

Supported functions are (will be) query-oriented features (`hasClass()`, `attr()` and the likes), as well as traversing functions that make sense given what we have to work with.

Syntax-wise, it is as close as possible to jQuery, with the same function names when possible, and that warm and fuzzy chainable interface.

## Installation

Since this library (and cascadia) depends on the experimental branch, this package must be installed first. Both GoQuery and Cascadia expect to find the experimental library with the `"exp/html"` import statement. To install it at this location, please [follow this guide][wikiexp].

Once this is done, install GoQuery:

`go get github.com/PuerkitoBio/goquery`

## API

Coming soon.

## License

Coming soon.

[jquery]: http://jquery.com/
[go]: http://golang.org/
[exphtml]: http://code.google.com/p/go/source/browse#hg%2Fsrc%2Fpkg%2Fexp
[cascadia]: http://code.google.com/p/cascadia/
[wikiexp]: http://code.google.com/p/go-wiki/wiki/InstallingExp
