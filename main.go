package main

import (
	"fmt"
	"github.com/jameshalsall/crawler/crawler"
	"github.com/jameshalsall/crawler/model"
	"github.com/jameshalsall/crawler/reader"
	"github.com/jameshalsall/crawler/registry"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	startingUrl := os.Args[1]

	sitem := &model.Sitemap{BaseURL: startingUrl, Children: map[string]model.Page{}}
	reg := registry.NewRegistry()
	pagech := make(chan model.Page)
	errch := make(chan error)

	c := crawler.NewCrawler(reg, pagech, errch, reader.NewReader())

	go c.Crawl(startingUrl, 2)

	go func() {
		select {
		case err := <- errch:
			log.Fatalln(err)
		}
	}()

	for p, ok := <-pagech; ok; p, ok = <-pagech {
		sitem.Children[p.URL] = p
	}

	fmt.Println(sitem)
	fmt.Println("Done.")
}

func usage() {
	fmt.Println(`Usage:
crawler <URL>  Start crawling on the specified URL, including protocol e.g. "crawler https://monzo.com"`)
}
