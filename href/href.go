// Package href provides helper functions for verifying and manipulation
// href attribute values from a crawled page
package href

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Normalize takes a URL path (href) value and returns it in a normalized format
// ready for use as a fully-qualifed URL, including protocol
func Normalize(baseUrl, path string) (string, error) {
	baseUrl = strings.TrimRight(baseUrl, "/")

	if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "http") {
		return "", errors.New("invalid href")
	}

	var u string
	if strings.HasPrefix(path, "http") {
		u = path
	} else {
		u = baseUrl + "/" + strings.TrimLeft(strings.TrimRight(path, "/"), "/")
	}

	_, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("error: %s is not a valid URL, skipping\n", u)
	}

	return u, nil
}

// UrlsHaveDifferentDomains checks whether the 2 provided URLs share the same domain or
// are on different domains.
// NOTE: It is not intelligent enough to realise that sub-domains are on the same higher
// level domain.
func UrlsHaveDifferentDomains(url1, url2 string) bool {
	u1, err := url.Parse(url1)
	if err != nil {
		return false
	}

	u2, err := url.Parse(url2)
	if err != nil {
		return false
	}

	return u1.Host != u2.Host
}
