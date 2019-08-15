# Simple web crawler

This is a simple implementation of a web crawler in Go. It makes the following assumptions
when crawling a URL:

* The same URL will not be visited twice
* Child pages are mapped to their parents purely on a first-crawled basis, so if a page is
linked to by more than one "parent" page, it will just be displayed beneath the first-crawled
parent rather than both of them.
* It has a hard-coded "depth" of 3, which means it will follow 2 child pages from the starting
URL provided when running the crawler (this can be edited in the [`main.go`](main.go) file if
necessary)

## Building

To build the binary you will need Go version >= 1.11 as module support is required.

    go build

## Usage

In order to crawl you need to specify the starting URL. This should contain the protocol as well
as the domain name. For example:

    ./webcrawler https://domain.com/

Alternatively you can run the `main.go` file using `go run` without building:

    go run main.go https://domain.com/

## Viewing results

Results are written as a simple string representation of a sitemap, and written directly to stdout. This
means that viewing results for larger sitemap trees can be difficult. A simple solution is to pipe the result
to something like `more` or `less` (pardon the pun), or written to a file

    ./webcrawler https://domain.com/ | more

    ./webcrawler https://domain.com/ > crawl.out
