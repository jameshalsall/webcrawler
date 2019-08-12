package model

import (
	"fmt"
	"strings"
)

type Sitemap struct {
	BaseURL string
	Children map[string]Page
}

func (s Sitemap) String() string {
	var str []string
	str = append(str, fmt.Sprintf("Sitemap for %s", s.BaseURL))

	for _, p := range s.Children {
		str = append(str, printChildren(&p, 1))
	}

	return strings.Join(str, "\n")
}

type Page struct {
	URL string
	Children map[string]Page
}

func (p Page) String() string {
	return "----" + p.URL
}

func printChildren(p *Page, indent int) string {
	str := []string{
		strings.Repeat("----", indent) + p.String(),
	}

	for _, cp := range p.Children {
		str = append(str, printChildren(&cp, indent + 1))
	}

	return strings.Join(str, "\n")
}
