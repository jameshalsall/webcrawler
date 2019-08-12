package reader

import (
	"github.com/PuerkitoBio/goquery"
	"io"
)

type HrefReader interface {
	ReadFrom(reader io.Reader) ([]string, error)
}

func NewReader() HrefReader {
	return &goqueryReader{}
}

type goqueryReader struct {
}

func (hr goqueryReader) ReadFrom(reader io.Reader) ([]string, error) {
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
