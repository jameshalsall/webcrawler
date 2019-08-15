package client

import (
	"github.com/jameshalsall/webcrawler/parser"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	p := parser.NewParser()
	type args struct {
		hr parser.HrefParser
	}
	tests := []struct {
		name string
		args args
		want HttpClient
	}{
		{
			name: "Returns a new instance of a client",
			args: args{hr: p},
			want: &httpClient{hr: p},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.hr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
