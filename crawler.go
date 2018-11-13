package main

import (
	"fmt"
	"sync"
)

// crawler defines a crawler object.
type crawler struct {
	baseUrl       string            // initial URL from which crawling should start
	concurrency   int               // maximum number of workers
	workqueue     chan string       // channel used to feed jobs to workers
	workerResults chan workerResult // channel used to collect workers' responses
	waitgroup     sync.WaitGroup    // waitgroup to keep track of number of in-progress requests

	crawledPages    int  // number of crawled pages
	maxCrawledPages int  // max number of pages that should be crawled
	debugLogs       bool // enable debug logs
}

// workerResult defines the response format used by a worker to return
// a slice of links for a given URL.
type workerResult struct {
	url   string   // requested URL
	links []string // links for the given URL
}

// NewCrawler creates a new crawler.
func NewCrawler(baseUrl string, concurrency int, maxCrawledPages int, debugLogs bool) *crawler {
	return &crawler{
		baseUrl:         baseUrl,
		concurrency:     concurrency,
		workqueue:       make(chan string, concurrency),
		workerResults:   make(chan workerResult, concurrency),
		maxCrawledPages: maxCrawledPages,
		debugLogs:       debugLogs,
	}
}

// spawnWorkers spawns 'concurrency' goroutines that will accept jobs
// over the 'workqueue' channel to crawl a given URL.
// The response is sent back to the collector by using the
// `workerResults' channel.
func (c *crawler) spawnWorkers() {
	for i := 0; i < c.concurrency; i++ {
		workerId := i
		go func() {
			for url := range c.workqueue {
				if c.debugLogs {
					fmt.Println("[+] worker", workerId, "fetching", url)
				}

				links := GetLinks(url)
				c.workerResults <- workerResult{url, links}
			}
		}()
	}
}

// enqueueJob adds a URL to the queue of URLs that should be crawled and
// increments the jobs waitgroup.
func (c *crawler) enqueueJob(url string) {
	c.crawledPages++
	c.waitgroup.Add(1)
	go func() { c.workqueue <- url }()
}

// collectResults collects responses from the workers and uses the
// responses to enqueue new jobs for links that have not been visited.
// It returns a map that maps a URL to a slice of links referenced by
// its HTML document.
func (c *crawler) collectResults() map[string][]string {
	visited := make(map[string][]string)

	go func() {
		for result := range c.workerResults {
			visited[result.url] = result.links

			for _, link := range result.links {
				if _, ok := visited[link]; !ok {
					// Mark the link as visited with an empty slice.
					// This will get overwritten once the worker returns its
					// response.
					visited[link] = []string{}

					if c.maxCrawledPages == 0 || c.crawledPages < c.maxCrawledPages {
						c.enqueueJob(link)
					}
				}
			}

			c.waitgroup.Done()
		}
	}()

	c.waitgroup.Wait()

	return visited
}

// crawl starts the crawler.
func (c *crawler) crawl() map[string][]string {
	c.spawnWorkers()
	c.enqueueJob(c.baseUrl)
	links := c.collectResults()

	return links
}
