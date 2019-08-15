package main

import (
	"fmt"
	"github.com/jameshalsall/crawler/client"
	"github.com/jameshalsall/crawler/crawler"
	"github.com/jameshalsall/crawler/model"
	"github.com/jameshalsall/crawler/parser"
	"github.com/jameshalsall/crawler/registry"
	"log"
	"os"
	"time"
)

const depth = 4

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	startingUrl := os.Args[1]

	fmt.Println("Starting crawler on", startingUrl, "up to", depth, "pages deep")

	reg := registry.NewRegistry()
	errch := make(chan error)

	go progress(reg)
	go listenForErrors(errch)

	c := crawler.NewCrawler(reg, errch, client.NewClient(parser.NewParser()))

	sitem := c.Crawl(startingUrl, depth)

	fmt.Println("")
	fmt.Println(model.SitemapAsString(sitem))
	fmt.Println("Done.")
}

func listenForErrors(errch chan error) {
	select {
	case err := <- errch:
		log.Fatalln(err)
	}
}

func usage() {
	fmt.Println(`Usage:
crawler <URL>  Start crawling on the specified URL, including protocol e.g. "crawler https://monzo.com"`)
}

func progress(reg registry.PageRegistry) {
	fmt.Println("")
	for {
		fmt.Println("Pages crawled:", reg.NumberOfPagesVisited())
		time.Sleep(time.Second * 1)
	}
}
