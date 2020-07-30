package goquery_test

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// This example shows how to use NewDocumentFromReader from a file.
func ExampleNewDocumentFromReader_file() {
	// create from a file
	f, err := os.Open("some/file.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}
	// use the goquery document...
	_ = doc.Find("h1")
}

// This example shows how to use NewDocumentFromReader from a string.
func ExampleNewDocumentFromReader_string() {
	// create from a string
	data := `
<html>
	<head>
		<title>My document</title>
	</head>
	<body>
		<h1>Header</h1>
	</body>
</html>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	header := doc.Find("h1").Text()
	fmt.Println(header)

	// Output: Header
}
