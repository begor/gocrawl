package main

import (
	"fmt"
	"log"

	"net/http"
	"net/url"

	"gocrawl/sitemap"

	"github.com/PuerkitoBio/goquery"
)



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
			links = append(links, base_u.ResolveReference(u).String())
		}
	}

	return links;
}


func Crawl(url_str string, sitemap *sitemap.Sitemap, concurrency_limiter chan int) {
	fmt.Println("Crawling: ", url_str)
	fmt.Println("Visited so far: ", sitemap.Size())
	
	if sitemap.Contains(url_str) {
		fmt.Println(url_str, " already visited")
		return
	}

	sitemap.Add(url_str)

	concurrency_limiter <- 1
	
	doc := FetchDocument(url_str)

	<-concurrency_limiter

	urls := GetLinks(url_str, doc)

	fmt.Println("Got ", len(urls), " links from ", url_str)

	children_queue := make(chan bool)

	for _, url := range urls {
		go func (u string) {
			Crawl(u, sitemap, concurrency_limiter)
			children_queue <- true
		}(url)
	}

	for range urls {
		<-children_queue
	}


	fmt.Println("Crawled ", url_str, "!")
}


func main() {
	url := "https://monzo.com"
	sitemap := sitemap.Create()
	concurrency_limiter := make(chan int, 1000)
	
	Crawl(url, &sitemap, concurrency_limiter)

	sitemap.Print()
}