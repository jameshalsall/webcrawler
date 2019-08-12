package crawler

import (
	"github.com/jameshalsall/crawler/model"
	"testing"
)

func TestNewCrawler(t *testing.T) {
	t.Run("NewCrawler creates a crawler", func(t *testing.T) {
		NewCrawler(&fakeReg{}, make(chan model.Page), make(chan error), &fakeClient{})
	})
}

type fakeReg struct {
}

func (fr fakeReg) HasBeenVisited(url string) bool {
	return true
}

func (fr fakeReg) Visit(url string) {
}

type fakeClient struct {
}

func (fc fakeClient) GetHrefsFromUrl(url string) ([]string, error) {
	return []string{""}, nil
}
