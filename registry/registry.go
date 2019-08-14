package registry

import (
	"sync"
)

type PageRegistry interface {
	HasBeenVisited(url string) bool
	Visit(url string)
	NumberOfPagesVisited() int
}

type pageRegistry struct {
	mux sync.RWMutex
	visited map[string]bool
}

func NewRegistry() PageRegistry {
	return &pageRegistry{
		visited: map[string]bool{},
	}
}

func (pr pageRegistry) HasBeenVisited(url string) bool {
	pr.mux.RLock()
	defer pr.mux.RUnlock()

	return pr.visited[url]
}

func (pr pageRegistry) Visit(url string) {
	pr.mux.Lock()
	defer pr.mux.Unlock()

	pr.visited[url] = true
}

func (pr pageRegistry) NumberOfPagesVisited() int {
	pr.mux.RLock()
	defer pr.mux.RUnlock()

	return len(pr.visited)
}
