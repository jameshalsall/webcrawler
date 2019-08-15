// Package parser provides interfaces for parsing href values from readers.
package parser

import (
	"github.com/PuerkitoBio/goquery"
	"io"
)

// HrefParser provides a way of parsing href attribute values from
// a HTML document source.
type HrefParser interface {
	ParseFrom(reader io.Reader) ([]string, error)
}

// NewParser creates a new HrefParser
func NewParser() HrefParser {
	return &goqueryHrefParser{}
}

type goqueryHrefParser struct {
}

// ParseFrom will read a HTML document body from the provided `io.Reader`
// instance and return any href values in a slice of strings.
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
