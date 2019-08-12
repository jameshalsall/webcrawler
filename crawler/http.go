package crawler

import (
	"github.com/jameshalsall/crawler/href"
	"github.com/jameshalsall/crawler/model"
	"github.com/jameshalsall/crawler/reader"
	"github.com/jameshalsall/crawler/registry"
	"net/http"
)

type crawler struct {
	reg    registry.PageRegistry
	ch     chan<- model.Page
	errch  chan<- error
	hr reader.HrefReader
}

func NewCrawler(reg registry.PageRegistry, ch chan<- model.Page, errch chan<- error, hr reader.HrefReader) Crawler {
	return &crawler{
		reg:   reg,
		ch:    ch,
		errch: errch,
		hr: hr,
	}
}

func (c *crawler) Crawl(url string, depth int) {
	c.crawlUrl(url, c.ch, depth)
}

func (c *crawler) crawlUrl(url string, ch chan<- model.Page, depth int) {
	defer close(ch)
	if depth <= 0 {
		return
	}

	res, err := http.Get(url)
	if err != nil {
		c.errch <- err
		return
	}

	if res.StatusCode != 200 {
		return
	}

	defer res.Body.Close()

	hrefs, err := c.hr.ReadFrom(res.Body)
	if err != nil {
		c.errch <- err
		return
	}

	c.reg.Visit(url)

	for _, h := range hrefs {
		u, err := href.Normalize(url, h)
		if err != nil {
			continue
		}

		if href.UrlsHaveDifferentDomains(url, u) || c.reg.HasBeenVisited(u) {
			continue
		}

		p := &model.Page{
			URL:      u,
			Children: map[string]model.Page{},
		}

		cpch := make(chan model.Page)
		go c.crawlUrl(u, cpch, depth-1)
		for cp, ok := <-cpch; ok; cp, ok = <-cpch {
			p.Children[cp.URL] = cp
		}

		ch <- *p
	}
}
