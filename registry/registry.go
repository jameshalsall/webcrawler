// Package registry provides interfaces for tracking visited page URLs
// in the crawler.
package registry

import (
	"sync"
)

// PageRegistry is responsible for keeping track of which URLs
// have been visited and is safe to use across goroutines.
type PageRegistry interface {
	HasBeenVisited(url string) bool
	Visit(url string)
}

type pageRegistry struct {
	mux     sync.RWMutex
	visited map[string]bool
}

// NewRegistry creates a new default registry instance that
// is initialised and ready to use
func NewRegistry() PageRegistry {
	return &pageRegistry{
		visited: map[string]bool{},
	}
}

func (pr *pageRegistry) HasBeenVisited(url string) bool {
	pr.mux.RLock()
	defer pr.mux.RUnlock()

	return pr.visited[url]
}

func (pr *pageRegistry) Visit(url string) {
	pr.mux.Lock()
	defer pr.mux.Unlock()

	pr.visited[url] = true
}
