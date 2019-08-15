// Package model provides structs that model sitemaps and their child pages
package model

import (
	"strings"
)

// Sitemap represents the root page that the crawler was started on
type Sitemap struct {
	BaseURL  string
	Children map[string]Page
}

// Page represents child pages in the sitemap
type Page struct {
	URL      string
	Children map[string]Page
}

// SitemapAsString returns a string representation of a Sitemap
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
		str = append(str, childrenAsString(&cp, indent+1))
	}

	return strings.Join(str, "\n")
}
