package parser

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestNewParser(t *testing.T) {
	tests := []struct {
		name string
		want HrefParser
	}{
		{
			name: "Parser is returned",
			want: &goqueryHrefParser{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewParser(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_HrefParser_ParseFrom(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Hrefs are parsed from a parser",
			args: args{
				reader: strings.NewReader(`
				<!doctype html>
				<html>
				<head></head>
					<a href="/foo">Foo</a>	
					<a href="/foo/baz">Baz</a>
				</html>`,
				)},
			want: []string{"/foo", "/foo/baz"},
		},
		{
			name: "<a> tags with no hrefs are not returned",
			args: args{
				reader: strings.NewReader(`
				<!doctype html>
				<html>
				<head></head>
					<a href="/foo">Foo</a>	
					<a>Baz</a>
				</html>`,
				)},
			want: []string{"/foo"},
		},
		{
			name: "Empty string returns no hrefs",
			args: args{
				reader: strings.NewReader(""),
			},
		},
		{
			name: "Invalid HTML document returns no hrefs",
			args: args{
				reader: strings.NewReader(">admkiwajdiaw989243493dakdm`s,FOO></>"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hp := goqueryHrefParser{}
			got, err := hp.ParseFrom(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFrom() got = %v, want %v", got, tt.want)
			}
		})
	}
}
