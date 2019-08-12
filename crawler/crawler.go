package crawler

type Crawler interface {
	Crawl(url string, depth int)
}

