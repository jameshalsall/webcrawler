package crawler

import (
	"github.com/jameshalsall/crawler/client"
	"github.com/jameshalsall/crawler/model"
	"github.com/jameshalsall/crawler/registry"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewCrawler(t *testing.T) {
	ch := make(chan model.Page)
	errch := make(chan error)

	type args struct {
		reg   registry.PageRegistry
		ch    chan<- model.Page
		errch chan<- error
		cl    client.HttpClient
	}
	tests := []struct {
		name string
		args args
		want Crawler
	}{
		{
			name: "NewCrawler() returns new instance of a crawler",
			args: args{
				reg:   registry.NewRegistry(),
				ch:    ch,
				errch: errch,
				cl:    &fakeClient{},
			},
			want: &crawler{
				reg:   registry.NewRegistry(),
				ch:    ch,
				errch: errch,
				cl:    &fakeClient{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCrawler(tt.args.reg, tt.args.ch, tt.args.errch, tt.args.cl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCrawler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_crawler_Crawl(t *testing.T) {
	errch := make(chan error)
	errstore := &errorStore{Errors: []error{}}

	go func() {
		select {
		case err := <-errch:
			errstore.add(err)
		}
	}()

	type args struct {
		url   string
		depth int
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Sitemap
		wantErr bool
	}{
		{
			name: "Crawl() sends correct pages to the page channel",
			args: args{
				url:   "https://foo.com/bar",
				depth: 3,
			},
			want: &model.Sitemap{
				BaseURL: "",
				Children: map[string]model.Page{
					"https://foo.com/bar/baz": {
						URL: "https://foo.com/bar/baz",
						Children: map[string]model.Page{
							"https://foo.com/foo/baz": {
								URL:      "https://foo.com/foo/baz",
								Children: map[string]model.Page{},
							},
							"https://foo.com/foo/bop": {
								URL:      "https://foo.com/foo/bop",
								Children: map[string]model.Page{},
							},
						},
					},
					"https://foo.com/bar/bop": {
						URL: "https://foo.com/bar/bop",
						Children: map[string]model.Page{
							"https://foo.com/foo/beep": {
								URL:      "https://foo.com/foo/beep",
								Children: map[string]model.Page{},
							},
						},
					},
				},
			},
		},
		{
			name: "Crawl() at depth 0 results in no pages",
			args: args{
				url:   "https://foo.com/",
				depth: 0,
			},
			want: &model.Sitemap{BaseURL:"", Children: map[string]model.Page{}},
		},
		{
			name: "Crawl() sends error from client fetching URL to the errch",
			args:args{
				url:   "https://non-existent-url/",
				depth: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		errstore.reset()
		ch := make(chan model.Page, 100)

		t.Run(tt.name, func(t *testing.T) {
			c := &crawler{
				reg:   registry.NewRegistry(),
				ch:    ch,
				errch: errch,
				cl:    &fakeClient{},
			}

			got := &model.Sitemap{Children: map[string]model.Page{}}
			c.Crawl(tt.args.url, tt.args.depth)
			for p, ok := <-ch; ok; p, ok = <-ch {
				got.Children[p.URL] = p
			}

			time.Sleep(time.Millisecond * 100)

			if tt.wantErr {
				if len(errstore.Errors) == 0 {
					t.Errorf("expected at least one error but received none in errch")
				}
			} else {
				if len(errstore.Errors) > 0 {
					t.Errorf("expected no errors but received following errors in errch: %v", errstore.Errors)
				} else if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("crawler.Crawl: Expected %v but got %v", tt.want, got)
				}
			}
		})
	}
	close(errch)
}

type errorStore struct {
	sync.Mutex
	Errors []error
}

func (e *errorStore) reset() {
	e.Errors = []error{}
}

func (e *errorStore) add(err error) {
	e.Lock()
	defer e.Unlock()
	e.Errors = append(e.Errors, err)
}
