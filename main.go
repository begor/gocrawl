package main

import (
	"flag"
	"fmt"
	"gocrawl/crawler"
	"gocrawl/sitemap"
)

func main() {
	parallel_fetchers := flag.Int("n", 32, "number of parallel coroutines")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("gocrawl is a simple Web Crawler.")
		fmt.Println("Usage:\n\tgocrawl -n <number of parallel coroutines> [site url]")
		return
	}

	url := flag.Args()[0]

	sitemap := sitemap.Create()
	fetch_limiter := make(chan int, *parallel_fetchers)

	crawler.Crawl(url, &sitemap, fetch_limiter)

	sitemap.Print()
}
