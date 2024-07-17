package api

import (
	"golang.org/x/net/html"
)

func extractLinks(n *html.Node) []*Link {
	var folders []*Link = make([]*Link, 0)
	var f func(*html.Node, bool)
	f = func(n *html.Node, inTable bool) {
		if n.Type == html.ElementNode {
			if n.Data == "table" {
				inTable = true
			}
			if inTable && n.Data == "a" {
				var href, text string
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						href = attr.Val
						break
					}
				}
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					text = n.FirstChild.Data
				}
				if href != "" && text != "" {
					folders = append(folders, &Link{Href: href, Name: text})
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, inTable)
		}
	}
	f(n, false)
	return folders
}

func extractNewLink(n *html.Node) *Link {
	var folder *Link
	var f func(*html.Node)
	f = func(n *html.Node) {
		if folder != nil {
			return
		}
		if n.Type == html.ElementNode {
			if n.Data == "a" {
				var href, text string
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						href = attr.Val
						break
					}
				}
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					text = n.FirstChild.Data
				}
				if href != "" && text != "" {
					folder = &Link{Href: href, Name: text}
					return
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return folder
}
