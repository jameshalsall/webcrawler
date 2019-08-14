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

	sitem := &model.Sitemap{BaseURL: startingUrl, Children: map[string]model.Page{}}
	reg := registry.NewRegistry()
	pagech := make(chan model.Page, 100)
	errch := make(chan error)

	go progress(reg)
	go listenForErrors(errch)

	c := crawler.NewCrawler(reg, pagech, errch, client.NewClient(parser.NewParser()))

	go c.Crawl(startingUrl, depth)

	for p, ok := <-pagech; ok; p, ok = <-pagech {
		sitem.Children[p.URL] = p
	}

	fmt.Println(model.SitemapAsAscii(sitem))
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
