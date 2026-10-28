package goquery

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/html"
)

// SelectorForNode returns a CSS selector for node, using the same strategy as
// Chrome DevTools' "Copy selector" command. It returns false for non-element
// nodes.
//
// The selector is optimized by stopping at an ID or at the html, head, or body
// element. Otherwise, each step uses the element name, classes when they
// distinguish same-name siblings, and :nth-child when classes are insufficient.
// IDs are assumed to be unique, as they are by Chrome DevTools. The algorithm
// mirrors the implementation at:
// https://github.com/ChromeDevTools/devtools-frontend/blob/main/front_end/panels/elements/DOMPath.ts
func SelectorForNode(node *html.Node) (string, bool) {
	if node == nil || node.Type != html.ElementNode {
		return "", false
	}

	var steps []string
	for current := node; current != nil; current = current.Parent {
		if current.Type == html.DocumentNode {
			break
		}

		step, optimized, ok := selectorStep(current, current == node)
		if !ok {
			return "", false
		}
		steps = append(steps, step)
		if optimized {
			break
		}
	}

	slices.Reverse(steps)
	return strings.Join(steps, " > "), true
}

func selectorStep(node *html.Node, isTarget bool) (value string, optimized, ok bool) {
	if node == nil || node.Type != html.ElementNode {
		return "", false, false
	}

	id := nodeAttribute(node, "id")
	if id != "" {
		return "#" + escapeCSSIdentifier(id), true, true
	}

	tag := node.Data
	switch strings.ToLower(tag) {
	case "html", "head", "body":
		return tag, true, true
	}

	parent := node.Parent
	if parent == nil || parent.Type == html.DocumentNode {
		return tag, true, true
	}

	classes := strings.Fields(nodeAttribute(node, "class"))
	needsClasses := false
	needsNthChild := false
	elementIndex := 0
	ownIndex := -1

	for sibling := parent.FirstChild; sibling != nil; sibling = sibling.NextSibling {
		if sibling.Type != html.ElementNode {
			continue
		}

		if sibling == node {
			ownIndex = elementIndex
			elementIndex++
			continue
		}
		elementIndex++

		if sibling.Data != tag {
			continue
		}

		needsClasses = true
		if len(classes) == 0 || nodeHasAllClasses(sibling, classes) {
			needsNthChild = true
		}
	}

	if ownIndex < 0 {
		return "", false, false
	}

	var result strings.Builder
	result.WriteString(tag)
	if isTarget && strings.EqualFold(tag, "input") && len(classes) == 0 {
		if inputType := nodeAttribute(node, "type"); inputType != "" {
			result.WriteString("[type=")
			result.WriteString(escapeCSSIdentifier(inputType))
			result.WriteByte(']')
		}
	}

	if needsNthChild {
		fmt.Fprintf(&result, ":nth-child(%d)", ownIndex+1)
	} else if needsClasses {
		for _, class := range classes {
			result.WriteByte('.')
			result.WriteString(escapeCSSIdentifier(class))
		}
	}

	return result.String(), false, true
}

func nodeAttribute(node *html.Node, name string) string {
	if attr := getAttributePtr(name, node); attr != nil {
		return attr.Val
	}
	return ""
}

func nodeHasAllClasses(node *html.Node, classes []string) bool {
	siblingClasses := make(map[string]struct{})
	for _, class := range strings.Fields(nodeAttribute(node, "class")) {
		siblingClasses[class] = struct{}{}
	}
	for _, class := range classes {
		if _, ok := siblingClasses[class]; !ok {
			return false
		}
	}
	return true
}

// escapeCSSIdentifier implements the CSSOM "serialize an identifier"
// algorithm used by CSS.escape.
func escapeCSSIdentifier(value string) string {
	if value == "" {
		return ""
	}

	runes := []rune(value)
	var result strings.Builder
	for i, r := range runes {
		switch {
		case r == 0:
			result.WriteRune(utf8.RuneError)
		case r >= 0x1 && r <= 0x1f || r == 0x7f:
			writeCSSCodePointEscape(&result, r)
		case i == 0 && r >= '0' && r <= '9':
			writeCSSCodePointEscape(&result, r)
		case i == 1 && r >= '0' && r <= '9' && runes[0] == '-':
			writeCSSCodePointEscape(&result, r)
		case i == 0 && r == '-' && len(runes) == 1:
			result.WriteString(`\-`)
		case r >= utf8.RuneSelf || r == '-' || r == '_' ||
			r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z':
			result.WriteRune(r)
		default:
			result.WriteByte('\\')
			result.WriteRune(r)
		}
	}
	return result.String()
}

func writeCSSCodePointEscape(result *strings.Builder, r rune) {
	result.WriteByte('\\')
	result.WriteString(strconv.FormatInt(int64(r), 16))
	result.WriteByte(' ')
}
