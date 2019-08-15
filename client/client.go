package client

import (
	"github.com/jameshalsall/webcrawler/parser"
	"net/http"
)

type HttpClient interface {
	GetHrefsFromUrl(url string) ([]string, error)
}

type httpClient struct {
	hr parser.HrefParser
}

func NewClient(hr parser.HrefParser) HttpClient {
	return &httpClient{
		hr: hr,
	}
}

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
