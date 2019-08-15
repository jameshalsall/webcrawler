package crawler

import (
	"github.com/jameshalsall/webcrawler/client"
	"github.com/jameshalsall/webcrawler/href"
	"github.com/jameshalsall/webcrawler/model"
	"github.com/jameshalsall/webcrawler/registry"
	"sync"
)

type Crawler interface {
	Crawl(url string, depth int) *model.Sitemap
}

type crawler struct {
	reg   registry.PageRegistry
	errch chan<- error
	cl    client.HttpClient
}

func NewCrawler(reg registry.PageRegistry, errch chan<- error, cl client.HttpClient) Crawler {
	return &crawler{
		reg:   reg,
		errch: errch,
		cl:    cl,
	}
}

func (c *crawler) Crawl(url string, depth int) *model.Sitemap {
	ch := make(chan model.Page, 100)
	sitem := &model.Sitemap{BaseURL: url, Children: map[string]model.Page{}}

	go c.crawlUrl(url, ch, depth)

	for p, ok := <-ch; ok; p, ok = <-ch {
		sitem.Children[p.URL] = p
	}

	return sitem
}

func (c *crawler) crawlUrl(url string, ch chan<- model.Page, depth int) {
	defer close(ch)

	if depth <= 0 {
		return
	}

	hrefs, err := c.cl.GetHrefsFromUrl(url)
	if err != nil {
		c.errch <- err
		return
	}

	c.reg.Visit(url)

	wg := &sync.WaitGroup{}
	for _, h := range hrefs {
		u, err := href.Normalize(url, h)
		if err != nil {
			continue
		}

		if href.UrlsHaveDifferentDomains(url, u) || c.reg.HasBeenVisited(u) {
			continue
		}

		cpch := make(chan model.Page)
		wg.Add(1)
		go func() {
			p := &model.Page{
				URL:      u,
				Children: map[string]model.Page{},
			}

			for cp, ok := <-cpch; ok; cp, ok = <-cpch {
				p.Children[cp.URL] = cp
			}
			ch <- *p
			wg.Done()
		}()

		go c.crawlUrl(u, cpch, depth-1)

		wg.Wait()
	}
}
