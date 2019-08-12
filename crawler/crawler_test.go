package crawler

import (
	"github.com/jameshalsall/crawler/model"
	"io"
	"testing"
)

func TestNewCrawler(t *testing.T) {
	t.Run("NewCrawler creates a crawler", func(t *testing.T) {
		NewCrawler(&fakeReg{}, make(chan model.Page), make(chan error), &fakeReader{})
	})
}

type fakeReg struct {
}

func (f fakeReg) HasBeenVisited(url string) bool {
	return true
}

func (f fakeReg) Visit(url string) {
}

type fakeReader struct {
}

func (f fakeReader) ParseFrom(reader io.Reader) ([]string, error) {
	return nil, nil
}