package goquery

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
