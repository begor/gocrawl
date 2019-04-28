package crawler

import (
	"fmt"
	"gocrawl/parser"
	"gocrawl/sitemap"
	"log"
	"net/http"
	"net/url"
)

// Main crawler routine:
// 0) if given url_str already in sitemap, return
// 1) fetch document from url_str if fetch_limiter is available or block 
// 2) parse links from document using parser package
// 3) run DFS on links using crawlChilderLinks
func Crawl(url_str string, sitemap *sitemap.Sitemap, fetch_limiter chan int) {
	fmt.Println("Crawling: ", url_str)

	current_url, err := url.Parse(url_str)

	if err != nil {
		log.Fatal("Unable to parse ", url_str, " to URL")
	}

	if sitemap.Contains(url_str) {
		fmt.Println(url_str, " already visited")
		return
	}

	sitemap.Add(url_str)

	response := fetch(url_str, fetch_limiter)
	links := parser.ExtractLinksWithCurrentHost(current_url, response.Body)

	crawlChilderLinks(links, sitemap, fetch_limiter)
}


func fetch(url string, fetch_limiter chan int) *http.Response {
	fetch_limiter <- 1

	res, err := http.Get(url)

	<-fetch_limiter

	if err != nil {
		log.Fatal("Unable to fetch ", url)
	}

	return res
}

func crawlChilderLinks(links []string, sitemap *sitemap.Sitemap, fetch_limiter chan int) {
	children_queue := make(chan bool)

	for _, link := range links {
		go func(link string) {
			Crawl(link, sitemap, fetch_limiter)
			children_queue <- true
		}(link)
	}

	for range links {
		<-children_queue
	}
}
