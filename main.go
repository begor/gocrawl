package main

import (
	"flag"
	"fmt"
	"github.com/op/go-logging"
	"gocrawl/crawler"
	"gocrawl/fetcher"
	"gocrawl/sitemap"
	"time"
)

var log = logging.MustGetLogger("main")

func setupLogging() {
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.1s} %{color:reset} %{message}`,
	)

	logging.SetFormatter(format)
}

func main() {
	parallel_fetchers := flag.Int("n", 32, "Number of parallel coroutines.")
	short_view := flag.Bool("s", false, "Short view (without page links).")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("gocrawl is a simple Web Crawler.")
		fmt.Println("Usage:\n\tgocrawl -n <number of parallel coroutines> [site url]")
		return
	}

	url := flag.Args()[0]

	setupLogging()

	sitemap := sitemap.Create()
	fetcher := fetcher.CreateHTTPFetcher(*parallel_fetchers)

	log.Info("Starting Crawler with ", *parallel_fetchers, " concurrent fetchers")

	start := time.Now()
	crawler.Crawl(url, &fetcher, &sitemap)
	fin := time.Now()

	elapsed := fin.Sub(start)

	log.Info("Crawler done at ", elapsed, " seconds!")

	sitemap.PrintReport(*short_view)
}
