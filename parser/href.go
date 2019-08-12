package parser

import (
	"github.com/PuerkitoBio/goquery"
	"io"
)

type HrefParser interface {
	ParseFrom(reader io.Reader) ([]string, error)
}

func NewParser() HrefParser {
	return &goqueryHrefParser{}
}

type goqueryHrefParser struct {
}

func (hp goqueryHrefParser) ParseFrom(reader io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return nil, err
	}

	var hrefs []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		h, ok := s.Attr("href")
		if !ok {
			return
		}

		hrefs = append(hrefs, h)
	})

	return hrefs, nil
}
