package model

import (
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

func SitemapAsString(s *Sitemap) string {
	var str []string

	for _, p := range s.Children {
		str = append(str, childrenAsString(&p, 1))
	}

	return strings.Join(str, "\n")
}

func childrenAsString(p *Page, indent int) string {
	str := []string{
		strings.Repeat("----", indent) + p.URL,
	}

	for _, cp := range p.Children {
		str = append(str, childrenAsString(&cp, indent + 1))
	}

	return strings.Join(str, "\n")
}
