package model

import (
	"strings"
	"testing"
)

func TestSitemapAsString(t *testing.T) {
	type args struct {
		s *Sitemap
	}
	tests := []struct {
		name          string
		args          args
		shouldContain []string
	}{
		{
			name: "Returns sitemap represented as a string",
			args:args{s:&Sitemap{
				BaseURL: "https://foo.com",
				Children: map[string]Page{
					"https://baz.com": {
						URL: "https://baz.com",
						Children: map[string]Page{
							"https://bar.com": {
								URL: "https://bar.com",
								Children: map[string]Page{},
							},
						},
					},
				},
			}},
			shouldContain: []string{"https://foo.com","https://baz.com","https://bar.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SitemapAsString(tt.args.s);
			for _, s := range tt.shouldContain {
				if !strings.Contains(got, s) {
					t.Errorf("SitemapAsString() value does not contain expected URL: %s (value received: %s)", s, got)
				}
			}
		})
	}
}
