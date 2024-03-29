package href

import "testing"

func TestNormalize(t *testing.T) {
	type args struct {
		baseUrl string
		path    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid URL",
			args: args{
				baseUrl: "https://foo.com/",
				path:    "/foo",
			},
			want: "https://foo.com/foo",
		},
		{
			name: "Valid URL with query string",
			args: args{
				baseUrl: "https://foo.com/",
				path:    "/foo/bar?baz=1",
			},
			want: "https://foo.com/foo/bar?baz=1",
		},
		{
			name: "Valid absolute URL in path",
			args: args{
				baseUrl: "https://foo.com/",
				path:    "https://foo.com/lending",
			},
			want: "https://foo.com/lending",
		},
		{
			name: "Trailing slash in resulting URL is removed",
			args: args{
				baseUrl: "https://foo.com/",
				path:    "/lending/",
			},
			want: "https://foo.com/lending",
		},
		{
			name: "Invalid URL: javascript: prefix",
			args: args{
				baseUrl: "https://foo.com/",
				path:    "javascript:void(0);",
			},
			wantErr: true,
		},
		{
			name: "Invalid URL: # link",
			args: args{
				baseUrl: "https://foo.com/",
				path:    "#",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Normalize(tt.args.baseUrl, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Normalize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Normalize() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlsHaveDifferentDomains(t *testing.T) {
	type args struct {
		url1 string
		url2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Returns false if URLs are for the same domain",
			args: args{
				url1: "https://foo.com/",
				url2: "https://foo.com/path/foo",
			},
			want: false,
		},
		{
			name: "Returns true if URLs are for a different domain",
			args: args{
				url1: "https://foo.com/foo",
				url2: "https://google.com/",
			},
			want: true,
		},
		{
			name: "Returns false if URLs are on different sub-domain",
			args: args{
				url1: "https://foo.com/foo",
				url2: "https://sub.foo.com/foo",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UrlsHaveDifferentDomains(tt.args.url1, tt.args.url2); got != tt.want {
				t.Errorf("UrlsHaveDifferentDomains() = %v, want %v", got, tt.want)
			}
		})
	}
}
