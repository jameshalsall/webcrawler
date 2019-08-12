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
	"runtime"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	startingUrl := os.Args[1]

	fmt.Println("Starting crawler on", startingUrl)
	go progress()

	sitem := &model.Sitemap{BaseURL: startingUrl, Children: map[string]model.Page{}}
	reg := registry.NewRegistry()
	pagech := make(chan model.Page, 100)
	errch := make(chan error)

	c := crawler.NewCrawler(reg, pagech, errch, client.NewClient(parser.NewParser()))

	go c.Crawl(startingUrl, 3)

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

func progress() {
	fmt.Println("")
	for {
		fmt.Print(".")
		fmt.Print(runtime.NumGoroutine())
		time.Sleep(time.Second * 1)
	}
}
