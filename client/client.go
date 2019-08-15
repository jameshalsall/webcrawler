// Package client provides a HTTP client for the crawler.
package client

import (
	"github.com/jameshalsall/webcrawler/parser"
	"net/http"
)

// HttpClient provides a way of fetching hrefs (links) from URLs
type HttpClient interface {
	GetHrefsFromUrl(url string) ([]string, error)
}

type httpClient struct {
	hr parser.HrefParser
}

// NewClient creates and returns a default HttpClient instance
// with the provided HrefParser, ready for use.
func NewClient(hr parser.HrefParser) HttpClient {
	return &httpClient{
		hr: hr,
	}
}

// GetHrefsFromUrl takes a url parameter, fetches it using a HTTP request
// in the Go stdlib, parses the response for href attributes in <a> tags
// using the client's configured parser.HrefParser instance and then returns
// those hrefs as a slice of strings.
func (h httpClient) GetHrefsFromUrl(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return []string{}, nil
	}

	defer res.Body.Close()

	return h.hr.ParseFrom(res.Body)
}
