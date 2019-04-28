package sitemap

import (
	"fmt"
	"sync"
)

// Sitemap represented as a thread-safe HashSet.
// All syncronization happens under the hood using RWMutex to keep sync of shared mutable state inside black box.
type Sitemap struct {
	urls map[string]bool
	lock sync.RWMutex
}

func Create() Sitemap {
	return Sitemap{urls: make(map[string]bool), lock: sync.RWMutex{}}
}

func (s *Sitemap) Add(url string) {
	s.lock.Lock()
	s.urls[url] = true
	s.lock.Unlock()
}

func (s *Sitemap) Contains(url string) bool {
	s.lock.RLock()
	_, ok := s.urls[url]
	s.lock.RUnlock()

	return ok
}

func (s *Sitemap) Size() int {
	s.lock.RLock()
	size := len(s.urls)
	s.lock.RUnlock()

	return size
}

func (s *Sitemap) Print() {
	s.lock.RLock()

	urls := make([]string, s.Size())
	for k, _ := range s.urls {
		urls = append(urls, k)
	}

	s.lock.RUnlock()

	fmt.Println(urls)
}
