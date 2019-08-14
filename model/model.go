package model

import (
	"fmt"
	"strings"
)

type Sitemap struct {
	BaseURL string
	Children map[string]Page
}

type Page struct {
	URL string
	Children map[string]Page
}

func SitemapAsAscii(s *Sitemap) string {
	var str []string
	str = append(str, fmt.Sprintf("Sitemap for %s", s.BaseURL))

	for _, p := range s.Children {
		str = append(str, childrenAsAscii(&p, 1))
	}

	return strings.Join(str, "\n")
}

func childrenAsAscii(p *Page, indent int) string {
	str := []string{
		strings.Repeat("----", indent) + p.URL,
	}

	for _, cp := range p.Children {
		str = append(str, childrenAsAscii(&cp, indent + 1))
	}

	return strings.Join(str, "\n")
}
