package goquery

import (
	"fmt"
	"log"

	// In real use, this import would be required (not in this example, since it
	// is part of the goquery package)
	//"github.com/PuerkitoBio/goquery"
)

// This example scrapes the reviews shown on the home page of metalsucks.net.
func ExampleScrape_MetalSucks() {
	// Load the HTML document (in real use, the type would be *goquery.Document)
	doc, err := NewDocument("http://metalsucks.net")
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items (the type of the Selection would be *goquery.Selection)
	doc.Find(".reviews-wrap article .review-rhs").Each(func(i int, s *Selection) {
		// For each item found, get the band and title
		band := s.Find("h3").Text()
		title := s.Find("i").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
	// To see the output of the Example while running the test suite (go test), simply
	// remove the leading "x" before Output on the next line. This will cause the
	// example to fail (all the "real" tests should pass).

	// xOutput: voluntarily fail the Example output.
}
