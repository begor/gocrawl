package main

import (
	"fmt"
	"log"
	"sync"

	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type SyncSet struct {
	state map[string]bool
	lock sync.RWMutex
}

func NewSyncSet() SyncSet {
	return SyncSet{state: make(map[string]bool), lock: sync.RWMutex{}}
}

func (s *SyncSet) Add(el string) {
	s.lock.Lock()
	s.state[el] = true
	s.lock.Unlock()
}

func (s *SyncSet) Contains(el string) bool {
	s.lock.RLock()
	_, ok := s.state[el]
	s.lock.RUnlock()
	return ok
}


func FetchDocument(url string) *goquery.Document {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal("Unable to fetch ", url)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal("Unable to get document from ", url)
	}

	return doc
}


func GetLinks(base_url_str string, doc *goquery.Document) ([]string) {
	links := make([]string, 0)

	for _, node := range doc.Find("a").Nodes {
		for _, attr := range node.Attr {
			if attr.Key != "href" {
				continue
			}

			href := attr.Val

			u, err := url.Parse(href)

			if err != nil {
				log.Fatal("Unable to parse URL ", href)
			}


			// NOTE: don't care about absolute URLs, only consider relative to current site
			if u.IsAbs() {
				continue
			}

			base_u, _ := url.Parse(base_url_str)
			new_url := url.URL{Scheme: base_u.Scheme, Host: base_u.Host, Path: u.Path}
			links = append(links, new_url.String())
		}
	}

	return links;
}


func Crawl(url_str string, visited *SyncSet, done_channel chan bool) {
	fmt.Println("Crawling: ", url_str)
	
	if visited.Contains(url_str) {
		fmt.Println(url_str, " already visited")
		done_channel <- true
		return
	}

	visited.Add(url_str)
	
	doc := FetchDocument(url_str)
	urls := GetLinks(url_str, doc)

	fmt.Println("Got ", len(urls), " links")

	child_done_channel := make(chan bool)

	for _, url := range urls {
		go Crawl(url, visited, child_done_channel)
	}

	for _, url := range urls {
		fmt.Println("Awaiting: ", url)
		<-child_done_channel
	}


	done_channel <- true

	fmt.Println("Crawled ", url_str, " !")
}


func main() {
	url := "https://monzo.com"
	visited := NewSyncSet()
	done_channel := make(chan bool)	
	
	Crawl(url, &visited, done_channel)

	<-done_channel
}