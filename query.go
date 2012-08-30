package goquery

import (
	"exp/html"
	"regexp"
	"strings"
)

// Returns true if at least one node in the selection has the given class
func (this *Selection) HasClass(class string) bool {
	var rx = regexp.MustCompile("[\t\r\n]")

	class = " " + class + " "
	for _, n := range this.Nodes {
		// Applies only to element nodes
		if n.Type == html.ElementNode {
			if elClass, ok := getAttributeValue("class", n); ok {
				elClass = rx.ReplaceAllString(" "+elClass+" ", " ")
				if strings.Index(elClass, class) > -1 {
					return true
				}
			}
		}
	}
	return false
}
