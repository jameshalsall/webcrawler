// href provides helper functions for verifying and manipulation
// href attribute values from a crawled page
package href

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func Normalize(baseUrl, path string) (string, error) {
	baseUrl = strings.TrimRight(baseUrl, "/")

	if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "http") {
		return "", errors.New("invalid href")
	}

	var u string
	if strings.HasPrefix(path, "http") {
		u = path
	} else {
		u = baseUrl + "/" + strings.TrimLeft(path, "/")
	}

	_, err := url.ParseRequestURI(u)
	if err != nil {
		return "", errors.New(fmt.Sprintf("error: %s is not a valid URL, skipping\n", u))
	}

	return u, nil
}

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
