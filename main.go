package main

import (
	"gocrawl/crawler"
	"gocrawl/sitemap"
)

func main() {
	url := "https://monzo.com"
	sitemap := sitemap.Create()

	concurrency_limiter := make(chan int, 1)

	crawler.Crawl(url, &sitemap, concurrency_limiter)

	sitemap.Print()
}
