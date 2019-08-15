package registry

import (
	"reflect"
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
				visited: map[string]bool{},
			},
			args: args{url: "https://foo.com/"},
			want: false,
		},
		{
			name: "HasBeenVisited() returns true for URL if it has been visited",
			fields: fields{
				visited: map[string]bool{"https://baz.com/": true},
			},
			args: args{url: "https://baz.com/"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := pageRegistry{
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
			visited: map[string]bool{},
		}

		url := "https://foo.com/"
		pr.Visit(url)
		if res := pr.HasBeenVisited(url); res != true {
			t.Errorf("Visit() did not mark URL as visited (got %v, wanted true)", res)
		}
	})
}
