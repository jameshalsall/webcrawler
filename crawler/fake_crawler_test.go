package crawler

import (
	"errors"
)

type fakeClient struct {
}

func (fc fakeClient) GetHrefsFromUrl(url string) ([]string, error) {
	hrefsForUrls := map[string][]string{
		"https://foo.com/bar": {
			"https://foo.com/bar/baz",
			"https://foo.com/bar/bop",
		},
		"https://foo.com/bar/baz": {
			"https://foo.com/foo/baz",
			"https://foo.com/foo/bop",
		},
		"https://foo.com/bar/bop": {
			"https://foo.com/foo/baz",
			"https://foo.com/foo/beep",
			"invalid",
		},
		"https://foo.com/foo/baz":  {},
		"https://foo.com/foo/bop":  {},
		"https://foo.com/foo/beep": {},
	}

	hrefs, ok := hrefsForUrls[url]
	if !ok {
		return nil, errors.New("client cannot fetch URL")
	}

	return hrefs, nil
}
