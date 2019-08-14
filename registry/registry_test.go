package registry

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	tests := []struct {
		name string
		want PageRegistry
	}{
		{
			name: "NewRegistry() returns a registry",
			want: &pageRegistry{visited: map[string]bool{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRegistry(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRegistry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pageRegistry_HasBeenVisited(t *testing.T) {
	type fields struct {
		mux     sync.RWMutex
		visited map[string]bool
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "HasBeenVisited() returns false for URL if not visited",
			fields: fields{
				mux:     sync.RWMutex{},
				visited: map[string]bool{},
			},
			args: args{url:"https://foo.com/"},
			want: false,
		},
		{
			name: "HasBeenVisited() returns true for URL if it has been visited",
			fields: fields{
				mux:     sync.RWMutex{},
				visited: map[string]bool{"https://baz.com/": true,},
			},
			args: args{url:"https://baz.com/"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := pageRegistry{
				mux:     tt.fields.mux,
				visited: tt.fields.visited,
			}
			if got := pr.HasBeenVisited(tt.args.url); got != tt.want {
				t.Errorf("HasBeenVisited() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pageRegistry_Visit(t *testing.T) {
	t.Run("Visit() marks a page as visited", func(t *testing.T) {
		pr := pageRegistry{
			mux:     sync.RWMutex{},
			visited: map[string]bool{},
		}

		url := "https://foo.com/"
		pr.Visit(url)
		if res := pr.HasBeenVisited(url); res != true {
			t.Errorf("Visit() did not mark URL as visited (got %v, wanted true)", res)
		}
	})
}

func Test_pageRegistry_NumberOfPagesVisited(t *testing.T) {
	type fields struct {
		mux     sync.RWMutex
		visited map[string]bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "NumberOfPagesVisited() returns true count",
			fields:fields{
				mux:     sync.RWMutex{},
				visited: map[string]bool{
					"https://foo.com/": true,
					"https://baz.com/": true,
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := pageRegistry{
				mux:     tt.fields.mux,
				visited: tt.fields.visited,
			}
			if got := pr.NumberOfPagesVisited(); got != tt.want {
				t.Errorf("NumberOfPagesVisited() = %v, want %v", got, tt.want)
			}
		})
	}
}
