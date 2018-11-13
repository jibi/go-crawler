package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

// OutputGraphviz takes as input the base URL, a map of all the links
// crawled and generates a file ("outfile") in graphviz format.
//
// Example output file:
// digraph {
//	node [shape = circle]; "https://example.com";
//	node [shape = box];
//	"https://example.com/test" -> "https://example.com/test2"
//	"https://example.com/test" -> "https://example.com/test3"
//      ..
// }
func OutputGraphviz(url string, links map[string][]string, outfile string) {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("digraph {\n"))
	buffer.WriteString(fmt.Sprintf("\tnode [shape = circle]; \"%s\";\n", url))
	buffer.WriteString(fmt.Sprintf("\tnode [shape = box];\n"))

	for url, links := range links {
		for _, link := range links {
			buffer.WriteString(fmt.Sprintf("\t\"%s\" -> \"%s\"\n", url, link))
		}
	}

	buffer.WriteString("}")

	ioutil.WriteFile(outfile, buffer.Bytes(), 0644)
}

// OutputCSV takes as input the base URL, a map of all the links
// crawled and generates a file ("outfile") in CSV format for each link,
// where the first column is the document URL and the second column is a
// link referenced in the document
func OutputCSV(url string, links map[string][]string, outfile string) {
	var buffer bytes.Buffer

	for url, links := range links {
		for _, link := range links {
			buffer.WriteString(fmt.Sprintf("%s;%s\n", url, link))
		}
	}

	ioutil.WriteFile(outfile, buffer.Bytes(), 0644)
}
