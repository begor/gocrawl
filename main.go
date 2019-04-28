package main

import (
	"flag"
	"fmt"
	"gocrawl/fetcher"
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
	fetcher := fetcher.CreateHTTPFetcher(*parallel_fetchers)

	crawler.Crawl(url, &fetcher, &sitemap)

	sitemap.Print()
}
