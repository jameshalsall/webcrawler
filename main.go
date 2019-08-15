package main

import (
	"fmt"
	"github.com/jameshalsall/webcrawler/client"
	"github.com/jameshalsall/webcrawler/crawler"
	"github.com/jameshalsall/webcrawler/model"
	"github.com/jameshalsall/webcrawler/parser"
	"github.com/jameshalsall/webcrawler/registry"
	"log"
	"os"
	"time"
)

const depth = 3

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	startingUrl := os.Args[1]

	fmt.Println("Running crawler on", startingUrl, "up to", depth, "pages deep")

	errch := make(chan error)

	go progress()
	go listenForErrors(errch)

	c := crawler.NewCrawler(registry.NewRegistry(), errch, client.NewClient(parser.NewParser()))

	sitem := c.Crawl(startingUrl, depth)

	fmt.Println("")
	fmt.Println(model.SitemapAsString(sitem))
	fmt.Println("Done.")
}

func listenForErrors(errch chan error) {
	select {
	case err := <-errch:
		log.Fatalln(err)
	}
}

func usage() {
	fmt.Println(`Usage:
webcrawler <URL>  Start crawling on the specified URL, including protocol e.g. "webcrawler https://monzo.com"`)
}

func progress() {
	fmt.Println("")
	for {
		fmt.Print(".")
		time.Sleep(time.Second * 1)
	}
}
