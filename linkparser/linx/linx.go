package linx

import (
	"bytes"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

type (
	// Link contains the link text and the url the link references.
	Link struct{ Href, Text string }

	// Links contains multiple links
	Links []Link

	parseFunc func(*html.Node)
)

func (l *Link) setHref(node *html.Node) {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			l.Href = attr.Val
			return
		}
	}
}

func (l *Link) setText(nodes []html.Node) {
	for _, node := range nodes {
		if node.Type == html.TextNode {
			l.Text += node.Data
		} else if node.Data == "strong" {
			l.Text += node.FirstChild.Data
		}
	}
	l.Text = removeUnwantedChars(l.Text)
}

// GetLinks returns all non-link-nested links for a html file path.
func GetLinks(data []byte) (Links, error) {
	var (
		link       Link
		links      Links
		linkParser parseFunc
	)

	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return links, err
	}

	linkParser = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var nestedNodes []html.Node

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Data == "a" {
					break
				}
				nestedNodes = append(nestedNodes, *c)
			}
			link.setHref(n)
			link.setText(nestedNodes)
			links = append(links, link)
			link = Link{}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			linkParser(c)
		}
	}
	linkParser(doc)
	return links, nil
}

func removeUnwantedChars(s string) string {
	for _, char := range s {
		if unicode.IsSpace(char) {
			s = strings.Trim(strings.ReplaceAll(s, string(char), " "), " ")
		}
	}
	return s
}
