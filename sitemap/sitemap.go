package sitemap

import (
	"fmt"
	"sync"
)

// Sitemap represented as a thread-safe HashSet.
// All syncronization happens under the hood using RWMutex to keep sync of shared mutable state inside black box.
type Sitemap struct {
	urls   map[string][]string
	errors map[string]error
	lock   sync.RWMutex
}

func Create() Sitemap {
	return Sitemap{urls: make(map[string][]string), errors: make(map[string]error), lock: sync.RWMutex{}}
}

func (s *Sitemap) MarkVisited(url string) {
	s.lock.Lock()

	if _, ok := s.urls[url]; !ok {
		s.urls[url] = make([]string, 0)
	}

	s.lock.Unlock()
}

func (s *Sitemap) SetLinks(url string, links []string) {
	s.lock.Lock()

	s.urls[url] = links

	s.lock.Unlock()
}

func (s *Sitemap) GetLinks(url string) ([]string, bool) {
	s.lock.RLock()
	val, ok := s.urls[url]
	s.lock.RUnlock()

	return val, ok
}

func (s *Sitemap) SetError(url string, err error) {
	s.lock.Lock()

	delete(s.urls, url)
	s.errors[url] = err

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

func (s *Sitemap) ErrSize() int {
	s.lock.RLock()
	size := len(s.errors)
	s.lock.RUnlock()

	return size
}

func (s *Sitemap) PrintReport(short_view bool) {
	s.lock.RLock()

	fmt.Println("Succesfully crawled ", s.Size(), " pages")

	for url, links := range s.urls {
		fmt.Println(url)

		if short_view {
			continue
		}

		for _, link := range links {
			fmt.Println("\t-", link)
		}
	}

	fmt.Println("\n\n\nGot ", s.ErrSize(), " errors")

	for url, error := range s.errors {
		fmt.Println(url, ": ", error)
	}

	s.lock.RUnlock()
}
