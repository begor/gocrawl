package sitemap

import (
	"fmt"
	"sync"
)

// Sitemap represented as two thread-safe HashMaps: 
// one for url -> [children_url] mapping, one for url -> error mapping for those url that can't be fetched.
// All syncronization happens under the hood (using RWMutex), clients can use it as black box without synchronization.
type Sitemap struct {
	urls   map[string][]string
	errors map[string]error
	lock   sync.RWMutex
}

func Create() Sitemap {
	return Sitemap{urls: make(map[string][]string), errors: make(map[string]error), lock: sync.RWMutex{}}
}

// Mark url as visited: useful for exactly-once fetching before we sure about actual links/error fetched from that url.
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
	// If we've got error we should clear it's links for safety
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


// Very simple printer for sitemap.
// NOTE: in "real production" that'll live in separate package with interface and various implementations. 
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
