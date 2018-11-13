package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	neturl "net/url"
	"os"
)

// findLinks takes an HTML document and extracts all the URLs present
// in its anchor nodes.
func findLinks(doc *html.Node) []string {
	var links []string

	nodes := []*html.Node{doc}
	for len(nodes) > 0 {
		var node *html.Node
		node, nodes = nodes[0], nodes[1:]

		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			nodes = append(nodes, child)
		}
	}

	return links
}

// filterLinks takes a list of URLs and a base URL, and returns all the
// URLs that have the same domain as the base URL.
// Moreover, all relative URLs are converted to absolute ones.
func filterLinks(links []string, url *neturl.URL) []string {
	var filteredLinks []string

	for _, link := range links {
		linkUrl, err := url.Parse(link)
		if err != nil {
			continue
		}

		if linkUrl.Scheme != "" && linkUrl.Scheme != "http" && linkUrl.Scheme != "https" {
			continue
		}

		if linkUrl.Host != "" && url.Host != linkUrl.Host {
			continue
		}

		fullLink := url.ResolveReference(linkUrl)
		fullLink.Fragment = ""
		filteredLinks = append(filteredLinks, fullLink.String())
	}

	return filteredLinks
}

// GetLinks takes a URL and returns all the URLs for the same domain
// referenced in the anchor nodes of the HTML document
func GetLinks(url string) []string {
	parsedUrl, err := neturl.Parse(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot parse URL:", err)
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot fetch URL:", err)
		return nil
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot parse HTML document:", err)
		return nil
	}

	links := findLinks(doc)
	links = filterLinks(links, parsedUrl)

	return links
}
