package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	url             string
	concurrency     int
	graphvizOutfile string
	csvOutfile      string
	maxPages        int
	debug           bool
)

func init() {
	flag.StringVar(&url, "u", "", "Target URL")
	flag.IntVar(&concurrency, "c", 32, "Maximum number of concurrent requests")
	flag.StringVar(&graphvizOutfile, "graphviz-outfile", "", "Output graphviz file")
	flag.StringVar(&csvOutfile, "csv-outfile", "", "Output CSV file")
	flag.IntVar(&maxPages, "p", 0, "Maximum number of pages that should be fetched (0 to disable it)")
	flag.BoolVar(&debug, "d", false, "Enable debug logs")
}

func main() {
	flag.Parse()

	if url == "" {
		fmt.Fprintln(os.Stderr, "No URL given")
		os.Exit(1)
	}

	if graphvizOutfile == "" && csvOutfile == "" {
		fmt.Fprintln(os.Stderr, "Please specify at least one of the -graphviz-outfile or -csv-outfile options")
		os.Exit(1)
	}

	if concurrency < 0 {
		fmt.Fprintln(os.Stderr, "-concurrency cannot be negative")
		os.Exit(1)
	}

	crawler := NewCrawler(url, concurrency, maxPages, debug)
	links := crawler.crawl()

	if graphvizOutfile != "" {
		OutputGraphviz(url, links, graphvizOutfile)
	}

	if csvOutfile != "" {
		OutputCSV(url, links, csvOutfile)
	}
}
