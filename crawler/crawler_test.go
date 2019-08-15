package crawler

import (
	"github.com/jameshalsall/webcrawler/client"
	"github.com/jameshalsall/webcrawler/model"
	"github.com/jameshalsall/webcrawler/registry"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewCrawler(t *testing.T) {
	errch := make(chan error)

	type args struct {
		reg   registry.PageRegistry
		errch chan<- error
		cl    client.HttpClient
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Returns new instance of a crawler",
			args: args{
				reg:   registry.NewRegistry(),
				errch: errch,
				cl:    &fakeClient{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCrawler(tt.args.reg, tt.args.errch, tt.args.cl)
			_, ok := got.(Crawler)
			if !ok {
				t.Error("NewCrawler() does not return an instance of a crawler")
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
			name: "Crawls and builds an expected, valid sitemap",
			args: args{
				url:   "https://foo.com/bar",
				depth: 3,
			},
			want: &model.Sitemap{
				BaseURL: "https://foo.com/bar",
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
			name: "With a depth of 0 it builds an empty sitemap",
			args: args{
				url:   "https://foo.com/",
				depth: 0,
			},
			want: &model.Sitemap{BaseURL: "https://foo.com/", Children: map[string]model.Page{}},
		},
		{
			name: "Sends error from client fetching URL to the errch",
			args: args{
				url:   "https://non-existent-url/",
				depth: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		errstore.reset()

		t.Run(tt.name, func(t *testing.T) {
			c := &crawler{
				reg:   registry.NewRegistry(),
				errch: errch,
				cl:    &fakeClient{},
			}

			got := c.Crawl(tt.args.url, tt.args.depth)

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
