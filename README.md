Crawler
=============

A simple concurrent web crawler for a single domain that produces a CSV or a grapviz .dot file as output.

## Build
Dependencies are handled using the Go 1.11 Modules tool.

Running `go build` should be enough to build the project.

## Usage

```
$ ./crawler -help
Usage of ./crawler:
  -c int
    	Maximum number of concurrent requests (default 32)
  -csv-outfile string
    	Output CSV file
  -d	Enable debug logs
  -graphviz-outfile string
    	Output graphviz file
  -p int
    	Maximum number of pages that should be fetched (0 to disable it)
  -u string
    	Target URL
```

Once the graphviz outfile has been generated, it's possible to render it as a direct graph using the graphviz `dot` tool:

```
$ dot graph.dot -T png -o graph.png
```
